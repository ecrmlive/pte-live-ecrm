package identity

import "time"

type SystemAdmin struct {
	AdminID       uint       `gorm:"column:admin_id;primaryKey" json:"admin_id"`
	Account       string     `gorm:"column:account" json:"account"`
	Pwd           string     `gorm:"column:pwd" json:"-"`
	RealName      string     `gorm:"column:real_name" json:"real_name"`
	Phone         string     `gorm:"column:phone" json:"phone"`
	Roles         string     `gorm:"column:roles" json:"roles"`
	LastIP        string     `gorm:"column:last_ip" json:"last_ip"`
	LastTime      *time.Time `gorm:"column:last_time" json:"last_time"`
	LoginCount    uint       `gorm:"column:login_count" json:"login_count"`
	Status        uint8      `gorm:"column:status" json:"status"`
	Level         uint8      `gorm:"column:level" json:"level"`
	RegionIDs     string     `gorm:"column:region_ids" json:"region_ids"`
	IsAgent       int8       `gorm:"column:is_agent" json:"is_agent"`
	CircleAgentID uint       `gorm:"column:circle_agent_id" json:"circle_agent_id"`
	IsDel         uint8      `gorm:"column:is_del" json:"-"`
}

func (SystemAdmin) TableName() string { return "qixi_m_admin_system_admin" }

type MerchantAdmin struct {
	MerchantAdminID uint       `gorm:"column:merchant_admin_id;primaryKey" json:"merchant_admin_id"`
	MerID           uint       `gorm:"column:mer_id" json:"mer_id"`
	Account         string     `gorm:"column:account" json:"account"`
	Pwd             string     `gorm:"column:pwd" json:"-"`
	RealName        string     `gorm:"column:real_name" json:"real_name"`
	Phone           string     `gorm:"column:phone" json:"phone"`
	Roles           string     `gorm:"column:roles" json:"roles"`
	LastIP          string     `gorm:"column:last_ip" json:"last_ip"`
	LastTime        *time.Time `gorm:"column:last_time" json:"last_time"`
	LoginCount      uint       `gorm:"column:login_count" json:"login_count"`
	Level           uint8      `gorm:"column:level" json:"level"`
	IsDel           uint8      `gorm:"column:is_del" json:"-"`
	Status          uint8      `gorm:"column:status" json:"status"`
}

func (MerchantAdmin) TableName() string { return "qixi_m_admin_merchant_admin" }

type Merchant struct {
	MerID    uint   `gorm:"column:mer_id;primaryKey" json:"mer_id"`
	MerName  string `gorm:"column:mer_name" json:"mer_name"`
	Status   int8   `gorm:"column:status" json:"status"`
	MerState int8   `gorm:"column:mer_state" json:"mer_state"`
	IsDel    int8   `gorm:"column:is_del" json:"-"`
}

func (Merchant) TableName() string { return "qixi_m_admin_merchant" }

// User C 端用户（登录/注册所需字段）
type User struct {
	UID           uint       `gorm:"column:uid;primaryKey" json:"uid"`
	Account       string     `gorm:"column:account" json:"account"`
	Pwd           string     `gorm:"column:pwd" json:"-"`
	Nickname      string     `gorm:"column:nickname" json:"nickname"`
	Avatar        string     `gorm:"column:avatar" json:"avatar"`
	Phone         string     `gorm:"column:phone" json:"phone"`
	NowMoney      float64    `gorm:"column:now_money" json:"now_money"`
	Integral      int        `gorm:"column:integral" json:"integral"`
	IsSvip        int8       `gorm:"column:is_svip" json:"is_svip"`
	SvipEndtime   *time.Time `gorm:"column:svip_endtime" json:"svip_endtime"`
	SvipSaveMoney float64    `gorm:"column:svip_save_money" json:"svip_save_money"`
	Status        int8       `gorm:"column:status" json:"status"`
	UserType      string     `gorm:"column:user_type" json:"user_type"`
	LastIP        string     `gorm:"column:last_ip" json:"last_ip"`
	LastTime      *time.Time `gorm:"column:last_time" json:"last_time"`
}

func (User) TableName() string { return "qixi_m_app_user" }

// StoreService 店员/客服（manager 门户登录）
type StoreService struct {
	ServiceID  uint      `gorm:"column:service_id;primaryKey" json:"service_id"`
	MerID      uint      `gorm:"column:mer_id" json:"mer_id"`
	UID        uint      `gorm:"column:uid" json:"uid"`
	Avatar     string    `gorm:"column:avatar" json:"avatar"`
	Nickname   string    `gorm:"column:nickname" json:"nickname"`
	Account    string    `gorm:"column:account" json:"account"`
	Pwd        string    `gorm:"column:pwd" json:"-"`
	IsOpen     int8      `gorm:"column:is_open" json:"is_open"`
	Status     int8      `gorm:"column:status" json:"status"`
	Customer   int8      `gorm:"column:customer" json:"customer"`
	IsVerify   int8      `gorm:"column:is_verify" json:"is_verify"`
	IsGoods    int8      `gorm:"column:is_goods" json:"is_goods"`
	Phone      string    `gorm:"column:phone" json:"phone"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	IsDel      int8      `gorm:"column:is_del" json:"-"`
}

func (StoreService) TableName() string { return "qixi_m_admin_store_service" }

type SystemRole struct {
	RoleID   uint   `gorm:"column:role_id;primaryKey" json:"role_id"`
	RoleName string `gorm:"column:role_name" json:"role_name"`
	Rules    string `gorm:"column:rules" json:"rules"`
	Status   uint8  `gorm:"column:status" json:"status"`
	MerID    uint   `gorm:"column:mer_id" json:"mer_id"`
	IsAgent  int8   `gorm:"column:is_agent" json:"is_agent"`
	CircleID uint   `gorm:"column:circle_id" json:"circle_id"`
}

func (SystemRole) TableName() string { return "qixi_m_admin_system_role" }

type SystemMenu struct {
	MenuID   uint   `gorm:"column:menu_id;primaryKey" json:"menu_id"`
	PID      uint   `gorm:"column:pid" json:"pid"`
	Path     string `gorm:"column:path" json:"path"`
	Icon     string `gorm:"column:icon" json:"icon"`
	MenuName string `gorm:"column:menu_name" json:"menu_name"`
	Route    string `gorm:"column:route" json:"route"`
	Params   string `gorm:"column:params" json:"params"`
	Sort     int8   `gorm:"column:sort" json:"sort"`
	IsShow   uint8  `gorm:"column:is_show" json:"is_show"`
	IsMer    uint8  `gorm:"column:is_mer" json:"is_mer"`
	IsMenu   uint8  `gorm:"column:is_menu" json:"is_menu"`
	IsAgent  int8   `gorm:"column:is_agent" json:"is_agent"`
}

func (SystemMenu) TableName() string { return "qixi_m_admin_system_menu" }

// MenuNode 前端侧栏/角色勾选树（含 is_menu=2 按钮节点）
type MenuNode struct {
	MenuID   uint        `json:"menu_id"`
	PID      uint        `json:"pid"`
	Path     string      `json:"path"`
	Icon     string      `json:"icon"`
	MenuName string      `json:"menu_name"`
	Route    string      `json:"route"`
	Sort     int8        `json:"sort"`
	IsMenu   uint8       `json:"is_menu"`
	Children []*MenuNode `json:"children,omitempty"`
}
