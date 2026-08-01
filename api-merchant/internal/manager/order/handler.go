package order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/trade"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct {
	trade *trade.Service
	id    *identity.Service
}

func NewHandler(tradeSvc *trade.Service, idSvc *identity.Service) *Handler {
	return &Handler{trade: tradeSvc, id: idSvc}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/orders", h.List)
	r.GET("/orders/code/:code", h.GetByCode)
	r.GET("/orders/:id", h.Get)
	r.POST("/orders/:id/verify", h.Verify)
	r.POST("/orders/:id/delivery", h.Deliver)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if c.Query("await_verify") == "1" {
		res, err := h.trade.MerchantListAwaitVerify(c.Request.Context(), middleware.MerID(c), page, limit)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询失败")
			return
		}
		response.OK(c, res)
		return
	}
	var paid, status *int8
	if s := c.Query("paid"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			vv := int8(v)
			paid = &vv
		}
	}
	if s := c.Query("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			vv := int8(v)
			status = &vv
		}
	}
	res, err := h.trade.MerchantList(c.Request.Context(), middleware.MerID(c), paid, status, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	o, err := h.trade.GetMerchantOrder(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, o)
}

func (h *Handler) GetByCode(c *gin.Context) {
	o, err := h.trade.GetStoreOrderByVerifyCode(c.Request.Context(), middleware.MerID(c), c.Param("code"))
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, o)
}

type verifyReq struct {
	VerifyCode string `json:"verify_code"`
}

func (h *Handler) Verify(c *gin.Context) {
	if err := h.id.RequireStoreVerifyPerm(c.Request.Context(), middleware.AdminID(c)); err != nil {
		writeIDErr(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req verifyReq
	_ = c.ShouldBindJSON(&req)
	if err := h.trade.VerifyByStaff(c.Request.Context(), middleware.MerID(c), uint(id), middleware.AdminID(c), req.VerifyCode); err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Deliver(c *gin.Context) {
	if err := h.id.RequireStoreDeliverPerm(c.Request.Context(), middleware.AdminID(c)); err != nil {
		writeIDErr(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in trade.DeliveryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.trade.Deliver(c.Request.Context(), middleware.MerID(c), uint(id), in); err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeIDErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, identity.ErrNoVerifyPerm),
		errors.Is(err, identity.ErrNoDeliverPerm),
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
	case errors.Is(err, trade.ErrBadStatus), errors.Is(err, trade.ErrNotPaid),
		errors.Is(err, trade.ErrBadParam), errors.Is(err, trade.ErrVerifyCodeMismatch),
		errors.Is(err, trade.ErrDeliveryParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
