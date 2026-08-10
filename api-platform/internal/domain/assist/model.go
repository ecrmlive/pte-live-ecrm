package assist

import "time"

// Set 状态：1 进行中 / 10 已完成可下单 / 20 已支付 / -1 失败
const (
	SetStatusRunning = 1
	SetStatusDone    = 10
	SetStatusPaid    = 20
	SetStatusFailed  = -1
)

type ProductAssist struct {
	ProductAssistID uint      `gorm:"column:product_assist_id;primaryKey" json:"product_assist_id"`
	StartTime       time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime         time.Time `gorm:"column:end_time" json:"end_time"`
	Status          int       `gorm:"column:status" json:"status"`
	PayCount        int       `gorm:"column:pay_count" json:"pay_count"`
	AssistCount     int       `gorm:"column:assist_count" json:"assist_count"`
	AssistUserCount int       `gorm:"column:assist_user_count" json:"assist_user_count"`
	ProductID       uint      `gorm:"column:product_id" json:"product_id"`
	AssistPrice     float64   `gorm:"column:assist_price" json:"assist_price"`
	Stock           int       `gorm:"column:stock" json:"stock"`
	IsShow          int8      `gorm:"column:is_show" json:"is_show"`
	StoreName       string    `gorm:"column:store_name" json:"store_name"`
	MerID           uint      `gorm:"column:mer_id" json:"mer_id"`
	StoreInfo       string    `gorm:"column:store_info" json:"store_info"`
	IsDel           int       `gorm:"column:is_del" json:"-"`
	CreateTime      time.Time `gorm:"column:create_time" json:"create_time"`
	ProductStatus   int       `gorm:"column:product_status" json:"product_status"`
	Refusal         string    `gorm:"column:refusal" json:"refusal"`
	ActionStatus    int       `gorm:"column:action_status" json:"action_status"`

	Image             string  `gorm:"-" json:"image,omitempty"`
	OtPrice           float64 `gorm:"-" json:"ot_price,omitempty"`
	MerName           string  `gorm:"-" json:"mer_name,omitempty"`
	IsTrader          int8    `gorm:"-" json:"is_trader"`
	TraderName        string  `gorm:"-" json:"trader_name,omitempty"`
	StockCount        int     `gorm:"-" json:"stock_count"`         // 限量（原始）≈ stock + pay
	Success           int     `gorm:"-" json:"success"`             // 助力成功人次（set status 10/20）
	All               int     `gorm:"-" json:"all"`                 // 参与助力人次（helpers）
	Pay               int     `gorm:"-" json:"pay"`                 // 已售/已支付人次（set status 20）
	AssistStatus      int8    `gorm:"-" json:"assist_status"`       // 0未开始 1进行中 2已结束
	AssistStatusText  string  `gorm:"-" json:"assist_status_text,omitempty"`
	ProductStatusName string  `gorm:"-" json:"product_status_name,omitempty"`
}

// 助力活动是 C 端业务事实，统一存放在 business 库；不得回退到 CRMEB
// 旧运行表。
func (ProductAssist) TableName() string { return "qixi_crm_b_assist" }

type AssistSet struct {
	ProductAssistSetID uint      `gorm:"column:product_assist_set_id;primaryKey" json:"product_assist_set_id"`
	ProductAssistID    uint      `gorm:"column:product_assist_id" json:"product_assist_id"`
	ProductID          uint      `gorm:"column:product_id" json:"product_id"`
	UID                uint      `gorm:"column:uid" json:"uid"`
	Status             int       `gorm:"column:status" json:"status"`
	AssistCount        int       `gorm:"column:assist_count" json:"assist_count"`
	AssistUserCount    int       `gorm:"column:assist_user_count" json:"assist_user_count"`
	YetAssistCount     int       `gorm:"column:yet_assist_count" json:"yet_assist_count"`
	CreateTime         time.Time `gorm:"column:create_time" json:"create_time"`
	MerID              uint      `gorm:"column:mer_id" json:"mer_id"`
	IsDel              int       `gorm:"column:is_del" json:"-"`

	StoreName   string       `gorm:"-" json:"store_name,omitempty"`
	AssistPrice float64      `gorm:"-" json:"assist_price,omitempty"`
	Nickname    string       `gorm:"-" json:"nickname,omitempty"`
	Image       string       `gorm:"-" json:"image,omitempty"`
	MerName     string       `gorm:"-" json:"mer_name,omitempty"`
	IsTrader    int8         `gorm:"-" json:"is_trader"`
	TraderName  string       `gorm:"-" json:"trader_name,omitempty"`
	StartTime   time.Time    `gorm:"-" json:"start_time,omitempty"`
	EndTime     time.Time    `gorm:"-" json:"end_time,omitempty"`
	Helpers     []AssistUser `gorm:"-" json:"helpers,omitempty"`
}

// AdminSetQuery 平台「助力活动」实例列表筛选（对齐 CRMEB StoreProductAssistSet）。
type AdminSetQuery struct {
	Page     int
	Limit    int
	MerID    *uint
	MerIDs   []uint // is_trader 反查得到的店铺 ID 集合
	Keyword  string // 商品名称 / 商品 ID / 实例 ID
	UserName string // 发起人昵称
	DateFrom string
	DateTo   string
	Status   *int
}

func (AssistSet) TableName() string { return "qixi_crm_b_assist_set" }

type AssistUser struct {
	ProductAssistUserID uint      `gorm:"column:product_assist_user_id;primaryKey" json:"product_assist_user_id"`
	ProductAssistSetID  uint      `gorm:"column:product_assist_set_id" json:"product_assist_set_id"`
	ProductAssistID     uint      `gorm:"column:product_assist_id" json:"product_assist_id"`
	UID                 uint      `gorm:"column:uid" json:"uid"`
	Nickname            string    `gorm:"column:nickname" json:"nickname"`
	AvatarImg           string    `gorm:"column:avatar_img" json:"avatar_img"`
	CreateTime          time.Time `gorm:"column:create_time" json:"create_time"`
}

func (AssistUser) TableName() string { return "qixi_crm_b_assist_user" }

type SaveInput struct {
	ProductID       uint    `json:"product_id"`
	StoreName       string  `json:"store_name"`
	StoreInfo       string  `json:"store_info"`
	AssistPrice     float64 `json:"assist_price"`
	AssistCount     int     `json:"assist_count"`
	AssistUserCount int     `json:"assist_user_count"`
	Stock           int     `json:"stock"`
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
	IsShow          *int    `json:"is_show"`
	Status          *int    `json:"status"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
