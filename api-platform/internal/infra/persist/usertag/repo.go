package usertag

import (
	"context"

	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/usertag"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListLabel(ctx context.Context, page, limit int) ([]usertag.Label, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&usertag.Label{}).Where("is_del = 0")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []usertag.Label
	err := q.Order("sort DESC, label_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetLabel(ctx context.Context, id uint) (*usertag.Label, error) {
	var row usertag.Label
	err := r.db.WithContext(ctx).Where("label_id = ? AND is_del = 0", id).First(&row).Error
	return &row, err
}

func (r *Repo) CreateLabel(ctx context.Context, row *usertag.Label) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateLabel(ctx context.Context, row *usertag.Label) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) SoftDeleteLabel(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&usertag.Label{}).Where("label_id = ?", id).Update("is_del", 1).Error
}

func (r *Repo) ListGroup(ctx context.Context, page, limit int) ([]usertag.Group, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&usertag.Group{}).Where("is_del = 0")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []usertag.Group
	err := q.Order("sort DESC, group_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetGroup(ctx context.Context, id uint) (*usertag.Group, error) {
	var row usertag.Group
	err := r.db.WithContext(ctx).Where("group_id = ? AND is_del = 0", id).First(&row).Error
	return &row, err
}

func (r *Repo) CreateGroup(ctx context.Context, row *usertag.Group) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateGroup(ctx context.Context, row *usertag.Group) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) SoftDeleteGroup(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&usertag.Group{}).Where("group_id = ?", id).Update("is_del", 1).Error
}

func (r *Repo) ListUserLabels(ctx context.Context, uid uint) ([]usertag.Label, error) {
	var rows []usertag.Label
	err := r.db.WithContext(ctx).Table("qixi_m_app_user_label AS l").
		Joins("INNER JOIN qixi_m_app_user_label_relation r ON r.label_id = l.label_id").
		Where("r.uid = ? AND l.is_del = 0", uid).
		Select("l.*").Order("l.sort DESC, l.label_id DESC").Find(&rows).Error
	return rows, err
}

func (r *Repo) ReplaceUserLabels(ctx context.Context, uid uint, labelIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("uid = ?", uid).Delete(&usertag.Relation{}).Error; err != nil {
			return err
		}
		for _, lid := range labelIDs {
			if lid == 0 {
				continue
			}
			if err := tx.Create(&usertag.Relation{UID: uid, LabelID: lid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
