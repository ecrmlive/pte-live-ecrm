package reservation

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCheckRejectsInvalidInputBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil).Register(router)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/order/reservation/check", strings.NewReader(`{"slot_id":1,"date":"2026-08-03"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestCreateRequiresIdempotencyKeyBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil).Register(router)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/order/reservation/create", strings.NewReader(`{"product_id":1001,"slot_id":11,"date":"2026-08-03","address_id":1}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
