package coupon

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/promotion"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct {
	svc *promotion.Service
	id  *identity.Service
}

func NewHandler(svc *promotion.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/coupons", h.List)
	r.GET("/coupons/records", h.ListRecords)
	r.GET("/coupons/sends", h.ListSends)
	r.POST("/coupons", middleware.RequireMerchantMenu(h.id, identity.MerPermCouponCreate), h.Create)
	r.POST("/coupons/:id/send", middleware.RequireMerchantMenu(h.id, identity.MerPermCouponSend), h.Send)
	r.PUT("/coupons/:id", h.Update)
	r.DELETE("/coupons/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermCouponDelete), h.Delete)
	r.POST("/coupons/:id/status", middleware.RequireMerchantMenu(h.id, identity.MerPermCouponToggle), h.Status)
}

func (h *Handler) ListRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var couponID *uint
	if raw := c.Query("coupon_id"); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err == nil && id > 0 {
			v := uint(id)
			couponID = &v
		}
	}
	res, err := h.svc.ListMerchantCouponUsers(c.Request.Context(), middleware.MerID(c), couponID, page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) ListSends(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListMerchantCouponSends(c.Request.Context(), middleware.MerID(c), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListAdmin(c.Request.Context(), middleware.MerID(c), promotion.CouponTypeStore, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	var in promotion.CouponSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateAdmin(c.Request.Context(), middleware.MerID(c), promotion.CouponTypeStore, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Send(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in promotion.CouponSendInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SendMerchantCoupon(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in promotion.CouponSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateAdmin(c.Request.Context(), middleware.MerID(c), promotion.CouponTypeStore, uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Status(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in promotion.StatusInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetStatus(c.Request.Context(), middleware.MerID(c), promotion.CouponTypeStore, uint(id), in.Status); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteAdmin(c.Request.Context(), middleware.MerID(c), promotion.CouponTypeStore, uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, promotion.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, promotion.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, promotion.ErrBadParam),
		errors.Is(err, promotion.ErrBadStatus),
		errors.Is(err, promotion.ErrClosed),
		errors.Is(err, promotion.ErrSoldOut),
		errors.Is(err, promotion.ErrAlreadyReceived):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
