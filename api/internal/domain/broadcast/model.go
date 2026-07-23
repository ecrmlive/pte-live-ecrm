package broadcast

import "time"

// LiveStatus 对齐 CRMEB：101 直播中 / 102 未开始 / 103 已结束。
const (
	LiveStatusLiving   int16 = 101
	LiveStatusPending  int16 = 102
	LiveStatusFinished int16 = 103
)

// AuditStatus 平台审核：0 待审 / 2 通过 / -1 驳回。
const (
	AuditPending  int8 = 0
	AuditApproved int8 = 2
	AuditRejected int8 = -1
)

type Room struct {
	BroadcastRoomID uint       `gorm:"column:broadcast_room_id;primaryKey" json:"broadcast_room_id"`
	MerID           uint       `gorm:"column:mer_id" json:"mer_id"`
	Name            string     `gorm:"column:name" json:"name"`
	CoverImg        string     `gorm:"column:cover_img" json:"cover_img"`
	FeedsImg        string     `gorm:"column:feeds_img" json:"feeds_img"`
	PlayURL         string     `gorm:"column:play_url" json:"play_url"`
	PushURL         string     `gorm:"column:push_url" json:"push_url"`
	StartTime       *time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime         *time.Time `gorm:"column:end_time" json:"end_time"`
	AnchorName      string     `gorm:"column:anchor_name" json:"anchor_name"`
	Phone           string     `gorm:"column:phone" json:"phone"`
	Status          int8       `gorm:"column:status" json:"status"`
	LiveStatus      int16      `gorm:"column:live_status" json:"live_status"`
	IsShow          int8       `gorm:"column:is_show" json:"is_show"`
	IsDel           int8       `gorm:"column:is_del" json:"-"`
	CreateTime      time.Time  `gorm:"column:create_time" json:"create_time"`
	Sort            int        `gorm:"column:sort" json:"sort"`
	Star            int        `gorm:"column:star" json:"star"`
	Mark            string     `gorm:"column:mark" json:"mark"`
	Refusal         string     `gorm:"column:refusal" json:"refusal"`

	MerName string      `gorm:"-" json:"mer_name,omitempty"`
	Goods   []RoomGoods `gorm:"-" json:"goods,omitempty"`
}

func (Room) TableName() string { return "qixi_broadcast_room" }

type RoomGoods struct {
	ID              uint  `gorm:"column:id;primaryKey" json:"id"`
	BroadcastRoomID uint  `gorm:"column:broadcast_room_id" json:"broadcast_room_id"`
	ProductID       uint  `gorm:"column:product_id" json:"product_id"`
	OnSale          int8  `gorm:"column:on_sale" json:"on_sale"`
	Sort            int   `gorm:"column:sort" json:"sort"`

	StoreName string  `gorm:"-" json:"store_name,omitempty"`
	Image     string  `gorm:"-" json:"image,omitempty"`
	Price     float64 `gorm:"-" json:"price,omitempty"`
}

func (RoomGoods) TableName() string { return "qixi_broadcast_room_goods" }

type SaveInput struct {
	Name       string `json:"name"`
	CoverImg   string `json:"cover_img"`
	FeedsImg   string `json:"feeds_img"`
	PlayURL    string `json:"play_url"`
	PushURL    string `json:"push_url"`
	AnchorName string `json:"anchor_name"`
	Phone      string `json:"phone"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	LiveStatus *int16 `json:"live_status"`
	IsShow     *int   `json:"is_show"`
	Sort       *int   `json:"sort"`
	Star       *int   `json:"star"`
	Mark       string `json:"mark"`
	ProductIDs []uint `json:"product_ids"`
}

type GoodsInput struct {
	ProductIDs []uint `json:"product_ids"`
}

type AuditInput struct {
	Status  int8   `json:"status"` // 2通过 -1驳回
	Refusal string `json:"refusal"`
	IsShow  *int   `json:"is_show"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
