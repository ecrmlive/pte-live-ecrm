package coupon

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/cart"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/promotion"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc     *promotion.Service
	cartSvc *cart.Service
}

func NewHandler(svc *promotion.Service, cartSvc *cart.Service) *Handler {
	return &Handler{svc: svc, cartSvc: cartSvc}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/coupons/center", h.Center)
	r.POST("/coupons/:id/receive", h.Receive)
	r.GET("/coupons/mine", h.Mine)
	r.GET("/coupons/usable", h.Usable)
	r.POST("/spread/bind", h.BindSpread)
	r.GET("/spread/me", h.SpreadMe)
	r.GET("/spread/bills", h.Bills)
}

func (h *Handler) Center(c *gin.Context) {
	uid := middleware.UID(c)
	var merID *uint
	if s := c.Query("mer_id"); s != "" {
		v, _ := strconv.ParseUint(s, 10, 64)
		m := uint(v)
		merID = &m
	}
	list, err := h.svc.ListReceivable(c.Request.Context(), uid, merID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	// 未指定 mer_id 时附带演示店铺券（mer=1）
	if merID == nil {
		m := uint(1)
		store, err := h.svc.ListReceivable(c.Request.Context(), uid, &m)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询失败")
			return
		}
		list = append(list, store...)
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) Receive(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	cu, err := h.svc.Receive(c.Request.Context(), middleware.UID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, cu)
}

func (h *Handler) Mine(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var status *int
	if s := c.Query("status"); s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			status = &v
		}
	}
	res, err := h.svc.MyCoupons(c.Request.Context(), middleware.UID(c), status, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Usable(c *gin.Context) {
	raw := strings.TrimSpace(c.Query("cart_ids"))
	if raw == "" {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	parts := strings.Split(raw, ",")
	ids := make([]uint64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseUint(p, 10, 64)
		if err != nil || v == 0 {
			response.Fail(c, http.StatusBadRequest, "参数错误")
			return
		}
		ids = append(ids, v)
	}
	if len(ids) == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	rows, err := h.cartSvc.LoadForCheckout(c.Request.Context(), middleware.UID(c), ids)
	if err != nil {
		writeCartErr(c, err)
		return
	}
	mers := map[uint]float64{}
	var total float64
	for _, r := range rows {
		sub := r.Price * float64(r.CartNum)
		mers[r.MerID] += sub
		total += sub
	}
	zero := 0
	mine, err := h.svc.MyCoupons(c.Request.Context(), middleware.UID(c), &zero, 1, 50)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	usable := make([]promotion.CouponUser, 0)
	for _, cu := range mine.List {
		ok := false
		if cu.CouponKind == promotion.CouponTypePlatform {
			ok = total >= float64(cu.UseMinPrice)
		} else {
			ok = mers[cu.MerID] >= float64(cu.UseMinPrice)
		}
		if ok {
			usable = append(usable, cu)
		}
	}
	response.OK(c, gin.H{"list": usable, "total_price": total})
}

func (h *Handler) BindSpread(c *gin.Context) {
	var in promotion.BindSpreadInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.BindSpread(c.Request.Context(), middleware.UID(c), in.SpreadUID); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) SpreadMe(c *gin.Context) {
	spreadUID, isPromoter, count, err := h.svc.SpreadMe(c.Request.Context(), middleware.UID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{
		"uid":          middleware.UID(c),
		"spread_uid":   spreadUID,
		"is_promoter":  isPromoter,
		"spread_count": count,
	})
}

func (h *Handler) Bills(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.MyBrokerage(c.Request.Context(), middleware.UID(c), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, promotion.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, promotion.ErrAlreadyReceived),
		errors.Is(err, promotion.ErrSoldOut),
		errors.Is(err, promotion.ErrClosed),
		errors.Is(err, promotion.ErrCouponInvalid),
		errors.Is(err, promotion.ErrCouponMinNotMet),
		errors.Is(err, promotion.ErrCouponConflict),
		errors.Is(err, promotion.ErrSpreadSelf),
		errors.Is(err, promotion.ErrSpreadInvalid),
		errors.Is(err, promotion.ErrSpreadBound),
		errors.Is(err, promotion.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}

func writeCartErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, cart.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, cart.ErrForbidden),
		errors.Is(err, cart.ErrProductOff),
		errors.Is(err, cart.ErrStockNotEnough):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
