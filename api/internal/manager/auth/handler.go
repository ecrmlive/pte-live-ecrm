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

func (h *Handler) Register(public, authed gin.IRoutes) {
	public.POST("/auth/login", h.Login)
	public.POST("/auth/refresh", h.Refresh)
	authed.GET("/auth/me", h.Me)
}

type loginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
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
	staff, mer, err := h.svc.LoginStoreService(c.Request.Context(), req.Account, req.Password)
	if err != nil {
		writeErr(c, err)
		return
	}
	pair, err := h.jwt.Issue(authjwt.PortalManager, staff.ServiceID, staff.MerID, 0, staff.Account)
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
		},
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	claims, err := h.jwt.ParseExpect(req.RefreshToken, authjwt.PortalManager, authjwt.TokenRefresh)
	if err != nil || claims.MerID == 0 || claims.AdminID == 0 {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	pair, err := h.jwt.Issue(authjwt.PortalManager, claims.AdminID, claims.MerID, 0, claims.Account)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair})
}

func (h *Handler) Me(c *gin.Context) {
	profile, err := h.svc.StoreServiceProfile(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	if profile.MerID != middleware.MerID(c) {
		response.Fail(c, http.StatusForbidden, "商户上下文不匹配")
		return
	}
	response.OK(c, profile)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, identity.ErrInvalidCredentials):
		response.Fail(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, identity.ErrAccountDisabled),
		errors.Is(err, identity.ErrMerchantDisabled),
		errors.Is(err, identity.ErrNoVerifyPerm):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, identity.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
