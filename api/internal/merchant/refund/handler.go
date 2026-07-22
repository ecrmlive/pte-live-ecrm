package refund

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/aftersale"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *aftersale.Service
	id  *identity.Service
}

func NewHandler(svc *aftersale.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/refunds", h.List)
	r.GET("/refunds/:id", h.Get)
	r.POST("/refunds/:id/approve", middleware.RequireMerchantMenu(h.id, identity.MerPermRefundApprove), h.Approve)
	r.POST("/refunds/:id/reject", middleware.RequireMerchantMenu(h.id, identity.MerPermRefundReject), h.Reject)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var status *int8
	if s := c.Query("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			vv := int8(v)
			status = &vv
		}
	}
	res, err := h.svc.ListMerchant(c.Request.Context(), middleware.MerID(c), status, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ro, err := h.svc.GetMerchant(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, ro)
}

func (h *Handler) Approve(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Approve(c.Request.Context(), middleware.MerID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Reject(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in aftersale.RejectInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.Reject(c.Request.Context(), middleware.MerID(c), uint(id), in); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, aftersale.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, aftersale.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, aftersale.ErrBadParam),
		errors.Is(err, aftersale.ErrBadStatus),
		errors.Is(err, aftersale.ErrAlreadyDone),
		errors.Is(err, aftersale.ErrRejectNeedMsg),
		errors.Is(err, aftersale.ErrProductInvalid):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
