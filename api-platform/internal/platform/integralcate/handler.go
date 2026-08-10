// Package integralcate exposes platform CRUD for points-mall categories
// (CRMEB admin.points.Category / eb_store_category.type=1).
package integralcate

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
	// RequireAdminMenu 仅认 kind=button；page 码 marketing.integral.classify 只用于导航。
	menuClassifyManage = "marketing.integral.classify.manage"
	cateType           = 1
)

type Handler struct {
	businessDB *gorm.DB
	adminDB    *gorm.DB
}

func NewHandler(businessDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	menu := middleware.RequireAdminMenu(h.adminDB, menuClassifyManage)
	r.GET("/points/categories", access, menu, h.List)
	r.GET("/points/categories/select", access, menu, h.Select)
	r.POST("/points/categories", access, menu, h.Create)
	r.PUT("/points/categories/:id", access, menu, h.Update)
	r.PUT("/points/categories/:id/status", access, menu, h.SwitchStatus)
	r.DELETE("/points/categories/:id", access, menu, h.Delete)
}

type categoryRow struct {
	StoreCategoryID uint64    `gorm:"column:store_category_id" json:"store_category_id"`
	PID             uint64    `gorm:"column:pid" json:"pid"`
	CateName        string    `gorm:"column:cate_name" json:"cate_name"`
	Path            string    `gorm:"column:path" json:"path"`
	Sort            int       `gorm:"column:sort" json:"sort"`
	Pic             string    `gorm:"column:pic" json:"pic"`
	IsShow          int8      `gorm:"column:is_show" json:"is_show"`
	Level           uint      `gorm:"column:level" json:"level"`
	MerID           uint64    `gorm:"column:mer_id" json:"mer_id"`
	CreateTime      time.Time `gorm:"column:create_time" json:"create_time"`
	IsHot           int8      `gorm:"column:is_hot" json:"is_hot"`
	Type            int8      `gorm:"column:type" json:"type"`
	HasProduct      int8      `json:"has_product"`
}

type saveInput struct {
	CateName string  `json:"cate_name"`
	IsShow   *int    `json:"is_show"`
	Sort     *int    `json:"sort"`
	Pic      string  `json:"pic"`
	PID      *uint64 `json:"pid"`
}

func (h *Handler) List(c *gin.Context) {
	rows, err := h.loadAll(c)
	if err != nil {
		fail(c, "积分商品分类查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": len(rows)})
}

func (h *Handler) Select(c *gin.Context) {
	rows := make([]categoryRow, 0)
	err := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_points_category").
		Where("mer_id = 0 AND type = ? AND is_del = 0 AND is_show = 1", cateType).
		Order("sort DESC, store_category_id DESC").
		Select("store_category_id,cate_name,sort,is_show").
		Scan(&rows).Error
	if err != nil {
		fail(c, "积分商品分类筛选查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) Create(c *gin.Context) {
	var in saveInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "积分商品分类参数错误")
		return
	}
	name, show, sort, pic, pid, ok := normalizeSave(in)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "请填写分类名称，并将排序限制在 0～99999")
		return
	}
	row := map[string]any{
		"pid":       pid,
		"cate_name": name,
		"path":      "/",
		"sort":      sort,
		"pic":       pic,
		"is_show":   show,
		"level":     0,
		"mer_id":    0,
		"is_hot":    0,
		"type":      cateType,
		"is_del":    0,
	}
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_category").Create(row).Error; err != nil {
		fail(c, "积分商品分类添加失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in saveInput
	if id == 0 || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "积分商品分类参数错误")
		return
	}
	name, show, sort, pic, pid, ok := normalizeSave(in)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "请填写分类名称，并将排序限制在 0～99999")
		return
	}
	if !h.exists(c, id) {
		response.Fail(c, http.StatusNotFound, "数据不存在")
		return
	}
	changes := map[string]any{
		"cate_name": name,
		"is_show":   show,
		"sort":      sort,
		"pic":       pic,
		"pid":       pid,
	}
	if err := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_points_category").
		Where("store_category_id = ? AND is_del = 0 AND type = ?", id, cateType).
		Updates(changes).Error; err != nil {
		fail(c, "积分商品分类编辑失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) SwitchStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in struct {
		Status *int `json:"status"`
		IsShow *int `json:"is_show"`
	}
	if id == 0 || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "积分商品分类状态参数错误")
		return
	}
	raw := in.IsShow
	if raw == nil {
		raw = in.Status
	}
	if raw == nil || (*raw != 0 && *raw != 1) {
		response.Fail(c, http.StatusBadRequest, "积分商品分类状态参数错误")
		return
	}
	if !h.exists(c, id) {
		response.Fail(c, http.StatusNotFound, "数据不存在")
		return
	}
	if err := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_points_category").
		Where("store_category_id = ? AND is_del = 0 AND type = ?", id, cateType).
		Update("is_show", *raw).Error; err != nil {
		fail(c, "积分商品分类状态修改失败")
		return
	}
	response.OK(c, gin.H{"ok": true, "is_show": *raw})
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "积分商品分类参数错误")
		return
	}
	if !h.exists(c, id) {
		response.Fail(c, http.StatusNotFound, "数据不存在")
		return
	}
	var child int64
	if err := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_points_category").
		Where("pid = ? AND is_del = 0 AND type = ?", id, cateType).
		Count(&child).Error; err != nil {
		fail(c, "积分商品分类删除失败")
		return
	}
	if child > 0 {
		response.Fail(c, http.StatusBadRequest, "该分类存在子集，请先处理子集")
		return
	}
	if err := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_points_category").
		Where("store_category_id = ? AND is_del = 0 AND type = ?", id, cateType).
		Update("is_del", 1).Error; err != nil {
		fail(c, "积分商品分类删除失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) loadAll(c *gin.Context) ([]categoryRow, error) {
	rows := make([]categoryRow, 0)
	err := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_points_category").
		Where("mer_id = 0 AND type = ? AND is_del = 0", cateType).
		Order("sort DESC, store_category_id DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return rows, nil
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.StoreCategoryID)
	}
	type countRow struct {
		CateID uint64 `gorm:"column:cate_id"`
		Cnt    int64  `gorm:"column:cnt"`
	}
	counts := make([]countRow, 0)
	_ = h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_points_product_view").
		Select("cate_id, COUNT(*) AS cnt").
		Where("cate_id IN ?", ids).
		Group("cate_id").
		Scan(&counts).Error
	byID := make(map[uint64]int64, len(counts))
	for _, item := range counts {
		byID[item.CateID] = item.Cnt
	}
	for i := range rows {
		if byID[rows[i].StoreCategoryID] > 0 {
			rows[i].HasProduct = 1
		}
	}
	return rows, nil
}

func (h *Handler) exists(c *gin.Context, id uint64) bool {
	var n int64
	_ = h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_points_category").
		Where("store_category_id = ? AND mer_id = 0 AND type = ? AND is_del = 0", id, cateType).
		Count(&n).Error
	return n > 0
}

func normalizeSave(in saveInput) (name string, show, sort int, pic string, pid uint64, ok bool) {
	name = strings.TrimSpace(in.CateName)
	if name == "" || utf8.RuneCountInString(name) > 100 {
		return "", 0, 0, "", 0, false
	}
	show = 1
	if in.IsShow != nil {
		if *in.IsShow != 0 && *in.IsShow != 1 {
			return "", 0, 0, "", 0, false
		}
		show = *in.IsShow
	}
	sort = 0
	if in.Sort != nil {
		if *in.Sort < 0 || *in.Sort > 99999 {
			return "", 0, 0, "", 0, false
		}
		sort = *in.Sort
	}
	pic = strings.TrimSpace(in.Pic)
	if utf8.RuneCountInString(pic) > 128 {
		return "", 0, 0, "", 0, false
	}
	if in.PID != nil {
		pid = *in.PID
	}
	return name, show, sort, pic, pid, true
}

func fail(c *gin.Context, message string) {
	response.Fail(c, http.StatusInternalServerError, message)
}
