package refund

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/aftersale"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/trade"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct {
	after *aftersale.Service
	trade *trade.Service
	id    *identity.Service
}

func NewHandler(afterSvc *aftersale.Service, tradeSvc *trade.Service, idSvc *identity.Service) *Handler {
	return &Handler{after: afterSvc, trade: tradeSvc, id: idSvc}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/refunds", h.Apply)
	r.GET("/refunds", h.List)
	r.POST("/refunds/:id/approve", h.Approve)
	r.POST("/refunds/:id/reject", h.Reject)
}

type applyReq struct {
	OrderID       uint   `json:"order_id" binding:"required"`
	RefundMessage string `json:"refund_message" binding:"required"`
	RefundType    int8   `json:"refund_type" binding:"required"`
}

func (h *Handler) Apply(c *gin.Context) {
	var req applyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if req.RefundType != aftersale.RefundTypeMoneyOnly {
		response.Fail(c, http.StatusBadRequest, "仅支持仅退款")
		return
	}
	ro, err := h.after.ApplyBehalf(c.Request.Context(), middleware.MerID(c), aftersale.ApplyInput{
		OrderID:       req.OrderID,
		RefundType:    aftersale.RefundTypeMoneyOnly,
		RefundMessage: req.RefundMessage,
	})
	if err != nil {
		writeAfterErr(c, err)
		return
	}
	response.OK(c, ro)
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
	res, err := h.after.ListMerchant(c.Request.Context(), middleware.MerID(c), status, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Approve(c *gin.Context) {
	if err := h.id.RequireStoreVerifyPerm(c.Request.Context(), middleware.AdminID(c)); err != nil {
		writeIDErr(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.after.Approve(c.Request.Context(), middleware.MerID(c), uint(id)); err != nil {
		writeAfterErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Reject(c *gin.Context) {
	if err := h.id.RequireStoreVerifyPerm(c.Request.Context(), middleware.AdminID(c)); err != nil {
		writeIDErr(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in aftersale.RejectInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.after.Reject(c.Request.Context(), middleware.MerID(c), uint(id), in); err != nil {
		writeAfterErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeIDErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, identity.ErrNoVerifyPerm),
		errors.Is(err, identity.ErrAccountDisabled):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, identity.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}

func writeTradeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, trade.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, trade.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}

func writeAfterErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, aftersale.ErrNotFound), errors.Is(err, aftersale.ErrOrderNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, aftersale.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, aftersale.ErrBadParam),
		errors.Is(err, aftersale.ErrBadStatus),
		errors.Is(err, aftersale.ErrAlreadyDone),
		errors.Is(err, aftersale.ErrOrderNotPaid),
		errors.Is(err, aftersale.ErrOrderRefunded),
		errors.Is(err, aftersale.ErrRefundInProgress),
		errors.Is(err, aftersale.ErrRejectNeedMsg),
		errors.Is(err, aftersale.ErrProductInvalid):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
