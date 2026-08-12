package logistics

import "time"

type Express struct {
	ExpressID  uint      `gorm:"column:express_id;primaryKey" json:"express_id"`
	Name       string    `gorm:"column:name" json:"name"`
	Code       string    `gorm:"column:code" json:"code"`
	Sort       int       `gorm:"column:sort" json:"sort"`
	IsShow     int8      `gorm:"column:is_show" json:"is_show"`
	IsDel      int8      `gorm:"column:is_del" json:"-"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Express) TableName() string { return "qixi_crm_a_express" }

type City struct {
	CityID   uint   `gorm:"column:city_id;primaryKey" json:"city_id"`
	ParentID uint   `gorm:"column:parent_id" json:"parent_id"`
	Name     string `gorm:"column:name" json:"name"`
	Level    int8   `gorm:"column:level" json:"level"`
	IsShow   int8   `gorm:"column:is_show" json:"is_show"`
}

func (City) TableName() string { return "qixi_crm_a_city" }

type ShippingTemplate struct {
	TemplateID uint      `gorm:"column:template_id;primaryKey" json:"template_id"`
	MerID      uint      `gorm:"column:mer_id" json:"mer_id"`
	Name       string    `gorm:"column:name" json:"name"`
	Type       int8      `gorm:"column:type" json:"type"`
	Appoint    int8      `gorm:"column:appoint" json:"appoint"`
	IsDefault  int8      `gorm:"column:is_default" json:"is_default"`
	Sort       int       `gorm:"column:sort" json:"sort"`
	IsDel      int8      `gorm:"column:is_del" json:"-"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	Regions    []Region  `gorm:"-" json:"regions,omitempty"`
}

// 运费模板属于店铺事实。当前平台端未注册其 CRUD 路由，后续商户端接入时必须连接 merchant 库。
func (ShippingTemplate) TableName() string { return "qixi_crm_m_shipping_template" }

type Region struct {
	RegionID      uint    `gorm:"column:region_id;primaryKey" json:"region_id"`
	TemplateID    uint    `gorm:"column:template_id" json:"template_id"`
	CityIDs       string  `gorm:"column:city_ids" json:"city_ids"`
	First         float64 `gorm:"column:first" json:"first"`
	FirstPrice    float64 `gorm:"column:first_price" json:"first_price"`
	Continue      float64 `gorm:"column:continue" json:"continue"`
	ContinuePrice float64 `gorm:"column:continue_price" json:"continue_price"`
}

func (Region) TableName() string { return "qixi_crm_m_shipping_template_region" }

type ExpressInput struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Sort   int    `json:"sort"`
	IsShow *int8  `json:"is_show"`
}

// ExpressSyncResult 是平台内置物流公司目录同步的结果。同步仅维护目录名称和
// 已删除状态，不覆盖运营人员维护的排序与是否显示。
type ExpressSyncResult struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Total   int `json:"total"`
}

type RegionInput struct {
	CityIDs       string  `json:"city_ids"`
	First         float64 `json:"first"`
	FirstPrice    float64 `json:"first_price"`
	Continue      float64 `json:"continue"`
	ContinuePrice float64 `json:"continue_price"`
}

type TemplateInput struct {
	Name    string        `json:"name"`
	Type    int8          `json:"type"`
	Appoint int8          `json:"appoint"`
	Sort    int           `json:"sort"`
	Regions []RegionInput `json:"regions"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
