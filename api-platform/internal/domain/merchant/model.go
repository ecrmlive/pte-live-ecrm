package merchant

import "time"

type Merchant struct {
	MerID           uint      `gorm:"column:mer_id;primaryKey" json:"mer_id"`
	CategoryID      uint      `gorm:"column:category_id" json:"category_id"`
	MerName         string    `gorm:"column:mer_name" json:"mer_name"`
	RealName        string    `gorm:"column:real_name" json:"real_name"`
	MerPhone        string    `gorm:"column:mer_phone" json:"mer_phone"`
	MerAddress      string    `gorm:"column:mer_address" json:"mer_address"`
	MerInfo         string    `gorm:"column:mer_info" json:"mer_info"`
	Mark            string    `gorm:"column:mark" json:"mark"`
	Status          int8      `gorm:"column:status" json:"status"`
	MerState        int8      `gorm:"column:mer_state" json:"mer_state"`
	IsDel           int8      `gorm:"column:is_del" json:"-"`
	IsAudit         int8      `gorm:"column:is_audit" json:"is_audit"`
	SvipCouponMerge int8      `gorm:"column:svip_coupon_merge" json:"svip_coupon_merge"`
	RegionID        uint      `gorm:"column:region_id" json:"region_id"`
	CreateTime      time.Time `gorm:"column:create_time" json:"create_time"`
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
	MerchantCategoryID uint    `gorm:"column:merchant_category_id;primaryKey" json:"merchant_category_id"`
	CommissionRate     float64 `gorm:"column:commission_rate" json:"commission_rate"`
	CategoryName       string  `gorm:"column:category_name" json:"category_name"`
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
