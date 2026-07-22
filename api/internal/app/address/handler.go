package address

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
	r.GET("/address", h.List)
	r.POST("/address", h.Create)
	r.PUT("/address/:id", h.Update)
	r.DELETE("/address/:id", h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	rows, err := h.svc.ListAddresses(c.Request.Context(), middleware.UID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) Create(c *gin.Context) {
	var in cart.AddressInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	a, err := h.svc.CreateAddress(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, a)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in cart.AddressInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	a, err := h.svc.UpdateAddress(c.Request.Context(), middleware.UID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, a)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteAddress(c.Request.Context(), middleware.UID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, cart.ErrAddrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, cart.ErrAddrInvalid):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
