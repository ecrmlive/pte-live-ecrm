package contentpersist

import (
	"context"
	"errors"

	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/content"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListNotices(ctx context.Context, onlyShow bool, page, limit int) ([]content.Notice, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Model(&content.Notice{}).Where("is_del = 0")
	if onlyShow {
		q = q.Where("is_show = 1")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []content.Notice
	err := q.Order("sort DESC, notice_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetNotice(ctx context.Context, id uint) (*content.Notice, error) {
	var row content.Notice
	err := r.db.WithContext(ctx).Where("notice_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateNotice(ctx context.Context, n *content.Notice) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *Repo) UpdateNotice(ctx context.Context, n *content.Notice) error {
	return r.db.WithContext(ctx).Model(n).Where("notice_id = ?", n.NoticeID).Updates(map[string]interface{}{
		"title":   n.Title,
		"content": n.Content,
		"is_show": n.IsShow,
		"sort":    n.Sort,
	}).Error
}

func (r *Repo) SoftDeleteNotice(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&content.Notice{}).
		Where("notice_id = ?", id).Update("is_del", 1).Error
}

func (r *Repo) GetCache(ctx context.Context, key string) (*content.Cache, error) {
	var row content.Cache
	err := r.db.WithContext(ctx).Where("`key` = ?", key).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) UpsertCache(ctx context.Context, row *content.Cache) error {
	var exist content.Cache
	err := r.db.WithContext(ctx).Where("`key` = ?", row.Key).First(&exist).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.WithContext(ctx).Create(row).Error
		}
		return err
	}
	return r.db.WithContext(ctx).Model(&content.Cache{}).Where("`key` = ?", row.Key).
		Updates(map[string]interface{}{
			"result":      row.Result,
			"expire_time": row.ExpireTime,
		}).Error
}

var _ content.Store = (*Repo)(nil)
