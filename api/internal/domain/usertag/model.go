package usertag

import "time"

type Label struct {
	LabelID    uint      `gorm:"column:label_id;primaryKey" json:"label_id"`
	LabelName  string    `gorm:"column:label_name" json:"label_name"`
	Sort       int       `gorm:"column:sort" json:"sort"`
	IsDel      int8      `gorm:"column:is_del" json:"-"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Label) TableName() string { return "qixi_user_label" }

type Group struct {
	GroupID    uint      `gorm:"column:group_id;primaryKey" json:"group_id"`
	GroupName  string    `gorm:"column:group_name" json:"group_name"`
	Sort       int       `gorm:"column:sort" json:"sort"`
	IsDel      int8      `gorm:"column:is_del" json:"-"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Group) TableName() string { return "qixi_user_group" }

type Relation struct {
	ID      uint `gorm:"column:id;primaryKey" json:"id"`
	UID     uint `gorm:"column:uid" json:"uid"`
	LabelID uint `gorm:"column:label_id" json:"label_id"`
}

func (Relation) TableName() string { return "qixi_user_label_relation" }

type LabelInput struct {
	LabelName string `json:"label_name"`
	Sort      int    `json:"sort"`
}

type GroupInput struct {
	GroupName string `json:"group_name"`
	Sort      int    `json:"sort"`
}

type MarkInput struct {
	UID      uint   `json:"uid"`
	LabelIDs []uint `json:"label_ids"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
