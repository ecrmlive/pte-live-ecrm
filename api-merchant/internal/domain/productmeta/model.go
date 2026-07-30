package productmeta

import "time"

type Label struct {
	LabelID    uint      `gorm:"column:label_id;primaryKey" json:"label_id"`
	MerID      uint      `gorm:"column:mer_id" json:"mer_id"`
	Name       string    `gorm:"column:name" json:"name"`
	Info       string    `gorm:"column:info" json:"info"`
	Sort       int       `gorm:"column:sort" json:"sort"`
	Status     int8      `gorm:"column:status" json:"status"`
	IsDel      int8      `gorm:"column:is_del" json:"-"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Label) TableName() string { return "qixi_m_admin_store_product_label" }

type Guarantee struct {
	GuaranteeID uint      `gorm:"column:guarantee_id;primaryKey" json:"guarantee_id"`
	MerID       uint      `gorm:"column:mer_id" json:"mer_id"`
	Name        string    `gorm:"column:name" json:"name"`
	Content     string    `gorm:"column:content" json:"content"`
	Sort        int       `gorm:"column:sort" json:"sort"`
	Status      int8      `gorm:"column:status" json:"status"`
	IsDel       int8      `gorm:"column:is_del" json:"-"`
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`
}

func (Guarantee) TableName() string { return "qixi_m_admin_store_product_guarantee" }

type AttrTemplate struct {
	TemplateID    uint      `gorm:"column:template_id;primaryKey" json:"template_id"`
	MerID         uint      `gorm:"column:mer_id" json:"mer_id"`
	TemplateName  string    `gorm:"column:template_name" json:"template_name"`
	TemplateValue string    `gorm:"column:template_value" json:"template_value"`
	Sort          int       `gorm:"column:sort" json:"sort"`
	IsDel         int8      `gorm:"column:is_del" json:"-"`
	CreateTime    time.Time `gorm:"column:create_time" json:"create_time"`
}

func (AttrTemplate) TableName() string { return "qixi_m_admin_store_product_attr_template" }

type LabelInput struct {
	Name   string `json:"name"`
	Info   string `json:"info"`
	Sort   int    `json:"sort"`
	Status *int8  `json:"status"`
}

type GuaranteeInput struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Sort    int    `json:"sort"`
	Status  *int8  `json:"status"`
}

type AttrTemplateInput struct {
	TemplateName  string `json:"template_name"`
	TemplateValue string `json:"template_value"`
	Sort          int    `json:"sort"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
