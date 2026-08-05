package userlist

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// statusAcceptanceUserID is reserved only in the disposable local acceptance
// database. It contains no personal data and is never loaded by seed data.
const statusAcceptanceUserID uint64 = 987650001

const statusAcceptancePromoterUserID uint64 = 987650003

const statusAcceptanceAdminID uint64 = 987650002

const statusAcceptanceExportStartID uint64 = 987660000

const statusAcceptanceExportEndID uint64 = 987665000

const statusAcceptanceAdminUsername = "status-platform"

const statusAcceptanceAdminPassword = "status-local-fake-password"

const statusAcceptanceCreatedAccount = "LOCAL-CREATE-中文-987650001"

const statusAcceptanceJWTSecret = "local-status-acceptance-not-a-production-secret"

type statusAcceptanceAdmin struct {
	ID               uint64 `gorm:"column:id"`
	AuthVersion      uint64 `gorm:"column:auth_version"`
	DataScopeVersion uint64 `gorm:"column:data_scope_version"`
	Username         string `gorm:"column:username"`
}

func statusAcceptanceDB(t *testing.T) (*gorm.DB, *gorm.DB) {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("ECRM_STATUS_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_STATUS_TEST_DSN to run the disposable MySQL acceptance test")
	}
	business, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open business acceptance database: %v", err)
	}
	adminDSN := strings.Replace(dsn, "qixi_crm_business", "qixi_crm_admin", 1)
	if adminDSN == dsn {
		t.Fatal("acceptance DSN must target qixi_crm_business")
	}
	admin, err := gorm.Open(mysql.Open(adminDSN), &gorm.Config{})
	if err != nil {
		t.Fatalf("open admin acceptance database: %v", err)
	}
	return business, admin
}

func statusAcceptanceRouter(t *testing.T, business, admin *gorm.DB) (*gin.Engine, string, statusAcceptanceAdmin) {
	t.Helper()
	var platformAdmin statusAcceptanceAdmin
	err := admin.Table("qixi_crm_a_admin_user").
		Select("id,auth_version,data_scope_version,username").
		Where("id=? AND status=1 AND deleted_at IS NULL", statusAcceptanceAdminID).Take(&platformAdmin).Error
	if err != nil {
		t.Fatalf("find platform acceptance admin: %v", err)
	}
	jwt := authjwt.NewManager(statusAcceptanceJWTSecret, time.Hour, time.Hour)
	pair, err := jwt.IssueAdminConsoleWithIdentityVersion(uint(platformAdmin.ID), platformAdmin.Username, []string{"platform"}, platformAdmin.DataScopeVersion, platformAdmin.AuthVersion)
	if err != nil {
		t.Fatalf("issue platform acceptance token: %v", err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.JWTRequired(jwt, authjwt.PortalPlatform), middleware.RequireAdminConsole(), middleware.RequireAdminSession(admin))
	New(business, admin).Register(r.Group("/api/platform/v1"))
	return r, pair.AccessToken, platformAdmin
}

func statusAcceptanceRequest(router http.Handler, token string, body any) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/platform/v1/user-list/987650001/status", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authori-zation", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func statusAcceptancePrepareAdmin(t *testing.T, admin *gorm.DB) {
	t.Helper()
	var role struct{ ID uint64 }
	if err := admin.Table("qixi_crm_a_role").Select("id").Where("code='platform' AND status=1").Take(&role).Error; err != nil {
		t.Fatalf("find platform role: %v", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(statusAcceptanceAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash local acceptance password: %v", err)
	}
	if err := admin.Exec("INSERT INTO qixi_crm_a_admin_user (id,username,password_hash,display_name,status,auth_version,data_scope_version) VALUES (?,?,?,?,1,1,1) ON DUPLICATE KEY UPDATE username=VALUES(username),password_hash=VALUES(password_hash),display_name=VALUES(display_name),status=1,auth_version=1,data_scope_version=1,deleted_at=NULL", statusAcceptanceAdminID, statusAcceptanceAdminUsername, string(hash), "中文启停验收平台主管").Error; err != nil {
		t.Fatalf("create local platform acceptance admin: %v", err)
	}
	if err := admin.Exec("INSERT IGNORE INTO qixi_crm_a_admin_user_role (admin_user_id,role_id) VALUES (?,?)", statusAcceptanceAdminID, role.ID).Error; err != nil {
		t.Fatalf("grant local platform acceptance role: %v", err)
	}
}

func statusAcceptancePrepareUser(t *testing.T, business *gorm.DB) {
	t.Helper()
	if err := business.Table("qixi_crm_b_user_notification_command_audit").Where("user_id=?", statusAcceptanceUserID).Delete(nil).Error; err != nil {
		t.Fatalf("clear notification audit: %v", err)
	}
	if err := business.Table("qixi_crm_b_notification").Where("user_id=? AND category=? AND title IN ?", statusAcceptanceUserID, "system", []string{"中文站内通知验收", "中文站内通知页面验收"}).Delete(nil).Error; err != nil {
		t.Fatalf("clear notification fixture: %v", err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user_status_command_audit WHERE user_id=?", statusAcceptanceUserID).Error; err != nil {
		t.Fatalf("clear status audit: %v", err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user_admin_command_audit WHERE user_id=?", statusAcceptanceUserID).Error; err != nil {
		t.Fatalf("clear user-admin command audit: %v", err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user_identity WHERE user_id=?", statusAcceptanceUserID).Error; err != nil {
		t.Fatalf("clear user identity: %v", err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user WHERE id=?", statusAcceptanceUserID).Error; err != nil {
		t.Fatalf("clear status user: %v", err)
	}
	if err := business.Exec("INSERT INTO qixi_crm_b_user (id,nickname,status,auth_version) VALUES (?,?,1,1)", statusAcceptanceUserID, "中文状态验收用户").Error; err != nil {
		t.Fatalf("create Chinese acceptance user: %v", err)
	}
}

func passwordAcceptanceRequest(router http.Handler, token string, body any) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/platform/v1/user-list/987650001/password", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authori-zation", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func createAcceptanceRequest(router http.Handler, token string, body any) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/platform/v1/user-list", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authori-zation", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func profileAcceptanceRequest(router http.Handler, token string, userID uint64, body any) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/platform/v1/user-list/"+strconv.FormatUint(userID, 10)+"/profile", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authori-zation", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func promoterAcceptanceRequest(router http.Handler, token string, body any) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/platform/v1/user-list/promoters/assign", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authori-zation", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func notificationAcceptanceRequest(router http.Handler, token string, userID uint64, body any) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/platform/v1/user-list/"+strconv.FormatUint(userID, 10)+"/notifications", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authori-zation", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func exportAcceptanceRequest(router http.Handler, token string, body any) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/platform/v1/user-list/export", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authori-zation", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestPromoterIntegrationRBACAtomicityAndCommissionBoundary(t *testing.T) {
	business, admin := statusAcceptanceDB(t)
	statusAcceptancePrepareAdmin(t, admin)
	statusAcceptancePrepareUser(t, business)
	if err := business.Table("qixi_crm_b_distribution_promoter_command_audit").Where("idempotency_key LIKE ?", "promoter-%987650001%").Delete(nil).Error; err != nil {
		t.Fatalf("clear promoter audit fixture: %v", err)
	}
	for _, table := range []string{"qixi_crm_b_distribution_promoter", "qixi_crm_b_commission_ledger"} {
		if err := business.Table(table).Where("user_id IN ?", []uint64{statusAcceptanceUserID, statusAcceptancePromoterUserID}).Delete(nil).Error; err != nil {
			t.Fatalf("clear promoter fixture %s: %v", table, err)
		}
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user WHERE id=?", statusAcceptancePromoterUserID).Error; err != nil {
		t.Fatal(err)
	}
	if err := business.Exec("INSERT INTO qixi_crm_b_user (id,nickname,status,auth_version) VALUES (?,?,1,1)", statusAcceptancePromoterUserID, "中文推广资格验收用户").Error; err != nil {
		t.Fatal(err)
	}
	if err := business.Exec("INSERT INTO qixi_crm_b_commission_ledger (user_id,order_id,amount,status,idempotency_key) VALUES (?,?,?,?,?)", statusAcceptanceUserID, 987650001, "18.88", "available", "acceptance-promoter-commission-987650001").Error; err != nil {
		t.Fatal(err)
	}
	router, platformToken, platformAdmin := statusAcceptanceRouter(t, business, admin)
	for _, role := range []string{"merchant", "region", "customer_service", "operations"} {
		jwt := authjwt.NewManager(statusAcceptanceJWTSecret, time.Hour, time.Hour)
		pair, err := jwt.IssueAdminConsoleWithIdentityVersion(uint(platformAdmin.ID), platformAdmin.Username, []string{role}, platformAdmin.DataScopeVersion, platformAdmin.AuthVersion)
		if err != nil {
			t.Fatal(err)
		}
		if response := promoterAcceptanceRequest(router, pair.AccessToken, gin.H{"user_ids": []uint64{statusAcceptanceUserID}, "status": 1, "reason": "虚构中文角色越权", "idempotency_key": "promoter-role-" + role}); response.Code != http.StatusForbidden {
			t.Fatalf("%s promoter=%d, want 403", role, response.Code)
		}
	}
	if response := promoterAcceptanceRequest(router, platformToken, gin.H{"user_ids": []uint64{statusAcceptanceUserID, 999999999}, "status": 1, "reason": "虚构中文部分缺失", "idempotency_key": "promoter-missing-987650001"}); response.Code != http.StatusNotFound {
		t.Fatalf("partial users=%d, want 404", response.Code)
	}
	var promoterCount int64
	if err := business.Table("qixi_crm_b_distribution_promoter").Where("user_id IN ?", []uint64{statusAcceptanceUserID, statusAcceptancePromoterUserID}).Count(&promoterCount).Error; err != nil {
		t.Fatal(err)
	}
	if promoterCount != 0 {
		t.Fatalf("partial request wrote %d promoter rows", promoterCount)
	}

	command := gin.H{"user_ids": []uint64{statusAcceptancePromoterUserID, statusAcceptanceUserID, statusAcceptancePromoterUserID}, "status": 1, "reason": "虚构中文开通推广资格", "idempotency_key": "promoter-enable-987650001"}
	const callers = 6
	responses := make(chan *httptest.ResponseRecorder, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() { defer wg.Done(); responses <- promoterAcceptanceRequest(router, platformToken, command) }()
	}
	wg.Wait()
	close(responses)
	created, replayed := 0, 0
	for response := range responses {
		if response.Code != http.StatusOK {
			t.Fatalf("concurrent promoter=%d: %s", response.Code, response.Body.String())
		}
		if strings.Contains(response.Body.String(), `"replayed":true`) {
			replayed++
		} else {
			created++
		}
	}
	if created != 1 || replayed != callers-1 {
		t.Fatalf("promoter created/replayed=%d/%d", created, replayed)
	}
	if response := promoterAcceptanceRequest(router, platformToken, gin.H{"user_ids": []uint64{statusAcceptanceUserID, statusAcceptancePromoterUserID}, "status": 1, "reason": "不同中文原因", "idempotency_key": "promoter-enable-987650001"}); response.Code != http.StatusConflict {
		t.Fatalf("changed promoter replay=%d", response.Code)
	}
	var enabledCount int64
	if err := business.Table("qixi_crm_b_distribution_promoter").Where("user_id IN ? AND status=1", []uint64{statusAcceptanceUserID, statusAcceptancePromoterUserID}).Count(&enabledCount).Error; err != nil {
		t.Fatal(err)
	}
	if enabledCount != 2 {
		t.Fatalf("enabled promoter count=%d", enabledCount)
	}
	if response := promoterAcceptanceRequest(router, platformToken, gin.H{"user_ids": []uint64{statusAcceptanceUserID, statusAcceptancePromoterUserID}, "status": 0, "reason": "虚构中文停用推广资格", "idempotency_key": "promoter-disable-987650001"}); response.Code != http.StatusOK {
		t.Fatalf("disable promoter=%d", response.Code)
	}
	var disabledCount int64
	if err := business.Table("qixi_crm_b_distribution_promoter").Where("user_id IN ? AND status=0", []uint64{statusAcceptanceUserID, statusAcceptancePromoterUserID}).Count(&disabledCount).Error; err != nil {
		t.Fatal(err)
	}
	if disabledCount != 2 {
		t.Fatalf("disabled promoter count=%d", disabledCount)
	}
	var ledger struct {
		Amount string `gorm:"column:amount"`
		Status string `gorm:"column:status"`
	}
	if err := business.Table("qixi_crm_b_commission_ledger").Select("amount,status").Where("idempotency_key=?", "acceptance-promoter-commission-987650001").Take(&ledger).Error; err != nil {
		t.Fatal(err)
	}
	if ledger.Amount != "18.88" || ledger.Status != "available" {
		t.Fatalf("commission ledger changed: %+v", ledger)
	}
	var auditCount int64
	if err := business.Table("qixi_crm_b_distribution_promoter_command_audit").Where("user_ids_json IS NOT NULL AND idempotency_key IN ?", []string{"promoter-enable-987650001", "promoter-disable-987650001"}).Count(&auditCount).Error; err != nil {
		t.Fatal(err)
	}
	if auditCount != 2 {
		t.Fatalf("promoter audit count=%d", auditCount)
	}
}

func TestCreateAndProfileIntegrationRBACAndAtomicity(t *testing.T) {
	business, admin := statusAcceptanceDB(t)
	statusAcceptancePrepareAdmin(t, admin)
	if err := business.Exec("DELETE FROM qixi_crm_b_user_admin_command_audit WHERE user_id IN (SELECT user_id FROM qixi_crm_b_user_identity WHERE channel='pc' AND subject=?)", statusAcceptanceCreatedAccount).Error; err != nil {
		t.Fatal(err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user_profile WHERE user_id IN (SELECT user_id FROM qixi_crm_b_user_identity WHERE channel='pc' AND subject=?)", statusAcceptanceCreatedAccount).Error; err != nil {
		t.Fatal(err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user WHERE id IN (SELECT user_id FROM qixi_crm_b_user_identity WHERE channel='pc' AND subject=?)", statusAcceptanceCreatedAccount).Error; err != nil {
		t.Fatal(err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user_identity WHERE channel='pc' AND subject=?", statusAcceptanceCreatedAccount).Error; err != nil {
		t.Fatal(err)
	}
	router, platformToken, platformAdmin := statusAcceptanceRouter(t, business, admin)
	create := gin.H{"account": statusAcceptanceCreatedAccount, "password": "本地虚构创建口令456", "nickname": "七禧中文创建验收用户", "reason": "虚构中文创建工单", "idempotency_key": "create-user-987650001"}
	for _, role := range []string{"merchant", "region", "customer_service", "operations"} {
		jwt := authjwt.NewManager(statusAcceptanceJWTSecret, time.Hour, time.Hour)
		pair, err := jwt.IssueAdminConsoleWithIdentityVersion(uint(platformAdmin.ID), platformAdmin.Username, []string{role}, platformAdmin.DataScopeVersion, platformAdmin.AuthVersion)
		if err != nil {
			t.Fatal(err)
		}
		if response := createAcceptanceRequest(router, pair.AccessToken, create); response.Code != http.StatusForbidden {
			t.Fatalf("%s create=%d, want 403: %s", role, response.Code, response.Body.String())
		}
	}
	const callers = 5
	responses := make(chan *httptest.ResponseRecorder, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() { defer wg.Done(); responses <- createAcceptanceRequest(router, platformToken, create) }()
	}
	wg.Wait()
	close(responses)
	succeeded, conflicted := 0, 0
	for response := range responses {
		switch response.Code {
		case http.StatusOK:
			succeeded++
		case http.StatusConflict:
			conflicted++
		default:
			t.Fatalf("concurrent create=%d: %s", response.Code, response.Body.String())
		}
	}
	if succeeded != 1 || conflicted != callers-1 {
		t.Fatalf("create success/conflict=%d/%d", succeeded, conflicted)
	}
	var created struct {
		UserID         uint64 `gorm:"column:user_id"`
		CredentialHash string `gorm:"column:credential_hash"`
		AuthVersion    uint64 `gorm:"column:auth_version"`
	}
	if err := business.Table("qixi_crm_b_user_identity AS i").Select("i.user_id,i.credential_hash,u.auth_version").Joins("JOIN qixi_crm_b_user AS u ON u.id=i.user_id").Where("i.channel='pc' AND i.subject=?", statusAcceptanceCreatedAccount).Take(&created).Error; err != nil {
		t.Fatal(err)
	}
	if created.AuthVersion != 1 || bcrypt.CompareHashAndPassword([]byte(created.CredentialHash), []byte(create["password"].(string))) != nil {
		t.Fatal("created identity mismatch")
	}
	if response := createAcceptanceRequest(router, platformToken, gin.H{"account": statusAcceptanceCreatedAccount, "password": "另一个虚构密码789", "nickname": "七禧中文创建验收用户", "reason": "另一条工单", "idempotency_key": "create-user-second-key"}); response.Code != http.StatusConflict {
		t.Fatalf("duplicate account=%d", response.Code)
	}

	profile := gin.H{"nickname": "七禧中文资料验收用户", "avatar_url": "/demo/avatars/中文验收.png", "gender": 2, "bio": "仅用于本地中文资料维护验收，不含真实个人信息。", "reason": "虚构中文资料工单", "idempotency_key": "profile-update-987650001"}
	for _, role := range []string{"merchant", "region", "customer_service", "operations"} {
		jwt := authjwt.NewManager(statusAcceptanceJWTSecret, time.Hour, time.Hour)
		pair, err := jwt.IssueAdminConsoleWithIdentityVersion(uint(platformAdmin.ID), platformAdmin.Username, []string{role}, platformAdmin.DataScopeVersion, platformAdmin.AuthVersion)
		if err != nil {
			t.Fatal(err)
		}
		if response := profileAcceptanceRequest(router, pair.AccessToken, created.UserID, profile); response.Code != http.StatusForbidden {
			t.Fatalf("%s profile=%d, want 403", role, response.Code)
		}
	}
	responses = make(chan *httptest.ResponseRecorder, callers)
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			responses <- profileAcceptanceRequest(router, platformToken, created.UserID, profile)
		}()
	}
	wg.Wait()
	close(responses)
	createdCount, replayed := 0, 0
	for response := range responses {
		if response.Code != http.StatusOK {
			t.Fatalf("concurrent profile=%d: %s", response.Code, response.Body.String())
		}
		if strings.Contains(response.Body.String(), `"replayed":true`) {
			replayed++
		} else {
			createdCount++
		}
	}
	if createdCount != 1 || replayed != callers-1 {
		t.Fatalf("profile created/replayed=%d/%d", createdCount, replayed)
	}
	var stored struct {
		Nickname  string `gorm:"column:nickname"`
		AvatarURL string `gorm:"column:avatar_url"`
		Gender    int    `gorm:"column:gender"`
		Bio       string `gorm:"column:bio"`
	}
	if err := business.Table("qixi_crm_b_user AS u").Select("u.nickname,p.avatar_url,p.gender,p.bio").Joins("JOIN qixi_crm_b_user_profile AS p ON p.user_id=u.id").Where("u.id=?", created.UserID).Take(&stored).Error; err != nil {
		t.Fatal(err)
	}
	if stored.Nickname != profile["nickname"].(string) || stored.AvatarURL != profile["avatar_url"].(string) || stored.Gender != 2 || stored.Bio != profile["bio"].(string) {
		t.Fatalf("stored profile=%+v", stored)
	}
	if response := profileAcceptanceRequest(router, platformToken, created.UserID, gin.H{"nickname": profile["nickname"], "avatar_url": profile["avatar_url"], "gender": 2, "bio": profile["bio"], "reason": "不同中文原因", "idempotency_key": "profile-update-987650001"}); response.Code != http.StatusConflict {
		t.Fatalf("profile changed replay=%d", response.Code)
	}
}

func TestPasswordResetIntegrationRBACAndSingleUse(t *testing.T) {
	business, admin := statusAcceptanceDB(t)
	statusAcceptancePrepareAdmin(t, admin)
	statusAcceptancePrepareUser(t, business)
	oldPassword := "本地虚构初始口令123"
	oldHash, err := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	if err := business.Exec("INSERT INTO qixi_crm_b_user_identity (user_id,channel,subject,credential_hash) VALUES (?,?,?,?)", statusAcceptanceUserID, "pc", "中文口令验收账号", string(oldHash)).Error; err != nil {
		t.Fatalf("create local PC identity: %v", err)
	}
	router, platformToken, platformAdmin := statusAcceptanceRouter(t, business, admin)
	for _, role := range []string{"merchant", "region", "customer_service", "operations"} {
		jwt := authjwt.NewManager(statusAcceptanceJWTSecret, time.Hour, time.Hour)
		pair, err := jwt.IssueAdminConsoleWithIdentityVersion(uint(platformAdmin.ID), platformAdmin.Username, []string{role}, platformAdmin.DataScopeVersion, platformAdmin.AuthVersion)
		if err != nil {
			t.Fatal(err)
		}
		response := passwordAcceptanceRequest(router, pair.AccessToken, gin.H{"password": "本地虚构重置口令456", "reason": "虚构中文角色越权", "idempotency_key": "password-role-" + role})
		if response.Code != http.StatusForbidden {
			t.Fatalf("%s password reset=%d, want 403: %s", role, response.Code, response.Body.String())
		}
	}

	command := gin.H{"password": "本地虚构重置口令456", "reason": "虚构中文风控重置", "idempotency_key": "password-reset-987650001"}
	const callers = 6
	responses := make(chan *httptest.ResponseRecorder, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			responses <- passwordAcceptanceRequest(router, platformToken, command)
		}()
	}
	wg.Wait()
	close(responses)
	succeeded, conflicted := 0, 0
	for response := range responses {
		switch response.Code {
		case http.StatusOK:
			succeeded++
		case http.StatusConflict:
			conflicted++
		default:
			t.Fatalf("concurrent password reset=%d: %s", response.Code, response.Body.String())
		}
	}
	if succeeded != 1 || conflicted != callers-1 {
		t.Fatalf("password reset success/conflict=%d/%d, want 1/%d", succeeded, conflicted, callers-1)
	}

	var user struct {
		AuthVersion uint64 `gorm:"column:auth_version"`
	}
	if err := business.Table("qixi_crm_b_user").Select("auth_version").Where("id=?", statusAcceptanceUserID).Take(&user).Error; err != nil {
		t.Fatal(err)
	}
	if user.AuthVersion != 2 {
		t.Fatalf("auth_version=%d, want 2", user.AuthVersion)
	}
	var identity struct {
		CredentialHash string `gorm:"column:credential_hash"`
	}
	if err := business.Table("qixi_crm_b_user_identity").Select("credential_hash").Where("user_id=? AND channel='pc'", statusAcceptanceUserID).Take(&identity).Error; err != nil {
		t.Fatal(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(identity.CredentialHash), []byte(oldPassword)) == nil {
		t.Fatal("old local password remains valid")
	}
	if bcrypt.CompareHashAndPassword([]byte(identity.CredentialHash), []byte(command["password"].(string))) != nil {
		t.Fatal("new local password is not valid")
	}
	var audits int64
	if err := business.Table("qixi_crm_b_user_admin_command_audit").Where("action='password_reset' AND user_id=?", statusAcceptanceUserID).Count(&audits).Error; err != nil {
		t.Fatal(err)
	}
	if audits != 1 {
		t.Fatalf("password reset audit count=%d, want 1", audits)
	}
}

func TestStatusIntegrationDisableAndRBAC(t *testing.T) {
	business, admin := statusAcceptanceDB(t)
	statusAcceptancePrepareAdmin(t, admin)
	statusAcceptancePrepareUser(t, business)
	router, platformToken, platformAdmin := statusAcceptanceRouter(t, business, admin)

	for _, role := range []string{"merchant", "region", "customer_service", "operations"} {
		jwt := authjwt.NewManager(statusAcceptanceJWTSecret, time.Hour, time.Hour)
		forged, err := jwt.IssueAdminConsoleWithIdentityVersion(uint(platformAdmin.ID), platformAdmin.Username, []string{role}, platformAdmin.DataScopeVersion, platformAdmin.AuthVersion)
		if err != nil {
			t.Fatal(err)
		}
		response := statusAcceptanceRequest(router, forged.AccessToken, gin.H{"status": 0, "reason": "虚构中文角色越权", "idempotency_key": "status-role-" + role})
		if response.Code != http.StatusForbidden {
			t.Fatalf("%s status command=%d, want 403: %s", role, response.Code, response.Body.String())
		}
	}

	command := gin.H{"status": 0, "reason": "虚构中文风控停用验收", "idempotency_key": "status-disable-concurrent-987650001"}
	const callers = 8
	results := make(chan *httptest.ResponseRecorder, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- statusAcceptanceRequest(router, platformToken, command)
		}()
	}
	wg.Wait()
	close(results)
	created, replayed := 0, 0
	for response := range results {
		if response.Code != http.StatusOK {
			t.Fatalf("concurrent status command=%d: %s", response.Code, response.Body.String())
		}
		if strings.Contains(response.Body.String(), `"replayed":true`) {
			replayed++
		} else {
			created++
		}
	}
	if created != 1 || replayed != callers-1 {
		t.Fatalf("created=%d replayed=%d, want 1/%d", created, replayed, callers-1)
	}

	var user struct {
		Status      int    `gorm:"column:status"`
		AuthVersion uint64 `gorm:"column:auth_version"`
	}
	if err := business.Table("qixi_crm_b_user").Select("status,auth_version").Where("id=?", statusAcceptanceUserID).Take(&user).Error; err != nil {
		t.Fatal(err)
	}
	if user.Status != 0 || user.AuthVersion != 2 {
		t.Fatalf("disabled user=%+v, want status=0 auth_version=2", user)
	}
	var audits int64
	if err := business.Table("qixi_crm_b_user_status_command_audit").Where("user_id=?", statusAcceptanceUserID).Count(&audits).Error; err != nil {
		t.Fatal(err)
	}
	if audits != 1 {
		t.Fatalf("status audit count=%d, want 1", audits)
	}

	conflict := statusAcceptanceRequest(router, platformToken, gin.H{"status": 0, "reason": "不同中文原因必须冲突", "idempotency_key": "status-disable-concurrent-987650001"})
	if conflict.Code != http.StatusConflict {
		t.Fatalf("changed replay status=%d, want 409: %s", conflict.Code, conflict.Body.String())
	}
}

func TestStatusIntegrationReenable(t *testing.T) {
	business, admin := statusAcceptanceDB(t)
	statusAcceptancePrepareAdmin(t, admin)
	router, platformToken, _ := statusAcceptanceRouter(t, business, admin)
	response := statusAcceptanceRequest(router, platformToken, gin.H{"status": 1, "reason": "虚构中文复核恢复验收", "idempotency_key": "status-enable-987650001"})
	if response.Code != http.StatusOK {
		t.Fatalf("reenable status=%d: %s", response.Code, response.Body.String())
	}
	var user struct {
		Status      int    `gorm:"column:status"`
		AuthVersion uint64 `gorm:"column:auth_version"`
	}
	if err := business.Table("qixi_crm_b_user").Select("status,auth_version").Where("id=?", statusAcceptanceUserID).Take(&user).Error; err != nil {
		t.Fatal(err)
	}
	if user.Status != 1 || user.AuthVersion != 3 {
		t.Fatalf("reenabled user=%+v, want status=1 auth_version=3", user)
	}
}

func TestNotificationIntegrationRBACAndReplay(t *testing.T) {
	business, admin := statusAcceptanceDB(t)
	statusAcceptancePrepareAdmin(t, admin)
	statusAcceptancePrepareUser(t, business)
	router, platformToken, platformAdmin := statusAcceptanceRouter(t, business, admin)

	command := gin.H{
		"title":           "中文站内通知验收",
		"body":            "这是仅用于本地验收的中文图文提醒，不包含真实个人信息或外部凭据。",
		"cover_url":       "/demo/notification/中文封面.png",
		"reason":          "虚构中文服务工单发送",
		"idempotency_key": "notification-send-987650001",
	}
	for _, role := range []string{"merchant", "region", "customer_service", "operations"} {
		jwt := authjwt.NewManager(statusAcceptanceJWTSecret, time.Hour, time.Hour)
		pair, err := jwt.IssueAdminConsoleWithIdentityVersion(uint(platformAdmin.ID), platformAdmin.Username, []string{role}, platformAdmin.DataScopeVersion, platformAdmin.AuthVersion)
		if err != nil {
			t.Fatal(err)
		}
		if response := notificationAcceptanceRequest(router, pair.AccessToken, statusAcceptanceUserID, command); response.Code != http.StatusForbidden {
			t.Fatalf("%s notification=%d, want 403: %s", role, response.Code, response.Body.String())
		}
	}
	if response := notificationAcceptanceRequest(router, platformToken, 999999999, command); response.Code != http.StatusNotFound {
		t.Fatalf("missing notification=%d, want 404: %s", response.Code, response.Body.String())
	}

	const callers = 6
	responses := make(chan *httptest.ResponseRecorder, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			responses <- notificationAcceptanceRequest(router, platformToken, statusAcceptanceUserID, command)
		}()
	}
	wg.Wait()
	close(responses)
	created, replayed := 0, 0
	for response := range responses {
		if response.Code != http.StatusOK {
			t.Fatalf("concurrent notification=%d: %s", response.Code, response.Body.String())
		}
		if strings.Contains(response.Body.String(), `"replayed":true`) {
			replayed++
		} else {
			created++
		}
	}
	if created != 1 || replayed != callers-1 {
		t.Fatalf("notification created/replayed=%d/%d, want 1/%d", created, replayed, callers-1)
	}
	var notification struct {
		ID       uint64     `gorm:"column:id"`
		Category string     `gorm:"column:category"`
		Title    string     `gorm:"column:title"`
		Body     string     `gorm:"column:body"`
		Payload  string     `gorm:"column:payload"`
		ReadAt   *time.Time `gorm:"column:read_at"`
	}
	if err := business.Table("qixi_crm_b_notification").Select("id,category,title,body,payload,read_at").Where("user_id=? AND category=? AND title=?", statusAcceptanceUserID, "system", command["title"]).Take(&notification).Error; err != nil {
		t.Fatal(err)
	}
	var payload map[string]string
	if err := json.Unmarshal([]byte(notification.Payload), &payload); err != nil {
		t.Fatalf("decode notification payload: %v", err)
	}
	if notification.Category != "system" || notification.Body != command["body"] || notification.ReadAt != nil || payload["kind"] != "image_text" || payload["cover_url"] != "/demo/notification/中文封面.png" || payload["source"] != "platform_manual" {
		t.Fatalf("unexpected notification=%+v", notification)
	}
	var audits int64
	if err := business.Table("qixi_crm_b_user_notification_command_audit").Where("user_id=? AND notification_id=? AND title=? AND reason=?", statusAcceptanceUserID, notification.ID, command["title"], command["reason"]).Count(&audits).Error; err != nil {
		t.Fatal(err)
	}
	if audits != 1 {
		t.Fatalf("notification audit count=%d, want 1", audits)
	}
	if response := notificationAcceptanceRequest(router, platformToken, statusAcceptanceUserID, gin.H{"title": command["title"], "body": command["body"], "cover_url": command["cover_url"], "reason": "不同中文原因必须冲突", "idempotency_key": command["idempotency_key"]}); response.Code != http.StatusConflict {
		t.Fatalf("changed notification replay=%d, want 409: %s", response.Code, response.Body.String())
	}
	if response := notificationAcceptanceRequest(router, platformToken, statusAcceptanceUserID, gin.H{"title": command["title"], "body": command["body"], "cover_url": "/untrusted/cover.png", "reason": "虚构中文非法封面", "idempotency_key": "notification-invalid-cover-987650001"}); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid cover notification=%d, want 400: %s", response.Code, response.Body.String())
	}
	var finalCount int64
	if err := business.Table("qixi_crm_b_notification").Where("user_id=? AND category=? AND title=?", statusAcceptanceUserID, "system", command["title"]).Count(&finalCount).Error; err != nil {
		t.Fatal(err)
	}
	if finalCount != 1 {
		t.Fatalf("notification count=%d, want 1", finalCount)
	}
}

func TestExportIntegrationRBACTruncationAndCSVSafety(t *testing.T) {
	business, admin := statusAcceptanceDB(t)
	statusAcceptancePrepareAdmin(t, admin)
	if err := business.Exec("DELETE FROM qixi_crm_b_user WHERE id BETWEEN ? AND ?", statusAcceptanceExportStartID, statusAcceptanceExportEndID).Error; err != nil {
		t.Fatal(err)
	}
	if err := business.Table("qixi_crm_b_user_export_audit").Where("reason LIKE ?", "虚构中文导出验收%").Delete(nil).Error; err != nil {
		t.Fatal(err)
	}
	fixtures := make([]map[string]any, 0, 5001)
	for index := uint64(0); index <= 5000; index++ {
		id := statusAcceptanceExportStartID + index
		nickname := "导出截断验收用户" + strconv.FormatUint(index, 10)
		mobile := any(nil)
		if index == 5000 {
			nickname, mobile = "=导出截断验收公式", "13912345678"
		}
		fixtures = append(fixtures, map[string]any{"id": id, "nickname": nickname, "mobile": mobile, "status": 1, "auth_version": 1})
	}
	if err := business.Table("qixi_crm_b_user").CreateInBatches(fixtures, 500).Error; err != nil {
		t.Fatal(err)
	}
	router, platformToken, platformAdmin := statusAcceptanceRouter(t, business, admin)
	command := gin.H{"keyword": "导出截断验收", "status": 1, "reason": "虚构中文导出验收：筛选与脱敏"}
	for _, role := range []string{"merchant", "region", "customer_service", "operations"} {
		jwt := authjwt.NewManager(statusAcceptanceJWTSecret, time.Hour, time.Hour)
		pair, err := jwt.IssueAdminConsoleWithIdentityVersion(uint(platformAdmin.ID), platformAdmin.Username, []string{role}, platformAdmin.DataScopeVersion, platformAdmin.AuthVersion)
		if err != nil {
			t.Fatal(err)
		}
		if response := exportAcceptanceRequest(router, pair.AccessToken, command); response.Code != http.StatusForbidden {
			t.Fatalf("%s export=%d, want 403", role, response.Code)
		}
	}
	response := exportAcceptanceRequest(router, platformToken, command)
	if response.Code != http.StatusOK {
		t.Fatalf("export=%d: %s", response.Code, response.Body.String())
	}
	var result struct {
		Data struct {
			Content   string `json:"content"`
			RowCount  int    `json:"row_count"`
			Truncated bool   `json:"truncated"`
		} `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Data.RowCount != 5000 || !result.Data.Truncated || !strings.HasPrefix(result.Data.Content, "\ufeff") {
		t.Fatalf("unexpected export metadata=%+v", result.Data)
	}
	reader := csv.NewReader(strings.NewReader(strings.TrimPrefix(result.Data.Content, "\ufeff")))
	records, err := reader.ReadAll()
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if len(records) != 5001 || records[0][0] != "用户ID" || records[1][1] != "'=导出截断验收公式" || records[1][2] != "139****5678" {
		t.Fatalf("unexpected CSV export header/first row=%v/%v", records[0], records[1])
	}
	var audits int64
	if err := business.Table("qixi_crm_b_user_export_audit").Where("reason=? AND row_count=? AND operator_admin_id=?", command["reason"], 5000, statusAcceptanceAdminID).Count(&audits).Error; err != nil {
		t.Fatal(err)
	}
	if audits != 1 {
		t.Fatalf("export audit=%d, want 1", audits)
	}
}

// TestStatusIntegrationCleanup is intentionally invoked explicitly after the
// multi-process acceptance sequence. It removes only the fixed local test IDs.
func TestStatusIntegrationCleanup(t *testing.T) {
	business, admin := statusAcceptanceDB(t)
	if err := business.Table("qixi_crm_b_user_export_audit").Where("reason LIKE ?", "虚构中文导出验收%").Delete(nil).Error; err != nil {
		t.Fatalf("clear export audit: %v", err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user WHERE id BETWEEN ? AND ?", statusAcceptanceExportStartID, statusAcceptanceExportEndID).Error; err != nil {
		t.Fatalf("clear export users: %v", err)
	}
	if err := business.Table("qixi_crm_b_user_notification_command_audit").Where("user_id=?", statusAcceptanceUserID).Delete(nil).Error; err != nil {
		t.Fatalf("clear notification acceptance audit: %v", err)
	}
	if err := business.Table("qixi_crm_b_notification").Where("user_id=? AND category=? AND title IN ?", statusAcceptanceUserID, "system", []string{"中文站内通知验收", "中文站内通知页面验收"}).Delete(nil).Error; err != nil {
		t.Fatalf("clear notification acceptance fixture: %v", err)
	}
	if err := business.Table("qixi_crm_b_distribution_promoter_command_audit").Where("idempotency_key LIKE ?", "promoter-%987650001%").Delete(nil).Error; err != nil {
		t.Fatalf("clear promoter acceptance audit: %v", err)
	}
	if err := business.Table("qixi_crm_b_distribution_promoter").Where("user_id IN ?", []uint64{statusAcceptanceUserID, statusAcceptancePromoterUserID}).Delete(nil).Error; err != nil {
		t.Fatalf("clear promoter status: %v", err)
	}
	if err := business.Table("qixi_crm_b_commission_ledger").Where("idempotency_key=?", "acceptance-promoter-commission-987650001").Delete(nil).Error; err != nil {
		t.Fatalf("clear promoter commission fixture: %v", err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user WHERE id=?", statusAcceptancePromoterUserID).Error; err != nil {
		t.Fatalf("clear promoter acceptance user: %v", err)
	}
	var createdIDs []uint64
	if err := business.Table("qixi_crm_b_user_identity").Where("channel='pc' AND subject=?", statusAcceptanceCreatedAccount).Pluck("user_id", &createdIDs).Error; err != nil {
		t.Fatalf("find created acceptance user: %v", err)
	}
	if len(createdIDs) > 0 {
		if err := business.Exec("DELETE FROM qixi_crm_b_user_admin_command_audit WHERE user_id IN ?", createdIDs).Error; err != nil {
			t.Fatalf("clear created user audit: %v", err)
		}
		if err := business.Exec("DELETE FROM qixi_crm_b_user_profile WHERE user_id IN ?", createdIDs).Error; err != nil {
			t.Fatalf("clear created user profile: %v", err)
		}
		if err := business.Exec("DELETE FROM qixi_crm_b_user_identity WHERE user_id IN ?", createdIDs).Error; err != nil {
			t.Fatalf("clear created user identity: %v", err)
		}
		if err := business.Exec("DELETE FROM qixi_crm_b_user WHERE id IN ?", createdIDs).Error; err != nil {
			t.Fatalf("clear created user: %v", err)
		}
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user_status_command_audit WHERE user_id=?", statusAcceptanceUserID).Error; err != nil {
		t.Fatalf("clear status acceptance audit: %v", err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user_admin_command_audit WHERE user_id=?", statusAcceptanceUserID).Error; err != nil {
		t.Fatalf("clear user-admin acceptance audit: %v", err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user_identity WHERE user_id=?", statusAcceptanceUserID).Error; err != nil {
		t.Fatalf("clear user identity: %v", err)
	}
	if err := business.Exec("DELETE FROM qixi_crm_b_user WHERE id=?", statusAcceptanceUserID).Error; err != nil {
		t.Fatalf("clear status acceptance user: %v", err)
	}
	if err := admin.Exec("DELETE FROM qixi_crm_a_admin_user_role WHERE admin_user_id=?", statusAcceptanceAdminID).Error; err != nil {
		t.Fatalf("clear status acceptance role: %v", err)
	}
	if err := admin.Exec("DELETE FROM qixi_crm_a_admin_user WHERE id=?", statusAcceptanceAdminID).Error; err != nil {
		t.Fatalf("clear status acceptance admin: %v", err)
	}
}
