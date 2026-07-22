package diy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/diy"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct{ svc *diy.Service }

func NewHandler(svc *diy.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/diy/home", h.Home)
}

func (h *Handler) Home(c *gin.Context) {
	p, err := h.svc.GetActiveHome(c.Request.Context(), 0)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	if p == nil {
		response.OK(c, gin.H{
			"id": 0, "title": "", "name": "",
			"banners": []any{}, "menus": []any{},
		})
		return
	}
	v := diy.PageValue{}
	if p.Parsed != nil {
		v = *p.Parsed
	} else {
		v = p.ParseValue()
	}
	response.OK(c, gin.H{
		"id": p.ID, "name": p.Name, "title": p.Title,
		"template_name": p.TemplateName,
		"banners":       v.Banners,
		"menus":         v.Menus,
		"value":         p.Value,
	})
}
