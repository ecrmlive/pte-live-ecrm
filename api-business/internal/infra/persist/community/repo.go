package communitypersist

import (
	"context"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/community"
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

func (r *Repo) ListTopics(ctx context.Context, onlyOn bool) ([]community.Topic, error) {
	q := r.db.WithContext(ctx).Model(&community.Topic{}).Where("is_del = 0")
	if onlyOn {
		q = q.Where("status = 1")
	}
	var rows []community.Topic
	err := q.Order("is_hot DESC, sort ASC, topic_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) ListPosts(ctx context.Context, f community.ListFilter) ([]community.Post, int64, error) {
	page, limit := f.Page, f.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Model(&community.Post{}).Where("is_del = 0")
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
		q = q.Where("title LIKE ? OR content LIKE ?", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
	}
	if f.TopicID > 0 {
		q = q.Where("topic_id = ?", f.TopicID)
	}
	if f.OnlyPublic {
		q = q.Where("status = ? AND is_show = 1", community.StatusApproved)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []community.Post
	err := q.Order("is_hot DESC, community_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
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
		"is_hot": p.IsHot, "content": p.Content, "refusal": p.Refusal, "status_time": p.StatusTime,
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

func (r *Repo) ListReplies(ctx context.Context, communityID uint, page, limit int) ([]community.Reply, int64, error) {
	q := r.db.WithContext(ctx).Model(&community.Reply{}).
		Where("community_id = ? AND is_del = 0 AND status = 1", communityID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []community.Reply
	err := q.Order("reply_id ASC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CreateReply(ctx context.Context, row *community.Reply) error {
	return r.db.WithContext(ctx).Create(row).Error
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
		MerID     uint    `gorm:"column:merchant_id"`
	}
	err = r.db.WithContext(ctx).Table("qixi_crm_b_product_view").
		Select("store_name, price, merchant_id").
		Where("product_id = ? AND sale_status = 1", productID).
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
