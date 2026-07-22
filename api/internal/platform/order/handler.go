package order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *trade.Service
}

func NewHandler(svc *trade.Service) *Handler { return &Handler{svc: svc} }

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
	res, err := h.svc.PlatformList(c.Request.Context(), paid, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	o, err := h.svc.PlatformGet(c.Request.Context(), uint(id))
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
