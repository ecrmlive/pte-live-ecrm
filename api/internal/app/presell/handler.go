package presell

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/presell"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc   *presell.Service
	trade *trade.Service
}

func NewHandler(svc *presell.Service, tradeSvc *trade.Service) *Handler {
	return &Handler{svc: svc, trade: tradeSvc}
}

func (h *Handler) RegisterPublic(r gin.IRoutes) {
	r.GET("/presell/actives", h.List)
	r.GET("/presell/actives/:id", h.Get)
}

func (h *Handler) RegisterAuthed(r gin.IRoutes) {
	r.POST("/order/presell/check", h.Check)
	r.POST("/order/presell/create", h.Create)
	r.GET("/presell/finals", h.ListFinals)
	r.GET("/presell/finals/:id", h.GetFinal)
	r.POST("/presell/pay/:id", h.FinalPay)
	r.POST("/order/presell/final/pay/:id", h.FinalPay)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListApp(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Check(c *gin.Context) {
	var in trade.PresellInput
	if err := c.ShouldBindJSON(&in); err != nil || in.ProductPresellID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.trade.PresellCheck(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	var in trade.PresellInput
	if err := c.ShouldBindJSON(&in); err != nil || in.ProductPresellID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.trade.PresellCreate(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) ListFinals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.trade.ListPresellFinals(c.Request.Context(), middleware.UID(c), page, limit)
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetFinal(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.trade.GetPresellFinal(c.Request.Context(), middleware.UID(c), uint(id))
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) FinalPay(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		Type    string `json:"type"`
		PayType string `json:"pay_type"`
	}
	_ = c.ShouldBindJSON(&body)
	payType := body.Type
	if payType == "" {
		payType = body.PayType
	}
	if payType == "" {
		payType = "mock"
	}
	row, err := h.trade.PresellFinalPay(c.Request.Context(), middleware.UID(c), uint(id), payType)
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, row)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, presell.ErrNotFound), errors.Is(err, presell.ErrFinalNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, presell.ErrInactive), errors.Is(err, presell.ErrSoldOut),
		errors.Is(err, presell.ErrNotFullPay), errors.Is(err, presell.ErrBadParam),
		errors.Is(err, presell.ErrFinalNotOpen), errors.Is(err, presell.ErrFinalPaid),
		errors.Is(err, presell.ErrFinalInvalid):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}

func writeTradeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, trade.ErrAddressRequired), errors.Is(err, trade.ErrBadParam),
		errors.Is(err, trade.ErrStockNotEnough), errors.Is(err, trade.ErrAlreadyPaid),
		errors.Is(err, trade.ErrBadStatus), errors.Is(err, trade.ErrBalanceNotEnough),
		errors.Is(err, trade.ErrInvalidPayType),
		errors.Is(err, trade.ErrPresellFinalNotOpen), errors.Is(err, trade.ErrPresellFinalPaid),
		errors.Is(err, trade.ErrPresellFinalInvalid), errors.Is(err, trade.ErrPresellFinalTimeout),
		errors.Is(err, trade.ErrPresellFinal),
		errors.Is(err, presell.ErrInactive), errors.Is(err, presell.ErrSoldOut),
		errors.Is(err, presell.ErrNotFullPay), errors.Is(err, presell.ErrBadParam),
		errors.Is(err, presell.ErrFinalNotOpen), errors.Is(err, presell.ErrFinalPaid),
		errors.Is(err, presell.ErrFinalInvalid):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, trade.ErrNotFound), errors.Is(err, trade.ErrForbidden),
		errors.Is(err, trade.ErrPresellFinalNotFound),
		errors.Is(err, presell.ErrNotFound), errors.Is(err, presell.ErrFinalNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
