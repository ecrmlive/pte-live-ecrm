package merchant

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	domainmerchant "github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/merchant"
	merchantpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/merchant"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMerchantCategoryIntegrationCRUDAndConflict(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_MERCHANT_CATEGORY_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_MERCHANT_CATEGORY_TEST_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		db.Table("qixi_crm_a_merchant_category").Where("name LIKE ?", "中文商户分类验收%").Delete(nil)
	}
	cleanup()
	defer cleanup()

	gin.SetMode(gin.TestMode)
	h := NewHandler(domainmerchant.NewService(merchantpersist.NewStoreAdapter(merchantpersist.NewRepo(db))), nil, db, nil)
	r := gin.New()
	r.POST("/merchant-categories", h.CreateCategory)
	r.PUT("/merchant-categories/:id", h.UpdateCategory)
	r.DELETE("/merchant-categories/:id", h.DeleteCategory)
	call := func(method, path string, input any) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(input)
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	if w := call(http.MethodPost, "/merchant-categories", gin.H{"category_name": "中文商户分类验收精度", "commission_rate": 8.123}); w.Code != http.StatusBadRequest {
		t.Fatalf("fractional commission=%d %s", w.Code, w.Body.String())
	}
	w := call(http.MethodPost, "/merchant-categories", gin.H{"category_name": "中文商户分类验收服饰", "commission_rate": 8.5})
	if w.Code != http.StatusOK {
		t.Fatalf("create=%d %s", w.Code, w.Body.String())
	}
	var created struct {
		Data domainmerchant.Category `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil || created.Data.MerchantCategoryID == 0 {
		t.Fatalf("created=%+v err=%v", created, err)
	}
	if w := call(http.MethodPost, "/merchant-categories", gin.H{"category_name": "中文商户分类验收服饰", "commission_rate": 8.5}); w.Code != http.StatusConflict {
		t.Fatalf("duplicate=%d %s", w.Code, w.Body.String())
	}
	if w := call(http.MethodPut, "/merchant-categories/999999999", gin.H{"category_name": "中文商户分类验收缺失", "commission_rate": 9}); w.Code != http.StatusNotFound {
		t.Fatalf("missing update=%d %s", w.Code, w.Body.String())
	}
	path := "/merchant-categories/" + strconv.FormatUint(uint64(created.Data.MerchantCategoryID), 10)
	if w := call(http.MethodPut, path, gin.H{"category_name": "中文商户分类验收家居", "commission_rate": 9.25}); w.Code != http.StatusOK {
		t.Fatalf("update=%d %s", w.Code, w.Body.String())
	}
	if w := call(http.MethodPut, path, gin.H{"category_name": "中文商户分类验收家居", "commission_rate": 9.25}); w.Code != http.StatusOK {
		t.Fatalf("idempotent update=%d %s", w.Code, w.Body.String())
	}
	if w := call(http.MethodDelete, "/merchant-categories/999999999", nil); w.Code != http.StatusNotFound {
		t.Fatalf("missing delete=%d %s", w.Code, w.Body.String())
	}
	if w := call(http.MethodDelete, path, nil); w.Code != http.StatusOK {
		t.Fatalf("delete=%d %s", w.Code, w.Body.String())
	}
	var total int64
	db.Table("qixi_crm_a_merchant_category").Where("id = ?", created.Data.MerchantCategoryID).Count(&total)
	if total != 0 {
		t.Fatalf("category remains=%d", total)
	}
}

func TestMerchantCategoryHTTPRBACAcrossFiveRoles(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_MERCHANT_CATEGORY_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_MERCHANT_CATEGORY_TEST_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	roles := []string{"platform", "merchant", "region", "operations", "customer_service"}
	type roleRecord struct {
		ID   uint
		Code string
	}
	var storedRoles []roleRecord
	if err := db.Table("qixi_crm_a_role").Select("id, code").Where("code IN ? AND status = 1", roles).Find(&storedRoles).Error; err != nil {
		t.Fatal(err)
	}
	roleIDs := make(map[string]uint, len(storedRoles))
	for _, item := range storedRoles {
		roleIDs[item.Code] = item.ID
	}
	if len(roleIDs) != len(roles) {
		t.Fatalf("five role fixture missing: %#v", roleIDs)
	}
	userIDs := make([]uint, 0, len(roles))
	for index := range roles {
		userIDs = append(userIDs, 987680130+uint(index))
	}
	cleanup := func() {
		db.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", userIDs).Delete(nil)
		db.Table("qixi_crm_a_admin_user").Where("id IN ?", userIDs).Delete(nil)
	}
	cleanup()
	defer cleanup()
	for index, role := range roles {
		userID := userIDs[index]
		if err := db.Table("qixi_crm_a_admin_user").Create(map[string]any{
			"id":                 userID,
			"username":           "category-rbac-" + role,
			"password_hash":      "not-used-by-rbac-test",
			"display_name":       "中文分类权限验收-" + role,
			"status":             1,
			"auth_version":       1,
			"data_scope_version": 1,
		}).Error; err != nil {
			t.Fatal(err)
		}
		if err := db.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": userID, "role_id": roleIDs[role]}).Error; err != nil {
			t.Fatal(err)
		}
	}

	gin.SetMode(gin.TestMode)
	jwtManager := authjwt.NewManager(strings.Repeat("x", 32), time.Minute, 2*time.Minute)
	r := gin.New()
	authed := r.Group("/api/platform/v1")
	authed.Use(
		middleware.JWTRequired(jwtManager, authjwt.PortalPlatform),
		middleware.RequireAdminConsole(),
		middleware.RequireAdminSession(db),
		middleware.RestrictRoleConsole(),
		middleware.RestrictRegionConsole(),
	)
	NewHandler(domainmerchant.NewService(merchantpersist.NewStoreAdapter(merchantpersist.NewRepo(db))), nil, db, nil).Register(authed)
	call := func(role string, userID uint) int {
		pair, issueErr := jwtManager.IssueAdminConsole(userID, "category-rbac-"+role, []string{role}, 1)
		if issueErr != nil {
			t.Fatal(issueErr)
		}
		req := httptest.NewRequest(http.MethodGet, "/api/platform/v1/merchant-categories", nil)
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}
	if got := call("platform", userIDs[0]); got != http.StatusOK {
		t.Fatalf("platform=%d", got)
	}
	for index, role := range roles[1:] {
		if got := call(role, userIDs[index+1]); got != http.StatusForbidden {
			t.Fatalf("%s=%d, want 403", role, got)
		}
	}
}
