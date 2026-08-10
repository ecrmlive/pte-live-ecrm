package article

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/article"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *article.Service
	adminDB *gorm.DB
}

func NewHandler(svc *article.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	platformOrOps := middleware.RequireAdminRoles("platform", "operations")
	// 分类列表供文章表单下拉复用：分类读 或 文章读/写 任一即可
	categoryRead := middleware.RequireAdminMenuAny(
		h.adminDB,
		"content.article_category.read",
		"content.article_category.manage",
		"content.article.read",
		"content.article.manage",
	)
	categoryWrite := middleware.RequireAdminMenu(h.adminDB, "content.article_category.manage")
	r.GET("/article/categories", platformOrOps, categoryRead, h.ListCategories)
	r.POST("/article/categories", platformOrOps, categoryWrite, h.CreateCategory)
	r.PUT("/article/categories/:id", platformOrOps, categoryWrite, h.UpdateCategory)
	r.PUT("/article/categories/:id/status", platformOrOps, categoryWrite, h.SetCategoryStatus)
	r.DELETE("/article/categories/:id", platformOrOps, categoryWrite, h.DeleteCategory)

	articleRead := middleware.RequireAdminMenuAny(h.adminDB, "content.article.read", "content.article.manage")
	articleWrite := middleware.RequireAdminMenu(h.adminDB, "content.article.manage")
	r.GET("/articles", platformOrOps, articleRead, h.List)
	r.GET("/articles/:id", platformOrOps, articleRead, h.Get)
	r.POST("/articles", platformOrOps, articleWrite, h.Create)
	r.PUT("/articles/:id", platformOrOps, articleWrite, h.Update)
	r.DELETE("/articles/:id", platformOrOps, articleWrite, h.Delete)
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

func (h *Handler) SetCategoryStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		Status *int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Status == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetCategoryStatus(c.Request.Context(), uint(id), *body.Status); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
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
	title := c.Query("title")
	if title == "" {
		title = c.Query("keyword")
	}
	res, err := h.svc.ListAdmin(c.Request.Context(), page, limit, uint(cid), title)
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
