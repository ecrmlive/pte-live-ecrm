package searchhistory

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNormalizeKeepsChineseKeyword(t *testing.T) {
	value, ok := normalize("  夏日防晒衣  ")
	if !ok || value != "夏日防晒衣" {
		t.Fatalf("normalize() = %q, %v", value, ok)
	}
	for _, raw := range []string{"", "   ", strings.Repeat("测", 129)} {
		if _, ok := normalize(raw); ok {
			t.Fatalf("%q should be rejected", raw)
		}
	}
}

func TestSearchHistoryRoutesRejectInvalidInputBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil).Register(router)

	invalidKeyword := httptest.NewRecorder()
	router.ServeHTTP(invalidKeyword, httptest.NewRequest(http.MethodPost, "/search-history", strings.NewReader(`{"keyword":"   "}`)))
	if invalidKeyword.Code != http.StatusBadRequest {
		t.Fatalf("invalid keyword status = %d; want %d", invalidKeyword.Code, http.StatusBadRequest)
	}

	invalidID := httptest.NewRecorder()
	router.ServeHTTP(invalidID, httptest.NewRequest(http.MethodDelete, "/search-history/0", nil))
	if invalidID.Code != http.StatusBadRequest {
		t.Fatalf("invalid id status = %d; want %d", invalidID.Code, http.StatusBadRequest)
	}
}
