package productmeta

import (
	"context"

	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/productmeta"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListLabel(ctx context.Context, merID uint, page, limit int) ([]productmeta.Label, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&productmeta.Label{}).Where("mer_id = ? AND is_del = 0", merID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []productmeta.Label
	err := q.Order("sort DESC, label_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListLabelPlatform(ctx context.Context, page, limit int) ([]productmeta.Label, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&productmeta.Label{}).Where("is_del = 0")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []productmeta.Label
	err := q.Order("label_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetLabel(ctx context.Context, merID, id uint) (*productmeta.Label, error) {
	var row productmeta.Label
	err := r.db.WithContext(ctx).Where("label_id = ? AND mer_id = ? AND is_del = 0", id, merID).First(&row).Error
	return &row, err
}

func (r *Repo) CreateLabel(ctx context.Context, row *productmeta.Label) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateLabel(ctx context.Context, row *productmeta.Label) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) SoftDeleteLabel(ctx context.Context, merID, id uint) error {
	return r.db.WithContext(ctx).Model(&productmeta.Label{}).Where("label_id = ? AND mer_id = ?", id, merID).Update("is_del", 1).Error
}

func (r *Repo) ListGuarantee(ctx context.Context, merID uint, page, limit int) ([]productmeta.Guarantee, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&productmeta.Guarantee{}).Where("mer_id = ? AND is_del = 0", merID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []productmeta.Guarantee
	err := q.Order("sort DESC, guarantee_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListGuaranteePlatform(ctx context.Context, page, limit int) ([]productmeta.Guarantee, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&productmeta.Guarantee{}).Where("is_del = 0")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []productmeta.Guarantee
	err := q.Order("guarantee_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetGuarantee(ctx context.Context, merID, id uint) (*productmeta.Guarantee, error) {
	var row productmeta.Guarantee
	err := r.db.WithContext(ctx).Where("guarantee_id = ? AND mer_id = ? AND is_del = 0", id, merID).First(&row).Error
	return &row, err
}

func (r *Repo) CreateGuarantee(ctx context.Context, row *productmeta.Guarantee) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateGuarantee(ctx context.Context, row *productmeta.Guarantee) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) SoftDeleteGuarantee(ctx context.Context, merID, id uint) error {
	return r.db.WithContext(ctx).Model(&productmeta.Guarantee{}).Where("guarantee_id = ? AND mer_id = ?", id, merID).Update("is_del", 1).Error
}

func (r *Repo) ListAttrTemplate(ctx context.Context, merID uint, page, limit int) ([]productmeta.AttrTemplate, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&productmeta.AttrTemplate{}).Where("mer_id = ? AND is_del = 0", merID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []productmeta.AttrTemplate
	err := q.Order("sort DESC, template_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetAttrTemplate(ctx context.Context, merID, id uint) (*productmeta.AttrTemplate, error) {
	var row productmeta.AttrTemplate
	err := r.db.WithContext(ctx).Where("template_id = ? AND mer_id = ? AND is_del = 0", id, merID).First(&row).Error
	return &row, err
}

func (r *Repo) CreateAttrTemplate(ctx context.Context, row *productmeta.AttrTemplate) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateAttrTemplate(ctx context.Context, row *productmeta.AttrTemplate) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) SoftDeleteAttrTemplate(ctx context.Context, merID, id uint) error {
	return r.db.WithContext(ctx).Model(&productmeta.AttrTemplate{}).Where("template_id = ? AND mer_id = ?", id, merID).Update("is_del", 1).Error
}
