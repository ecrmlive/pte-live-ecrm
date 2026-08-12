package notification

import "time"

type Audience string

const (
	AudienceMember Audience = "member"
	AudienceStore  Audience = "store"
)

type Channel string

const (
	ChannelWechat      Channel = "wechat"
	ChannelMiniProgram Channel = "mini_program"
	ChannelSMS         Channel = "sms"
)

// Config 是一个固定业务场景在三种外发渠道上的默认发送策略。
// 文本只定义投递默认值；真正投递仍由对应业务事件和已配置的通道适配器执行。
type Config struct {
	NotificationID uint      `gorm:"column:notification_id;primaryKey" json:"notification_id"`
	Audience       Audience  `gorm:"column:audience" json:"audience"`
	NoticeType     string    `gorm:"column:notice_type" json:"notice_type"`
	Scene          string    `gorm:"column:scene" json:"scene"`
	WechatEnabled  int8      `gorm:"column:wechat_enabled" json:"wechat_enabled"`
	MiniEnabled    int8      `gorm:"column:mini_program_enabled" json:"mini_program_enabled"`
	SMSEnabled     int8      `gorm:"column:sms_enabled" json:"sms_enabled"`
	WechatText     string    `gorm:"column:wechat_text" json:"wechat_text"`
	MiniText       string    `gorm:"column:mini_program_text" json:"mini_program_text"`
	SMSText        string    `gorm:"column:sms_text" json:"sms_text"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Config) TableName() string { return "qixi_crm_a_notification_config" }

type SaveInput struct {
	WechatEnabled int8   `json:"wechat_enabled"`
	MiniEnabled   int8   `json:"mini_program_enabled"`
	SMSEnabled    int8   `json:"sms_enabled"`
	WechatText    string `json:"wechat_text"`
	MiniText      string `json:"mini_program_text"`
	SMSText       string `json:"sms_text"`
}

type PageResult struct {
	List  []Config `json:"list"`
	Total int64    `json:"total"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
}
