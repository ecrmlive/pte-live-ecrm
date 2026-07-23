package community

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/community"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct{ svc *community.Service }

func NewHandler(svc *community.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterPublic(r gin.IRoutes) {
	r.GET("/community/categories", h.ListCategories)
	r.GET("/community/topics", h.ListTopics)
	r.GET("/community/topics/hot", h.ListHotTopics)
	r.GET("/community/posts", h.ListPosts)
	r.GET("/community/posts/:id", h.GetPost)
	r.GET("/community/posts/:id/replies", h.ListReplies)
}

func (h *Handler) RegisterAuthed(r gin.IRoutes) {
	r.POST("/community/posts", h.CreatePost)
	r.POST("/community/posts/:id/replies", h.CreateReply)
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

func (h *Handler) ListHotTopics(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	list, err := h.svc.ListHotTopics(c.Request.Context(), limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	topicID, _ := strconv.ParseUint(c.DefaultQuery("topic_id", "0"), 10, 64)
	res, err := h.svc.ListApp(c.Request.Context(), uint(topicID), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetPost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.Get(c.Request.Context(), uint(id), true)
	if err != nil {
		writeErr(c, err)
		return
	}
	if row.Status != community.StatusApproved || row.IsShow != 1 {
		response.Fail(c, http.StatusNotFound, community.ErrNotFound.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handler) CreatePost(c *gin.Context) {
	var in community.CreatePostInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreatePost(c.Request.Context(), middleware.UID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListReplies(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	res, err := h.svc.ListReplies(c.Request.Context(), uint(id), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) CreateReply(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in community.CreateReplyInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateReply(c.Request.Context(), middleware.UID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, community.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, community.ErrBadParam), errors.Is(err, community.ErrForbidden):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
