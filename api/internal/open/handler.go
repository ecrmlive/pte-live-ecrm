package open

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/catalog"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/openapi"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

const ctxOpenAuth = "open_auth"

type Handler struct {
	openSvc *openapi.Service
	catSvc  *catalog.Service
	trade   *trade.Service
}

func NewHandler(openSvc *openapi.Service, catSvc *catalog.Service, tradeSvc *trade.Service) *Handler {
	return &Handler{openSvc: openSvc, catSvc: catSvc, trade: tradeSvc}
}

func (h *Handler) Register(public, authed gin.IRoutes) {
	public.POST("/auth", h.Auth)
	authed.GET("/order/list", h.OrderList)
	authed.GET("/order/detail/:id", h.OrderDetail)
	authed.GET("/product/list", h.ProductList)
	authed.GET("/product/detail/:id", h.ProductDetail)
	authed.POST("/product/create", h.ProductCreate)
}

// RequireAuth 校验开放 JWT，并加载凭证权限（挂到 context）。
func (h *Handler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := middleware.ClaimsFrom(c)
		if claims == nil || claims.Account == "" {
			response.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}
		row, err := h.openSvc.LoadAuth(c.Request.Context(), claims.Account)
		if err != nil {
			writeOpenErr(c, err)
			c.Abort()
			return
		}
		if row.MerID != claims.MerID {
			response.Fail(c, http.StatusForbidden, "商户上下文不匹配")
			c.Abort()
			return
		}
		c.Set(ctxOpenAuth, row)
		c.Next()
	}
}

func openAuthFrom(c *gin.Context) *openapi.OpenAuth {
	v, ok := c.Get(ctxOpenAuth)
	if !ok {
		return nil
	}
	a, _ := v.(*openapi.OpenAuth)
	return a
}

func (h *Handler) Auth(c *gin.Context) {
	var in openapi.AuthInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.openSvc.Authenticate(c.Request.Context(), in, c.ClientIP())
	if err != nil {
		writeOpenErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) OrderList(c *gin.Context) {
	a := openAuthFrom(c)
	if a == nil || !a.AllowOrder() {
		response.Fail(c, http.StatusForbidden, "无权访问订单")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var paid, status *int8
	if s := c.Query("paid"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			vv := int8(v)
			paid = &vv
		}
	}
	if s := c.Query("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			vv := int8(v)
			status = &vv
		}
	}
	res, err := h.trade.MerchantList(c.Request.Context(), middleware.MerID(c), paid, status, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) OrderDetail(c *gin.Context) {
	a := openAuthFrom(c)
	if a == nil || !a.AllowOrder() {
		response.Fail(c, http.StatusForbidden, "无权访问订单")
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	o, err := h.trade.GetMerchantOrder(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeTradeErr(c, err)
		return
	}
	response.OK(c, o)
}

func (h *Handler) ProductList(c *gin.Context) {
	a := openAuthFrom(c)
	if a == nil || !a.AllowProduct() {
		response.Fail(c, http.StatusForbidden, "无权访问商品")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var statusPtr *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		statusPtr = &st
	}
	res, err := h.catSvc.ListMerchantProducts(c.Request.Context(), middleware.MerID(c), statusPtr, c.Query("keyword"), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) ProductDetail(c *gin.Context) {
	a := openAuthFrom(c)
	if a == nil || !a.AllowProduct() {
		response.Fail(c, http.StatusForbidden, "无权访问商品")
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.catSvc.GetMerchantProduct(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeCatErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) ProductCreate(c *gin.Context) {
	a := openAuthFrom(c)
	if a == nil || !a.AllowProduct() {
		response.Fail(c, http.StatusForbidden, "无权访问商品")
		return
	}
	var in catalog.ProductSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	p, err := h.catSvc.CreateMerchantProduct(c.Request.Context(), middleware.MerID(c), in)
	if err != nil {
		writeCatErr(c, err)
		return
	}
	response.OK(c, p)
}

func writeOpenErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, openapi.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, openapi.ErrUnauthorized), errors.Is(err, openapi.ErrDisabled):
		response.Fail(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, openapi.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}

func writeTradeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, trade.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, trade.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}

func writeCatErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, catalog.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, catalog.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, catalog.ErrInvalidPrice),
		errors.Is(err, catalog.ErrNameRequired),
		errors.Is(err, catalog.ErrCateRequired):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
