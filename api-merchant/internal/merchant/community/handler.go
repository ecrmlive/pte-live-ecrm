package community

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/community"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct {
	svc *community.Service
	id  *identity.Service
}

func NewHandler(svc *community.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/community/posts", h.List)
	r.GET("/community/posts/:id", h.Get)
	r.POST("/community/posts", middleware.RequireMerchantMenu(h.id, identity.MerPermCommunityCreate), h.Create)
	r.PUT("/community/posts/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermCommunityUpdate), h.Update)
	r.DELETE("/community/posts/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermCommunityDelete), h.Delete)
	r.GET("/community/posts/:id/replies", h.ListReplies)
	r.GET("/community/categories", h.ListCategories)
	r.GET("/community/topics", h.ListTopics)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListMerchant(c.Request.Context(), middleware.MerID(c), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.Get(c.Request.Context(), uint(id), false)
	if err != nil {
		writeErr(c, err)
		return
	}
	if row.MerID != middleware.MerID(c) {
		response.Fail(c, http.StatusNotFound, community.ErrNotFound.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handler) Create(c *gin.Context) {
	var in community.CreatePostInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateMerchantPost(c.Request.Context(), middleware.MerID(c), 0, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in community.CreatePostInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateMerchantPost(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteMerchantPost(c.Request.Context(), middleware.MerID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListReplies(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	post, err := h.svc.Get(c.Request.Context(), uint(id), false)
	if err != nil {
		writeErr(c, err)
		return
	}
	if post.MerID != middleware.MerID(c) {
		response.Fail(c, http.StatusNotFound, community.ErrNotFound.Error())
		return
	}
	res, err := h.svc.ListReplies(c.Request.Context(), uint(id), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) ListCategories(c *gin.Context) {
	list, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) ListTopics(c *gin.Context) {
	list, err := h.svc.ListTopics(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, community.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, community.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, community.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
