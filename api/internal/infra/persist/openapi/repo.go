package openapipersist

import (
	"context"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/openapi"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) GetByAccessKey(ctx context.Context, accessKey string) (*openapi.OpenAuth, error) {
	var row openapi.OpenAuth
	err := r.db.WithContext(ctx).
		Where("access_key = ? AND is_del = 0", accessKey).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) TouchLogin(ctx context.Context, id uint, ip string, at time.Time) error {
	return r.db.WithContext(ctx).Model(&openapi.OpenAuth{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_ip":   ip,
			"last_time": at,
		}).Error
}

type StoreAdapter struct{ *Repo }

func NewStoreAdapter(repo *Repo) *StoreAdapter { return &StoreAdapter{Repo: repo} }

var _ openapi.Store = (*StoreAdapter)(nil)
