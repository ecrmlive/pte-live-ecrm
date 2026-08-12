package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/cloudconfig"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/captchaclient"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/smsclient"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/wechatmini"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc     *Service
	jwt     *authjwt.Manager
	captcha *captchaclient.Client
	cloud   *cloudconfig.Service
	wechat  *wechatmini.Client
	sms     *smsclient.Client
}

func NewHandler(svc *Service, jwt *authjwt.Manager, captcha *captchaclient.Client, cloud ...*cloudconfig.Service) *Handler {
	var configService *cloudconfig.Service
	if len(cloud) > 0 {
		configService = cloud[0]
	}
	return &Handler{svc: svc, jwt: jwt, captcha: captcha, cloud: configService, wechat: &wechatmini.Client{}, sms: smsclient.New()}
}

func (h *Handler) Register(public, authed gin.IRoutes) {
	public.POST("/auth/login", h.Login)
	public.POST("/auth/wechat/mini-program", h.MiniProgramLogin)
	public.POST("/auth/register", h.RegisterAccount)
	public.POST("/auth/refresh", h.Refresh)
	public.POST("/auth/sms/send", h.SendSMS)
	public.POST("/auth/sms/login", h.SMSLogin)
	public.POST("/auth/password/reset", h.ResetPasswordBySMS)
	authed.POST("/auth/mobile/bind", h.BindMobile)
	authed.POST("/auth/mobile/change", h.ChangeMobile)
	authed.GET("/auth/me", h.Me)
	authed.PATCH("/auth/me", h.UpdateMe)
	authed.POST("/auth/store-context", h.IssueStoreContext)
}

type miniProgramLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// MiniProgramLogin is the silent login exchange for a WeChat Mini Program.
// It accepts only the short-lived code from uni.login; app_secret stays in the
// encrypted platform cloud configuration.
func (h *Handler) MiniProgramLogin(c *gin.Context) {
	if h.cloud == nil || h.wechat == nil {
		response.Fail(c, http.StatusServiceUnavailable, "微信小程序登录未配置")
		return
	}
	var req miniProgramLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "微信登录参数错误")
		return
	}
	values, err := h.cloud.Values(c.Request.Context(), "wechat_mini_program")
	if err != nil || strings.TrimSpace(values["enabled"]) != "true" {
		response.Fail(c, http.StatusServiceUnavailable, "微信小程序登录未配置")
		return
	}
	session, err := h.wechat.Code2Session(c.Request.Context(), wechatmini.Config{AppID: values["app_id"], AppSecret: values["app_secret"]}, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, wechatmini.ErrNotConfigured):
			response.Fail(c, http.StatusServiceUnavailable, "微信小程序登录未配置")
		case errors.Is(err, wechatmini.ErrInvalidCode):
			response.Fail(c, http.StatusUnauthorized, "微信登录凭证已失效，请重试")
		default:
			response.Fail(c, http.StatusBadGateway, "微信登录服务暂不可用")
		}
		return
	}
	profile, err := h.svc.LoginOrRegisterExternal(c.Request.Context(), session.OpenID, "微信用户", ChannelMiniProgram)
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
	if err := h.validateCaptcha(c, req.CaptchaToken, "register"); err != nil {
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

type smsSendRequest struct {
	Mobile       string `json:"mobile" binding:"required"`
	Purpose      string `json:"purpose" binding:"required"`
	CaptchaToken string `json:"captcha_token" binding:"required"`
}
type mobileBindRequest struct {
	Mobile string `json:"mobile" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

func (h *Handler) SendSMS(c *gin.Context) {
	var req smsSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "短信参数错误")
		return
	}
	if req.Purpose != "binding" && req.Purpose != "login" && req.Purpose != "reset_password" && req.Purpose != "change_mobile" {
		response.Fail(c, http.StatusBadRequest, "短信用途错误")
		return
	}
	if err := h.validateCaptcha(c, req.CaptchaToken, "login_sms"); err != nil {
		return
	}
	if h.cloud == nil {
		response.Fail(c, http.StatusServiceUnavailable, "短信服务未配置")
		return
	}
	values, err := h.cloud.Values(c.Request.Context(), "tencent_sms")
	if err != nil || strings.TrimSpace(values["enabled"]) != "true" {
		response.Fail(c, http.StatusServiceUnavailable, "短信服务未配置")
		return
	}
	code, err := h.svc.CreateSMSCode(c.Request.Context(), req.Mobile, req.Purpose)
	if err != nil {
		writeError(c, err)
		return
	}
	if err = h.sms.SendTencent(c.Request.Context(), smsclient.TencentConfig{
		AppKey:      values["app_key"],
		SDKAppID:    values["sdk_app_id"],
		SignContent: values["sign_content"],
		TemplateID:  values["template_id"],
	}, strings.TrimSpace(req.Mobile), code); err != nil {
		_ = h.svc.DiscardSMSCode(c.Request.Context(), req.Mobile, req.Purpose, code)
		response.Fail(c, http.StatusBadGateway, "短信发送失败，请稍后重试")
		return
	}
	response.OK(c, gin.H{"sent": true, "expires_in": int(smsCodeTTL.Seconds())})
}

type smsLoginRequest struct {
	Mobile  string `json:"mobile" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Channel string `json:"channel" binding:"required"`
}

func (h *Handler) SMSLogin(c *gin.Context) {
	var req smsLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "短信登录参数错误")
		return
	}
	channel, err := ParseChannel(req.Channel)
	if err != nil || channel != ChannelH5 {
		response.Fail(c, http.StatusBadRequest, "短信登录仅支持 H5")
		return
	}
	if err = h.svc.ConsumeSMSCode(c.Request.Context(), req.Mobile, "login", req.Code); err != nil {
		writeError(c, err)
		return
	}
	profile, err := h.svc.LoginOrRegisterMobile(c.Request.Context(), req.Mobile, channel)
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

type resetPasswordRequest struct {
	Mobile          string `json:"mobile" binding:"required"`
	Code            string `json:"code" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

func (h *Handler) ResetPasswordBySMS(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.NewPassword != req.ConfirmPassword {
		response.Fail(c, http.StatusBadRequest, "新密码确认不一致")
		return
	}
	if err := h.svc.ConsumeSMSCode(c.Request.Context(), req.Mobile, "reset_password", req.Code); err != nil {
		writeError(c, err)
		return
	}
	if err := h.svc.ResetPasswordByMobile(c.Request.Context(), req.Mobile, req.NewPassword); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"reset": true})
}

type mobileChangeRequest struct {
	OldMobile string `json:"old_mobile" binding:"required"`
	OldCode   string `json:"old_code" binding:"required"`
	NewMobile string `json:"new_mobile" binding:"required"`
	NewCode   string `json:"new_code" binding:"required"`
}

func (h *Handler) ChangeMobile(c *gin.Context) {
	var req mobileChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "换绑参数错误")
		return
	}
	if err := h.svc.ChangeMobile(c.Request.Context(), uint64(middleware.UID(c)), req.OldMobile, req.OldCode, req.NewMobile, req.NewCode); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"mobile": strings.TrimSpace(req.NewMobile)})
}

func (h *Handler) BindMobile(c *gin.Context) {
	var req mobileBindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "绑定参数错误")
		return
	}
	if err := h.svc.ConsumeSMSCode(c.Request.Context(), req.Mobile, "binding", req.Code); err != nil {
		writeError(c, err)
		return
	}
	if err := h.svc.BindMobile(c.Request.Context(), uint64(middleware.UID(c)), req.Mobile); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"mobile": strings.TrimSpace(req.Mobile)})
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

type updateProfileRequest struct {
	Nickname string `json:"nickname"`
}

func (h *Handler) UpdateMe(c *gin.Context) {
	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "资料参数错误")
		return
	}
	claims := middleware.ClaimsFrom(c)
	channel, err := ParseChannel(claims.ClientPlatform)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "登录来源无效")
		return
	}
	if err := h.svc.UpdateNickname(c.Request.Context(), uint64(middleware.UID(c)), req.Nickname); err != nil {
		writeError(c, err)
		return
	}
	profile, err := h.svc.Profile(c.Request.Context(), uint64(middleware.UID(c)), channel)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, profile)
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
	case errors.Is(err, ErrSMSRateLimited):
		response.Fail(c, http.StatusTooManyRequests, err.Error())
	case errors.Is(err, ErrSMSInvalid):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
