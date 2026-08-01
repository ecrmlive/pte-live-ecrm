package circle

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/circle"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
)

type Handler struct {
	svc *circle.Service
	id  *identity.Service
}

func NewHandler(svc *circle.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/business-zones", h.ListCircles)
	r.POST("/business-zones", middleware.RequirePlatformMenu(h.id, identity.PlatPermCircleManage), h.CreateCircle)
	r.GET("/business-zones/:id", h.GetCircle)
	r.PUT("/business-zones/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermCircleManage), h.UpdateCircle)
	r.DELETE("/business-zones/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermCircleManage), h.DeleteCircle)
	r.GET("/business-zone-agents", h.ListAgents)
	r.GET("/business-zone-agents/:id", h.GetAgent)
	r.POST("/business-zone-agents", middleware.RequirePlatformMenu(h.id, identity.PlatPermCircleAgentManage), h.CreateAgent)
	r.PUT("/business-zone-agents/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermCircleAgentManage), h.UpdateAgent)
	r.POST("/business-zone-agents/:id/audit", middleware.RequirePlatformMenu(h.id, identity.PlatPermCircleAgentReview), h.AuditAgent)
}

func page(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	return p, l
}
func optionalStatus(c *gin.Context) (*int8, error) {
	raw := c.Query("status")
	if raw == "" {
		return nil, nil
	}
	n, err := strconv.ParseInt(raw, 10, 8)
	if err != nil {
		return nil, circle.ErrBadParam
	}
	v := int8(n)
	return &v, nil
}
func id(c *gin.Context) (uint, error) {
	n, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || n == 0 {
		return 0, circle.ErrBadParam
	}
	return uint(n), nil
}

func (h *Handler) ListCircles(c *gin.Context) {
	status, err := optionalStatus(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	p, l := page(c)
	out, err := h.svc.ListCircles(c.Request.Context(), c.Query("keyword"), status, p, l)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, out)
}
func (h *Handler) GetCircle(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	row, err := h.svc.GetCircle(c.Request.Context(), key)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}
func (h *Handler) CreateCircle(c *gin.Context) {
	var in circle.CircleInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	row, err := h.svc.CreateCircle(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}
func (h *Handler) UpdateCircle(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	var in circle.CircleInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	row, err := h.svc.UpdateCircle(c.Request.Context(), key, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}
func (h *Handler) DeleteCircle(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	if err = h.svc.DeleteCircle(c.Request.Context(), key); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}
func (h *Handler) ListAgents(c *gin.Context) {
	status, err := optionalStatus(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	p, l := page(c)
	out, err := h.svc.ListAgents(c.Request.Context(), c.Query("keyword"), status, p, l)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, out)
}
func (h *Handler) GetAgent(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	row, err := h.svc.GetAgent(c.Request.Context(), key)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}
func (h *Handler) CreateAgent(c *gin.Context) {
	var in circle.AgentInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	row, err := h.svc.CreateAgent(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}
func (h *Handler) UpdateAgent(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	var in circle.AgentInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	row, err := h.svc.UpdateAgent(c.Request.Context(), key, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}
func (h *Handler) AuditAgent(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	var in circle.AuditInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	if err = h.svc.AuditAgent(c.Request.Context(), key, in, middleware.AdminID(c)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, circle.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, circle.ErrBadParam), errors.Is(err, circle.ErrHasChildren), errors.Is(err, circle.ErrAlreadyAudited):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "区域代理服务异常")
	}
}
