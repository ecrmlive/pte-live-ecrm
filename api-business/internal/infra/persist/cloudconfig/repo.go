package cloudconfigpersist

import (
	"context"
	"errors"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/cloudconfig"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListByGroup(ctx context.Context, group string) ([]cloudconfig.Config, error) {
	var rows []cloudconfig.Config
	err := r.db.WithContext(ctx).Where("group_key = ?", group).Order("config_key ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) Upsert(ctx context.Context, row *cloudconfig.Config) error {
	var old cloudconfig.Config
	err := r.db.WithContext(ctx).Where("group_key = ? AND config_key = ?", row.GroupKey, row.ConfigKey).First(&old).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.WithContext(ctx).Create(row).Error
	}
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&cloudconfig.Config{}).Where("config_id = ?", old.ConfigID).Updates(map[string]any{
		"ciphertext": row.Ciphertext, "is_secret": row.IsSecret, "updated_by": row.UpdatedBy, "update_time": time.Now(),
	}).Error
}

var _ cloudconfig.Store = (*Repo)(nil)
