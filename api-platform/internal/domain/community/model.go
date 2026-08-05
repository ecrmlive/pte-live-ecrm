package community

import "time"

// 审核：0 待审 / 1 通过 / -1 驳回
const (
	StatusPending  int8 = 0
	StatusApproved int8 = 1
	StatusRejected int8 = -1
)

type Category struct {
	CategoryID uint   `gorm:"column:category_id;primaryKey" json:"category_id"`
	CateName   string `gorm:"column:cate_name" json:"cate_name"`
	PID        int    `gorm:"column:pid" json:"pid"`
	IsShow     int8   `gorm:"column:is_show" json:"is_show"`
	Sort       int    `gorm:"column:sort" json:"sort"`
}

func (Category) TableName() string { return "qixi_crm_b_social_category" }

type Topic struct {
	TopicID    uint      `gorm:"column:topic_id;primaryKey" json:"topic_id"`
	TopicName  string    `gorm:"column:topic_name" json:"topic_name"`
	Status     int8      `gorm:"column:status" json:"status"`
	IsHot      int8      `gorm:"column:is_hot" json:"is_hot"`
	CategoryID uint      `gorm:"column:category_id" json:"category_id"`
	IsDel      int8      `gorm:"column:is_del" json:"-"`
	CountUse   int       `gorm:"column:count_use" json:"count_use"`
	Sort       int       `gorm:"column:sort" json:"sort"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Topic) TableName() string { return "qixi_crm_b_social_topic" }

type Post struct {
	CommunityID uint       `gorm:"column:community_id;primaryKey" json:"community_id"`
	Title       string     `gorm:"column:title" json:"title"`
	Image       string     `gorm:"column:image" json:"image"`
	CategoryID  uint       `gorm:"column:category_id" json:"category_id"`
	TopicID     uint       `gorm:"column:topic_id" json:"topic_id"`
	UID         uint       `gorm:"column:uid" json:"uid"`
	MerID       uint       `gorm:"column:mer_id" json:"mer_id"`
	ProductID   uint       `gorm:"column:product_id" json:"product_id"`
	CountStart  int        `gorm:"column:count_start" json:"count_start"`
	CountReply  int        `gorm:"column:count_reply" json:"count_reply"`
	Status      int8       `gorm:"column:status" json:"status"`
	IsShow      int8       `gorm:"column:is_show" json:"is_show"`
	IsHot       int8       `gorm:"column:is_hot" json:"is_hot"`
	IsType      int8       `gorm:"column:is_type" json:"is_type"`
	Content     string     `gorm:"column:content" json:"content"`
	Refusal     string     `gorm:"column:refusal" json:"refusal"`
	PV          int        `gorm:"column:pv" json:"pv"`
	IsDel       int8       `gorm:"column:is_del" json:"-"`
	CreateTime  time.Time  `gorm:"column:create_time" json:"create_time"`
	StatusTime  *time.Time `gorm:"column:status_time" json:"status_time,omitempty"`

	Nickname     string  `gorm:"-" json:"nickname,omitempty"`
	TopicName    string  `gorm:"-" json:"topic_name,omitempty"`
	CateName     string  `gorm:"-" json:"cate_name,omitempty"`
	ProductName  string  `gorm:"-" json:"product_name,omitempty"`
	ProductPrice float64 `gorm:"-" json:"product_price,omitempty"`
}

func (Post) TableName() string { return "qixi_crm_b_social_post" }

type Reply struct {
	ReplyID     uint      `gorm:"column:reply_id;primaryKey" json:"reply_id"`
	Content     string    `gorm:"column:content" json:"content"`
	PID         uint      `gorm:"column:pid" json:"pid"`
	UID         uint      `gorm:"column:uid" json:"uid"`
	CommunityID uint      `gorm:"column:community_id" json:"community_id"`
	Status      int8      `gorm:"column:status" json:"status"`
	IsDel       int8      `gorm:"column:is_del" json:"-"`
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`

	Nickname  string `gorm:"-" json:"nickname,omitempty"`
	PostTitle string `gorm:"-" json:"post_title,omitempty"`
}

func (Reply) TableName() string { return "qixi_crm_b_social_reply" }

type CreatePostInput struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Image      string `json:"image"`
	CategoryID uint   `json:"category_id"`
	TopicID    uint   `json:"topic_id"`
	ProductID  uint   `json:"product_id"`
	MerID      uint   `json:"mer_id"`
}

type CreateReplyInput struct {
	Content string `json:"content"`
	PID     uint   `json:"pid"`
}

type AuditInput struct {
	Status  int8   `json:"status"` // 1通过 -1驳回 0=仅调整已通过帖的展示/置顶
	Refusal string `json:"refusal"`
	IsShow  *int   `json:"is_show"`
	IsHot   *int   `json:"is_hot"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
