package recharge

import (
	"bytes"
	"encoding/json"
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

func TestRechargePlanHTTPRBACVersionAndOrderSnapshot(t *testing.T) {
	adminDSN, businessDSN := strings.TrimSpace(os.Getenv("ECRM_RECHARGE_ADMIN_TEST_DSN")), strings.TrimSpace(os.Getenv("ECRM_RECHARGE_BUSINESS_TEST_DSN"))
	if adminDSN == "" || businessDSN == "" {
		t.Skip("set ECRM_RECHARGE_ADMIN_TEST_DSN and ECRM_RECHARGE_BUSINESS_TEST_DSN to run recharge integration test")
	}
	adminDB, err := gorm.Open(mysql.Open(adminDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	businessDB, err := gorm.Open(mysql.Open(businessDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	var versionColumnCount int64
	if err := businessDB.Raw("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_recharge_plan' AND COLUMN_NAME = 'version'").Scan(&versionColumnCount).Error; err != nil {
		t.Fatal(err)
	}
	if versionColumnCount == 0 {
		if err := businessDB.Exec("ALTER TABLE qixi_crm_b_recharge_plan ADD COLUMN version bigint unsigned NOT NULL DEFAULT 1 AFTER sort").Error; err != nil {
			t.Fatal(err)
		}
	}
	roles := []string{"platform", "merchant", "region", "operations", "customer_service"}
	var roleRows []struct {
		ID   uint
		Code string
	}
	if err := adminDB.Table("qixi_crm_a_role").Select("id,code").Where("code IN ? AND status=1", roles).Find(&roleRows).Error; err != nil {
		t.Fatal(err)
	}
	roleIDs := map[string]uint{}
	for _, item := range roleRows {
		roleIDs[item.Code] = item.ID
	}
	if len(roleIDs) != len(roles) {
		t.Fatalf("roles missing: %#v", roleIDs)
	}
	const planID uint = 987674001
	const orderID uint = 987674011
	userIDs := []uint{987674101, 987674102, 987674103, 987674104, 987674105}
	cleanup := func() {
		_ = businessDB.Table("qixi_crm_b_recharge_order").Where("id=?", orderID).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_recharge_plan").Where("id=?", planID).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", userIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user").Where("id IN ?", userIDs).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	fixtures := []struct {
		ID   uint
		Role string
	}{{userIDs[0], "platform"}, {userIDs[1], "merchant"}, {userIDs[2], "region"}, {userIDs[3], "operations"}, {userIDs[4], "customer_service"}}
	for _, item := range fixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{"id": item.ID, "username": "recharge-rbac-" + item.Role, "password_hash": "not-used", "display_name": "中文充值权限验收-" + item.Role, "status": 1, "auth_version": 1, "data_scope_version": 1}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": item.ID, "role_id": roleIDs[item.Role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := businessDB.Table("qixi_crm_b_recharge_plan").Create(map[string]any{"id": planID, "name": "中文并发充值计划", "amount": "100.00", "bonus_amount": "8.00", "status": 1, "sort": 1, "version": 1}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_recharge_order").Create(map[string]any{"id": orderID, "recharge_no": "R-ACCEPTANCE-SNAPSHOT-001", "user_id": 987674901, "amount": "100.00", "bonus_amount": "8.00", "status": "pending", "idempotency_key": "recharge-acceptance-snapshot-001"}).Error; err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	jwtManager := authjwt.NewManager(strings.Repeat("x", 32), time.Minute, 2*time.Minute)
	r := gin.New()
	authed := r.Group("/api/platform/v1")
	authed.Use(middleware.JWTRequired(jwtManager, authjwt.PortalPlatform), middleware.RequireAdminConsole(), middleware.RequireAdminSession(adminDB), middleware.RestrictRoleConsole(), middleware.RestrictRegionConsole())
	NewHandler(businessDB, adminDB).Register(authed)
	call := func(item struct {
		ID   uint
		Role string
	}, method, path string, body any) int {
		raw, _ := json.Marshal(body)
		pair, issueErr := jwtManager.IssueAdminConsole(item.ID, "recharge-rbac-"+item.Role, []string{item.Role}, 1)
		if issueErr != nil {
			t.Fatal(issueErr)
		}
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}
	invalidMoney := gin.H{"name": "中文并发充值计划", "amount": 100.001, "bonus_amount": 8, "status": 1, "sort": 1, "version": 1}
	if got := call(fixtures[3], http.MethodPut, "/api/platform/v1/recharge/plans/987674001", invalidMoney); got != http.StatusBadRequest {
		t.Fatalf("invalid money update=%d, want 400", got)
	}
	update := gin.H{"name": "中文并发充值计划（运营已更新）", "amount": 200, "bonus_amount": 10, "status": 0, "sort": 2, "version": 1}
	if got := call(fixtures[3], http.MethodPut, "/api/platform/v1/recharge/plans/987674001", update); got != http.StatusOK {
		t.Fatalf("operations update=%d", got)
	}
	if got := call(fixtures[0], http.MethodPut, "/api/platform/v1/recharge/plans/987674001", update); got != http.StatusConflict {
		t.Fatalf("stale platform update=%d, want 409", got)
	}
	for _, item := range []struct {
		ID   uint
		Role string
	}{fixtures[1], fixtures[2], fixtures[4]} {
		if got := call(item, http.MethodPut, "/api/platform/v1/recharge/plans/987674001", update); got != http.StatusForbidden {
			t.Fatalf("%s update=%d", item.Role, got)
		}
	}
	var planRow struct {
		Name    string `gorm:"column:name"`
		Amount  string `gorm:"column:amount"`
		Bonus   string `gorm:"column:bonus_amount"`
		Status  int    `gorm:"column:status"`
		Version uint64 `gorm:"column:version"`
	}
	if err := businessDB.Table("qixi_crm_b_recharge_plan").Select("name,amount,bonus_amount,status,version").Where("id=?", planID).Take(&planRow).Error; err != nil || planRow.Name != "中文并发充值计划（运营已更新）" || planRow.Amount != "200.00" || planRow.Bonus != "10.00" || planRow.Status != 0 || planRow.Version != 2 {
		t.Fatalf("plan=%#v err=%v", planRow, err)
	}
	var orderRow struct {
		Amount string `gorm:"column:amount"`
		Bonus  string `gorm:"column:bonus_amount"`
		Status string `gorm:"column:status"`
	}
	if err := businessDB.Table("qixi_crm_b_recharge_order").Select("amount,bonus_amount,status").Where("id=?", orderID).Take(&orderRow).Error; err != nil || orderRow.Amount != "100.00" || orderRow.Bonus != "8.00" || orderRow.Status != "pending" {
		t.Fatalf("order snapshot=%#v err=%v", orderRow, err)
	}
}
