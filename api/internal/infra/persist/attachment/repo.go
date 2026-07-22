package attachmentpersist

import (
	"context"

	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/attachment"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

var _ attachment.Store = (*Repo)(nil)

func (r *Repo) ListCategories(ctx context.Context, merID uint) ([]attachment.Category, error) {
	var rows []attachment.Category
	err := r.db.WithContext(ctx).Where("mer_id = ?", merID).
		Order("sort DESC, attachment_category_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) CreateCategory(ctx context.Context, c *attachment.Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repo) UpdateCategory(ctx context.Context, c *attachment.Category) error {
	return r.db.WithContext(ctx).Model(c).Where("attachment_category_id = ?", c.AttachmentCategoryID).
		Updates(map[string]interface{}{
			"pid": c.PID, "attachment_category_name": c.AttachmentCategoryName,
			"attachment_category_enname": c.AttachmentCategoryEnname, "sort": c.Sort,
		}).Error
}

func (r *Repo) DeleteCategory(ctx context.Context, id, merID uint) error {
	return r.db.WithContext(ctx).Where("attachment_category_id = ? AND mer_id = ?", id, merID).
		Delete(&attachment.Category{}).Error
}

func (r *Repo) GetCategory(ctx context.Context, id uint) (*attachment.Category, error) {
	var row attachment.Category
	err := r.db.WithContext(ctx).Where("attachment_category_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) List(ctx context.Context, userType int, cateID uint, page, limit int) ([]attachment.Attachment, int64, error) {
	q := r.db.WithContext(ctx).Model(&attachment.Attachment{}).Where("user_type = ?", userType)
	if cateID > 0 {
		q = q.Where("attachment_category_id = ?", cateID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []attachment.Attachment
	err := q.Order("attachment_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) Get(ctx context.Context, id uint) (*attachment.Attachment, error) {
	var row attachment.Attachment
	err := r.db.WithContext(ctx).Where("attachment_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) Create(ctx context.Context, a *attachment.Attachment) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *Repo) Delete(ctx context.Context, id uint, userType int) error {
	res := r.db.WithContext(ctx).Where("attachment_id = ? AND user_type = ?", id, userType).
		Delete(&attachment.Attachment{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
