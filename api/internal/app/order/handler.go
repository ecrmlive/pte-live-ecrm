package order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/cart"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/promotion"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *trade.Service
}

func NewHandler(svc *trade.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/v2/order/check", h.V2Check)
	r.POST("/v2/order/create", h.V2Create)
	r.POST("/order/pay/:id", h.Pay)
	r.GET("/orders", h.ListGroup)
	r.GET("/orders/:id", h.GetGroup)
	r.GET("/order/:id", h.GetOrder)
}

func (h *Handler) V2Check(c *gin.Context) {
	var in trade.CheckInput
	if err := c.ShouldBindJSON(&in); err != nil || len(in.CartIDs) == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.svc.V2Check(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) V2Create(c *gin.Context) {
	var in trade.CreateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.svc.V2Create(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) Pay(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		PayType string `json:"pay_type"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.PayType == "" {
		req.PayType = "mock"
	}
	uid := middleware.UID(c)
	if trade.IsChannelPayType(req.PayType) {
		intent, err := h.svc.CreatePayIntent(c.Request.Context(), uid, uint(id), req.PayType)
		if err != nil {
			writeErr(c, err)
			return
		}
		response.OK(c, intent)
		return
	}
	g, err := h.svc.PaySuccess(c.Request.Context(), uid, uint(id), req.PayType)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) ListGroup(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListGroupOrders(c.Request.Context(), middleware.UID(c), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetGroup(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	g, err := h.svc.GetGroupOrder(c.Request.Context(), middleware.UID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) GetOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	o, err := h.svc.GetStoreOrderForUser(c.Request.Context(), middleware.UID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, o)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, trade.ErrNotFound),
		errors.Is(err, cart.ErrNotFound),
		errors.Is(err, cart.ErrAddrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, trade.ErrForbidden), errors.Is(err, cart.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, trade.ErrInvalidPayType),
		errors.Is(err, trade.ErrChannelDisabled),
		errors.Is(err, trade.ErrBadNotify),
		errors.Is(err, trade.ErrBalanceNotEnough),
		errors.Is(err, trade.ErrIntegralNotEnough),
		errors.Is(err, trade.ErrMerIntegralOff),
		errors.Is(err, trade.ErrStockNotEnough),
		errors.Is(err, trade.ErrEmptyCart),
		errors.Is(err, trade.ErrBadStatus),
		errors.Is(err, trade.ErrAlreadyPaid),
		errors.Is(err, trade.ErrAddressRequired),
		errors.Is(err, trade.ErrCoupon),
		errors.Is(err, trade.ErrBadParam),
		errors.Is(err, trade.ErrPointsAlone),
		errors.Is(err, trade.ErrNotPointsProduct),
		errors.Is(err, trade.ErrPointsProductMix),
		errors.Is(err, trade.ErrSeckillLimit),
		errors.Is(err, trade.ErrIntegralOnActivity),
		errors.Is(err, promotion.ErrCouponInvalid),
		errors.Is(err, promotion.ErrCouponMinNotMet),
		errors.Is(err, promotion.ErrCouponConflict),
		errors.Is(err, promotion.ErrPlatformOnly),
		errors.Is(err, cart.ErrProductOff),
		errors.Is(err, cart.ErrStockNotEnough),
		errors.Is(err, cart.ErrAddrInvalid):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
