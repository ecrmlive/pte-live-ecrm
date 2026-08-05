package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestChangePasswordRejectsMismatchedConfirmationBeforeService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil, nil, nil).RegisterPassword(router)
	request := httptest.NewRequest(http.MethodPost, "/auth/password", strings.NewReader(`{"current_password":"旧密码1234","new_password":"新密码1234","confirm_password":"新密码5678"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}
