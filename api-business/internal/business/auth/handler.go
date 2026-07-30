package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/authjwt"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
)

type Handler struct {
	svc *Service
	jwt *authjwt.Manager
}

func NewHandler(svc *Service, jwt *authjwt.Manager) *Handler { return &Handler{svc: svc, jwt: jwt} }

func (h *Handler) Register(public, authed gin.IRoutes) {
	public.POST("/auth/login", h.Login)
	public.POST("/auth/register", h.RegisterAccount)
	public.POST("/auth/refresh", h.Refresh)
	authed.GET("/auth/me", h.Me)
	authed.POST("/auth/store-context", h.IssueStoreContext)
}

type credentialRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Channel  string `json:"channel" binding:"required"`
	Nickname string `json:"nickname"`
}

func (h *Handler) Login(c *gin.Context) {
	var req credentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	channel, err := ParseChannel(req.Channel)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "不支持的登录来源")
		return
	}
	profile, err := h.svc.Login(c.Request.Context(), req.Account, req.Password, channel)
	if err != nil {
		writeError(c, err)
		return
	}
	pair, err := h.jwt.IssueCUserWithIdentityVersion(uint(profile.ID), profile.Subject, string(profile.Channel), profile.AuthVersion)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair, "user": profile})
}

// IssueStoreContext 用当前全局 C 端 JWT 加 X-AppId 换取店铺上下文 JWT。
// AppId 与 IM SDK AppId 仅由业务库投影解析，客户端不能提交 sdk_app_id。
func (h *Handler) IssueStoreContext(c *gin.Context) {
	claims := middleware.ClaimsFrom(c)
	channel, err := ParseChannel(claims.ClientPlatform)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "登录来源无效")
		return
	}
	profile, err := h.svc.Profile(c.Request.Context(), uint64(claims.UID), channel)
	if err != nil || profile.AuthVersion != claims.IdentityVersion {
		response.Fail(c, http.StatusUnauthorized, "登录已失效")
		return
	}
	context, err := h.svc.ResolveStoreContext(c.Request.Context(), c.GetHeader("X-AppId"))
	if err != nil {
		if errors.Is(err, ErrBadParam) {
			response.Fail(c, http.StatusBadRequest, "缺少 X-AppId")
			return
		}
		if errors.Is(err, ErrNotFound) {
			response.Fail(c, http.StatusNotFound, "店铺不存在或未启用")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "服务异常")
		return
	}
	pair, err := h.jwt.IssueCUserStoreContext(uint(profile.ID), profile.Subject, string(profile.Channel), profile.AuthVersion, uint(context.MerchantID), uint(context.StoreID), context.MerchantAppID, context.IMSDKAppID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair, "context": context})
}

func (h *Handler) issueForContext(c *gin.Context, profile *Profile, claims *authjwt.Claims) (*authjwt.Pair, error) {
	if claims.AuthContext != authjwt.ContextStore {
		return h.jwt.IssueCUserWithIdentityVersion(uint(profile.ID), profile.Subject, string(profile.Channel), profile.AuthVersion)
	}
	context, err := h.svc.ResolveStoreContext(c.Request.Context(), claims.MerchantAppID)
	if err != nil || uint(context.StoreID) != claims.StoreID || uint(context.MerchantID) != claims.MerID || context.IMSDKAppID != claims.IMSDKAppID {
		return nil, authjwt.ErrInvalidToken
	}
	return h.jwt.IssueCUserStoreContext(uint(profile.ID), profile.Subject, string(profile.Channel), profile.AuthVersion, uint(context.MerchantID), uint(context.StoreID), context.MerchantAppID, context.IMSDKAppID)
}

func (h *Handler) RegisterAccount(c *gin.Context) {
	var req credentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	channel, err := ParseChannel(req.Channel)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "不支持的登录来源")
		return
	}
	profile, err := h.svc.Register(c.Request.Context(), req.Account, req.Password, req.Nickname, channel)
	if err != nil {
		writeError(c, err)
		return
	}
	pair, err := h.jwt.IssueCUserWithIdentityVersion(uint(profile.ID), profile.Subject, string(profile.Channel), profile.AuthVersion)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair, "user": profile})
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := authjwt.BearerToken(c.GetHeader("Authori-zation"))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	claims, err := h.jwt.ParseExpect(refreshToken, authjwt.PortalApp, authjwt.TokenRefresh)
	if err != nil || claims.UID == 0 || claims.Scope != authjwt.ScopeCUser || claims.Channel == "" {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	channel, err := ParseChannel(claims.ClientPlatform)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	profile, err := h.svc.Profile(c.Request.Context(), uint64(claims.UID), channel)
	if err != nil || profile.AuthVersion != claims.IdentityVersion {
		response.Fail(c, http.StatusUnauthorized, "登录已失效")
		return
	}
	pair, err := h.issueForContext(c, profile, claims)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair})
}

func (h *Handler) Me(c *gin.Context) {
	claims := middleware.ClaimsFrom(c)
	channel, err := ParseChannel(claims.ClientPlatform)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "登录来源无效")
		return
	}
	profile, err := h.svc.Profile(c.Request.Context(), uint64(middleware.UID(c)), channel)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, profile)
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		response.Fail(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrAccountDisabled):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrAccountExists), errors.Is(err, ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
