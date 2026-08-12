package notification

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/notification"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const menuManage = "setting.notice.config.manage"

type Handler struct {
	svc     *notification.Service
	adminDB *gorm.DB
}

func NewHandler(svc *notification.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	platform := middleware.RequireAdminRoles("platform")
	manage := middleware.RequireAdminMenu(h.adminDB, menuManage)
	r.GET("/notification-configs", platform, manage, h.List)
	r.GET("/notification-configs/:id", platform, manage, h.Get)
	r.PUT("/notification-configs/:id", platform, manage, h.Save)
	r.POST("/notification-configs/sync", platform, manage, h.Sync)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	result, err := h.svc.List(c.Request.Context(), notification.Audience(c.Query("audience")), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	result, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Save(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var input notification.SaveInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	result, err := h.svc.Save(c.Request.Context(), uint(id), input)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Sync(c *gin.Context) {
	var input struct {
		Audience notification.Audience `json:"audience"`
		Channel  notification.Channel  `json:"channel"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.Sync(c.Request.Context(), input.Audience, input.Channel); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, notification.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, notification.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, notification.ErrSyncUnavailable):
		response.Fail(c, http.StatusConflict, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "保存通知配置失败")
	}
}
