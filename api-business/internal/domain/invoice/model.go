package invoice

import "time"

const (
	ProfilePersonal   = "personal"
	ProfileEnterprise = "enterprise"

	StatusRequested = "requested"
	StatusIssued    = "issued"
	StatusRejected  = "rejected"
	StatusVoided    = "voided"
)

// InvoiceProfile is a user-owned invoice header. An order invoice keeps a
// snapshot so later profile changes do not alter a submitted application.
type InvoiceProfile struct {
	ID        uint64    `gorm:"column:id;primaryKey" json:"id"`
	UserID    uint64    `gorm:"column:user_id" json:"-"`
	Type      string    `gorm:"column:type" json:"type"`
	Title     string    `gorm:"column:title" json:"title"`
	TaxNo     string    `gorm:"column:tax_no" json:"tax_no"`
	Email     string    `gorm:"column:email" json:"email"`
	IsDefault bool      `gorm:"column:is_default" json:"is_default"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (InvoiceProfile) TableName() string { return "qixi_crm_b_user_invoice_profile" }

// Invoice is a C-end invoice application. Status transitions are owned by the
// merchant invoice workflow; the user endpoint only creates and reads it.
type Invoice struct {
	ID               uint64     `gorm:"column:id;primaryKey" json:"id"`
	OrderID          uint64     `gorm:"column:order_id" json:"order_id"`
	InvoiceProfileID uint64     `gorm:"column:invoice_profile_id" json:"invoice_profile_id"`
	ProfileType      string     `gorm:"column:profile_type" json:"profile_type"`
	Title            string     `gorm:"column:title" json:"title"`
	TaxNo            string     `gorm:"column:tax_no" json:"tax_no"`
	Email            string     `gorm:"column:email" json:"email"`
	Status           string     `gorm:"column:status" json:"status"`
	InvoiceNo        string     `gorm:"column:invoice_no" json:"invoice_no"`
	FileURL          string     `gorm:"column:file_url" json:"file_url"`
	RejectionReason  string     `gorm:"column:rejection_reason" json:"rejection_reason"`
	RequestedAt      time.Time  `gorm:"column:requested_at" json:"requested_at"`
	IssuedAt         *time.Time `gorm:"column:issued_at" json:"issued_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (Invoice) TableName() string { return "qixi_crm_b_order_invoice" }

type ProfileInput struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	TaxNo     string `json:"tax_no"`
	Email     string `json:"email"`
	IsDefault bool   `json:"is_default"`
}

type ApplyInput struct {
	OrderID          uint64 `json:"order_id"`
	InvoiceProfileID uint64 `json:"invoice_profile_id"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
