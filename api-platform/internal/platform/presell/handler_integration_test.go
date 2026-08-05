package presell

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	domain "github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/presell"
	persist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/presell"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestPresellHTTPRBACAndStatusOnlyUpdate(t *testing.T) {
	adminDSN, businessDSN := strings.TrimSpace(os.Getenv("ECRM_PRESELL_ADMIN_TEST_DSN")), strings.TrimSpace(os.Getenv("ECRM_PRESELL_BUSINESS_TEST_DSN"))
	if adminDSN == "" || businessDSN == "" {
		t.Skip("set ECRM_PRESELL_ADMIN_TEST_DSN and ECRM_PRESELL_BUSINESS_TEST_DSN")
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
	for _, r := range dbRoles {
		roleID[r.Code] = r.ID
	}
	fixtures := []struct {
		ID   uint
		Role string
	}{{987677101, "platform"}, {987677102, "merchant"}, {987677103, "region"}, {987677104, "operations"}, {987677105, "customer_service"}}
	userIDs := []uint{987677101, 987677102, 987677103, 987677104, 987677105}
	const presellID uint64 = 987677001
	cleanup := func() {
		_ = businessDB.Table("qixi_crm_b_presell").Where("product_presell_id=?", presellID).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", userIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user").Where("id IN ?", userIDs).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	for _, f := range fixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{"id": f.ID, "username": "presell-rbac-" + f.Role, "password_hash": "not-used", "display_name": "中文预售验收-" + f.Role, "status": 1, "auth_version": 1, "data_scope_version": 1}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": f.ID, "role_id": roleID[f.Role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now()
	if err := businessDB.Table("qixi_crm_b_presell").Create(map[string]any{"product_presell_id": presellID, "start_time": now.Add(-time.Hour), "end_time": now.Add(time.Hour), "status": 1, "presell_type": 1, "product_id": 987677901, "price": "299.00", "stock": 12, "is_show": 1, "store_name": "中文预售快照验收商品", "mer_id": 1, "store_info": "中文模拟预售活动。", "is_del": 0, "product_status": 1, "action_status": 1}).Error; err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager(strings.Repeat("x", 32), time.Minute, 2*time.Minute)
	router := gin.New()
	group := router.Group("/api/platform/v1")
	group.Use(middleware.JWTRequired(jwt, authjwt.PortalPlatform), middleware.RequireAdminConsole(), middleware.RequireAdminSession(adminDB), middleware.RestrictRoleConsole(), middleware.RestrictRegionConsole())
	NewHandler(domain.NewService(persist.NewRepo(businessDB)), adminDB).Register(group)
	call := func(f struct {
		ID   uint
		Role string
	}, method, path string, body any) *httptest.ResponseRecorder {
		var raw []byte
		if body != nil {
			raw, _ = json.Marshal(body)
		}
		pair, e := jwt.IssueAdminConsole(f.ID, "presell-rbac-"+f.Role, []string{f.Role}, 1)
		if e != nil {
			t.Fatal(e)
		}
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	if got := call(fixtures[3], http.MethodGet, "/api/platform/v1/presell/actives", nil); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "中文预售快照验收商品") {
		t.Fatalf("operations list code=%d body=%s", got.Code, got.Body.String())
	}
	for _, f := range []struct {
		ID   uint
		Role string
	}{fixtures[1], fixtures[2], fixtures[4]} {
		if got := call(f, http.MethodGet, "/api/platform/v1/presell/actives", nil).Code; got != http.StatusForbidden {
			t.Fatalf("%s list=%d", f.Role, got)
		}
	}
	if got := call(fixtures[0], http.MethodPut, "/api/platform/v1/presell/actives/987677001", gin.H{"price": 1}).Code; got != http.StatusBadRequest {
		t.Fatalf("price-only=%d", got)
	}
	if got := call(fixtures[3], http.MethodPut, "/api/platform/v1/presell/actives/987677001", gin.H{"status": 0, "price": 1, "stock": 1}).Code; got != http.StatusOK {
		t.Fatalf("operations stop=%d", got)
	}
	var row struct {
		Status int
		Price  float64
		Stock  int
	}
	if err := businessDB.Table("qixi_crm_b_presell").Select("status,price,stock").Where("product_presell_id=?", presellID).Take(&row).Error; err != nil || row.Status != 0 || row.Price != 299 || row.Stock != 12 {
		t.Fatalf("row=%#v err=%v", row, err)
	}
	if _, err := domain.NewService(persist.NewRepo(businessDB)).Quote(context.Background(), uint(presellID)); err == nil {
		t.Fatal("stopped presell must reject quote")
	}
}
