package seckill

import "time"

type TimeSlot struct {
	SeckillTimeID uint      `gorm:"column:seckill_time_id;primaryKey" json:"seckill_time_id"`
	Title         string    `gorm:"column:title" json:"title"`
	StartTime     int       `gorm:"column:start_time" json:"start_time"`
	EndTime       int       `gorm:"column:end_time" json:"end_time"`
	Status        int8      `gorm:"column:status" json:"status"`
	CreateTime    time.Time `gorm:"column:create_time" json:"create_time"`
	Pic           string    `gorm:"column:pic" json:"pic"`
}

func (TimeSlot) TableName() string { return "qixi_crm_b_seckill_time" }

type Active struct {
	SeckillActiveID uint    `gorm:"column:seckill_active_id;primaryKey" json:"seckill_active_id"`
	Name            string  `gorm:"column:name" json:"name"`
	SeckillTimeIDs  string  `gorm:"column:seckill_time_ids" json:"seckill_time_ids"`
	StartDay        string  `gorm:"column:start_day" json:"start_day"`
	EndDay          string  `gorm:"column:end_day" json:"end_day"`
	MerID           uint    `gorm:"column:mer_id" json:"mer_id"`
	ProductID       uint    `gorm:"column:product_id" json:"product_id"`
	SeckillPrice    float64 `gorm:"column:seckill_price" json:"seckill_price"`
	OncePayCount    int     `gorm:"column:once_pay_count" json:"once_pay_count"`
	AllPayCount     int     `gorm:"column:all_pay_count" json:"all_pay_count"`
	ActiveStatus    int8    `gorm:"column:active_status" json:"active_status"`
	Status          int8    `gorm:"column:status" json:"status"`
	CreateTime      int64   `gorm:"column:create_time" json:"create_time"`
	UpdateTime      int64   `gorm:"column:update_time" json:"update_time"`
	DeleteTime      *int64  `gorm:"column:delete_time" json:"-"`

	StoreName string  `gorm:"-" json:"store_name,omitempty"`
	Image     string  `gorm:"-" json:"image,omitempty"`
	Price     float64 `gorm:"-" json:"price,omitempty"`
	MerName   string  `gorm:"-" json:"mer_name,omitempty"`
	InWindow  bool    `gorm:"-" json:"in_window"`
}

func (Active) TableName() string { return "qixi_crm_b_seckill_active" }

type ActiveInput struct {
	Name           string  `json:"name"`
	SeckillTimeIDs string  `json:"seckill_time_ids"`
	StartDay       string  `json:"start_day"`
	EndDay         string  `json:"end_day"`
	ProductID      uint    `json:"product_id"`
	SeckillPrice   float64 `json:"seckill_price"`
	OncePayCount   int     `json:"once_pay_count"`
	Status         *int8   `json:"status"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
