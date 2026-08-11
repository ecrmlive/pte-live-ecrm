package cloudconfig

import "time"

// Config 是后台云服务配置项。后台修改后的密钥使用密文保存；仅被 Git 忽略的
// local/test *_key.sql 可在首次初始化时写入受控引导值。
type Config struct {
	GroupKey   string    `gorm:"column:provider;primaryKey" json:"group_key"`
	ConfigKey  string    `gorm:"column:config_key;primaryKey" json:"config_key"`
	Ciphertext string    `gorm:"column:ciphertext" json:"-"`
	KeyVersion string    `gorm:"column:key_version" json:"-"`
	IsSecret   bool      `gorm:"-" json:"is_secret"`
	UpdatedBy  uint      `gorm:"column:updated_by" json:"updated_by"`
	UpdateTime time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Config) TableName() string { return "qixi_crm_a_cloud_config" }

type FieldMeta struct {
	Key       string   `json:"key"`
	Label     string   `json:"label"`
	Secret    bool     `json:"secret"`
	Required  bool     `json:"required"`
	Hint      string   `json:"hint,omitempty"`
	InputType string   `json:"input_type,omitempty"`
	Options   []string `json:"options,omitempty"`
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
