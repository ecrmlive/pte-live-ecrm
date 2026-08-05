package city

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(nil).Register(r)
	for _, path := range []string{"/system/city/lst", "/system/city/lst/:pid"} {
		found := false
		for _, route := range r.Routes() {
			if route.Method == http.MethodGet && route.Path == path {
				found = true
			}
		}
		if !found {
			t.Fatalf("route %s is not registered", path)
		}
	}
}

func TestListRejectsInvalidParentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(nil).Register(r)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/system/city/lst/not-a-number", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
