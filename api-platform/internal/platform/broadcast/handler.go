package broadcast

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/broadcast"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/response"
)

type Handler struct {
	svc *broadcast.Service
	id  *identity.Service
}

func NewHandler(svc *broadcast.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/broadcast/rooms", h.List)
	r.GET("/broadcast/rooms/:id", h.Get)
	r.PUT("/broadcast/rooms/:id/status", middleware.RequirePlatformMenu(h.id, identity.PlatPermBroadcastAudit), h.UpdateStatus)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListPlatform(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	room, err := h.svc.Get(c.Request.Context(), uint(id), true)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, room)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in broadcast.AuditInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	room, err := h.svc.Audit(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, room)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, broadcast.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, broadcast.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
