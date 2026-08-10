package usertag

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/usertag"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	// RequireAdminMenu 仅认 kind=button；page 码 user.label / user.group 不能用于接口鉴权。
	menuUserLabelRead   = "user.label.read"
	menuUserLabelManage = "user.label.manage"
	menuUserGroupRead   = "user.group.read"
	menuUserGroupManage = "user.group.manage"
)

type Handler struct {
	svc     *usertag.Service
	adminDB *gorm.DB
}

func NewHandler(svc *usertag.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	labelRead := middleware.RequireAdminMenu(h.adminDB, menuUserLabelRead)
	labelManage := middleware.RequireAdminMenu(h.adminDB, menuUserLabelManage)
	groupRead := middleware.RequireAdminMenu(h.adminDB, menuUserGroupRead)
	groupManage := middleware.RequireAdminMenu(h.adminDB, menuUserGroupManage)
	r.GET("/user/labels", access, labelRead, h.ListLabels)
	r.POST("/user/labels", access, labelManage, h.CreateLabel)
	r.PUT("/user/labels/:id", access, labelManage, h.UpdateLabel)
	r.DELETE("/user/labels/:id", access, labelManage, h.DeleteLabel)

	r.GET("/user/groups", access, groupRead, h.ListGroups)
	r.POST("/user/groups", access, groupManage, h.CreateGroup)
	r.PUT("/user/groups/:id", access, groupManage, h.UpdateGroup)
	r.DELETE("/user/groups/:id", access, groupManage, h.DeleteGroup)

	r.GET("/user/:uid/labels", access, labelRead, h.ListUserLabels)
	r.PUT("/user/:uid/labels", access, labelManage, h.MarkUser)
}

func (h *Handler) ListLabels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListLabels(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) CreateLabel(c *gin.Context) {
	var in usertag.LabelInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateLabel(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateLabel(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in usertag.LabelInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateLabel(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteLabel(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteLabel(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListGroups(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) CreateGroup(c *gin.Context) {
	var in usertag.GroupInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateGroup(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateGroup(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in usertag.GroupInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateGroup(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteGroup(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteGroup(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListUserLabels(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Param("uid"), 10, 64)
	list, err := h.svc.ListUserLabels(c.Request.Context(), uint(uid))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) MarkUser(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Param("uid"), 10, 64)
	var body struct {
		LabelIDs []uint `json:"label_ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.MarkUser(c.Request.Context(), usertag.MarkInput{UID: uint(uid), LabelIDs: body.LabelIDs}); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usertag.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, usertag.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
