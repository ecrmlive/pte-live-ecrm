package aftersalepersist

import (
	"context"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/aftersale"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/catalog"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) WithTx(fn func(tx aftersale.Store) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&Repo{db: tx})
	})
}

func (r *Repo) GetStoreOrder(ctx context.Context, orderID uint) (*aftersale.StoreOrderBrief, error) {
	var o aftersale.StoreOrderBrief
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repo) ListOrderProducts(ctx context.Context, orderID uint) ([]aftersale.OrderProductLine, error) {
	var rows []aftersale.OrderProductLine
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&rows).Error
	return rows, err
}

func (r *Repo) ListOrderProductsByIDs(ctx context.Context, orderID uint, ids []uint) ([]aftersale.OrderProductLine, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []aftersale.OrderProductLine
	err := r.db.WithContext(ctx).Where("order_id = ? AND order_product_id IN ?", orderID, ids).Find(&rows).Error
	return rows, err
}

func (r *Repo) HasActiveRefund(ctx context.Context, orderID uint) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&aftersale.RefundOrder{}).
		Where("order_id = ? AND is_del = 0 AND status IN ?", orderID,
			[]int8{aftersale.StatusWait, aftersale.StatusBack, aftersale.StatusReceive, aftersale.StatusPlatform}).
		Count(&n).Error
	return n > 0, err
}

func (r *Repo) CreateRefundOrder(ctx context.Context, o *aftersale.RefundOrder) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *Repo) CreateRefundProducts(ctx context.Context, rows []aftersale.RefundProduct) error {
	if len(rows) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&rows).Error
}

func (r *Repo) CreateRefundStatus(ctx context.Context, log *aftersale.RefundStatusLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *Repo) GetRefund(ctx context.Context, id uint) (*aftersale.RefundOrder, error) {
	var o aftersale.RefundOrder
	err := r.db.WithContext(ctx).Where("refund_order_id = ?", id).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repo) ListRefunds(ctx context.Context, filter aftersale.ListFilter) ([]aftersale.RefundOrder, int64, error) {
	page, limit := normalize(filter.Page, filter.Limit)
	q := r.db.WithContext(ctx).Model(&aftersale.RefundOrder{}).Where("is_del = 0")
	if filter.UID != nil {
		q = q.Where("uid = ?", *filter.UID)
	}
	if filter.MerID != nil {
		q = q.Where("mer_id = ?", *filter.MerID)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []aftersale.RefundOrder
	err := q.Order("refund_order_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

// ListPlatformRefundsByRegions 在 SQL 层以商户所属区域过滤退款单，区域账号不得读取范围外数据。
func (r *Repo) ListPlatformRefundsByRegions(ctx context.Context, regionIDs []uint, status *int8, page, limit int) ([]aftersale.RefundOrder, int64, error) {
	page, limit = normalize(page, limit)
	if len(regionIDs) == 0 {
		return []aftersale.RefundOrder{}, 0, nil
	}
	q := r.db.WithContext(ctx).Table("qixi_m_app_store_refund_order AS r").
		Joins("JOIN qixi_m_admin_merchant AS m ON m.mer_id = r.mer_id").
		Where("r.is_del = 0 AND m.is_del = 0 AND m.region_id IN ?", regionIDs)
	if status != nil {
		q = q.Where("r.status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []aftersale.RefundOrder
	err := q.Select("r.*").Order("r.refund_order_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	return rows, total, err
}

// GetPlatformRefundByRegions 在 SQL 层按退款单所属商户区域过滤详情。
func (r *Repo) GetPlatformRefundByRegions(ctx context.Context, id uint, regionIDs []uint) (*aftersale.RefundOrder, error) {
	if len(regionIDs) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var row aftersale.RefundOrder
	err := r.db.WithContext(ctx).Table("qixi_m_app_store_refund_order AS r").
		Select("r.*").Joins("JOIN qixi_m_admin_merchant AS m ON m.mer_id = r.mer_id").
		Where("r.refund_order_id = ? AND r.is_del = 0 AND m.is_del = 0 AND m.region_id IN ?", id, regionIDs).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) ListRefundProducts(ctx context.Context, refundOrderID uint) ([]aftersale.RefundProduct, error) {
	var rows []aftersale.RefundProduct
	err := r.db.WithContext(ctx).Where("refund_order_id = ?", refundOrderID).
		Order("refund_product_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) UpdateRefundStatus(ctx context.Context, id uint, fromStatus, toStatus int8, failMessage string) (bool, error) {
	now := time.Now()
	updates := map[string]interface{}{
		"status":      toStatus,
		"status_time": now,
	}
	if failMessage != "" {
		updates["fail_message"] = failMessage
	}
	res := r.db.WithContext(ctx).Model(&aftersale.RefundOrder{}).
		Where("refund_order_id = ? AND status = ? AND is_del = 0", id, fromStatus).
		Updates(updates)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *Repo) AddUserBalance(ctx context.Context, uid uint, amount float64) error {
	return r.db.WithContext(ctx).Table("qixi_m_app_user").
		Where("uid = ?", uid).
		Update("now_money", gorm.Expr("now_money + ?", amount)).Error
}

func (r *Repo) AddProductStock(ctx context.Context, productID uint, num uint) error {
	return r.db.WithContext(ctx).Model(&catalog.Product{}).
		Where("product_id = ?", productID).
		Updates(map[string]interface{}{
			"stock": gorm.Expr("stock + ?", num),
			"sales": gorm.Expr("CASE WHEN sales >= ? THEN sales - ? ELSE 0 END", num, num),
		}).Error
}

func (r *Repo) AddSKUStock(ctx context.Context, productID uint, unique string, num uint) error {
	q := r.db.WithContext(ctx).Model(&catalog.AttrValue{}).
		Where("product_id = ?", productID)
	if unique != "" {
		q = q.Where("`unique` = ?", unique)
	}
	return q.Updates(map[string]interface{}{
		"stock": gorm.Expr("stock + ?", num),
		"sales": gorm.Expr("CASE WHEN sales >= ? THEN sales - ? ELSE 0 END", num, num),
	}).Error
}

func (r *Repo) MarkOrderProductRefund(ctx context.Context, orderProductID uint, addNum int, isRefund int8) error {
	updates := map[string]interface{}{
		"is_refund": isRefund,
	}
	if addNum > 0 {
		updates["refund_num"] = gorm.Expr("refund_num + ?", addNum)
	}
	return r.db.WithContext(ctx).Model(&aftersale.OrderProductLine{}).
		Where("order_product_id = ?", orderProductID).
		Updates(updates).Error
}

func (r *Repo) UpdateOrderStatus(ctx context.Context, orderID uint, status int8) error {
	return r.db.WithContext(ctx).Model(&aftersale.StoreOrderBrief{}).
		Where("order_id = ?", orderID).Update("status", status).Error
}

func (r *Repo) CountUnrefundedProducts(ctx context.Context, orderID uint) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&aftersale.OrderProductLine{}).
		Where("order_id = ? AND (is_refund <> ? OR refund_num < product_num)", orderID, aftersale.OrderProductRefundFull).
		Count(&n).Error
	return n, err
}

func normalize(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}

var _ aftersale.Store = (*Repo)(nil)
