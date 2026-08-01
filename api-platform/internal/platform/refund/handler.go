package refund

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/aftersale"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
)

type Handler struct {
	svc *aftersale.Service
	id  *identity.Service
}

func NewHandler(svc *aftersale.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/refunds", h.List)
	r.GET("/refunds/:id", h.Get)
	r.POST("/refunds/:id/approve", middleware.RequirePlatformMenu(h.id, identity.PlatPermRefundApprove), h.Approve)
	r.POST("/refunds/:id/reject", middleware.RequirePlatformMenu(h.id, identity.PlatPermRefundReject), h.Reject)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var status *int8
	if s := c.Query("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			vv := int8(v)
			status = &vv
		}
	}
	regionIDs, err := h.id.PlatformRegionScope(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "区域数据范围读取失败")
		return
	}
	res, err := h.svc.ListPlatformByRegions(c.Request.Context(), regionIDs, status, page, limit)
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
		response.Fail(c, http.StatusInternalServerError, "区域数据范围读取失败")
		return
	}
	ro, err := h.svc.GetPlatformByRegions(c.Request.Context(), uint(id), regionIDs)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, ro)
}

func (h *Handler) Approve(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	// merID=0：平台代审
	if err := h.svc.Approve(c.Request.Context(), 0, uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Reject(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in aftersale.RejectInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.Reject(c.Request.Context(), 0, uint(id), in); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, aftersale.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, aftersale.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, aftersale.ErrBadParam),
		errors.Is(err, aftersale.ErrBadStatus),
		errors.Is(err, aftersale.ErrAlreadyDone),
		errors.Is(err, aftersale.ErrRejectNeedMsg),
		errors.Is(err, aftersale.ErrProductInvalid):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
