package seckill

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/seckill"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *seckill.Service
	id  *identity.Service
}

func NewHandler(svc *seckill.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/seckill/times", h.Times)
	r.GET("/seckill/actives", h.List)
	r.POST("/seckill/actives", middleware.RequireMerchantMenu(h.id, identity.MerPermSeckillCreate), h.Create)
	r.PUT("/seckill/actives/:id", h.Update)
	r.PUT("/seckill/actives/:id/status", middleware.RequireMerchantMenu(h.id, identity.MerPermSeckillToggle), h.SetStatus)
	r.DELETE("/seckill/actives/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermSeckillDelete), h.Delete)
}

func (h *Handler) Times(c *gin.Context) {
	rows, err := h.svc.ListTimes(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
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
	var in seckill.ActiveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	a, err := h.svc.Create(c.Request.Context(), middleware.MerID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, a)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in seckill.ActiveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if in.Status != nil {
		if err := h.id.RequireMerchantMenu(c.Request.Context(), middleware.AdminID(c), identity.MerPermSeckillToggle); err != nil {
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		}
	}
	a, err := h.svc.Update(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, a)
}

type statusReq struct {
	Status *int8 `json:"status" binding:"required"`
}

func (h *Handler) SetStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req statusReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	a, err := h.svc.SetStatus(c.Request.Context(), middleware.MerID(c), uint(id), *req.Status)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, a)
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
	case errors.Is(err, seckill.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, seckill.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
