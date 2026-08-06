// Package nativesettlement owns store-scoped merchant settlement applications.
// It never reads legacy qixi_m_admin_financial tables and never receives bank
// account or payment credentials from the store console.
package nativesettlement

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	merchantsettlement "github.com/crmlive/pte-live-ecrm/api-merchant/internal/event/merchantsettlement"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	errConflict             = errors.New("settlement status conflict")
	errIdempotencyKeyReused = errors.New("idempotency key belongs to another settlement")
	errNotFound             = errors.New("settlement not found")
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/settlements", h.list)
	r.GET("/settlements/:id", h.get)
	r.POST("/settlements/:id/apply", middleware.RequireStorePermission(h.db, "finance.settlement.apply"), h.apply)
}

type settlement struct {
	ID              uint64     `gorm:"column:id"`
	StoreID         uint64     `gorm:"column:store_id"`
	MerchantID      uint64     `gorm:"column:merchant_id"`
	MerchantName    string     `gorm:"column:merchant_name"`
	RegionID        *uint64    `gorm:"column:region_id"`
	PeriodStart     time.Time  `gorm:"column:period_start"`
	PeriodEnd       time.Time  `gorm:"column:period_end"`
	Amount          float64    `gorm:"column:amount"`
	Status          string     `gorm:"column:status"`
	IdempotencyKey  *string    `gorm:"column:idempotency_key"`
	ApplicationNo   *string    `gorm:"column:application_no"`
	AppliedAt       *time.Time `gorm:"column:applied_at"`
	ReviewNote      string     `gorm:"column:review_note"`
	PayoutReference *string    `gorm:"column:payout_reference"`
	PaidAt          *time.Time `gorm:"column:paid_at"`
	Version         uint64     `gorm:"column:version"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
}

type applyInput struct {
	IdempotencyKey string `json:"idempotency_key"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := paging(c)
	q := h.base(c)
	if status, ok := settlementStatus(c.Query("status")); !ok {
		response.Fail(c, http.StatusBadRequest, "结算状态错误")
		return
	} else if status != "" {
		q = q.Where("b.status = ?", status)
	}
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where("b.updated_at >= ?", t)
		}
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where("b.updated_at < ?", t.AddDate(0, 0, 1))
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	rows := make([]settlement, 0)
	if err := q.Order("b.updated_at DESC,b.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"list": views(rows), "total": total, "page": page, "limit": limit})
}

func (h *Handler) get(c *gin.Context) {
	id := id(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "结算单 ID 错误")
		return
	}
	var row settlement
	if err := h.base(c).Where("b.id = ?", id).Scan(&row).Error; err != nil {
		fail(c)
		return
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "结算单不存在")
		return
	}
	response.OK(c, view(row))
}

func (h *Handler) apply(c *gin.Context) {
	id := id(c)
	var in applyInput
	if id == 0 || c.ShouldBindJSON(&in) != nil || !validIdempotencyKey(in.IdempotencyKey) {
		response.Fail(c, http.StatusBadRequest, "结算申请参数错误")
		return
	}
	key := strings.TrimSpace(in.IdempotencyKey)
	var out settlement
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var existing settlement
		if err := h.baseTx(c, tx).Where("b.idempotency_key = ?", key).Scan(&existing).Error; err != nil {
			return err
		}
		if existing.ID != 0 {
			if existing.ID == id {
				out = existing
				return nil
			}
			return errIdempotencyKeyReused
		}
		var row settlement
		err := h.baseTx(c, tx).Where("b.id = ?", id).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&row).Error
		if err != nil {
			return err
		}
		if row.ID == 0 {
			return errNotFound
		}
		if row.IdempotencyKey != nil && *row.IdempotencyKey == key {
			out = row
			return nil
		}
		if row.Status != "bill_frozen" || row.Amount <= 0 {
			return errConflict
		}
		now := time.Now()
		applicationNo := fmt.Sprintf("SET-%d", row.ID)
		result := tx.Table("qixi_crm_m_settlement_bill").Where("id = ? AND store_id = ? AND status = ? AND idempotency_key IS NULL", row.ID, middleware.StoreID(c), "bill_frozen").Updates(map[string]any{
			"status": "withdraw_applied", "idempotency_key": key, "application_no": applicationNo,
			"applied_by_account_id": middleware.AdminID(c), "applied_at": now, "version": gorm.Expr("version + 1"),
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errConflict
		}
		row.Status, row.IdempotencyKey, row.ApplicationNo, row.AppliedAt, row.Version, row.UpdatedAt = "withdraw_applied", &key, &applicationNo, &now, row.Version+1, now
		if err := merchantsettlement.Enqueue(c.Request.Context(), tx, projection(row)); err != nil {
			return err
		}
		out = row
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, view(out))
	case errors.Is(err, errNotFound):
		response.Fail(c, http.StatusNotFound, "结算单不存在")
	case errors.Is(err, errConflict):
		response.Fail(c, http.StatusConflict, "当前结算单不可申请或已发生变化")
	case errors.Is(err, errIdempotencyKeyReused):
		response.Fail(c, http.StatusConflict, "幂等键已用于其他结算单")
	default:
		fail(c)
	}
}

func (h *Handler) base(c *gin.Context) *gorm.DB { return h.baseTx(c, h.db) }
func (h *Handler) baseTx(c *gin.Context, db *gorm.DB) *gorm.DB {
	return db.WithContext(c.Request.Context()).Table("qixi_crm_m_settlement_bill AS b").
		Select("b.id,b.store_id,b.merchant_id,m.name AS merchant_name,m.region_id,b.period_start,b.period_end,b.amount,b.status,b.idempotency_key,b.application_no,b.applied_at,b.review_note,b.payout_reference,b.paid_at,b.version,b.updated_at").
		Joins("JOIN qixi_crm_m_merchant AS m ON m.id = b.merchant_id").
		Where("b.store_id = ? AND b.merchant_id = ?", middleware.StoreID(c), middleware.MerID(c))
}

func projection(row settlement) merchantsettlement.Payload {
	return merchantsettlement.Payload{SettlementID: row.ID, MerchantID: row.MerchantID, StoreID: row.StoreID, MerchantName: row.MerchantName, RegionID: row.RegionID, PeriodStart: row.PeriodStart, PeriodEnd: row.PeriodEnd, Amount: row.Amount, Status: row.Status, UpdatedAt: row.UpdatedAt}
}
func views(rows []settlement) []gin.H {
	out := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		out = append(out, view(row))
	}
	return out
}
func view(row settlement) gin.H {
	return gin.H{"settlement_id": row.ID, "store_id": row.StoreID, "mer_id": row.MerchantID, "merchant_name": row.MerchantName, "period_start": row.PeriodStart, "period_end": row.PeriodEnd, "amount": row.Amount, "status": row.Status, "application_no": row.ApplicationNo, "review_note": row.ReviewNote, "paid_at": row.PaidAt, "updated_at": row.UpdatedAt}
}
func settlementStatus(raw string) (string, bool) {
	switch strings.TrimSpace(raw) {
	case "", "bill_pending", "bill_frozen", "withdraw_applied", "approved", "paid", "rejected":
		return strings.TrimSpace(raw), true
	default:
		return "", false
	}
}
func validIdempotencyKey(raw string) bool {
	length := len([]rune(strings.TrimSpace(raw)))
	return length >= 8 && length <= 128
}
func paging(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}
func id(c *gin.Context) uint64 { value, _ := strconv.ParseUint(c.Param("id"), 10, 64); return value }
func fail(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "结算申请处理失败")
}
