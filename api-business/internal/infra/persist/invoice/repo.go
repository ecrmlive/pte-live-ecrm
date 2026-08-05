package invoice

import (
	"context"
	"errors"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/invoice"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListProfiles(ctx context.Context, userID uint64) ([]invoice.InvoiceProfile, error) {
	var rows []invoice.InvoiceProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC, id DESC").Find(&rows).Error
	return rows, err
}

func (r *Repo) GetProfile(ctx context.Context, userID, id uint64) (*invoice.InvoiceProfile, error) {
	var row invoice.InvoiceProfile
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&row).Error
	return &row, err
}

func (r *Repo) CreateProfile(ctx context.Context, row *invoice.InvoiceProfile) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateProfile(ctx context.Context, row *invoice.InvoiceProfile) error {
	return r.db.WithContext(ctx).Model(&invoice.InvoiceProfile{}).Where("id = ? AND user_id = ?", row.ID, row.UserID).Updates(map[string]any{
		"type": row.Type, "title": row.Title, "tax_no": row.TaxNo, "email": row.Email, "is_default": row.IsDefault,
	}).Error
}

func (r *Repo) DeleteProfile(ctx context.Context, userID, id uint64) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&invoice.InvoiceProfile{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repo) SetDefaultProfile(ctx context.Context, userID, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&invoice.InvoiceProfile{}).Where("id = ? AND user_id = ?", id, userID).Update("is_default", true)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Model(&invoice.InvoiceProfile{}).Where("user_id = ? AND id <> ?", userID, id).Update("is_default", false).Error
	})
}

func (r *Repo) ListByUID(ctx context.Context, userID uint64, page, limit int) ([]invoice.Invoice, int64, error) {
	base := r.db.WithContext(ctx).Table("qixi_crm_b_order_invoice AS oi").Joins("JOIN qixi_crm_b_order AS o ON o.id = oi.order_id").Where("o.user_id = ?", userID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []invoice.Invoice
	err := base.Select("oi.id, oi.order_id, oi.invoice_profile_id, oi.profile_type, oi.title, oi.tax_no, oi.email, oi.status, oi.invoice_no, oi.file_url, oi.rejection_reason, oi.requested_at, oi.issued_at, oi.updated_at").Order("oi.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	return rows, total, err
}

func (r *Repo) GetByUID(ctx context.Context, userID, id uint64) (*invoice.Invoice, error) {
	var row invoice.Invoice
	err := r.db.WithContext(ctx).Table("qixi_crm_b_order_invoice AS oi").
		Select("oi.id, oi.order_id, oi.invoice_profile_id, oi.profile_type, oi.title, oi.tax_no, oi.email, oi.status, oi.invoice_no, oi.file_url, oi.rejection_reason, oi.requested_at, oi.issued_at, oi.updated_at").
		Joins("JOIN qixi_crm_b_order AS o ON o.id = oi.order_id").Where("oi.id = ? AND o.user_id = ?", id, userID).Scan(&row).Error
	if err == nil && row.ID == 0 {
		err = gorm.ErrRecordNotFound
	}
	return &row, err
}

func (r *Repo) FindByOrder(ctx context.Context, orderID uint64) (*invoice.Invoice, error) {
	var row invoice.Invoice
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&row).Error
	return &row, err
}

func (r *Repo) Create(ctx context.Context, row *invoice.Invoice) error {
	err := r.db.WithContext(ctx).Create(row).Error
	if err != nil && (errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate")) {
		return invoice.ErrExists
	}
	return err
}

func (r *Repo) LoadOrder(ctx context.Context, orderID uint64) (*invoice.OrderMeta, error) {
	var row invoice.OrderMeta
	err := r.db.WithContext(ctx).Table("qixi_crm_b_order").Select("id AS order_id, user_id, status").Where("id = ?", orderID).Scan(&row).Error
	if err == nil && row.OrderID == 0 {
		err = gorm.ErrRecordNotFound
	}
	return &row, err
}
