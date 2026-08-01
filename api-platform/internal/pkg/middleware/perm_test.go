package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
)

func TestRequireAdminConsole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaimsKey, &authjwt.Claims{AdminID: 8, PrincipalID: 8, PrincipalType: authjwt.PrincipalAdminUser, Scope: authjwt.ScopeAdminConsole, Roles: []string{"operations"}})
	})
	r.GET("/ok", RequireAdminConsole(), func(c *gin.Context) { c.Status(http.StatusNoContent) })

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusNoContent)
	}
}

func TestRestrictRegionConsoleAllowsMigratedOrderRead(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaimsKey, &authjwt.Claims{AdminID: 9, PrincipalID: 9, PrincipalType: authjwt.PrincipalAdminUser, Scope: authjwt.ScopeAdminConsole, Roles: []string{"region"}})
	})
	r.GET("/api/platform/v1/orders", RestrictRegionConsole(), func(c *gin.Context) { c.Status(http.StatusNoContent) })

	req := httptest.NewRequest(http.MethodGet, "/api/platform/v1/orders", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusNoContent)
	}
}

func TestRequireStoreConsole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaimsKey, &authjwt.Claims{
			AdminID: 4, PrincipalID: 4, PrincipalType: authjwt.PrincipalStoreAccount, MerchantID: 5, MerID: 5, StoreID: 6, StoreAppID: "qixi.store.demo.1", MerchantAppID: "qixi.store.demo.1",
			Scope: authjwt.ScopeStoreConsole, Roles: []string{"owner"},
		})
	})
	r.GET("/store", RequireStoreConsole(), func(c *gin.Context) { c.Status(http.StatusNoContent) })

	req := httptest.NewRequest(http.MethodGet, "/store", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusNoContent)
	}
}

func TestRequireStoreAppID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaimsKey, &authjwt.Claims{MerchantAppID: "qixi.store.demo.1"})
	})
	r.GET("/store", RequireStoreAppID(), func(c *gin.Context) { c.Status(http.StatusNoContent) })

	for _, tc := range []struct {
		name   string
		appID  string
		status int
	}{
		{name: "matched", appID: "qixi.store.demo.1", status: http.StatusNoContent},
		{name: "missing", status: http.StatusForbidden},
		{name: "mismatched", appID: "qixi.store.other", status: http.StatusForbidden},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/store", nil)
			if tc.appID != "" {
				req.Header.Set("X-AppId", tc.appID)
			}
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			if resp.Code != tc.status {
				t.Fatalf("status = %d, want %d", resp.Code, tc.status)
			}
		})
	}
}
