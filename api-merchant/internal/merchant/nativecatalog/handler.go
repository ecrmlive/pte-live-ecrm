// Package nativecatalog owns the store product console after the qixi_crm_m_
// migration.  It deliberately does not depend on the legacy catalog domain.
package nativecatalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	eventProductUpsert = "merchant.catalog.product.upsert"
	eventProductRemove = "merchant.catalog.product.remove"
)

type Handler struct {
	merchantDB *gorm.DB
	businessDB *gorm.DB
	identity   *identity.Service
}

func NewHandler(merchantDB, businessDB *gorm.DB, id *identity.Service) *Handler {
	return &Handler{merchantDB: merchantDB, businessDB: businessDB, identity: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/product-categories", h.categories)
	r.GET("/products", h.list)
	r.GET("/products/:id", h.get)
	r.POST("/products", middleware.RequireMerchantMenu(h.identity, identity.MerPermProductCreate), h.create)
	r.PUT("/products/:id", h.update)
	r.DELETE("/products/:id", middleware.RequireMerchantMenu(h.identity, identity.MerPermProductDelete), h.remove)
	r.PUT("/products/:id/show", middleware.RequireMerchantMenu(h.identity, identity.MerPermProductShow), h.setShow)
	r.PUT("/products/:id/stock", middleware.RequireMerchantMenu(h.identity, identity.MerPermProductStock), h.setStock)
}

type product struct {
	ID         uint64    `gorm:"column:id;primaryKey"`
	StoreID    uint64    `gorm:"column:store_id"`
	Title      string    `gorm:"column:title"`
	CategoryID uint64    `gorm:"column:category_id"`
	Status     string    `gorm:"column:status"`
	Version    uint64    `gorm:"column:version"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (product) TableName() string { return "qixi_crm_m_product" }

type productDetail struct {
	ProductID     uint64   `gorm:"column:product_id;primaryKey"`
	Brief         string   `gorm:"column:brief"`
	Keyword       string   `gorm:"column:keyword"`
	UnitName      string   `gorm:"column:unit_name"`
	CoverURL      string   `gorm:"column:cover_url"`
	OriginalPrice *float64 `gorm:"column:original_price"`
}

func (productDetail) TableName() string { return "qixi_crm_m_product_detail" }

type sku struct {
	ID        uint64  `gorm:"column:id;primaryKey"`
	ProductID uint64  `gorm:"column:product_id"`
	SpecJSON  string  `gorm:"column:spec_json"`
	Price     float64 `gorm:"column:price"`
	Stock     int     `gorm:"column:stock"`
	Status    int8    `gorm:"column:status"`
}

func (sku) TableName() string { return "qixi_crm_m_product_sku" }

type store struct {
	ID         uint64 `gorm:"column:id;primaryKey"`
	MerchantID uint64 `gorm:"column:merchant_id"`
	Name       string `gorm:"column:name"`
}

func (store) TableName() string { return "qixi_crm_m_store" }

type category struct {
	ID       uint64 `gorm:"column:category_id"`
	ParentID uint64 `gorm:"column:parent_id"`
	Name     string `gorm:"column:name"`
	Sort     int    `gorm:"column:sort"`
}

type outboxEvent struct {
	EventType     string `gorm:"column:event_type"`
	AggregateType string `gorm:"column:aggregate_type"`
	AggregateID   string `gorm:"column:aggregate_id"`
	Payload       string `gorm:"column:payload"`
}

func (outboxEvent) TableName() string { return "qixi_crm_m_outbox" }

type saveRequest struct {
	CateID      uint64  `json:"cate_id" binding:"required"`
	StoreName   string  `json:"store_name" binding:"required"`
	StoreInfo   string  `json:"store_info"`
	Keyword     string  `json:"keyword"`
	UnitName    string  `json:"unit_name"`
	Price       float64 `json:"price" binding:"gte=0"`
	OtPrice     float64 `json:"ot_price" binding:"gte=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	Image       string  `json:"image"`
	SliderImage string  `json:"slider_image"`
}

type productResponse struct {
	ProductID  uint64  `json:"product_id"`
	CateID     uint64  `json:"cate_id"`
	CateName   string  `json:"cate_name"`
	StoreName  string  `json:"store_name"`
	StoreInfo  string  `json:"store_info"`
	Keyword    string  `json:"keyword"`
	UnitName   string  `json:"unit_name"`
	Image      string  `json:"image"`
	Price      float64 `json:"price"`
	OtPrice    float64 `json:"ot_price"`
	Stock      int     `json:"stock"`
	Status     int8    `json:"status"`
	IsShow     int8    `json:"is_show"`
	CreateTime string  `json:"create_time"`
}

func (h *Handler) categories(c *gin.Context) {
	var rows []category
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_category_view").Where("status = 1").Order("sort, category_id").Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品分类失败")
		return
	}
	byParent := map[uint64][]gin.H{}
	for _, row := range rows {
		byParent[row.ParentID] = append(byParent[row.ParentID], gin.H{"store_category_id": row.ID, "cate_name": row.Name, "is_show": 1})
	}
	var makeTree func(uint64) []gin.H
	makeTree = func(parentID uint64) []gin.H {
		items := byParent[parentID]
		for i := range items {
			id := items[i]["store_category_id"].(uint64)
			items[i]["children"] = makeTree(id)
		}
		return items
	}
	response.OK(c, gin.H{"list": makeTree(0)})
}

func (h *Handler) list(c *gin.Context) {
	page, limit := page(c)
	query := h.merchantDB.WithContext(c.Request.Context()).Model(&product{}).Where("store_id = ?", middleware.StoreID(c))
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("status = ?", statusName(status))
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品失败")
		return
	}
	var rows []product
	if err := query.Order("updated_at DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品失败")
		return
	}
	items, err := h.responses(c, rows)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载商品详情失败")
		return
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func (h *Handler) get(c *gin.Context) {
	p, err := h.owned(c, parseID(c))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品失败")
		return
	}
	items, err := h.responses(c, []product{p})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载商品详情失败")
		return
	}
	response.OK(c, items[0])
}

func (h *Handler) create(c *gin.Context) {
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil || !valid(req) {
		response.Fail(c, http.StatusBadRequest, "商品参数不合法")
		return
	}
	var created product
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		created = product{StoreID: uint64(middleware.StoreID(c)), Title: strings.TrimSpace(req.StoreName), CategoryID: req.CateID, Status: "pending_review", Version: 1}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		return h.saveFields(tx, created, req)
	})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建商品失败")
		return
	}
	items, err := h.responses(c, []product{created})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载商品详情失败")
		return
	}
	response.OK(c, items[0])
}

func (h *Handler) update(c *gin.Context) {
	id := parseID(c)
	var req saveRequest
	if id == 0 || c.ShouldBindJSON(&req) != nil || !valid(req) {
		response.Fail(c, http.StatusBadRequest, "商品参数不合法")
		return
	}
	var updated product
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND store_id = ?", id, middleware.StoreID(c)).First(&updated).Error; err != nil {
			return err
		}
		status := updated.Status
		if status == "on_sale" || status == "off_sale" {
			status = "pending_review"
		}
		if err := tx.Model(&product{}).Where("id = ? AND store_id = ?", id, middleware.StoreID(c)).Updates(map[string]any{"title": strings.TrimSpace(req.StoreName), "category_id": req.CateID, "status": status, "version": gorm.Expr("version + 1")}).Error; err != nil {
			return err
		}
		updated.Title, updated.CategoryID, updated.Status, updated.Version = strings.TrimSpace(req.StoreName), req.CateID, status, updated.Version+1
		return h.saveFields(tx, updated, req)
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新商品失败")
		return
	}
	items, err := h.responses(c, []product{updated})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载商品详情失败")
		return
	}
	response.OK(c, items[0])
}

func (h *Handler) remove(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "商品 ID 错误")
		return
	}
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var p product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND store_id = ?", id, middleware.StoreID(c)).First(&p).Error; err != nil {
			return err
		}
		if err := tx.Delete(&product{}, id).Error; err != nil {
			return err
		}
		return enqueue(tx, eventProductRemove, id, gin.H{"product_id": id})
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除商品失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) setShow(c *gin.Context) {
	id := parseID(c)
	var req struct {
		IsShow *int8 `json:"is_show"`
	}
	if id == 0 || c.ShouldBindJSON(&req) != nil || req.IsShow == nil || (*req.IsShow != 0 && *req.IsShow != 1) {
		response.Fail(c, http.StatusBadRequest, "上架状态错误")
		return
	}
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var p product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND store_id = ?", id, middleware.StoreID(c)).First(&p).Error; err != nil {
			return err
		}
		if *req.IsShow == 1 && p.Status != "off_sale" {
			return errInvalidTransition
		}
		next := "off_sale"
		if *req.IsShow == 1 {
			next = "pending_review"
		}
		if err := tx.Model(&product{}).Where("id = ?", id).Updates(map[string]any{"status": next, "version": gorm.Expr("version + 1")}).Error; err != nil {
			return err
		}
		if next == "off_sale" {
			return enqueue(tx, eventProductRemove, id, gin.H{"product_id": id})
		}
		return nil
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	if errors.Is(err, errInvalidTransition) {
		response.Fail(c, http.StatusConflict, "当前商品不可直接上架，需要平台审核")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新上架状态失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) setStock(c *gin.Context) {
	id := parseID(c)
	var req struct {
		Stock *int `json:"stock"`
	}
	if id == 0 || c.ShouldBindJSON(&req) != nil || req.Stock == nil || *req.Stock < 0 {
		response.Fail(c, http.StatusBadRequest, "库存错误")
		return
	}
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var p product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND store_id = ?", id, middleware.StoreID(c)).First(&p).Error; err != nil {
			return err
		}
		return tx.Model(&sku{}).Where("product_id = ?", id).Update("stock", *req.Stock).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新库存失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) saveFields(tx *gorm.DB, p product, req saveRequest) error {
	cover := strings.TrimSpace(req.Image)
	if cover == "" {
		cover = strings.TrimSpace(req.SliderImage)
	}
	detail := productDetail{ProductID: p.ID, Brief: strings.TrimSpace(req.StoreInfo), Keyword: strings.TrimSpace(req.Keyword), UnitName: nonEmpty(req.UnitName, "件"), CoverURL: cover}
	if req.OtPrice > 0 {
		value := req.OtPrice
		detail.OriginalPrice = &value
	}
	if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "product_id"}}, DoUpdates: clause.AssignmentColumns([]string{"brief", "keyword", "unit_name", "cover_url", "original_price"})}).Create(&detail).Error; err != nil {
		return err
	}
	row := sku{ProductID: p.ID, SpecJSON: `{}`, Price: req.Price, Stock: req.Stock, Status: 1}
	if err := tx.Where("product_id = ?", p.ID).FirstOrCreate(&row).Error; err != nil {
		return err
	}
	return tx.Model(&sku{}).Where("product_id = ?", p.ID).Updates(map[string]any{"price": req.Price, "stock": req.Stock, "status": 1}).Error
}

func (h *Handler) owned(c *gin.Context, id uint64) (product, error) {
	var p product
	err := h.merchantDB.WithContext(c.Request.Context()).Where("id = ? AND store_id = ?", id, middleware.StoreID(c)).First(&p).Error
	return p, err
}

func (h *Handler) responses(c *gin.Context, rows []product) ([]productResponse, error) {
	if len(rows) == 0 {
		return []productResponse{}, nil
	}
	ids := make([]uint64, 0, len(rows))
	cats := make([]uint64, 0, len(rows))
	for _, p := range rows {
		ids = append(ids, p.ID)
		cats = append(cats, p.CategoryID)
	}
	var details []productDetail
	if err := h.merchantDB.WithContext(c.Request.Context()).Where("product_id IN ?", ids).Find(&details).Error; err != nil {
		return nil, err
	}
	var skus []sku
	if err := h.merchantDB.WithContext(c.Request.Context()).Where("product_id IN ? AND status = 1", ids).Find(&skus).Error; err != nil {
		return nil, err
	}
	var categories []category
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_category_view").Where("category_id IN ?", cats).Find(&categories).Error; err != nil {
		return nil, err
	}
	detailByID := map[uint64]productDetail{}
	for _, x := range details {
		detailByID[x.ProductID] = x
	}
	skuByID := map[uint64]sku{}
	for _, x := range skus {
		skuByID[x.ProductID] = x
	}
	nameByID := map[uint64]string{}
	for _, x := range categories {
		nameByID[x.ID] = x.Name
	}
	out := make([]productResponse, 0, len(rows))
	for _, p := range rows {
		d, s := detailByID[p.ID], skuByID[p.ID]
		ot := 0.0
		if d.OriginalPrice != nil {
			ot = *d.OriginalPrice
		}
		out = append(out, productResponse{ProductID: p.ID, CateID: p.CategoryID, CateName: nameByID[p.CategoryID], StoreName: p.Title, StoreInfo: d.Brief, Keyword: d.Keyword, UnitName: nonEmpty(d.UnitName, "件"), Image: d.CoverURL, Price: s.Price, OtPrice: ot, Stock: s.Stock, Status: statusCode(p.Status), IsShow: showCode(p.Status), CreateTime: p.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	return out, nil
}

func enqueue(tx *gorm.DB, kind string, id uint64, body any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return tx.Create(&outboxEvent{EventType: kind, AggregateType: "product", AggregateID: strconv.FormatUint(id, 10), Payload: string(raw)}).Error
}
func page(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if p < 1 {
		p = 1
	}
	if l < 1 {
		l = 20
	}
	if l > 100 {
		l = 100
	}
	return p, l
}
func parseID(c *gin.Context) uint64 { id, _ := strconv.ParseUint(c.Param("id"), 10, 64); return id }
func nonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
func valid(req saveRequest) bool {
	return req.CateID > 0 && strings.TrimSpace(req.StoreName) != "" && req.Price >= 0 && req.Stock >= 0
}
func statusCode(status string) int8 {
	switch status {
	case "on_sale":
		return 1
	case "rejected":
		return -1
	case "off_sale":
		return -2
	default:
		return 0
	}
}
func showCode(status string) int8 {
	if status == "on_sale" {
		return 1
	}
	return 0
}
func statusName(value string) string {
	switch value {
	case "1":
		return "on_sale"
	case "-1":
		return "rejected"
	case "-2":
		return "off_sale"
	default:
		return "pending_review"
	}
}

var errInvalidTransition = errors.New("invalid product state transition")
