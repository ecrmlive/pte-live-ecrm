package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/authjwt"
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

func TestJWTOptionalAddsUserContextOnlyForValidAuthoriZation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := authjwt.NewManager("test-secret", time.Hour, time.Hour)
	pair, err := mgr.IssueCUser(42, "13500000001", "pc")
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	r := gin.New()
	r.GET("/optional", JWTOptional(mgr, authjwt.PortalApp), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"uid": UID(c)})
	})

	guest := httptest.NewRecorder()
	r.ServeHTTP(guest, httptest.NewRequest(http.MethodGet, "/optional", nil))
	if guest.Code != http.StatusOK || !strings.Contains(guest.Body.String(), `"uid":0`) {
		t.Fatalf("guest response = %d %s", guest.Code, guest.Body.String())
	}

	request := httptest.NewRequest(http.MethodGet, "/optional", nil)
	request.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
	authed := httptest.NewRecorder()
	r.ServeHTTP(authed, request)
	if authed.Code != http.StatusOK || !strings.Contains(authed.Body.String(), `"uid":42`) {
		t.Fatalf("authed response = %d %s", authed.Code, authed.Body.String())
	}
}
