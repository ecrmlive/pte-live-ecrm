package combination

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/combination"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *combination.Service
	id  *identity.Service
}

func NewHandler(svc *combination.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/combination/groups", h.List)
	r.POST("/combination/groups", middleware.RequireMerchantMenu(h.id, identity.MerPermCombinationCreate), h.Create)
	r.PUT("/combination/groups/:id", h.Update)
	r.PUT("/combination/groups/:id/show", middleware.RequireMerchantMenu(h.id, identity.MerPermCombinationToggle), h.SetShow)
	r.DELETE("/combination/groups/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermCombinationDelete), h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	merID := middleware.MerID(c)
	res, err := h.svc.ListAdmin(c.Request.Context(), &merID, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	var in combination.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.svc.Create(c.Request.Context(), middleware.MerID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in combination.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if in.IsShow != nil {
		if err := h.id.RequireMerchantMenu(c.Request.Context(), middleware.AdminID(c), identity.MerPermCombinationToggle); err != nil {
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		}
	}
	g, err := h.svc.Update(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

type showReq struct {
	IsShow *int `json:"is_show" binding:"required"`
}

func (h *Handler) SetShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req showReq
	if err := c.ShouldBindJSON(&req); err != nil || req.IsShow == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.svc.SetShow(c.Request.Context(), middleware.MerID(c), uint(id), *req.IsShow)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), middleware.MerID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, combination.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, combination.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
