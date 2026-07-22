package reservation

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/reservation"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc   *reservation.Service
	trade *trade.Service
}

func NewHandler(svc *reservation.Service, tradeSvc *trade.Service) *Handler {
	return &Handler{svc: svc, trade: tradeSvc}
}

func (h *Handler) RegisterPublic(r gin.IRoutes) {
	r.GET("/reservation/products", h.List)
	r.GET("/reservation/products/:id", h.Get)
	r.GET("/reservation/products/:id/slots", h.Slots)
}

func (h *Handler) RegisterAuthed(r gin.IRoutes) {
	r.POST("/order/reservation/check", h.Check)
	r.POST("/order/reservation/create", h.Create)
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
	p, err := h.svc.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Slots(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	date := c.Query("date")
	rows, err := h.svc.DaySlots(c.Request.Context(), uint(id), date)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"list": rows, "date": date})
}

func (h *Handler) Check(c *gin.Context) {
	var in trade.ReservationInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.trade.ReservationCheck(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	var in trade.ReservationInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.trade.ReservationCreate(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, g)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, reservation.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, reservation.ErrBadParam), errors.Is(err, reservation.ErrBadDate),
		errors.Is(err, reservation.ErrNoSlot), errors.Is(err, reservation.ErrFull):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, reservation.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}

func writeTradeErr(c *gin.Context, err error) {
	if errors.Is(err, trade.ErrAddressRequired) || errors.Is(err, trade.ErrBadParam) {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	writeErr(c, err)
}
