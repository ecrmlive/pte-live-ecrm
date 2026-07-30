package callback

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
)

type PayHandler struct {
	svc *trade.Service
}

func NewPayHandler(svc *trade.Service) *PayHandler { return &PayHandler{svc: svc} }

func (h *PayHandler) Register(r gin.IRoutes) {
	r.POST("/pay/wechat", h.Wechat)
	r.POST("/pay/alipay", h.Alipay)
}

func (h *PayHandler) Mock(c *gin.Context) {
	var body struct {
		GroupOrderID uint `json:"group_order_id"`
		UID          uint `json:"uid"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.GroupOrderID == 0 || body.UID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.svc.PayGroup(c.Request.Context(), body.UID, body.GroupOrderID, "mock")
	if err != nil {
		writePayErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *PayHandler) Wechat(c *gin.Context) { h.channel(c, "wechat") }
func (h *PayHandler) Alipay(c *gin.Context) { h.channel(c, "alipay") }

func (h *PayHandler) channel(c *gin.Context, channel string) {
	var body trade.ChannelNotifyInput
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.svc.NotifyChannelPay(c.Request.Context(), channel, body)
	if err != nil {
		writePayErr(c, err)
		return
	}
	response.OK(c, g)
}

func writePayErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, trade.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, trade.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, trade.ErrBadNotify),
		errors.Is(err, trade.ErrInvalidPayType),
		errors.Is(err, trade.ErrChannelDisabled),
		errors.Is(err, trade.ErrPaymentConfig),
		errors.Is(err, trade.ErrBalanceNotEnough),
		errors.Is(err, trade.ErrBadStatus),
		errors.Is(err, trade.ErrStockNotEnough),
		errors.Is(err, trade.ErrNotPointsProduct):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusBadRequest, err.Error())
	}
}
