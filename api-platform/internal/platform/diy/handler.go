package diy

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/diy"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/pkg/response"
)

type Handler struct {
	svc *diy.Service
	id  *identity.Service
}

func NewHandler(svc *diy.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/diy/pages", h.List)
	r.GET("/diy/pages/:id", h.Get)
	r.GET("/diy/editor/bootstrap", h.Bootstrap)
	r.GET("/diy/editor/bootstrap/:id", h.Bootstrap)
	r.POST("/diy/pages", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyCreate), h.Create)
	r.PUT("/diy/pages/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyUpdate), h.Update)
	r.POST("/diy/pages/:id/active", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyActive), h.SetActive)
	r.POST("/diy/pages/:id/copy", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyCreate), h.Copy)
	r.POST("/diy/pages/:id/recovery", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyUpdate), h.Recovery)
	r.DELETE("/diy/pages/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyDelete), h.Delete)
	r.GET("/diy/page-categories", h.ListCategories)
	r.POST("/diy/page-categories", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyCreate), h.CreateCategory)
	r.PUT("/diy/page-categories/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyUpdate), h.UpdateCategory)
	r.DELETE("/diy/page-categories/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyDelete), h.DeleteCategory)
	r.GET("/diy/page-links", h.ListLinks)
	r.POST("/diy/page-links", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyCreate), h.CreateLink)
	r.PUT("/diy/page-links/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyUpdate), h.UpdateLink)
	r.DELETE("/diy/page-links/:id", middleware.RequirePlatformMenu(h.id, identity.PlatPermDiyDelete), h.DeleteLink)
}

func platformLinkScope(c *gin.Context) int8 {
	if c.DefaultQuery("scope", "platform") == "merchant" {
		return 1
	}
	return 0
}

func (h *Handler) ListCategories(c *gin.Context) {
	rows, err := h.svc.ListCategories(c.Request.Context(), platformLinkScope(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var in diy.CategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateCategory(c.Request.Context(), platformLinkScope(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in diy.CategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateCategory(c.Request.Context(), uint(id), platformLinkScope(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteCategory(c.Request.Context(), uint(id), platformLinkScope(c)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListLinks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	f := diy.LinkListFilter{IsMer: platformLinkScope(c), Page: page, Limit: limit, Name: c.Query("name")}
	if value := c.Query("status"); value != "" {
		n, _ := strconv.ParseInt(value, 10, 8)
		status := int8(n)
		f.Status = &status
	}
	rows, err := h.svc.ListLinks(c.Request.Context(), f)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, rows)
}

func (h *Handler) CreateLink(c *gin.Context) {
	var in diy.LinkInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateLink(c.Request.Context(), platformLinkScope(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateLink(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in diy.LinkInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateLink(c.Request.Context(), uint(id), platformLinkScope(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteLink(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteLink(c.Request.Context(), uint(id), platformLinkScope(c)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	f := diy.ListFilter{MerID: 0, Page: page, Limit: limit, Name: c.Query("name")}
	if v := c.Query("is_diy"); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		iv := int8(n)
		f.IsDiy = &iv
	}
	if v := c.Query("status"); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		iv := int8(n)
		f.Status = &iv
	}
	res, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	if p.MerID != 0 {
		response.Fail(c, http.StatusNotFound, diy.ErrNotFound.Error())
		return
	}
	response.OK(c, p)
}

func (h *Handler) Bootstrap(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		if q := c.Query("id"); q != "" {
			id, _ = strconv.ParseUint(q, 10, 64)
		}
	}
	boot, err := h.svc.EditorBootstrap(c.Request.Context(), uint(id), 0)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, boot)
}

func (h *Handler) Create(c *gin.Context) {
	var in diy.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	p, err := h.svc.Create(c.Request.Context(), 0, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in diy.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	p, err := h.svc.Update(c.Request.Context(), uint(id), 0, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) SetActive(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.SetActive(c.Request.Context(), uint(id), 0)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Copy(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Copy(c.Request.Context(), uint(id), 0)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Recovery(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.Recovery(c.Request.Context(), uint(id), 0)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), uint(id), 0); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, diy.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, diy.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
