// Package nativecatalog implements the new-schema platform product audit flow.
package nativecatalog

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"gorm.io/gorm"
)

type Handler struct {
	adminDB    *gorm.DB
	merchantDB *gorm.DB
	businessDB *gorm.DB
	identity   *identity.Service
}

func NewHandler(adminDB, merchantDB, businessDB *gorm.DB, id *identity.Service) *Handler {
	return &Handler{adminDB: adminDB, merchantDB: merchantDB, businessDB: businessDB, identity: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	h.RegisterMeta(r)
	r.GET("/products", h.list)
	r.GET("/products/:id", h.get)
	r.POST("/products/:id/audit", middleware.RequirePlatformMenu(h.identity, identity.PlatPermProductAudit), h.audit)
}

type productRow struct {
	ID           uint64    `gorm:"column:id"`
	StoreID      uint64    `gorm:"column:store_id"`
	MerchantID   uint64    `gorm:"column:merchant_id"`
	MerchantName string    `gorm:"column:merchant_name"`
	StoreName    string    `gorm:"column:store_name"`
	Title        string    `gorm:"column:title"`
	CategoryID   uint64    `gorm:"column:category_id"`
	Status       string    `gorm:"column:status"`
	Version      uint64    `gorm:"column:version"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

type detailRow struct {
	ProductID     uint64   `gorm:"column:product_id"`
	Brief         string   `gorm:"column:brief"`
	CoverURL      string   `gorm:"column:cover_url"`
	OriginalPrice *float64 `gorm:"column:original_price"`
}

type skuRow struct {
	ProductID uint64  `gorm:"column:product_id"`
	Price     float64 `gorm:"column:price"`
	Stock     int     `gorm:"column:stock"`
}

type productReview struct {
	ID         uint64 `gorm:"column:id;primaryKey"`
	ProductID  uint64 `gorm:"column:product_id"`
	StoreID    uint64 `gorm:"column:store_id"`
	Status     string `gorm:"column:status"`
	Reason     string `gorm:"column:reason"`
	ReviewedBy uint64 `gorm:"column:reviewed_by"`
}

func (productReview) TableName() string { return "qixi_crm_a_product_review" }

type categoryRow struct {
	ID   uint64 `gorm:"column:category_id"`
	Name string `gorm:"column:name"`
}

type productResponse struct {
	ProductID  uint64  `json:"product_id"`
	MerID      uint64  `json:"mer_id"`
	MerName    string  `json:"mer_name"`
	StoreName  string  `json:"store_name"`
	StoreInfo  string  `json:"store_info"`
	CateName   string  `json:"cate_name"`
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
	query := h.base(c)
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		query = query.Where("p.title LIKE ?", "%"+keyword+"%")
	}
	if merchantID, _ := strconv.ParseUint(c.Query("mer_id"), 10, 64); merchantID > 0 {
		query = query.Where("s.merchant_id = ?", merchantID)
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("p.status = ?", statusName(status))
	}
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
	var row productRow
	if err := h.base(c).Where("p.id = ?", id).Scan(&row).Error; err != nil {
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
	var row productRow
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_m_product AS p").Select("p.id,p.store_id,s.merchant_id,m.name AS merchant_name,s.name AS store_name,p.title,p.category_id,p.status,p.version,p.created_at").Joins("JOIN qixi_crm_m_store AS s ON s.id = p.store_id").Joins("JOIN qixi_crm_m_merchant AS m ON m.id = s.merchant_id").Where("p.id = ?", id).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID == 0 {
			return gorm.ErrRecordNotFound
		}
		if row.Status != "pending_review" && row.Status != "draft" {
			return errAuditState
		}
		next := "rejected"
		if req.Status == 1 {
			next = "on_sale"
		}
		return tx.Table("qixi_crm_m_product").Where("id = ?", id).Updates(map[string]any{"status": next, "version": gorm.Expr("version + 1")}).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	if errors.Is(err, errAuditState) {
		response.Fail(c, http.StatusConflict, "商品当前不在待审核状态")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新商品审核状态失败")
		return
	}
	reviewStatus := "rejected"
	if req.Status == 1 {
		reviewStatus = "approved"
	}
	review := productReview{ProductID: row.ID, StoreID: row.StoreID, Status: reviewStatus, Reason: strings.TrimSpace(req.Refusal), ReviewedBy: uint64(middleware.AdminID(c))}
	if err := h.adminDB.WithContext(c.Request.Context()).Create(&review).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "记录审核结果失败")
		return
	}
	if req.Status == -1 {
		_ = h.businessDB.WithContext(c.Request.Context()).Exec("DELETE FROM qixi_crm_b_product_view WHERE product_id = ?", id).Error
		response.OK(c, gin.H{"ok": true})
		return
	}
	row.Status, row.Version = "on_sale", row.Version+1
	items, err := h.responses(c, []productRow{row})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载商品详情失败")
		return
	}
	item := items[0]
	err = h.businessDB.WithContext(c.Request.Context()).Exec(`INSERT INTO qixi_crm_b_product_view
      (product_id, merchant_id, store_id, merchant_name, store_name, category_id, title, cover_url, price, original_price, product_type, sales, stock, sale_status, version, updated_at)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 0, ?, 1, ?, NOW())
      ON DUPLICATE KEY UPDATE merchant_id=VALUES(merchant_id), store_id=VALUES(store_id), merchant_name=VALUES(merchant_name), store_name=VALUES(store_name), category_id=VALUES(category_id), title=VALUES(title), cover_url=VALUES(cover_url), price=VALUES(price), original_price=VALUES(original_price), stock=VALUES(stock), sale_status=1, version=VALUES(version), updated_at=NOW()`,
		row.ID, row.MerchantID, row.StoreID, row.MerchantName, row.StoreName, row.CategoryID, row.Title, item.Image, item.Price, nullableOriginal(item.OtPrice), item.Stock, row.Version).Error
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "同步商品消费视图失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) base(c *gin.Context) *gorm.DB {
	return h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product AS p").Select("p.id,p.store_id,s.merchant_id,m.name AS merchant_name,s.name AS store_name,p.title,p.category_id,p.status,p.version,p.created_at").Joins("JOIN qixi_crm_m_store AS s ON s.id = p.store_id").Joins("JOIN qixi_crm_m_merchant AS m ON m.id = s.merchant_id")
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
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_sku").Where("product_id IN ? AND status = 1", ids).Find(&skus).Error; err != nil {
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
		out = append(out, productResponse{ProductID: row.ID, MerID: row.MerchantID, MerName: row.MerchantName, StoreName: row.Title, StoreInfo: d.Brief, CateName: nameByID[row.CategoryID], Image: d.CoverURL, Price: s.Price, OtPrice: ot, Stock: s.Stock, Sales: 0, Status: statusCode(row.Status), IsShow: showCode(row.Status), Refusal: reviewByID[row.ID].Reason, CreateTime: row.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	return out, nil
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

var errAuditState = errors.New("product audit state invalid")
