// Package nativeproductmeta owns store product metadata in qixi_crm_m_ tables.
package nativeproductmeta

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }
func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/store-categories", h.listCategories)
	r.POST("/store-categories", middleware.RequireStorePermission(h.db, "product.category.create"), h.createCategory)
	r.PUT("/store-categories/:id", middleware.RequireStorePermission(h.db, "product.category.update"), h.updateCategory)
	r.DELETE("/store-categories/:id", middleware.RequireStorePermission(h.db, "product.category.delete"), h.deleteCategory)
	r.GET("/product/labels", h.listLabels)
	r.POST("/product/labels", middleware.RequireStorePermission(h.db, "product.label.create"), h.createLabel)
	r.PUT("/product/labels/:id", middleware.RequireStorePermission(h.db, "product.label.update"), h.updateLabel)
	r.DELETE("/product/labels/:id", middleware.RequireStorePermission(h.db, "product.label.delete"), h.deleteLabel)

	// 店铺商品参数模板（CRMEB merchantStoreParameterTemplate*；写操作按 store_id 隔离）
	r.GET("/product/parameter-templates", h.listParameterTemplates)
	r.GET("/product/parameter-templates/:id", h.detailParameterTemplate)
	r.POST("/product/parameter-templates", h.createParameterTemplate)
	r.PUT("/product/parameter-templates/:id", h.updateParameterTemplate)
	r.DELETE("/product/parameter-templates/:id", h.deleteParameterTemplate)
}

type tag struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	StoreID   uint64    `gorm:"column:store_id"`
	Name      string    `gorm:"column:name"`
	Info      string    `gorm:"column:info"`
	Sort      int       `gorm:"column:sort"`
	Status    int8      `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (tag) TableName() string { return "qixi_crm_m_product_tag" }

type tagInput struct {
	Name   string `json:"name"`
	Info   string `json:"info"`
	Sort   int    `json:"sort"`
	Status *int8  `json:"status"`
}

func (h *Handler) listLabels(c *gin.Context) {
	page, limit := pagination(c)
	q := h.db.WithContext(c.Request.Context()).Model(&tag{}).Where("store_id = ?", middleware.StoreID(c))
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	if statusRaw := strings.TrimSpace(c.Query("status")); statusRaw != "" {
		if status, err := strconv.Atoi(statusRaw); err == nil && (status == 0 || status == 1) {
			q = q.Where("status = ?", status)
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品标签失败")
		return
	}
	var rows []tag
	if err := q.Order("sort DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品标签失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, tagJSON(row))
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}
func (h *Handler) createLabel(c *gin.Context) {
	var req tagInput
	if err := c.ShouldBindJSON(&req); err != nil || !valid(req) {
		response.Fail(c, http.StatusBadRequest, "商品标签参数不正确")
		return
	}
	row := tag{StoreID: uint64(middleware.StoreID(c)), Name: strings.TrimSpace(req.Name), Info: strings.TrimSpace(req.Info), Sort: req.Sort, Status: enabled(req.Status)}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "商品标签名称已存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "新增商品标签失败")
		return
	}
	response.OK(c, tagJSON(row))
}
func (h *Handler) updateLabel(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "商品标签编号不正确")
		return
	}
	var req tagInput
	if err := c.ShouldBindJSON(&req); err != nil || !valid(req) {
		response.Fail(c, http.StatusBadRequest, "商品标签参数不正确")
		return
	}
	row, err := h.owned(c, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品标签不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品标签失败")
		return
	}
	row.Name, row.Info, row.Sort, row.Status = strings.TrimSpace(req.Name), strings.TrimSpace(req.Info), req.Sort, enabled(req.Status)
	if err := h.db.WithContext(c.Request.Context()).Save(&row).Error; err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "商品标签名称已存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "更新商品标签失败")
		return
	}
	response.OK(c, tagJSON(row))
}
func (h *Handler) deleteLabel(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "商品标签编号不正确")
		return
	}
	row, err := h.owned(c, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品标签不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品标签失败")
		return
	}
	var bindings int64
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_m_product_tag_binding").Where("tag_id = ?", row.ID).Count(&bindings).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "检查商品标签引用失败")
		return
	}
	if bindings > 0 {
		response.Fail(c, http.StatusConflict, "该标签已关联商品，不能删除")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Delete(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除商品标签失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}
func (h *Handler) owned(c *gin.Context, tagID uint64) (tag, error) {
	var row tag
	err := h.db.WithContext(c.Request.Context()).Where("id = ? AND store_id = ?", tagID, middleware.StoreID(c)).Take(&row).Error
	return row, err
}
func tagJSON(row tag) gin.H {
	return gin.H{"label_id": row.ID, "name": row.Name, "info": row.Info, "sort": row.Sort, "status": row.Status, "create_time": row.CreatedAt.Format("2006-01-02 15:04:05")}
}
func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}
func id(c *gin.Context) (uint64, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return value, err == nil && value > 0
}
func enabled(value *int8) int8 {
	if value != nil && *value == 0 {
		return 0
	}
	return 1
}
func valid(req tagInput) bool {
	return len(strings.TrimSpace(req.Name)) > 0 && len(strings.TrimSpace(req.Name)) <= 64 && len(strings.TrimSpace(req.Info)) <= 255 && req.Sort >= 0
}
func isDuplicate(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
