package merchantdeposit

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/profitsharing"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const financeRBACAcceptanceUserBase uint = 987680100

func TestMerchantFinanceHTTPRBACAcrossFiveRoles(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_FINANCE_RBAC_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_FINANCE_RBAC_TEST_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	roles := []string{"platform", "merchant", "region", "operations", "customer_service"}
	roleIDs := make(map[string]uint)
	type roleRecord struct {
		ID   uint
		Code string
	}
	var storedRoles []roleRecord
	if err := db.Table("qixi_crm_a_role").Select("id, code").Where("code IN ? AND status = 1", roles).Find(&storedRoles).Error; err != nil {
		t.Fatal(err)
	}
	for _, item := range storedRoles {
		roleIDs[item.Code] = item.ID
	}
	if len(roleIDs) != len(roles) {
		t.Fatalf("five role fixture missing: %#v", roleIDs)
	}
	userIDs := make([]uint, 0, len(roles))
	for index := range roles {
		userIDs = append(userIDs, financeRBACAcceptanceUserBase+uint(index))
	}
	cleanup := func() {
		db.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", userIDs).Delete(nil)
		db.Table("qixi_crm_a_admin_user").Where("id IN ?", userIDs).Delete(nil)
	}
	cleanup()
	defer cleanup()
	for index, role := range roles {
		userID := userIDs[index]
		if err := db.Table("qixi_crm_a_admin_user").Create(map[string]any{
			"id":                 userID,
			"username":           "finance-rbac-" + role,
			"password_hash":      "not-used-by-rbac-test",
			"display_name":       "中文财务权限验收-" + role,
			"status":             1,
			"auth_version":       1,
			"data_scope_version": 1,
		}).Error; err != nil {
			t.Fatal(err)
		}
		if err := db.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": userID, "role_id": roleIDs[role]}).Error; err != nil {
			t.Fatal(err)
		}
	}

	gin.SetMode(gin.TestMode)
	jwtManager := authjwt.NewManager(strings.Repeat("x", 32), time.Minute, 2*time.Minute)
	r := gin.New()
	authed := r.Group("/api/platform/v1")
	authed.Use(
		middleware.JWTRequired(jwtManager, authjwt.PortalPlatform),
		middleware.RequireAdminConsole(),
		middleware.RequireAdminSession(db),
		middleware.RestrictRoleConsole(),
		middleware.RestrictRegionConsole(),
	)
	NewHandler(db).Register(authed)
	profitsharing.NewHandler(db).Register(authed)

	call := func(path, role string, userID uint) int {
		pair, issueErr := jwtManager.IssueAdminConsole(userID, "finance-rbac-"+role, []string{role}, 1)
		if issueErr != nil {
			t.Fatal(issueErr)
		}
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}
	for _, path := range []string{"/api/platform/v1/merchant-deposits", "/api/platform/v1/merchant-profitsharing-applications"} {
		if got := call(path, "platform", userIDs[0]); got != http.StatusOK {
			t.Fatalf("platform %s=%d", path, got)
		}
		for index, role := range roles[1:] {
			if got := call(path, role, userIDs[index+1]); got != http.StatusForbidden {
				t.Fatalf("%s %s=%d, want 403", role, path, got)
			}
		}
	}
	unauthenticated := httptest.NewRequest(http.MethodGet, "/api/platform/v1/merchant-deposits", nil)
	unauthenticatedRecorder := httptest.NewRecorder()
	r.ServeHTTP(unauthenticatedRecorder, unauthenticated)
	if unauthenticatedRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated=%d", unauthenticatedRecorder.Code)
	}
}
