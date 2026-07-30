package trade

import "time"

// PayType 存储值：0=balance，1=wechat，2=alipay，7=mock，8=integral
// API 入参为 "balance"|"wechat"|"alipay"|"mock"|"integral"
const (
	PayTypeBalance  int8 = 0
	PayTypeWechat   int8 = 1
	PayTypeAlipay   int8 = 2
	PayTypeMock     int8 = 7
	PayTypeIntegral int8 = 8
)

const (
	OrderStatusAwaitShip int8 = 0 // 待发货
	OrderStatusShipped   int8 = 1 // 待收货
	OrderStatusDone      int8 = 3 // 完成
)

// ActivityTypeSeckill 秒杀（订单行 product_type=1）
const ActivityTypeSeckill int8 = 1

// ActivityTypePoints 积分商城单（FUNCTIONAL-TRUTH：activity_type=20）
const ActivityTypePoints int8 = 20

const GoodsTypeNormal uint8 = 0
const GoodsTypePoints uint8 = 1

type GroupOrder struct {
	GroupOrderID  uint       `gorm:"column:group_order_id;primaryKey" json:"group_order_id"`
	GroupOrderSN  string     `gorm:"column:group_order_sn" json:"group_order_sn"`
	UID           uint       `gorm:"column:uid" json:"uid"`
	TotalPostage  float64    `gorm:"column:total_postage" json:"total_postage"`
	TotalPrice    float64    `gorm:"column:total_price" json:"total_price"`
	TotalNum      int        `gorm:"column:total_num" json:"total_num"`
	CouponPrice   float64    `gorm:"column:coupon_price" json:"coupon_price"`
	CouponID      uint       `gorm:"column:coupon_id" json:"coupon_id"`
	Integral      int        `gorm:"column:integral" json:"integral"`
	IntegralPrice float64    `gorm:"column:integral_price" json:"integral_price"`
	GiveIntegral  int        `gorm:"column:give_integral" json:"give_integral"`
	RealName      string     `gorm:"column:real_name" json:"real_name"`
	UserPhone     string     `gorm:"column:user_phone" json:"user_phone"`
	UserAddress   string     `gorm:"column:user_address" json:"user_address"`
	PayPrice      float64    `gorm:"column:pay_price" json:"pay_price"`
	PayPostage    float64    `gorm:"column:pay_postage" json:"pay_postage"`
	Cost          float64    `gorm:"column:cost" json:"cost"`
	Paid          int8       `gorm:"column:paid" json:"paid"`
	PayTime       *time.Time `gorm:"column:pay_time" json:"pay_time"`
	PayType       int8       `gorm:"column:pay_type" json:"pay_type"`
	ActivityType  int8       `gorm:"column:activity_type" json:"activity_type"`
	CreateTime    time.Time  `gorm:"column:create_time" json:"create_time"`
	IsDel         int8       `gorm:"column:is_del" json:"-"`

	Orders []StoreOrder `gorm:"-" json:"orders,omitempty"`
}

func (GroupOrder) TableName() string { return "qixi_m_app_store_group_order" }

type StoreOrder struct {
	OrderID             uint       `gorm:"column:order_id;primaryKey" json:"order_id"`
	GroupOrderID        uint       `gorm:"column:group_order_id" json:"group_order_id"`
	OrderSN             string     `gorm:"column:order_sn" json:"order_sn"`
	UID                 uint       `gorm:"column:uid" json:"uid"`
	RealName            string     `gorm:"column:real_name" json:"real_name"`
	UserPhone           string     `gorm:"column:user_phone" json:"user_phone"`
	UserAddress         string     `gorm:"column:user_address" json:"user_address"`
	CartID              string     `gorm:"column:cart_id" json:"cart_id"`
	TotalNum            int        `gorm:"column:total_num" json:"total_num"`
	TotalPrice          float64    `gorm:"column:total_price" json:"total_price"`
	TotalPostage        float64    `gorm:"column:total_postage" json:"total_postage"`
	PayPrice            float64    `gorm:"column:pay_price" json:"pay_price"`
	PayPostage          float64    `gorm:"column:pay_postage" json:"pay_postage"`
	Paid                int8       `gorm:"column:paid" json:"paid"`
	PayTime             *time.Time `gorm:"column:pay_time" json:"pay_time"`
	PayType             int8       `gorm:"column:pay_type" json:"pay_type"`
	CreateTime          time.Time  `gorm:"column:create_time" json:"create_time"`
	Status              int8       `gorm:"column:status" json:"status"`
	DeliveryType        string     `gorm:"column:delivery_type" json:"delivery_type"`
	DeliveryName        string     `gorm:"column:delivery_name" json:"delivery_name"`
	DeliveryID          string     `gorm:"column:delivery_id" json:"delivery_id"`
	Mark                string     `gorm:"column:mark" json:"mark"`
	MerID               uint       `gorm:"column:mer_id" json:"mer_id"`
	Cost                float64    `gorm:"column:cost" json:"cost"`
	CouponID            string     `gorm:"column:coupon_id" json:"coupon_id"`
	CouponPrice         float64    `gorm:"column:coupon_price" json:"coupon_price"`
	PlatformCouponPrice float64    `gorm:"column:platform_coupon_price" json:"platform_coupon_price"`
	SvipDiscount        float64    `gorm:"column:svip_discount" json:"svip_discount"`
	Integral            int        `gorm:"column:integral" json:"integral"`
	IntegralPrice       float64    `gorm:"column:integral_price" json:"integral_price"`
	GiveIntegral        int        `gorm:"column:give_integral" json:"give_integral"`
	ActivityType        int8       `gorm:"column:activity_type" json:"activity_type"`
	VerifyCode          string     `gorm:"column:verify_code" json:"verify_code,omitempty"`
	VerifyServiceID     *uint      `gorm:"column:verify_service_id" json:"verify_service_id,omitempty"`
	VerifyTime          *time.Time `gorm:"column:verify_time" json:"verify_time,omitempty"`
	ReservationDate     string     `gorm:"column:reservation_date" json:"reservation_date,omitempty"`
	ReservationID       uint       `gorm:"column:reservation_id" json:"reservation_id,omitempty"`
	ReservationTimePart string     `gorm:"column:reservation_time_part" json:"reservation_time_part,omitempty"`
	SpreadUID           uint       `gorm:"-" json:"spread_uid"`
	ExtensionOne        float64    `gorm:"-" json:"extension_one"`
	IsDel               int8       `gorm:"column:is_del" json:"-"`

	MerName  string         `gorm:"-" json:"mer_name,omitempty"`
	Products []OrderProduct `gorm:"-" json:"products,omitempty"`
}

func (StoreOrder) TableName() string { return "qixi_m_app_store_order" }

type OrderProduct struct {
	OrderProductID uint    `gorm:"column:order_product_id;primaryKey" json:"order_product_id"`
	OrderID        uint    `gorm:"column:order_id" json:"order_id"`
	UID            uint    `gorm:"column:uid" json:"uid"`
	CartID         int     `gorm:"column:cart_id" json:"cart_id"`
	ProductID      uint    `gorm:"column:product_id" json:"product_id"`
	ProductSKU     string  `gorm:"column:product_sku" json:"product_sku"`
	ProductNum     int     `gorm:"column:product_num" json:"product_num"`
	ProductType    int8    `gorm:"column:product_type" json:"product_type"`
	ActivityID     uint    `gorm:"column:activity_id" json:"activity_id"`
	Cost           float64 `gorm:"column:cost" json:"cost"`
	ProductPrice   float64 `gorm:"column:product_price" json:"product_price"`
	TotalPrice     float64 `gorm:"column:total_price" json:"total_price"`
	ProductInfo    string  `gorm:"column:cart_info" json:"product_info"`
}

type CheckInput struct {
	CartIDs       []uint64 `json:"cart_ids"`
	CouponUserID  uint     `json:"coupon_user_id"`  // P0 平台券
	CouponUserIDs []uint   `json:"coupon_user_ids"` // 扩展：多券（店铺+平台）
	UseIntegral   int      `json:"use_integral"`
}

// NormalizedCouponUserIDs 合并 singular/plural 选券入参。
func (in CheckInput) NormalizedCouponUserIDs() []uint {
	ids := append([]uint{}, in.CouponUserIDs...)
	if in.CouponUserID > 0 {
		ids = append(ids, in.CouponUserID)
	}
	return ids
}

type PayInput struct {
	PayType string `json:"pay_type"`
}

func (OrderProduct) TableName() string { return "qixi_m_app_store_order_product" }

type CreateInput struct {
	CartIDs       []uint64 `json:"cart_ids"`
	AddressID     uint     `json:"address_id"`
	Mark          string   `json:"mark"`
	CouponUserID  uint     `json:"coupon_user_id"`
	CouponUserIDs []uint   `json:"coupon_user_ids"`
	UseIntegral   int      `json:"use_integral"`
}

func (in CreateInput) NormalizedCouponUserIDs() []uint {
	ids := append([]uint{}, in.CouponUserIDs...)
	if in.CouponUserID > 0 {
		ids = append(ids, in.CouponUserID)
	}
	return ids
}

type CheckMerchant struct {
	MerID               uint        `json:"mer_id"`
	MerName             string      `json:"mer_name"`
	Postage             float64     `json:"postage"`
	TotalPrice          float64     `json:"total_price"`
	TotalNum            int         `json:"total_num"`
	CouponPrice         float64     `json:"coupon_price"`
	PlatformCouponPrice float64     `json:"platform_coupon_price"`
	Integral            int         `json:"integral"`
	IntegralPrice       float64     `json:"integral_price"`
	GiveIntegral        int         `json:"give_integral"`
	PayPrice            float64     `json:"pay_price"`
	Items               []CheckItem `json:"items"`
}

// UserBill 积分/资金流水（表 qixi_m_app_user_bill）
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

const (
	BillPMOut int8 = 0
	BillPMIn  int8 = 1
)

// PointsProductView 积分商城下单用商品快照
type PointsProductView struct {
	ProductID   uint
	MerID       uint
	StoreName   string
	MerName     string
	Image       string
	Price       float64
	Cost        float64
	Stock       uint
	Integral    int
	Unique      string
	ProductType uint8
}

type CheckItem struct {
	CartID            uint64  `json:"cart_id"`
	ProductID         uint    `json:"product_id"`
	ProductAttrUnique string  `json:"product_attr_unique"`
	StoreName         string  `json:"store_name"`
	Image             string  `json:"image"`
	Price             float64 `json:"price"`
	CartNum           uint    `json:"cart_num"`
	Subtotal          float64 `json:"subtotal"`
}

type CheckResult struct {
	Merchants            []CheckMerchant `json:"merchants"`
	TotalPrice           float64         `json:"total_price"`
	TotalPostage         float64         `json:"total_postage"`
	CouponPrice          float64         `json:"coupon_price"`
	PlatformCouponUserID uint            `json:"platform_coupon_user_id"`
	SvipDiscount         float64         `json:"svip_discount"`
	UsedSvip             bool            `json:"used_svip"`
	Integral             int             `json:"integral"`
	IntegralPrice        float64         `json:"integral_price"`
	UserIntegral         int             `json:"user_integral"`
	PayPrice             float64         `json:"pay_price"`
	TotalNum             int             `json:"total_num"`
	GiveIntegral         int             `json:"give_integral"`

	merCouponUserIDs map[uint]uint `json:"-"`
}

type DeliveryInput struct {
	DeliveryName string `json:"delivery_name"`
	DeliveryID   string `json:"delivery_id"`
	DeliveryType string `json:"delivery_type"`
	ExpressID    uint   `json:"express_id"`    // 可选：平台快递公司 id，用于回填 delivery_name
	TemplateName string `json:"template_name"` // 可选备注，不影响发货主流程
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
