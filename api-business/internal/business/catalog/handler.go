package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

// Handler 是六端共用的商品消费视图读取接口。它只读取由店铺域经事件同步的
// qixi_crm_b_product_view，禁止 C 端跨库读取 qixi_crm_m_product。
type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/catalog/home", h.Home)
	r.GET("/catalog/categories", h.Categories)
	r.GET("/catalog/stores", h.Stores)
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

type storeDirectoryItem struct {
	StoreID      uint64 `gorm:"column:store_id"`
	MerchantID   uint64 `gorm:"column:merchant_id"`
	StoreName    string `gorm:"column:store_name"`
	ProductCount int64  `gorm:"column:product_count"`
	SalesCount   int64  `gorm:"column:sales_count"`
	CoverURL     string `gorm:"column:cover_url"`
}

type categoryView struct {
	CategoryID uint64 `gorm:"column:category_id"`
	ParentID   uint64 `gorm:"column:parent_id"`
	Name       string `gorm:"column:name"`
	Sort       int    `gorm:"column:sort"`
	Status     int8   `gorm:"column:status"`
}

func (categoryView) TableName() string { return "qixi_crm_b_category_view" }

func (h *Handler) Home(c *gin.Context) {
	scope, err := h.resolveStoreScope(c)
	if err != nil {
		writeScopeError(c, err)
		return
	}
	rows, err := h.list(c, scope.MerchantID, 0, "", 1, 12, "default", "desc")
	if err != nil {
		fail(c, err)
		return
	}
	hot := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		hot = append(hot, responseProduct(row))
	}
	decoration, err := h.platformHomeDecoration(c)
	if err != nil {
		fail(c, err)
		return
	}
	response.OK(c, gin.H{
		"diy_id":        decoration.PageID,
		"diy_title":     decoration.Title,
		"banners":       decoration.Banners,
		"menus":         decoration.Menus,
		"display_types": decoration.DisplayTypes,
		"hot":           hot,
	})
}

// platformHomeDecoration 读取运营后台已经发布的平台首页装修投影。C 端首页只读
// business 库，不能跨库访问统一后台的装修表。
func (h *Handler) platformHomeDecoration(c *gin.Context) (homeDecoration, error) {
	var row struct {
		PageID   uint64 `gorm:"column:page_id"`
		Name     string `gorm:"column:name"`
		Document string `gorm:"column:document"`
	}
	err := h.db.WithContext(c.Request.Context()).
		Table("qixi_crm_b_diy_page_view").
		Select("page_id,name,document").
		Where("source = ? AND store_id = 0 AND page_type = ? AND status = ? AND is_active = ?", "platform", "home", "published", 1).
		Order("updated_at DESC,page_id DESC").
		First(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return homeDecoration{Banners: []gin.H{}, Menus: []gin.H{}, DisplayTypes: []homeDisplayType{}}, nil
		}
		return homeDecoration{}, err
	}
	return parseHomeDecoration(row.PageID, row.Name, row.Document), nil
}

type homeDecoration struct {
	PageID       uint64
	Title        string
	Banners      []gin.H
	Menus        []gin.H
	DisplayTypes []homeDisplayType
}

// homeDisplayType is sourced from a platform home DIY "商品组" whose source is
// "自动获取". The admin-selected category is the single source of truth for the
// PC, H5 and mini-program home display type; clients only decide when to page.
type homeDisplayType struct {
	CategoryID   uint64 `json:"category_id"`
	InitialLimit int    `json:"initial_limit"`
	Sort         string `json:"sort"`
}

// parseHomeDecoration 保持与 /diy/home 的组件文档兼容，同时过滤未配置图片的
// Banner，避免客户端回退到伪造的渐变占位画面。
func parseHomeDecoration(pageID uint64, name, document string) homeDecoration {
	result := homeDecoration{PageID: pageID, Title: name, Banners: []gin.H{}, Menus: []gin.H{}, DisplayTypes: []homeDisplayType{}}
	var doc struct {
		Page struct {
			Params struct {
				Title string `json:"title"`
			} `json:"params"`
		} `json:"page"`
		Items []struct {
			Type   string `json:"type"`
			Params struct {
				Source string `json:"source"`
				Auto   struct {
					Category    uint64 `json:"category"`
					ProductSort string `json:"productSort"`
					ShowNum     int    `json:"showNum"`
				} `json:"auto"`
			} `json:"params"`
			Data []struct {
				Name  string `json:"imgName"`
				Text  string `json:"text"`
				Image string `json:"imgUrl"`
				URL   string `json:"linkUrl"`
			} `json:"data"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(document), &doc); err != nil {
		return result
	}
	if title := strings.TrimSpace(doc.Page.Params.Title); title != "" {
		result.Title = title
	}
	seenDisplayTypes := make(map[uint64]struct{})
	for _, item := range doc.Items {
		for index, value := range item.Data {
			label := strings.TrimSpace(value.Name)
			if label == "" {
				label = strings.TrimSpace(value.Text)
			}
			switch item.Type {
			case "banner":
				if image := strings.TrimSpace(value.Image); image != "" {
					result.Banners = append(result.Banners, gin.H{"id": index + 1, "title": label, "image": image, "url": value.URL})
				}
			case "navBar":
				if label != "" {
					result.Menus = append(result.Menus, gin.H{"id": index + 1, "title": label, "image": value.Image, "url": value.URL})
				}
			}
		}
		if item.Type != "product" || item.Params.Source != "auto" || item.Params.Auto.Category == 0 {
			continue
		}
		if _, exists := seenDisplayTypes[item.Params.Auto.Category]; exists {
			continue
		}
		seenDisplayTypes[item.Params.Auto.Category] = struct{}{}
		limit := item.Params.Auto.ShowNum
		if limit < 4 {
			limit = 8
		}
		if limit > 24 {
			limit = 24
		}
		sort, _ := productSort(item.Params.Auto.ProductSort, "desc")
		result.DisplayTypes = append(result.DisplayTypes, homeDisplayType{
			CategoryID:   item.Params.Auto.Category,
			InitialLimit: limit,
			Sort:         sort,
		})
	}
	return result
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

// Stores supplies the public store directory used by the PC marketplace.  It
// reads only business projections, so a public user never traverses into the
// merchant database to render the directory.
func (h *Handler) Stores(c *gin.Context) {
	scope, err := h.resolveStoreScope(c)
	if err != nil {
		writeScopeError(c, err)
		return
	}

	query := h.db.WithContext(c.Request.Context()).
		Table("qixi_crm_b_store_view AS s").
		Select(`s.store_id,s.merchant_id,s.store_name,
			COUNT(p.product_id) AS product_count,
			COALESCE(SUM(p.sales), 0) AS sales_count,
			COALESCE(MAX(p.cover_url), '') AS cover_url`).
		Joins("LEFT JOIN qixi_crm_b_product_view AS p ON p.store_id = s.store_id AND p.sale_status = 1").
		Where("s.status = 1")
	if scope.MerchantID != 0 {
		query = query.Where("s.merchant_id = ?", scope.MerchantID)
	}

	rows := make([]storeDirectoryItem, 0)
	if err := query.Group("s.store_id,s.merchant_id,s.store_name").
		Order("sales_count DESC,s.store_id DESC").
		Find(&rows).Error; err != nil {
		fail(c, err)
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{
			"store_id":      row.StoreID,
			"mer_id":        row.MerchantID,
			"name":          row.StoreName,
			"product_count": row.ProductCount,
			"sales_count":   row.SalesCount,
			"cover_url":     row.CoverURL,
		})
	}
	response.OK(c, gin.H{"list": items, "total": len(items)})
}

func (h *Handler) Products(c *gin.Context) {
	page, limit := pageParams(c)
	merchantID, _ := strconv.ParseUint(c.Query("mer_id"), 10, 64)
	scope, err := h.resolveStoreScope(c)
	if err != nil {
		writeScopeError(c, err)
		return
	}
	if scope.MerchantID != 0 {
		if merchantID != 0 && merchantID != scope.MerchantID {
			response.Fail(c, http.StatusForbidden, "X-AppId 与商户范围不一致")
			return
		}
		merchantID = scope.MerchantID
	}
	categoryID, err := queryUint(c, "cate_id")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "分类 ID 错误")
		return
	}
	minPrice, maxPrice, err := priceRange(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "价格范围错误")
		return
	}
	sort, order := productSort(c.Query("sort"), c.Query("order"))
	rows, total, err := h.listPage(c, merchantID, categoryID, strings.TrimSpace(c.Query("keyword")), minPrice, maxPrice, page, limit, sort, order)
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
	scope, scopeErr := h.resolveStoreScope(c)
	if scopeErr != nil {
		writeScopeError(c, scopeErr)
		return
	}
	query := h.db.WithContext(c.Request.Context()).Where("product_id = ? AND sale_status = 1", id)
	if scope.MerchantID != 0 {
		query = query.Where("merchant_id = ?", scope.MerchantID)
	}
	var row productView
	err = query.First(&row).Error
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
	scope, scopeErr := h.resolveStoreScope(c)
	if scopeErr != nil {
		writeScopeError(c, scopeErr)
		return
	}
	if scope.MerchantID != 0 && merchantID != scope.MerchantID {
		response.Fail(c, http.StatusForbidden, "X-AppId 与商户范围不一致")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	rows, total, err := h.listPage(c, merchantID, 0, "", nil, nil, page, 20, "default", "desc")
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
	rows, _, err := h.listPage(c, merchantID, categoryID, keyword, nil, nil, page, limit, sort, order)
	return rows, err
}
func (h *Handler) listPage(c *gin.Context, merchantID, categoryID uint64, keyword string, minPrice, maxPrice *float64, page, limit int, sort, order string) ([]productView, int64, error) {
	query := h.db.WithContext(c.Request.Context()).Model(&productView{}).Where("sale_status = 1")
	if merchantID != 0 {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if categoryID != 0 {
		categoryIDs, err := h.categoryIDsIncludingDescendants(c, categoryID)
		if err != nil {
			return nil, 0, err
		}
		query = query.Where("category_id IN ?", categoryIDs)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if minPrice != nil {
		query = query.Where("price >= ?", *minPrice)
	}
	if maxPrice != nil {
		query = query.Where("price <= ?", *maxPrice)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	rows := make([]productView, 0)
	err := query.Order(productOrder(sort, order)).Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

// priceRange accepts optional non-negative decimal values and rejects a
// reversed range before the database query is composed.
func priceRange(c *gin.Context) (min, max *float64, err error) {
	parse := func(name string) (*float64, error) {
		raw := strings.TrimSpace(c.Query(name))
		if raw == "" {
			return nil, nil
		}
		value, parseErr := strconv.ParseFloat(raw, 64)
		if parseErr != nil || value < 0 {
			return nil, fmt.Errorf("invalid %s", name)
		}
		return &value, nil
	}
	if min, err = parse("min_price"); err != nil {
		return nil, nil, err
	}
	if max, err = parse("max_price"); err != nil {
		return nil, nil, err
	}
	if min != nil && max != nil && *min > *max {
		return nil, nil, fmt.Errorf("reversed range")
	}
	return min, max, nil
}

// storeScope is derived exclusively from a verified qixi_crm_b_store_view row.
// A client-provided X-AppId narrows public catalog access; it never expands it.
type storeScope struct {
	MerchantID uint64
	StoreID    uint64
}

func (h *Handler) resolveStoreScope(c *gin.Context) (storeScope, error) {
	appID := strings.TrimSpace(c.GetHeader("X-AppId"))
	if appID == "" {
		return storeScope{}, nil
	}
	var store storeView
	err := h.db.WithContext(c.Request.Context()).
		Where("store_app_id = ? AND status = 1", appID).
		First(&store).Error
	if err != nil {
		return storeScope{}, err
	}
	return storeScope{MerchantID: store.MerchantID, StoreID: store.StoreID}, nil
}

func writeScopeError(c *gin.Context, err error) {
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, http.StatusNotFound, "X-AppId 对应店铺不存在或已停用")
		return
	}
	fail(c, err)
}

// categoryIDsIncludingDescendants makes a root-category filter cover all of
// its visible descendants. The category projection is intentionally read in
// one query to avoid recursive N+1 lookups on the public product list.
func (h *Handler) categoryIDsIncludingDescendants(c *gin.Context, rootID uint64) ([]uint64, error) {
	rows := make([]categoryView, 0)
	if err := h.db.WithContext(c.Request.Context()).
		Where("status = ?", 1).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return descendantCategoryIDs(rows, rootID), nil
}

func descendantCategoryIDs(rows []categoryView, rootID uint64) []uint64 {
	children := make(map[uint64][]uint64, len(rows))
	known := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		known[row.CategoryID] = struct{}{}
		children[row.ParentID] = append(children[row.ParentID], row.CategoryID)
	}
	if _, ok := known[rootID]; !ok {
		return []uint64{rootID}
	}
	out := make([]uint64, 0, 8)
	seen := map[uint64]bool{rootID: true}
	queue := []uint64{rootID}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		out = append(out, current)
		for _, child := range children[current] {
			if !seen[child] {
				seen[child] = true
				queue = append(queue, child)
			}
		}
	}
	return out
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
	return gin.H{"id": row.ProductID, "title": row.Title, "mer_id": row.MerchantID, "mer_name": row.MerchantName, "store_id": row.StoreID, "store_name": row.Title, "shop_name": row.StoreName, "category_id": row.CategoryID, "image": row.CoverURL, "price": fmt.Sprintf("%.2f", row.Price), "ot_price": original, "sales": row.Sales, "stock": row.Stock}
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
