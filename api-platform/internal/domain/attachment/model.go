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
	// IsSystem=1：客户端/装修用的系统预置素材；挂在系统分类下的运营图应为 0。
	IsSystem int8 `gorm:"column:is_system" json:"is_system"`
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
// Kind: 0=图片 1=视频。分类 is_system=1 仅表示固定类目不可增删改；
// 侧栏「系统素材」还要求素材行 is_system=1，运营图可挂在此类目下但行级为 0。
type SystemCategorySpec struct {
	EnName string
	Name   string
	Sort   int
	Kind   int8 // 0 image, 1 video
}

var SystemCategories = []SystemCategorySpec{
	{EnName: "store_cover", Name: "店铺封面", Sort: 90, Kind: 0},
	{EnName: "pay_icon", Name: "支付图标", Sort: 80, Kind: 0},
	{EnName: "logistics_icon", Name: "物流图标", Sort: 70, Kind: 0},
	{EnName: "service_icon", Name: "客服图标", Sort: 60, Kind: 0},
	{EnName: "product_image", Name: "商品图片", Sort: 50, Kind: 0},
	{EnName: "background_image", Name: "背景图片", Sort: 40, Kind: 0},
	{EnName: "list_icon", Name: "列表图标", Sort: 30, Kind: 0},
	{EnName: "other_image", Name: "其他图片", Sort: 20, Kind: 0},
	{EnName: "store_video", Name: "店铺视频", Sort: 19, Kind: 1},
	{EnName: "product_video", Name: "商品视频", Sort: 18, Kind: 1},
	{EnName: "other_video", Name: "其他视频", Sort: 17, Kind: 1},
}

func IsSystemCategoryEnname(en string) bool {
	for _, spec := range SystemCategories {
		if spec.EnName == en {
			return true
		}
	}
	return false
}

// SystemCategoryEnnames 返回系统分类 enname；mediaType 为 nil 时返回全部，0=图片，1=视频。
func SystemCategoryEnnames(mediaType *int8) []string {
	out := make([]string, 0, len(SystemCategories))
	for _, spec := range SystemCategories {
		if mediaType != nil && spec.Kind != *mediaType {
			continue
		}
		out = append(out, spec.EnName)
	}
	return out
}
