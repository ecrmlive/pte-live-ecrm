package promotion

import "time"

const (
	CouponTypeStore    int = 0  // 店铺券
	CouponTypeProduct  int = 1  // 商品券
	CouponTypePlatform int = 10 // 平台通用券
)

const (
	TemplateDays  int = 0 // 领取后 N 天有效
	TemplateFixed int = 1 // 固定时段
)

const (
	UserUnused  int = 0
	UserUsed    int = 1
	UserExpired int = 2
)

const (
	BillPMOut int8 = 0
	BillPMIn  int8 = 1
)

// 演示一级佣金比例（支付金额 × 5%）
const SpreadRate = 0.05

type Coupon struct {
	CouponID      uint       `gorm:"column:coupon_id;primaryKey" json:"coupon_id"`
	MerID         uint       `gorm:"column:mer_id" json:"mer_id"`
	IsTimeout     int8       `gorm:"column:is_timeout" json:"is_timeout"`
	StartTime     *time.Time `gorm:"column:start_time" json:"start_time,omitempty"`
	EndTime       *time.Time `gorm:"column:end_time" json:"end_time,omitempty"`
	IsLimited     int8       `gorm:"column:is_limited" json:"is_limited"`
	TotalCount    uint       `gorm:"column:total_count" json:"total_count"`
	RemainCount   uint       `gorm:"column:remain_count" json:"remain_count"`
	SendType      int8       `gorm:"column:send_type" json:"send_type"`
	FullReduction float64    `gorm:"column:full_reduction" json:"full_reduction"`
	Title         string     `gorm:"column:title" json:"title"`
	CouponPrice   float64    `gorm:"column:coupon_price" json:"coupon_price"`
	UseMinPrice   int        `gorm:"column:use_min_price" json:"use_min_price"`
	CouponType    int8       `gorm:"column:coupon_type" json:"coupon_type"`
	CouponTime    uint       `gorm:"column:coupon_time" json:"coupon_time"`
	UseStartTime  *time.Time `gorm:"column:use_start_time" json:"use_start_time,omitempty"`
	UseEndTime    *time.Time `gorm:"column:use_end_time" json:"use_end_time,omitempty"`
	Sort          uint       `gorm:"column:sort" json:"sort"`
	Status        int8       `gorm:"column:status" json:"status"`
	CreateTime    time.Time  `gorm:"column:create_time" json:"create_time"`
	IsDel         int8       `gorm:"column:is_del" json:"-"`
	Type          int        `gorm:"column:type" json:"type"`

	Received bool `gorm:"-" json:"received,omitempty"`
}

func (Coupon) TableName() string { return "qixi_crm_b_store_coupon" }

type CouponUser struct {
	CouponUserID uint       `gorm:"column:coupon_user_id;primaryKey" json:"coupon_user_id"`
	CouponID     uint       `gorm:"column:coupon_id" json:"coupon_id"`
	MerID        uint       `gorm:"column:mer_id" json:"mer_id"`
	UID          uint       `gorm:"column:uid" json:"uid"`
	CouponTitle  string     `gorm:"column:coupon_title" json:"coupon_title"`
	CouponPrice  float64    `gorm:"column:coupon_price" json:"coupon_price"`
	UseMinPrice  int        `gorm:"column:use_min_price" json:"use_min_price"`
	CreateTime   time.Time  `gorm:"column:create_time" json:"create_time"`
	StartTime    *time.Time `gorm:"column:start_time" json:"start_time,omitempty"`
	EndTime      *time.Time `gorm:"column:end_time" json:"end_time,omitempty"`
	UseTime      *time.Time `gorm:"column:use_time" json:"use_time,omitempty"`
	Type         string     `gorm:"column:type" json:"type"`
	SendID       uint       `gorm:"column:send_id" json:"send_id"`
	Status       int        `gorm:"column:status" json:"status"`
	IsFail       int8       `gorm:"column:is_fail" json:"is_fail"`

	CouponKind int `gorm:"-" json:"coupon_kind,omitempty"` // 模板 type：0/10
}

func (CouponUser) TableName() string { return "qixi_m_app_store_coupon_user" }

// CouponSend 商户后台定向发券的批次记录。
type CouponSend struct {
	CouponSendID uint      `gorm:"column:coupon_send_id;primaryKey" json:"coupon_send_id"`
	MerID        uint      `gorm:"column:mer_id" json:"mer_id"`
	CouponID     uint      `gorm:"column:coupon_id" json:"coupon_id"`
	CouponNum    uint      `gorm:"column:coupon_num" json:"coupon_num"`
	Mark         string    `gorm:"column:mark" json:"mark"`
	CreateTime   time.Time `gorm:"column:create_time" json:"create_time"`
	Status       int8      `gorm:"column:status" json:"status"`
}

func (CouponSend) TableName() string { return "qixi_crm_b_store_coupon_send" }

// CouponSendListFilter 平台发送记录筛选。
type CouponSendListFilter struct {
	DateFrom    string
	DateTo      string
	CouponType  *int // 模板 type：10 平台通用券等
	CouponName  string
	SendStatus  *int8
}

// CouponSendListItem 发送记录列表行（对齐 CRMEB systemCouponSendLst）。
type CouponSendListItem struct {
	CouponSendID   uint       `json:"coupon_send_id"`
	CouponID       uint       `json:"coupon_id"`
	Title          string     `json:"title"`
	Type           int        `json:"type"`
	CouponTypeName string     `json:"coupon_type_name"`
	CreateTime     time.Time  `json:"create_time"`
	CouponType     int8       `json:"coupon_type"`
	CouponTime     uint       `json:"coupon_time"`
	UseStartTime   *time.Time `json:"use_start_time,omitempty"`
	UseEndTime     *time.Time `json:"use_end_time,omitempty"`
	ValidityText   string     `json:"validity_text"`
	Mark           string     `json:"mark"`
	FilterText     string     `json:"filter_text"`
	CouponNum      uint       `json:"coupon_num"`
	UseCount       int64      `json:"use_count"`
	SendStatus     int8       `json:"send_status"`
}

// CouponSendDetail 发送记录详情。
type CouponSendDetail struct {
	CouponSendListItem
	CouponPrice   float64    `json:"coupon_price"`
	UseMinPrice   int        `json:"use_min_price"`
	IsTimeout     int8       `json:"is_timeout"`
	StartTime     *time.Time `json:"start_time,omitempty"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	SendType      int8       `json:"send_type"`
	SendTypeName  string     `json:"send_type_name"`
	IsLimited     int8       `json:"is_limited"`
	TotalCount    uint       `json:"total_count"`
	RemainCount   uint       `json:"remain_count"`
	Status        int8       `json:"status"`
	Sort          uint       `json:"sort"`
	SentTotal     uint       `json:"sent_total"`
	UsedTotal     int64      `json:"used_total"`
}

// CouponSendUserRow 某次发送的使用记录行。
type CouponSendUserRow struct {
	UserID     uint64 `json:"user_id"`
	Nickname   string `json:"nickname"`
	AvatarURL  string `json:"avatar_url"`
	Source     string `json:"source"`
	SourceName string `json:"source_name"`
	Status     string `json:"status"`
	StatusName string `json:"status_name"`
}

type IssueUser struct {
	UID        uint      `gorm:"column:uid;primaryKey" json:"uid"`
	CouponID   uint      `gorm:"column:coupon_id;primaryKey" json:"coupon_id"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (IssueUser) TableName() string { return "qixi_m_app_store_coupon_issue_user" }

type SpreadLog struct {
	UserSpreadLogID uint      `gorm:"column:user_spread_log_id;primaryKey" json:"user_spread_log_id"`
	UID             uint      `gorm:"column:uid" json:"uid"`
	OldSpreadUID    uint      `gorm:"column:old_spread_uid" json:"old_spread_uid"`
	SpreadUID       uint      `gorm:"column:spread_uid" json:"spread_uid"`
	CreateTime      time.Time `gorm:"column:create_time" json:"create_time"`
}

func (SpreadLog) TableName() string { return "qixi_m_app_user_spread_log" }

type UserBill struct {
	BillID     uint      `gorm:"column:bill_id;primaryKey" json:"bill_id"`
	UID        uint      `gorm:"column:uid" json:"uid"`
	LinkID     string    `gorm:"column:link_id" json:"link_id"`
	PM         int8      `gorm:"column:pm" json:"pm"`
	Title      string    `gorm:"column:title" json:"title"`
	Category   string    `gorm:"column:category" json:"category"`
	Type       string    `gorm:"column:type" json:"type"`
	Number     float64   `gorm:"column:number" json:"number"`
	Balance    float64   `gorm:"column:balance" json:"balance"`
	Mark       string    `gorm:"column:mark" json:"mark"`
	MerID      uint      `gorm:"column:mer_id" json:"mer_id"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	Status     int8      `gorm:"column:status" json:"status"`
}

func (UserBill) TableName() string { return "qixi_m_app_user_bill" }

// 对齐 CRMEB StoreCouponRepository 获取方式。
const (
	SendTypeReceive int8 = 0 // 领取
	SendTypeAdmin   int8 = 6 // 后台发放
	SendTypeGift    int8 = 3 // 赠送券
)

type CreateCouponInput struct {
	Title         string  `json:"title"`
	CouponPrice   float64 `json:"coupon_price"`
	UseMinPrice   int     `json:"use_min_price"`
	UseType       *int8   `json:"use_type"` // 0无门槛 1有门槛；nil 时按 use_min_price 推断
	CouponType    *int8   `json:"coupon_type"`
	CouponTime    uint    `json:"coupon_time"`
	UseStartTime  string  `json:"use_start_time"`
	UseEndTime    string  `json:"use_end_time"`
	IsTimeout     *int8   `json:"is_timeout"`
	StartTime     string  `json:"start_time"`
	EndTime       string  `json:"end_time"`
	SendType      *int8   `json:"send_type"`
	FullReduction float64 `json:"full_reduction"`
	TotalCount    uint    `json:"total_count"`
	IsLimited     int8    `json:"is_limited"`
	Sort          uint    `json:"sort"`
	Status        *int8   `json:"status"`
}

// CouponSaveInput 管理端创建/更新入参别名。
type CouponSaveInput = CreateCouponInput

// CouponListFilter 管理端列表筛选。
type CouponListFilter struct {
	Keyword   string
	Status    *int8
	SendType  *int8
	StoreOnly bool   // mer_id > 0（商户优惠券）
	MerIDs    []uint // 店铺类别等二次过滤
}

// CouponDetail 管理端优惠券详情（含领取/使用统计）。
type CouponDetail struct {
	Coupon
	ReceivedTotal  int64  `json:"received_total"`
	UsedTotal      int64  `json:"used_total"`
	UseType        int8   `json:"use_type"`
	MerName        string `json:"mer_name,omitempty"`
	IsTrader       int8   `json:"is_trader,omitempty"`
	TraderName     string `json:"trader_name,omitempty"`
	CouponTypeName string `json:"coupon_type_name,omitempty"`
	ClaimText      string `json:"claim_text,omitempty"`
	ValidityText   string `json:"validity_text,omitempty"`
}

// StoreCouponListItem 平台「商户优惠券」列表行。
type StoreCouponListItem struct {
	Coupon
	MerName        string `json:"mer_name"`
	IsTrader       int8   `json:"is_trader"`
	TraderName     string `json:"trader_name"`
	CouponTypeName string `json:"coupon_type_name"`
	ClaimText      string `json:"claim_text"`
	ValidityText   string `json:"validity_text"`
	ReceivedTotal  int64  `json:"received_total"`
	UsedTotal      int64  `json:"used_total"`
}

type CouponSendInput struct {
	Mark string `json:"mark"`
	UIDs []uint `json:"uids"`
}

type StatusInput struct {
	Status int8 `json:"status"`
}

type BindSpreadInput struct {
	SpreadUID uint `json:"spread_uid"`
}

type MerTotal struct {
	MerID      uint
	TotalPrice float64
}

type QuoteInput struct {
	MerTotals          []MerTotal
	CouponUserIDs      []uint
	SkipStoreCoupon    bool // 秒杀/拼团/SVIP(未合并)：清空店铺券
	SkipPlatformCoupon bool // 非普通单：清空平台券（FUNCTIONAL-TRUTH）
}

type QuoteResult struct {
	CouponPrice          float64          `json:"coupon_price"`
	PlatformCouponUserID uint             `json:"platform_coupon_user_id"`
	PlatformDiscount     float64          `json:"platform_discount"`
	MerStoreDiscount     map[uint]float64 `json:"mer_store_discount"`
	MerCouponUserID      map[uint]uint    `json:"mer_coupon_user_id"`
	MerPlatformShare     map[uint]float64 `json:"mer_platform_share"`
	PayPrice             float64          `json:"pay_price"`
	TotalPrice           float64          `json:"total_price"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
