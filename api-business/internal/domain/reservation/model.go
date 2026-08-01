package reservation

import "time"

// ProductTypeReservation 对齐 CRMEB DEFINE_TYPE_RESERVATION = 4（商品 type 字段）
const ProductTypeReservation uint8 = 4

type Config struct {
	ProductReservationID uint      `gorm:"column:product_reservation_id;primaryKey" json:"product_reservation_id"`
	ProductID            uint      `gorm:"column:product_id" json:"product_id"`
	MerID                uint      `gorm:"column:merchant_id" json:"mer_id"`
	StoreID              uint      `gorm:"column:store_id" json:"store_id"`
	ReservationType      int8      `gorm:"column:reservation_type" json:"reservation_type"`
	ShowReservationDays  int       `gorm:"column:show_reservation_days" json:"show_reservation_days"`
	IsCancelReservation  int8      `gorm:"column:is_cancel_reservation" json:"is_cancel_reservation"`
	TimePeriod           string    `gorm:"column:time_period" json:"time_period"`
	CreateTime           time.Time `gorm:"column:create_time" json:"create_time"`
	Status               int8      `gorm:"column:status" json:"status"`
}

func (Config) TableName() string { return "qixi_crm_b_reservation_activity" }

type Slot struct {
	AttrReservationID uint   `gorm:"column:attr_reservation_id;primaryKey" json:"attr_reservation_id"`
	ProductID         uint   `gorm:"column:product_id" json:"product_id"`
	Unique            string `gorm:"column:slot_key" json:"unique"`
	StartTime         string `gorm:"column:start_time" json:"start_time"`
	EndTime           string `gorm:"column:end_time" json:"end_time"`
	Stock             int    `gorm:"column:stock" json:"stock"`
	UseNum            int    `gorm:"column:use_num" json:"use_num"`
}

func (Slot) TableName() string { return "qixi_crm_b_reservation_slot" }

// Booking 是预约时段的已占用账本。创建预约订单时必须在交易完成前落账，
// 日历余量只能基于本表计算，不能再扫描历史订单。
type Booking struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	ProductID uint      `gorm:"column:product_id" json:"product_id"`
	SlotID    uint      `gorm:"column:slot_id" json:"slot_id"`
	Date      string    `gorm:"column:booking_date" json:"date"`
	OrderID   uint      `gorm:"column:order_id" json:"order_id"`
	UID       uint      `gorm:"column:user_id" json:"user_id"`
	Status    int8      `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Booking) TableName() string { return "qixi_crm_b_reservation_booking" }

type ProductView struct {
	ProductID   uint    `json:"product_id"`
	MerID       uint    `json:"mer_id"`
	StoreID     uint    `json:"store_id"`
	MerName     string  `json:"mer_name"`
	StoreName   string  `json:"store_name"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	OtPrice     float64 `json:"ot_price"`
	Stock       uint    `json:"stock"`
	ShowDays    int     `json:"show_reservation_days"`
	ReserveType int8    `json:"reservation_type"`
}

type SlotDayView struct {
	AttrReservationID uint   `json:"attr_reservation_id"`
	StartTime         string `json:"start_time"`
	EndTime           string `json:"end_time"`
	Stock             int    `json:"stock"`
	Booked            int64  `json:"booked"`
	Remain            int    `json:"remain"`
	Label             string `json:"label"`
}

type ConfigSaveInput struct {
	ProductID           uint   `json:"product_id"`
	ReservationType     int8   `json:"reservation_type"`
	ShowReservationDays int    `json:"show_reservation_days"`
	Slots               []Slot `json:"slots"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
