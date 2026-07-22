package diy

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/diy"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *diy.Service
	id  *identity.Service
}

func NewHandler(svc *diy.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/diy/pages", h.List)
	r.GET("/diy/pages/:id", h.Get)
	r.POST("/diy/pages", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyCreate), h.Create)
	r.PUT("/diy/pages/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyUpdate), h.Update)
	r.POST("/diy/pages/:id/active", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyActive), h.SetActive)
	r.DELETE("/diy/pages/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyDelete), h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.List(c.Request.Context(), 0, page, limit)
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
	response.OK(c, p)
}

func (h *Handler) Create(c *gin.Context) {
	var in diy.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	p, err := h.svc.Create(c.Request.Context(), 0, in)
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
	p, err := h.svc.Update(c.Request.Context(), uint(id), 0, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) SetActive(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.SetActive(c.Request.Context(), uint(id), 0)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), uint(id), 0); err != nil {
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
