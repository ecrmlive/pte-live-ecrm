// Package nativecatalog implements the new-schema platform product audit flow.
package nativecatalog

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/adminscope"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	adminDB    *gorm.DB
	merchantDB *gorm.DB
	businessDB *gorm.DB
}

func NewHandler(adminDB, merchantDB, businessDB *gorm.DB) *Handler {
	return &Handler{adminDB: adminDB, merchantDB: merchantDB, businessDB: businessDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	h.RegisterMeta(r)
	r.GET("/products", h.list)
	r.GET("/products/:id", h.get)
	r.POST("/products/:id/audit", middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.adminDB, "product.audit.submit"), h.audit)
}

type productRow struct {
	ID            uint64    `gorm:"column:id"`
	StoreID       uint64    `gorm:"column:store_id"`
	MerchantID    uint64    `gorm:"column:merchant_id"`
	MerchantName  string    `gorm:"column:merchant_name"`
	StoreName     string    `gorm:"column:store_name"`
	Title         string    `gorm:"column:title"`
	CategoryID    uint64    `gorm:"column:category_id"`
	BrandName     string    `gorm:"column:brand_name"`
	SVIPPriceType int8      `gorm:"column:svip_price_type"`
	SVIPPrice     float64   `gorm:"column:svip_price"`
	Status        string    `gorm:"column:status"`
	Version       uint64    `gorm:"column:version"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

type detailRow struct {
	ProductID     uint64   `gorm:"column:product_id"`
	Brief         string   `gorm:"column:brief"`
	CoverURL      string   `gorm:"column:cover_url"`
	OriginalPrice *float64 `gorm:"column:original_price"`
}

type skuRow struct {
	ID        uint64  `gorm:"column:id"`
	ProductID uint64  `gorm:"column:product_id"`
	SpecJSON  []byte  `gorm:"column:spec_json"`
	Price     float64 `gorm:"column:price"`
	Stock     int     `gorm:"column:stock"`
}

type productReview struct {
	ID            uint64  `gorm:"column:id;primaryKey"`
	ProductID     uint64  `gorm:"column:product_id"`
	StoreID       uint64  `gorm:"column:store_id"`
	SourceEventID *uint64 `gorm:"column:source_event_id"`
	Status        string  `gorm:"column:status"`
	Reason        string  `gorm:"column:reason"`
	ReviewedBy    uint64  `gorm:"column:reviewed_by"`
}

func (productReview) TableName() string { return "qixi_crm_a_product_review" }

// productProjectionOutbox keeps the cross-database projection command beside
// the platform review fact. A failed business projection is therefore durable
// and can be retried instead of being lost after the merchant status changes.
type productProjectionOutbox struct {
	ID            uint64    `gorm:"column:id;primaryKey"`
	ProductID     uint64    `gorm:"column:product_id"`
	SourceEventID *uint64   `gorm:"column:source_event_id"`
	Action        string    `gorm:"column:action"`
	Payload       []byte    `gorm:"column:payload"`
	Status        string    `gorm:"column:status"`
	Attempts      uint      `gorm:"column:attempts"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (productProjectionOutbox) TableName() string { return "qixi_crm_a_product_projection_outbox" }

type categoryRow struct {
	ID   uint64 `gorm:"column:category_id"`
	Name string `gorm:"column:name"`
}

type productResponse struct {
	ProductID  uint64  `json:"product_id"`
	MerID      uint64  `json:"mer_id"`
	MerName    string  `json:"mer_name"`
	StoreName  string  `json:"store_name"`
	Title      string  `json:"title"`
	StoreInfo  string  `json:"store_info"`
	CateName   string  `json:"cate_name"`
	BrandName  string  `json:"brand_name"`
	Image      string  `json:"image"`
	Price      float64 `json:"price"`
	OtPrice    float64 `json:"ot_price"`
	Stock      int     `json:"stock"`
	Sales      int     `json:"sales"`
	Status     int8    `json:"status"`
	IsShow     int8    `json:"is_show"`
	Refusal    string  `json:"refusal,omitempty"`
	CreateTime string  `json:"create_time"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := page(c)
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置商品监管数据范围")
		return
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.OK(c, gin.H{"list": []productResponse{}, "total": 0, "page": page, "limit": limit})
		return
	}
	query := h.base(c, merchantIDs)
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		query = query.Where("p.title LIKE ?", "%"+keyword+"%")
	}
	if merchantID, _ := strconv.ParseUint(c.Query("mer_id"), 10, 64); merchantID > 0 {
		query = query.Where("s.merchant_id = ?", merchantID)
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("p.status = ?", statusName(status))
	}
	query = queryfilter.ApplyCreatedAtRange(query, c, "p.created_at")
	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品失败")
		return
	}
	var rows []productRow
	if err := query.Order("p.created_at DESC,p.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
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
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "商品 ID 错误")
		return
	}
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置商品监管数据范围")
		return
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	var row productRow
	if err := h.base(c, merchantIDs).Where("p.id = ?", id).Scan(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品失败")
		return
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	items, err := h.responses(c, []productRow{row})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载商品详情失败")
		return
	}
	response.OK(c, items[0])
}

func (h *Handler) audit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status  int8   `json:"status"`
		Refusal string `json:"refusal"`
	}
	if id == 0 || c.ShouldBindJSON(&req) != nil || (req.Status != 1 && req.Status != -1) || (req.Status == -1 && strings.TrimSpace(req.Refusal) == "") {
		response.Fail(c, http.StatusBadRequest, "审核参数不合法")
		return
	}
	reviewStatus, action, next := "rejected", "delete", "rejected"
	if req.Status == 1 {
		reviewStatus, action, next = "approved", "upsert", "on_sale"
	}
	var row productRow
	var command productAuditOutbox
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("qixi_crm_m_product AS p").Select("p.id,p.store_id,s.merchant_id,m.name AS merchant_name,s.name AS store_name,p.title,p.category_id,p.brand_name,p.svip_price_type,p.svip_price,p.status,p.version,p.created_at").Joins("JOIN qixi_crm_m_store AS s ON s.id = p.store_id").Joins("JOIN qixi_crm_m_merchant AS m ON m.id = s.merchant_id").Where("p.id = ?", id).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID == 0 {
			return gorm.ErrRecordNotFound
		}
		if row.Status != "pending_review" && row.Status != "draft" {
			return errAuditState
		}
		if req.Status == 1 {
			var sellableSKUs int64
			if err := tx.Table("qixi_crm_m_product_sku").Where("product_id = ? AND status = 1", id).Count(&sellableSKUs).Error; err != nil {
				return err
			}
			if sellableSKUs == 0 {
				return errMissingSellableSKU
			}
		}
		if err := tx.Table("qixi_crm_m_product").Where("id = ? AND status = ?", id, row.Status).Updates(map[string]any{"status": next, "version": gorm.Expr("version + 1")}).Error; err != nil {
			return err
		}
		command = productAuditOutbox{ProductID: row.ID, StoreID: row.StoreID, Action: action, ReviewStatus: reviewStatus, Reason: strings.TrimSpace(req.Refusal), ReviewedBy: uint64(middleware.AdminID(c)), Status: "pending"}
		return tx.Create(&command).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	if errors.Is(err, errAuditState) {
		response.Fail(c, http.StatusConflict, "商品当前不在待审核状态")
		return
	}
	if errors.Is(err, errMissingSellableSKU) {
		response.Fail(c, http.StatusUnprocessableEntity, "商品缺少可售规格，不能审核上架")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新商品审核状态失败")
		return
	}
	processed, err := h.processAuditOutbox(c.Request.Context(), command)
	if err != nil {
		response.OK(c, gin.H{"ok": true, "projection_pending": true})
		return
	}
	if !processed {
		response.OK(c, gin.H{"ok": true, "projection_pending": true})
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) base(c *gin.Context, merchantIDs []uint64) *gorm.DB {
	q := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product AS p").Select("p.id,p.store_id,s.merchant_id,m.name AS merchant_name,s.name AS store_name,p.title,p.category_id,p.brand_name,p.svip_price_type,p.svip_price,p.status,p.version,p.created_at").Joins("JOIN qixi_crm_m_store AS s ON s.id = p.store_id").Joins("JOIN qixi_crm_m_merchant AS m ON m.id = s.merchant_id")
	if merchantIDs != nil {
		q = q.Where("s.merchant_id IN ?", merchantIDs)
	}
	return q
}

// merchantScope maps unified-admin direct merchant and region assignments to
// the current merchant store. nil is the platform's full supervision scope.
func (h *Handler) merchantScope(c *gin.Context) ([]uint64, error) {
	scope, err := adminscope.ResolveMerchantScope(c.Request.Context(), h.adminDB, middleware.ClaimsFrom(c))
	if err != nil {
		return nil, err
	}
	if scope.Full {
		return nil, nil
	}
	ids := append([]uint64{}, scope.MerchantIDs...)
	if len(scope.RegionIDs) == 0 {
		return ids, nil
	}
	var rows []struct{ ID uint64 }
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_merchant").Select("id").Where("region_id IN ?", scope.RegionIDs).Find(&rows).Error; err != nil {
		return nil, err
	}
	regionMerchantIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		regionMerchantIDs = append(regionMerchantIDs, row.ID)
	}
	return mergeMerchantIDs(ids, regionMerchantIDs), nil
}

func mergeMerchantIDs(direct, regional []uint64) []uint64 {
	result := make([]uint64, 0, len(direct)+len(regional))
	seen := make(map[uint64]struct{}, len(direct)+len(regional))
	for _, list := range [][]uint64{direct, regional} {
		for _, id := range list {
			if id == 0 {
				continue
			}
			if _, exists := seen[id]; exists {
				continue
			}
			seen[id] = struct{}{}
			result = append(result, id)
		}
	}
	return result
}

func (h *Handler) responses(c *gin.Context, rows []productRow) ([]productResponse, error) {
	if len(rows) == 0 {
		return []productResponse{}, nil
	}
	ids, categories := make([]uint64, 0, len(rows)), make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
		categories = append(categories, row.CategoryID)
	}
	var details []detailRow
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_detail").Where("product_id IN ?", ids).Find(&details).Error; err != nil {
		return nil, err
	}
	var skus []skuRow
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_sku").Where("product_id IN ? AND status = 1", ids).Order("id ASC").Find(&skus).Error; err != nil {
		return nil, err
	}
	var categoriesRows []categoryRow
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_category_view").Where("category_id IN ?", categories).Find(&categoriesRows).Error; err != nil {
		return nil, err
	}
	var reviews []productReview
	if err := h.adminDB.WithContext(c.Request.Context()).Where("product_id IN ?", ids).Order("id DESC").Find(&reviews).Error; err != nil {
		return nil, err
	}
	detailByID := map[uint64]detailRow{}
	for _, x := range details {
		detailByID[x.ProductID] = x
	}
	skuByID := map[uint64]skuRow{}
	for _, x := range skus {
		skuByID[x.ProductID] = x
	}
	nameByID := map[uint64]string{}
	for _, x := range categoriesRows {
		nameByID[x.ID] = x.Name
	}
	reviewByID := map[uint64]productReview{}
	for _, x := range reviews {
		if _, ok := reviewByID[x.ProductID]; !ok {
			reviewByID[x.ProductID] = x
		}
	}
	out := make([]productResponse, 0, len(rows))
	for _, row := range rows {
		d, s := detailByID[row.ID], skuByID[row.ID]
		ot := 0.0
		if d.OriginalPrice != nil {
			ot = *d.OriginalPrice
		}
		out = append(out, productResponse{ProductID: row.ID, MerID: row.MerchantID, MerName: row.MerchantName, StoreName: row.StoreName, Title: row.Title, StoreInfo: d.Brief, CateName: nameByID[row.CategoryID], BrandName: row.BrandName, Image: d.CoverURL, Price: s.Price, OtPrice: ot, Stock: s.Stock, Sales: 0, Status: statusCode(row.Status), IsShow: showCode(row.Status), Refusal: reviewByID[row.ID].Reason, CreateTime: row.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	return out, nil
}

// syncSKUProjection copies only sellable SKU facts into the business-owned
// consumption view. The stored merchant SKU primary key is the later
// inventory-command identity; it is never inferred from product_id or text.
func (h *Handler) syncSKUProjection(c *gin.Context, productID uint64) error {
	var skus []skuRow
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_sku").Where("product_id = ? AND status = 1", productID).Order("id ASC").Find(&skus).Error; err != nil {
		return err
	}
	if len(skus) == 0 {
		return errors.New("product has no sellable sku")
	}
	return h.businessDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM qixi_crm_b_product_sku_view WHERE product_id = ?", productID).Error; err != nil {
			return err
		}
		for _, sku := range skus {
			spec := string(sku.SpecJSON)
			if strings.TrimSpace(spec) == "" {
				spec = "{}"
			}
			if err := tx.Exec(`INSERT INTO qixi_crm_b_product_sku_view
                (merchant_sku_id,product_id,sku_key,spec_snapshot,price,stock,sale_status,version,updated_at)
                VALUES (?,?,?,?,?,?,1,1,NOW())`, sku.ID, productID, strconv.FormatUint(sku.ID, 10), spec, sku.Price, sku.Stock).Error; err != nil {
				return err
			}
		}
		return nil
	})
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
func statusName(v string) string {
	switch v {
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
func statusCode(v string) int8 {
	switch v {
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
func showCode(v string) int8 {
	if v == "on_sale" {
		return 1
	}
	return 0
}
func nullableOriginal(value float64) any {
	if value <= 0 {
		return nil
	}
	return value
}

var (
	errAuditState         = errors.New("product audit state invalid")
	errMissingSellableSKU = errors.New("product has no sellable sku")
)
