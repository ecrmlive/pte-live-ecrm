package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/gin-gonic/gin"
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

func TestMatchesAdminSessionRejectsStaleVersions(t *testing.T) {
	claims := &authjwt.Claims{IdentityVersion: 3, DataScopeVersion: 8}
	if !matchesAdminSession(claims, adminSessionVersion{Status: 1, AuthVersion: 3, DataScopeVersion: 8}) {
		t.Fatal("matching enabled session should pass")
	}
	for _, session := range []adminSessionVersion{
		{Status: 0, AuthVersion: 3, DataScopeVersion: 8},
		{Status: 1, AuthVersion: 4, DataScopeVersion: 8},
		{Status: 1, AuthVersion: 3, DataScopeVersion: 9},
	} {
		if matchesAdminSession(claims, session) {
			t.Fatalf("stale session should fail: %#v", session)
		}
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

func TestRestrictRegionConsoleAllowsScopedMerchantSettlementRead(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaimsKey, &authjwt.Claims{AdminID: 9, PrincipalID: 9, PrincipalType: authjwt.PrincipalAdminUser, Scope: authjwt.ScopeAdminConsole, Roles: []string{"region"}})
	})
	r.GET("/api/platform/v1/finance/merchant-settlements", RestrictRegionConsole(), func(c *gin.Context) { c.Status(http.StatusNoContent) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/platform/v1/finance/merchant-settlements", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestRestrictRegionConsoleAllowsScopedRefundRead(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaimsKey, &authjwt.Claims{AdminID: 9, PrincipalID: 9, PrincipalType: authjwt.PrincipalAdminUser, Scope: authjwt.ScopeAdminConsole, Roles: []string{"region"}})
	})
	r.GET("/api/platform/v1/refunds", RestrictRegionConsole(), func(c *gin.Context) { c.Status(http.StatusNoContent) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/platform/v1/refunds", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestRestrictRoleConsoleUsesDenyByDefaultForNonPlatformRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name   string
		roles  []string
		path   string
		status int
	}{
		{name: "platform can supervise finance", roles: []string{"platform"}, path: "/api/platform/v1/finance/withdraws", status: http.StatusNoContent},
		{name: "customer service can use own queue", roles: []string{"customer_service"}, path: "/api/platform/v1/customer-service/threads", status: http.StatusNoContent},
		{name: "customer service can handle user feedback", roles: []string{"customer_service"}, path: "/api/platform/v1/user-feedback", status: http.StatusNoContent},
		{name: "customer service cannot supervise finance", roles: []string{"customer_service"}, path: "/api/platform/v1/finance/withdraws", status: http.StatusForbidden},
		{name: "merchant can read scoped orders", roles: []string{"merchant"}, path: "/api/platform/v1/orders", status: http.StatusNoContent},
		{name: "merchant can read scoped products", roles: []string{"merchant"}, path: "/api/platform/v1/products", status: http.StatusNoContent},
		{name: "merchant cannot call diy", roles: []string{"merchant"}, path: "/api/platform/v1/diy/pages", status: http.StatusForbidden},
		{name: "operations can call article", roles: []string{"operations"}, path: "/api/platform/v1/articles", status: http.StatusNoContent},
		{name: "operations can supervise assist", roles: []string{"operations"}, path: "/api/platform/v1/assist/actives", status: http.StatusNoContent},
		{name: "operations can supervise points mall", roles: []string{"operations"}, path: "/api/platform/v1/points/products", status: http.StatusNoContent},
		{name: "operations can supervise recharge plans", roles: []string{"operations"}, path: "/api/platform/v1/recharge/plans", status: http.StatusNoContent},
		{name: "merchant cannot supervise points mall", roles: []string{"merchant"}, path: "/api/platform/v1/points/products", status: http.StatusForbidden},
		{name: "merchant cannot supervise recharge plans", roles: []string{"merchant"}, path: "/api/platform/v1/recharge/plans", status: http.StatusForbidden},
		{name: "operations can call community moderation", roles: []string{"operations"}, path: "/api/platform/v1/community/posts", status: http.StatusNoContent},
		{name: "operations can call attachment library", roles: []string{"operations"}, path: "/api/platform/v1/attachments", status: http.StatusNoContent},
		{name: "operations can call broadcast moderation", roles: []string{"operations"}, path: "/api/platform/v1/broadcast/rooms", status: http.StatusNoContent},
		{name: "operations can read distribution supervision", roles: []string{"operations"}, path: "/api/platform/v1/distribution/promoters", status: http.StatusNoContent},
		{name: "operations cannot alter admin settings", roles: []string{"operations"}, path: "/api/platform/v1/setting/admins", status: http.StatusForbidden},
		{name: "operations cannot manage SVIP users", roles: []string{"operations"}, path: "/api/platform/v1/users/1/svip", status: http.StatusForbidden},
		{name: "customer service cannot manage SVIP users", roles: []string{"customer_service"}, path: "/api/platform/v1/users/1/svip", status: http.StatusForbidden},
		{name: "all roles can refresh own console profile", roles: []string{"customer_service"}, path: "/api/platform/v1/auth/me", status: http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set(CtxClaimsKey, &authjwt.Claims{AdminID: 9, PrincipalID: 9, PrincipalType: authjwt.PrincipalAdminUser, Scope: authjwt.ScopeAdminConsole, Roles: tc.roles})
			})
			r.GET(tc.path, RestrictRoleConsole(), func(c *gin.Context) { c.Status(http.StatusNoContent) })
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if w.Code != tc.status {
				t.Fatalf("status = %d, want %d", w.Code, tc.status)
			}
		})
	}
}

func TestRestrictRegionConsoleDoesNotLimitPlatformRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaimsKey, &authjwt.Claims{AdminID: 9, PrincipalID: 9, PrincipalType: authjwt.PrincipalAdminUser, Scope: authjwt.ScopeAdminConsole, Roles: []string{"platform", "region"}})
	})
	r.GET("/api/platform/v1/finance/withdraws", RestrictRegionConsole(), func(c *gin.Context) { c.Status(http.StatusNoContent) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/platform/v1/finance/withdraws", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
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

func TestRequireAdminRolesRejectsNonPlatformRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaimsKey, &authjwt.Claims{Roles: []string{"operations"}})
	})
	r.GET("/setting/admins", RequireAdminRoles("platform"), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/setting/admins", nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestRequireAdminRolesAllowsPlatformRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaimsKey, &authjwt.Claims{Roles: []string{"platform"}})
	})
	r.GET("/setting/admins", RequireAdminRoles("platform"), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/setting/admins", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestRequireAdminMenuFailsClosedWithoutUnifiedDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaimsKey, &authjwt.Claims{AdminID: 9, PrincipalID: 9, PrincipalType: authjwt.PrincipalAdminUser, Scope: authjwt.ScopeAdminConsole, Roles: []string{"platform"}})
	})
	r.POST("/products/1/audit", RequireAdminMenu(nil, "product.audit.submit"), func(c *gin.Context) { c.Status(http.StatusNoContent) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/products/1/audit", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
