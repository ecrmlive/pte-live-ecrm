package seckill

import "time"

type TimeSlot struct {
	SeckillTimeID uint      `gorm:"column:seckill_time_id;primaryKey" json:"seckill_time_id"`
	Title         string    `gorm:"column:title" json:"title"`
	StartTime     int       `gorm:"column:start_time" json:"start_time"`
	EndTime       int       `gorm:"column:end_time" json:"end_time"`
	Status        int8      `gorm:"column:status" json:"status"`
	CreateTime    time.Time `gorm:"column:create_time" json:"create_time"`
	Pic           string    `gorm:"column:pic" json:"pic"`
}

func (TimeSlot) TableName() string { return "qixi_crm_b_seckill_time" }

type Active struct {
	SeckillActiveID uint    `gorm:"column:seckill_active_id;primaryKey" json:"seckill_active_id"`
	ActivityID      uint    `gorm:"column:activity_id" json:"activity_id"`
	Name            string  `gorm:"column:name" json:"name"`
	SeckillTimeIDs  string  `gorm:"column:seckill_time_ids" json:"seckill_time_ids"`
	StartDay        string  `gorm:"column:start_day" json:"start_day"`
	EndDay          string  `gorm:"column:end_day" json:"end_day"`
	MerID           uint    `gorm:"column:mer_id" json:"mer_id"`
	ProductID       uint    `gorm:"column:product_id" json:"product_id"`
	SeckillPrice    float64 `gorm:"column:seckill_price" json:"seckill_price"`
	OncePayCount    int     `gorm:"column:once_pay_count" json:"once_pay_count"`
	AllPayCount     int     `gorm:"column:all_pay_count" json:"all_pay_count"`
	ActiveStatus    int8    `gorm:"column:active_status" json:"active_status"`
	Status          int8    `gorm:"column:status" json:"status"`
	IsShow          int8    `gorm:"column:is_show" json:"is_show"`
	ProductStatus   int8    `gorm:"column:product_status" json:"product_status"`
	Star            int8    `gorm:"column:star" json:"star"`
	Sort            int     `gorm:"column:sort" json:"sort"`
	Stock           int     `gorm:"column:stock" json:"stock"`
	Sales           int     `gorm:"column:sales" json:"sales"`
	SysLabels       string  `gorm:"column:sys_labels" json:"sys_labels"`
	Refusal         string  `gorm:"column:refusal" json:"refusal"`
	CreateTime      int64   `gorm:"column:create_time" json:"create_time"`
	UpdateTime      int64   `gorm:"column:update_time" json:"update_time"`
	DeleteTime      *int64  `gorm:"column:delete_time" json:"-"`

	StoreName          string  `gorm:"-" json:"store_name,omitempty"`
	Image              string  `gorm:"-" json:"image,omitempty"`
	Price              float64 `gorm:"-" json:"price,omitempty"`
	MerName            string  `gorm:"-" json:"mer_name,omitempty"`
	IsTrader           int8    `gorm:"-" json:"is_trader"`
	TraderName         string  `gorm:"-" json:"trader_name,omitempty"`
	TimeTitles         string  `gorm:"-" json:"time_titles,omitempty"`
	ProductStatusName  string  `gorm:"-" json:"product_status_name"`
	ActivityStatus     int8    `gorm:"-" json:"activity_status"` // 0未开始 1进行中 -1已结束
	ActivityStatusText string  `gorm:"-" json:"activity_status_text"`
	InWindow           bool    `gorm:"-" json:"in_window"`
}

func (Active) TableName() string { return "qixi_crm_b_seckill_active" }

type ActiveInput struct {
	Name           string  `json:"name"`
	SeckillTimeIDs string  `json:"seckill_time_ids"`
	StartDay       string  `json:"start_day"`
	EndDay         string  `json:"end_day"`
	ProductID      uint    `json:"product_id"`
	SeckillPrice   float64 `json:"seckill_price"`
	OncePayCount   int     `json:"once_pay_count"`
	AllPayCount    *int    `json:"all_pay_count"`
	Status         *int8   `json:"status"`
	IsShow         *int8   `json:"is_show"`
	ProductStatus  *int8   `json:"product_status"`
	Star           *int8   `json:"star"`
	Sort           *int    `json:"sort"`
	Stock          *int    `json:"stock"`
	SysLabels      *string `json:"sys_labels"`
	Refusal        string  `json:"refusal"`
}

// TimeSlotInput 秒杀配置（场次）写入。
type TimeSlotInput struct {
	Title     string `json:"title"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	Status    *int8  `json:"status"`
	Pic       string `json:"pic"`
}

type TimeSlotQuery struct {
	Status *int8
	Page   int
	Limit  int
}

// ActiveAdminQuery 秒杀管理列表筛选（对齐 CRMEB StoreProductSeckill）。
type ActiveAdminQuery struct {
	Type       int // 1出售中 2仓库中 5回收站 6待审核 7审核未通过
	MerID      *uint
	MerIDs     []uint // 店铺类别等预解析后的商户 ID
	ActiveName string
	Keyword    string
	IsTrader   *int8
	Star       *int8
	UsStatus   *int8 // 活动状态投影
	SysLabels  string
	Page       int
	Limit      int
}

type StatusFilterItem struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// Activity 平台秒杀活动（对齐 CRMEB store_seckill_active）。
type Activity struct {
	SeckillActivityID  uint       `gorm:"column:seckill_activity_id;primaryKey" json:"seckill_activity_id"`
	Name               string     `gorm:"column:name" json:"name"`
	SeckillTimeIDs     string     `gorm:"column:seckill_time_ids" json:"seckill_time_ids"`
	StartDay           string     `gorm:"column:start_day" json:"start_day"`
	EndDay             string     `gorm:"column:end_day" json:"end_day"`
	OncePayCount       int        `gorm:"column:once_pay_count" json:"once_pay_count"`
	AllPayCount        int        `gorm:"column:all_pay_count" json:"all_pay_count"`
	ProductCategoryIDs string     `gorm:"column:product_category_ids" json:"product_category_ids"`
	BorderPic          string     `gorm:"column:border_pic" json:"border_pic"`
	Status             int8       `gorm:"column:status" json:"status"`
	ActiveStatus       int8       `gorm:"column:active_status" json:"active_status"`
	ProductCount       int        `gorm:"column:product_count" json:"product_count"`
	MerchantCount      int        `gorm:"column:merchant_count" json:"merchant_count"`
	CreateTime         time.Time  `gorm:"column:create_time" json:"create_time"`
	UpdateTime         time.Time  `gorm:"column:update_time" json:"update_time"`
	DeleteTime         *time.Time `gorm:"column:delete_time" json:"-"`

	StatusText       string   `gorm:"-" json:"status_text"`
	SeckillTimeTexts []string `gorm:"-" json:"seckill_time_texts"`
}

func (Activity) TableName() string { return "qixi_crm_b_seckill_activity" }

type ActivityInput struct {
	Name               string `json:"name"`
	SeckillTimeIDs     string `json:"seckill_time_ids"`
	StartDay           string `json:"start_day"`
	EndDay             string `json:"end_day"`
	OncePayCount       int    `json:"once_pay_count"`
	AllPayCount        int    `json:"all_pay_count"`
	ProductCategoryIDs string `json:"product_category_ids"`
	BorderPic          string `json:"border_pic"`
	Status             *int8  `json:"status"`
}

// ActivityProductQuery 活动「已加商品」分页筛选。
type ActivityProductQuery struct {
	ActivityID    uint
	ProductStatus *int8 // 0待审 1通过 -1未通过；nil 不过滤
	Keyword       string
	Page          int
	Limit         int
}

// ActivityProductSKU 展开行（优先 qixi_crm_b_product_sku_view；无则商品级「默认」）。
type ActivityProductSKU struct {
	SKU          string  `json:"sku"`
	Image        string  `json:"image,omitempty"`
	Price        float64 `json:"price,omitempty"` // 规格售价
	SeckillPrice float64 `json:"seckill_price"`
	Stock        int     `json:"stock"`       // 规格/商品库存
	LimitStock   int     `json:"limit_stock"` // 秒杀限量剩余
}

// ProductSKURow 商品规格投影（business product_sku_view）。
type ProductSKURow struct {
	MerchantSKUID uint
	SKUKey        string
	SpecSnapshot  string
	Price         float64
	Stock         int
}

// ActivityProductSaveItem 编辑 Drawer「秒杀商品」草稿行写入。
type ActivityProductSaveItem struct {
	ProductID    uint    `json:"product_id"`
	SeckillPrice float64 `json:"seckill_price"`
	Stock        int     `json:"stock"` // 限量
	Status       *int8   `json:"status"`
	Sort         int     `json:"sort"`
}

// ActivityProductsSaveInput 批量挂载/更新活动秒杀商品。
type ActivityProductsSaveInput struct {
	Products []ActivityProductSaveItem `json:"products"`
}

// ActivityProductItem 活动已挂载秒杀商品（查看 Drawer「已加商品」）。
type ActivityProductItem struct {
	SeckillActiveID    uint                 `json:"seckill_active_id"`
	ProductID          uint                 `json:"product_id"`
	Name               string               `json:"name"`
	StoreName          string               `json:"store_name"`
	Image              string               `json:"image"`
	CateName           string               `json:"cate_name"`
	MerID              uint                 `json:"mer_id"`
	MerName            string               `json:"mer_name"`
	Price              float64              `json:"price"`
	SeckillPrice       float64              `json:"seckill_price"`
	ProductStock       int                  `json:"product_stock"`
	Stock              int                  `json:"stock"`
	Sort               int                  `json:"sort"`
	ProductStatus      int8                 `json:"product_status"`
	ProductStatusName  string               `json:"product_status_name"`
	Children           []ActivityProductSKU `json:"children"`
}

type ActivityQuery struct {
	Name         string
	DateFrom     string
	DateTo       string
	ActiveStatus *int8 // 0未开始 1进行中 -1已结束
	Status       *int8 // 开关
	Page         int
	Limit        int
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// ActivityProductRow 活动商品统计行（对齐 CRMEB chartProduct）
type ActivityProductRow struct {
	SeckillActiveID  uint     `json:"seckill_active_id"`
	ProductID        uint     `json:"product_id"`
	Name             string   `json:"name"`
	Image            string   `json:"image,omitempty"`
	CategoryName     string   `json:"category_name,omitempty"`
	MerID            uint     `json:"mer_id"`
	MerName          string   `json:"mer_name,omitempty"`
	Price            float64  `json:"price"`
	SeckillPrice     float64  `json:"seckill_price"`
	Stock            int      `json:"stock"`
	Sales            int      `json:"sales"`
	SeckillTimeTexts []string `json:"seckill_time_texts,omitempty"`
}

// ActivityStats 活动 KPI 面板（对齐 CRMEB chartPanel）
type ActivityStats struct {
	SeckillActivityID   uint    `json:"seckill_activity_id"`
	Name                string  `json:"name"`
	OrdersPeopleCount   int64   `json:"orders_people_count"`
	PayOrderMoney       float64 `json:"pay_order_money"`
	PayOrderPeopleCount int64   `json:"pay_order_people_count"`
	PayOrderCount       int64   `json:"pay_order_count"`
}

// ActivityStatPeople 活动参与人（投影表；订单链路未接前供演示）
type ActivityStatPeople struct {
	UID          uint      `gorm:"column:uid" json:"uid"`
	Nickname     string    `gorm:"column:nickname" json:"nickname"`
	Phone        string    `gorm:"column:phone" json:"phone"`
	MerID        uint      `gorm:"column:mer_id" json:"mer_id"`
	SumTotalNum  int       `gorm:"column:sum_total_num" json:"sum_total_num"`
	OrderCount   int       `gorm:"column:order_count" json:"order_count"`
	SumPayPrice  float64   `gorm:"column:sum_pay_price" json:"sum_pay_price"`
	LastJoinTime time.Time `gorm:"column:last_join_time" json:"last_join_time"`
}

func (ActivityStatPeople) TableName() string {
	return "qixi_crm_b_seckill_activity_stat_people"
}

// ActivityStatOrder 活动订单（投影表）
type ActivityStatOrder struct {
	OrderSN    string     `gorm:"column:order_sn" json:"order_sn"`
	UID        uint       `gorm:"column:uid" json:"uid"`
	Nickname   string     `gorm:"column:nickname" json:"nickname"`
	MerID      uint       `gorm:"column:mer_id" json:"mer_id"`
	Status     int8       `gorm:"column:status" json:"status"`
	StatusText string     `gorm:"column:status_text" json:"status_text"`
	PayPrice   float64    `gorm:"column:pay_price" json:"pay_price"`
	TotalNum   int        `gorm:"column:total_num" json:"total_num"`
	Paid       int8       `gorm:"column:paid" json:"paid"`
	CreateTime time.Time  `gorm:"column:create_time" json:"create_time"`
	PayTime    *time.Time `gorm:"column:pay_time" json:"pay_time,omitempty"`
}

func (ActivityStatOrder) TableName() string {
	return "qixi_crm_b_seckill_activity_stat_order"
}

type ActivityStatQuery struct {
	Keyword  string
	DateFrom string
	DateTo   string
	MerID    *uint
	Status   *int8 // 订单状态
	Page     int
	Limit    int
}
