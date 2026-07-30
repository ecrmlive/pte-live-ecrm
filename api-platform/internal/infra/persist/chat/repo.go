package chat

import (
	"context"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/chat"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) FindThreadByMerUID(ctx context.Context, merID, uid uint) (*chat.Thread, error) {
	var row chat.Thread
	err := r.db.WithContext(ctx).Where("mer_id = ? AND uid = ? AND is_del = 0", merID, uid).First(&row).Error
	return &row, err
}

func (r *Repo) GetThread(ctx context.Context, id uint) (*chat.Thread, error) {
	var row chat.Thread
	err := r.db.WithContext(ctx).Where("thread_id = ?", id).First(&row).Error
	return &row, err
}

func (r *Repo) CreateThread(ctx context.Context, t *chat.Thread) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *Repo) UpdateThread(ctx context.Context, t *chat.Thread) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *Repo) ListThreadsByMer(ctx context.Context, merID uint, page, limit int) ([]chat.Thread, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&chat.Thread{}).Where("mer_id = ? AND is_del = 0", merID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []chat.Thread
	err := q.Order("IFNULL(last_time, create_time) DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListThreadsByUID(ctx context.Context, uid uint, page, limit int) ([]chat.Thread, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&chat.Thread{}).Where("uid = ? AND is_del = 0", uid)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []chat.Thread
	err := q.Order("IFNULL(last_time, create_time) DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CreateMessage(ctx context.Context, m *chat.Message) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *Repo) ListMessages(ctx context.Context, threadID uint, page, limit int) ([]chat.Message, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&chat.Message{}).Where("thread_id = ?", threadID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []chat.Message
	err := q.Order("msg_id ASC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) EnsureIdentity(ctx context.Context, portal string, localID uint, imUserID string, imUserNum int64) (*chat.Identity, error) {
	var row chat.Identity
	err := r.db.WithContext(ctx).Where("portal = ? AND local_id = ?", portal, localID).First(&row).Error
	if err == nil {
		changed := false
		if imUserID != "" && row.ImUserID != imUserID {
			row.ImUserID = imUserID
			changed = true
		}
		if imUserNum > 0 && row.ImUserNum != imUserNum {
			row.ImUserNum = imUserNum
			changed = true
		}
		if changed {
			_ = r.db.WithContext(ctx).Save(&row).Error
		}
		return &row, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	row = chat.Identity{
		Portal: portal, LocalID: localID, ImUserID: imUserID, ImUserNum: imUserNum, CreateTime: time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) FindCustomerServiceID(ctx context.Context, merID uint) (uint, error) {
	var id uint
	err := r.db.WithContext(ctx).Table("qixi_m_admin_store_service").
		Select("service_id").
		Where("mer_id = ? AND is_del = 0 AND status = 1 AND is_open = 1 AND customer = 1", merID).
		Order("service_id ASC").Limit(1).Scan(&id).Error
	return id, err
}

func (r *Repo) LoadUserNickname(ctx context.Context, uid uint) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("qixi_m_app_user").Select("nickname").Where("uid = ?", uid).Scan(&name).Error
	return name, err
}

func (r *Repo) LoadMerName(ctx context.Context, merID uint) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("qixi_m_admin_merchant").Select("mer_name").Where("mer_id = ?", merID).Scan(&name).Error
	return name, err
}
