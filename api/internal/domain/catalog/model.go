package catalog

import "time"

type Category struct {
	StoreCategoryID uint   `gorm:"column:store_category_id;primaryKey" json:"store_category_id"`
	PID             uint   `gorm:"column:pid" json:"pid"`
	CateName        string `gorm:"column:cate_name" json:"cate_name"`
	Path            string `gorm:"column:path" json:"path"`
	Sort            int    `gorm:"column:sort" json:"sort"`
	Pic             string `gorm:"column:pic" json:"pic"`
	IsShow          int8   `gorm:"column:is_show" json:"is_show"`
	Level           uint   `gorm:"column:level" json:"level"`
	MerID           uint   `gorm:"column:mer_id" json:"mer_id"`
	IsHot           int8   `gorm:"column:is_hot" json:"is_hot"`
	Type            int8   `gorm:"column:type" json:"type"`
}

func (Category) TableName() string { return "qixi_store_category" }

type Brand struct {
	BrandID         uint   `gorm:"column:brand_id;primaryKey" json:"brand_id"`
	BrandCategoryID uint   `gorm:"column:brand_category_id" json:"brand_category_id"`
	BrandName       string `gorm:"column:brand_name" json:"brand_name"`
	Sort            int    `gorm:"column:sort" json:"sort"`
	Pic             string `gorm:"column:pic" json:"pic"`
	IsShow          int8   `gorm:"column:is_show" json:"is_show"`
}

func (Brand) TableName() string { return "qixi_store_brand" }

// ProductTypePoints 积分商城商品（对齐 CRMEB activity/product_type=20）
const ProductTypePoints uint8 = 20

type Product struct {
	ProductID    uint      `gorm:"column:product_id;primaryKey" json:"product_id"`
	MerID        uint      `gorm:"column:mer_id" json:"mer_id"`
	StoreName    string    `gorm:"column:store_name" json:"store_name"`
	StoreInfo    string    `gorm:"column:store_info" json:"store_info"`
	Keyword      string    `gorm:"column:keyword" json:"keyword"`
	BrandID      int       `gorm:"column:brand_id" json:"brand_id"`
	IsShow       uint8     `gorm:"column:is_show" json:"is_show"`
	Status       int8      `gorm:"column:status" json:"status"`
	IsDel        uint8     `gorm:"column:is_del" json:"-"`
	MerStatus    int8      `gorm:"column:mer_status" json:"mer_status"`
	CateID       int       `gorm:"column:cate_id" json:"cate_id"`
	UnitName     string    `gorm:"column:unit_name" json:"unit_name"`
	Price         float64   `gorm:"column:price" json:"price"`
	Integral      int       `gorm:"column:integral" json:"integral"`
	OtPrice       float64   `gorm:"column:ot_price" json:"ot_price"`
	SvipPriceType int8      `gorm:"column:svip_price_type" json:"svip_price_type"`
	SvipPrice     float64   `gorm:"column:svip_price" json:"svip_price"`
	MerSvipStatus int8      `gorm:"column:mer_svip_status" json:"mer_svip_status"`
	Stock         uint      `gorm:"column:stock" json:"stock"`
	Sales         uint      `gorm:"column:sales" json:"sales"`
	ProductType   uint8     `gorm:"column:product_type" json:"product_type"`
	SpecType      int8      `gorm:"column:spec_type" json:"spec_type"`
	Refusal       string    `gorm:"column:refusal" json:"refusal"`
	Image         string    `gorm:"column:image" json:"image"`
	SliderImage   string    `gorm:"column:slider_image" json:"slider_image"`
	DeliveryWay   string    `gorm:"column:delivery_way" json:"delivery_way"`
	Type          uint8     `gorm:"column:type" json:"type"`
	CreateTime    time.Time `gorm:"column:create_time" json:"create_time"`
	MerName       string    `gorm:"-" json:"mer_name,omitempty"`
	CateName      string    `gorm:"-" json:"cate_name,omitempty"`
}

func (Product) TableName() string { return "qixi_store_product" }

const (
	ProductStatusPending  int8 = 0
	ProductStatusApproved int8 = 1
	ProductStatusRejected int8 = -1
	ProductStatusOff      int8 = -2
)

type CategoryNode struct {
	StoreCategoryID uint           `json:"store_category_id"`
	PID             uint           `json:"pid"`
	CateName        string         `json:"cate_name"`
	Sort            int            `json:"sort"`
	IsShow          int8           `json:"is_show"`
	Level           uint           `json:"level"`
	Children        []CategoryNode `json:"children,omitempty"`
}

// AttrValue 单规格默认 SKU（阶段2）；多规格后续扩展
type AttrValue struct {
	ValueID   uint    `gorm:"column:value_id;primaryKey" json:"value_id"`
	ProductID uint    `gorm:"column:product_id" json:"product_id"`
	Detail    string  `gorm:"column:detail" json:"detail"`
	SKU       string  `gorm:"column:sku" json:"sku"`
	Stock     uint    `gorm:"column:stock" json:"stock"`
	Sales     uint    `gorm:"column:sales" json:"sales"`
	Image     string  `gorm:"column:image" json:"image"`
	Cost      float64 `gorm:"column:cost" json:"cost"`
	OtPrice   float64 `gorm:"column:ot_price" json:"ot_price"`
	Price     float64 `gorm:"column:price" json:"price"`
	SvipPrice float64 `gorm:"column:svip_price" json:"svip_price"`
	Unique    string  `gorm:"column:unique" json:"unique"`
	IsShow    int8    `gorm:"column:is_show" json:"is_show"`
}

func (AttrValue) TableName() string { return "qixi_store_product_attr_value" }

type ProductSaveInput struct {
	StoreName   string  `json:"store_name"`
	StoreInfo   string  `json:"store_info"`
	Keyword     string  `json:"keyword"`
	CateID      int     `json:"cate_id"`
	BrandID     int     `json:"brand_id"`
	UnitName    string  `json:"unit_name"`
	Price       float64 `json:"price"`
	OtPrice     float64 `json:"ot_price"`
	Cost        float64 `json:"cost"`
	Stock       uint    `json:"stock"`
	Image       string  `json:"image"`
	SliderImage string  `json:"slider_image"`
	DeliveryWay string  `json:"delivery_way"`
	Type          uint8   `json:"type"`
	SpecType      int8    `json:"spec_type"`
	IsShow        *uint8  `json:"is_show"`
	SvipPriceType *int8   `json:"svip_price_type"`
	SvipPrice     *float64 `json:"svip_price"`
	MerSvipStatus *int8   `json:"mer_svip_status"`
}
