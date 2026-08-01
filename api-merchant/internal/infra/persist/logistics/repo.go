package logistics

import (
	"context"

	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/logistics"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListExpress(ctx context.Context, page, limit int, showOnly bool) ([]logistics.Express, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&logistics.Express{}).Where("is_del = 0")
	if showOnly {
		q = q.Where("is_show = 1")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []logistics.Express
	err := q.Order("sort DESC, express_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetExpress(ctx context.Context, id uint) (*logistics.Express, error) {
	var row logistics.Express
	err := r.db.WithContext(ctx).Where("express_id = ? AND is_del = 0", id).First(&row).Error
	return &row, err
}

func (r *Repo) CreateExpress(ctx context.Context, row *logistics.Express) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateExpress(ctx context.Context, row *logistics.Express) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) SoftDeleteExpress(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&logistics.Express{}).Where("express_id = ?", id).Update("is_del", 1).Error
}

func (r *Repo) ListCity(ctx context.Context, parentID *uint) ([]logistics.City, error) {
	q := r.db.WithContext(ctx).Model(&logistics.City{}).Where("is_show = 1")
	if parentID != nil {
		q = q.Where("parent_id = ?", *parentID)
	}
	var rows []logistics.City
	err := q.Order("city_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) ListTemplate(ctx context.Context, merID uint, page, limit int) ([]logistics.ShippingTemplate, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&logistics.ShippingTemplate{}).Where("mer_id = ? AND is_del = 0", merID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []logistics.ShippingTemplate
	err := q.Order("sort DESC, template_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetTemplate(ctx context.Context, merID, id uint) (*logistics.ShippingTemplate, error) {
	var row logistics.ShippingTemplate
	err := r.db.WithContext(ctx).Where("template_id = ? AND mer_id = ? AND is_del = 0", id, merID).First(&row).Error
	return &row, err
}

func (r *Repo) ListRegions(ctx context.Context, templateID uint) ([]logistics.Region, error) {
	var rows []logistics.Region
	err := r.db.WithContext(ctx).Where("template_id = ?", templateID).Order("region_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) CreateTemplate(ctx context.Context, row *logistics.ShippingTemplate, regions []logistics.Region) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(row).Error; err != nil {
			return err
		}
		for i := range regions {
			regions[i].TemplateID = row.TemplateID
			regions[i].RegionID = 0
			if err := tx.Create(&regions[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) UpdateTemplate(ctx context.Context, row *logistics.ShippingTemplate, regions []logistics.Region) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(row).Error; err != nil {
			return err
		}
		if err := tx.Where("template_id = ?", row.TemplateID).Delete(&logistics.Region{}).Error; err != nil {
			return err
		}
		for i := range regions {
			regions[i].TemplateID = row.TemplateID
			regions[i].RegionID = 0
			if err := tx.Create(&regions[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) SetDefaultTemplate(ctx context.Context, merID, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&logistics.ShippingTemplate{}).
			Where("template_id = ? AND mer_id = ? AND is_del = 0", id, merID).
			Update("is_default", 1)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Model(&logistics.ShippingTemplate{}).
			Where("mer_id = ? AND template_id <> ? AND is_del = 0", merID, id).
			Update("is_default", 0).Error
	})
}

func (r *Repo) SoftDeleteTemplate(ctx context.Context, merID, id uint) error {
	return r.db.WithContext(ctx).Model(&logistics.ShippingTemplate{}).
		Where("template_id = ? AND mer_id = ?", id, merID).Update("is_del", 1).Error
}
