package content

import "time"

type Notice struct {
	NoticeID   uint      `gorm:"column:notice_id;primaryKey" json:"notice_id"`
	Title      string    `gorm:"column:title" json:"title"`
	Content    string    `gorm:"column:content" json:"content"`
	IsShow     int8      `gorm:"column:is_show" json:"is_show"`
	Sort       int       `gorm:"column:sort" json:"sort"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	IsDel      int8      `gorm:"column:is_del" json:"-"`
}

func (Notice) TableName() string { return "qixi_crm_a_notice" }

type NoticeInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	IsShow  *int8  `json:"is_show"`
	Sort    int    `json:"sort"`
}

// Cache 保存统一后台可维护的协议与安全配置 stub。
type Cache struct {
	Key        string    `gorm:"column:key;primaryKey" json:"key"`
	ExpireTime int       `gorm:"column:expire_time" json:"expire_time"`
	Result     string    `gorm:"column:result" json:"result"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Cache) TableName() string { return "qixi_crm_a_setting_cache" }

type AgreeMeta struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type AgreeSaveInput struct {
	Content string `json:"content"`
}

type AgreeView struct {
	Key     string `json:"key"`
	Label   string `json:"label"`
	Content string `json:"content"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// CacheListItem 价格说明、活动标签等 setting_cache 列表 stub。
type CacheListItem struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Remark  string `json:"remark"`
}
