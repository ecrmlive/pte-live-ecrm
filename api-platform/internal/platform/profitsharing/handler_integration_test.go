package profitsharing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const profitsharingAcceptanceNo = "ECRM-PS-987680001"

func TestProfitsharingReviewIntegrationConcurrencyAndAudit(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_PROFITSHARING_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_PROFITSHARING_TEST_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		var ids []uint
		db.Table(tab).Where("application_no = ?", profitsharingAcceptanceNo).Pluck("id", &ids)
		if len(ids) > 0 {
			db.Table("qixi_crm_a_merchant_profitsharing_audit").Where("application_id IN ?", ids).Delete(nil)
		}
		db.Table(tab).Where("application_no = ?", profitsharingAcceptanceNo).Delete(nil)
	}
	cleanup()
	defer cleanup()
	if err := db.Table(tab).Create(map[string]any{
		"merchant_id":    987680002,
		"application_no": profitsharingAcceptanceNo,
		"status":         "applied",
		"description":    "中文模拟分账申请：仅验收审核状态机，不包含渠道账户信息",
	}).Error; err != nil {
		t.Fatal(err)
	}
	var applicationID uint
	if err := db.Table(tab).Where("application_no = ?", profitsharingAcceptanceNo).Pluck("id", &applicationID).Error; err != nil || applicationID == 0 {
		t.Fatalf("load application id=%d err=%v", applicationID, err)
	}

	gin.SetMode(gin.TestMode)
	h := NewHandler(db)
	r := gin.New()
	r.POST("/applications/:id/review", func(c *gin.Context) {
		c.Set(middleware.CtxAdminID, uint(987680013))
		h.Review(c)
	})
	r.PUT("/applications/:id/note", func(c *gin.Context) {
		c.Set(middleware.CtxAdminID, uint(987680013))
		h.Note(c)
	})
	call := func(method, path string, input any) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(input)
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	path := "/applications/" + strconv.FormatUint(uint64(applicationID), 10) + "/review"
	var wg sync.WaitGroup
	statuses := make(chan int, 6)
	for range 6 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			statuses <- call(http.MethodPost, path, gin.H{"approved": true, "note": "中文模拟审核通过：资料完整"}).Code
		}()
	}
	wg.Wait()
	close(statuses)
	var succeeded, conflicted int
	for status := range statuses {
		switch status {
		case http.StatusOK:
			succeeded++
		case http.StatusConflict:
			conflicted++
		default:
			t.Fatalf("unexpected review status=%d", status)
		}
	}
	if succeeded != 1 || conflicted != 5 {
		t.Fatalf("review concurrency succeeded=%d conflicted=%d", succeeded, conflicted)
	}
	var application struct {
		Status     string
		ReviewNote string `gorm:"column:review_note"`
		ReviewedBy uint   `gorm:"column:reviewed_by"`
	}
	if err := db.Table(tab).Where("id = ?", applicationID).Take(&application).Error; err != nil {
		t.Fatal(err)
	}
	if application.Status != "approved" || application.ReviewNote != "中文模拟审核通过：资料完整" || application.ReviewedBy != 987680013 {
		t.Fatalf("application=%+v", application)
	}
	var auditCount int64
	if err := db.Table("qixi_crm_a_merchant_profitsharing_audit").Where("application_id = ? AND from_status = 'applied' AND to_status = 'approved' AND note = ?", applicationID, "中文模拟审核通过：资料完整").Count(&auditCount).Error; err != nil || auditCount != 1 {
		t.Fatalf("audit=%d err=%v", auditCount, err)
	}
	if w := call(http.MethodPut, "/applications/"+strconv.FormatUint(uint64(applicationID), 10)+"/note", gin.H{"note": "中文模拟内部备注：已归档"}); w.Code != http.StatusOK {
		t.Fatalf("note=%d %s", w.Code, w.Body.String())
	}
	if w := call(http.MethodPost, path, gin.H{"approved": true, "note": "中文模拟重复审核"}); w.Code != http.StatusConflict {
		t.Fatalf("repeat review=%d %s", w.Code, w.Body.String())
	}
	if w := call(http.MethodPut, "/applications/"+strconv.FormatUint(uint64(applicationID), 10)+"/note", gin.H{"note": strings.Repeat("审", 501)}); w.Code != http.StatusBadRequest {
		t.Fatalf("oversize note=%d %s", w.Code, w.Body.String())
	}
	var retainedAuditNote string
	if err := db.Table("qixi_crm_a_merchant_profitsharing_audit").Where("application_id = ?", applicationID).Pluck("note", &retainedAuditNote).Error; err != nil || retainedAuditNote != "中文模拟审核通过：资料完整" {
		t.Fatalf("retained audit note=%q err=%v", retainedAuditNote, err)
	}
}
