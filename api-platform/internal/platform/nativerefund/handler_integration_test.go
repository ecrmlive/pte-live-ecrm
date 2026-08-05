package nativerefund

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestRefundHTTPRBACAndMerchantRegionScopes uses deliberately fake Chinese
// merchants and refunds. It proves that refund supervision does not treat a
// hidden Vben menu as authorization: merchant and region accounts can only
// read their assigned merchant facts, and neither can approve a refund.
func TestRefundHTTPRBACAndMerchantRegionScopes(t *testing.T) {
	adminDSN := strings.TrimSpace(os.Getenv("ECRM_REFUND_ADMIN_TEST_DSN"))
	businessDSN := strings.TrimSpace(os.Getenv("ECRM_REFUND_BUSINESS_TEST_DSN"))
	merchantDSN := strings.TrimSpace(os.Getenv("ECRM_REFUND_MERCHANT_TEST_DSN"))
	if adminDSN == "" || businessDSN == "" || merchantDSN == "" {
		t.Skip("set ECRM_REFUND_ADMIN_TEST_DSN, ECRM_REFUND_BUSINESS_TEST_DSN and ECRM_REFUND_MERCHANT_TEST_DSN")
	}
	adminDB, err := gorm.Open(mysql.Open(adminDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	businessDB, err := gorm.Open(mysql.Open(businessDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	merchantDB, err := gorm.Open(mysql.Open(merchantDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	roles := []string{"platform", "merchant", "region", "operations", "customer_service"}
	var storedRoles []struct {
		ID   uint
		Code string
	}
	if err := adminDB.Table("qixi_crm_a_role").Select("id,code").Where("code IN ? AND status=1", roles).Find(&storedRoles).Error; err != nil || len(storedRoles) != len(roles) {
		t.Fatalf("roles=%v err=%v", storedRoles, err)
	}
	roleIDs := make(map[string]uint, len(storedRoles))
	for _, item := range storedRoles {
		roleIDs[item.Code] = item.ID
	}

	type fixture struct {
		ID   uint
		Role string
	}
	fixtures := []fixture{
		{987678401, "platform"},
		{987678402, "merchant"},
		{987678403, "region"},
		{987678404, "operations"},
		{987678405, "customer_service"},
		{987678406, "merchant"}, // intentionally has no data scope
	}
	userIDs := []uint{987678401, 987678402, 987678403, 987678404, 987678405, 987678406}
	const (
		insideMerchantID  uint64 = 987678411
		outsideMerchantID uint64 = 987678412
		insideOrderID     uint64 = 987678431
		outsideOrderID    uint64 = 987678432
		insideRefundID    uint64 = 987678441
		outsideRefundID   uint64 = 987678442
		insideRegionID    uint64 = 987678421
		outsideRegionID   uint64 = 987678422
	)
	cleanup := func() {
		ctx := context.Background()
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_refund_export_audit").Where("operator_admin_id IN ?", userIDs).Delete(nil).Error
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_refund_event").Where("refund_id IN ?", []uint64{insideRefundID, outsideRefundID}).Delete(nil).Error
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_refund").Where("id IN ?", []uint64{insideRefundID, outsideRefundID}).Delete(nil).Error
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_order").Where("id IN ?", []uint64{insideOrderID, outsideOrderID}).Delete(nil).Error
		_ = merchantDB.WithContext(ctx).Table("qixi_crm_m_merchant").Where("id IN ?", []uint64{insideMerchantID, outsideMerchantID}).Delete(nil).Error
		_ = adminDB.WithContext(ctx).Table("qixi_crm_a_data_scope").Where("admin_user_id IN ?", userIDs).Delete(nil).Error
		_ = adminDB.WithContext(ctx).Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", userIDs).Delete(nil).Error
		_ = adminDB.WithContext(ctx).Table("qixi_crm_a_admin_user").Where("id IN ?", userIDs).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	for _, item := range fixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{
			"id": item.ID, "username": "refund-scope-" + strconv.Itoa(int(item.ID)), "password_hash": "not-used",
			"display_name": "中文退款范围验收-" + item.Role, "status": 1, "auth_version": 1, "data_scope_version": 1,
		}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": item.ID, "role_id": roleIDs[item.Role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := adminDB.Table("qixi_crm_a_data_scope").Create([]map[string]any{
		{"admin_user_id": fixtures[1].ID, "scope_type": "merchant", "scope_value": `{"merchant_ids":[987678411]}`, "version": 1},
		{"admin_user_id": fixtures[2].ID, "scope_type": "region", "scope_value": `[987678421]`, "version": 1},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := merchantDB.Table("qixi_crm_m_merchant").Create([]map[string]any{
		{"id": insideMerchantID, "name": "七禧退款区域内茶铺", "status": 1, "region_id": insideRegionID},
		{"id": outsideMerchantID, "name": "七禧退款区域外花店", "status": 1, "region_id": outsideRegionID},
	}).Error; err != nil {
		t.Fatal(err)
	}
	now := time.Now().Add(-time.Minute).Truncate(time.Second)
	if err := businessDB.Table("qixi_crm_b_order").Create([]map[string]any{
		{"id": insideOrderID, "group_order_id": insideOrderID, "order_no": "RBAC-中文-内-001", "merchant_id": insideMerchantID, "merchant_name_snapshot": "七禧退款区域内茶铺", "store_id": insideMerchantID, "store_name_snapshot": "七禧退款区域内茶铺总店", "user_id": 987678451, "total_amount": "66.50", "pay_amount": "66.50", "total_quantity": 1, "recipient_snapshot": `{}`, "status": "aftersale", "paid_at": now, "created_at": now, "updated_at": now},
		{"id": outsideOrderID, "group_order_id": outsideOrderID, "order_no": "RBAC-中文-外-002", "merchant_id": outsideMerchantID, "merchant_name_snapshot": "七禧退款区域外花店", "store_id": outsideMerchantID, "store_name_snapshot": "七禧退款区域外花店总店", "user_id": 987678452, "total_amount": "88.80", "pay_amount": "88.80", "total_quantity": 1, "recipient_snapshot": `{}`, "status": "aftersale", "paid_at": now, "created_at": now, "updated_at": now},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_refund").Create([]map[string]any{
		{"id": insideRefundID, "order_id": insideOrderID, "refund_no": "RF-RBAC-中文-内-001", "reason": "中文模拟退款：茶叶破损", "amount": "66.50", "refund_type": "money_only", "order_status_before": "paid", "status": "applied", "idempotency_key": "refund-rbac-inside-001", "created_at": now, "updated_at": now},
		{"id": outsideRefundID, "order_id": outsideOrderID, "refund_no": "RF-RBAC-中文-外-002", "reason": "中文模拟退款：鲜花缺水", "amount": "88.80", "refund_type": "money_only", "order_status_before": "paid", "status": "applied", "idempotency_key": "refund-rbac-outside-002", "created_at": now, "updated_at": now},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_refund_event").Create(map[string]any{
		"refund_id": insideRefundID, "from_status": "merchant_handling", "to_status": "applied", "actor_type": "merchant", "actor_id": insideMerchantID,
		"reason": "中文模拟日志：商户提交平台仲裁", "idempotency_key": "refund-rbac-event-inside-001", "created_at": now,
	}).Error; err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager(strings.Repeat("x", 32), time.Minute, 2*time.Minute)
	router := gin.New()
	group := router.Group("/api/platform/v1")
	group.Use(middleware.JWTRequired(jwt, authjwt.PortalPlatform), middleware.RequireAdminConsole(), middleware.RequireAdminSession(adminDB), middleware.RestrictRoleConsole(), middleware.RestrictRegionConsole())
	NewHandler(businessDB, merchantDB, adminDB).Register(group)
	call := func(item fixture, method, path string, body any) *httptest.ResponseRecorder {
		pair, issueErr := jwt.IssueAdminConsole(item.ID, "refund-scope-"+strconv.Itoa(int(item.ID)), []string{item.Role}, 1)
		if issueErr != nil {
			t.Fatal(issueErr)
		}
		var raw []byte
		if body != nil {
			raw, issueErr = json.Marshal(body)
			if issueErr != nil {
				t.Fatal(issueErr)
			}
		}
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}

	if got := call(fixtures[0], http.MethodGet, "/api/platform/v1/refunds", nil); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "RF-RBAC-中文-内-001") || !strings.Contains(got.Body.String(), "RF-RBAC-中文-外-002") {
		t.Fatalf("platform list code=%d body=%s", got.Code, got.Body.String())
	}
	if got := call(fixtures[0], http.MethodGet, "/api/platform/v1/refunds/987678442", nil); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "中文模拟退款：鲜花缺水") {
		t.Fatalf("platform cross-region detail code=%d body=%s", got.Code, got.Body.String())
	}
	if got := call(fixtures[0], http.MethodGet, "/api/platform/v1/refunds/987678441/events", nil); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "中文模拟日志：商户提交平台仲裁") {
		t.Fatalf("platform events code=%d body=%s", got.Code, got.Body.String())
	}
	exportResult := call(fixtures[0], http.MethodPost, "/api/platform/v1/refunds/export", gin.H{"reason": "中文模拟导出：售后例行核对", "status": "applied"})
	if exportResult.Code != http.StatusOK || !strings.Contains(exportResult.Body.String(), "退款单号") || strings.Contains(exportResult.Body.String(), "中文模拟退款：茶叶破损") {
		t.Fatalf("platform export code=%d body=%s", exportResult.Code, exportResult.Body.String())
	}
	var exportAudit struct {
		RowCount        int    `gorm:"column:row_count"`
		Reason          string `gorm:"column:reason"`
		OperatorAdminID uint   `gorm:"column:operator_admin_id"`
	}
	if err := businessDB.Table("qixi_crm_b_refund_export_audit").Select("row_count,reason,operator_admin_id").Where("operator_admin_id = ?", fixtures[0].ID).Order("id DESC").Take(&exportAudit).Error; err != nil || exportAudit.RowCount != 2 || exportAudit.Reason != "中文模拟导出：售后例行核对" || exportAudit.OperatorAdminID != fixtures[0].ID {
		t.Fatalf("export audit=%#v err=%v", exportAudit, err)
	}
	for _, item := range []fixture{fixtures[1], fixtures[2]} {
		if got := call(item, http.MethodGet, "/api/platform/v1/refunds", nil); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "RF-RBAC-中文-内-001") || strings.Contains(got.Body.String(), "RF-RBAC-中文-外-002") {
			t.Fatalf("%s scoped list code=%d body=%s", item.Role, got.Code, got.Body.String())
		}
		if got := call(item, http.MethodGet, "/api/platform/v1/refunds/987678442", nil).Code; got != http.StatusNotFound {
			t.Fatalf("%s cross-region detail=%d", item.Role, got)
		}
		if got := call(item, http.MethodPost, "/api/platform/v1/refunds/987678441/approve", nil).Code; got != http.StatusForbidden {
			t.Fatalf("%s approve=%d", item.Role, got)
		}
		if got := call(item, http.MethodGet, "/api/platform/v1/refunds/987678441/events", nil).Code; got != http.StatusForbidden {
			t.Fatalf("%s events=%d", item.Role, got)
		}
		if got := call(item, http.MethodPost, "/api/platform/v1/refunds/export", gin.H{"reason": "中文模拟越权导出"}).Code; got != http.StatusForbidden {
			t.Fatalf("%s export=%d", item.Role, got)
		}
	}
	for _, item := range []fixture{fixtures[3], fixtures[4], fixtures[5]} {
		if got := call(item, http.MethodGet, "/api/platform/v1/refunds", nil).Code; got != http.StatusForbidden {
			t.Fatalf("%s list=%d", item.Role, got)
		}
		if got := call(item, http.MethodGet, "/api/platform/v1/refunds/987678441/events", nil).Code; got != http.StatusForbidden {
			t.Fatalf("%s events=%d", item.Role, got)
		}
		if got := call(item, http.MethodPost, "/api/platform/v1/refunds/export", gin.H{"reason": "中文模拟越权导出"}).Code; got != http.StatusForbidden {
			t.Fatalf("%s export=%d", item.Role, got)
		}
	}
}
