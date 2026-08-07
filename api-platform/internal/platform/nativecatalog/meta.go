package nativecatalog

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type platformCategory struct {
	ID       uint64 `gorm:"column:id;primaryKey" json:"store_category_id"`
	ParentID uint64 `gorm:"column:parent_id" json:"pid"`
	Name     string `gorm:"column:name" json:"cate_name"`
	Sort     int    `gorm:"column:sort" json:"sort"`
	Status   int8   `gorm:"column:status" json:"is_show"`
}

func (platformCategory) TableName() string { return "qixi_crm_a_platform_category" }

type platformBrand struct {
	ID         uint64 `gorm:"column:id;primaryKey" json:"brand_id"`
	CategoryID uint64 `gorm:"column:category_id" json:"category_id"`
	Name       string `gorm:"column:name" json:"brand_name"`
	Sort       int    `gorm:"column:sort" json:"sort"`
	Status     int8   `gorm:"column:status" json:"is_show"`
}

func (platformBrand) TableName() string { return "qixi_crm_a_platform_brand" }

type platformBrandCategory struct {
	ID       uint64 `gorm:"column:id;primaryKey" json:"brand_category_id"`
	ParentID uint64 `gorm:"column:parent_id" json:"pid"`
	Name     string `gorm:"column:name" json:"cate_name"`
	Sort     int    `gorm:"column:sort" json:"sort"`
	Status   int8   `gorm:"column:status" json:"is_show"`
}

func (platformBrandCategory) TableName() string { return "qixi_crm_a_platform_brand_category" }

type categoryRequest struct {
	PID    uint64 `json:"pid"`
	Name   string `json:"cate_name"`
	Sort   int    `json:"sort"`
	IsShow *int8  `json:"is_show"`
}
type brandRequest struct {
	CategoryID uint64 `json:"category_id"`
	Name       string `json:"brand_name"`
	Sort       int    `json:"sort"`
	IsShow     *int8  `json:"is_show"`
}
type brandCategoryRequest struct {
	PID    uint64 `json:"pid"`
	Name   string `json:"cate_name"`
	Sort   int    `json:"sort"`
	IsShow *int8  `json:"is_show"`
}

func (h *Handler) RegisterMeta(r gin.IRoutes) {
	platformOnly := middleware.RequireAdminRoles("platform")
	r.GET("/product-categories", platformOnly, h.categories)
	r.POST("/product-categories", platformOnly, middleware.RequireAdminMenu(h.adminDB, "product.category.manage"), h.createCategory)
	r.PUT("/product-categories/:id", platformOnly, middleware.RequireAdminMenu(h.adminDB, "product.category.manage"), h.updateCategory)
	r.DELETE("/product-categories/:id", platformOnly, middleware.RequireAdminMenu(h.adminDB, "product.category.manage"), h.deleteCategory)
	r.GET("/brand-categories", platformOnly, h.brandCategories)
	r.POST("/brand-categories", platformOnly, middleware.RequireAdminMenu(h.adminDB, "product.brand.manage"), h.createBrandCategory)
	r.PUT("/brand-categories/:id", platformOnly, middleware.RequireAdminMenu(h.adminDB, "product.brand.manage"), h.updateBrandCategory)
	r.DELETE("/brand-categories/:id", platformOnly, middleware.RequireAdminMenu(h.adminDB, "product.brand.manage"), h.deleteBrandCategory)
	r.GET("/brands", platformOnly, h.brands)
	r.POST("/brands", platformOnly, middleware.RequireAdminMenu(h.adminDB, "product.brand.manage"), h.createBrand)
	r.PUT("/brands/:id", platformOnly, middleware.RequireAdminMenu(h.adminDB, "product.brand.manage"), h.updateBrand)
	r.DELETE("/brands/:id", platformOnly, middleware.RequireAdminMenu(h.adminDB, "product.brand.manage"), h.deleteBrand)
}

func (h *Handler) categories(c *gin.Context) {
	var rows []platformCategory
	if err := h.adminDB.WithContext(c.Request.Context()).
		Select("id,parent_id,name,sort,status").
		Order("sort ASC,id ASC").
		Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品分类失败")
		return
	}
	byParent := map[uint64][]gin.H{}
	for _, row := range rows {
		byParent[row.ParentID] = append(byParent[row.ParentID], gin.H{
			"store_category_id": row.ID,
			"pid":               row.ParentID,
			"cate_name":         row.Name,
			"sort":              row.Sort,
			"is_show":           row.Status,
		})
	}
	var build func(uint64) []gin.H
	build = func(parent uint64) []gin.H {
		items := byParent[parent]
		for i := range items {
			id := items[i]["store_category_id"].(uint64)
			if children := build(id); len(children) > 0 {
				items[i]["children"] = children
			}
		}
		return items
	}
	response.OK(c, gin.H{"list": build(0)})
}

func (h *Handler) createCategory(c *gin.Context) {
	var req categoryRequest
	if c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Name) == "" {
		response.Fail(c, http.StatusBadRequest, "分类参数不合法")
		return
	}
	status := int8(1)
	if req.IsShow != nil {
		status = *req.IsShow
	}
	if status != 0 && status != 1 {
		response.Fail(c, http.StatusBadRequest, "分类状态错误")
		return
	}
	if req.PID > 0 {
		var parent platformCategory
		if err := h.adminDB.WithContext(c.Request.Context()).Where("id = ?", req.PID).First(&parent).Error; err != nil {
			response.Fail(c, http.StatusBadRequest, "父分类不存在")
			return
		}
	}
	row := platformCategory{ParentID: req.PID, Name: strings.TrimSpace(req.Name), Sort: req.Sort, Status: status}
	if err := h.adminDB.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建分类失败")
		return
	}
	if err := h.syncCategory(c, row); err != nil {
		response.Fail(c, http.StatusInternalServerError, "同步消费分类失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) updateCategory(c *gin.Context) {
	id := parseUInt(c.Param("id"))
	var req categoryRequest
	if id == 0 || c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Name) == "" {
		response.Fail(c, http.StatusBadRequest, "分类参数不合法")
		return
	}
	var row platformCategory
	if err := h.adminDB.WithContext(c.Request.Context()).Where("id = ?", id).First(&row).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "分类不存在")
		return
	}
	row.Name = strings.TrimSpace(req.Name)
	row.Sort = req.Sort
	if req.IsShow != nil {
		row.Status = *req.IsShow
	}
	if row.Status != 0 && row.Status != 1 {
		response.Fail(c, http.StatusBadRequest, "分类状态错误")
		return
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Save(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新分类失败")
		return
	}
	if err := h.syncCategory(c, row); err != nil {
		response.Fail(c, http.StatusInternalServerError, "同步消费分类失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) deleteCategory(c *gin.Context) {
	id := parseUInt(c.Param("id"))
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "分类 ID 错误")
		return
	}
	var childCount, productCount int64
	if err := h.adminDB.WithContext(c.Request.Context()).Model(&platformCategory{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "检查分类失败")
		return
	}
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product").Where("category_id = ?", id).Count(&productCount).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "检查关联商品失败")
		return
	}
	if childCount > 0 || productCount > 0 {
		response.Fail(c, http.StatusConflict, "分类存在子分类或关联商品，不能删除")
		return
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Delete(&platformCategory{}, id).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除分类失败")
		return
	}
	if err := h.businessDB.WithContext(c.Request.Context()).Exec("DELETE FROM qixi_crm_b_category_view WHERE category_id = ?", id).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "同步消费分类失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) syncCategory(c *gin.Context, row platformCategory) error {
	return h.businessDB.WithContext(c.Request.Context()).Exec(`INSERT INTO qixi_crm_b_category_view (category_id,parent_id,name,sort,status,updated_at) VALUES (?, ?, ?, ?, ?, NOW()) ON DUPLICATE KEY UPDATE parent_id=VALUES(parent_id),name=VALUES(name),sort=VALUES(sort),status=VALUES(status),updated_at=NOW()`, row.ID, row.ParentID, row.Name, row.Sort, row.Status).Error
}

func (h *Handler) brands(c *gin.Context) {
	q := h.adminDB.WithContext(c.Request.Context()).Order("sort ASC,id ASC")
	if categoryID := parseUInt(c.Query("category_id")); categoryID > 0 {
		q = q.Where("category_id = ?", categoryID)
	}
	var rows []platformBrand
	if err := q.Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询品牌失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func (h *Handler) createBrand(c *gin.Context) {
	var req brandRequest
	if c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Name) == "" {
		response.Fail(c, http.StatusBadRequest, "品牌参数不合法")
		return
	}
	status := int8(1)
	if req.IsShow != nil {
		status = *req.IsShow
	}
	if req.CategoryID > 0 {
		var cat platformBrandCategory
		if err := h.adminDB.WithContext(c.Request.Context()).Where("id = ?", req.CategoryID).First(&cat).Error; err != nil {
			response.Fail(c, http.StatusBadRequest, "品牌分类不存在")
			return
		}
	}
	row := platformBrand{CategoryID: req.CategoryID, Name: strings.TrimSpace(req.Name), Sort: req.Sort, Status: status}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_platform_brand").Create(map[string]any{
		"category_id": row.CategoryID, "name": row.Name, "sort": row.Sort, "status": row.Status,
	}).Error; err != nil {
		response.Fail(c, http.StatusConflict, "品牌名称已存在")
		return
	}
	response.OK(c, row)
}
func (h *Handler) updateBrand(c *gin.Context) {
	id := parseUInt(c.Param("id"))
	var req brandRequest
	if id == 0 || c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Name) == "" {
		response.Fail(c, http.StatusBadRequest, "品牌参数不合法")
		return
	}
	if req.CategoryID > 0 {
		var cat platformBrandCategory
		if err := h.adminDB.WithContext(c.Request.Context()).Where("id = ?", req.CategoryID).First(&cat).Error; err != nil {
			response.Fail(c, http.StatusBadRequest, "品牌分类不存在")
			return
		}
	}
	updates := map[string]any{"name": strings.TrimSpace(req.Name), "sort": req.Sort, "category_id": req.CategoryID}
	if req.IsShow != nil {
		updates["status"] = *req.IsShow
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_platform_brand").Where("id = ?", id).Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新品牌失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}
func (h *Handler) deleteBrand(c *gin.Context) {
	id := parseUInt(c.Param("id"))
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "品牌 ID 错误")
		return
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Delete(&platformBrand{}, id).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除品牌失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) brandCategories(c *gin.Context) {
	var rows []platformBrandCategory
	if err := h.adminDB.WithContext(c.Request.Context()).Order("sort ASC,id ASC").Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询品牌分类失败")
		return
	}
	byParent := map[uint64][]gin.H{}
	for _, row := range rows {
		byParent[row.ParentID] = append(byParent[row.ParentID], gin.H{
			"brand_category_id": row.ID, "pid": row.ParentID, "cate_name": row.Name, "sort": row.Sort, "is_show": row.Status,
		})
	}
	var build func(uint64) []gin.H
	build = func(parent uint64) []gin.H {
		items := byParent[parent]
		for i := range items {
			id := items[i]["brand_category_id"].(uint64)
			items[i]["children"] = build(id)
		}
		return items
	}
	response.OK(c, gin.H{"list": build(0)})
}

func (h *Handler) createBrandCategory(c *gin.Context) {
	var req brandCategoryRequest
	if c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Name) == "" {
		response.Fail(c, http.StatusBadRequest, "品牌分类参数不合法")
		return
	}
	status := int8(1)
	if req.IsShow != nil {
		status = *req.IsShow
	}
	if status != 0 && status != 1 {
		response.Fail(c, http.StatusBadRequest, "品牌分类状态错误")
		return
	}
	if req.PID > 0 {
		var parent platformBrandCategory
		if err := h.adminDB.WithContext(c.Request.Context()).Where("id = ?", req.PID).First(&parent).Error; err != nil {
			response.Fail(c, http.StatusBadRequest, "父分类不存在")
			return
		}
	}
	row := platformBrandCategory{ParentID: req.PID, Name: strings.TrimSpace(req.Name), Sort: req.Sort, Status: status}
	if err := h.adminDB.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		response.Fail(c, http.StatusConflict, "品牌分类名称已存在")
		return
	}
	response.OK(c, row)
}

func (h *Handler) updateBrandCategory(c *gin.Context) {
	id := parseUInt(c.Param("id"))
	var req brandCategoryRequest
	if id == 0 || c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Name) == "" {
		response.Fail(c, http.StatusBadRequest, "品牌分类参数不合法")
		return
	}
	if req.PID == id {
		response.Fail(c, http.StatusBadRequest, "父分类不能是自身")
		return
	}
	var row platformBrandCategory
	if err := h.adminDB.WithContext(c.Request.Context()).Where("id = ?", id).First(&row).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "品牌分类不存在")
		return
	}
	updates := map[string]any{"parent_id": req.PID, "name": strings.TrimSpace(req.Name), "sort": req.Sort}
	if req.IsShow != nil {
		updates["status"] = *req.IsShow
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Model(&row).Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusConflict, "更新品牌分类失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) deleteBrandCategory(c *gin.Context) {
	id := parseUInt(c.Param("id"))
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "品牌分类 ID 错误")
		return
	}
	var childCount, brandCount int64
	ctx := c.Request.Context()
	if err := h.adminDB.WithContext(ctx).Model(&platformBrandCategory{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除品牌分类失败")
		return
	}
	if childCount > 0 {
		response.Fail(c, http.StatusConflict, "请先删除子分类")
		return
	}
	if err := h.adminDB.WithContext(ctx).Model(&platformBrand{}).Where("category_id = ?", id).Count(&brandCount).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除品牌分类失败")
		return
	}
	if brandCount > 0 {
		response.Fail(c, http.StatusConflict, "分类下仍有品牌，无法删除")
		return
	}
	if err := h.adminDB.WithContext(ctx).Delete(&platformBrandCategory{}, id).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除品牌分类失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func parseUInt(value string) uint64 { id, _ := strconv.ParseUint(value, 10, 64); return id }

var _ = gorm.ErrRecordNotFound
