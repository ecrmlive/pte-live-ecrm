package cloudconfig

import "time"

// Config 是统一后台拥有的云服务配置项。C 端服务只读该受控配置源，
// 禁止在业务库复制 COS 或支付原始密钥。
type Config struct {
	GroupKey   string    `gorm:"column:provider;primaryKey" json:"group_key"`
	ConfigKey  string    `gorm:"column:config_key;primaryKey" json:"config_key"`
	Ciphertext string    `gorm:"column:ciphertext" json:"-"`
	KeyVersion string    `gorm:"column:key_version" json:"-"`
	IsSecret   bool      `gorm:"-" json:"is_secret"`
	UpdatedBy  uint      `gorm:"column:updated_by" json:"updated_by"`
	UpdateTime time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Config) TableName() string { return "qixi_crm_admin.qixi_crm_a_cloud_config" }

type FieldMeta struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Secret   bool   `json:"secret"`
	Required bool   `json:"required"`
	Hint     string `json:"hint,omitempty"`
}

type GroupMeta struct {
	Key    string      `json:"key"`
	Label  string      `json:"label"`
	Fields []FieldMeta `json:"fields"`
}

type GroupView struct {
	GroupKey   string            `json:"group_key"`
	Label      string            `json:"label"`
	Fields     []FieldMeta       `json:"fields"`
	Values     map[string]string `json:"values"`
	UpdatedAt  *time.Time        `json:"updated_at,omitempty"`
	Configured bool              `json:"configured"`
}

type SaveInput struct {
	Values map[string]string `json:"values"`
}

const SecretMasked = "********"
