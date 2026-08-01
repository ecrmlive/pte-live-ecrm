package serviceportal

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/chat"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/cs"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/trade"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
)

// Handler 客服工作台：正式 JWT + 查单 + 快捷回复 + 会话桥
type Handler struct {
	id   *identity.Service
	jwt  *authjwt.Manager
	ord  *trade.Service
	cs   *cs.Service
	chat *chat.Service
}

func NewHandler(id *identity.Service, jwt *authjwt.Manager, ord *trade.Service, csSvc *cs.Service, chatSvc *chat.Service) *Handler {
	return &Handler{id: id, jwt: jwt, ord: ord, cs: csSvc, chat: chatSvc}
}

func (h *Handler) Register(public, authed gin.IRoutes) {
	public.POST("/auth/login", h.Login)
	public.POST("/auth/refresh", h.Refresh)
	public.GET("/ping", h.Ping)
	authed.GET("/auth/me", h.Me)
	authed.GET("/orders/:id", h.OrderDetail)
	authed.GET("/quick-replies", h.QuickReplies)
	authed.GET("/threads", h.ListThreads)
	authed.POST("/threads/:id/claim", h.ClaimThread)
	authed.GET("/threads/:id/messages", h.ListMessages)
	authed.POST("/threads/:id/messages", h.SendMessage)
	authed.GET("/im/credential", h.IMCredential)
}

type loginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Ping(c *gin.Context) {
	response.OK(c, gin.H{"prefix": "/api/service/v1", "ok": true, "auth": "jwt"})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	staff, mer, err := h.id.LoginCustomerService(c.Request.Context(), req.Account, req.Password)
	if err != nil {
		writeIDErr(c, err)
		return
	}
	// 客服属于统一后台角色，不单独签发 service JWT；数据范围仍由 mer_id 和 RBAC 约束。
	pair, err := h.jwt.Issue(authjwt.PortalPlatform, staff.ServiceID, staff.MerID, 0, staff.Account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{
		"token": pair,
		"user": identity.StoreServiceProfile{
			ServiceID: staff.ServiceID,
			MerID:     staff.MerID,
			MerName:   mer.MerName,
			Account:   staff.Account,
			Nickname:  staff.Nickname,
			IsVerify:  staff.IsVerify,
			IsGoods:   staff.IsGoods,
			Customer:  staff.Customer,
		},
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := authjwt.BearerToken(c.GetHeader("Authori-zation"))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	claims, err := h.jwt.ParseExpect(refreshToken, authjwt.PortalPlatform, authjwt.TokenRefresh)
	if err != nil || claims.MerID == 0 || claims.AdminID == 0 {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	pair, err := h.jwt.Issue(authjwt.PortalPlatform, claims.AdminID, claims.MerID, 0, claims.Account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair})
}

func (h *Handler) Me(c *gin.Context) {
	profile, err := h.id.StoreServiceProfile(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		writeIDErr(c, err)
		return
	}
	if profile.MerID != middleware.MerID(c) {
		response.Fail(c, http.StatusForbidden, "商户上下文不匹配")
		return
	}
	if profile.Customer != 1 {
		response.Fail(c, http.StatusForbidden, "无客服权限")
		return
	}
	response.OK(c, profile)
}

func (h *Handler) OrderDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	merID := middleware.MerID(c)
	if id == 0 || merID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	o, err := h.ord.GetMerchantOrder(c.Request.Context(), merID, uint(id))
	if err != nil {
		if errors.Is(err, trade.ErrNotFound) || errors.Is(err, trade.ErrForbidden) {
			response.Fail(c, http.StatusNotFound, "订单不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, o)
}

func (h *Handler) QuickReplies(c *gin.Context) {
	merID := middleware.MerID(c)
	rows, err := h.cs.ListEnabled(c.Request.Context(), merID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		list = append(list, map[string]any{
			"id": r.ServiceReplyID, "title": r.Keyword, "keyword": r.Keyword,
			"content": r.Content, "type": r.Type,
		})
	}
	response.OK(c, gin.H{"list": list, "mer_id": merID})
}

func (h *Handler) ListThreads(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.chat.ListServiceThreads(c.Request.Context(), middleware.MerID(c), page, limit)
	if err != nil {
		writeChatErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) ClaimThread(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.chat.Claim(c.Request.Context(), middleware.MerID(c), middleware.AdminID(c), uint(id))
	if err != nil {
		writeChatErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListMessages(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	res, err := h.chat.ListMessages(c.Request.Context(), middleware.MerID(c), 0, middleware.AdminID(c), uint(id), page, limit)
	if err != nil {
		writeChatErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) SendMessage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in chat.SendInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	msg, err := h.chat.SendService(c.Request.Context(), middleware.MerID(c), middleware.AdminID(c), uint(id), in)
	if err != nil {
		writeChatErr(c, err)
		return
	}
	response.OK(c, msg)
}

func (h *Handler) IMCredential(c *gin.Context) {
	threadID, _ := strconv.ParseUint(c.DefaultQuery("thread_id", "0"), 10, 64)
	cred, err := h.chat.IssueCredentialForThread(c.Request.Context(), "service", middleware.AdminID(c), uint(threadID))
	if err != nil {
		writeChatErr(c, err)
		return
	}
	response.OK(c, cred)
}

func writeIDErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, identity.ErrInvalidCredentials):
		response.Fail(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, identity.ErrAccountDisabled), errors.Is(err, identity.ErrMerchantDisabled),
		errors.Is(err, identity.ErrNoCustomerPerm):
		response.Fail(c, http.StatusForbidden, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}

func writeChatErr(c *gin.Context, err error) {
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
