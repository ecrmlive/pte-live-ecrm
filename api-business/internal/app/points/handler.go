package points

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/cart"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/catalog"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
)

type Handler struct {
	trade   *trade.Service
	catalog *catalog.Service
}

func NewHandler(tradeSvc *trade.Service, catSvc *catalog.Service) *Handler {
	return &Handler{trade: tradeSvc, catalog: catSvc}
}

func (h *Handler) RegisterPublic(r gin.IRoutes) {
	r.GET("/points/products", h.ListProducts)
	r.GET("/points/products/:id", h.ProductDetail)
}

func (h *Handler) RegisterAuthed(r gin.IRoutes) {
	r.GET("/integral", h.Integral)
	r.POST("/order/v3/check", h.V3Check)
	r.POST("/order/v3/create", h.V3Create)
	r.POST("/order/points/pay/:id", h.PointsPay)
}

func (h *Handler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.catalog.ListPointsProducts(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	list := make([]gin.H, 0, len(res.List))
	for _, p := range res.List {
		list = append(list, toPointsProduct(p))
	}
	response.OK(c, gin.H{"list": list, "total": res.Total, "page": res.Page, "limit": res.Limit})
}

func (h *Handler) ProductDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.catalog.GetAppProduct(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, catalog.ErrNotFound) || errors.Is(err, catalog.ErrNotOnSale) {
			response.Fail(c, http.StatusNotFound, "商品不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	if p.ProductType != catalog.ProductTypePoints {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	item := toPointsProduct(*p)
	item["unit_name"] = p.UnitName
	item["store_info"] = p.StoreInfo
	item["slider_image"] = splitImages(p.SliderImage, p.Image)
	response.OK(c, item)
}

func (h *Handler) Integral(c *gin.Context) {
	n, err := h.trade.UserIntegral(c.Request.Context(), middleware.UID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"integral": n})
}

func (h *Handler) V3Check(c *gin.Context) {
	var in trade.PointsInput
	if err := c.ShouldBindJSON(&in); err != nil || in.ProductID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.trade.V3Check(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) V3Create(c *gin.Context) {
	var in trade.PointsInput
	if err := c.ShouldBindJSON(&in); err != nil || in.ProductID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.trade.V3Create(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) PointsPay(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	g, err := h.trade.PointsPay(c.Request.Context(), middleware.UID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func toPointsProduct(p catalog.Product) gin.H {
	integral := p.Integral
	if integral <= 0 {
		integral = int(p.OtPrice)
	}
	return gin.H{
		"id":           p.ProductID,
		"mer_id":       p.MerID,
		"mer_name":     p.MerName,
		"store_name":   p.StoreName,
		"image":        p.Image,
		"price":        fmt.Sprintf("%.2f", p.Price),
		"ot_price":     fmt.Sprintf("%.2f", p.OtPrice),
		"integral":     integral,
		"sales":        p.Sales,
		"stock":        p.Stock,
		"product_type": p.ProductType,
	}
}

func splitImages(slider, cover string) []string {
	parts := strings.Split(slider, ",")
	out := make([]string, 0, len(parts)+1)
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 && cover != "" {
		out = append(out, cover)
	}
	return out
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, trade.ErrNotFound),
		errors.Is(err, cart.ErrNotFound),
		errors.Is(err, cart.ErrAddrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, trade.ErrForbidden), errors.Is(err, cart.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, trade.ErrIntegralNotEnough),
		errors.Is(err, trade.ErrMerIntegralOff),
		errors.Is(err, trade.ErrStockNotEnough),
		errors.Is(err, trade.ErrAddressRequired),
		errors.Is(err, trade.ErrBadParam),
		errors.Is(err, trade.ErrNotPointsProduct),
		errors.Is(err, trade.ErrPointsAlone),
		errors.Is(err, trade.ErrAlreadyPaid),
		errors.Is(err, cart.ErrAddrInvalid):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
