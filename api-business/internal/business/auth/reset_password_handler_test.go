package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResetPasswordRejectsMismatchedConfirmationBeforeConsumingCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil, nil, nil).Register(router, router)

	req := httptest.NewRequest(http.MethodPost, "/auth/password/reset", strings.NewReader(`{"mobile":"13800000000","code":"123456","new_password":"安全密码1234","confirm_password":"安全密码5678"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", resp.Code, resp.Body.String())
	}
}
