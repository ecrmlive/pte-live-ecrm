package presellpersist

import (
	"context"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/presell"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) List(ctx context.Context, merID *uint, onlyOn bool, page, limit int) ([]presell.Presell, int64, error) {
	q := r.db.WithContext(ctx).Model(&presell.Presell{}).Where("is_del = 0")
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	if onlyOn {
		q = q.Where("status = 1 AND is_show = 1 AND product_status = 1 AND action_status = 1")
		q = q.Where("start_time <= NOW() AND end_time >= NOW()")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []presell.Presell
	err := q.Order("product_presell_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) Get(ctx context.Context, id uint) (*presell.Presell, error) {
	var row presell.Presell
	err := r.db.WithContext(ctx).Where("product_presell_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetByProduct(ctx context.Context, productID uint) (*presell.Presell, error) {
	var row presell.Presell
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_del = 0 AND status = 1 AND is_show = 1 AND product_status = 1 AND action_status = 1", productID).
		Where("start_time <= NOW() AND end_time >= NOW()").
		Order("product_presell_id DESC").First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) Create(ctx context.Context, p *presell.Presell) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repo) Update(ctx context.Context, p *presell.Presell) error {
	return r.db.WithContext(ctx).Model(p).Where("product_presell_id = ?", p.ProductPresellID).Updates(map[string]interface{}{
		"start_time": p.StartTime, "end_time": p.EndTime,
		"final_start_time": p.FinalStartTime, "final_end_time": p.FinalEndTime,
		"status": p.Status, "presell_type": p.PresellType, "pay_count": p.PayCount,
		"delivery_type": p.DeliveryType, "delivery_day": p.DeliveryDay,
		"price": p.Price, "down_price": p.DownPrice, "final_price": p.FinalPrice, "stock": p.Stock,
		"is_show": p.IsShow, "store_name": p.StoreName,
		"store_info": p.StoreInfo, "product_status": p.ProductStatus,
		"action_status": p.ActionStatus, "refusal": p.Refusal,
	}).Error
}

func (r *Repo) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&presell.Presell{}).Where("product_presell_id = ?", id).
		Updates(map[string]interface{}{"is_del": 1, "is_show": 0, "action_status": -1}).Error
}

func (r *Repo) DecStock(ctx context.Context, id uint, num int) error {
	res := r.db.WithContext(ctx).Model(&presell.Presell{}).
		Where("product_presell_id = ? AND is_del = 0 AND stock >= ?", id, num).
		Update("stock", gorm.Expr("stock - ?", num))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repo) IncStock(ctx context.Context, id uint, num int) error {
	if num <= 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&presell.Presell{}).
		Where("product_presell_id = ? AND is_del = 0", id).
		Update("stock", gorm.Expr("stock + ?", num)).Error
}

func (r *Repo) IncSeles(ctx context.Context, id uint, num int) error {
	return r.db.WithContext(ctx).Model(&presell.Presell{}).
		Where("product_presell_id = ?", id).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Update("seles", gorm.Expr("seles + ?", num)).Error
}

func (r *Repo) LoadProductMeta(ctx context.Context, productID uint) (storeName, image, merName string, price, cost float64, merID uint, err error) {
	var row struct {
		StoreName string  `gorm:"column:store_name"`
		Image     string  `gorm:"column:image"`
		Price     float64 `gorm:"column:price"`
		Cost      float64 `gorm:"column:cost"`
		MerID     uint    `gorm:"column:mer_id"`
		MerName   string  `gorm:"column:mer_name"`
	}
	err = r.db.WithContext(ctx).Table("qixi_m_admin_store_product AS p").
		Select("p.store_name, p.image, p.price, p.cost, p.mer_id, m.mer_name").
		Joins("LEFT JOIN qixi_m_admin_merchant m ON m.mer_id = p.mer_id").
		Where("p.product_id = ? AND p.is_del = 0", productID).
		Scan(&row).Error
	if err != nil {
		return "", "", "", 0, 0, 0, err
	}
	if row.MerID == 0 {
		return "", "", "", 0, 0, 0, gorm.ErrRecordNotFound
	}
	return row.StoreName, row.Image, row.MerName, row.Price, row.Cost, row.MerID, nil
}

func (r *Repo) CreatePresellOrder(ctx context.Context, o *presell.PresellOrder) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *Repo) GetPresellOrder(ctx context.Context, id uint) (*presell.PresellOrder, error) {
	var row presell.PresellOrder
	err := r.db.WithContext(ctx).Where("presell_order_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetPresellOrderByOrderID(ctx context.Context, orderID uint) (*presell.PresellOrder, error) {
	var row presell.PresellOrder
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) ListPresellOrdersByUID(ctx context.Context, uid uint, unpaidOnly bool, page, limit int) ([]presell.PresellOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Model(&presell.PresellOrder{}).Where("uid = ?", uid)
	if unpaidOnly {
		q = q.Where("paid = 0 AND status = 1")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []presell.PresellOrder
	err := q.Order("presell_order_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) MarkPresellOrderPaid(ctx context.Context, id uint, payType int8, at time.Time) error {
	return r.db.WithContext(ctx).Model(&presell.PresellOrder{}).
		Where("presell_order_id = ? AND paid = 0", id).
		Updates(map[string]interface{}{"paid": 1, "pay_type": payType, "pay_time": at}).Error
}

func (r *Repo) InvalidatePresellOrder(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&presell.PresellOrder{}).
		Where("presell_order_id = ? AND paid = 0", id).
		Update("status", 0).Error
}

var _ presell.Store = (*Repo)(nil)
