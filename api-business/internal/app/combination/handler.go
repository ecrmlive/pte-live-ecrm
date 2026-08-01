package combination

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/combination"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/trade"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/response"
)

type Handler struct {
	svc   *combination.Service
	trade *trade.Service
}

func NewHandler(svc *combination.Service, tradeSvc *trade.Service) *Handler {
	return &Handler{svc: svc, trade: tradeSvc}
}

func (h *Handler) RegisterPublic(r gin.IRoutes) {
	r.GET("/combination/groups", h.List)
	r.GET("/combination/groups/:id", h.Get)
	r.GET("/combination/groups/:id/buyings", h.Buyings)
	r.GET("/combination/buyings/:id", h.BuyingDetail)
}

func (h *Handler) RegisterAuthed(r gin.IRoutes) {
	r.POST("/order/group/check", h.Check)
	r.POST("/order/group/create", h.Create)
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
	g, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) Buyings(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	rows, err := h.svc.ListBuyings(c.Request.Context(), uint(id), 20)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) BuyingDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	b, err := h.svc.GetBuying(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, b)
}

func (h *Handler) Check(c *gin.Context) {
	var in trade.GroupInput
	if err := c.ShouldBindJSON(&in); err != nil || in.ProductGroupID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.trade.GroupCheck(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	var in trade.GroupInput
	if err := c.ShouldBindJSON(&in); err != nil || in.ProductGroupID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.trade.GroupCreate(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, g)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, combination.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, combination.ErrBadParam),
		errors.Is(err, combination.ErrInactive),
		errors.Is(err, combination.ErrBuyingFull),
		errors.Is(err, combination.ErrBuyingClosed),
		errors.Is(err, combination.ErrAlreadyJoined):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}

func writeTradeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, trade.ErrAddressRequired), errors.Is(err, trade.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, combination.ErrInactive),
		errors.Is(err, combination.ErrBuyingFull),
		errors.Is(err, combination.ErrBuyingClosed),
		errors.Is(err, combination.ErrAlreadyJoined),
		errors.Is(err, combination.ErrNotFound),
		errors.Is(err, combination.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusBadRequest, err.Error())
	}
}
