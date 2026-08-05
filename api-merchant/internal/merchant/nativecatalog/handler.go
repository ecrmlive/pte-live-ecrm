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

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
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
}

func NewHandler(merchantDB, businessDB *gorm.DB) *Handler {
	return &Handler{merchantDB: merchantDB, businessDB: businessDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/product-categories", h.categories)
	r.GET("/products", h.list)
	r.GET("/products/recycle-bin", h.listRecycleBin)
	r.GET("/products/:id", h.get)
	r.POST("/products", middleware.RequireStorePermission(h.merchantDB, "product.create"), h.create)
	r.PUT("/products/:id", middleware.RequireStorePermission(h.merchantDB, "product.update"), h.update)
	r.DELETE("/products/:id", middleware.RequireStorePermission(h.merchantDB, "product.delete"), h.remove)
	r.POST("/products/:id/restore", middleware.RequireStorePermission(h.merchantDB, "product.restore"), h.restore)
	r.PUT("/products/:id/show", middleware.RequireStorePermission(h.merchantDB, "product.show"), h.setShow)
	r.PUT("/products/:id/stock", middleware.RequireStorePermission(h.merchantDB, "product.stock"), h.setStock)
}

type product struct {
	ID            uint64    `gorm:"column:id;primaryKey"`
	StoreID       uint64    `gorm:"column:store_id"`
	Title         string    `gorm:"column:title"`
	CategoryID    uint64    `gorm:"column:category_id"`
	BrandName     string    `gorm:"column:brand_name"`
	SVIPPriceType int8      `gorm:"column:svip_price_type"`
	SVIPPrice     float64   `gorm:"column:svip_price"`
	Status        string    `gorm:"column:status"`
	Version       uint64    `gorm:"column:version"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
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

type recycleBin struct {
	ProductID          uint64    `gorm:"column:product_id;primaryKey"`
	StoreID            uint64    `gorm:"column:store_id"`
	DeletedByAccountID uint64    `gorm:"column:deleted_by_account_id"`
	DeletedAt          time.Time `gorm:"column:deleted_at"`
	RestoreUntil       time.Time `gorm:"column:restore_until"`
}

func (recycleBin) TableName() string { return "qixi_crm_m_product_recycle_bin" }

type stockLedger struct {
	SKUId           uint64 `gorm:"column:sku_id"`
	ChangeQuantity  int    `gorm:"column:change_quantity"`
	BalanceQuantity int    `gorm:"column:balance_quantity"`
	ReasonType      string `gorm:"column:reason_type"`
	ReferenceType   string `gorm:"column:reference_type"`
	ReferenceID     string `gorm:"column:reference_id"`
	IdempotencyKey  string `gorm:"column:idempotency_key"`
}

func (stockLedger) TableName() string { return "qixi_crm_m_stock_ledger" }

type saveRequest struct {
	CateID        uint64  `json:"cate_id" binding:"required"`
	BrandName     string  `json:"brand_name"`
	StoreName     string  `json:"store_name" binding:"required"`
	StoreInfo     string  `json:"store_info"`
	Keyword       string  `json:"keyword"`
	UnitName      string  `json:"unit_name"`
	Price         float64 `json:"price" binding:"gte=0"`
	OtPrice       float64 `json:"ot_price" binding:"gte=0"`
	Stock         int     `json:"stock" binding:"gte=0"`
	Image         string  `json:"image"`
	SliderImage   string  `json:"slider_image"`
	SVIPPriceType int8    `json:"svip_price_type"`
	SVIPPrice     float64 `json:"svip_price"`
}

type productResponse struct {
	ProductID     uint64  `json:"product_id"`
	CateID        uint64  `json:"cate_id"`
	CateName      string  `json:"cate_name"`
	BrandName     string  `json:"brand_name"`
	StoreName     string  `json:"store_name"`
	StoreInfo     string  `json:"store_info"`
	Keyword       string  `json:"keyword"`
	UnitName      string  `json:"unit_name"`
	Image         string  `json:"image"`
	Price         float64 `json:"price"`
	OtPrice       float64 `json:"ot_price"`
	SVIPPriceType int8    `json:"svip_price_type"`
	SVIPPrice     float64 `json:"svip_price"`
	Stock         int     `json:"stock"`
	Status        int8    `json:"status"`
	IsShow        int8    `json:"is_show"`
	CreateTime    string  `json:"create_time"`
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
	query := h.merchantDB.WithContext(c.Request.Context()).Model(&product{}).
		Where("store_id = ?", middleware.StoreID(c)).
		Where("NOT EXISTS (SELECT 1 FROM qixi_crm_m_product_recycle_bin AS rb WHERE rb.product_id = qixi_crm_m_product.id)")
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

func (h *Handler) listRecycleBin(c *gin.Context) {
	page, limit := page(c)
	base := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_recycle_bin AS rb").
		Joins("INNER JOIN qixi_crm_m_product AS p ON p.id = rb.product_id").
		Where("rb.store_id = ?", middleware.StoreID(c))
	var total int64
	if err := base.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询回收站失败")
		return
	}
	var rows []struct {
		ProductID uint64    `gorm:"column:product_id"`
		Title     string    `gorm:"column:title"`
		DeletedAt time.Time `gorm:"column:deleted_at"`
		RestoreAt time.Time `gorm:"column:restore_until"`
	}
	if err := base.Select("rb.product_id,p.title,rb.deleted_at,rb.restore_until").Order("rb.deleted_at DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询回收站失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{"product_id": row.ProductID, "store_name": row.Title, "deleted_at": row.DeletedAt.Format("2006-01-02 15:04:05"), "restore_until": row.RestoreAt.Format("2006-01-02 15:04:05")})
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
		created = product{StoreID: uint64(middleware.StoreID(c)), Title: strings.TrimSpace(req.StoreName), CategoryID: req.CateID, BrandName: strings.TrimSpace(req.BrandName), SVIPPriceType: req.SVIPPriceType, SVIPPrice: req.SVIPPrice, Status: "pending_review", Version: 1}
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
		if err := tx.Model(&product{}).Where("id = ? AND store_id = ?", id, middleware.StoreID(c)).Updates(map[string]any{"title": strings.TrimSpace(req.StoreName), "category_id": req.CateID, "brand_name": strings.TrimSpace(req.BrandName), "svip_price_type": req.SVIPPriceType, "svip_price": req.SVIPPrice, "status": status, "version": gorm.Expr("version + 1")}).Error; err != nil {
			return err
		}
		updated.Title, updated.CategoryID, updated.BrandName, updated.SVIPPriceType, updated.SVIPPrice, updated.Status, updated.Version = strings.TrimSpace(req.StoreName), req.CateID, strings.TrimSpace(req.BrandName), req.SVIPPriceType, req.SVIPPrice, status, updated.Version+1
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
		now := time.Now()
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "product_id"}},
			DoUpdates: clause.Assignments(map[string]any{"store_id": middleware.StoreID(c), "deleted_by_account_id": middleware.AdminID(c), "deleted_at": now, "restore_until": now.Add(30 * 24 * time.Hour)}),
		}).Create(&recycleBin{ProductID: id, StoreID: uint64(middleware.StoreID(c)), DeletedByAccountID: uint64(middleware.AdminID(c)), DeletedAt: now, RestoreUntil: now.Add(30 * 24 * time.Hour)}).Error; err != nil {
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

func (h *Handler) restore(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "商品 ID 错误")
		return
	}
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var row recycleBin
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_id = ? AND store_id = ?", id, middleware.StoreID(c)).First(&row).Error; err != nil {
			return err
		}
		if time.Now().After(row.RestoreUntil) {
			return errRestoreExpired
		}
		if err := tx.Model(&product{}).Where("id = ? AND store_id = ?", id, middleware.StoreID(c)).Updates(map[string]any{"status": "pending_review", "version": gorm.Expr("version + 1")}).Error; err != nil {
			return err
		}
		return tx.Delete(&recycleBin{}, "product_id = ?", id).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "回收站商品不存在")
		return
	}
	if errors.Is(err, errRestoreExpired) {
		response.Fail(c, http.StatusConflict, "商品已超过可恢复期限")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "恢复商品失败")
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
	idempotencyKey := strings.TrimSpace(c.GetHeader("X-Idempotency-Key"))
	if id == 0 || c.ShouldBindJSON(&req) != nil || req.Stock == nil || *req.Stock < 0 || idempotencyKey == "" || len(idempotencyKey) > 128 {
		response.Fail(c, http.StatusBadRequest, "库存错误")
		return
	}
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var p product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND store_id = ?", id, middleware.StoreID(c)).First(&p).Error; err != nil {
			return err
		}
		var current sku
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_id = ? AND status = 1", id).First(&current).Error; err != nil {
			return err
		}
		var prior stockLedger
		if err := tx.Where("sku_id = ? AND idempotency_key = ?", current.ID, idempotencyKey).Take(&prior).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := tx.Model(&sku{}).Where("id = ?", current.ID).Update("stock", *req.Stock).Error; err != nil {
			return err
		}
		return tx.Create(&stockLedger{SKUId: current.ID, ChangeQuantity: *req.Stock - current.Stock, BalanceQuantity: *req.Stock, ReasonType: "merchant_manual_adjust", ReferenceType: "product", ReferenceID: strconv.FormatUint(id, 10), IdempotencyKey: idempotencyKey}).Error
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
		out = append(out, productResponse{ProductID: p.ID, CateID: p.CategoryID, CateName: nameByID[p.CategoryID], BrandName: p.BrandName, StoreName: p.Title, StoreInfo: d.Brief, Keyword: d.Keyword, UnitName: nonEmpty(d.UnitName, "件"), Image: d.CoverURL, Price: s.Price, OtPrice: ot, SVIPPriceType: p.SVIPPriceType, SVIPPrice: p.SVIPPrice, Stock: s.Stock, Status: statusCode(p.Status), IsShow: showCode(p.Status), CreateTime: p.CreatedAt.Format("2006-01-02 15:04:05")})
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
	if req.CateID == 0 || strings.TrimSpace(req.StoreName) == "" || req.Price < 0 || req.Stock < 0 {
		return false
	}
	if len([]rune(strings.TrimSpace(req.BrandName))) > 64 || strings.ContainsAny(req.BrandName, "\x00\r\n") {
		return false
	}
	if req.SVIPPriceType < 0 || req.SVIPPriceType > 2 || req.SVIPPrice < 0 {
		return false
	}
	return req.SVIPPriceType != 2 || (req.SVIPPrice > 0 && req.SVIPPrice < req.Price)
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
var errRestoreExpired = errors.New("product restore expired")
