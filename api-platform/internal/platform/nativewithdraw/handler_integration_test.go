package nativewithdraw

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

func TestWithdrawHTTPRBACAndPayoutIdempotency(t *testing.T) {
	adminDSN, businessDSN := strings.TrimSpace(os.Getenv("ECRM_WITHDRAW_ADMIN_TEST_DSN")), strings.TrimSpace(os.Getenv("ECRM_WITHDRAW_BUSINESS_TEST_DSN"))
	if adminDSN == "" || businessDSN == "" {
		t.Skip("set ECRM_WITHDRAW_ADMIN_TEST_DSN and ECRM_WITHDRAW_BUSINESS_TEST_DSN")
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
	for _, row := range dbRoles {
		roleID[row.Code] = row.ID
	}
	fixtures := []struct {
		ID   uint
		Role string
	}{{987676101, "platform"}, {987676102, "merchant"}, {987676103, "region"}, {987676104, "operations"}, {987676105, "customer_service"}}
	users := []uint{987676101, 987676102, 987676103, 987676104, 987676105}
	const withdrawalID uint64 = 987676001
	cleanup := func() {
		_ = businessDB.Table("qixi_crm_b_withdrawal_application").Where("id=?", withdrawalID).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", users).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user").Where("id IN ?", users).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	for _, f := range fixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{"id": f.ID, "username": "withdraw-rbac-" + f.Role, "password_hash": "not-used", "display_name": "中文提现验收-" + f.Role, "status": 1, "auth_version": 1, "data_scope_version": 1}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": f.ID, "role_id": roleID[f.Role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := businessDB.Table("qixi_crm_b_withdrawal_application").Create(map[string]any{"id": withdrawalID, "withdrawal_no": "WD-ACCEPTANCE-中文-001", "user_id": 987676901, "amount": "88.50", "channel": "wechat", "account_snapshot": `{"account_name":"虚构验收用户","account_no":"已脱敏"}`, "status": "applied", "idempotency_key": "withdraw-acceptance-create-001"}).Error; err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager(strings.Repeat("x", 32), time.Minute, 2*time.Minute)
	r := gin.New()
	g := r.Group("/api/platform/v1")
	g.Use(middleware.JWTRequired(jwt, authjwt.PortalPlatform), middleware.RequireAdminConsole(), middleware.RequireAdminSession(adminDB), middleware.RestrictRoleConsole(), middleware.RestrictRegionConsole())
	NewHandler(businessDB, adminDB).Register(g)
	call := func(f struct {
		ID   uint
		Role string
	}, method, path string, body any) *httptest.ResponseRecorder {
		var raw []byte
		if body != nil {
			raw, _ = json.Marshal(body)
		}
		pair, e := jwt.IssueAdminConsole(f.ID, "withdraw-rbac-"+f.Role, []string{f.Role}, 1)
		if e != nil {
			t.Fatal(e)
		}
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	if got := call(fixtures[0], http.MethodPost, "/api/platform/v1/finance/withdraws/987676001/approve", nil).Code; got != http.StatusOK {
		t.Fatalf("approve=%d", got)
	}
	for _, f := range fixtures[1:] {
		if got := call(f, http.MethodGet, "/api/platform/v1/finance/withdraws", nil).Code; got != http.StatusForbidden {
			t.Fatalf("%s list=%d", f.Role, got)
		}
		if got := call(f, http.MethodPost, "/api/platform/v1/finance/withdraws/987676001/mark-paid", gin.H{"idempotency_key": "withdraw-acceptance-paid-001", "payout_reference": "本地模拟凭证-001"}).Code; got != http.StatusForbidden {
			t.Fatalf("%s=%d", f.Role, got)
		}
	}
	payout := gin.H{"idempotency_key": "withdraw-acceptance-paid-001", "payout_reference": "本地模拟凭证-001"}
	if got := call(fixtures[0], http.MethodPost, "/api/platform/v1/finance/withdraws/987676001/mark-paid", payout).Code; got != http.StatusOK {
		t.Fatalf("mark paid=%d", got)
	}
	if got := call(fixtures[0], http.MethodPost, "/api/platform/v1/finance/withdraws/987676001/mark-paid", payout).Code; got != http.StatusOK {
		t.Fatalf("same replay=%d", got)
	}
	if got := call(fixtures[0], http.MethodPost, "/api/platform/v1/finance/withdraws/987676001/mark-paid", gin.H{"idempotency_key": "withdraw-acceptance-paid-001", "payout_reference": "本地模拟篡改凭证-002"}).Code; got != http.StatusConflict {
		t.Fatalf("changed replay=%d", got)
	}
	var row struct {
		Status          string
		PayoutReference string `gorm:"column:payout_reference"`
		PaidBy          uint   `gorm:"column:paid_by"`
	}
	if err := businessDB.Table("qixi_crm_b_withdrawal_application").Select("status,payout_reference,paid_by").Where("id=?", withdrawalID).Take(&row).Error; err != nil || row.Status != "paid" || row.PayoutReference != "本地模拟凭证-001" || row.PaidBy != fixtures[0].ID {
		t.Fatalf("row=%#v err=%v", row, err)
	}
}
