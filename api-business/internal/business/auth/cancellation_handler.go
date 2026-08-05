package auth

import (
	"errors"
	"net/http"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type cancelAccountRequest struct {
	ConfirmationToken string `json:"confirmation_token"`
}

func (h *Handler) RegisterCancellation(r gin.IRoutes) { r.POST("/auth/cancel", h.CancelAccount) }

func (h *Handler) CancelAccount(c *gin.Context) {
	var req cancelAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "注销参数错误")
		return
	}
	result, err := h.svc.CancelAccount(c.Request.Context(), uint64(middleware.UID(c)), req.ConfirmationToken)
	if err != nil {
		if errors.Is(err, ErrCancellationConfirmation) {
			response.Fail(c, http.StatusConflict, err.Error())
			return
		}
		writeError(c, err)
		return
	}
	response.OK(c, result)
}
