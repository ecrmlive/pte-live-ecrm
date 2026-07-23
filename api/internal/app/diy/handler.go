package diy

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/diy"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct{ svc *diy.Service }

func NewHandler(svc *diy.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/diy/home", h.Home)
	r.GET("/diy/pages/:id", h.Get)
}

func (h *Handler) Home(c *gin.Context) {
	merID := uint(0)
	if v := c.Query("mer_id"); v != "" {
		n, _ := strconv.ParseUint(v, 10, 64)
		merID = uint(n)
	}
	p, err := h.svc.GetActiveHome(c.Request.Context(), merID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	if p == nil {
		response.OK(c, gin.H{
			"id": 0, "title": "", "name": "",
			"page": map[string]any{}, "items": []any{},
			"banners": []any{}, "menus": []any{},
		})
		return
	}
	response.OK(c, pagePayload(p))
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		response.Fail(c, http.StatusNotFound, "装修页不存在")
		return
	}
	response.OK(c, pagePayload(p))
}

func pagePayload(p *diy.Page) gin.H {
	doc := diy.PageDoc{}
	if p.Doc != nil {
		doc = *p.Doc
	} else {
		doc = p.ParseDoc()
	}
	banners, menus := extractLegacyCompat(doc)
	return gin.H{
		"id":            p.ID,
		"name":          p.Name,
		"title":         p.Title,
		"template_name": p.TemplateName,
		"cover_image":   p.CoverImage,
		"is_bg_color":   p.IsBgColor,
		"is_bg_pic":     p.IsBgPic,
		"color_picker":  p.ColorPicker,
		"bg_pic":        p.BgPic,
		"page":          doc.Page,
		"items":         doc.Items,
		"banners":       banners,
		"menus":         menus,
		"value":         p.Value,
	}
}

func extractLegacyCompat(doc diy.PageDoc) ([]map[string]any, []map[string]any) {
	banners := []map[string]any{}
	menus := []map[string]any{}
	for _, item := range doc.Items {
		t, _ := item["type"].(string)
		data, _ := item["data"].([]any)
		if data == nil {
			if arr, ok := item["data"].([]map[string]any); ok {
				for _, row := range arr {
					data = append(data, row)
				}
			}
		}
		switch t {
		case "banner":
			for i, row := range data {
				m, _ := row.(map[string]any)
				if m == nil {
					continue
				}
				banners = append(banners, map[string]any{
					"id": i + 1, "title": m["imgName"], "image": m["imgUrl"], "url": m["linkUrl"],
				})
			}
		case "navBar", "option":
			for i, row := range data {
				m, _ := row.(map[string]any)
				if m == nil {
					continue
				}
				text := m["text"]
				if text == nil {
					text = m["title"]
				}
				menus = append(menus, map[string]any{
					"id": i + 1, "name": text, "icon": m["imgUrl"], "url": m["linkUrl"],
				})
			}
		}
	}
	return banners, menus
}
