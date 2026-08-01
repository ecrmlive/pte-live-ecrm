package invoice

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/invoice"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct{ svc *invoice.Service }

func NewHandler(svc *invoice.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/invoices", h.List)
	r.PUT("/invoices/:id/audit", h.Audit)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListMerchant(c.Request.Context(), middleware.MerID(c), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Audit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in invoice.AuditInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Audit(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, invoice.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, invoice.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, invoice.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
