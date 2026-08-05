package attachment

import (
	"errors"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/attachment"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/upload"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *attachment.Service
	adminDB *gorm.DB
	upload  upload.Store
}

func NewHandler(svc *attachment.Service, adminDB *gorm.DB, up upload.Store) *Handler {
	return &Handler{svc: svc, adminDB: adminDB, upload: up}
}

func (h *Handler) Register(r gin.IRoutes) {
	write := middleware.RequireAdminRoles("platform", "operations")
	manage := middleware.RequireAdminMenu(h.adminDB, "content.attachment.manage")
	r.GET("/attachments/categories", h.ListCategories)
	r.POST("/attachments/categories", write, manage, h.CreateCategory)
	r.PUT("/attachments/categories/:id", write, manage, h.UpdateCategory)
	r.DELETE("/attachments/categories/:id", write, manage, h.DeleteCategory)
	r.GET("/attachments", h.List)
	r.POST("/attachments/upload", write, manage, h.Upload)
	r.DELETE("/attachments/:id", write, manage, h.Delete)
}

func (h *Handler) ListCategories(c *gin.Context) {
	list, err := h.svc.ListCategories(c.Request.Context(), 0)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var in attachment.CategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateCategory(c.Request.Context(), 0, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in attachment.CategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateCategory(c.Request.Context(), 0, uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteCategory(c.Request.Context(), 0, uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cateID, _ := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 64)
	typeValue, err := attachmentType(c.Query("type"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "素材类型错误")
		return
	}
	res, err := h.svc.List(c.Request.Context(), 0, uint(cateID), typeValue, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "请上传 file")
		return
	}
	cateID, _ := strconv.ParseUint(c.DefaultPostForm("category_id", "0"), 10, 64)
	src, name, err := h.upload.Save("platform", fh)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	row, err := h.svc.CreateFile(c.Request.Context(), 0, middleware.AdminID(c), uint(cateID), name, src, uploadType(fh.Filename, fh.Header.Get("Content-Type")))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func attachmentType(raw string) (*int8, error) {
	if raw == "" || raw == "all" {
		return nil, nil
	}
	value := int8(0)
	switch raw {
	case "image":
		return &value, nil
	case "video":
		value = 1
		return &value, nil
	default:
		return nil, errors.New("bad attachment type")
	}
}

func uploadType(filename, contentType string) int8 {
	if contentType == "" || contentType == "application/octet-stream" {
		contentType = mime.TypeByExtension(filepath.Ext(filename))
	}
	if len(contentType) >= 6 && contentType[:6] == "video/" {
		return 1
	}
	return 0
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
	case errors.Is(err, attachment.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, attachment.ErrBadParam), errors.Is(err, attachment.ErrForbidden):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
