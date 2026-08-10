package presellpersist

import (
	"context"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/presell"
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

func (r *Repo) applyAdminFilters(q *gorm.DB, query presell.AdminQuery, presellType int) *gorm.DB {
	q = q.Where("is_del = 0")
	if query.MerID != nil {
		q = q.Where("mer_id = ?", *query.MerID)
	}
	if len(query.MerIDs) > 0 {
		q = q.Where("mer_id IN ?", query.MerIDs)
	}
	if star := query.Star; star != nil {
		q = q.Where("star = ?", *star)
	}
	if ps := query.ProductStatus; ps != nil {
		if *ps == -1 {
			q = q.Where("product_status IN ?", []int{-1, -2})
		} else {
			q = q.Where("product_status = ?", *ps)
		}
	}
	if us := query.UsStatus; us != nil {
		switch *us {
		case 1:
			q = q.Where("product_status = 1 AND status = 1 AND is_show = 1")
		case 0:
			q = q.Where("product_status = 1 AND status = 1 AND is_show = 0")
		case -1:
			q = q.Where("(status = 0 OR product_status <> 1)")
		}
	}
	if act := query.ActivityType; act != nil {
		switch *act {
		case 0: // 未开始
			q = q.Where("action_status = 1").Where(`(
				start_time > NOW()
				OR (
					start_time <= NOW() AND end_time > NOW()
					AND (product_status <> 1 OR is_show <> 1 OR status <> 1)
				)
			)`)
		case 1: // 进行中
			q = q.Where("action_status = 1 AND start_time <= NOW() AND end_time > NOW()").
				Where("product_status = 1 AND status = 1 AND is_show = 1")
		case 2: // 已结束
			q = q.Where("(action_status = -1 OR end_time <= NOW())")
		}
	}
	if labels := strings.TrimSpace(query.SysLabels); labels != "" {
		q = q.Where("FIND_IN_SET(?, REPLACE(sys_labels,' ',''))", labels)
	}
	pt := presellType
	if pt == 0 {
		pt = query.PresellType
	}
	if pt == 1 || pt == 2 {
		q = q.Where("presell_type = ?", pt)
	}
	if kw := strings.TrimSpace(query.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("(store_name LIKE ? OR CAST(product_id AS CHAR) LIKE ? OR CAST(product_presell_id AS CHAR) LIKE ?)", like, like, like)
	}
	return q
}

func (r *Repo) ListAdmin(ctx context.Context, query presell.AdminQuery) ([]presell.Presell, int64, error) {
	q := r.db.WithContext(ctx).Model(&presell.Presell{})
	q = r.applyAdminFilters(q, query, query.PresellType)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []presell.Presell
	err := q.Order("star DESC, product_presell_id DESC").
		Offset((query.Page - 1) * query.Limit).Limit(query.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CountAdmin(ctx context.Context, query presell.AdminQuery, presellType int) (int64, error) {
	q := r.db.WithContext(ctx).Model(&presell.Presell{})
	q = r.applyAdminFilters(q, query, presellType)
	var total int64
	err := q.Count(&total).Error
	return total, err
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
		"star": p.Star, "sys_labels": p.SysLabels, "stock_count": p.StockCount,
		"attend_num": p.AttendNum, "success_num": p.SuccessNum,
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
	err = r.db.WithContext(ctx).Table("qixi_crm_b_product_view AS p").
		Select("p.title AS store_name, p.cover_url AS image, p.price, p.merchant_id AS mer_id, p.merchant_name AS mer_name, COALESCE(p.original_price, p.price) AS cost").
		Where("p.product_id = ? AND p.sale_status = 1", productID).
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
