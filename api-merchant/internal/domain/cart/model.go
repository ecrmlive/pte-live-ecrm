package cart

import (
	"strings"
	"time"
)

type Cart struct {
	CartID            uint64    `gorm:"column:cart_id;primaryKey" json:"cart_id"`
	UID               uint      `gorm:"column:uid" json:"uid"`
	MerID             uint      `gorm:"column:mer_id" json:"mer_id"`
	ProductType       int8      `gorm:"column:product_type" json:"product_type"`
	ProductID         uint      `gorm:"column:product_id" json:"product_id"`
	ProductAttrUnique string    `gorm:"column:product_attr_unique" json:"product_attr_unique"`
	CartNum           uint      `gorm:"column:cart_num" json:"cart_num"`
	IsPay             int8      `gorm:"column:is_pay" json:"is_pay"`
	IsDel             int8      `gorm:"column:is_del" json:"-"`
	IsNew             int8      `gorm:"column:is_new" json:"is_new"`
	IsFail            int8      `gorm:"column:is_fail" json:"is_fail"`
	CreateTime        time.Time `gorm:"column:create_time" json:"create_time"`

	StoreName       string  `gorm:"-" json:"store_name,omitempty"`
	MerName         string  `gorm:"-" json:"mer_name,omitempty"`
	Image           string  `gorm:"-" json:"image,omitempty"`
	Price           float64 `gorm:"-" json:"price,omitempty"`
	OtPrice         float64 `gorm:"-" json:"ot_price,omitempty"`
	Cost            float64 `gorm:"-" json:"cost,omitempty"`
	Stock           uint    `gorm:"-" json:"stock,omitempty"`
	GoodsType       uint8   `gorm:"-" json:"goods_type,omitempty"`
	SeckillActiveID uint    `gorm:"-" json:"seckill_active_id,omitempty"`
	PresellActiveID uint    `gorm:"-" json:"presell_active_id,omitempty"`
	OncePayCount    int     `gorm:"-" json:"once_pay_count,omitempty"`
	SvipPriceType   int8    `gorm:"-" json:"svip_price_type,omitempty"`
	SvipPrice       float64 `gorm:"-" json:"svip_price,omitempty"`
	MerSvipStatus   int8    `gorm:"-" json:"mer_svip_status,omitempty"`
	UsedSvip        bool    `gorm:"-" json:"used_svip,omitempty"`
	SvipDiscount    float64 `gorm:"-" json:"svip_discount,omitempty"` // 行总优惠
}

func (Cart) TableName() string { return "qixi_m_app_store_cart" }

type AddInput struct {
	ProductID         uint   `json:"product_id"`
	ProductAttrUnique string `json:"product_attr_unique"`
	CartNum           uint   `json:"cart_num"`
	IsNew             int8   `json:"is_new"`
}

type MerchantBucket struct {
	MerID    uint    `json:"mer_id"`
	MerName  string  `json:"mer_name"`
	Subtotal float64 `json:"subtotal"`
	Items    []Cart  `json:"items"`
}

type Address struct {
	AddressID uint   `gorm:"column:address_id;primaryKey" json:"address_id"`
	UID       uint   `gorm:"column:uid" json:"uid"`
	RealName  string `gorm:"column:real_name" json:"real_name"`
	Phone     string `gorm:"column:phone" json:"phone"`
	Province  string `gorm:"column:province" json:"province"`
	City      string `gorm:"column:city" json:"city"`
	District  string `gorm:"column:district" json:"district"`
	Detail    string `gorm:"column:detail" json:"detail"`
	PostCode  int    `gorm:"column:post_code" json:"post_code"`
	IsDefault int8   `gorm:"column:is_default" json:"is_default"`
	IsDel     int8   `gorm:"column:is_del" json:"-"`
}

func (Address) TableName() string { return "qixi_m_app_user_address" }

func (a Address) FullAddress() string {
	return strings.TrimSpace(a.Province + a.City + a.District + a.Detail)
}

type AddressInput struct {
	RealName  string `json:"real_name"`
	Phone     string `json:"phone"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Detail    string `json:"detail"`
	PostCode  int    `json:"post_code"`
	IsDefault *int8  `json:"is_default"`
}
