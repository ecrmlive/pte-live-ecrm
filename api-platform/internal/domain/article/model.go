package article

import "time"

type Category struct {
	CID    uint   `gorm:"column:cid;primaryKey" json:"cid"`
	Title  string `gorm:"column:title" json:"title"`
	Status int8   `gorm:"column:status" json:"status"`
	Sort   int    `gorm:"column:sort" json:"sort"`
	IsDel  int8   `gorm:"column:is_del" json:"-"`
}

func (Category) TableName() string { return "qixi_crm_a_article_category" }

type Article struct {
	ArticleID  uint      `gorm:"column:article_id;primaryKey" json:"article_id"`
	CID        uint      `gorm:"column:cid" json:"cid"`
	Title      string    `gorm:"column:title" json:"title"`
	Author     string    `gorm:"column:author" json:"author"`
	Image      string    `gorm:"column:image" json:"image"`
	Synopsis   string    `gorm:"column:synopsis" json:"synopsis"`
	Content    string    `gorm:"column:content" json:"content"`
	Visit      int       `gorm:"column:visit" json:"visit"`
	Sort       int       `gorm:"column:sort" json:"sort"`
	Status     int8      `gorm:"column:status" json:"status"`
	IsDel      int8      `gorm:"column:is_del" json:"-"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Article) TableName() string { return "qixi_crm_a_article" }

type CategoryInput struct {
	Title  string `json:"title"`
	Status *int8  `json:"status"`
	Sort   int    `json:"sort"`
}

type ArticleInput struct {
	CID      uint   `json:"cid"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Image    string `json:"image"`
	Synopsis string `json:"synopsis"`
	Content  string `json:"content"`
	Sort     int    `json:"sort"`
	Status   *int8  `json:"status"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
