package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/authjwt"
)

func TestJWTRequiredOnlyAcceptsAuthoriZation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := authjwt.NewManager("test-secret", time.Hour, 2*time.Hour)
	pair, err := mgr.IssueCUser(7, "user", "h5")
	if err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	r.GET("/protected", JWTRequired(mgr, authjwt.PortalApp), func(c *gin.Context) { c.Status(http.StatusNoContent) })

	wrong := httptest.NewRequest(http.MethodGet, "/protected", nil)
	wrong.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	wrongResp := httptest.NewRecorder()
	r.ServeHTTP(wrongResp, wrong)
	if wrongResp.Code != http.StatusUnauthorized {
		t.Fatalf("Authorization status = %d, want %d", wrongResp.Code, http.StatusUnauthorized)
	}

	valid := httptest.NewRequest(http.MethodGet, "/protected", nil)
	valid.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
	validResp := httptest.NewRecorder()
	r.ServeHTTP(validResp, valid)
	if validResp.Code != http.StatusNoContent {
		t.Fatalf("Authori-zation status = %d, want %d", validResp.Code, http.StatusNoContent)
	}
}
