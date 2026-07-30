package order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/response"
)

type Handler struct {
	svc *trade.Service
	id  *identity.Service
}

func NewHandler(svc *trade.Service, id *identity.Service) *Handler { return &Handler{svc: svc, id: id} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/orders", h.List)
	r.GET("/orders/:id", h.Get)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var paid *int8
	if s := c.Query("paid"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		p := int8(v)
		paid = &p
	}
	regionIDs, err := h.id.PlatformRegionScope(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "登录已失效")
		return
	}
	res, err := h.svc.PlatformListByRegions(c.Request.Context(), paid, regionIDs, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	regionIDs, err := h.id.PlatformRegionScope(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "登录已失效")
		return
	}
	o, err := h.svc.PlatformGetByRegions(c.Request.Context(), uint(id), regionIDs)
	if err != nil {
		if errors.Is(err, trade.ErrNotFound) {
			response.Fail(c, http.StatusNotFound, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, o)
}
