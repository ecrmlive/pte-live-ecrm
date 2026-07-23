package diy

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/diy"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
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
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	f := diy.ListFilter{MerID: middleware.MerID(c), Page: page, Limit: limit, Name: c.Query("name")}
	if v := c.Query("is_diy"); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		iv := int8(n)
		f.IsDiy = &iv
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
	if p.MerID != middleware.MerID(c) {
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
	boot, err := h.svc.EditorBootstrap(c.Request.Context(), uint(id), middleware.MerID(c))
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
	p, err := h.svc.Create(c.Request.Context(), middleware.MerID(c), in)
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
	p, err := h.svc.Update(c.Request.Context(), uint(id), middleware.MerID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) SetActive(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.SetActive(c.Request.Context(), uint(id), middleware.MerID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Copy(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Copy(c.Request.Context(), uint(id), middleware.MerID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) ApplyDefault(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.ApplyDefault(c.Request.Context(), uint(id), middleware.MerID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Recovery(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Recovery(c.Request.Context(), uint(id), middleware.MerID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), uint(id), middleware.MerID(c)); err != nil {
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
