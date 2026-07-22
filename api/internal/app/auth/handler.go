package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/authjwt"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *identity.Service
	jwt *authjwt.Manager
}

func NewHandler(svc *identity.Service, jwt *authjwt.Manager) *Handler {
	return &Handler{svc: svc, jwt: jwt}
}

func (h *Handler) Register(r gin.IRoutes, authed gin.IRoutes) {
	r.POST("/auth/login", h.Login)
	r.POST("/auth/register", h.RegisterAccount)
	r.POST("/auth/refresh", h.Refresh)
	authed.GET("/auth/me", h.Me)
}

type loginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	user, err := h.svc.LoginApp(c.Request.Context(), req.Account, req.Password, c.ClientIP())
	if err != nil {
		writeIdentityErr(c, err)
		return
	}
	pair, err := h.jwt.Issue(authjwt.PortalApp, 0, 0, user.UID, user.Account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{
		"token": pair,
		"user": identity.AppProfile{
			UID:      user.UID,
			Account:  user.Account,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Phone:    user.Phone,
		},
	})
}

func (h *Handler) RegisterAccount(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	user, err := h.svc.RegisterApp(c.Request.Context(), req.Account, req.Password, req.Nickname, c.ClientIP())
	if err != nil {
		writeIdentityErr(c, err)
		return
	}
	pair, err := h.jwt.Issue(authjwt.PortalApp, 0, 0, user.UID, user.Account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{
		"token": pair,
		"user": identity.AppProfile{
			UID:      user.UID,
			Account:  user.Account,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Phone:    user.Phone,
		},
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	claims, err := h.jwt.ParseExpect(req.RefreshToken, authjwt.PortalApp, authjwt.TokenRefresh)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	pair, err := h.jwt.Issue(authjwt.PortalApp, 0, 0, claims.UID, claims.Account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair})
}

func (h *Handler) Me(c *gin.Context) {
	profile, err := h.svc.AppProfile(c.Request.Context(), middleware.UID(c))
	if err != nil {
		writeIdentityErr(c, err)
		return
	}
	response.OK(c, profile)
}

func writeIdentityErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, identity.ErrInvalidCredentials),
		errors.Is(err, identity.ErrBadPassword):
		response.Fail(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, identity.ErrAccountDisabled),
		errors.Is(err, identity.ErrMerchantDisabled):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, identity.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, identity.ErrWeakPassword),
		errors.Is(err, identity.ErrAccountExists):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
