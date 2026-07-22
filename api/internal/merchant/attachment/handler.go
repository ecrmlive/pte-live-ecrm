package attachment

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/attachment"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/upload"
)

type Handler struct {
	svc    *attachment.Service
	id     *identity.Service
	upload upload.Local
}

func NewHandler(svc *attachment.Service, id *identity.Service, up upload.Local) *Handler {
	return &Handler{svc: svc, id: id, upload: up}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/attachments/categories", h.ListCategories)
	r.POST("/attachments/categories", middleware.RequireMerchantMenu(h.id, identity.MerPermAttachmentUpload), h.CreateCategory)
	r.PUT("/attachments/categories/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermAttachmentUpload), h.UpdateCategory)
	r.DELETE("/attachments/categories/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermAttachmentDelete), h.DeleteCategory)
	r.GET("/attachments", h.List)
	r.POST("/attachments/upload", middleware.RequireMerchantMenu(h.id, identity.MerPermAttachmentUpload), h.Upload)
	r.DELETE("/attachments/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermAttachmentDelete), h.Delete)
}

func (h *Handler) ListCategories(c *gin.Context) {
	list, err := h.svc.ListCategories(c.Request.Context(), middleware.MerID(c))
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
	row, err := h.svc.CreateCategory(c.Request.Context(), middleware.MerID(c), in)
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
	row, err := h.svc.UpdateCategory(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteCategory(c.Request.Context(), middleware.MerID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cateID, _ := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 64)
	merID := int(middleware.MerID(c))
	res, err := h.svc.List(c.Request.Context(), merID, uint(cateID), page, limit)
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
	merID := middleware.MerID(c)
	src, name, err := h.upload.Save(fmt.Sprintf("merchant/%d", merID), fh)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	row, err := h.svc.CreateFile(c.Request.Context(), int(merID), middleware.AdminID(c), uint(cateID), name, src)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), int(middleware.MerID(c)), uint(id)); err != nil {
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
