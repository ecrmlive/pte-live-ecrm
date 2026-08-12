package contentpersist

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/content"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListNotices(ctx context.Context, onlyShow bool, filter content.NoticeListFilter) ([]content.Notice, int64, error) {
	q := r.db.WithContext(ctx).Model(&content.Notice{}).Where("is_del = 0")
	if onlyShow {
		q = q.Where("is_show = 1")
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		q = q.Where("title LIKE ?", "%"+keyword+"%")
	}
	if filter.IsShow != nil {
		q = q.Where("is_show = ?", *filter.IsShow)
	}
	if from := strings.TrimSpace(filter.DateFrom); from != "" {
		q = q.Where("create_time >= ?", from)
	}
	if to := strings.TrimSpace(filter.DateTo); to != "" {
		q = q.Where("create_time <= ?", to)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []content.Notice
	err := q.Order("create_time DESC, notice_id DESC").Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).Find(&rows).Error
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

func (r *Repo) CreateNotice(ctx context.Context, n *content.Notice, scopeIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(n).Error; err != nil {
			return err
		}
		return replaceNoticeScopes(ctx, tx, n.NoticeID, n.ScopeType, scopeIDs)
	})
}

func (r *Repo) UpdateNotice(ctx context.Context, n *content.Notice, scopeIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(n).Where("notice_id = ?", n.NoticeID).Updates(map[string]interface{}{
			"title":      n.Title,
			"content":    n.Content,
			"is_show":    n.IsShow,
			"scope_type": n.ScopeType,
		}).Error; err != nil {
			return err
		}
		return replaceNoticeScopes(ctx, tx, n.NoticeID, n.ScopeType, scopeIDs)
	})
}

func (r *Repo) UpdateNoticeStatus(ctx context.Context, id uint, isShow int8) error {
	return r.db.WithContext(ctx).Model(&content.Notice{}).
		Where("notice_id = ? AND is_del = 0", id).Update("is_show", isShow).Error
}

func (r *Repo) ListNoticeScopes(ctx context.Context, noticeIDs []uint) ([]content.NoticeScope, error) {
	if len(noticeIDs) == 0 {
		return []content.NoticeScope{}, nil
	}
	var rows []content.NoticeScope
	if err := r.db.WithContext(ctx).Where("notice_id IN ?", noticeIDs).
		Order("notice_id ASC, scope_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	idsByType := map[content.NoticeScopeType][]uint{}
	for _, row := range rows {
		idsByType[row.ScopeType] = append(idsByType[row.ScopeType], row.ScopeID)
	}
	names := make(map[content.NoticeScopeType]map[uint]string)
	for scopeType, ids := range idsByType {
		names[scopeType] = lookupScopeNames(ctx, r.db, scopeType, ids)
	}
	for i := range rows {
		rows[i].Name = names[rows[i].ScopeType][rows[i].ScopeID]
	}
	return rows, nil
}

func replaceNoticeScopes(ctx context.Context, db *gorm.DB, noticeID uint, scopeType content.NoticeScopeType, scopeIDs []uint) error {
	if err := db.WithContext(ctx).Where("notice_id = ?", noticeID).Delete(&content.NoticeScope{}).Error; err != nil {
		return err
	}
	if scopeType == content.NoticeScopeAll {
		return nil
	}
	if !scopeIDsExist(ctx, db, scopeType, scopeIDs) {
		return content.ErrBadParam
	}
	rows := make([]content.NoticeScope, 0, len(scopeIDs))
	for _, id := range scopeIDs {
		rows = append(rows, content.NoticeScope{NoticeID: noticeID, ScopeID: id, ScopeType: scopeType})
	}
	return db.WithContext(ctx).Create(&rows).Error
}

func scopeIDsExist(ctx context.Context, db *gorm.DB, scopeType content.NoticeScopeType, ids []uint) bool {
	if len(ids) == 0 {
		return false
	}
	var count int64
	var result *gorm.DB
	switch scopeType {
	case content.NoticeScopeStoreName:
		result = db.WithContext(ctx).Table("qixi_crm_a_merchant_view").Where("merchant_id IN ?", ids).Count(&count)
	case content.NoticeScopeStoreType:
		result = db.WithContext(ctx).Table("qixi_crm_a_merchant_type").Where("id IN ?", ids).Count(&count)
	case content.NoticeScopeStoreCategory:
		result = db.WithContext(ctx).Table("qixi_crm_a_merchant_category").Where("id IN ?", ids).Count(&count)
	default:
		return false
	}
	return result.Error == nil && count == int64(len(ids))
}

func lookupScopeNames(ctx context.Context, db *gorm.DB, scopeType content.NoticeScopeType, ids []uint) map[uint]string {
	result := make(map[uint]string, len(ids))
	if len(ids) == 0 {
		return result
	}
	var rows []struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var query *gorm.DB
	switch scopeType {
	case content.NoticeScopeStoreName:
		query = db.WithContext(ctx).Table("qixi_crm_a_merchant_view").Select("merchant_id AS id, merchant_name AS name").Where("merchant_id IN ?", ids)
	case content.NoticeScopeStoreType:
		query = db.WithContext(ctx).Table("qixi_crm_a_merchant_type").Select("id, name").Where("id IN ?", ids)
	case content.NoticeScopeStoreCategory:
		query = db.WithContext(ctx).Table("qixi_crm_a_merchant_category").Select("id, name").Where("id IN ?", ids)
	default:
		return result
	}
	if query.Find(&rows).Error != nil {
		return result
	}
	for _, row := range rows {
		result[row.ID] = row.Name
	}
	return result
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
			// 首次写入必须带合法 create_time；零值会被 GORM 写成 '0000-00-00' 触发 MySQL 1292。
			if row.CreateTime.IsZero() {
				row.CreateTime = time.Now()
			}
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
