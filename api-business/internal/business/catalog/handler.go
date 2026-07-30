package catalog

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

// Handler 是六端共用的商品消费视图读取接口。它只读取由店铺域经事件同步的
// qixi_crm_b_product_view，禁止 C 端跨库读取 qixi_crm_m_product。
type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/catalog/home", h.Home)
	r.GET("/catalog/categories", h.Categories)
	r.GET("/catalog/products", h.Products)
	r.GET("/catalog/products/:id", h.ProductDetail)
	r.GET("/catalog/stores/:id", h.StoreHome)
}

type productView struct {
	ProductID     uint64   `gorm:"column:product_id"`
	MerchantID    uint64   `gorm:"column:merchant_id"`
	StoreID       uint64   `gorm:"column:store_id"`
	CategoryID    uint64   `gorm:"column:category_id"`
	MerchantName  string   `gorm:"column:merchant_name"`
	StoreName     string   `gorm:"column:store_name"`
	Title         string   `gorm:"column:title"`
	CoverURL      string   `gorm:"column:cover_url"`
	Price         float64  `gorm:"column:price"`
	OriginalPrice *float64 `gorm:"column:original_price"`
	Sales         int      `gorm:"column:sales"`
	Stock         int      `gorm:"column:stock"`
	SaleStatus    int8     `gorm:"column:sale_status"`
}

func (productView) TableName() string { return "qixi_crm_b_product_view" }

type storeView struct {
	StoreID       uint64 `gorm:"column:store_id"`
	MerchantID    uint64 `gorm:"column:merchant_id"`
	MerchantAppID string `gorm:"column:store_app_id"`
	StoreName     string `gorm:"column:store_name"`
}

func (storeView) TableName() string { return "qixi_crm_b_store_view" }

type categoryView struct {
	CategoryID uint64 `gorm:"column:category_id"`
	ParentID   uint64 `gorm:"column:parent_id"`
	Name       string `gorm:"column:name"`
	Sort       int    `gorm:"column:sort"`
	Status     int8   `gorm:"column:status"`
}

func (categoryView) TableName() string { return "qixi_crm_b_category_view" }

func (h *Handler) Home(c *gin.Context) {
	rows, err := h.list(c, 0, 0, "", 1, 12, "default", "desc")
	if err != nil {
		fail(c, err)
		return
	}
	hot := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		hot = append(hot, responseProduct(row))
	}
	response.OK(c, gin.H{"diy_id": 0, "diy_title": "", "banners": []gin.H{}, "menus": []gin.H{}, "hot": hot})
}

func (h *Handler) Categories(c *gin.Context) {
	rows := make([]categoryView, 0)
	if err := h.db.WithContext(c.Request.Context()).Where("status = ?", 1).Order("sort ASC,category_id ASC").Find(&rows).Error; err != nil {
		fail(c, err)
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{"id": row.CategoryID, "pid": row.ParentID, "name": row.Name})
	}
	response.OK(c, items)
}

func (h *Handler) Products(c *gin.Context) {
	page, limit := pageParams(c)
	merchantID, _ := strconv.ParseUint(c.Query("mer_id"), 10, 64)
	categoryID, err := queryUint(c, "cate_id")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "分类 ID 错误")
		return
	}
	sort, order := productSort(c.Query("sort"), c.Query("order"))
	rows, total, err := h.listPage(c, merchantID, categoryID, strings.TrimSpace(c.Query("keyword")), page, limit, sort, order)
	if err != nil {
		fail(c, err)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, responseProduct(row))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) ProductDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "商品 ID 错误")
		return
	}
	var row productView
	err = h.db.WithContext(c.Request.Context()).Where("product_id = ? AND sale_status = 1", id).First(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "商品不存在")
			return
		}
		fail(c, err)
		return
	}
	item := responseProduct(row)
	item["unit_name"] = "件"
	item["store_info"] = row.StoreName
	if row.CoverURL == "" {
		item["slider_image"] = []string{}
	} else {
		item["slider_image"] = []string{row.CoverURL}
	}
	item["spec_type"] = 0
	item["delivery_way"] = "express"
	store, err := h.findStoreView(c, row.MerchantID)
	if err != nil {
		fail(c, err)
		return
	}
	if store != nil {
		item["merchant_app_id"] = store.MerchantAppID
	}
	response.OK(c, item)
}

func (h *Handler) StoreHome(c *gin.Context) {
	merchantID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || merchantID == 0 {
		response.Fail(c, http.StatusBadRequest, "商户 ID 错误")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	rows, total, err := h.listPage(c, merchantID, 0, "", page, 20, "default", "desc")
	if err != nil {
		fail(c, err)
		return
	}
	items := make([]gin.H, 0, len(rows))
	merchantName := ""
	for _, row := range rows {
		if merchantName == "" {
			merchantName = row.MerchantName
		}
		items = append(items, responseProduct(row))
	}
	store, err := h.findStoreView(c, merchantID)
	if err != nil {
		fail(c, err)
		return
	}
	if merchantName == "" {
		if store != nil {
			merchantName = store.StoreName
		}
	}
	storeID, merchantAppID := uint64(0), ""
	if store != nil {
		storeID, merchantAppID = store.StoreID, store.MerchantAppID
	}
	response.OK(c, gin.H{
		"mer_id": merchantID, "mer_name": merchantName,
		"store_id": storeID, "merchant_app_id": merchantAppID,
		"products": items, "total": total,
	})
}

// findStoreView returns no value for the platform's self-operated goods
// (merchant_id=0), and never makes the C-end API reach into merchant storage.
func (h *Handler) findStoreView(c *gin.Context, merchantID uint64) (*storeView, error) {
	if merchantID == 0 {
		return nil, nil
	}
	var store storeView
	err := h.db.WithContext(c.Request.Context()).
		Where("merchant_id = ? AND status = 1", merchantID).
		First(&store).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (h *Handler) list(c *gin.Context, merchantID, categoryID uint64, keyword string, page, limit int, sort, order string) ([]productView, error) {
	rows, _, err := h.listPage(c, merchantID, categoryID, keyword, page, limit, sort, order)
	return rows, err
}
func (h *Handler) listPage(c *gin.Context, merchantID, categoryID uint64, keyword string, page, limit int, sort, order string) ([]productView, int64, error) {
	query := h.db.WithContext(c.Request.Context()).Model(&productView{}).Where("sale_status = 1")
	if merchantID != 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if categoryID != 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	rows := make([]productView, 0)
	err := query.Order(productOrder(sort, order)).Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func queryUint(c *gin.Context, key string) (uint64, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0, nil
	}
	return strconv.ParseUint(raw, 10, 64)
}

// productSort is whitelisted before values reach GORM's Order expression.
func productSort(sort, order string) (string, string) {
	if sort != "sales" && sort != "price" {
		sort = "default"
	}
	if order != "asc" {
		order = "desc"
	}
	return sort, order
}

func productOrder(sort, order string) string {
	sort, order = productSort(sort, order)
	switch sort {
	case "sales":
		return "sales " + strings.ToUpper(order) + ",updated_at DESC,product_id DESC"
	case "price":
		return "price " + strings.ToUpper(order) + ",updated_at DESC,product_id DESC"
	default:
		return "sales DESC,updated_at DESC,product_id DESC"
	}
}
func responseProduct(row productView) gin.H {
	original := ""
	if row.OriginalPrice != nil {
		original = fmt.Sprintf("%.2f", *row.OriginalPrice)
	}
	// store_name 是旧 C 端契约中实际承载商品标题的字段；保留其兼容值，
	// 新页面使用 title，真实店铺名称使用 shop_name。
	return gin.H{"id": row.ProductID, "title": row.Title, "mer_id": row.MerchantID, "mer_name": row.MerchantName, "store_id": row.StoreID, "store_name": row.Title, "shop_name": row.StoreName, "image": row.CoverURL, "price": fmt.Sprintf("%.2f", row.Price), "ot_price": original, "sales": row.Sales, "stock": row.Stock}
}
func pageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
func fail(c *gin.Context, err error) {
	response.Fail(c, http.StatusInternalServerError, "商品消费视图查询失败")
}
