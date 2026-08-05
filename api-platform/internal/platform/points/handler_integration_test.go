package points

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

func TestPointsHTTPRBACVersionAndOrderSnapshot(t *testing.T) {
	adminDSN, businessDSN := strings.TrimSpace(os.Getenv("ECRM_POINTS_ADMIN_TEST_DSN")), strings.TrimSpace(os.Getenv("ECRM_POINTS_BUSINESS_TEST_DSN"))
	if adminDSN == "" || businessDSN == "" {
		t.Skip("set ECRM_POINTS_ADMIN_TEST_DSN and ECRM_POINTS_BUSINESS_TEST_DSN")
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
	var rows []struct {
		ID   uint
		Code string
	}
	if err := adminDB.Table("qixi_crm_a_role").Select("id,code").Where("code IN ? AND status=1", roles).Find(&rows).Error; err != nil || len(rows) != len(roles) {
		t.Fatalf("roles=%v err=%v", rows, err)
	}
	ids := map[string]uint{}
	for _, r := range rows {
		ids[r.Code] = r.ID
	}
	users := []uint{987675101, 987675102, 987675103, 987675104, 987675105}
	const productID uint64 = 987675001
	const orderID uint64 = 987675011
	cleanup := func() {
		_ = businessDB.Table("qixi_crm_b_group_order").Where("id=?", orderID).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_points_product_view").Where("product_id=?", productID).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", users).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user").Where("id IN ?", users).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	fixtures := []struct {
		ID   uint
		Role string
	}{{users[0], "platform"}, {users[1], "merchant"}, {users[2], "region"}, {users[3], "operations"}, {users[4], "customer_service"}}
	for _, f := range fixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{"id": f.ID, "username": "points-rbac-" + f.Role, "password_hash": "not-used", "display_name": "中文积分验收-" + f.Role, "status": 1, "auth_version": 1, "data_scope_version": 1}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": f.ID, "role_id": ids[f.Role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := businessDB.Table("qixi_crm_b_points_product_view").Create(map[string]any{"product_id": productID, "merchant_id": 1, "store_id": 1, "merchant_name": "中文积分验收商户", "store_name": "中文积分验收店", "title": "中文积分并发商品", "original_price": "88.00", "points_required": 120, "stock": 5, "sale_status": 1, "version": 1}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_group_order").Create(map[string]any{"id": orderID, "order_no": "POINTS-ACCEPTANCE-001", "user_id": 987675901, "total_amount": 0, "discount_amount": 0, "freight_amount": 0, "pay_amount": 0, "total_quantity": 2, "activity_type": 20, "points_amount": 240, "recipient_snapshot": "{}", "pay_status": "closed", "idempotency_key": "points-acceptance-001"}).Error; err != nil {
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
		pair, e := jwt.IssueAdminConsole(f.ID, "points-rbac-"+f.Role, []string{f.Role}, 1)
		if e != nil {
			t.Fatal(e)
		}
		q := httptest.NewRequest(method, path, bytes.NewReader(raw))
		q.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		q.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, q)
		return w
	}
	update := gin.H{"points_required": 150, "stock": 4, "sale_status": 0, "version": 1}
	if got := call(fixtures[3], http.MethodPut, "/api/platform/v1/points/products/987675001", update).Code; got != http.StatusOK {
		t.Fatalf("operations=%d", got)
	}
	if got := call(fixtures[0], http.MethodPut, "/api/platform/v1/points/products/987675001", update).Code; got != http.StatusConflict {
		t.Fatalf("stale=%d", got)
	}
	for _, f := range []struct {
		ID   uint
		Role string
	}{fixtures[1], fixtures[2], fixtures[4]} {
		if got := call(f, http.MethodPut, "/api/platform/v1/points/products/987675001", update).Code; got != http.StatusForbidden {
			t.Fatalf("%s=%d", f.Role, got)
		}
	}
	closedOrders := call(fixtures[3], http.MethodGet, "/api/platform/v1/points/orders?pay_status=closed", nil)
	if closedOrders.Code != http.StatusOK || !strings.Contains(closedOrders.Body.String(), "POINTS-ACCEPTANCE-001") {
		t.Fatalf("closed order filter code=%d body=%s", closedOrders.Code, closedOrders.Body.String())
	}
	var product struct {
		Points  int64 `gorm:"column:points_required"`
		Stock   int
		Sale    int `gorm:"column:sale_status"`
		Version uint64
	}
	if err := businessDB.Table("qixi_crm_b_points_product_view").Select("points_required,stock,sale_status,version").Where("product_id=?", productID).Take(&product).Error; err != nil || product.Points != 150 || product.Stock != 4 || product.Sale != 0 || product.Version != 2 {
		t.Fatalf("product=%#v err=%v", product, err)
	}
	var order struct {
		Points   int64  `gorm:"column:points_amount"`
		Quantity int    `gorm:"column:total_quantity"`
		Status   string `gorm:"column:pay_status"`
	}
	if err := businessDB.Table("qixi_crm_b_group_order").Select("points_amount,total_quantity,pay_status").Where("id=?", orderID).Take(&order).Error; err != nil || order.Points != 240 || order.Quantity != 2 || order.Status != "closed" {
		t.Fatalf("snapshot=%#v err=%v", order, err)
	}
}
