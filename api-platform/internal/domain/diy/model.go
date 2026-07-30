package diy

import (
	"encoding/json"
	"strings"
	"time"
)

// PageDoc 可视化装修文档（对齐 pte-live-shop {page,items[]}）。
type PageDoc struct {
	Page  map[string]any   `json:"page"`
	Items []map[string]any `json:"items"`
}

// LegacyValue 阶段6 轻量协议，读取时自动升级。
type LegacyValue struct {
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
	IsBgColor    int8      `gorm:"column:is_bg_color" json:"is_bg_color"`
	IsBgPic      int8      `gorm:"column:is_bg_pic" json:"is_bg_pic"`
	ColorPicker  string    `gorm:"column:color_picker" json:"color_picker"`
	BgPic        string    `gorm:"column:bg_pic" json:"bg_pic"`
	BgTabVal     int8      `gorm:"column:bg_tab_val" json:"bg_tab_val"`
	MerID        uint      `gorm:"column:mer_id" json:"mer_id"`
	IsDefault    int8      `gorm:"column:is_default" json:"is_default"`
	Value        string    `gorm:"column:value" json:"value"`
	IsDel        int8      `gorm:"column:is_del" json:"-"`
	AddTime      time.Time `gorm:"column:add_time" json:"add_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"update_time"`

	Doc *PageDoc `gorm:"-" json:"doc,omitempty"`
}

func (Page) TableName() string { return "qixi_m_admin_diy" }

// PageCategory 是装修链接选择器的分组。is_mer=0 表示平台商城链接，1 表示商户商城链接。
// 它不是商户私有数据：商户端只能读取平台配置的商户链接分组。
type PageCategory struct {
	ID       uint           `gorm:"column:id;primaryKey" json:"id"`
	PID      uint           `gorm:"column:pid" json:"pid"`
	Type     string         `gorm:"column:type" json:"type"`
	Name     string         `gorm:"column:name" json:"name"`
	Sort     int            `gorm:"column:sort" json:"sort"`
	Status   int8           `gorm:"column:status" json:"status"`
	Level    int8           `gorm:"column:level" json:"level"`
	IsMer    int8           `gorm:"column:is_mer" json:"is_mer"`
	AddTime  time.Time      `gorm:"column:add_time" json:"add_time"`
	Children []PageCategory `gorm:"-" json:"children,omitempty"`
}

func (PageCategory) TableName() string { return "qixi_m_admin_page_category" }

// PageLink 是可由装修组件选择的固定页面路径或小程序路径。
type PageLink struct {
	ID       uint          `gorm:"column:id;primaryKey" json:"id"`
	CateID   uint          `gorm:"column:cate_id" json:"cate_id"`
	Type     int8          `gorm:"column:type" json:"type"`
	Name     string        `gorm:"column:name" json:"name"`
	URL      string        `gorm:"column:url" json:"url"`
	Param    string        `gorm:"column:param" json:"param"`
	Example  string        `gorm:"column:example" json:"example"`
	Status   int8          `gorm:"column:status" json:"status"`
	Sort     int           `gorm:"column:sort" json:"sort"`
	IsMer    int8          `gorm:"column:is_mer" json:"is_mer"`
	AddTime  time.Time     `gorm:"column:add_time" json:"add_time"`
	Category *PageCategory `gorm:"-" json:"category,omitempty"`
}

func (PageLink) TableName() string { return "qixi_m_admin_page_link" }

func (p *Page) ParseDoc() PageDoc {
	if p == nil || p.Value == "" {
		return emptyDoc()
	}
	var doc PageDoc
	if err := json.Unmarshal([]byte(p.Value), &doc); err == nil && doc.Items != nil {
		if doc.Page == nil {
			doc.Page = map[string]any{}
		}
		return doc
	}
	var legacy LegacyValue
	if err := json.Unmarshal([]byte(p.Value), &legacy); err == nil &&
		(len(legacy.Banners) > 0 || len(legacy.Menus) > 0 || looksLegacy(p.Value)) {
		return legacyToDoc(legacy, p.Name, p.Title)
	}
	return emptyDoc()
}

func looksLegacy(raw string) bool {
	return (strings.Contains(raw, `"banners"`) || strings.Contains(raw, `"menus"`)) &&
		!strings.Contains(raw, `"items"`)
}

func emptyDoc() PageDoc {
	return PageDoc{Page: map[string]any{}, Items: []map[string]any{}}
}

func legacyToDoc(v LegacyValue, name, title string) PageDoc {
	bannerData := make([]map[string]any, 0, len(v.Banners))
	for _, b := range v.Banners {
		bannerData = append(bannerData, map[string]any{
			"imgUrl": b.Image, "linkUrl": b.URL, "imgName": b.Title,
		})
	}
	navData := make([]map[string]any, 0, len(v.Menus))
	for _, m := range v.Menus {
		navData = append(navData, map[string]any{
			"text": m.Name, "imgUrl": m.Icon, "linkUrl": m.URL,
		})
	}
	items := []map[string]any{}
	if len(bannerData) > 0 {
		items = append(items, map[string]any{
			"type": "banner", "name": "轮播图", "params": map[string]any{},
			"style": map[string]any{"btnColor": "#ffffff", "btnShape": "round", "indicator": "1"},
			"data":  bannerData,
		})
	}
	if len(navData) > 0 {
		items = append(items, map[string]any{
			"type": "navBar", "name": "导航组", "params": map[string]any{},
			"style": map[string]any{"background": "#ffffff", "rowsNum": "4", "show_title": "1"},
			"data":  navData,
		})
	}
	if name == "" {
		name = "首页"
	}
	if title == "" {
		title = name
	}
	return PageDoc{
		Page: map[string]any{
			"type": "page", "name": "页面设置",
			"params": map[string]any{"name": name, "title": title, "share_title": title},
			"style":  map[string]any{"titleTextColor": "black", "titleBackgroundColor": "#ffffff"},
		},
		Items: items,
	}
}

type SaveInput struct {
	Name         string          `json:"name"`
	Title        string          `json:"title"`
	TemplateName string          `json:"template_name"`
	CoverImage   string          `json:"cover_image"`
	IsDiy        *int8           `json:"is_diy"`
	Type         *int8           `json:"type"`
	IsShow       *int8           `json:"is_show"`
	IsBgColor    *int8           `json:"is_bg_color"`
	IsBgPic      *int8           `json:"is_bg_pic"`
	ColorPicker  string          `json:"color_picker"`
	BgPic        string          `json:"bg_pic"`
	BgTabVal     *int8           `json:"bg_tab_val"`
	Doc          *PageDoc        `json:"doc"`
	Value        json.RawMessage `json:"value"`
	Status       *int8           `json:"status"`
}

type ListFilter struct {
	MerID  uint
	IsDiy  *int8
	Page   int
	Limit  int
	Name   string
	Status *int8
}

type PageResult struct {
	List  []Page `json:"list"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type CategoryInput struct {
	PID    uint   `json:"pid"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Sort   int    `json:"sort"`
	Status *int8  `json:"status"`
}

type LinkInput struct {
	CateID  uint   `json:"cate_id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Param   string `json:"param"`
	Example string `json:"example"`
	Sort    int    `json:"sort"`
	Status  *int8  `json:"status"`
}

type LinkListFilter struct {
	IsMer  int8
	Status *int8
	Page   int
	Limit  int
	Name   string
}

type LinkResult struct {
	List  []PageLink `json:"list"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}

type EditorBootstrap struct {
	DefaultData map[string]any `json:"defaultData"`
	DefaultPage map[string]any `json:"defaultPage"`
	JSONData    PageDoc        `json:"jsonData"`
	Opts        map[string]any `json:"opts"`
}
