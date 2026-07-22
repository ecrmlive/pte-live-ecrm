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
	r.POST("/auth/refresh", h.Refresh)
	authed.GET("/auth/me", h.Me)
	authed.GET("/auth/menus", h.Menus)
	authed.GET("/auth/permissions", h.Permissions)
	authed.PUT("/auth/password", h.ChangePassword)
}

type loginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type passwordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	admin, mer, err := h.svc.LoginMerchant(c.Request.Context(), req.Account, req.Password, c.ClientIP())
	if err != nil {
		writeIdentityErr(c, err)
		return
	}
	pair, err := h.jwt.Issue(authjwt.PortalMerchant, admin.MerchantAdminID, admin.MerID, 0, admin.Account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{
		"token": pair,
		"user": identity.MerchantProfile{
			MerchantAdminID: admin.MerchantAdminID,
			MerID:           admin.MerID,
			MerName:         mer.MerName,
			Account:         admin.Account,
			RealName:        admin.RealName,
			Phone:           admin.Phone,
			Roles:           admin.Roles,
			Level:           admin.Level,
		},
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	claims, err := h.jwt.ParseExpect(req.RefreshToken, authjwt.PortalMerchant, authjwt.TokenRefresh)
	if err != nil || claims.MerID == 0 {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	pair, err := h.jwt.Issue(authjwt.PortalMerchant, claims.AdminID, claims.MerID, 0, claims.Account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair})
}

func (h *Handler) Me(c *gin.Context) {
	profile, err := h.svc.MerchantProfile(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		writeIdentityErr(c, err)
		return
	}
	if profile.MerID != middleware.MerID(c) {
		response.Fail(c, http.StatusForbidden, "商户上下文不匹配")
		return
	}
	response.OK(c, profile)
}

func (h *Handler) Menus(c *gin.Context) {
	profile, err := h.svc.MerchantProfile(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		writeIdentityErr(c, err)
		return
	}
	if profile.MerID != middleware.MerID(c) {
		response.Fail(c, http.StatusForbidden, "商户上下文不匹配")
		return
	}
	menus, err := h.svc.MenusForMerchant(c.Request.Context(), profile.Roles, profile.Level)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载菜单失败")
		return
	}
	response.OK(c, gin.H{"menus": menus, "mer_id": profile.MerID})
}

func (h *Handler) Permissions(c *gin.Context) {
	profile, err := h.svc.MerchantProfile(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		writeIdentityErr(c, err)
		return
	}
	if profile.MerID != middleware.MerID(c) {
		response.Fail(c, http.StatusForbidden, "商户上下文不匹配")
		return
	}
	perms, err := h.svc.MerchantPermissionPaths(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载权限失败")
		return
	}
	response.OK(c, gin.H{"permissions": perms})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req passwordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	err := h.svc.ChangeMerchantPassword(c.Request.Context(), middleware.AdminID(c), req.OldPassword, req.NewPassword)
	if err != nil {
		writeIdentityErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
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
	case errors.Is(err, identity.ErrWeakPassword):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
