package cspersist

import (
	"context"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/cs"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListReplies(ctx context.Context, merID uint, onlyOn bool, page, limit int) ([]cs.Reply, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Model(&cs.Reply{}).Where("mer_id = ?", merID)
	if onlyOn {
		q = q.Where("status = 1")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []cs.Reply
	err := q.Order("service_reply_id ASC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetReply(ctx context.Context, merID, id uint) (*cs.Reply, error) {
	var row cs.Reply
	err := r.db.WithContext(ctx).Where("service_reply_id = ? AND mer_id = ?", id, merID).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateReply(ctx context.Context, row *cs.Reply) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateReply(ctx context.Context, row *cs.Reply) error {
	return r.db.WithContext(ctx).Model(row).Where("service_reply_id = ? AND mer_id = ?", row.ServiceReplyID, row.MerID).
		Updates(map[string]interface{}{
			"type":    row.Type,
			"keyword": row.Keyword,
			"content": row.Content,
			"status":  row.Status,
		}).Error
}

func (r *Repo) DeleteReply(ctx context.Context, merID, id uint) error {
	return r.db.WithContext(ctx).Where("service_reply_id = ? AND mer_id = ?", id, merID).Delete(&cs.Reply{}).Error
}
