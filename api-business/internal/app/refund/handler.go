package refund

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/aftersale"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/response"
)

type Handler struct {
	svc *aftersale.Service
}

func NewHandler(svc *aftersale.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/refund/apply", h.Apply)
	r.GET("/refunds", h.List)
	r.GET("/refunds/:id", h.Get)
	r.POST("/refunds/:id/cancel", h.Cancel)
	r.POST("/refunds/:id/platform", h.RequestPlatform)
}

func (h *Handler) Apply(c *gin.Context) {
	var in aftersale.ApplyInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	ro, err := h.svc.Apply(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, ro)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListUser(c.Request.Context(), middleware.UID(c), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ro, err := h.svc.GetUser(c.Request.Context(), middleware.UID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, ro)
}

func (h *Handler) Cancel(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Cancel(c.Request.Context(), middleware.UID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) RequestPlatform(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.RequestPlatform(c.Request.Context(), middleware.UID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, aftersale.ErrNotFound), errors.Is(err, aftersale.ErrOrderNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, aftersale.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, aftersale.ErrBadParam),
		errors.Is(err, aftersale.ErrBadStatus),
		errors.Is(err, aftersale.ErrOrderNotPaid),
		errors.Is(err, aftersale.ErrOrderRefunded),
		errors.Is(err, aftersale.ErrRefundInProgress),
		errors.Is(err, aftersale.ErrProductInvalid),
		errors.Is(err, aftersale.ErrAlreadyDone),
		errors.Is(err, aftersale.ErrRejectNeedMsg):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
