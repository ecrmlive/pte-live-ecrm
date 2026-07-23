package catalog

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/catalog"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/diy"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *catalog.Service
	diy *diy.Service
}

func NewHandler(svc *catalog.Service, diySvc *diy.Service) *Handler {
	return &Handler{svc: svc, diy: diySvc}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/catalog/home", h.Home)
	r.GET("/catalog/categories", h.Categories)
	r.GET("/catalog/products", h.Products)
	r.GET("/catalog/points/products", h.PointsProducts)
	r.GET("/catalog/products/:id", h.ProductDetail)
	r.GET("/catalog/stores/:id", h.StoreHome)
}

func (h *Handler) Home(c *gin.Context) {
	_, hot, err := h.svc.AppHome(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	items := make([]gin.H, 0, len(hot))
	for _, p := range hot {
		items = append(items, toAppProduct(p))
	}

	banners := []gin.H{}
	menus := []gin.H{}
	diyID := uint(0)
	diyTitle := ""
	if h.diy != nil {
		page, err := h.diy.GetActiveHome(c.Request.Context(), 0)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询失败")
			return
		}
		if page != nil {
			diyID = page.ID
			diyTitle = page.Title
			doc := page.ParseDoc()
			for _, item := range doc.Items {
				t, _ := item["type"].(string)
				rawData, _ := item["data"].([]any)
				switch t {
				case "banner":
					for i, row := range rawData {
						m, _ := row.(map[string]any)
						if m == nil {
							continue
						}
						banners = append(banners, gin.H{
							"id": i + 1, "title": m["imgName"], "image": m["imgUrl"], "url": m["linkUrl"],
						})
					}
				case "navBar", "option":
					for i, row := range rawData {
						m, _ := row.(map[string]any)
						if m == nil {
							continue
						}
						text := m["text"]
						if text == nil {
							text = m["title"]
						}
						menus = append(menus, gin.H{
							"id": i + 1, "name": text, "icon": m["imgUrl"], "url": m["linkUrl"],
						})
					}
				}
			}
		}
	}
	if len(banners) == 0 {
		banners = []gin.H{
			{"id": 1, "title": "多商户入驻 · 精选好物", "image": "", "url": ""},
			{"id": 2, "title": "同城配送 · 品质生活", "image": "", "url": ""},
		}
	}

	response.OK(c, gin.H{
		"diy_id":    diyID,
		"diy_title": diyTitle,
		"banners":   banners,
		"menus":     menus,
		"hot":       items,
	})
}

func (h *Handler) Categories(c *gin.Context) {
	rows, err := h.svc.ListAppCategories(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	out := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		out = append(out, gin.H{
			"id":   r.StoreCategoryID,
			"name": r.CateName,
			"pid":  r.PID,
			"pic":  r.Pic,
		})
	}
	response.OK(c, out)
}

func (h *Handler) Products(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cateID, _ := strconv.Atoi(c.Query("cate_id"))
	var merPtr *uint
	if m := c.Query("mer_id"); m != "" {
		v, _ := strconv.ParseUint(m, 10, 64)
		id := uint(v)
		merPtr = &id
	}
	res, err := h.svc.ListAppProducts(c.Request.Context(), cateID, c.Query("keyword"), merPtr, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	list := make([]gin.H, 0, len(res.List))
	for _, p := range res.List {
		list = append(list, toAppProduct(p))
	}
	response.OK(c, gin.H{"list": list, "total": res.Total, "page": res.Page, "limit": res.Limit})
}

func (h *Handler) PointsProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListPointsProducts(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	list := make([]gin.H, 0, len(res.List))
	for _, p := range res.List {
		item := toAppProduct(p)
		item["integral"] = p.Integral
		item["product_type"] = p.ProductType
		list = append(list, item)
	}
	response.OK(c, gin.H{"list": list, "total": res.Total, "page": res.Page, "limit": res.Limit})
}

func (h *Handler) ProductDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.GetAppProduct(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, catalog.ErrNotFound) || errors.Is(err, catalog.ErrNotOnSale) {
			response.Fail(c, http.StatusNotFound, "商品不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	item := toAppProduct(*p)
	item["unit_name"] = p.UnitName
	item["store_info"] = p.StoreInfo
	item["slider_image"] = splitImages(p.SliderImage, p.Image)
	item["spec_type"] = p.SpecType
	item["delivery_way"] = p.DeliveryWay
	response.OK(c, item)
}

func (h *Handler) StoreHome(c *gin.Context) {
	merID64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	merID := uint(merID64)
	if merID == 0 {
		response.Fail(c, http.StatusBadRequest, "商户无效")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	res, err := h.svc.ListAppProducts(c.Request.Context(), 0, "", &merID, page, 20)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	list := make([]gin.H, 0, len(res.List))
	merName := ""
	for _, p := range res.List {
		if merName == "" {
			merName = p.MerName
		}
		list = append(list, toAppProduct(p))
	}
	response.OK(c, gin.H{
		"mer_id":   merID,
		"mer_name": merName,
		"products": list,
		"total":    res.Total,
	})
}

func toAppProduct(p catalog.Product) gin.H {
	return gin.H{
		"id":              p.ProductID,
		"mer_id":          p.MerID,
		"mer_name":        p.MerName,
		"store_name":      p.StoreName,
		"image":           p.Image,
		"price":           fmt.Sprintf("%.2f", p.Price),
		"ot_price":        fmt.Sprintf("%.2f", p.OtPrice),
		"svip_price_type": p.SvipPriceType,
		"svip_price":      fmt.Sprintf("%.2f", p.SvipPrice),
		"sales":           p.Sales,
		"stock":           p.Stock,
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
