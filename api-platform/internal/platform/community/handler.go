package community

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/community"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
)

type Handler struct {
	svc *community.Service
	id  *identity.Service
}

func NewHandler(svc *community.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/community/posts", h.List)
	r.GET("/community/posts/:id", h.Get)
	r.GET("/community/posts/:id/replies", h.ListReplies)
	r.POST("/community/posts/:id/audit", middleware.RequirePlatformMenu(h.id, identity.PlatPermCommunityAudit), h.Audit)
	r.DELETE("/community/posts/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermCommunityDelete), h.Delete)
	r.DELETE("/community/replies/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermCommunityDelete), h.DeleteReply)
	r.GET("/community/categories", h.ListCategories)
	r.GET("/community/topics", h.ListTopics)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var status *int8
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "状态参数错误")
			return
		}
		value := int8(parsed)
		status = &value
	}
	res, err := h.svc.ListPlatform(c.Request.Context(), status, c.Query("keyword"), page, limit)
	if err != nil {
		if errors.Is(err, community.ErrBadParam) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.Get(c.Request.Context(), uint(id), false)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListReplies(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListReplies(c.Request.Context(), uint(id), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Audit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in community.AuditInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Audit(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeletePost(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) DeleteReply(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteReply(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListCategories(c *gin.Context) {
	list, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) ListTopics(c *gin.Context) {
	list, err := h.svc.ListTopics(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, community.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, community.ErrBadParam), errors.Is(err, community.ErrForbidden):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
