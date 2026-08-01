package article

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/article"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/pkg/response"
)

type Handler struct{ svc *article.Service }

func NewHandler(svc *article.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/article/categories", h.ListCategories)
	r.POST("/article/categories", h.CreateCategory)
	r.PUT("/article/categories/:id", h.UpdateCategory)
	r.DELETE("/article/categories/:id", h.DeleteCategory)

	r.GET("/articles", h.List)
	r.GET("/articles/:id", h.Get)
	r.POST("/articles", h.Create)
	r.PUT("/articles/:id", h.Update)
	r.DELETE("/articles/:id", h.Delete)
}

func (h *Handler) ListCategories(c *gin.Context) {
	list, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var in article.CategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateCategory(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in article.CategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateCategory(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cid, _ := strconv.ParseUint(c.DefaultQuery("cid", "0"), 10, 64)
	res, err := h.svc.ListAdmin(c.Request.Context(), page, limit, uint(cid))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.GetAdmin(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Create(c *gin.Context) {
	var in article.ArticleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in article.ArticleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Update(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, article.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, article.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
