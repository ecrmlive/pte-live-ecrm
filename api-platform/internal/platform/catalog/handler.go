package catalog

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/catalog"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
)

type Handler struct {
	svc *catalog.Service
	id  *identity.Service
}

func NewHandler(svc *catalog.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	h.RegisterMeta(r)
}

// RegisterMeta keeps category and brand routes during the staged migration.
// Product audit routes are registered by nativecatalog against qixi_crm_*.
func (h *Handler) RegisterMeta(r gin.IRoutes) {
	r.GET("/product-categories", h.CategoryTree)
	r.POST("/product-categories", middleware.RequirePlatformMenu(h.id, identity.PlatPermCategoryManage), h.CreateCategory)
	r.PUT("/product-categories/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermCategoryManage), h.UpdateCategory)
	r.DELETE("/product-categories/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermCategoryManage), h.DeleteCategory)

	r.GET("/brands", h.ListBrands)
	r.POST("/brands", middleware.RequirePlatformMenu(h.id, identity.PlatPermBrandManage), h.CreateBrand)
	r.PUT("/brands/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermBrandManage), h.UpdateBrand)
	r.DELETE("/brands/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermBrandManage), h.DeleteBrand)

}

func (h *Handler) CategoryTree(c *gin.Context) {
	tree, err := h.svc.CategoryTree(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": tree})
}

type categoryReq struct {
	PID    uint   `json:"pid"`
	Name   string `json:"cate_name"`
	Sort   int    `json:"sort"`
	IsShow *int8  `json:"is_show"`
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var req categoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	isShow := int8(1)
	if req.IsShow != nil {
		isShow = *req.IsShow
	}
	row, err := h.svc.CreateCategory(c.Request.Context(), req.PID, req.Name, req.Sort, isShow)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req categoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	isShow := int8(1)
	if req.IsShow != nil {
		isShow = *req.IsShow
	}
	if err := h.svc.UpdateCategory(c.Request.Context(), uint(id), req.Name, req.Sort, isShow); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListBrands(c *gin.Context) {
	list, err := h.svc.ListBrands(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

type brandReq struct {
	Name   string `json:"brand_name"`
	Sort   int    `json:"sort"`
	IsShow *int8  `json:"is_show"`
}

func (h *Handler) CreateBrand(c *gin.Context) {
	var req brandReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	isShow := int8(1)
	if req.IsShow != nil {
		isShow = *req.IsShow
	}
	row, err := h.svc.CreateBrand(c.Request.Context(), req.Name, req.Sort, isShow)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateBrand(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req brandReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	isShow := int8(1)
	if req.IsShow != nil {
		isShow = *req.IsShow
	}
	if err := h.svc.UpdateBrand(c.Request.Context(), uint(id), req.Name, req.Sort, isShow); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) DeleteBrand(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteBrand(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var statusPtr *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		statusPtr = &st
	}
	var merPtr *uint
	if m := c.Query("mer_id"); m != "" {
		v, _ := strconv.ParseUint(m, 10, 64)
		id := uint(v)
		merPtr = &id
	}
	res, err := h.svc.ListProducts(c.Request.Context(), statusPtr, c.Query("keyword"), merPtr, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

type auditReq struct {
	Status  int8   `json:"status"`
	Refusal string `json:"refusal"`
}

func (h *Handler) AuditProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req auditReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.AuditProduct(c.Request.Context(), uint(id), req.Status, req.Refusal); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, catalog.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, catalog.ErrBadStatus),
		errors.Is(err, catalog.ErrRejectNeedMsg),
		errors.Is(err, catalog.ErrNameRequired):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
