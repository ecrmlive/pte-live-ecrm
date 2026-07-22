package diy

import (
	"encoding/json"
	"time"
)

// PageValue 本仓库轻量装修协议（非 CRMEB 可视化组件全量）。
type PageValue struct {
	Banners []BannerItem `json:"banners"`
	Menus   []MenuItem   `json:"menus"`
}

type BannerItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	URL   string `json:"url"`
}

type MenuItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	URL  string `json:"url"`
}

type Page struct {
	ID           uint      `gorm:"column:id;primaryKey" json:"id"`
	Version      string    `gorm:"column:version" json:"version"`
	Name         string    `gorm:"column:name" json:"name"`
	Title        string    `gorm:"column:title" json:"title"`
	CoverImage   string    `gorm:"column:cover_image" json:"cover_image"`
	TemplateName string    `gorm:"column:template_name" json:"template_name"`
	Status       int8      `gorm:"column:status" json:"status"`
	Type         int8      `gorm:"column:type" json:"type"`
	IsShow       int8      `gorm:"column:is_show" json:"is_show"`
	IsDiy        int8      `gorm:"column:is_diy" json:"is_diy"`
	MerID        uint      `gorm:"column:mer_id" json:"mer_id"`
	IsDefault    int8      `gorm:"column:is_default" json:"is_default"`
	Value        string    `gorm:"column:value" json:"value"`
	IsDel        int8      `gorm:"column:is_del" json:"-"`
	AddTime      time.Time `gorm:"column:add_time" json:"add_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"update_time"`

	Parsed *PageValue `gorm:"-" json:"parsed,omitempty"`
}

func (Page) TableName() string { return "qixi_diy" }

func (p *Page) ParseValue() PageValue {
	var v PageValue
	if p == nil || p.Value == "" {
		return v
	}
	_ = json.Unmarshal([]byte(p.Value), &v)
	return v
}

type SaveInput struct {
	Name         string    `json:"name"`
	Title        string    `json:"title"`
	TemplateName string    `json:"template_name"`
	Value        PageValue `json:"value"`
	Status       *int8     `json:"status"`
}

type PageResult struct {
	List  []Page `json:"list"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}
