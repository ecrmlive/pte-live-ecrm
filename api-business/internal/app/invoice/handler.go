package invoice

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/invoice"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *invoice.Service }

func NewHandler(svc *invoice.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/invoice-profiles", h.ListProfiles)
	r.POST("/invoice-profiles", h.CreateProfile)
	r.PUT("/invoice-profiles/:id", h.UpdateProfile)
	r.DELETE("/invoice-profiles/:id", h.DeleteProfile)
	r.POST("/invoice-profiles/:id/default", h.SetDefaultProfile)
	r.GET("/invoices", h.List)
	r.GET("/invoices/:id", h.Get)
	r.POST("/invoices", h.Apply)
}

func (h *Handler) ListProfiles(c *gin.Context) {
	rows, err := h.svc.ListProfiles(c.Request.Context(), uint64(middleware.UID(c)))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) CreateProfile(c *gin.Context) {
	var in invoice.ProfileInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateProfile(c.Request.Context(), uint64(middleware.UID(c)), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "发票抬头 ID 错误")
		return
	}
	var in invoice.ProfileInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateProfile(c.Request.Context(), uint64(middleware.UID(c)), id, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteProfile(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "发票抬头 ID 错误")
		return
	}
	if err := h.svc.DeleteProfile(c.Request.Context(), uint64(middleware.UID(c)), id); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) SetDefaultProfile(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "发票抬头 ID 错误")
		return
	}
	if err := h.svc.SetDefaultProfile(c.Request.Context(), uint64(middleware.UID(c)), id); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListMine(c.Request.Context(), uint64(middleware.UID(c)), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "发票 ID 错误")
		return
	}
	row, err := h.svc.GetMine(c.Request.Context(), uint64(middleware.UID(c)), id)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Apply(c *gin.Context) {
	var in invoice.ApplyInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Apply(c.Request.Context(), uint64(middleware.UID(c)), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func positiveID(raw string) (uint64, bool) {
	id, err := strconv.ParseUint(raw, 10, 64)
	return id, err == nil && id > 0
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, invoice.ErrBadParam), errors.Is(err, invoice.ErrExists), errors.Is(err, invoice.ErrOrder):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, invoice.ErrNotFound), errors.Is(err, invoice.ErrProfileNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
