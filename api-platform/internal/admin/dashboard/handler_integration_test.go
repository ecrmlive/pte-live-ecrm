package dashboard

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

// TestDashboardHTTPFiveRoleScopeAndRank verifies the registered summary route
// with three isolated databases. In particular, a store with two followers
// must still count each paid order exactly once in the sales ranking.
func TestDashboardHTTPFiveRoleScopeAndRank(t *testing.T) {
	adminDSN := strings.TrimSpace(os.Getenv("ECRM_DASHBOARD_ADMIN_TEST_DSN"))
	businessDSN := strings.TrimSpace(os.Getenv("ECRM_DASHBOARD_BUSINESS_TEST_DSN"))
	merchantDSN := strings.TrimSpace(os.Getenv("ECRM_DASHBOARD_MERCHANT_TEST_DSN"))
	if adminDSN == "" || businessDSN == "" || merchantDSN == "" {
		t.Skip("set ECRM_DASHBOARD_ADMIN_TEST_DSN, ECRM_DASHBOARD_BUSINESS_TEST_DSN and ECRM_DASHBOARD_MERCHANT_TEST_DSN to run dashboard integration test")
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

	const (
		platformID        uint   = 988850001
		merchantUserID    uint   = 988850002
		regionUserID      uint   = 988850003
		operationsID      uint   = 988850004
		serviceID         uint   = 988850005
		merchantNoScopeID uint   = 988850006
		merchantA         uint64 = 988850011
		merchantB         uint64 = 988850012
		storeA            uint64 = 988850021
		storeB            uint64 = 988850022
		orderA            uint64 = 988850031
		orderB            uint64 = 988850032
	)
	adminIDs := []uint{platformID, merchantUserID, regionUserID, operationsID, serviceID, merchantNoScopeID}
	merchantIDs := []uint64{merchantA, merchantB}
	storeIDs := []uint64{storeA, storeB}
	orderIDs := []uint64{orderA, orderB}
	userIDs := []uint64{988850101, 988850102, 988850103}
	cleanup := func() {
		_ = businessDB.Table("qixi_crm_b_user_follow_store").Where("store_id IN ?", storeIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_user_browse_history").Where("store_id IN ?", storeIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_customer_service_binding").Where("store_id IN ?", storeIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_refund").Where("order_id IN ?", orderIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_order").Where("id IN ?", orderIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_product_view").Where("product_id IN ?", []uint64{988850041, 988850042}).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_store_view").Where("store_id IN ?", storeIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_user").Where("id IN ?", userIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_product_review").Where("product_id IN ?", []uint64{988850041, 988850042}).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_merchant_application").Where("id = ?", 988850051).Delete(nil).Error
		_ = merchantDB.Table("qixi_crm_m_store").Where("id IN ?", storeIDs).Delete(nil).Error
		_ = merchantDB.Table("qixi_crm_m_merchant").Where("id IN ?", merchantIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_data_scope").Where("admin_user_id IN ?", adminIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", adminIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user").Where("id IN ?", adminIDs).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)

	roles := []string{"platform", "merchant", "region", "operations", "customer_service"}
	var roleRows []struct {
		ID   uint
		Code string
	}
	if err := adminDB.Table("qixi_crm_a_role").Select("id,code").Where("code IN ? AND status = 1", roles).Find(&roleRows).Error; err != nil {
		t.Fatal(err)
	}
	roleIDs := map[string]uint{}
	for _, row := range roleRows {
		roleIDs[row.Code] = row.ID
	}
	if len(roleIDs) != len(roles) {
		t.Fatalf("five role fixture missing: %#v", roleIDs)
	}
	type adminFixture struct {
		id       uint
		role     string
		username string
	}
	fixtures := []adminFixture{
		{platformID, "platform", "dashboard-platform-acceptance"},
		{merchantUserID, "merchant", "dashboard-merchant-acceptance"},
		{regionUserID, "region", "dashboard-region-acceptance"},
		{operationsID, "operations", "dashboard-operations-acceptance"},
		{serviceID, "customer_service", "dashboard-service-acceptance"},
		{merchantNoScopeID, "merchant", "dashboard-merchant-no-scope"},
	}
	for _, item := range fixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{
			"id": item.id, "username": item.username, "password_hash": "not-used-by-dashboard-integration-test",
			"display_name": "中文控制台验收员", "status": 1, "auth_version": 1, "data_scope_version": 1,
		}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": item.id, "role_id": roleIDs[item.role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := adminDB.Table("qixi_crm_a_data_scope").Create([]map[string]any{
		{"admin_user_id": merchantUserID, "scope_type": "merchant", "scope_value": `{"merchant_ids":[988850011]}`, "version": 1},
		{"admin_user_id": regionUserID, "scope_type": "region", "scope_value": `[10]`, "version": 1},
		{"admin_user_id": serviceID, "scope_type": "service_queue", "scope_value": `{"store_ids":[988850021]}`, "version": 1},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := merchantDB.Table("qixi_crm_m_merchant").Create([]map[string]any{
		{"id": merchantA, "name": "七禧区域十茶铺", "region_id": 10, "status": 1},
		{"id": merchantB, "name": "七禧区域二十花店", "region_id": 20, "status": 1},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := merchantDB.Table("qixi_crm_m_store").Create([]map[string]any{
		{"id": storeA, "merchant_id": merchantA, "app_id": "dashboard-a", "name": "七禧区域十店", "status": 1},
		{"id": storeB, "merchant_id": merchantB, "app_id": "dashboard-b", "name": "七禧区域二十店", "status": 1},
	}).Error; err != nil {
		t.Fatal(err)
	}
	now := time.Now().Local().Truncate(time.Second)
	if err := businessDB.Table("qixi_crm_b_store_view").Create([]map[string]any{
		{"store_id": storeA, "merchant_id": merchantA, "store_app_id": "dashboard-a", "store_name": "七禧区域十店", "status": 1},
		{"store_id": storeB, "merchant_id": merchantB, "store_app_id": "dashboard-b", "store_name": "七禧区域二十店", "status": 1},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_product_view").Create([]map[string]any{
		{"product_id": 988850041, "merchant_id": merchantA, "store_id": storeA, "merchant_name": "七禧区域十茶铺", "store_name": "七禧区域十店", "title": "中文控制台茶具", "price": 100.0, "stock": 8, "sale_status": 1, "version": 1, "updated_at": now},
		{"product_id": 988850042, "merchant_id": merchantB, "store_id": storeB, "merchant_name": "七禧区域二十花店", "store_name": "七禧区域二十店", "title": "中文控制台花艺", "price": 50.0, "stock": 5, "sale_status": 1, "version": 1, "updated_at": now},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_user").Create([]map[string]any{
		{"id": userIDs[0], "nickname": "中文下单用户小满", "status": 1, "auth_version": 1, "created_at": now},
		{"id": userIDs[1], "nickname": "中文下单用户小雪", "status": 1, "auth_version": 1, "created_at": now},
		{"id": userIDs[2], "nickname": "中文未下单用户小春", "status": 1, "auth_version": 1, "created_at": now},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_order").Create([]map[string]any{
		{"id": orderA, "group_order_id": 1, "order_no": "DASHBOARD-ORDER-A", "merchant_id": merchantA, "merchant_name_snapshot": "七禧区域十茶铺", "store_id": storeA, "store_name_snapshot": "七禧区域十店", "user_id": userIDs[0], "total_amount": 100.0, "pay_amount": 100.0, "total_quantity": 2, "recipient_snapshot": `{}`, "status": "paid", "paid_at": now, "created_at": now},
		{"id": orderB, "group_order_id": 2, "order_no": "DASHBOARD-ORDER-B", "merchant_id": merchantB, "merchant_name_snapshot": "七禧区域二十花店", "store_id": storeB, "store_name_snapshot": "七禧区域二十店", "user_id": userIDs[1], "total_amount": 50.0, "pay_amount": 50.0, "total_quantity": 1, "recipient_snapshot": `{}`, "status": "shipped", "paid_at": now, "created_at": now},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_refund").Create(map[string]any{"order_id": orderA, "refund_no": "DASHBOARD-REFUND-A", "reason": "中文虚构待处理退款", "amount": 100.0, "refund_type": "money_only", "order_status_before": "paid", "status": "applied", "idempotency_key": "dashboard-refund-a"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_customer_service_binding").Create([]map[string]any{
		{"user_id": userIDs[0], "store_id": storeA, "im_conversation_id": "dashboard-service-a", "status": "open", "last_msg": "中文店铺 A 待处理咨询"},
		{"user_id": userIDs[1], "store_id": storeB, "im_conversation_id": "dashboard-service-b", "status": "open", "last_msg": "中文店铺 B 待处理咨询"},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_user_browse_history").Create([]map[string]any{
		{"user_id": userIDs[0], "product_id": 988850041, "store_id": storeA, "viewed_at": now},
		{"user_id": userIDs[2], "product_id": 988850041, "store_id": storeA, "viewed_at": now},
		{"user_id": userIDs[1], "product_id": 988850042, "store_id": storeB, "viewed_at": now},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_user_follow_store").Create([]map[string]any{
		{"user_id": userIDs[0], "store_id": storeA}, {"user_id": userIDs[2], "store_id": storeA}, {"user_id": userIDs[1], "store_id": storeB},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := adminDB.Table("qixi_crm_a_product_review").Create([]map[string]any{
		{"product_id": 988850041, "store_id": storeA, "status": "pending"}, {"product_id": 988850042, "store_id": storeB, "status": "pending"},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := adminDB.Table("qixi_crm_a_merchant_application").Create(map[string]any{"id": 988850051, "merchant_name": "七禧待审核中文店铺", "contact_name": "虚构联系人", "contact_mobile": "13900000041", "status": "pending"}).Error; err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager(strings.Repeat("d", 32), time.Minute, time.Minute)
	router := gin.New()
	authed := router.Group("/api/platform/v1")
	authed.Use(
		middleware.JWTRequired(jwt, authjwt.PortalPlatform),
		middleware.RequireAdminConsole(),
		middleware.RequireAdminSession(adminDB),
		middleware.RestrictRoleConsole(),
		middleware.RestrictRegionConsole(),
	)
	NewHandler(adminDB, businessDB, merchantDB).Register(authed)
	call := func(item adminFixture) *httptest.ResponseRecorder {
		pair, issueErr := jwt.IssueAdminConsole(item.id, item.username, []string{item.role}, 1)
		if issueErr != nil {
			t.Fatal(issueErr)
		}
		req := httptest.NewRequest(http.MethodGet, "/api/platform/v1/dashboard/summary", bytes.NewReader(nil))
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	decode := func(w *httptest.ResponseRecorder) Summary {
		var body struct {
			Status int     `json:"status"`
			Data   Summary `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body.Status != http.StatusOK {
			t.Fatalf("dashboard response=%d body=%s err=%v", w.Code, w.Body.String(), err)
		}
		return body.Data
	}
	platform, merchant, region, operations, service, merchantNoScope := fixtures[0], fixtures[1], fixtures[2], fixtures[3], fixtures[4], fixtures[5]
	platformSummary := decode(call(platform))
	if platformSummary.Scope != "all" || platformSummary.StoreTotal != 2 || platformSummary.PaidOrder != 2 || platformSummary.PendingRefund != 1 || platformSummary.PendingStoreAudit != 1 || platformSummary.PendingProductAudit != 2 || platformSummary.PendingDelivery != 1 || platformSummary.PendingService != 2 || platformSummary.PageViews.Today != 3 || platformSummary.Visitors.Today != 3 || len(platformSummary.StoreSalesRank) != 2 {
		t.Fatalf("platform dashboard=%#v", platformSummary)
	}
	var platformStoreA *StoreSalesRank
	for index := range platformSummary.StoreSalesRank {
		if platformSummary.StoreSalesRank[index].StoreID == storeA {
			platformStoreA = &platformSummary.StoreSalesRank[index]
		}
	}
	if platformStoreA == nil || platformStoreA.FollowerCount != 2 || platformStoreA.SaleCount != 2 || platformStoreA.SaleAmount != 100 {
		t.Fatalf("store A rank must not multiply orders by follower count: %#v", platformStoreA)
	}
	for _, item := range []adminFixture{merchant, region, service} {
		summary := decode(call(item))
		if summary.Scope != "store" || summary.StoreTotal != 1 || summary.PaidOrder != 1 || summary.PendingRefund != 1 || summary.PendingStoreAudit != 0 || summary.PendingProductAudit != 1 || summary.PendingDelivery != 1 || summary.PendingService != 1 || summary.PageViews.Today != 2 || summary.Visitors.Today != 2 || len(summary.StoreSalesRank) != 1 || summary.StoreSalesRank[0].StoreID != storeA || summary.StoreSalesRank[0].SaleAmount != 100 {
			t.Fatalf("%s dashboard scope=%#v", item.role, summary)
		}
	}
	for _, item := range []adminFixture{operations, merchantNoScope} {
		if got := call(item).Code; got != http.StatusForbidden {
			t.Fatalf("%s dashboard=%d, want 403", item.role, got)
		}
	}
}
