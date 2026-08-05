package circle

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	domaincircle "github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/circle"
	circlepersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/circle"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAgentRevocationHTTPRBACAcrossFiveRoles(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_CIRCLE_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_CIRCLE_TEST_DSN to run circle agent RBAC integration test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	roles := []string{"platform", "merchant", "region", "operations", "customer_service"}
	type roleRow struct {
		ID   uint
		Code string
	}
	var storedRoles []roleRow
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

	const agentID uint = 987670091
	userIDs := make([]uint, len(roles))
	for index := range roles {
		userIDs[index] = 987670091 + uint(index+1)
	}
	cleanup := func() {
		_ = db.Table("qixi_crm_a_business_zone_agent_command_audit").Where("circle_agent_id = ?", agentID).Delete(nil).Error
		_ = db.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", userIDs).Delete(nil).Error
		_ = db.Table("qixi_crm_a_admin_user").Where("id IN ?", userIDs).Delete(nil).Error
		_ = db.Table("qixi_crm_a_business_zone_agent").Where("circle_agent_id = ?", agentID).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	if err := db.Create(&domaincircle.Agent{CircleAgentID: agentID, Name: "中文五角色撤销验收代理", Phone: "13800000091", Status: domaincircle.AgentApproved, CreateTime: time.Now(), UpdateTime: time.Now()}).Error; err != nil {
		t.Fatal(err)
	}
	for index, role := range roles {
		userID := userIDs[index]
		if err := db.Table("qixi_crm_a_admin_user").Create(map[string]any{
			"id": userID, "username": "circle-revoke-rbac-" + role, "password_hash": "not-used-by-rbac-test",
			"display_name": "中文区域代理权限验收-" + role, "status": 1, "auth_version": 1, "data_scope_version": 1,
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
	NewHandler(domaincircle.NewService(circlepersist.NewRepo(db)), db).Register(authed)
	call := func(role string, userID uint) int {
		pair, issueErr := jwtManager.IssueAdminConsole(userID, "circle-revoke-rbac-"+role, []string{role}, 1)
		if issueErr != nil {
			t.Fatal(issueErr)
		}
		raw, _ := json.Marshal(gin.H{"reason": "中文五角色撤销验收", "idempotency_key": "circle-rbac-revoke-001"})
		req := httptest.NewRequest(http.MethodDelete, "/api/platform/v1/business-zone-agents/987670091", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
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
	request := httptest.NewRequest(http.MethodDelete, "/api/platform/v1/business-zone-agents/987670091", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous=%d, want 401", w.Code)
	}
}

func TestAgentPasswordResetHTTPRBACAndSessionInvalidation(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_CIRCLE_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_CIRCLE_TEST_DSN to run circle agent password RBAC integration test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	roles := []string{"platform", "merchant", "region", "operations", "customer_service"}
	type roleRow struct {
		ID   uint
		Code string
	}
	var storedRoles []roleRow
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

	const agentID uint = 987670094
	const boundAdminID uint = 987670095
	userIDs := make([]uint, len(roles))
	for index := range roles {
		userIDs[index] = 987680130 + uint(index)
	}
	cleanup := func() {
		_ = db.Table("qixi_crm_a_business_zone_agent_password_reset_audit").Where("circle_agent_id = ?", agentID).Delete(nil).Error
		_ = db.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", append(userIDs, boundAdminID)).Delete(nil).Error
		_ = db.Table("qixi_crm_a_admin_user").Where("id IN ?", append(userIDs, boundAdminID)).Delete(nil).Error
		_ = db.Table("qixi_crm_a_business_zone_agent").Where("circle_agent_id = ?", agentID).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	oldHash, err := bcrypt.GenerateFromPassword([]byte("LocalOnlyOldAgentPassword01"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&domaincircle.Agent{CircleAgentID: agentID, Name: "中文五角色口令验收代理", Phone: "13800000094", Status: domaincircle.AgentApproved, CreateTime: time.Now(), UpdateTime: time.Now()}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Table("qixi_crm_a_admin_user").Create(map[string]any{
		"id": boundAdminID, "username": "circle-bound-region-user", "password_hash": string(oldHash),
		"display_name": "中文绑定区域账号", "status": 1, "auth_version": 1, "data_scope_version": 1, "circle_agent_id": agentID,
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": boundAdminID, "role_id": roleIDs["region"]}).Error; err != nil {
		t.Fatal(err)
	}
	for index, role := range roles {
		userID := userIDs[index]
		if err := db.Table("qixi_crm_a_admin_user").Create(map[string]any{
			"id": userID, "username": "circle-password-rbac-" + role, "password_hash": "not-used-by-rbac-test",
			"display_name": "中文区域口令权限验收-" + role, "status": 1, "auth_version": 1, "data_scope_version": 1,
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
	NewHandler(domaincircle.NewService(circlepersist.NewRepo(db)), db).Register(authed)
	staleBoundSession, err := jwtManager.IssueAdminConsoleWithIdentityVersion(boundAdminID, "circle-bound-region-user", []string{"region"}, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	call := func(role string, userID uint) int {
		pair, issueErr := jwtManager.IssueAdminConsole(userID, "circle-password-rbac-"+role, []string{role}, 1)
		if issueErr != nil {
			t.Fatal(issueErr)
		}
		raw, _ := json.Marshal(gin.H{"password": "LocalOnlyNewAgentPassword02", "reason": "中文五角色密码重置验收", "idempotency_key": "circle-rbac-password-001"})
		req := httptest.NewRequest(http.MethodPost, "/api/platform/v1/business-zone-agents/987670094/password", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
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
	staleRequest := httptest.NewRequest(http.MethodGet, "/api/platform/v1/business-zone-agents/settings", nil)
	staleRequest.Header.Set("Authori-zation", "Bearer "+staleBoundSession.AccessToken)
	staleResponse := httptest.NewRecorder()
	r.ServeHTTP(staleResponse, staleRequest)
	if staleResponse.Code != http.StatusUnauthorized {
		t.Fatalf("stale bound admin session=%d, want 401", staleResponse.Code)
	}
	anonymous := httptest.NewRequest(http.MethodPost, "/api/platform/v1/business-zone-agents/987670094/password", nil)
	anonymousResponse := httptest.NewRecorder()
	r.ServeHTTP(anonymousResponse, anonymous)
	if anonymousResponse.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous=%d, want 401", anonymousResponse.Code)
	}
	var updated struct {
		Hash    string `gorm:"column:password_hash"`
		Version uint64 `gorm:"column:auth_version"`
	}
	if err := db.Table("qixi_crm_a_admin_user").Select("password_hash, auth_version").Where("id = ?", boundAdminID).Take(&updated).Error; err != nil {
		t.Fatal(err)
	}
	if updated.Version != 2 || bcrypt.CompareHashAndPassword([]byte(updated.Hash), []byte("LocalOnlyNewAgentPassword02")) != nil {
		t.Fatalf("password reset projection is invalid: version=%d", updated.Version)
	}
}
