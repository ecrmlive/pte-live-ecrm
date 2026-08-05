package chat

import "time"

const (
	RoleUser    = "user"
	RoleService = "service"

	MsgText   = "text"
	MsgOrder  = "order"
	MsgSystem = "system"

	StatusOpen   = "open"
	StatusClosed = "closed"
)

type Thread struct {
	ThreadID         uint       `gorm:"column:id;primaryKey" json:"thread_id"`
	MerID            uint       `gorm:"column:mer_id;->" json:"mer_id"`
	StoreID          uint       `gorm:"column:store_id" json:"store_id"`
	UID              uint       `gorm:"column:user_id" json:"uid"`
	ServiceID        uint       `gorm:"column:assigned_admin_id" json:"service_id"`
	AssignedAt       *time.Time `gorm:"column:assigned_at" json:"assigned_at,omitempty"`
	ImConversationID string     `gorm:"column:im_conversation_id" json:"im_conversation_id"`
	LastMsg          string     `gorm:"column:last_msg" json:"last_msg"`
	LastTime         *time.Time `gorm:"column:last_time" json:"last_time"`
	UserUnread       uint       `gorm:"column:user_unread" json:"user_unread"`
	ServiceUnread    uint       `gorm:"column:service_unread" json:"service_unread"`
	Status           string     `gorm:"column:status" json:"status"`
	CreateTime       time.Time  `gorm:"column:created_at" json:"create_time"`

	UserNickname string `gorm:"-" json:"user_nickname,omitempty"`
	MerName      string `gorm:"-" json:"mer_name,omitempty"`
}

func (Thread) TableName() string { return "qixi_crm_b_customer_service_binding" }

type Message struct {
	MsgID      uint64    `gorm:"column:id;primaryKey" json:"msg_id"`
	ThreadID   uint      `gorm:"column:binding_id" json:"thread_id"`
	MerID      uint      `gorm:"column:merchant_id" json:"mer_id"`
	SenderRole string    `gorm:"column:sender_role" json:"sender_role"`
	SenderID   uint      `gorm:"column:sender_id" json:"sender_id"`
	MsgType    string    `gorm:"column:msg_type" json:"msg_type"`
	Content    string    `gorm:"column:content" json:"content"`
	CreateTime time.Time `gorm:"column:created_at" json:"create_time"`
}

func (Message) TableName() string { return "qixi_crm_b_customer_service_message" }

type Identity struct {
	ID         uint      `gorm:"column:id;primaryKey" json:"id"`
	Portal     string    `gorm:"column:portal" json:"portal"`
	LocalID    uint      `gorm:"column:local_id" json:"local_id"`
	ImUserID   string    `gorm:"column:im_user_id" json:"im_user_id"`
	ImUserNum  int64     `gorm:"column:im_user_num" json:"im_user_num"`
	CreateTime time.Time `gorm:"column:created_at" json:"create_time"`
}

func (Identity) TableName() string { return "qixi_crm_b_im_identity" }

type SendInput struct {
	MsgType string `json:"msg_type"`
	Content string `json:"content"`
}

type OpenInput struct {
	MerID uint `json:"mer_id"`
}

type Credential struct {
	Mode             string `json:"mode"`
	AppID            string `json:"app_id"`
	SDKAppID         string `json:"sdk_app_id,omitempty"`
	ImUserID         string `json:"im_user_id"`
	Identifier       string `json:"identifier,omitempty"`
	UserSig          string `json:"user_sig"`
	ExpireAt         int64  `json:"expire_at"`
	APIURL           string `json:"api_url,omitempty"`
	WSURL            string `json:"ws_url,omitempty"`
	WSHint           string `json:"ws_hint,omitempty"`
	ImConversationID uint64 `json:"im_conversation_id,omitempty"`
	Note             string `json:"note,omitempty"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
