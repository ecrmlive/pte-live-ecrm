package nativeproductmeta

import (
	"errors"
	"net/http"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type category struct {
	ID       uint64 `gorm:"column:id;primaryKey"`
	StoreID  uint64 `gorm:"column:store_id"`
	ParentID uint64 `gorm:"column:parent_id"`
	Name     string `gorm:"column:name"`
	Sort     int    `gorm:"column:sort"`
	Status   int8   `gorm:"column:status"`
}

func (category) TableName() string { return "qixi_crm_m_category" }

type categoryInput struct {
	ParentID uint64 `json:"parent_id"`
	Name     string `json:"name"`
	Sort     int    `json:"sort"`
	Status   *int8  `json:"status"`
}

func (h *Handler) listCategories(c *gin.Context) {
	var rows []category
	if err := h.db.WithContext(c.Request.Context()).Where("store_id = ?", middleware.StoreID(c)).Order("sort DESC,id DESC").Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品分类失败")
		return
	}
	byParent := map[uint64][]gin.H{}
	for _, row := range rows {
		byParent[row.ParentID] = append(byParent[row.ParentID], categoryJSON(row))
	}
	var tree func(uint64) []gin.H
	tree = func(parent uint64) []gin.H {
		items := byParent[parent]
		for i := range items {
			items[i]["children"] = tree(items[i]["category_id"].(uint64))
		}
		return items
	}
	response.OK(c, gin.H{"list": tree(0)})
}
func (h *Handler) createCategory(c *gin.Context) {
	var req categoryInput
	if err := c.ShouldBindJSON(&req); err != nil || !validCategory(req) {
		response.Fail(c, http.StatusBadRequest, "商品分类参数不正确")
		return
	}
	if req.ParentID > 0 {
		if _, err := h.ownedCategory(c, req.ParentID); err != nil {
			response.Fail(c, http.StatusBadRequest, "上级分类不存在或不属于当前店铺")
			return
		}
	}
	row := category{StoreID: uint64(middleware.StoreID(c)), ParentID: req.ParentID, Name: strings.TrimSpace(req.Name), Sort: req.Sort, Status: enabled(req.Status)}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "新增商品分类失败")
		return
	}
	response.OK(c, categoryJSON(row))
}
func (h *Handler) updateCategory(c *gin.Context) {
	categoryID, ok := id(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "商品分类编号不正确")
		return
	}
	var req categoryInput
	if err := c.ShouldBindJSON(&req); err != nil || !validCategory(req) || req.ParentID == categoryID {
		response.Fail(c, http.StatusBadRequest, "商品分类参数不正确")
		return
	}
	row, err := h.ownedCategory(c, categoryID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品分类不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品分类失败")
		return
	}
	if req.ParentID > 0 {
		parent, parentErr := h.ownedCategory(c, req.ParentID)
		if parentErr != nil || h.categoryWouldCycle(c, categoryID, parent) {
			response.Fail(c, http.StatusBadRequest, "上级分类不合法")
			return
		}
	}
	row.ParentID, row.Name, row.Sort, row.Status = req.ParentID, strings.TrimSpace(req.Name), req.Sort, enabled(req.Status)
	if err := h.db.WithContext(c.Request.Context()).Save(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新商品分类失败")
		return
	}
	response.OK(c, categoryJSON(row))
}

func (h *Handler) categoryWouldCycle(c *gin.Context, categoryID uint64, parent category) bool {
	for depth := 0; depth < 32; depth++ {
		if parent.ID == categoryID {
			return true
		}
		if parent.ParentID == 0 {
			return false
		}
		next, err := h.ownedCategory(c, parent.ParentID)
		if err != nil {
			return true
		}
		parent = next
	}
	return true
}

func (h *Handler) deleteCategory(c *gin.Context) {
	categoryID, ok := id(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "商品分类编号不正确")
		return
	}
	row, err := h.ownedCategory(c, categoryID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品分类不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品分类失败")
		return
	}
	var children, products int64
	if err := h.db.WithContext(c.Request.Context()).Model(&category{}).Where("store_id = ? AND parent_id = ?", row.StoreID, row.ID).Count(&children).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "检查子分类失败")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_m_product").Where("store_id = ? AND category_id = ?", row.StoreID, row.ID).Count(&products).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "检查商品引用失败")
		return
	}
	if children > 0 || products > 0 {
		response.Fail(c, http.StatusConflict, "分类下存在子分类或商品，不能删除")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Delete(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除商品分类失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}
func (h *Handler) ownedCategory(c *gin.Context, categoryID uint64) (category, error) {
	var row category
	err := h.db.WithContext(c.Request.Context()).Where("id = ? AND store_id = ?", categoryID, middleware.StoreID(c)).Take(&row).Error
	return row, err
}
func categoryJSON(row category) gin.H {
	return gin.H{"category_id": row.ID, "parent_id": row.ParentID, "name": row.Name, "sort": row.Sort, "status": row.Status}
}
func validCategory(req categoryInput) bool {
	return len(strings.TrimSpace(req.Name)) > 0 && len(strings.TrimSpace(req.Name)) <= 128 && req.Sort >= 0
}
