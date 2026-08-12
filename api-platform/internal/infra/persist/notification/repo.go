package notification

import (
	"context"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/notification"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) List(ctx context.Context, audience notification.Audience, page, limit int) ([]notification.Config, int64, error) {
	q := r.db.WithContext(ctx).Model(&notification.Config{}).Where("audience = ?", audience)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []notification.Config
	err := q.Order("notification_id ASC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) Get(ctx context.Context, id uint) (*notification.Config, error) {
	var row notification.Config
	err := r.db.WithContext(ctx).Where("notification_id = ?", id).First(&row).Error
	return &row, err
}

func (r *Repo) Save(ctx context.Context, config *notification.Config) error {
	return r.db.WithContext(ctx).Model(&notification.Config{}).
		Where("notification_id = ?", config.NotificationID).
		Updates(map[string]any{
			"wechat_enabled":       config.WechatEnabled,
			"mini_program_enabled": config.MiniEnabled,
			"sms_enabled":          config.SMSEnabled,
			"wechat_text":          config.WechatText,
			"mini_program_text":    config.MiniText,
			"sms_text":             config.SMSText,
		}).Error
}
