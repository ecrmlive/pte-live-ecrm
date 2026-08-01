package fulfillment

import (
	"context"

	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/fulfillment"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListStaff(ctx context.Context, merID uint, page, limit int) ([]fulfillment.Staff, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&fulfillment.Staff{}).Where("mer_id = ? AND is_del = 0", merID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []fulfillment.Staff
	err := q.Order("staff_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetStaff(ctx context.Context, merID, id uint) (*fulfillment.Staff, error) {
	var row fulfillment.Staff
	err := r.db.WithContext(ctx).Where("staff_id = ? AND mer_id = ? AND is_del = 0", id, merID).First(&row).Error
	return &row, err
}

func (r *Repo) CreateStaff(ctx context.Context, row *fulfillment.Staff) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateStaff(ctx context.Context, row *fulfillment.Staff) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) SoftDeleteStaff(ctx context.Context, merID, id uint) error {
	return r.db.WithContext(ctx).Model(&fulfillment.Staff{}).
		Where("staff_id = ? AND mer_id = ?", id, merID).
		Update("is_del", 1).Error
}
