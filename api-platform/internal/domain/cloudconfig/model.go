package cloudconfig

import "time"

// Config 加密存储的云服务配置项；value 永不以明文形式落库。
type Config struct {
	ConfigID   uint      `gorm:"column:config_id;primaryKey" json:"config_id"`
	GroupKey   string    `gorm:"column:group_key" json:"group_key"`
	ConfigKey  string    `gorm:"column:config_key" json:"config_key"`
	Ciphertext string    `gorm:"column:ciphertext" json:"-"`
	IsSecret   bool      `gorm:"column:is_secret" json:"is_secret"`
	UpdatedBy  uint      `gorm:"column:updated_by" json:"updated_by"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (Config) TableName() string { return "qixi_m_admin_cloud_config" }

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
