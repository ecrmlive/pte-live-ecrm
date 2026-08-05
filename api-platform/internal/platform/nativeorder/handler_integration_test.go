package nativeorder

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

// TestOrderHTTPRBACAndRegionalScope runs the registered unified-admin routes
// against three disposable databases. The fixture represents two Chinese
// merchants in separate regions and never exercises payment or fulfilment
// writes: this platform feature is supervision-only.
func TestOrderHTTPRBACAndRegionalScope(t *testing.T) {
	adminDSN := strings.TrimSpace(os.Getenv("ECRM_ORDER_ADMIN_TEST_DSN"))
	merchantDSN := strings.TrimSpace(os.Getenv("ECRM_ORDER_MERCHANT_TEST_DSN"))
	businessDSN := strings.TrimSpace(os.Getenv("ECRM_ORDER_BUSINESS_TEST_DSN"))
	if adminDSN == "" || merchantDSN == "" || businessDSN == "" {
		t.Skip("set ECRM_ORDER_ADMIN_TEST_DSN, ECRM_ORDER_MERCHANT_TEST_DSN and ECRM_ORDER_BUSINESS_TEST_DSN to run order integration test")
	}
	adminDB, err := gorm.Open(mysql.Open(adminDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	merchantDB, err := gorm.Open(mysql.Open(merchantDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	businessDB, err := gorm.Open(mysql.Open(businessDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	const (
		platformID        uint   = 988840001
		merchantUserID    uint   = 988840002
		regionUserID      uint   = 988840003
		operationsID      uint   = 988840004
		serviceID         uint   = 988840005
		merchantNoScopeID uint   = 988840006
		merchantA         uint64 = 988840011
		merchantB         uint64 = 988840012
		storeA            uint64 = 988840021
		storeB            uint64 = 988840022
		groupA            uint64 = 988840031
		groupB            uint64 = 988840032
		orderA            uint64 = 988840041
		orderB            uint64 = 988840042
	)
	adminIDs := []uint{platformID, merchantUserID, regionUserID, operationsID, serviceID, merchantNoScopeID}
	merchantIDs := []uint64{merchantA, merchantB}
	storeIDs := []uint64{storeA, storeB}
	groupIDs := []uint64{groupA, groupB}
	orderIDs := []uint64{orderA, orderB}
	cleanup := func() {
		_ = businessDB.Table("qixi_crm_b_order_delivery").Where("order_id IN ?", orderIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_order_item").Where("order_id IN ?", orderIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_order").Where("id IN ?", orderIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_group_order").Where("id IN ?", groupIDs).Delete(nil).Error
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
		{platformID, "platform", "order-platform-acceptance"},
		{merchantUserID, "merchant", "order-merchant-acceptance"},
		{regionUserID, "region", "order-region-acceptance"},
		{operationsID, "operations", "order-operations-acceptance"},
		{serviceID, "customer_service", "order-service-acceptance"},
		{merchantNoScopeID, "merchant", "order-merchant-no-scope"},
	}
	for _, item := range fixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{
			"id": item.id, "username": item.username, "password_hash": "not-used-by-order-integration-test",
			"display_name": "中文订单监管验收员", "status": 1, "auth_version": 1, "data_scope_version": 1,
		}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": item.id, "role_id": roleIDs[item.role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := adminDB.Table("qixi_crm_a_data_scope").Create([]map[string]any{
		{"admin_user_id": merchantUserID, "scope_type": "merchant", "scope_value": `{"merchant_ids":[988840011]}`, "version": 1},
		{"admin_user_id": regionUserID, "scope_type": "region", "scope_value": `[10]`, "version": 1},
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
		{"id": storeA, "merchant_id": merchantA, "app_id": "order-acceptance-a", "name": "七禧区域十店", "status": 1},
		{"id": storeB, "merchant_id": merchantB, "app_id": "order-acceptance-b", "name": "七禧区域二十店", "status": 1},
	}).Error; err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 4, 18, 30, 0, 0, time.Local)
	if err := businessDB.Table("qixi_crm_b_group_order").Create([]map[string]any{
		{"id": groupA, "order_no": "CS-ORDER-GROUP-A", "user_id": 988840101, "total_amount": 128.50, "pay_amount": 128.50, "total_quantity": 2, "recipient_snapshot": `{"recipient":"中文验收用户小满","mobile":"13900000031","province":"上海市","city":"上海市区","district":"黄浦区","detail":"虚构订单地址"}`, "pay_channel": "mock", "pay_status": "pending", "idempotency_key": "order-acceptance-group-a", "created_at": now},
		{"id": groupB, "order_no": "CS-ORDER-GROUP-B", "user_id": 988840102, "total_amount": 88.00, "pay_amount": 88.00, "total_quantity": 1, "recipient_snapshot": `{"recipient":"中文验收用户小雪","mobile":"13900000032","province":"江苏省","city":"苏州市","district":"姑苏区","detail":"范围外虚构订单地址"}`, "pay_channel": "wechat", "pay_status": "paid", "paid_at": now, "idempotency_key": "order-acceptance-group-b", "created_at": now.Add(time.Minute)},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_order").Create([]map[string]any{
		{"id": orderA, "group_order_id": groupA, "order_no": "CS-ORDER-A", "merchant_id": merchantA, "merchant_name_snapshot": "七禧区域十茶铺", "store_id": storeA, "store_name_snapshot": "七禧区域十店", "user_id": 988840101, "total_amount": 128.50, "pay_amount": 128.50, "total_quantity": 2, "recipient_snapshot": `{"mobile":"13900000031","province":"上海市","city":"上海市区","district":"黄浦区","detail":"虚构订单地址"}`, "status": "pending_pay", "created_at": now},
		{"id": orderB, "group_order_id": groupB, "order_no": "CS-ORDER-B", "merchant_id": merchantB, "merchant_name_snapshot": "七禧区域二十花店", "store_id": storeB, "store_name_snapshot": "七禧区域二十店", "user_id": 988840102, "total_amount": 88.00, "pay_amount": 88.00, "total_quantity": 1, "recipient_snapshot": `{"mobile":"13900000032","province":"江苏省","city":"苏州市","district":"姑苏区","detail":"范围外虚构订单地址"}`, "status": "shipped", "paid_at": now, "created_at": now.Add(time.Minute)},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_order_item").Create([]map[string]any{
		{"order_id": orderA, "product_id": 988840201, "merchant_sku_id": 988840211, "sku_key": "颜色:青瓷", "title_snapshot": "中文青瓷茶具套装", "spec_snapshot": `{"颜色":"青瓷"}`, "unit_price": 64.25, "quantity": 2},
		{"order_id": orderB, "product_id": 988840202, "merchant_sku_id": 988840212, "sku_key": "香型:桂花", "title_snapshot": "中文范围外香薰礼盒", "spec_snapshot": `{"香型":"桂花"}`, "unit_price": 88.00, "quantity": 1},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_order_delivery").Create([]map[string]any{
		{"order_id": orderA, "delivery_type": "express", "carrier_code": "qixi_demo_express", "tracking_no": "CS-ORDER-TRACK-A", "status": "pending"},
		{"order_id": orderB, "delivery_type": "express", "carrier_code": "qixi_demo_express", "tracking_no": "CS-ORDER-TRACK-B", "status": "shipped", "delivered_at": now},
	}).Error; err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager(strings.Repeat("o", 32), time.Minute, time.Minute)
	router := gin.New()
	authed := router.Group("/api/platform/v1")
	authed.Use(
		middleware.JWTRequired(jwt, authjwt.PortalPlatform),
		middleware.RequireAdminConsole(),
		middleware.RequireAdminSession(adminDB),
		middleware.RestrictRoleConsole(),
		middleware.RestrictRegionConsole(),
	)
	NewHandler(businessDB, merchantDB, adminDB).Register(authed)
	call := func(item adminFixture, method, path string) *httptest.ResponseRecorder {
		pair, issueErr := jwt.IssueAdminConsole(item.id, item.username, []string{item.role}, 1)
		if issueErr != nil {
			t.Fatal(issueErr)
		}
		req := httptest.NewRequest(method, path, bytes.NewReader(nil))
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	platform, merchant, region := fixtures[0], fixtures[1], fixtures[2]
	operations, service, merchantNoScope := fixtures[3], fixtures[4], fixtures[5]

	platformList := call(platform, http.MethodGet, "/api/platform/v1/orders")
	if platformList.Code != http.StatusOK || !strings.Contains(platformList.Body.String(), "988840041") || !strings.Contains(platformList.Body.String(), "988840042") || !strings.Contains(platformList.Body.String(), "七禧区域十茶铺") || !strings.Contains(platformList.Body.String(), "七禧区域二十花店") {
		t.Fatalf("platform order list=%d body=%s", platformList.Code, platformList.Body.String())
	}
	for _, item := range []adminFixture{merchant, region} {
		list := call(item, http.MethodGet, "/api/platform/v1/orders")
		if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), "988840041") || strings.Contains(list.Body.String(), "988840042") || !strings.Contains(list.Body.String(), "七禧区域十茶铺") {
			t.Fatalf("%s order scope=%d body=%s", item.role, list.Code, list.Body.String())
		}
		if got := call(item, http.MethodGet, "/api/platform/v1/orders/988840042").Code; got != http.StatusNotFound {
			t.Fatalf("%s out-of-scope detail=%d, want 404", item.role, got)
		}
	}
	for _, item := range []adminFixture{operations, service, merchantNoScope} {
		if got := call(item, http.MethodGet, "/api/platform/v1/orders").Code; got != http.StatusForbidden {
			t.Fatalf("%s order list=%d, want 403", item.role, got)
		}
	}
	pending := call(platform, http.MethodGet, "/api/platform/v1/orders?paid=0")
	if pending.Code != http.StatusOK || !strings.Contains(pending.Body.String(), "988840041") || strings.Contains(pending.Body.String(), "988840042") {
		t.Fatalf("unpaid filter=%d body=%s", pending.Code, pending.Body.String())
	}
	paid := call(platform, http.MethodGet, "/api/platform/v1/orders?paid=1")
	if paid.Code != http.StatusOK || strings.Contains(paid.Body.String(), "988840041") || !strings.Contains(paid.Body.String(), "988840042") {
		t.Fatalf("paid filter=%d body=%s", paid.Code, paid.Body.String())
	}
	detail := call(merchant, http.MethodGet, "/api/platform/v1/orders/988840041")
	if detail.Code != http.StatusOK || !strings.Contains(detail.Body.String(), "中文青瓷茶具套装") || !strings.Contains(detail.Body.String(), "CS-ORDER-TRACK-A") || !strings.Contains(detail.Body.String(), "七禧区域十店") {
		t.Fatalf("merchant scoped order detail=%d body=%s", detail.Code, detail.Body.String())
	}
	if got := call(platform, http.MethodGet, "/api/platform/v1/orders/not-a-number").Code; got != http.StatusBadRequest {
		t.Fatalf("invalid order id=%d, want 400", got)
	}
	if got := call(platform, http.MethodPost, "/api/platform/v1/orders").Code; got != http.StatusNotFound {
		t.Fatalf("read-only order route must not expose mutation=%d, want 404", got)
	}
}
