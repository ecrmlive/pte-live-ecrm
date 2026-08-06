package seckillpersist

import (
	"context"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/seckill"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/listquery"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListTimes(ctx context.Context) ([]seckill.TimeSlot, error) {
	var rows []seckill.TimeSlot
	err := r.db.WithContext(ctx).Where("status = 1").Order("start_time ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) ListActives(ctx context.Context, merID *uint, onlyOn bool, page, limit int, filter listquery.AdminFilter) ([]seckill.Active, int64, error) {
	q := r.db.WithContext(ctx).Model(&seckill.Active{}).Where("delete_time IS NULL")
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	if onlyOn {
		q = q.Where("status = 1 AND active_status = 1")
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	q = listquery.ApplyKeywordLike(q, filter.Keyword, "name")
	q = listquery.ApplyUnixColumnDateRange(q, "create_time", filter.DateFrom, filter.DateTo)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seckill.Active
	err := q.Order("seckill_active_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetActive(ctx context.Context, id uint) (*seckill.Active, error) {
	var row seckill.Active
	err := r.db.WithContext(ctx).Where("seckill_active_id = ? AND delete_time IS NULL", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetActiveByProduct(ctx context.Context, productID uint) (*seckill.Active, error) {
	var row seckill.Active
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND status = 1 AND active_status = 1 AND delete_time IS NULL", productID).
		Order("seckill_active_id DESC").First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateActive(ctx context.Context, a *seckill.Active) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *Repo) UpdateActive(ctx context.Context, a *seckill.Active) error {
	return r.db.WithContext(ctx).Model(a).Where("seckill_active_id = ?", a.SeckillActiveID).Updates(map[string]interface{}{
		"name": a.Name, "seckill_time_ids": a.SeckillTimeIDs, "start_day": a.StartDay, "end_day": a.EndDay,
		"seckill_price": a.SeckillPrice, "once_pay_count": a.OncePayCount, "status": a.Status, "update_time": a.UpdateTime,
	}).Error
}

func (r *Repo) SoftDeleteActive(ctx context.Context, id uint) error {
	now := time.Now().Unix()
	return r.db.WithContext(ctx).Model(&seckill.Active{}).Where("seckill_active_id = ?", id).
		Updates(map[string]interface{}{"delete_time": now, "status": 0}).Error
}

func (r *Repo) LoadProductMeta(ctx context.Context, productID uint) (storeName, image, merName string, price float64, merID uint, err error) {
	var row struct {
		StoreName string  `gorm:"column:store_name"`
		Image     string  `gorm:"column:image"`
		Price     float64 `gorm:"column:price"`
		MerID     uint    `gorm:"column:mer_id"`
		MerName   string  `gorm:"column:mer_name"`
	}
	err = r.db.WithContext(ctx).Table("qixi_m_admin_store_product AS p").
		Select("p.store_name, p.image, p.price, p.mer_id, m.mer_name").
		Joins("LEFT JOIN qixi_m_admin_merchant m ON m.mer_id = p.mer_id").
		Where("p.product_id = ? AND p.is_del = 0", productID).
		Scan(&row).Error
	if err != nil {
		return "", "", "", 0, 0, err
	}
	if row.MerID == 0 {
		return "", "", "", 0, 0, gorm.ErrRecordNotFound
	}
	return row.StoreName, row.Image, row.MerName, row.Price, row.MerID, nil
}

var _ seckill.Store = (*Repo)(nil)
