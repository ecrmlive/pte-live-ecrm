package content

import "time"

type NoticeScopeType string

const (
	NoticeScopeAll           NoticeScopeType = "all"
	NoticeScopeStoreName     NoticeScopeType = "store_name"
	NoticeScopeStoreType     NoticeScopeType = "store_type"
	NoticeScopeStoreCategory NoticeScopeType = "store_category"
)

type Notice struct {
	NoticeID   uint            `gorm:"column:notice_id;primaryKey" json:"notice_id"`
	Title      string          `gorm:"column:title" json:"title"`
	Content    string          `gorm:"column:content" json:"content"`
	IsShow     int8            `gorm:"column:is_show" json:"is_show"`
	Sort       int             `gorm:"column:sort" json:"-"`
	ScopeType  NoticeScopeType `gorm:"column:scope_type" json:"scope_type"`
	ScopeIDs   []uint          `gorm:"-" json:"scope_ids"`
	ScopeItems []NoticeScope   `gorm:"-" json:"scope_items"`
	CreateTime time.Time       `gorm:"column:create_time" json:"create_time"`
	IsDel      int8            `gorm:"column:is_del" json:"-"`
}

func (Notice) TableName() string { return "qixi_crm_a_notice" }

func (NoticeScope) TableName() string { return "qixi_crm_a_notice_scope" }

type NoticeInput struct {
	Title     string          `json:"title"`
	Content   string          `json:"content"`
	IsShow    *int8           `json:"is_show"`
	ScopeType NoticeScopeType `json:"scope_type"`
	ScopeIDs  []uint          `json:"scope_ids"`
}

type NoticeStatusInput struct {
	IsShow int8 `json:"is_show"`
}

type NoticeListFilter struct {
	Page     int
	Limit    int
	Keyword  string
	IsShow   *int8
	DateFrom string
	DateTo   string
}

// NoticeScope 是公告投放范围的关联记录。名称在服务层根据真实关联数据动态补全。
type NoticeScope struct {
	NoticeID  uint            `gorm:"column:notice_id" json:"-"`
	ScopeID   uint            `gorm:"column:scope_id" json:"id"`
	ScopeType NoticeScopeType `gorm:"column:scope_kind" json:"scope_type"`
	Name      string          `gorm:"-" json:"name"`
}

// Cache 保存统一后台可维护的协议与安全配置 stub。
type Cache struct {
	Key        string `gorm:"column:key;primaryKey" json:"key"`
	ExpireTime int    `gorm:"column:expire_time" json:"expire_time"`
	Result     string `gorm:"column:result" json:"result"`
	// autoCreateTime：避免零值写入 MySQL 非法 '0000-00-00'（严格模式 Error 1292）。
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
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
