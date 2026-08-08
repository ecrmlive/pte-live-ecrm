package circle

import "time"

// Circle 区域代理可经营、结算的商圈层级。type=0 为区域，type=1 为商户型商圈。
type Circle struct {
	CircleID              uint      `gorm:"column:circle_id;primaryKey" json:"circle_id"`
	PID                   uint      `gorm:"column:pid" json:"pid"`
	Path                  string    `gorm:"column:path" json:"path"`
	Name                  string    `gorm:"column:name" json:"name"`
	CircleAgentID         uint      `gorm:"column:circle_agent_id" json:"circle_agent_id"`
	CommissionType        int8      `gorm:"column:commission_type" json:"commission_type"`
	CommissionRate        float64   `gorm:"column:commission_rate" json:"commission_rate"`
	Level                 uint8     `gorm:"column:level" json:"level"`
	Remark                string    `gorm:"column:remark" json:"remark"`
	Sort                  int       `gorm:"column:sort" json:"sort"`
	Status                int8      `gorm:"column:status" json:"status"`
	Type                  int8      `gorm:"column:type" json:"type"`
	RoleID                uint      `gorm:"column:role_id" json:"role_id"`
	BusinessStoreCategory uint      `gorm:"column:business_store_category" json:"business_store_category"`
	BusinessStoreType     uint      `gorm:"column:business_store_type" json:"business_store_type"`
	GoodsType             string    `gorm:"column:goods_type" json:"goods_type"`
	PlatformCategoryIDs   string    `gorm:"column:platform_category_ids" json:"platform_category_ids"`
	CreateTime            time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime            time.Time `gorm:"column:update_time" json:"update_time"`
}

func (Circle) TableName() string { return "qixi_crm_a_business_zone" }

// Agent 区域或商户型商圈的代理申请与结算信息。结算账号资料只写入，不经统一后台列表或详情接口回传。
// Extend 对齐 CRMEB extend（JSON）；其中 avatar 为「区域代理」标识图。
type Agent struct {
	CircleAgentID         uint       `gorm:"column:circle_agent_id;primaryKey" json:"circle_agent_id"`
	UID                   uint       `gorm:"column:uid" json:"uid"`
	Name                  string     `gorm:"column:name" json:"name"`
	Phone                 string     `gorm:"column:phone" json:"phone"`
	Qualification         string     `gorm:"column:qualification" json:"qualification"`
	Remark                string     `gorm:"column:remark" json:"remark"`
	Extend                string     `gorm:"column:extend" json:"extend"`
	AuditAdminID          uint       `gorm:"column:audit_admin_id" json:"audit_admin_id"`
	AuditReason           string     `gorm:"column:audit_reason" json:"audit_reason"`
	AuditTime             *time.Time `gorm:"column:audit_time" json:"audit_time,omitempty"`
	Status                int8       `gorm:"column:status" json:"status"`
	PaymentMethod         uint8      `gorm:"column:payment_method" json:"payment_method"`
	PaymentName           string     `gorm:"column:payment_name" json:"payment_name"`
	PaymentAccount        string     `gorm:"column:payment_account" json:"-"`
	PaymentBank           string     `gorm:"column:payment_bank" json:"-"`
	PaymentQRImg          string     `gorm:"column:payment_qr_img" json:"-"`
	PaymentConfigured     bool       `gorm:"-" json:"payment_configured"`
	Balance               float64    `gorm:"column:balance" json:"balance"`
	Type                  int8       `gorm:"column:type" json:"type"`
	BusinessName          string     `gorm:"column:business_name" json:"business_name"`
	BusinessStoreCategory uint       `gorm:"column:business_store_category" json:"business_store_category"`
	BusinessStoreType     uint       `gorm:"column:business_store_type" json:"business_store_type"`
	CreateTime            time.Time  `gorm:"column:create_time" json:"create_time"`
	UpdateTime            time.Time  `gorm:"column:update_time" json:"update_time"`
}

func (Agent) TableName() string { return "qixi_crm_a_business_zone_agent" }

const (
	AgentPending  int8 = 0
	AgentApproved int8 = 1
	AgentRejected int8 = -1
	AgentRevoked  int8 = -2
)

type CircleListFilter struct {
	Keyword       string
	Status        *int8
	Type          *int8
	CircleAgentID uint
	// PID 非 nil 时按父级筛选：0=一级区域，>0=指定父级的下级。
	PID *uint
}

type CircleInput struct {
	PID                   uint    `json:"pid"`
	Name                  string  `json:"name"`
	CircleAgentID         uint    `json:"circle_agent_id"`
	CommissionType        int8    `json:"commission_type"`
	CommissionRate        float64 `json:"commission_rate"`
	Remark                string  `json:"remark"`
	Sort                  int     `json:"sort"`
	Status                int8    `json:"status"`
	Type                  int8    `json:"type"`
	RoleID                uint    `json:"role_id"`
	BusinessStoreCategory uint    `json:"business_store_category"`
	BusinessStoreType     uint    `json:"business_store_type"`
	GoodsTypes            []int   `json:"goods_type"`
	PlatformCategoryIDs   []uint  `json:"platform_category_ids"`
	MerchantIDs           []uint  `json:"merchant_ids"`
}

type AgentListFilter struct {
	Keyword   string
	Name      string
	Phone     string
	Status    *int8
	Type      *int8
	UID       *uint
	// UIDs 非 nil 时按 uid IN 过滤；空切片表示无匹配（返回空列表）。
	UIDs      *[]uint
	DateFrom  string
	DateTo    string
}

type AgentInput struct {
	UID                   uint   `json:"uid"`
	Name                  string `json:"name"`
	Phone                 string `json:"phone"`
	Qualification         string `json:"qualification"`
	Remark                string `json:"remark"`
	// Avatar 区域代理标识图；写入 extend.avatar（对齐 CRMEB extend）。
	Avatar                string `json:"avatar"`
	Extend                string `json:"extend"`
	PaymentMethod         uint8  `json:"payment_method"`
	PaymentName           string `json:"payment_name"`
	PaymentAccount        string `json:"payment_account"`
	PaymentBank           string `json:"payment_bank"`
	PaymentQRImg          string `json:"payment_qr_img"`
	Type                  int8   `json:"type"`
	BusinessName          string `json:"business_name"`
	BusinessStoreCategory uint   `json:"business_store_category"`
	BusinessStoreType     uint   `json:"business_store_type"`
	Account               string `json:"account"`
	Password              string `json:"password"`
	// AutoApprove 区域表单内「添加代理人」时立即通过，便于立刻绑定区域。
	AutoApprove bool `json:"auto_approve"`
	// CircleIDs 负责商户（type=1 商圈）多选；nil 表示不改绑定。
	CircleIDs *[]uint `json:"circle_ids"`
}

type AuditInput struct {
	Status      int8   `json:"status"`
	AuditReason string `json:"audit_reason"`
}
