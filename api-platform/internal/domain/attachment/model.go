package attachment

import "time"

type Category struct {
	AttachmentCategoryID     uint      `gorm:"column:attachment_category_id;primaryKey" json:"attachment_category_id"`
	PID                      uint      `gorm:"column:pid" json:"pid"`
	Path                     string    `gorm:"column:path" json:"path"`
	AttachmentCategoryName   string    `gorm:"column:attachment_category_name" json:"attachment_category_name"`
	AttachmentCategoryEnname string    `gorm:"column:attachment_category_enname" json:"attachment_category_enname"`
	Sort                     int       `gorm:"column:sort" json:"sort"`
	MerID                    uint      `gorm:"column:mer_id" json:"mer_id"`
	IsSystem                 int8      `gorm:"column:is_system" json:"is_system"`
	CreateTime               time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Category) TableName() string { return "qixi_crm_a_attachment_category" }

type Attachment struct {
	AttachmentID         uint      `gorm:"column:attachment_id;primaryKey" json:"attachment_id"`
	AttachmentCategoryID uint      `gorm:"column:attachment_category_id" json:"attachment_category_id"`
	AttachmentName       string    `gorm:"column:attachment_name" json:"attachment_name"`
	AttachmentSrc        string    `gorm:"column:attachment_src" json:"attachment_src"`
	UploadType           int8      `gorm:"column:upload_type" json:"upload_type"`
	UserType             int       `gorm:"column:user_type" json:"user_type"`
	UserID               uint      `gorm:"column:user_id" json:"user_id"`
	CreateTime           time.Time `gorm:"column:create_time" json:"create_time"`
	AttachmentType       int8      `gorm:"column:attachment_type" json:"attachment_type"`
}

func (Attachment) TableName() string { return "qixi_crm_a_attachment_asset" }

type CategoryInput struct {
	Name   string `json:"attachment_category_name"`
	EnName string `json:"attachment_category_enname"`
	PID    uint   `json:"pid"`
	Sort   int    `json:"sort"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// SystemCategorySpec 平台素材库固定系统分类（侧栏「全部素材」为虚拟入口，不落库）。
// 用途：客户端（H5/小程序/App）与装修页常用的图标、图片、背景。
type SystemCategorySpec struct {
	EnName string
	Name   string
	Sort   int
}

var SystemCategories = []SystemCategorySpec{
	{EnName: "store_cover", Name: "店铺封面", Sort: 90},
	{EnName: "pay_icon", Name: "支付图标", Sort: 80},
	{EnName: "logistics_icon", Name: "物流图标", Sort: 70},
	{EnName: "service_icon", Name: "客服图标", Sort: 60},
	{EnName: "product_image", Name: "商品图片", Sort: 50},
	{EnName: "background_image", Name: "背景图片", Sort: 40},
	{EnName: "list_icon", Name: "列表图标", Sort: 30},
	{EnName: "other_image", Name: "其他图片", Sort: 20},
}

func IsSystemCategoryEnname(en string) bool {
	for _, spec := range SystemCategories {
		if spec.EnName == en {
			return true
		}
	}
	return false
}
