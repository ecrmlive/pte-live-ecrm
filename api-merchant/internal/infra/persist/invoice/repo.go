package invoice

import (
	"context"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/invoice"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListByUID(ctx context.Context, uid uint, page, limit int) ([]invoice.Invoice, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&invoice.Invoice{}).Where("uid = ? AND is_del = 0", uid)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []invoice.Invoice
	err := q.Order("invoice_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListByMer(ctx context.Context, merID uint, page, limit int) ([]invoice.Invoice, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&invoice.Invoice{}).Where("mer_id = ? AND is_del = 0", merID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []invoice.Invoice
	err := q.Order("invoice_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) Get(ctx context.Context, id uint) (*invoice.Invoice, error) {
	var row invoice.Invoice
	err := r.db.WithContext(ctx).Where("invoice_id = ?", id).First(&row).Error
	return &row, err
}

func (r *Repo) FindByOrder(ctx context.Context, orderID uint) (*invoice.Invoice, error) {
	var row invoice.Invoice
	err := r.db.WithContext(ctx).Where("order_id = ? AND is_del = 0", orderID).First(&row).Error
	return &row, err
}

func (r *Repo) Create(ctx context.Context, row *invoice.Invoice) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) Update(ctx context.Context, row *invoice.Invoice) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) LoadOrder(ctx context.Context, orderID uint) (*invoice.OrderMeta, error) {
	var row struct {
		OrderID uint `gorm:"column:order_id"`
		UID     uint `gorm:"column:uid"`
		MerID   uint `gorm:"column:mer_id"`
		Paid    int8 `gorm:"column:paid"`
		Status  int8 `gorm:"column:status"`
		IsDel   int8 `gorm:"column:is_del"`
	}
	err := r.db.WithContext(ctx).Table("qixi_m_app_store_order").
		Select("order_id, uid, mer_id, paid, status, is_del").
		Where("order_id = ?", orderID).First(&row).Error
	if err != nil {
		return nil, err
	}
	if row.IsDel == 1 {
		return nil, gorm.ErrRecordNotFound
	}
	return &invoice.OrderMeta{
		OrderID: row.OrderID, UID: row.UID, MerID: row.MerID,
		Paid: row.Paid == 1, Status: row.Status,
	}, nil
}
