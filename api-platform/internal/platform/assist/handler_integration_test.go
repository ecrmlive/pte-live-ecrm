package assist

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

	domain "github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/assist"
	persist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/assist"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAssistHTTPRBACAndVisibilityOnlyUpdate(t *testing.T) {
	adminDSN, businessDSN := strings.TrimSpace(os.Getenv("ECRM_ASSIST_ADMIN_TEST_DSN")), strings.TrimSpace(os.Getenv("ECRM_ASSIST_BUSINESS_TEST_DSN"))
	if adminDSN == "" || businessDSN == "" {
		t.Skip("set ECRM_ASSIST_ADMIN_TEST_DSN and ECRM_ASSIST_BUSINESS_TEST_DSN")
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
	}{{987677201, "platform"}, {987677202, "merchant"}, {987677203, "region"}, {987677204, "operations"}, {987677205, "customer_service"}}
	userIDs := []uint{987677201, 987677202, 987677203, 987677204, 987677205}
	const assistID uint = 987677201
	cleanup := func() {
		_ = businessDB.Table("qixi_crm_b_assist").Where("product_assist_id=?", assistID).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", userIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user").Where("id IN ?", userIDs).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	for _, f := range fixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{"id": f.ID, "username": "assist-rbac-" + f.Role, "password_hash": "not-used", "display_name": "中文助力验收-" + f.Role, "status": 1, "auth_version": 1, "data_scope_version": 1}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": f.ID, "role_id": roleID[f.Role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now()
	if err := businessDB.Table("qixi_crm_b_assist").Create(map[string]any{"product_assist_id": assistID, "start_time": now.Add(-time.Hour), "end_time": now.Add(time.Hour), "status": 1, "assist_count": 2, "assist_user_count": 1, "product_id": 987677901, "assist_price": "199.00", "stock": 12, "is_show": 1, "store_name": "中文好友助力快照验收商品", "mer_id": 1, "store_info": "中文模拟好友助力活动。", "is_del": 0, "product_status": 1, "action_status": 1}).Error; err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager(strings.Repeat("x", 32), time.Minute, 2*time.Minute)
	router := gin.New()
	group := router.Group("/api/platform/v1")
	group.Use(middleware.JWTRequired(jwt, authjwt.PortalPlatform), middleware.RequireAdminConsole(), middleware.RequireAdminSession(adminDB), middleware.RestrictRoleConsole(), middleware.RestrictRegionConsole())
	svc := domain.NewService(persist.NewRepo(businessDB))
	NewHandler(svc, adminDB).Register(group)
	call := func(f struct {
		ID   uint
		Role string
	}, method, path string, body any) *httptest.ResponseRecorder {
		var raw []byte
		if body != nil {
			raw, _ = json.Marshal(body)
		}
		pair, e := jwt.IssueAdminConsole(f.ID, "assist-rbac-"+f.Role, []string{f.Role}, 1)
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
	if got := call(fixtures[3], http.MethodGet, "/api/platform/v1/assist/actives?mer_id=1", nil); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "中文好友助力快照验收商品") {
		t.Fatalf("operations list code=%d body=%s", got.Code, got.Body.String())
	}
	if got := call(fixtures[0], http.MethodGet, "/api/platform/v1/assist/actives/987677201", nil); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "中文好友助力快照验收商品") {
		t.Fatalf("platform detail code=%d body=%s", got.Code, got.Body.String())
	}
	for _, f := range []struct {
		ID   uint
		Role string
	}{fixtures[1], fixtures[2], fixtures[4]} {
		if got := call(f, http.MethodGet, "/api/platform/v1/assist/actives", nil).Code; got != http.StatusForbidden {
			t.Fatalf("%s list=%d", f.Role, got)
		}
	}
	if got := call(fixtures[0], http.MethodPut, "/api/platform/v1/assist/actives/987677201", gin.H{"assist_price": 1}).Code; got != http.StatusBadRequest {
		t.Fatalf("price-only=%d", got)
	}
	if got := call(fixtures[3], http.MethodPut, "/api/platform/v1/assist/actives/987677201", gin.H{"is_show": 0, "assist_price": 1, "stock": 1}).Code; got != http.StatusOK {
		t.Fatalf("operations hide=%d", got)
	}
	var row struct {
		IsShow      int
		AssistPrice float64
		Stock       int
	}
	if err := businessDB.Table("qixi_crm_b_assist").Select("is_show,assist_price,stock").Where("product_assist_id=?", assistID).Take(&row).Error; err != nil || row.IsShow != 0 || row.AssistPrice != 199 || row.Stock != 12 {
		t.Fatalf("row=%#v err=%v", row, err)
	}
	if _, err := svc.StartSet(context.Background(), 987677801, assistID); err == nil {
		t.Fatal("hidden assist must reject start")
	}
}
