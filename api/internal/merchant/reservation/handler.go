package reservation

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/reservation"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *reservation.Service
	id  *identity.Service
}

func NewHandler(svc *reservation.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/reservation/products", h.List)
	r.GET("/reservation/products/:id/config", h.GetConfig)
	r.PUT("/reservation/products/:id/config", middleware.RequireMerchantMenu(h.id, identity.MerPermReservationConfig), h.SaveConfig)
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

func (h *Handler) GetConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	cfg, slots, err := h.svc.GetConfig(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": cfg, "slots": slots})
}

func (h *Handler) SaveConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in reservation.ConfigSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	in.ProductID = uint(id)
	if err := h.svc.SaveConfig(c.Request.Context(), middleware.MerID(c), in); err != nil {
		writeErr(c, err)
		return
	}
	cfg, slots, err := h.svc.GetConfig(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"config": cfg, "slots": slots})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, reservation.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, reservation.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, reservation.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
