package cloudconfigpersist

import (
	"context"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/cloudconfig"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListByGroup(ctx context.Context, group string) ([]cloudconfig.Config, error) {
	var rows []cloudconfig.Config
	err := r.db.WithContext(ctx).Where("provider = ?", group).Order("config_key ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) Upsert(ctx context.Context, row *cloudconfig.Config) error {
	if row.KeyVersion == "" {
		row.KeyVersion = "v1"
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider"}, {Name: "config_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"ciphertext", "key_version", "updated_by", "updated_at"}),
	}).Create(row).Error
}

var _ cloudconfig.Store = (*Repo)(nil)
