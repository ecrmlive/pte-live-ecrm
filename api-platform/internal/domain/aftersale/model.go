package aftersale

import "time"

// 退款单状态（FUNCTIONAL-TRUTH §6）
const (
	StatusWait     int8 = 0  // 待审核
	StatusReject   int8 = -1 // 审核拒绝
	StatusBack     int8 = 1  // 待退货
	StatusReceive  int8 = 2  // 待收货
	StatusRefunded int8 = 3  // 已退款
	StatusPlatform int8 = 4  // 平台介入
	StatusCancel   int8 = -2 // 用户取消
)

const (
	RefundTypeMoneyOnly int8 = 1 // 仅退款
	RefundTypeReturn    int8 = 2 // 退货退款（P0 不实现）
)

const (
	OrderProductRefundNone    int8 = 0
	OrderProductRefunding     int8 = 1
	OrderProductRefundPartial int8 = 2
	OrderProductRefundFull    int8 = 3
)

const OrderStatusRefunded int8 = -1

type RefundOrder struct {
	RefundOrderID       uint      `gorm:"column:refund_order_id;primaryKey" json:"refund_order_id"`
	RefundOrderSN       string    `gorm:"column:refund_order_sn" json:"refund_order_sn"`
	OrderID             uint      `gorm:"column:order_id" json:"order_id"`
	UID                 uint      `gorm:"column:uid" json:"uid"`
	MerID               uint      `gorm:"column:mer_id" json:"mer_id"`
	DeliveryType        string    `gorm:"column:delivery_type" json:"delivery_type,omitempty"`
	DeliveryID          string    `gorm:"column:delivery_id" json:"delivery_id,omitempty"`
	DeliveryMark        string    `gorm:"column:delivery_mark" json:"delivery_mark,omitempty"`
	Phone               string    `gorm:"column:phone" json:"phone,omitempty"`
	Mark                string    `gorm:"column:mark" json:"mark"`
	MerMark             string    `gorm:"column:mer_mark" json:"mer_mark"`
	AdminMark           string    `gorm:"column:admin_mark" json:"admin_mark"`
	Pics                string    `gorm:"column:pics" json:"pics,omitempty"`
	RefundType          int8      `gorm:"column:refund_type" json:"refund_type"`
	RefundMessage       string    `gorm:"column:refund_message" json:"refund_message"`
	RefundPrice         float64   `gorm:"column:refund_price" json:"refund_price"`
	PlatformRefundPrice float64   `gorm:"column:platform_refund_price" json:"platform_refund_price"`
	RefundPostage       float64   `gorm:"column:refund_postage" json:"refund_postage"`
	RefundNum           int       `gorm:"column:refund_num" json:"refund_num"`
	FailMessage         string    `gorm:"column:fail_message" json:"fail_message,omitempty"`
	Status              int8      `gorm:"column:status" json:"status"`
	StatusTime          time.Time `gorm:"column:status_time" json:"status_time"`
	CreateTime          time.Time `gorm:"column:create_time" json:"create_time"`
	IsDel               int8      `gorm:"column:is_del" json:"-"`
	IsSystemDel         int8      `gorm:"column:is_system_del" json:"-"`

	Products []RefundProduct `gorm:"-" json:"products,omitempty"`
}

func (RefundOrder) TableName() string { return "qixi_m_app_store_refund_order" }

type RefundProduct struct {
	RefundProductID uint      `gorm:"column:refund_product_id;primaryKey" json:"refund_product_id"`
	RefundOrderID   uint      `gorm:"column:refund_order_id" json:"refund_order_id"`
	OrderProductID  uint      `gorm:"column:order_product_id" json:"order_product_id"`
	RefundPrice     float64   `gorm:"column:refund_price" json:"refund_price"`
	RefundNum       int       `gorm:"column:refund_num" json:"refund_num"`
	CreateTime      time.Time `gorm:"column:create_time" json:"create_time"`
}

func (RefundProduct) TableName() string { return "qixi_m_app_store_refund_product" }

type RefundStatusLog struct {
	ID            uint      `gorm:"column:id;primaryKey" json:"id"`
	RefundOrderID uint      `gorm:"column:refund_order_id" json:"refund_order_id"`
	ChangeType    string    `gorm:"column:change_type" json:"change_type"`
	ChangeMessage string    `gorm:"column:change_message" json:"change_message"`
	ChangeTime    time.Time `gorm:"column:change_time" json:"change_time"`
}

func (RefundStatusLog) TableName() string { return "qixi_m_app_store_refund_status" }

// OrderProductLine 订单商品行（含退款字段；不依赖 trade.OrderProduct）
type OrderProductLine struct {
	OrderProductID uint    `gorm:"column:order_product_id;primaryKey" json:"order_product_id"`
	OrderID        uint    `gorm:"column:order_id" json:"order_id"`
	UID            uint    `gorm:"column:uid" json:"uid"`
	ProductID      uint    `gorm:"column:product_id" json:"product_id"`
	ProductSKU     string  `gorm:"column:product_sku" json:"product_sku"`
	ProductNum     int     `gorm:"column:product_num" json:"product_num"`
	TotalPrice     float64 `gorm:"column:total_price" json:"total_price"`
	IsRefund       int8    `gorm:"column:is_refund" json:"is_refund"`
	RefundNum      int     `gorm:"column:refund_num" json:"refund_num"`
	CartInfo       string  `gorm:"column:cart_info" json:"cart_info,omitempty"`
}

func (OrderProductLine) TableName() string { return "qixi_m_app_store_order_product" }

type StoreOrderBrief struct {
	OrderID  uint    `gorm:"column:order_id;primaryKey"`
	UID      uint    `gorm:"column:uid"`
	MerID    uint    `gorm:"column:mer_id"`
	Paid     int8    `gorm:"column:paid"`
	Status   int8    `gorm:"column:status"`
	PayPrice float64 `gorm:"column:pay_price"`
	IsDel    int8    `gorm:"column:is_del"`
}

func (StoreOrderBrief) TableName() string { return "qixi_m_app_store_order" }

type ApplyInput struct {
	OrderID         uint   `json:"order_id"`
	RefundType      int8   `json:"refund_type"`
	RefundMessage   string `json:"refund_message"`
	OrderProductIDs []uint `json:"order_product_ids"`
}

type RejectInput struct {
	FailMessage string `json:"fail_message"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
