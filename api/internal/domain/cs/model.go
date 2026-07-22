package cs

import "time"

// Reply 客服快捷回复（qixi_store_service_reply）。
type Reply struct {
	ServiceReplyID uint      `gorm:"column:service_reply_id;primaryKey" json:"service_reply_id"`
	MerID          uint      `gorm:"column:mer_id" json:"mer_id"`
	Type           int8      `gorm:"column:type" json:"type"` // 1 文字 2 图片
	Keyword        string    `gorm:"column:keyword" json:"keyword"`
	Content        string    `gorm:"column:content" json:"content"`
	Status         int8      `gorm:"column:status" json:"status"`
	CreateTime     time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Reply) TableName() string { return "qixi_store_service_reply" }

type ReplyInput struct {
	Type    int8   `json:"type"`
	Keyword string `json:"keyword"`
	Content string `json:"content"`
	Status  *int8  `json:"status"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
