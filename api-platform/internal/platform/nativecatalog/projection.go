package nativecatalog

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const productProjectionLease = time.Minute

// productAuditOutbox is the source command: it is written in the same
// merchant-database transaction as the product state transition. Platform
// review and C-end views are therefore recoverable projections, never a
// prerequisite for preserving the audit decision.
type productAuditOutbox struct {
	ID           uint64    `gorm:"column:id;primaryKey"`
	ProductID    uint64    `gorm:"column:product_id"`
	StoreID      uint64    `gorm:"column:store_id"`
	Action       string    `gorm:"column:action"`
	ReviewStatus string    `gorm:"column:review_status"`
	Reason       string    `gorm:"column:reason"`
	ReviewedBy   uint64    `gorm:"column:reviewed_by"`
	Status       string    `gorm:"column:status"`
	Attempts     uint      `gorm:"column:attempts"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (productAuditOutbox) TableName() string { return "qixi_crm_m_product_audit_outbox" }

// StartAuditOutboxDispatcher projects committed merchant audit commands into
// the platform review and C-end view. It shares the same lease/retry semantics
// as product view projection and is safe to run in multiple platform instances.
func (h *Handler) StartAuditOutboxDispatcher(ctx context.Context) {
	go func() {
		_ = h.DispatchPendingAuditOutboxes(ctx)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = h.DispatchPendingAuditOutboxes(ctx)
			}
		}
	}()
}

func (h *Handler) DispatchPendingAuditOutboxes(ctx context.Context) error {
	var rows []productAuditOutbox
	if err := h.merchantDB.WithContext(ctx).
		Where("status IN ? OR (status = ? AND updated_at < ?)", []string{"pending", "failed"}, "processing", time.Now().Add(-productProjectionLease)).
		Order("id ASC").Limit(50).Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if _, err := h.processAuditOutbox(ctx, row); err != nil {
			continue
		}
	}
	return nil
}

func (h *Handler) processAuditOutbox(ctx context.Context, row productAuditOutbox) (bool, error) {
	claimed, err := h.claimAuditOutbox(ctx, row.ID)
	if err != nil || !claimed {
		return claimed, err
	}
	if err := h.persistAuditProjections(ctx, row); err != nil {
		return true, h.markAuditOutboxFailed(ctx, row.ID, err)
	}
	return true, h.markAuditOutboxPublished(ctx, row.ID)
}

func (h *Handler) claimAuditOutbox(ctx context.Context, id uint64) (bool, error) {
	result := h.merchantDB.WithContext(ctx).Table("qixi_crm_m_product_audit_outbox").
		Where("id = ? AND (status IN ? OR (status = ? AND updated_at < ?))", id, []string{"pending", "failed"}, "processing", time.Now().Add(-productProjectionLease)).
		Updates(map[string]any{"status": "processing", "last_error": "", "attempts": gorm.Expr("attempts + 1")})
	return result.RowsAffected == 1, result.Error
}

func (h *Handler) markAuditOutboxPublished(ctx context.Context, id uint64) error {
	return h.merchantDB.WithContext(ctx).Table("qixi_crm_m_product_audit_outbox").Where("id = ? AND status = ?", id, "processing").Updates(map[string]any{"status": "published", "processed_at": time.Now(), "last_error": ""}).Error
}

func (h *Handler) markAuditOutboxFailed(ctx context.Context, id uint64, cause error) error {
	message := strings.TrimSpace(cause.Error())
	if len(message) > 500 {
		message = message[:500]
	}
	return h.merchantDB.WithContext(ctx).Table("qixi_crm_m_product_audit_outbox").Where("id = ? AND status = ?", id, "processing").Updates(map[string]any{"status": "failed", "last_error": message}).Error
}

func (h *Handler) persistAuditProjections(ctx context.Context, row productAuditOutbox) error {
	payload, err := json.Marshal(map[string]any{"product_id": row.ProductID, "action": row.Action, "source_event_id": row.ID})
	if err != nil {
		return err
	}
	sourceEventID := row.ID
	review := productReview{ProductID: row.ProductID, StoreID: row.StoreID, SourceEventID: &sourceEventID, Status: row.ReviewStatus, Reason: row.Reason, ReviewedBy: row.ReviewedBy}
	outbox := productProjectionOutbox{ProductID: row.ProductID, SourceEventID: &sourceEventID, Action: row.Action, Payload: payload, Status: "pending"}
	if err := h.adminDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "source_event_id"}}, DoNothing: true}).Create(&review).Error; err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "source_event_id"}}, DoNothing: true}).Create(&outbox).Error
	}); err != nil {
		return err
	}
	var projection productProjectionOutbox
	if err := h.adminDB.WithContext(ctx).Where("source_event_id = ?", row.ID).First(&projection).Error; err != nil {
		return err
	}
	if projection.Status == "published" {
		return nil
	}
	processed, err := h.processProjection(ctx, projection)
	if err != nil {
		return err
	}
	if !processed {
		return errors.New("商品投影命令正在处理")
	}
	return nil
}

// StartProjectionDispatcher retries durable cross-database product projections.
// A row is atomically leased before it is processed, so platform instances do
// not process the same command concurrently. A process crash is recovered once
// its one-minute lease expires.
func (h *Handler) StartProjectionDispatcher(ctx context.Context) {
	go func() {
		_ = h.DispatchPendingProjections(ctx)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = h.DispatchPendingProjections(ctx)
			}
		}
	}()
}

// DispatchPendingProjections is public for a deterministic acceptance test and
// processes at most 50 commands per run. Business writes are idempotent, but
// the lease still prevents redundant work when several platform instances run.
func (h *Handler) DispatchPendingProjections(ctx context.Context) error {
	var rows []productProjectionOutbox
	if err := h.adminDB.WithContext(ctx).
		Where("source_event_id IS NULL AND (status IN ? OR (status = ? AND updated_at < ?))", []string{"pending", "failed"}, "processing", time.Now().Add(-productProjectionLease)).
		Order("id ASC").Limit(50).Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if _, err := h.processProjection(ctx, row); err != nil {
			// Continue with later commands: one broken merchant product must not
			// starve unrelated product projections.
			continue
		}
	}
	return nil
}

func (h *Handler) processProjection(ctx context.Context, row productProjectionOutbox) (bool, error) {
	claimed, err := h.claimProjection(ctx, row.ID)
	if err != nil || !claimed {
		return claimed, err
	}
	if err := h.dispatchProjection(ctx, row.ProductID, row.Action); err != nil {
		return true, h.markProjectionFailed(ctx, row.ID, err)
	}
	return true, h.markProjectionPublished(ctx, row.ID)
}

func (h *Handler) claimProjection(ctx context.Context, id uint64) (bool, error) {
	result := h.adminDB.WithContext(ctx).Table("qixi_crm_a_product_projection_outbox").
		Where("id = ? AND (status IN ? OR (status = ? AND updated_at < ?))", id, []string{"pending", "failed"}, "processing", time.Now().Add(-productProjectionLease)).
		Updates(map[string]any{"status": "processing", "last_error": "", "attempts": gorm.Expr("attempts + 1")})
	return result.RowsAffected == 1, result.Error
}

func (h *Handler) markProjectionPublished(ctx context.Context, id uint64) error {
	return h.adminDB.WithContext(ctx).Table("qixi_crm_a_product_projection_outbox").Where("id = ? AND status = ?", id, "processing").Updates(map[string]any{"status": "published", "processed_at": time.Now(), "last_error": ""}).Error
}

func (h *Handler) markProjectionFailed(ctx context.Context, id uint64, cause error) error {
	message := strings.TrimSpace(cause.Error())
	if len(message) > 500 {
		message = message[:500]
	}
	return h.adminDB.WithContext(ctx).Table("qixi_crm_a_product_projection_outbox").Where("id = ? AND status = ?", id, "processing").Updates(map[string]any{"status": "failed", "last_error": message}).Error
}

// dispatchProjection rebuilds a C-end product projection from merchant facts.
// It is intentionally idempotent so a durable outbox command can be retried
// after a process or business-database failure.
func (h *Handler) dispatchProjection(ctx context.Context, productID uint64, action string) error {
	if action == "delete" {
		return h.businessDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("DELETE FROM qixi_crm_b_product_sku_view WHERE product_id = ?", productID).Error; err != nil {
				return err
			}
			return tx.Exec("DELETE FROM qixi_crm_b_product_view WHERE product_id = ?", productID).Error
		})
	}
	if action != "upsert" {
		return errors.New("unsupported product projection action")
	}
	var row productRow
	if err := h.merchantDB.WithContext(ctx).Table("qixi_crm_m_product AS p").Select("p.id,p.store_id,s.merchant_id,m.name AS merchant_name,s.name AS store_name,p.title,p.category_id,p.brand_name,p.svip_price_type,p.svip_price,p.status,p.version,p.created_at").Joins("JOIN qixi_crm_m_store AS s ON s.id=p.store_id").Joins("JOIN qixi_crm_m_merchant AS m ON m.id=s.merchant_id").Where("p.id=? AND p.status='on_sale'", productID).Scan(&row).Error; err != nil {
		return err
	}
	if row.ID == 0 {
		return errors.New("approved product is unavailable for projection")
	}
	var detail detailRow
	if err := h.merchantDB.WithContext(ctx).Table("qixi_crm_m_product_detail").Where("product_id=?", productID).Scan(&detail).Error; err != nil {
		return err
	}
	var skus []skuRow
	if err := h.merchantDB.WithContext(ctx).Table("qixi_crm_m_product_sku").Where("product_id=? AND status=1", productID).Order("id ASC").Find(&skus).Error; err != nil {
		return err
	}
	if len(skus) == 0 {
		return errors.New("approved product has no sellable sku")
	}
	first := skus[0]
	return h.businessDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`INSERT INTO qixi_crm_b_product_view (product_id,merchant_id,store_id,merchant_name,store_name,category_id,brand_name,title,cover_url,price,original_price,svip_price_type,svip_price,product_type,sales,stock,sale_status,version,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,0,0,?,1,?,?) ON DUPLICATE KEY UPDATE merchant_id=VALUES(merchant_id),store_id=VALUES(store_id),merchant_name=VALUES(merchant_name),store_name=VALUES(store_name),category_id=VALUES(category_id),brand_name=VALUES(brand_name),title=VALUES(title),cover_url=VALUES(cover_url),price=VALUES(price),original_price=VALUES(original_price),svip_price_type=VALUES(svip_price_type),svip_price=VALUES(svip_price),stock=VALUES(stock),sale_status=1,version=VALUES(version),updated_at=VALUES(updated_at)`, row.ID, row.MerchantID, row.StoreID, row.MerchantName, row.StoreName, row.CategoryID, row.BrandName, row.Title, detail.CoverURL, first.Price, detail.OriginalPrice, row.SVIPPriceType, row.SVIPPrice, first.Stock, row.Version, time.Now()).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM qixi_crm_b_product_sku_view WHERE product_id=?", productID).Error; err != nil {
			return err
		}
		for _, sku := range skus {
			spec := string(sku.SpecJSON)
			if strings.TrimSpace(spec) == "" {
				spec = "{}"
			}
			if err := tx.Exec("INSERT INTO qixi_crm_b_product_sku_view (merchant_sku_id,product_id,sku_key,spec_snapshot,price,stock,sale_status,version,updated_at) VALUES (?,?,?,?,?,?,1,1,?)", sku.ID, productID, strconv.FormatUint(sku.ID, 10), spec, sku.Price, sku.Stock, time.Now()).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
