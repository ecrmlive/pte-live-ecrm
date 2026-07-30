package article

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/article"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
)

type Handler struct{ svc *article.Service }

func NewHandler(svc *article.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/articles", h.List)
	r.GET("/articles/:id", h.Get)
	r.GET("/article/categories", h.ListCategories)
}

func (h *Handler) ListCategories(c *gin.Context) {
	list, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cid, _ := strconv.ParseUint(c.DefaultQuery("cid", "0"), 10, 64)
	res, err := h.svc.ListApp(c.Request.Context(), page, limit, uint(cid))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.GetApp(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, article.ErrNotFound) {
			response.Fail(c, http.StatusNotFound, "文章不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, row)
}
