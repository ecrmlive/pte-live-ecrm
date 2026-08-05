package auth

import (
	"net/http"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

func (h *Handler) RegisterPassword(r gin.IRoutes) {
	r.POST("/auth/password", h.ChangePassword)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.NewPassword != req.ConfirmPassword {
		response.Fail(c, http.StatusBadRequest, "新密码确认不一致")
		return
	}
	claims := middleware.ClaimsFrom(c)
	channel, err := ParseChannel(claims.ClientPlatform)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "登录来源无效")
		return
	}
	profile, err := h.svc.ChangePassword(c.Request.Context(), uint64(middleware.UID(c)), channel, req.CurrentPassword, req.NewPassword)
	if err != nil {
		writeError(c, err)
		return
	}
	pair, err := h.issueForContext(c, profile, claims)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair, "user": profile})
}
