// Package servicemeal provides platform CRUD for 一号通服务套餐。
package servicemeal

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	tableName = "qixi_crm_a_serve_meal"

	menuRead   = "systemServeMerMealLst"
	menuDetail = "systemServeMealDetail"
	menuCreate = "systemServeMealCreate"
	menuUpdate = "systemServeMealUpdate"
	menuDelete = "systemServeMealDelete"
	menuStatus = "systemServeMealStatus"
)

type Handler struct {
	adminDB *gorm.DB
}

func NewHandler(adminDB *gorm.DB) *Handler {
	return &Handler{adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform")
	read := middleware.RequireAdminMenuAny(h.adminDB, menuRead, menuDetail)

	r.GET("/service-meals", access, read, h.List)
	r.GET("/service-meals/:id", access, read, h.Detail)
	r.POST("/service-meals", access, middleware.RequireAdminMenu(h.adminDB, menuCreate), h.Create)
	r.PUT("/service-meals/:id", access, middleware.RequireAdminMenu(h.adminDB, menuUpdate), h.Update)
	r.PUT("/service-meals/:id/status", access, middleware.RequireAdminMenu(h.adminDB, menuStatus), h.SetStatus)
	r.DELETE("/service-meals/:id", access, middleware.RequireAdminMenu(h.adminDB, menuDelete), h.Delete)
}

type mealRow struct {
	MealID     uint64    `gorm:"column:meal_id" json:"meal_id"`
	Name       string    `gorm:"column:name" json:"name"`
	Type       int8      `gorm:"column:type" json:"type"`
	Price      float64   `gorm:"column:price" json:"price"`
	Num        int       `gorm:"column:num" json:"num"`
	Sort       int       `gorm:"column:sort" json:"sort"`
	Status     int8      `gorm:"column:status" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

type saveRequest struct {
	Name   string   `json:"name"`
	Type   *int8    `json:"type"`
	Price  *float64 `json:"price"`
	Num    *int     `json:"num"`
	Sort   *int     `json:"sort"`
	Status *int8    `json:"status"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := pageLimit(c)
	q := h.adminDB.WithContext(c.Request.Context()).Table(tableName).Where("is_del = ?", 0)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询服务套餐失败")
		return
	}
	var rows []mealRow
	if err := q.Order("sort ASC, meal_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询服务套餐失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Detail(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var row mealRow
	err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).Where("meal_id = ? AND is_del = ?", id, 0).Take(&row).Error
	if err != nil {
		respondDBError(c, err, "服务套餐不存在")
		return
	}
	response.OK(c, row)
}

func (h *Handler) Create(c *gin.Context) {
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	values, ok := normalize(req, c)
	if !ok {
		return
	}
	values["create_time"] = time.Now()
	if err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).Create(values).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "新增服务套餐失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	values, ok := normalize(req, c)
	if !ok {
		return
	}
	query := h.adminDB.WithContext(c.Request.Context()).Table(tableName).Where("meal_id = ? AND is_del = ?", id, 0)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询服务套餐失败")
		return
	}
	if total == 0 {
		response.Fail(c, http.StatusNotFound, "服务套餐不存在")
		return
	}

	result := query.Updates(values)
	if result.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "保存服务套餐失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) SetStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req struct {
		Status *int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == nil || (*req.Status != 0 && *req.Status != 1) {
		response.Fail(c, http.StatusBadRequest, "状态参数错误")
		return
	}
	result := h.adminDB.WithContext(c.Request.Context()).Table(tableName).Where("meal_id = ? AND is_del = ?", id, 0).Update("status", *req.Status)
	if result.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "更新服务套餐状态失败")
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "服务套餐不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	result := h.adminDB.WithContext(c.Request.Context()).Table(tableName).Where("meal_id = ? AND is_del = ?", id, 0).Update("is_del", 1)
	if result.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "删除服务套餐失败")
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "服务套餐不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func normalize(req saveRequest, c *gin.Context) (map[string]any, bool) {
	name := strings.TrimSpace(req.Name)
	if name == "" || utf8.RuneCountInString(name) > 30 {
		response.Fail(c, http.StatusBadRequest, "套餐名称不能为空且不能超过 30 个字符")
		return nil, false
	}
	if req.Type == nil || (*req.Type != 1 && *req.Type != 2) {
		response.Fail(c, http.StatusBadRequest, "套餐类型错误")
		return nil, false
	}
	if req.Price == nil || *req.Price < 0 || *req.Price > 999_999.99 {
		response.Fail(c, http.StatusBadRequest, "价格范围错误")
		return nil, false
	}
	if req.Num == nil || *req.Num < 0 {
		response.Fail(c, http.StatusBadRequest, "数量不能小于 0")
		return nil, false
	}
	if req.Sort == nil {
		response.Fail(c, http.StatusBadRequest, "请填写排序")
		return nil, false
	}
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}
	if status != 0 && status != 1 {
		response.Fail(c, http.StatusBadRequest, "状态参数错误")
		return nil, false
	}
	return map[string]any{
		"name": name, "type": *req.Type, "price": *req.Price, "num": *req.Num,
		"sort": *req.Sort, "status": status,
	}, true
}

func pageLimit(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return page, limit
}

func parseID(c *gin.Context) (uint64, bool) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return 0, false
	}
	return id, true
}

func respondDBError(c *gin.Context, err error, notFound string) {
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, http.StatusNotFound, notFound)
		return
	}
	response.Fail(c, http.StatusInternalServerError, "查询服务套餐失败")
}
