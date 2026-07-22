package diypersist

import (
	"context"

	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/diy"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func NewStoreAdapter(repo *Repo) *Repo { return repo }

func (r *Repo) List(ctx context.Context, merID uint, page, limit int) ([]diy.Page, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Model(&diy.Page{}).Where("is_del = 0 AND mer_id = ?", merID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []diy.Page
	err := q.Order("status DESC, id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) Get(ctx context.Context, id uint) (*diy.Page, error) {
	var row diy.Page
	err := r.db.WithContext(ctx).Where("id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetActiveHome(ctx context.Context, merID uint) (*diy.Page, error) {
	var row diy.Page
	err := r.db.WithContext(ctx).
		Where("mer_id = ? AND status = 1 AND type = 0 AND is_del = 0 AND is_show = 1", merID).
		Order("id DESC").
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) Create(ctx context.Context, p *diy.Page) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repo) Update(ctx context.Context, p *diy.Page) error {
	return r.db.WithContext(ctx).Model(&diy.Page{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
		"name":          p.Name,
		"title":         p.Title,
		"template_name": p.TemplateName,
		"value":         p.Value,
		"status":        p.Status,
		"is_default":    p.IsDefault,
	}).Error
}

func (r *Repo) ClearActive(ctx context.Context, merID uint) error {
	return r.db.WithContext(ctx).Model(&diy.Page{}).
		Where("mer_id = ? AND status = 1 AND is_del = 0", merID).
		Updates(map[string]interface{}{"status": 0, "is_default": 0}).Error
}

func (r *Repo) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&diy.Page{}).
		Where("id = ?", id).Update("is_del", 1).Error
}

var _ diy.Store = (*Repo)(nil)
