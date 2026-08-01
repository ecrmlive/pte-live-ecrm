package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/authjwt"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/captchaclient"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/response"
)

type Handler struct {
	svc     *Service
	jwt     *authjwt.Manager
	captcha *captchaclient.Client
}

func NewHandler(svc *Service, jwt *authjwt.Manager, captcha *captchaclient.Client) *Handler {
	return &Handler{svc: svc, jwt: jwt, captcha: captcha}
}

func (h *Handler) Register(public, authed gin.IRoutes) {
	public.POST("/auth/login", h.Login)
	public.POST("/auth/register", h.RegisterAccount)
	public.POST("/auth/refresh", h.Refresh)
	authed.GET("/auth/me", h.Me)
	authed.POST("/auth/store-context", h.IssueStoreContext)
}

// RegisterCaptchaGateway keeps the pte-tools-captcha HMAC protocol entirely
// server-side. The browser SDK only reaches these BFF endpoints.
func (h *Handler) RegisterCaptchaGateway(r gin.IRoutes) {
	r.POST("/api/v1/challenges", h.CreateCaptchaChallenge)
	r.POST("/api/v1/challenges/:id/verify", h.VerifyCaptchaChallenge)
}

type captchaChallengeRequest struct {
	Action        string          `json:"action" binding:"required"`
	PreferredMode string          `json:"preferred_mode"`
	Locale        string          `json:"locale"`
	Theme         string          `json:"theme"`
	Client        json.RawMessage `json:"client"`
}

// Captcha actions are deliberately whitelisted. PTE token claims bind a
// completed challenge to exactly one auth purpose.
func captchaAction(raw string) (string, bool) {
	switch strings.TrimSpace(raw) {
	case "login_password", "login_sms", "register":
		return strings.TrimSpace(raw), true
	default:
		return "", false
	}
}

func (h *Handler) CreateCaptchaChallenge(c *gin.Context) {
	if h.captcha == nil {
		captchaGatewayError(c, http.StatusServiceUnavailable, "验证码服务未配置")
		return
	}
	var req captchaChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		captchaGatewayError(c, http.StatusBadRequest, "验证码请求参数错误")
		return
	}
	action, ok := captchaAction(req.Action)
	if !ok {
		captchaGatewayError(c, http.StatusBadRequest, "不支持的验证码用途")
		return
	}
	body := map[string]any{
		"action":         action,
		"preferred_mode": req.PreferredMode,
		"locale":         req.Locale,
		"theme":          req.Theme,
	}
	if len(req.Client) > 0 {
		var client any
		if err := json.Unmarshal(req.Client, &client); err != nil {
			captchaGatewayError(c, http.StatusBadRequest, "验证码客户端信息错误")
			return
		}
		body["client"] = client
	}
	result, err := h.captcha.Create(c.Request.Context(), body)
	if err != nil {
		captchaGatewayError(c, http.StatusBadGateway, "验证码服务暂不可用")
		return
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", result)
}

func (h *Handler) VerifyCaptchaChallenge(c *gin.Context) {
	if h.captcha == nil {
		captchaGatewayError(c, http.StatusServiceUnavailable, "验证码服务未配置")
		return
	}
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		captchaGatewayError(c, http.StatusBadRequest, "验证码校验参数错误")
		return
	}
	result, err := h.captcha.Verify(c.Request.Context(), c.Param("id"), body)
	if err != nil {
		captchaGatewayError(c, http.StatusBadGateway, "验证码校验服务暂不可用")
		return
	}
	var verified struct {
		Status            string `json:"status"`
		VerificationToken string `json:"verification_token"`
	}
	if err := json.Unmarshal(result, &verified); err != nil {
		captchaGatewayError(c, http.StatusBadGateway, "验证码服务响应错误")
		return
	}
	if verified.Status == "verified" {
		action, ok := captchaTokenAction(verified.VerificationToken)
		if !ok {
			captchaGatewayError(c, http.StatusBadGateway, "验证码令牌响应错误")
			return
		}
		if err := h.svc.RecordCaptchaToken(c.Request.Context(), verified.VerificationToken, action); err != nil {
			captchaGatewayError(c, http.StatusInternalServerError, "保存验证码状态失败")
			return
		}
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", result)
}

func captchaGatewayError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// captchaTokenAction only reads the pte response payload in order to persist
// its purpose. It is not used as proof: every consumer calls PTE Validate
// again before the local one-time row is consumed.
func captchaTokenAction(token string) (string, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", false
	}
	var claims struct {
		Action string `json:"action"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", false
	}
	return captchaAction(claims.Action)
}

type credentialRequest struct {
	Account      string `json:"account" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Channel      string `json:"channel" binding:"required"`
	Nickname     string `json:"nickname"`
	CaptchaToken string `json:"captcha_token"`
}

func (h *Handler) Login(c *gin.Context) {
	var req credentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.validateCaptcha(c, req.CaptchaToken, "login_password"); err != nil {
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

func (h *Handler) validateCaptcha(c *gin.Context, token, action string) error {
	if h.captcha == nil {
		response.Fail(c, http.StatusServiceUnavailable, "验证码服务未配置")
		return errors.New("captcha unavailable")
	}
	if strings.TrimSpace(token) == "" {
		response.Fail(c, http.StatusBadRequest, "请先完成安全验证")
		return errors.New("captcha token required")
	}
	if err := h.captcha.Validate(c.Request.Context(), token, action); err != nil {
		response.Fail(c, http.StatusBadRequest, "安全验证已失效，请重新验证")
		return err
	}
	if err := h.svc.ConsumeCaptchaToken(c.Request.Context(), token, action); err != nil {
		response.Fail(c, http.StatusBadRequest, "安全验证已使用或已失效，请重新验证")
		return err
	}
	return nil
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
	case errors.Is(err, ErrCaptchaUnavailable):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
