package seckill

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/seckill"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/response"
)

type Handler struct{ svc *seckill.Service }

func NewHandler(svc *seckill.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/seckill/times", h.Times)
	r.GET("/seckill/actives", h.List)
	r.GET("/seckill/actives/:id", h.Get)
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
	res, err := h.svc.ListApp(c.Request.Context(), page, limit)
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
		if errors.Is(err, seckill.ErrNotFound) {
			response.Fail(c, http.StatusNotFound, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, a)
}
