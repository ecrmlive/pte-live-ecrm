package assist

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/assist"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct {
	svc *assist.Service
	id  *identity.Service
}

func NewHandler(svc *assist.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/assist/actives", h.List)
	r.POST("/assist/actives", middleware.RequireMerchantMenu(h.id, identity.MerPermAssistCreate), h.Create)
	r.PUT("/assist/actives/:id", h.Update)
	r.PUT("/assist/actives/:id/show", middleware.RequireMerchantMenu(h.id, identity.MerPermAssistToggle), h.SetShow)
	r.DELETE("/assist/actives/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermAssistDelete), h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	merID := middleware.MerID(c)
	res, err := h.svc.ListAdmin(c.Request.Context(), &merID, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	var in assist.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Create(c.Request.Context(), middleware.MerID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in assist.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if in.IsShow != nil {
		if err := h.id.RequireMerchantMenu(c.Request.Context(), middleware.AdminID(c), identity.MerPermAssistToggle); err != nil {
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		}
	}
	row, err := h.svc.Update(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

type showReq struct {
	IsShow *int `json:"is_show" binding:"required"`
}

func (h *Handler) SetShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req showReq
	if err := c.ShouldBindJSON(&req); err != nil || req.IsShow == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SetShow(c.Request.Context(), middleware.MerID(c), uint(id), *req.IsShow)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), middleware.MerID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, assist.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, assist.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
