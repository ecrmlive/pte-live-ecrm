package community

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/community"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *community.Service
	adminDB *gorm.DB
}

func NewHandler(svc *community.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	write := middleware.RequireAdminRoles("platform", "operations")
	listRead := middleware.RequireAdminMenuAny(
		h.adminDB,
		"content.community_list.read",
		"content.community_list.manage",
		"content.community.audit",
		"content.community.delete",
	)
	listManage := middleware.RequireAdminMenuAny(
		h.adminDB,
		"content.community_list.manage",
		"content.community.audit",
		"content.community.delete",
	)
	replyRead := middleware.RequireAdminMenuAny(
		h.adminDB,
		"content.community_reply.read",
		"content.community_reply.manage",
		"content.community_list.read",
		"content.community_list.manage",
		"content.community.delete",
	)
	replyManage := middleware.RequireAdminMenuAny(
		h.adminDB,
		"content.community_reply.manage",
		"content.community_list.manage",
		"content.community.delete",
	)
	categoryRead := middleware.RequireAdminMenuAny(
		h.adminDB,
		"content.community_category.read",
		"content.community_category.manage",
		"content.community_topic.read",
		"content.community_topic.manage",
		"content.community_list.read",
		"content.community_list.manage",
	)
	categoryWrite := middleware.RequireAdminMenu(h.adminDB, "content.community_category.manage")
	topicRead := middleware.RequireAdminMenuAny(
		h.adminDB,
		"content.community_topic.read",
		"content.community_topic.manage",
		"content.community_list.read",
		"content.community_list.manage",
	)
	topicWrite := middleware.RequireAdminMenu(h.adminDB, "content.community_topic.manage")
	r.GET("/community/posts", write, listRead, h.List)
	r.GET("/community/posts/:id", write, listRead, h.Get)
	r.GET("/community/posts/:id/replies", write, replyRead, h.ListReplies)
	r.POST("/community/posts/:id/audit", write, listManage, h.Audit)
	r.PUT("/community/posts/:id/star", write, listManage, h.UpdateStar)
	r.PUT("/community/posts/:id/show", write, listManage, h.SwitchShow)
	r.DELETE("/community/posts/:id", write, listManage, h.Delete)
	r.GET("/community/replies", write, replyRead, h.ListAllReplies)
	r.POST("/community/replies/:id/audit", write, replyManage, h.AuditReply)
	r.DELETE("/community/replies/:id", write, replyManage, h.DeleteReply)
	r.GET("/community/categories", write, categoryRead, h.ListCategories)
	r.POST("/community/categories", write, categoryWrite, h.CreateCategory)
	r.PUT("/community/categories/:id", write, categoryWrite, h.UpdateCategory)
	r.PUT("/community/categories/:id/status", write, categoryWrite, h.SetCategoryStatus)
	r.DELETE("/community/categories/:id", write, categoryWrite, h.DeleteCategory)
	r.GET("/community/topics", write, topicRead, h.ListTopics)
	r.POST("/community/topics", write, topicWrite, h.CreateTopic)
	r.PUT("/community/topics/:id", write, topicWrite, h.UpdateTopic)
	r.PUT("/community/topics/:id/status", write, topicWrite, h.SetTopicStatus)
	r.PUT("/community/topics/:id/hot", write, topicWrite, h.SetTopicHot)
	r.DELETE("/community/topics/:id", write, topicWrite, h.DeleteTopic)
}

func parseOptionalInt8(raw string) (*int8, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}
	value := int8(parsed)
	return &value, nil
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status, err := parseOptionalInt8(c.Query("status"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "状态参数错误")
		return
	}
	isShow, err := parseOptionalInt8(c.Query("is_show"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "显示参数错误")
		return
	}
	isType, err := parseOptionalInt8(c.Query("is_type"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "类型参数错误")
		return
	}
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)
	topicID, _ := strconv.ParseUint(c.Query("topic_id"), 10, 64)
	res, err := h.svc.ListPlatform(c.Request.Context(), community.ListFilter{
		Status: status, Keyword: c.Query("keyword"), Page: page, Limit: limit,
		CategoryID: uint(categoryID), TopicID: uint(topicID),
		IsShow: isShow, IsType: isType,
		AuthorType: c.Query("author_type"), AuthorKW: c.Query("author"),
	})
	if err != nil {
		if errors.Is(err, community.ErrBadParam) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
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
	response.OK(c, row)
}

func (h *Handler) ListReplies(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListReplies(c.Request.Context(), uint(id), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Audit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in community.AuditInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Audit(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateStar(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in community.StarInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateStar(c.Request.Context(), uint(id), in.Start)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SwitchShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		IsShow *int8 `json:"is_show"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.IsShow == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SwitchShow(c.Request.Context(), uint(id), *body.IsShow)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeletePost(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) DeleteReply(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteReply(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListAllReplies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListAllReplies(c.Request.Context(), community.ReplyListFilter{
		Keyword:  c.Query("keyword"),
		Username: c.Query("username"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     page,
		Limit:    limit,
	})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) AuditReply(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in community.ReplyAuditInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.AuditReply(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
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
	var in community.CategoryInput
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
	var in community.CategoryInput
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
		IsShow *int8 `json:"is_show"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.IsShow == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetCategoryShow(c.Request.Context(), uint(id), *body.IsShow); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, community.ErrForbidden) {
			response.Fail(c, http.StatusBadRequest, "分类下仍有话题或帖子，无法删除")
			return
		}
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListTopics(c *gin.Context) {
	list, err := h.svc.ListTopics(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) CreateTopic(c *gin.Context) {
	var in community.TopicInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateTopic(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateTopic(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in community.TopicInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateTopic(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SetTopicStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		Status *int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Status == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetTopicStatus(c.Request.Context(), uint(id), *body.Status); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) SetTopicHot(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		IsHot *int8 `json:"is_hot"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.IsHot == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetTopicHot(c.Request.Context(), uint(id), *body.IsHot); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) DeleteTopic(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteTopic(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, community.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, community.ErrDuplicate):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, community.ErrBadParam), errors.Is(err, community.ErrForbidden):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
