package cart

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/cart"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *cart.Service
}

func NewHandler(svc *cart.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/cart", h.List)
	r.POST("/cart", h.Add)
	r.PUT("/cart/:id", h.Update)
	r.DELETE("/cart/:id", h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context(), middleware.UID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	var totalNum uint
	var totalPrice float64
	for _, b := range list {
		totalPrice += b.Subtotal
		for _, it := range b.Items {
			if it.IsFail == 0 {
				totalNum += it.CartNum
			}
		}
	}
	response.OK(c, gin.H{"list": list, "total_num": totalNum, "total_price": totalPrice})
}

func (h *Handler) Add(c *gin.Context) {
	var in cart.AddInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	item, err := h.svc.Add(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, item)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		CartNum *uint `json:"cart_num" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.CartNum == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetNum(c.Request.Context(), middleware.UID(c), id, *req.CartNum); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), middleware.UID(c), id); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, cart.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, cart.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, cart.ErrProductOff),
		errors.Is(err, cart.ErrStockNotEnough),
		errors.Is(err, cart.ErrInvalidNum):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
