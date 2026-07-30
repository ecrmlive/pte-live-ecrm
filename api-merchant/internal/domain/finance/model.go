package finance

import "time"

const (
	TypeExtract int = 0 // 余额提现
	TypeDeposit int = 1 // 保证金
)

const (
	StatusWait   int = 0
	StatusPass   int = 1
	StatusReject int = -1
)

const (
	FinancialStatusUnpaid int = 0
	FinancialStatusPaid   int = 1
)

const (
	AccountBank   int = 1
	AccountWechat int = 2
	AccountAlipay int = 3
)

const (
	RecordPMOut int8 = 0 // 支出
	RecordPMIn  int8 = 1 // 收入
)

type Financial struct {
	FinancialID      uint       `gorm:"column:financial_id;primaryKey" json:"financial_id"`
	FinancialSN      string     `gorm:"column:financial_sn" json:"financial_sn"`
	MerMoney         float64    `gorm:"column:mer_money" json:"mer_money"`
	ExtractMoney     float64    `gorm:"column:extract_money" json:"extract_money"`
	FinancialType    int        `gorm:"column:financial_type" json:"financial_type"`
	FinancialAccount string     `gorm:"column:financial_account" json:"financial_account"`
	FinancialStatus  int        `gorm:"column:financial_status" json:"financial_status"`
	Status           int        `gorm:"column:status" json:"status"`
	Refusal          string     `gorm:"column:refusal" json:"refusal,omitempty"`
	MerID            uint       `gorm:"column:mer_id" json:"mer_id"`
	Image            string     `gorm:"column:image" json:"image,omitempty"`
	AdminID          *int       `gorm:"column:admin_id" json:"admin_id,omitempty"`
	CreateTime       *time.Time `gorm:"column:create_time" json:"create_time"`
	StatusTime       *time.Time `gorm:"column:status_time" json:"status_time,omitempty"`
	UpdateTime       *time.Time `gorm:"column:update_time" json:"update_time,omitempty"`
	IsDel            int        `gorm:"column:is_del" json:"-"`
	Mark             string     `gorm:"column:mark" json:"mark,omitempty"`
	AdminMark        string     `gorm:"column:admin_mark" json:"admin_mark,omitempty"`
	MerAdminID       *int       `gorm:"column:mer_admin_id" json:"mer_admin_id,omitempty"`
	Type             int        `gorm:"column:type" json:"type"`
}

func (Financial) TableName() string { return "qixi_m_admin_financial" }

type FinancialRecord struct {
	FinancialRecordID uint      `gorm:"column:financial_record_id;primaryKey" json:"financial_record_id"`
	FinancialRecordSN string    `gorm:"column:financial_record_sn" json:"financial_record_sn"`
	OrderID           uint      `gorm:"column:order_id" json:"order_id"`
	OrderSN           string    `gorm:"column:order_sn" json:"order_sn"`
	UserInfo          string    `gorm:"column:user_info" json:"user_info"`
	UserID            uint      `gorm:"column:user_id" json:"user_id"`
	FinancialType     string    `gorm:"column:financial_type" json:"financial_type"`
	FinancialPM       int8      `gorm:"column:financial_pm" json:"financial_pm"`
	Number            float64   `gorm:"column:number" json:"number"`
	Type              int8      `gorm:"column:type" json:"type"`
	MerID             uint      `gorm:"column:mer_id" json:"mer_id"`
	CreateTime        time.Time `gorm:"column:create_time" json:"create_time"`
	PayType           int       `gorm:"column:pay_type" json:"pay_type"`
}

func (FinancialRecord) TableName() string { return "qixi_m_admin_financial_record" }

type Balance struct {
	MerID    uint    `json:"mer_id"`
	MerMoney float64 `json:"mer_money"`
}

type WithdrawInput struct {
	ExtractMoney     float64 `json:"extract_money"`
	FinancialType    int     `json:"financial_type"`
	FinancialAccount string  `json:"financial_account"`
	Mark             string  `json:"mark"`
}

type RejectInput struct {
	Refusal string `json:"refusal"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
