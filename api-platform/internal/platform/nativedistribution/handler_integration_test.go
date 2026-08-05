package nativedistribution

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestDistributionHTTPRBACAndReadModel(t *testing.T) {
	adminDSN, businessDSN := strings.TrimSpace(os.Getenv("ECRM_DISTRIBUTION_ADMIN_TEST_DSN")), strings.TrimSpace(os.Getenv("ECRM_DISTRIBUTION_BUSINESS_TEST_DSN"))
	if adminDSN == "" || businessDSN == "" {
		t.Skip("set ECRM_DISTRIBUTION_ADMIN_TEST_DSN and ECRM_DISTRIBUTION_BUSINESS_TEST_DSN")
	}
	adminDB, err := gorm.Open(mysql.Open(adminDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	businessDB, err := gorm.Open(mysql.Open(businessDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	roles := []string{"platform", "merchant", "region", "operations", "customer_service"}
	var dbRoles []struct {
		ID   uint
		Code string
	}
	if err := adminDB.Table("qixi_crm_a_role").Select("id,code").Where("code IN ? AND status=1", roles).Find(&dbRoles).Error; err != nil || len(dbRoles) != len(roles) {
		t.Fatalf("roles=%v err=%v", dbRoles, err)
	}
	roleID := map[string]uint{}
	for _, role := range dbRoles {
		roleID[role.Code] = role.ID
	}
	fixtures := []struct {
		ID   uint
		Role string
	}{{987677301, "platform"}, {987677302, "merchant"}, {987677303, "region"}, {987677304, "operations"}, {987677305, "customer_service"}}
	userIDs := []uint{987677301, 987677302, 987677303, 987677304, 987677305}
	const promoterID uint64 = 987677301
	cleanup := func() {
		_ = businessDB.Table("qixi_crm_b_commission_ledger").Where("idempotency_key IN ?", []string{"distribution-rbac-pending", "distribution-rbac-available", "distribution-rbac-settled"}).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_distribution_relation").Where("user_id=?", promoterID+1).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_distribution_promoter").Where("user_id=?", promoterID).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", userIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user").Where("id IN ?", userIDs).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	for _, fixture := range fixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{"id": fixture.ID, "username": "distribution-rbac-" + fixture.Role, "password_hash": "not-used", "display_name": "中文分销验收-" + fixture.Role, "status": 1, "auth_version": 1, "data_scope_version": 1}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": fixture.ID, "role_id": roleID[fixture.Role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now()
	if err := businessDB.Table("qixi_crm_b_distribution_promoter").Create(map[string]any{"user_id": promoterID, "status": 1, "updated_at": now}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_distribution_relation").Create(map[string]any{"user_id": promoterID + 1, "parent_user_id": promoterID, "bound_at": now}).Error; err != nil {
		t.Fatal(err)
	}
	entries := []map[string]any{
		{"user_id": promoterID, "order_id": 0, "amount": "6.50", "status": "pending", "idempotency_key": "distribution-rbac-pending", "created_at": now},
		{"user_id": promoterID, "order_id": 0, "amount": "12.80", "status": "available", "idempotency_key": "distribution-rbac-available", "available_at": now, "created_at": now},
		{"user_id": promoterID, "order_id": 0, "amount": "3.20", "status": "settled", "idempotency_key": "distribution-rbac-settled", "available_at": now, "created_at": now},
	}
	if err := businessDB.Table("qixi_crm_b_commission_ledger").Create(&entries).Error; err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager(strings.Repeat("x", 32), time.Minute, 2*time.Minute)
	router := gin.New()
	group := router.Group("/api/platform/v1")
	group.Use(middleware.JWTRequired(jwt, authjwt.PortalPlatform), middleware.RequireAdminConsole(), middleware.RequireAdminSession(adminDB), middleware.RestrictRoleConsole(), middleware.RestrictRegionConsole())
	NewHandler(businessDB, adminDB).Register(group)
	call := func(fixture struct {
		ID   uint
		Role string
	}, path string) *httptest.ResponseRecorder {
		pair, issueErr := jwt.IssueAdminConsole(fixture.ID, "distribution-rbac-"+fixture.Role, []string{fixture.Role}, 1)
		if issueErr != nil {
			t.Fatal(issueErr)
		}
		req := httptest.NewRequest(http.MethodGet, path, bytes.NewReader(nil))
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	if got := call(fixtures[3], "/api/platform/v1/distribution/promoters?user_id=987677301"); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), `"direct_user_count":1`) || !strings.Contains(got.Body.String(), `"available_commission":12.8`) {
		t.Fatalf("operations promoter code=%d body=%s", got.Code, got.Body.String())
	}
	if got := call(fixtures[0], "/api/platform/v1/distribution/commissions?user_id=987677301&status=available"); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), `"amount":12.8`) {
		t.Fatalf("platform commission code=%d body=%s", got.Code, got.Body.String())
	}
	if got := call(fixtures[0], "/api/platform/v1/distribution/summary"); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), `"pending_commission":`) || !strings.Contains(got.Body.String(), `"settled_commission":`) {
		t.Fatalf("summary code=%d body=%s", got.Code, got.Body.String())
	}
	if got := call(fixtures[0], "/api/platform/v1/distribution/commissions?status=已结算").Code; got != http.StatusBadRequest {
		t.Fatalf("invalid status=%d", got)
	}
	for _, fixture := range []struct {
		ID   uint
		Role string
	}{fixtures[1], fixtures[2], fixtures[4]} {
		if got := call(fixture, "/api/platform/v1/distribution/summary").Code; got != http.StatusForbidden {
			t.Fatalf("%s summary=%d", fixture.Role, got)
		}
	}
}
