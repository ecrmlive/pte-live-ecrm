package nativecatalog

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

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestProductHTTPRBACAndRegionalScope verifies the production unified-admin
// route stack, not a direct handler call. It keeps all IDs and Chinese data in
// the disposable three-database fixture.
func TestProductHTTPRBACAndRegionalScope(t *testing.T) {
	adminDSN, merchantDSN, businessDSN := strings.TrimSpace(os.Getenv("ECRM_PRODUCT_ADMIN_TEST_DSN")), strings.TrimSpace(os.Getenv("ECRM_PRODUCT_MERCHANT_TEST_DSN")), strings.TrimSpace(os.Getenv("ECRM_PRODUCT_BUSINESS_TEST_DSN"))
	if adminDSN == "" || merchantDSN == "" || businessDSN == "" {
		t.Skip("set ECRM_PRODUCT_ADMIN_TEST_DSN, ECRM_PRODUCT_MERCHANT_TEST_DSN and ECRM_PRODUCT_BUSINESS_TEST_DSN")
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
	const platformID uint = 988830001
	const merchantUserID uint = 988830002
	const regionUserID uint = 988830003
	const operationsID uint = 988830004
	const serviceID uint = 988830005
	const merchantA uint64 = 988830011
	const merchantB uint64 = 988830012
	const storeA uint64 = 988830021
	const storeB uint64 = 988830022
	const productA uint64 = 988830031
	const productB uint64 = 988830032
	ids := []uint{platformID, merchantUserID, regionUserID, operationsID, serviceID}
	cleanup := func() {
		_ = merchantDB.Table("qixi_crm_m_product_audit_outbox").Where("product_id IN ?", []uint64{productA, productB}).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_product_projection_outbox").Where("product_id IN ?", []uint64{productA, productB}).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_product_review").Where("product_id IN ?", []uint64{productA, productB}).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_data_scope").Where("admin_user_id IN ?", ids).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", ids).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user").Where("id IN ?", ids).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_product_sku_view").Where("product_id IN ?", []uint64{productA, productB}).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_product_view").Where("product_id IN ?", []uint64{productA, productB}).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_category_view").Where("category_id = ?", 988830041).Delete(nil).Error
		_ = merchantDB.Table("qixi_crm_m_product_detail").Where("product_id IN ?", []uint64{productA, productB}).Delete(nil).Error
		_ = merchantDB.Table("qixi_crm_m_product_sku").Where("product_id IN ?", []uint64{productA, productB}).Delete(nil).Error
		_ = merchantDB.Table("qixi_crm_m_product").Where("id IN ?", []uint64{productA, productB}).Delete(nil).Error
		_ = merchantDB.Table("qixi_crm_m_store").Where("id IN ?", []uint64{storeA, storeB}).Delete(nil).Error
		_ = merchantDB.Table("qixi_crm_m_merchant").Where("id IN ?", []uint64{merchantA, merchantB}).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	var roleRows []struct {
		ID   uint
		Code string
	}
	roles := []string{"platform", "merchant", "region", "operations", "customer_service"}
	if err := adminDB.Table("qixi_crm_a_role").Select("id,code").Where("code IN ? AND status = 1", roles).Find(&roleRows).Error; err != nil {
		t.Fatal(err)
	}
	roleIDs := map[string]uint{}
	for _, row := range roleRows {
		roleIDs[row.Code] = row.ID
	}
	if len(roleIDs) != 5 {
		t.Fatalf("roles missing: %#v", roleIDs)
	}
	fixtures := []struct {
		id             uint
		role, username string
	}{{platformID, "platform", "product-platform-acceptance"}, {merchantUserID, "merchant", "product-merchant-acceptance"}, {regionUserID, "region", "product-region-acceptance"}, {operationsID, "operations", "product-operations-acceptance"}, {serviceID, "customer_service", "product-service-acceptance"}}
	for _, f := range fixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{"id": f.id, "username": f.username, "password_hash": "not-used", "display_name": "中文商品验收", "status": 1, "auth_version": 1, "data_scope_version": 1}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": f.id, "role_id": roleIDs[f.role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := adminDB.Table("qixi_crm_a_data_scope").Create([]map[string]any{{"admin_user_id": merchantUserID, "scope_type": "merchant", "scope_value": `{"merchant_ids":[988830011]}`, "version": 1}, {"admin_user_id": regionUserID, "scope_type": "region", "scope_value": `[10]`, "version": 1}}).Error; err != nil {
		t.Fatal(err)
	}
	if err := merchantDB.Table("qixi_crm_m_merchant").Create([]map[string]any{{"id": merchantA, "name": "七禧区域十茶铺", "region_id": 10, "status": 1}, {"id": merchantB, "name": "七禧区域二十花店", "region_id": 20, "status": 1}}).Error; err != nil {
		t.Fatal(err)
	}
	if err := merchantDB.Table("qixi_crm_m_store").Create([]map[string]any{{"id": storeA, "merchant_id": merchantA, "app_id": "product-a", "name": "七禧区域十店", "status": 1}, {"id": storeB, "merchant_id": merchantB, "app_id": "product-b", "name": "七禧区域二十店", "status": 1}}).Error; err != nil {
		t.Fatal(err)
	}
	if err := merchantDB.Table("qixi_crm_m_product").Create([]map[string]any{{"id": productA, "store_id": storeA, "title": "中文区域茶具验收商品", "category_id": 988830041, "status": "pending_review", "version": 1}, {"id": productB, "store_id": storeB, "title": "中文范围外花艺商品", "category_id": 988830041, "status": "pending_review", "version": 1}}).Error; err != nil {
		t.Fatal(err)
	}
	if err := merchantDB.Table("qixi_crm_m_product_detail").Create([]map[string]any{{"product_id": productA, "brief": "中文验收详情", "cover_url": "/demo/product-a.png", "original_price": 99.0}, {"product_id": productB, "brief": "范围外详情", "cover_url": "/demo/product-b.png", "original_price": 88.0}}).Error; err != nil {
		t.Fatal(err)
	}
	if err := merchantDB.Table("qixi_crm_m_product_sku").Create([]map[string]any{{"id": 988830051, "product_id": productA, "spec_json": `{"规格":"标准"}`, "price": 66.5, "stock": 8, "status": 1}, {"id": 988830052, "product_id": productB, "spec_json": `{"规格":"标准"}`, "price": 55.5, "stock": 6, "status": 1}}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_category_view").Create(map[string]any{"category_id": 988830041, "name": "中文验收分类", "status": 1}).Error; err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager(strings.Repeat("p", 32), time.Minute, time.Minute)
	r := gin.New()
	authed := r.Group("/api/platform/v1")
	authed.Use(middleware.JWTRequired(jwt, authjwt.PortalPlatform), middleware.RequireAdminConsole(), middleware.RequireAdminSession(adminDB), middleware.RestrictRoleConsole(), middleware.RestrictRegionConsole())
	handler := NewHandler(adminDB, merchantDB, businessDB)
	handler.Register(authed)
	call := func(f struct {
		id             uint
		role, username string
	}, method, path string, body any) *httptest.ResponseRecorder {
		var raw []byte
		if body != nil {
			raw, _ = json.Marshal(body)
		}
		pair, e := jwt.IssueAdminConsole(f.id, f.username, []string{f.role}, 1)
		if e != nil {
			t.Fatal(e)
		}
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	platform, merchant, region, operations, service := fixtures[0], fixtures[1], fixtures[2], fixtures[3], fixtures[4]
	if got := call(platform, http.MethodGet, "/api/platform/v1/products", nil); got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "988830031") || !strings.Contains(got.Body.String(), "988830032") || !strings.Contains(got.Body.String(), "中文区域茶具验收商品") || !strings.Contains(got.Body.String(), "中文范围外花艺商品") || !strings.Contains(got.Body.String(), "七禧区域十茶铺") || !strings.Contains(got.Body.String(), "七禧区域二十花店") {
		t.Fatalf("platform list=%d %s", got.Code, got.Body.String())
	}
	for _, f := range []struct {
		id             uint
		role, username string
	}{merchant, region} {
		got := call(f, http.MethodGet, "/api/platform/v1/products", nil)
		if got.Code != http.StatusOK || !strings.Contains(got.Body.String(), "988830031") || strings.Contains(got.Body.String(), "988830032") || !strings.Contains(got.Body.String(), "七禧区域十茶铺") {
			t.Fatalf("%s scope=%d %s", f.role, got.Code, got.Body.String())
		}
		if got := call(f, http.MethodGet, "/api/platform/v1/products/988830032", nil); got.Code != http.StatusNotFound {
			t.Fatalf("%s outer detail=%d", f.role, got.Code)
		}
	}
	for _, f := range []struct {
		id             uint
		role, username string
	}{operations, service} {
		if got := call(f, http.MethodGet, "/api/platform/v1/products", nil); got.Code != http.StatusForbidden {
			t.Fatalf("%s list=%d", f.role, got.Code)
		}
	}
	for _, f := range []struct {
		id             uint
		role, username string
	}{merchant, region, operations, service} {
		if got := call(f, http.MethodPost, "/api/platform/v1/products/988830031/audit", gin.H{"status": 1}); got.Code != http.StatusForbidden {
			t.Fatalf("%s audit=%d", f.role, got.Code)
		}
	}
	if got := call(platform, http.MethodPost, "/api/platform/v1/products/988830031/audit", gin.H{"status": 1}); got.Code != http.StatusOK {
		t.Fatalf("platform audit=%d %s", got.Code, got.Body.String())
	}
	var status string
	var reviews, views, skuViews, projections, auditCommands int64
	_ = merchantDB.Table("qixi_crm_m_product").Select("status").Where("id=?", productA).Scan(&status).Error
	_ = adminDB.Table("qixi_crm_a_product_review").Where("product_id=? AND status='approved'", productA).Count(&reviews).Error
	_ = businessDB.Table("qixi_crm_b_product_view").Where("product_id=? AND sale_status=1", productA).Count(&views).Error
	_ = businessDB.Table("qixi_crm_b_product_sku_view").Where("product_id=? AND sale_status=1", productA).Count(&skuViews).Error
	_ = adminDB.Table("qixi_crm_a_product_projection_outbox").Where("product_id=? AND action='upsert' AND status='published'", productA).Count(&projections).Error
	_ = merchantDB.Table("qixi_crm_m_product_audit_outbox").Where("product_id=? AND status='published'", productA).Count(&auditCommands).Error
	if status != "on_sale" || reviews != 1 || views != 1 || skuViews != 1 || projections != 1 || auditCommands != 1 {
		t.Fatalf("audit facts status=%s reviews=%d views=%d skuViews=%d projections=%d auditCommands=%d", status, reviews, views, skuViews, projections, auditCommands)
	}
	// Simulate a business-database projection loss after the approved merchant
	// fact is durable. The dispatcher must rebuild the C-end view from merchant
	// facts and mark the same command published again, without another audit.
	if err := merchantDB.Table("qixi_crm_m_product_audit_outbox").Where("product_id=?", productA).Updates(map[string]any{"status": "failed", "last_error": "模拟业务库暂不可用"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := adminDB.Table("qixi_crm_a_product_projection_outbox").Where("product_id=? AND action='upsert'", productA).Updates(map[string]any{"status": "failed", "last_error": "模拟业务库暂不可用"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_product_sku_view").Where("product_id=?", productA).Delete(nil).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_product_view").Where("product_id=?", productA).Delete(nil).Error; err != nil {
		t.Fatal(err)
	}
	if err := handler.DispatchPendingAuditOutboxes(context.Background()); err != nil {
		t.Fatal(err)
	}
	var retryStatus string
	var retryAttempts, projectionAttempts uint
	_ = merchantDB.Table("qixi_crm_m_product_audit_outbox").Select("status,attempts").Where("product_id=?", productA).Row().Scan(&retryStatus, &retryAttempts)
	_ = adminDB.Table("qixi_crm_a_product_projection_outbox").Select("attempts").Where("product_id=? AND action='upsert'", productA).Row().Scan(&projectionAttempts)
	_ = businessDB.Table("qixi_crm_b_product_view").Where("product_id=? AND sale_status=1", productA).Count(&views).Error
	_ = businessDB.Table("qixi_crm_b_product_sku_view").Where("product_id=? AND sale_status=1", productA).Count(&skuViews).Error
	if retryStatus != "published" || retryAttempts != 2 || projectionAttempts != 2 || views != 1 || skuViews != 1 {
		t.Fatalf("projection retry status=%s attempts=%d projectionAttempts=%d views=%d skuViews=%d", retryStatus, retryAttempts, projectionAttempts, views, skuViews)
	}
}
