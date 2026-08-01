package combination

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/combination"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/pkg/response"
)

type Handler struct{ svc *combination.Service }

func NewHandler(svc *combination.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/combination/groups", h.List)
	r.GET("/combination/groups/:id", h.Get)
	r.PUT("/combination/groups/:id", h.Update)
	r.DELETE("/combination/groups/:id", h.Delete)
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
	g, err := h.svc.Get(c.Request.Context(), uint(id))
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
	g, err := h.svc.Update(c.Request.Context(), 0, uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
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
	case errors.Is(err, combination.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, combination.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
