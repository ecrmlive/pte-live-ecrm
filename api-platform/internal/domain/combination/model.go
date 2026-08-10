package combination

import "time"

type ProductGroup struct {
	ProductGroupID uint      `gorm:"column:product_group_id;primaryKey" json:"product_group_id"`
	ProductID      uint      `gorm:"column:product_id" json:"product_id"`
	StartTime      time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime        time.Time `gorm:"column:end_time" json:"end_time"`
	Time           int       `gorm:"column:time" json:"time"` // 开团时长小时
	BuyingCountNum int       `gorm:"column:buying_count_num" json:"buying_count_num"`
	BuyingNum      int       `gorm:"column:buying_num" json:"buying_num"`
	PayCount       int       `gorm:"column:pay_count" json:"pay_count"`
	OncePayCount   int       `gorm:"column:once_pay_count" json:"once_pay_count"`
	Status         int       `gorm:"column:status" json:"status"`
	MerID          uint      `gorm:"column:mer_id" json:"mer_id"`
	IsShow         int       `gorm:"column:is_show" json:"is_show"`
	IsDel          int       `gorm:"column:is_del" json:"-"`
	SuccessNum     int       `gorm:"column:success_num" json:"success_num"`
	ProductStatus  int       `gorm:"column:product_status" json:"product_status"`
	Refusal        string    `gorm:"column:refusal" json:"refusal"`
	Price          float64   `gorm:"column:price" json:"price"`
	ActionStatus   int       `gorm:"column:action_status" json:"action_status"`
	CreateTime     time.Time `gorm:"column:create_time" json:"create_time"`

	StoreName string  `gorm:"-" json:"store_name,omitempty"`
	Image     string  `gorm:"-" json:"image,omitempty"`
	OtPrice   float64 `gorm:"-" json:"ot_price,omitempty"`
	MerName   string  `gorm:"-" json:"mer_name,omitempty"`
}

func (ProductGroup) TableName() string { return "qixi_crm_b_combination_group" }

// 开团状态：0 未完成 / 10 已完成 / -1 已失败（对齐 CRMEB ProductGroupBuying）
const (
	BuyingStatusRunning = 0
	BuyingStatusDone    = 10
	BuyingStatusFailed  = -1
)

type Buying struct {
	GroupBuyingID  uint      `gorm:"column:group_buying_id;primaryKey" json:"group_buying_id"`
	ProductGroupID uint      `gorm:"column:product_group_id" json:"product_group_id"`
	Status         int       `gorm:"column:status" json:"status"`
	BuyingCountNum int       `gorm:"column:buying_count_num" json:"buying_count_num"`
	BuyingNum      int       `gorm:"column:buying_num" json:"buying_num"`
	YetBuyingNum   int       `gorm:"column:yet_buying_num" json:"yet_buying_num"`
	IsDel          int       `gorm:"column:is_del" json:"-"`
	MerID          uint      `gorm:"column:mer_id" json:"mer_id"`
	EndTime        int64     `gorm:"column:end_time" json:"end_time"`
	CreateTime     time.Time `gorm:"column:create_time" json:"create_time"`

	Members        []Member `gorm:"-" json:"members,omitempty"`
	Remain         int      `gorm:"-" json:"remain"`
	ProductID      uint     `gorm:"-" json:"product_id,omitempty"`
	StoreName      string   `gorm:"-" json:"store_name,omitempty"`
	Image          string   `gorm:"-" json:"image,omitempty"`
	Price          float64  `gorm:"-" json:"price,omitempty"`
	MerName        string   `gorm:"-" json:"mer_name,omitempty"`
	IsTrader       int8     `gorm:"-" json:"is_trader"`
	TraderName     string   `gorm:"-" json:"trader_name,omitempty"`
	Nickname       string   `gorm:"-" json:"nickname,omitempty"`
	UID            uint     `gorm:"-" json:"uid,omitempty"`
	StopTime       string   `gorm:"-" json:"stop_time,omitempty"`
	StatusText     string   `gorm:"-" json:"status_text,omitempty"`
}

// AdminBuyingQuery 平台「拼团活动列表」开团记录筛选（对齐 CRMEB StoreProductGroupBuying）。
type AdminBuyingQuery struct {
	Page     int
	Limit    int
	MerID    *uint
	MerIDs   []uint // is_trader 反查得到的店铺 ID 集合
	Keyword  string // 商品名称 / 商品 ID
	UserName string // 团长昵称 / UID
	DateFrom string
	DateTo   string
	Status   *int
}

func (Buying) TableName() string { return "qixi_crm_b_combination_buying" }

type Member struct {
	ID             uint      `gorm:"column:id;primaryKey" json:"id"`
	GroupBuyingID  uint      `gorm:"column:group_buying_id" json:"group_buying_id"`
	ProductGroupID uint      `gorm:"column:product_group_id" json:"product_group_id"`
	Status         int       `gorm:"column:status" json:"status"`
	IsInitiator    int       `gorm:"column:is_initiator" json:"is_initiator"`
	OrderID        uint      `gorm:"column:order_id" json:"order_id"`
	UID            uint      `gorm:"column:uid" json:"uid"`
	Nickname       string    `gorm:"column:nickname" json:"nickname"`
	Avatar         string    `gorm:"column:avatar" json:"avatar"`
	IsDel          int       `gorm:"column:is_del" json:"-"`
	CreateTime     time.Time `gorm:"column:create_time" json:"create_time"`
	IsLeader       int8      `gorm:"column:is_leader" json:"is_leader"`
}

func (Member) TableName() string { return "qixi_crm_b_combination_member" }

type SaveInput struct {
	ProductID      uint    `json:"product_id"`
	Price          float64 `json:"price"`
	BuyingCountNum int     `json:"buying_count_num"`
	Time           int     `json:"time"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	IsShow         *int    `json:"is_show"`
	Status         *int    `json:"status"`
	ProductStatus  *int    `json:"product_status"`
	Refusal        string  `json:"refusal"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
