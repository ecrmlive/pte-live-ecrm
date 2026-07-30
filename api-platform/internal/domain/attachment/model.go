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
	CreateTime               time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Category) TableName() string { return "qixi_m_admin_system_attachment_category" }

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

func (Attachment) TableName() string { return "qixi_m_admin_system_attachment" }

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
