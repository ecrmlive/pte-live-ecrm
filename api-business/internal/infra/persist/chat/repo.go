package chat

import (
	"context"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/chat"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) threadQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Table("qixi_crm_b_customer_service_binding AS b").
		Select("b.*, s.merchant_id AS mer_id").
		Joins("JOIN qixi_crm_b_store_view AS s ON s.store_id = b.store_id")
}

func (r *Repo) FindThreadByMerUID(ctx context.Context, merID, uid uint) (*chat.Thread, error) {
	var row chat.Thread
	err := r.threadQuery(ctx).Where("s.merchant_id = ? AND b.user_id = ?", merID, uid).Order("b.id DESC").Take(&row).Error
	return &row, err
}

func (r *Repo) GetThread(ctx context.Context, id uint) (*chat.Thread, error) {
	var row chat.Thread
	err := r.threadQuery(ctx).Where("b.id = ?", id).Take(&row).Error
	return &row, err
}

func (r *Repo) CreateThread(ctx context.Context, t *chat.Thread) error {
	if t.StoreID == 0 {
		var storeID uint
		if err := r.db.WithContext(ctx).Table("qixi_crm_b_store_view").
			Select("store_id").Where("merchant_id = ? AND status = ?", t.MerID, 1).Order("store_id ASC").Take(&storeID).Error; err != nil {
			return err
		}
		t.StoreID = storeID
	}
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *Repo) UpdateThread(ctx context.Context, t *chat.Thread) error {
	return r.db.WithContext(ctx).Model(&chat.Thread{}).Where("id = ?", t.ThreadID).Updates(map[string]any{
		"assigned_admin_id": t.ServiceID, "assigned_at": t.AssignedAt, "im_conversation_id": t.ImConversationID,
		"last_msg": t.LastMsg, "last_time": t.LastTime, "user_unread": t.UserUnread,
		"service_unread": t.ServiceUnread, "status": t.Status,
	}).Error
}

func (r *Repo) ListThreadsByMer(ctx context.Context, merID uint, page, limit int) ([]chat.Thread, int64, error) {
	q := r.threadQuery(ctx).Where("s.merchant_id = ?", merID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []chat.Thread
	err := q.Order("COALESCE(b.last_time, b.created_at) DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	return rows, total, err
}

func (r *Repo) ListThreadsByUID(ctx context.Context, uid uint, page, limit int) ([]chat.Thread, int64, error) {
	q := r.threadQuery(ctx).Where("b.user_id = ?", uid)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []chat.Thread
	err := q.Order("COALESCE(b.last_time, b.created_at) DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	return rows, total, err
}

func (r *Repo) CreateMessage(ctx context.Context, m *chat.Message) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *Repo) ListMessages(ctx context.Context, threadID uint, page, limit int) ([]chat.Message, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&chat.Message{}).Where("binding_id = ?", threadID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []chat.Message
	err := q.Order("id ASC").Offset((page - 1) * limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) EnsureIdentity(ctx context.Context, portal string, localID uint, imUserID string, imUserNum int64) (*chat.Identity, error) {
	var row chat.Identity
	err := r.db.WithContext(ctx).Where("portal = ? AND local_id = ?", portal, localID).First(&row).Error
	if err == nil {
		if row.ImUserID != imUserID || row.ImUserNum != imUserNum {
			if err := r.db.WithContext(ctx).Model(&row).Updates(map[string]any{"im_user_id": imUserID, "im_user_num": imUserNum}).Error; err != nil {
				return nil, err
			}
			row.ImUserID, row.ImUserNum = imUserID, imUserNum
		}
		return &row, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	row = chat.Identity{Portal: portal, LocalID: localID, ImUserID: imUserID, ImUserNum: imUserNum, CreateTime: time.Now()}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) FindCustomerServiceID(ctx context.Context, merID uint) (uint, error) {
	var id uint
	err := r.db.WithContext(ctx).Table("qixi_crm_b_customer_service_agent_view AS a").
		Joins("JOIN qixi_crm_b_store_view AS s ON s.store_id = a.store_id").
		Select("a.admin_id").Where("s.merchant_id = ? AND s.status = ? AND a.status = ?", merID, 1, 1).
		Order("a.available_at ASC, a.admin_id ASC").Limit(1).Scan(&id).Error
	return id, err
}

func (r *Repo) LoadUserNickname(ctx context.Context, uid uint) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("qixi_crm_b_user").Select("nickname").Where("id = ?", uid).Scan(&name).Error
	return name, err
}

func (r *Repo) LoadMerName(ctx context.Context, merID uint) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("qixi_crm_b_store_view").Select("store_name").Where("merchant_id = ? AND status = ?", merID, 1).Order("store_id ASC").Scan(&name).Error
	return name, err
}

func (r *Repo) LoadMerchantIMConfig(ctx context.Context, merID uint) (*chat.MerchantIMConfig, error) {
	var row chat.MerchantIMConfig
	err := r.db.WithContext(ctx).Table("qixi_crm_b_merchant_im_sdk_app_view").
		Select("sdk_app_id, api_public_url, ws_public_url").Where("merchant_id = ?", merID).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}
