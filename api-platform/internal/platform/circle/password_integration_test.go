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

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/circle"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestResetAgentPasswordIntegrationSingleSubmit(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_CIRCLE_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_CIRCLE_TEST_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	const agentID, adminID = 987670011, 987670012
	defer func() {
		db.Exec("DELETE FROM qixi_crm_a_business_zone_agent_password_reset_audit WHERE circle_agent_id=?", agentID)
		db.Exec("DELETE FROM qixi_crm_a_admin_user WHERE id=?", adminID)
		db.Exec("DELETE FROM qixi_crm_a_business_zone_agent WHERE circle_agent_id=?", agentID)
	}()
	db.Exec("DELETE FROM qixi_crm_a_business_zone_agent_password_reset_audit WHERE circle_agent_id=?", agentID)
	db.Exec("DELETE FROM qixi_crm_a_admin_user WHERE id=?", adminID)
	db.Exec("DELETE FROM qixi_crm_a_business_zone_agent WHERE circle_agent_id=?", agentID)
	old := "LocalOldPassword12"
	oldHash, _ := bcrypt.GenerateFromPassword([]byte(old), bcrypt.DefaultCost)
	now := time.Now()
	if err := db.Create(&circle.Agent{CircleAgentID: agentID, Name: "中文口令验收代理", Phone: "13800000011", Status: circle.AgentApproved, CreateTime: now, UpdateTime: now}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO qixi_crm_a_admin_user (id,username,password_hash,display_name,status,auth_version,data_scope_version,circle_agent_id) VALUES (?,?,?,?,1,1,1,?)", adminID, "agent-password-local", string(oldHash), "中文代理后台账号", agentID).Error; err != nil {
		t.Fatal(err)
	}
	h := NewHandler(nil, db)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/reset/:id", func(c *gin.Context) { c.Set(middleware.CtxAdminID, uint(987670013)); h.ResetAgentPassword(c) })
	body := gin.H{"password": "LocalNewPassword34", "reason": "虚构中文密码重置验收", "idempotency_key": "agent-password-integration-001"}
	raw, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/reset/987670011", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("reset=%d %s", w.Code, w.Body.String())
	}
	var admin struct {
		Hash    string `gorm:"column:password_hash"`
		Version uint64 `gorm:"column:auth_version"`
	}
	if err := db.Table("qixi_crm_a_admin_user").Select("password_hash,auth_version").Where("id=?", adminID).Take(&admin).Error; err != nil {
		t.Fatal(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(admin.Hash), []byte(old)) == nil || bcrypt.CompareHashAndPassword([]byte(admin.Hash), []byte(body["password"].(string))) != nil || admin.Version != 2 {
		t.Fatal("password/auth version not updated")
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/reset/987670011", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("retry=%d", w.Code)
	}
}
