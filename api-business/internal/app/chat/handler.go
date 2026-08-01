package chat

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/chat"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
)

type Handler struct{ svc *chat.Service }

func NewHandler(svc *chat.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/cs/threads", h.Open)
	r.GET("/cs/threads", h.List)
	r.GET("/cs/threads/:id/messages", h.ListMessages)
	r.POST("/cs/threads/:id/messages", h.Send)
	r.GET("/cs/im/credential", h.IMCredential)
}

func (h *Handler) Open(c *gin.Context) {
	var in chat.OpenInput
	if err := c.ShouldBindJSON(&in); err != nil || in.MerID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.OpenUserThread(c.Request.Context(), middleware.UID(c), in.MerID)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListUserThreads(c.Request.Context(), middleware.UID(c), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) ListMessages(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	res, err := h.svc.ListMessages(c.Request.Context(), 0, middleware.UID(c), 0, uint(id), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Send(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in chat.SendInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	msg, err := h.svc.SendUser(c.Request.Context(), middleware.UID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, msg)
}

func (h *Handler) IMCredential(c *gin.Context) {
	threadID, _ := strconv.ParseUint(c.DefaultQuery("thread_id", "0"), 10, 64)
	cred, err := h.svc.IssueCredentialForThread(c.Request.Context(), "app", middleware.UID(c), uint(threadID))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, cred)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, chat.ErrBadParam), errors.Is(err, chat.ErrTextViaCS), errors.Is(err, chat.ErrIMRemoteRequired):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, chat.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, chat.ErrForbidden), errors.Is(err, chat.ErrClosed):
		response.Fail(c, http.StatusForbidden, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
