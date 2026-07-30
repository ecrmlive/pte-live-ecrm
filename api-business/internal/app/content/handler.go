package content

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/content"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
)

type Handler struct {
	svc *content.Service
}

func NewHandler(svc *content.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/notices", h.List)
	r.GET("/notices/:id", h.Get)
	r.GET("/agreements/:key", h.GetAgreement)
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
	n, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		response.Fail(c, http.StatusNotFound, "公告不存在")
		return
	}
	if n.IsShow != 1 {
		response.Fail(c, http.StatusNotFound, "公告不存在")
		return
	}
	response.OK(c, n)
}

func (h *Handler) GetAgreement(c *gin.Context) {
	row, err := h.svc.GetAgreement(c.Request.Context(), c.Param("key"))
	if err != nil {
		if errors.Is(err, content.ErrAgreeNotFound) {
			response.Fail(c, http.StatusNotFound, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, row)
}
