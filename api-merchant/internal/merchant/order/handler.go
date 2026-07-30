package order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/logistics"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/response"
)

type Handler struct {
	svc       *trade.Service
	id        *identity.Service
	logistics *logistics.Service
}

func NewHandler(svc *trade.Service, id *identity.Service, logisticsSvc *logistics.Service) *Handler {
	return &Handler{svc: svc, id: id, logistics: logisticsSvc}
}

func (h *Handler) Register(r gin.IRoutes) {
	h.RegisterVerify(r)
}

// RegisterVerify remains temporarily isolated while pickup verification is
// migrated to qixi_crm_b_order_verification. List/detail/delivery are native.
func (h *Handler) RegisterVerify(r gin.IRoutes) {
	r.POST("/orders/:id/verify", middleware.RequireMerchantMenu(h.id, identity.MerPermOrderVerify), h.Verify)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
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
	res, err := h.svc.MerchantList(c.Request.Context(), middleware.MerID(c), paid, status, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	o, err := h.svc.GetMerchantOrder(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, o)
}

func (h *Handler) Deliver(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in trade.DeliveryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if in.ExpressID > 0 && in.DeliveryName == "" && h.logistics != nil {
		if name, err := h.logistics.GetExpressName(c.Request.Context(), in.ExpressID); err == nil {
			in.DeliveryName = name
		}
	}
	if err := h.svc.Deliver(c.Request.Context(), middleware.MerID(c), uint(id), in); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

type verifyReq struct {
	VerifyCode string `json:"verify_code"`
}

func (h *Handler) Verify(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req verifyReq
	_ = c.ShouldBindJSON(&req)
	code := req.VerifyCode
	if code == "" {
		// 商户后台按单核销：未传码时使用订单自带码
		o, err := h.svc.GetMerchantOrder(c.Request.Context(), middleware.MerID(c), uint(id))
		if err != nil {
			writeErr(c, err)
			return
		}
		code = o.VerifyCode
	}
	if err := h.svc.Verify(c.Request.Context(), middleware.MerID(c), uint(id), code); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, trade.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, trade.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, trade.ErrBadStatus), errors.Is(err, trade.ErrDeliveryParam),
		errors.Is(err, trade.ErrNotPaid), errors.Is(err, trade.ErrBadParam),
		errors.Is(err, trade.ErrVerifyCodeMismatch):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
