package communitypersist

import (
	"context"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/community"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListCategories(ctx context.Context, onlyShow bool) ([]community.Category, error) {
	q := r.db.WithContext(ctx).Model(&community.Category{})
	if onlyShow {
		q = q.Where("is_show = 1")
	}
	var rows []community.Category
	err := q.Order("sort ASC, category_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) GetCategory(ctx context.Context, id uint) (*community.Category, error) {
	var row community.Category
	err := r.db.WithContext(ctx).Where("category_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateCategory(ctx context.Context, row *community.Category) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateCategory(ctx context.Context, row *community.Category) error {
	return r.db.WithContext(ctx).Model(row).Where("category_id = ?", row.CategoryID).Updates(map[string]interface{}{
		"cate_name": row.CateName,
		"pid":       row.PID,
		"is_show":   row.IsShow,
		"sort":      row.Sort,
	}).Error
}

func (r *Repo) UpdateCategoryShow(ctx context.Context, id uint, isShow int8) error {
	return r.db.WithContext(ctx).Model(&community.Category{}).
		Where("category_id = ?", id).
		Update("is_show", isShow).Error
}

func (r *Repo) DeleteCategory(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("category_id = ?", id).Delete(&community.Category{}).Error
}

func (r *Repo) CountCategoryUsage(ctx context.Context, id uint) (posts int64, topics int64, err error) {
	if err = r.db.WithContext(ctx).Model(&community.Post{}).
		Where("category_id = ? AND is_del = 0", id).Count(&posts).Error; err != nil {
		return
	}
	if err = r.db.WithContext(ctx).Model(&community.Topic{}).
		Where("category_id = ? AND is_del = 0", id).Count(&topics).Error; err != nil {
		return
	}
	return
}

func (r *Repo) ListTopics(ctx context.Context, onlyOn bool) ([]community.Topic, error) {
	q := r.db.WithContext(ctx).Model(&community.Topic{}).Where("is_del = 0")
	if onlyOn {
		q = q.Where("status = 1")
	}
	var rows []community.Topic
	err := q.Order("is_hot DESC, sort ASC, topic_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) GetTopic(ctx context.Context, id uint) (*community.Topic, error) {
	var row community.Topic
	err := r.db.WithContext(ctx).Where("topic_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateTopic(ctx context.Context, row *community.Topic) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateTopic(ctx context.Context, row *community.Topic) error {
	return r.db.WithContext(ctx).Model(row).Where("topic_id = ? AND is_del = 0", row.TopicID).Updates(map[string]interface{}{
		"topic_name":  row.TopicName,
		"pic":         row.Pic,
		"category_id": row.CategoryID,
		"sort":        row.Sort,
		"status":      row.Status,
		"is_hot":      row.IsHot,
	}).Error
}

func (r *Repo) UpdateTopicStatus(ctx context.Context, id uint, status int8) error {
	return r.db.WithContext(ctx).Model(&community.Topic{}).
		Where("topic_id = ? AND is_del = 0", id).
		Update("status", status).Error
}

func (r *Repo) UpdateTopicHot(ctx context.Context, id uint, isHot int8) error {
	return r.db.WithContext(ctx).Model(&community.Topic{}).
		Where("topic_id = ? AND is_del = 0", id).
		Update("is_hot", isHot).Error
}

func (r *Repo) SoftDeleteTopic(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&community.Topic{}).
		Where("topic_id = ?", id).
		Update("is_del", 1).Error
}

func (r *Repo) TopicNameExists(ctx context.Context, name string, excludeID uint) (bool, error) {
	q := r.db.WithContext(ctx).Model(&community.Topic{}).
		Where("topic_name = ? AND is_del = 0", name)
	if excludeID > 0 {
		q = q.Where("topic_id <> ?", excludeID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *Repo) applyPostFilter(q *gorm.DB, f community.ListFilter) *gorm.DB {
	q = q.Where("is_del = 0")
	if f.MerID != nil {
		q = q.Where("mer_id = ?", *f.MerID)
	}
	if f.UID != nil {
		q = q.Where("uid = ?", *f.UID)
	}
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.Keyword != "" {
		if f.TitleOnly {
			q = q.Where("title LIKE ?", "%"+f.Keyword+"%")
		} else {
			q = q.Where("title LIKE ? OR content LIKE ?", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
		}
	}
	if f.TopicID > 0 {
		q = q.Where("topic_id = ?", f.TopicID)
	}
	if f.CategoryID > 0 {
		q = q.Where("category_id = ?", f.CategoryID)
	}
	if f.IsShow != nil {
		q = q.Where("is_show = ?", *f.IsShow)
	}
	if f.IsType != nil {
		if *f.IsType == community.TypeImage {
			q = q.Where("is_type IN ?", []int8{0, community.TypeImage})
		} else {
			q = q.Where("is_type = ?", *f.IsType)
		}
	}
	if kw := f.AuthorKW; kw != "" {
		switch f.AuthorType {
		case "uid":
			q = q.Where("uid = ?", kw)
		case "phone":
			q = q.Where("uid IN (SELECT id FROM qixi_crm_b_user WHERE mobile LIKE ?)", "%"+kw+"%")
		default: // nickname
			q = q.Where("uid IN (SELECT id FROM qixi_crm_b_user WHERE nickname LIKE ?)", "%"+kw+"%")
		}
	}
	if f.OnlyPublic {
		q = q.Where("status = ? AND is_show = 1", community.StatusApproved)
	}
	return q
}

func (r *Repo) ListPosts(ctx context.Context, f community.ListFilter) ([]community.Post, int64, error) {
	page, limit := f.Page, f.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	q := r.applyPostFilter(r.db.WithContext(ctx).Model(&community.Post{}), f)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []community.Post
	err := q.Order("`start` DESC, is_hot DESC, community_id DESC").
		Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CountPosts(ctx context.Context, f community.ListFilter) (int64, error) {
	q := r.applyPostFilter(r.db.WithContext(ctx).Model(&community.Post{}), f)
	var total int64
	err := q.Count(&total).Error
	return total, err
}

func (r *Repo) GetPost(ctx context.Context, id uint) (*community.Post, error) {
	var row community.Post
	err := r.db.WithContext(ctx).Where("community_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreatePost(ctx context.Context, p *community.Post) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repo) UpdatePost(ctx context.Context, p *community.Post) error {
	return r.db.WithContext(ctx).Model(p).Where("community_id = ?", p.CommunityID).Updates(map[string]interface{}{
		"title": p.Title, "image": p.Image, "category_id": p.CategoryID, "topic_id": p.TopicID,
		"product_id": p.ProductID, "mer_id": p.MerID, "status": p.Status, "is_show": p.IsShow,
		"is_hot": p.IsHot, "start": p.Start, "is_type": p.IsType, "video_link": p.VideoLink,
		"content": p.Content, "refusal": p.Refusal, "status_time": p.StatusTime,
	}).Error
}

func (r *Repo) SoftDeletePost(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&community.Post{}).Where("community_id = ?", id).
		Updates(map[string]interface{}{"is_del": 1, "is_show": 0}).Error
}

func (r *Repo) IncPV(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&community.Post{}).Where("community_id = ?", id).
		Update("pv", gorm.Expr("pv + 1")).Error
}

func (r *Repo) IncReplyCount(ctx context.Context, id uint, delta int) error {
	if delta == 0 {
		return nil
	}
	expr := gorm.Expr("count_reply + ?", delta)
	if delta < 0 {
		expr = gorm.Expr("CASE WHEN count_reply > 0 THEN count_reply - 1 ELSE 0 END")
	}
	return r.db.WithContext(ctx).Model(&community.Post{}).Where("community_id = ?", id).
		Update("count_reply", expr).Error
}

func (r *Repo) ListAllReplies(ctx context.Context, f community.ReplyListFilter) ([]community.Reply, int64, error) {
	page, limit := f.Page, f.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Model(&community.Reply{}).Where("is_del = 0")
	if f.Keyword != "" {
		q = q.Where("content LIKE ?", "%"+f.Keyword+"%")
	}
	if f.Username != "" {
		q = q.Where("uid IN (SELECT id FROM qixi_crm_b_user WHERE nickname LIKE ?)", "%"+f.Username+"%")
	}
	if f.DateFrom != "" {
		q = q.Where("create_time >= ?", f.DateFrom+" 00:00:00")
	}
	if f.DateTo != "" {
		q = q.Where("create_time <= ?", f.DateTo+" 23:59:59")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []community.Reply
	err := q.Order("reply_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListReplies(ctx context.Context, communityID uint, page, limit int) ([]community.Reply, int64, error) {
	q := r.db.WithContext(ctx).Model(&community.Reply{}).
		Where("community_id = ? AND is_del = 0", communityID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []community.Reply
	err := q.Order("reply_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CreateReply(ctx context.Context, row *community.Reply) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateReply(ctx context.Context, row *community.Reply) error {
	return r.db.WithContext(ctx).Model(row).Where("reply_id = ?", row.ReplyID).Updates(map[string]interface{}{
		"status": row.Status, "refusal": row.Refusal,
	}).Error
}

func (r *Repo) SoftDeleteReply(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&community.Reply{}).Where("reply_id = ?", id).
		Update("is_del", 1).Error
}

func (r *Repo) GetReply(ctx context.Context, id uint) (*community.Reply, error) {
	var row community.Reply
	err := r.db.WithContext(ctx).Where("reply_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) LoadUserNickname(ctx context.Context, uid uint) (string, error) {
	var nick string
	err := r.db.WithContext(ctx).Table("qixi_crm_b_user").Select("nickname").Where("id = ?", uid).Scan(&nick).Error
	if err != nil {
		return "", err
	}
	if nick == "" {
		nick = "用户"
	}
	return nick, nil
}

func (r *Repo) LoadTopicName(ctx context.Context, id uint) (string, error) {
	if id == 0 {
		return "", nil
	}
	var name string
	err := r.db.WithContext(ctx).Table("qixi_crm_b_social_topic").Select("topic_name").
		Where("topic_id = ? AND is_del = 0", id).Scan(&name).Error
	return name, err
}

func (r *Repo) LoadCateName(ctx context.Context, id uint) (string, error) {
	if id == 0 {
		return "", nil
	}
	var name string
	err := r.db.WithContext(ctx).Table("qixi_crm_b_social_category").Select("cate_name").
		Where("category_id = ?", id).Scan(&name).Error
	return name, err
}

func (r *Repo) LoadProductMeta(ctx context.Context, productID uint) (name string, price float64, merID uint, err error) {
	var row struct {
		StoreName string  `gorm:"column:store_name"`
		Price     float64 `gorm:"column:price"`
		MerID     uint    `gorm:"column:mer_id"`
	}
	err = r.db.WithContext(ctx).Table("qixi_crm_b_product_view").
		Select("store_name, price, merchant_id AS mer_id").
		Where("product_id = ? AND is_del = 0", productID).
		Limit(1).Scan(&row).Error
	if err != nil {
		return "", 0, 0, err
	}
	if row.StoreName == "" && row.MerID == 0 {
		return "", 0, 0, gorm.ErrRecordNotFound
	}
	return row.StoreName, row.Price, row.MerID, nil
}

var _ community.Store = (*Repo)(nil)
