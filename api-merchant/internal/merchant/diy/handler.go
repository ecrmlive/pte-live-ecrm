package diy

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/diy"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct {
	svc *diy.Service
}

func NewHandler(svc *diy.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/diy/pages", h.List)
	r.GET("/diy/pages/:id", h.Get)
	r.GET("/diy/defaults", h.ListDefaults)
	r.GET("/diy/editor/bootstrap", h.Bootstrap)
	r.GET("/diy/editor/bootstrap/:id", h.Bootstrap)
	r.POST("/diy/pages", h.Create)
	r.PUT("/diy/pages/:id", h.Update)
	r.POST("/diy/pages/:id/active", h.SetActive)
	r.POST("/diy/pages/:id/copy", h.Copy)
	r.POST("/diy/pages/:id/recovery", h.Recovery)
	r.POST("/diy/defaults/:id/apply", h.ApplyDefault)
	r.DELETE("/diy/pages/:id", h.Delete)
	// 商户端只能读取平台维护的商户商城链接，不具备配置权限。
	r.GET("/diy/page-categories", h.ListLinkCategories)
	r.GET("/diy/page-links", h.ListPageLinks)
}

func (h *Handler) ListLinkCategories(c *gin.Context) {
	rows, err := h.svc.ListCategories(c.Request.Context(), 1)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) ListPageLinks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	f := diy.LinkListFilter{IsMer: 1, Page: page, Limit: limit, Name: c.Query("name")}
	active := int8(1)
	f.Status = &active
	rows, err := h.svc.ListLinks(c.Request.Context(), f)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, rows)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	f := diy.ListFilter{MerID: middleware.StoreID(c), Page: page, Limit: limit, Name: c.Query("name")}
	if v := c.Query("is_diy"); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		iv := int8(n)
		f.IsDiy = &iv
	}
	if v := c.Query("status"); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		iv := int8(n)
		f.Status = &iv
	}
	res, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) ListDefaults(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListDefaults(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	if p.StoreID != middleware.StoreID(c) {
		response.Fail(c, http.StatusNotFound, diy.ErrNotFound.Error())
		return
	}
	response.OK(c, p)
}

func (h *Handler) Bootstrap(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		if q := c.Query("id"); q != "" {
			id, _ = strconv.ParseUint(q, 10, 64)
		}
	}
	boot, err := h.svc.EditorBootstrap(c.Request.Context(), uint(id), middleware.StoreID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, boot)
}

func (h *Handler) Create(c *gin.Context) {
	var in diy.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	p, err := h.svc.Create(c.Request.Context(), middleware.StoreID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in diy.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	p, err := h.svc.Update(c.Request.Context(), uint(id), middleware.StoreID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) SetActive(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.SetActive(c.Request.Context(), uint(id), middleware.StoreID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Copy(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Copy(c.Request.Context(), uint(id), middleware.StoreID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) ApplyDefault(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.ApplyDefault(c.Request.Context(), uint(id), middleware.StoreID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Recovery(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Recovery(c.Request.Context(), uint(id), middleware.StoreID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), uint(id), middleware.StoreID(c)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, diy.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, diy.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
