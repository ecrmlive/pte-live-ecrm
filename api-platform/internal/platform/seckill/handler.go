package seckill

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/seckill"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *seckill.Service
	adminDB *gorm.DB
}

func NewHandler(svc *seckill.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/seckill/times", h.Times)
	r.GET("/seckill/actives", h.List)
	r.GET("/seckill/actives/:id", h.Get)
	write := middleware.RequireAdminRoles("platform", "operations")
	manage := middleware.RequireAdminMenu(h.adminDB, "marketing.seckill.manage")
	r.PUT("/seckill/actives/:id", write, manage, h.Update)
	r.DELETE("/seckill/actives/:id", write, manage, h.Delete)
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
	var merID *uint
	if s := c.Query("mer_id"); s != "" {
		v, _ := strconv.ParseUint(s, 10, 64)
		u := uint(v)
		merID = &u
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
	a, err := h.svc.Get(c.Request.Context(), uint(id))
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
	a, err := h.svc.Update(c.Request.Context(), 0, uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, a)
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
	case errors.Is(err, seckill.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, seckill.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
