package merchant

import "time"

type Merchant struct {
	MerID             uint      `gorm:"column:mer_id;primaryKey" json:"mer_id"`
	CategoryID        uint      `gorm:"column:category_id" json:"category_id"`
	TypeID            uint      `gorm:"column:type_id" json:"type_id"`
	BusinessID        uint      `gorm:"column:business_id" json:"business_id"`
	MerName           string    `gorm:"column:mer_name" json:"mer_name"`
	OwnerName         string    `gorm:"column:owner_name" json:"owner_name"`
	RealName          string    `gorm:"column:real_name" json:"real_name"`
	MerPhone          string    `gorm:"column:mer_phone" json:"mer_phone"`
	MerAddress        string    `gorm:"column:mer_address" json:"mer_address"`
	MerInfo           string    `gorm:"column:mer_info" json:"mer_info"`
	MerKeyword        string    `gorm:"column:mer_keyword" json:"mer_keyword"`
	Mark              string    `gorm:"column:mark" json:"mark"`
	Status            int8      `gorm:"column:status" json:"status"`
	MerState          int8      `gorm:"column:mer_state" json:"mer_state"`
	IsDel             int8      `gorm:"column:is_del" json:"-"`
	IsAudit           int8      `gorm:"column:is_audit" json:"is_audit"`
	IsBest            int8      `gorm:"column:is_best" json:"is_best"`
	OfflinePay        int8      `gorm:"column:offline_pay" json:"offline_pay"`
	IsTrader          int8      `gorm:"column:is_trader" json:"is_trader"`
	IsBroRoom         int8      `gorm:"column:is_bro_room" json:"is_bro_room"`
	IsBroGoods        int8      `gorm:"column:is_bro_goods" json:"is_bro_goods"`
	CommissionSwitch  int8      `gorm:"column:commission_switch" json:"commission_switch"`
	CommissionRate    float64   `gorm:"column:commission_rate" json:"commission_rate"`
	MerAccount        string    `gorm:"column:mer_account" json:"mer_account"`
	SubMchid          string    `gorm:"column:sub_mchid" json:"sub_mchid"`
	ApplymentID       string    `gorm:"column:applyment_id" json:"applyment_id"`
	CareCount         int       `gorm:"column:care_count" json:"care_count"`
	CareFicti         int       `gorm:"column:care_ficti" json:"care_ficti"`
	GoodsType         string    `gorm:"column:goods_type" json:"goods_type"`
	PlatformCategoryIDs string  `gorm:"column:platform_category_ids" json:"platform_category_ids"`
	MerStar           int8      `gorm:"column:mer_star" json:"mer_star"`
	MerAvatar         string    `gorm:"column:mer_avatar" json:"mer_avatar"`
	Sort              int       `gorm:"column:sort" json:"sort"`
	SvipCouponMerge   int8      `gorm:"column:svip_coupon_merge" json:"svip_coupon_merge"`
	RegionID          uint      `gorm:"column:region_id" json:"region_id"`
	CreateTime        time.Time `gorm:"column:create_time" json:"create_time"`
	// 列表/详情关联展示字段（非 merchant_view 原列）。
	CategoryName     string  `gorm:"column:category_name" json:"category_name"`
	TypeName         string  `gorm:"column:type_name" json:"type_name"`
	RegionName       string  `gorm:"column:region_name" json:"region_name"`
	DepositState     string  `gorm:"column:deposit_state" json:"deposit_state"`
	DepositRequired  float64 `gorm:"column:deposit_required" json:"deposit_required"`
	DepositAvailable float64 `gorm:"column:deposit_available" json:"deposit_available"`
	TypeMargin       float64 `gorm:"column:type_margin" json:"type_margin"`
	TypeIsMargin     int8    `gorm:"column:type_is_margin" json:"type_is_margin"`
	StoreGroupIDs    []uint  `gorm:"-" json:"store_group_ids"`
	GoodsTypes       []int   `gorm:"-" json:"goods_types"`
	PlatformCategoryIDList []uint `gorm:"-" json:"platform_category_id_list"`
}

type SvipConfig struct {
	MerID           uint `json:"mer_id"`
	SvipCouponMerge int8 `json:"svip_coupon_merge"`
}

// Merchant is the platform-owned, event-fed merchant supervision projection.
// It must never point at the legacy qixi_m_* namespace or a merchant-owned
// qixi_crm_m_* table.
func (Merchant) TableName() string { return "qixi_crm_a_merchant_view" }

type Category struct {
	MerchantCategoryID uint      `gorm:"column:merchant_category_id;primaryKey" json:"merchant_category_id"`
	CommissionRate     float64   `gorm:"column:commission_rate" json:"commission_rate"`
	CategoryName       string    `gorm:"column:category_name" json:"category_name"`
	CreateTime         time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Category) TableName() string { return "qixi_crm_a_merchant_category" }

type Intention struct {
	MerIntentionID      uint       `gorm:"column:mer_intention_id;primaryKey" json:"mer_intention_id"`
	SourceApplicationID uint       `gorm:"column:source_application_id" json:"-"`
	UID                 uint       `gorm:"column:uid" json:"uid"`
	Phone               string     `gorm:"column:phone" json:"phone"`
	MerName             string     `gorm:"column:mer_name" json:"mer_name"`
	Name                string     `gorm:"column:name" json:"name"`
	CreateTime          *time.Time `gorm:"column:create_time" json:"create_time"`
	Status              int8       `gorm:"column:status" json:"status"`
	FailMsg             string     `gorm:"column:fail_msg" json:"fail_msg"`
	IsDel               int8       `gorm:"column:is_del" json:"-"`
	Mark                string     `gorm:"column:mark" json:"mark"`
	MerID               uint       `gorm:"column:mer_id" json:"mer_id"`
	Images              string     `gorm:"column:images" json:"images"`
	CategoryName        string     `gorm:"column:category_name" json:"category_name"`
	TypeName            string     `gorm:"column:type_name" json:"type_name"`
	MerchantCategoryID  uint       `gorm:"column:merchant_category_id" json:"merchant_category_id"`
	MerTypeID           uint       `gorm:"column:mer_type_id" json:"mer_type_id"`
	CircleID            uint       `gorm:"column:circle_id" json:"circle_id"`
}

// Intention is the platform projection of a business-side merchant application.
func (Intention) TableName() string { return "qixi_crm_a_merchant_application" }

const (
	IntentionPending  int8 = 0
	IntentionApproved int8 = 1
	IntentionRejected int8 = 2
)
