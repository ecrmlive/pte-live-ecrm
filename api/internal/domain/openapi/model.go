package openapi

import "time"

// AuthRuleProduct / AuthRuleOrder 对齐 CRMEB open_auth.auth：1=商品 2=订单
const (
	AuthRuleProduct = "1"
	AuthRuleOrder   = "2"
)

type OpenAuth struct {
	ID         uint       `gorm:"column:id;primaryKey" json:"id"`
	Title      string     `gorm:"column:title" json:"title"`
	AccessKey  string     `gorm:"column:access_key" json:"access_key"`
	SecretKey  string     `gorm:"column:secret_key" json:"-"`
	Status     int8       `gorm:"column:status" json:"status"`
	Mark       string     `gorm:"column:mark" json:"mark"`
	MerID      uint       `gorm:"column:mer_id" json:"mer_id"`
	Auth       string     `gorm:"column:auth" json:"auth"`
	Sort       int        `gorm:"column:sort" json:"sort"`
	IsDel      int8       `gorm:"column:is_del" json:"-"`
	CreateTime time.Time  `gorm:"column:create_time" json:"create_time"`
	LastIP     string     `gorm:"column:last_ip" json:"last_ip,omitempty"`
	LastTime   *time.Time `gorm:"column:last_time" json:"last_time,omitempty"`
}

func (OpenAuth) TableName() string { return "qixi_open_auth" }

// AuthInput 对齐 CRMEB POST /openapi/auth 四要素；secret_key 仅本地演示捷径（有签名时忽略）。
type AuthInput struct {
	Unique     string `json:"unique"`
	Expiration int64  `json:"expiration"`
	AccessKey  string `json:"access_key"`
	Signature  string `json:"signature"`
	SecretKey  string `json:"secret_key"`
}

type AuthResult struct {
	Token     string `json:"token"`
	Exp       int64  `json:"exp"`
	MerID     uint   `json:"mer_id"`
	AccessKey string `json:"access_key"`
}
