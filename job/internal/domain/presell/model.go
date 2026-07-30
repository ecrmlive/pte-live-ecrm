package presell

import "time"

// ProductPresell 预售活动：presell_type 1全款 / 2定金。
type ProductPresell struct {
	ProductPresellID uint      `gorm:"column:product_presell_id;primaryKey" json:"product_presell_id"`
	StartTime        time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime          time.Time `gorm:"column:end_time" json:"end_time"`
	FinalStartTime   string    `gorm:"column:final_start_time" json:"final_start_time"`
	FinalEndTime     string    `gorm:"column:final_end_time" json:"final_end_time"`
	Status           int       `gorm:"column:status" json:"status"`
	PresellType      int       `gorm:"column:presell_type" json:"presell_type"`
	PayCount         int       `gorm:"column:pay_count" json:"pay_count"`
	DeliveryType     int       `gorm:"column:delivery_type" json:"delivery_type"`
	DeliveryDay      int       `gorm:"column:delivery_day" json:"delivery_day"`
	ProductID        uint      `gorm:"column:product_id" json:"product_id"`
	Price            float64   `gorm:"column:price" json:"price"`
	DownPrice        float64   `gorm:"column:down_price" json:"down_price"`
	FinalPrice       float64   `gorm:"column:final_price" json:"final_price"`
	Stock            int       `gorm:"column:stock" json:"stock"`
	IsShow           int       `gorm:"column:is_show" json:"is_show"`
	StoreName        string    `gorm:"column:store_name" json:"store_name"`
	MerID            uint      `gorm:"column:mer_id" json:"mer_id"`
	StoreInfo        string    `gorm:"column:store_info" json:"store_info"`
	IsDel            int       `gorm:"column:is_del" json:"-"`
	CreateTime       time.Time `gorm:"column:create_time" json:"create_time"`
	ProductStatus    int       `gorm:"column:product_status" json:"product_status"`
	Refusal          string    `gorm:"column:refusal" json:"refusal"`
	ActionStatus     int       `gorm:"column:action_status" json:"action_status"`
	Seles            int       `gorm:"column:seles" json:"seles"`

	Image    string  `gorm:"-" json:"image,omitempty"`
	OtPrice  float64 `gorm:"-" json:"ot_price,omitempty"`
	MerName  string  `gorm:"-" json:"mer_name,omitempty"`
	InWindow bool    `gorm:"-" json:"in_window"`
}

func (ProductPresell) TableName() string { return "qixi_m_admin_store_product_presell" }

// Presell 与 ProductPresell 同义（service 命名）。
type Presell = ProductPresell

type SaveInput struct {
	ProductID      uint    `json:"product_id"`
	StoreName      string  `json:"store_name"`
	StoreInfo      string  `json:"store_info"`
	Price          float64 `json:"price"`
	DownPrice      float64 `json:"down_price"`
	FinalPrice     float64 `json:"final_price"`
	PresellType    int     `json:"presell_type"` // 1全款 2定金
	Stock          int     `json:"stock"`
	PayCount       int     `json:"pay_count"`
	DeliveryType   int     `json:"delivery_type"`
	DeliveryDay    int     `json:"delivery_day"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	FinalStartTime string  `json:"final_start_time"`
	FinalEndTime   string  `json:"final_end_time"`
	IsShow         *int    `json:"is_show"`
	Status         *int    `json:"status"`
}

// PresellOrder 定金预售尾款单。
type PresellOrder struct {
	PresellOrderID   uint       `gorm:"column:presell_order_id;primaryKey" json:"presell_order_id"`
	PresellOrderSN   string     `gorm:"column:presell_order_sn" json:"presell_order_sn"`
	UID              uint       `gorm:"column:uid" json:"uid"`
	MerID            uint       `gorm:"column:mer_id" json:"mer_id"`
	OrderID          uint       `gorm:"column:order_id" json:"order_id"`
	ProductPresellID uint       `gorm:"column:product_presell_id" json:"product_presell_id"`
	FinalStartTime   time.Time  `gorm:"column:final_start_time" json:"final_start_time"`
	FinalEndTime     time.Time  `gorm:"column:final_end_time" json:"final_end_time"`
	Paid             int8       `gorm:"column:paid" json:"paid"`
	Status           int8       `gorm:"column:status" json:"status"`
	PayType          int8       `gorm:"column:pay_type" json:"pay_type"`
	PayPrice         float64    `gorm:"column:pay_price" json:"pay_price"`
	PayTime          *time.Time `gorm:"column:pay_time" json:"pay_time,omitempty"`
	CreateTime       time.Time  `gorm:"column:create_time" json:"create_time"`

	StoreName string `gorm:"-" json:"store_name,omitempty"`
	OrderSN   string `gorm:"-" json:"order_sn,omitempty"`
}

func (PresellOrder) TableName() string { return "qixi_m_app_presell_order" }

// Input 兼容 handler。
type Input = SaveInput

// TradeQuote 下单核价结果（供 trade.PresellHook）。
type TradeQuote struct {
	ProductPresellID uint
	ProductID        uint
	MerID            uint
	StoreName        string
	Image            string
	MerName          string
	Price            float64
	DownPrice        float64
	FinalPrice       float64
	PayUnit          float64 // 首期单价：全款=price，定金=down_price
	Stock            int
	PresellType      int
	FinalStartTime   string
	FinalEndTime     string
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
