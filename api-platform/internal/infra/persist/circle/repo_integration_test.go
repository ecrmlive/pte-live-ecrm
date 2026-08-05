package circlepersist

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/circle"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const revokeAcceptanceAgentID uint = 987670001

func revokeAcceptanceDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("ECRM_CIRCLE_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_CIRCLE_TEST_DSN to run disposable MySQL circle acceptance")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func clearRevokeAcceptance(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, q := range []string{
		"DELETE FROM qixi_crm_a_business_zone_agent_command_audit WHERE circle_agent_id=987670001",
		"DELETE FROM qixi_crm_a_business_zone WHERE circle_agent_id=987670001",
		"DELETE FROM qixi_crm_a_admin_user WHERE circle_agent_id=987670001",
		"DELETE FROM qixi_crm_a_business_zone_agent WHERE circle_agent_id=987670001",
	} {
		if err := db.Exec(q).Error; err != nil {
			t.Fatal(err)
		}
	}
}

func prepareRevokeAgent(t *testing.T, db *gorm.DB, status int8, balance float64) {
	t.Helper()
	clearRevokeAcceptance(t, db)
	now := time.Now()
	if err := db.Create(&circle.Agent{CircleAgentID: revokeAcceptanceAgentID, Name: "中文撤销验收代理", Phone: "13800000001", Status: status, Balance: balance, CreateTime: now, UpdateTime: now}).Error; err != nil {
		t.Fatal(err)
	}
}

func TestRevokeAgentIntegrationGuardsAndReplay(t *testing.T) {
	db := revokeAcceptanceDB(t)
	defer clearRevokeAcceptance(t, db)
	repo := NewRepo(db)
	ctx := context.Background()
	prepareRevokeAgent(t, db, circle.AgentPending, 0)
	if _, err := repo.RevokeAgent(ctx, revokeAcceptanceAgentID, "虚构中文待审撤销", "revoke-pending-001", 1, time.Now()); !errors.Is(err, circle.ErrAgentNotApproved) {
		t.Fatalf("pending=%v", err)
	}
	prepareRevokeAgent(t, db, circle.AgentApproved, 9.99)
	if _, err := repo.RevokeAgent(ctx, revokeAcceptanceAgentID, "虚构中文余额撤销", "revoke-balance-001", 1, time.Now()); !errors.Is(err, circle.ErrAgentBalance) {
		t.Fatalf("balance=%v", err)
	}
	prepareRevokeAgent(t, db, circle.AgentApproved, 0)
	if err := db.Exec("INSERT INTO qixi_crm_a_business_zone (circle_id,name,circle_agent_id,status,type) VALUES (987670001,'中文关联区域',987670001,1,0)").Error; err != nil {
		t.Fatal(err)
	}
	if _, err := repo.RevokeAgent(ctx, revokeAcceptanceAgentID, "虚构中文区域关联", "revoke-zone-001", 1, time.Now()); !errors.Is(err, circle.ErrAgentBound) {
		t.Fatalf("zone=%v", err)
	}
	if err := db.Exec("DELETE FROM qixi_crm_a_business_zone WHERE circle_agent_id=?", revokeAcceptanceAgentID).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO qixi_crm_a_admin_user (id,username,password_hash,display_name,status,circle_agent_id) VALUES (987670002,'circle-revoke-local','not-a-real-password','中文绑定验收账号',1,?)", revokeAcceptanceAgentID).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := repo.RevokeAgent(ctx, revokeAcceptanceAgentID, "虚构中文后台绑定", "revoke-admin-001", 1, time.Now()); !errors.Is(err, circle.ErrAgentAdminBound) {
		t.Fatalf("admin binding=%v", err)
	}
	if err := db.Exec("DELETE FROM qixi_crm_a_admin_user WHERE circle_agent_id=?", revokeAcceptanceAgentID).Error; err != nil {
		t.Fatal(err)
	}
	const callers = 6
	results := make(chan bool, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			replay, err := repo.RevokeAgent(ctx, revokeAcceptanceAgentID, "虚构中文并发撤销", "revoke-concurrent-001", 1, time.Now())
			if err != nil {
				t.Errorf("revoke=%v", err)
				return
			}
			results <- replay
		}()
	}
	wg.Wait()
	close(results)
	replayed := 0
	for v := range results {
		if v {
			replayed++
		}
	}
	if replayed != callers-1 {
		t.Fatalf("replayed=%d", replayed)
	}
	if _, err := repo.RevokeAgent(ctx, revokeAcceptanceAgentID, "不同中文原因", "revoke-concurrent-001", 1, time.Now()); !errors.Is(err, circle.ErrCommandConflict) {
		t.Fatalf("reason=%v", err)
	}
	if _, err := repo.RevokeAgent(ctx, revokeAcceptanceAgentID, "虚构中文并发撤销", "revoke-concurrent-001", 2, time.Now()); !errors.Is(err, circle.ErrCommandConflict) {
		t.Fatalf("operator=%v", err)
	}
	var audit int64
	if err := db.Table("qixi_crm_a_business_zone_agent_command_audit").Where("circle_agent_id=?", revokeAcceptanceAgentID).Count(&audit).Error; err != nil || audit != 1 {
		t.Fatalf("audit=%d err=%v", audit, err)
	}
}
