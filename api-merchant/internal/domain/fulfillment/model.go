package fulfillment

import "time"

type Staff struct {
	StaffID    uint      `gorm:"column:staff_id;primaryKey" json:"staff_id"`
	MerID      uint      `gorm:"column:mer_id" json:"mer_id"`
	Name       string    `gorm:"column:name" json:"name"`
	Phone      string    `gorm:"column:phone" json:"phone"`
	Status     int8      `gorm:"column:status" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	IsDel      int8      `gorm:"column:is_del" json:"-"`
}

func (Staff) TableName() string { return "qixi_m_admin_delivery_staff" }

type StaffInput struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Status *int8  `json:"status"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
