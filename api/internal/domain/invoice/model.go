package invoice

import "time"

const (
	StatusPending  int8 = 0
	StatusIssued   int8 = 1
	StatusRejected int8 = -1
)

type Invoice struct {
	InvoiceID   uint      `gorm:"column:invoice_id;primaryKey" json:"invoice_id"`
	UID         uint      `gorm:"column:uid" json:"uid"`
	OrderID     uint      `gorm:"column:order_id" json:"order_id"`
	MerID       uint      `gorm:"column:mer_id" json:"mer_id"`
	InvoiceType int8      `gorm:"column:invoice_type" json:"invoice_type"`
	HeaderType  int8      `gorm:"column:header_type" json:"header_type"`
	Header      string    `gorm:"column:header" json:"header"`
	TaxNo       string    `gorm:"column:tax_no" json:"tax_no"`
	Email       string    `gorm:"column:email" json:"email"`
	Status      int8      `gorm:"column:status" json:"status"`
	Mark        string    `gorm:"column:mark" json:"mark"`
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`
	IsDel       int8      `gorm:"column:is_del" json:"-"`
}

func (Invoice) TableName() string { return "qixi_store_order_invoice" }

type ApplyInput struct {
	OrderID     uint   `json:"order_id"`
	InvoiceType int8   `json:"invoice_type"`
	HeaderType  int8   `json:"header_type"`
	Header      string `json:"header"`
	TaxNo       string `json:"tax_no"`
	Email       string `json:"email"`
}

type AuditInput struct {
	Status int8   `json:"status"` // 1已开 -1驳回
	Mark   string `json:"mark"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
