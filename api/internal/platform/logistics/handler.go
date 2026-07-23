package logistics

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/logistics"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

// 菜单：快递公司/城市可挂 CRMEB 导入节点（sql/043）；本刀 JWT 即可。
type Handler struct{ svc *logistics.Service }

func NewHandler(svc *logistics.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/express", h.ListExpress)
	r.POST("/express", h.CreateExpress)
	r.PUT("/express/:id", h.UpdateExpress)
	r.DELETE("/express/:id", h.DeleteExpress)
	r.GET("/city", h.ListCity)
}

func (h *Handler) ListExpress(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListExpress(c.Request.Context(), page, limit, false)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) CreateExpress(c *gin.Context) {
	var in logistics.ExpressInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateExpress(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateExpress(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in logistics.ExpressInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateExpress(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteExpress(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteExpress(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListCity(c *gin.Context) {
	var parent *uint
	if raw := c.Query("parent_id"); raw != "" {
		v, err := strconv.ParseUint(raw, 10, 64)
		if err == nil {
			u := uint(v)
			parent = &u
		}
	}
	list, err := h.svc.ListCity(c.Request.Context(), parent)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, logistics.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, logistics.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
