package finance

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/finance"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct {
	svc *finance.Service
}

func NewHandler(svc *finance.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/finance/balance", h.Balance)
	r.GET("/finance/withdraws", h.List)
	r.POST("/finance/withdraw", h.Apply)
	r.GET("/finance/withdraws/:id", h.Get)
}

func (h *Handler) Balance(c *gin.Context) {
	bal, err := h.svc.Balance(c.Request.Context(), middleware.MerID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, bal)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListMerchant(c.Request.Context(), middleware.MerID(c), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Apply(c *gin.Context) {
	var in finance.WithdrawInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	f, err := h.svc.ApplyWithdraw(c.Request.Context(), middleware.MerID(c), middleware.AdminID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, f)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	f, err := h.svc.GetMerchant(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, f)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, finance.ErrNotFound), errors.Is(err, finance.ErrMerchantNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, finance.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, finance.ErrBadParam),
		errors.Is(err, finance.ErrBadStatus),
		errors.Is(err, finance.ErrBalanceNotEnough),
		errors.Is(err, finance.ErrAlreadyDone),
		errors.Is(err, finance.ErrRejectNeedMsg):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
