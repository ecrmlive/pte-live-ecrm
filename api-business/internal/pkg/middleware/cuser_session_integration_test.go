package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/authjwt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const statusAcceptanceUserID uint64 = 987650001

const statusAcceptanceJWTSecret = "local-status-acceptance-not-a-production-secret"

func cUserSessionAcceptanceDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("ECRM_STATUS_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_STATUS_TEST_DSN to run the disposable MySQL acceptance test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open business acceptance database: %v", err)
	}
	return db
}

func cUserSessionAcceptanceRequest(router http.Handler, token string) int {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authori-zation", "Bearer "+token)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, req)
	return writer.Code
}

func cUserSessionAcceptanceRouter(db *gorm.DB) (*gin.Engine, *authjwt.Manager) {
	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager(statusAcceptanceJWTSecret, time.Hour, time.Hour)
	r := gin.New()
	r.GET("/protected", JWTRequired(jwt, authjwt.PortalApp), CUserSessionRequired(db), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	return r, jwt
}

func TestCUserSessionIntegrationDisabled(t *testing.T) {
	db := cUserSessionAcceptanceDB(t)
	var user struct {
		Status      int    `gorm:"column:status"`
		AuthVersion uint64 `gorm:"column:auth_version"`
	}
	if err := db.Table("qixi_crm_b_user").Select("status,auth_version").Where("id=?", statusAcceptanceUserID).Take(&user).Error; err != nil {
		t.Fatalf("read disabled acceptance user: %v", err)
	}
	if user.Status != 0 || user.AuthVersion != 2 {
		t.Fatalf("disabled state=%+v, want status=0 auth_version=2", user)
	}
	router, jwt := cUserSessionAcceptanceRouter(db)
	for _, version := range []uint64{1, 2} {
		pair, err := jwt.IssueCUserWithIdentityVersion(uint(statusAcceptanceUserID), "中文状态验收用户", "pc", version)
		if err != nil {
			t.Fatal(err)
		}
		if got := cUserSessionAcceptanceRequest(router, pair.AccessToken); got != http.StatusUnauthorized {
			t.Fatalf("disabled token version=%d status=%d, want 401", version, got)
		}
	}
}

func TestCUserSessionIntegrationReenabled(t *testing.T) {
	db := cUserSessionAcceptanceDB(t)
	var user struct {
		Status      int    `gorm:"column:status"`
		AuthVersion uint64 `gorm:"column:auth_version"`
	}
	if err := db.Table("qixi_crm_b_user").Select("status,auth_version").Where("id=?", statusAcceptanceUserID).Take(&user).Error; err != nil {
		t.Fatalf("read reenabled acceptance user: %v", err)
	}
	if user.Status != 1 || user.AuthVersion != 3 {
		t.Fatalf("reenabled state=%+v, want status=1 auth_version=3", user)
	}
	router, jwt := cUserSessionAcceptanceRouter(db)
	for _, version := range []uint64{1, 2} {
		pair, err := jwt.IssueCUserWithIdentityVersion(uint(statusAcceptanceUserID), "中文状态验收用户", "pc", version)
		if err != nil {
			t.Fatal(err)
		}
		if got := cUserSessionAcceptanceRequest(router, pair.AccessToken); got != http.StatusUnauthorized {
			t.Fatalf("old token version=%d status=%d, want 401", version, got)
		}
	}
	current, err := jwt.IssueCUserWithIdentityVersion(uint(statusAcceptanceUserID), "中文状态验收用户", "pc", 3)
	if err != nil {
		t.Fatal(err)
	}
	if got := cUserSessionAcceptanceRequest(router, current.AccessToken); got != http.StatusNoContent {
		t.Fatalf("current token status=%d, want 204", got)
	}
}

func TestCUserSessionIntegrationPasswordReset(t *testing.T) {
	db := cUserSessionAcceptanceDB(t)
	var user struct {
		Status      int    `gorm:"column:status"`
		AuthVersion uint64 `gorm:"column:auth_version"`
	}
	if err := db.Table("qixi_crm_b_user").Select("status,auth_version").Where("id=?", statusAcceptanceUserID).Take(&user).Error; err != nil {
		t.Fatalf("read password-reset acceptance user: %v", err)
	}
	if user.Status != 1 || user.AuthVersion != 2 {
		t.Fatalf("password-reset state=%+v, want status=1 auth_version=2", user)
	}
	router, jwt := cUserSessionAcceptanceRouter(db)
	oldPair, err := jwt.IssueCUserWithIdentityVersion(uint(statusAcceptanceUserID), "中文口令验收账号", "pc", 1)
	if err != nil {
		t.Fatal(err)
	}
	if got := cUserSessionAcceptanceRequest(router, oldPair.AccessToken); got != http.StatusUnauthorized {
		t.Fatalf("pre-reset token status=%d, want 401", got)
	}
	newPair, err := jwt.IssueCUserWithIdentityVersion(uint(statusAcceptanceUserID), "中文口令验收账号", "pc", 2)
	if err != nil {
		t.Fatal(err)
	}
	if got := cUserSessionAcceptanceRequest(router, newPair.AccessToken); got != http.StatusNoContent {
		t.Fatalf("post-reset token status=%d, want 204", got)
	}
}
