package assist

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/assist"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *assist.Service
	adminDB *gorm.DB
}

type visibilityInput struct {
	IsShow *int `json:"is_show"`
}

func NewHandler(svc *assist.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := []gin.HandlerFunc{
		middleware.RequireAdminRoles("platform", "operations"),
		middleware.RequireAdminMenu(h.adminDB, "marketing.assist.manage"),
	}
	r.GET("/assist/actives", append(access, h.List)...)
	r.GET("/assist/actives/:id", append(access, h.Get)...)
	r.PUT("/assist/actives/:id", append(access, h.Update)...)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	merID, err := optionalMerID(c.Query("mer_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "商户 ID 参数错误")
		return
	}
	res, err := h.svc.ListAdmin(c.Request.Context(), merID, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "活动 ID 参数错误")
		return
	}
	row, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func optionalMerID(raw string) (*uint, error) {
	if raw == "" {
		return nil, nil
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return nil, errors.New("invalid merchant id")
	}
	value := uint(id)
	return &value, nil
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in visibilityInput
	if err := c.ShouldBindJSON(&in); err != nil || !validVisibility(&in) {
		response.Fail(c, http.StatusBadRequest, "仅允许更新助力活动展示状态")
		return
	}
	row, err := h.svc.SetShow(c.Request.Context(), 0, uint(id), *in.IsShow)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func validVisibility(in *visibilityInput) bool {
	return in != nil && in.IsShow != nil && (*in.IsShow == 0 || *in.IsShow == 1)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), 0, uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, assist.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, assist.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
