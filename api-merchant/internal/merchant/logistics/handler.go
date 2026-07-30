package logistics

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/logistics"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/response"
)

// 菜单：运费模板可挂 sql/043 导入节点；本刀 JWT + mer_id 隔离。
type Handler struct{ svc *logistics.Service }

func NewHandler(svc *logistics.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/express", h.ListExpress)
	r.GET("/city", h.ListCity)
	r.GET("/shipping/templates", h.ListTemplate)
	r.GET("/shipping/templates/:id", h.GetTemplate)
	r.POST("/shipping/templates", h.CreateTemplate)
	r.PUT("/shipping/templates/:id", h.UpdateTemplate)
	r.POST("/shipping/templates/:id/default", h.SetDefaultTemplate)
	r.DELETE("/shipping/templates/:id", h.DeleteTemplate)
}

func (h *Handler) ListExpress(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	res, err := h.svc.ListExpress(c.Request.Context(), page, limit, true)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) ListCity(c *gin.Context) {
	var parent *uint
	if raw := c.Query("parent_id"); raw != "" {
		v, err := strconv.ParseUint(raw, 10, 64)
		if err == nil {
			u := uint(v)
			parent = &u
		}
	}
	list, err := h.svc.ListCity(c.Request.Context(), parent)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) ListTemplate(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListTemplate(c.Request.Context(), middleware.MerID(c), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.GetTemplate(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	var in logistics.TemplateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateTemplate(c.Request.Context(), middleware.MerID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in logistics.TemplateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateTemplate(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SetDefaultTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.SetDefaultTemplate(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteTemplate(c.Request.Context(), middleware.MerID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, logistics.ErrBadParam), errors.Is(err, logistics.ErrDefaultTemplate):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, logistics.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
