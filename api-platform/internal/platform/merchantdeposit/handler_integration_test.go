package merchantdeposit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const depositAcceptanceMerchantID uint = 987680001
const depositDeductAcceptanceMerchantID uint = 987680004
const depositReviewAcceptanceMerchantID uint = 987680005

func TestMarkPaidIntegrationIdempotencyAndBalance(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_DEPOSIT_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_DEPOSIT_TEST_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		db.Table("qixi_crm_a_merchant_deposit_ledger").Where("merchant_id = ?", depositAcceptanceMerchantID).Delete(nil)
		db.Table("qixi_crm_a_merchant_deposit_refund").Where("merchant_id = ?", depositAcceptanceMerchantID).Delete(nil)
		db.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id = ?", depositAcceptanceMerchantID).Delete(nil)
	}
	cleanup()
	defer cleanup()
	if err := db.Table("qixi_crm_a_merchant_deposit_account").Create(map[string]any{
		"merchant_id":      depositAcceptanceMerchantID,
		"required_amount":  500,
		"available_amount": 500,
		"state":            "funded",
		"version":          1,
	}).Error; err != nil {
		t.Fatal(err)
	}
	createRefund := func(reason string) uint {
		if err := db.Table("qixi_crm_a_merchant_deposit_refund").Create(map[string]any{
			"merchant_id": depositAcceptanceMerchantID,
			"amount":      100,
			"status":      "approved",
			"reason":      reason,
		}).Error; err != nil {
			t.Fatal(err)
		}
		var id uint
		if err := db.Table("qixi_crm_a_merchant_deposit_refund").Where("merchant_id = ? AND reason = ?", depositAcceptanceMerchantID, reason).Order("id DESC").Limit(1).Pluck("id", &id).Error; err != nil || id == 0 {
			t.Fatalf("load refund id=%d err=%v", id, err)
		}
		return id
	}

	gin.SetMode(gin.TestMode)
	h := NewHandler(db)
	r := gin.New()
	r.POST("/refunds/:id/mark-paid", func(c *gin.Context) {
		c.Set(middleware.CtxAdminID, uint(987680013))
		h.MarkPaid(c)
	})
	call := func(id uint, key, reference string) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(gin.H{"idempotency_key": key, "payout_reference": reference})
		req := httptest.NewRequest(http.MethodPost, "/refunds/"+strconv.FormatUint(uint64(id), 10)+"/mark-paid", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	first := createRefund("中文模拟退款：首笔打款登记")
	key := "deposit-paid-987680001"
	if w := call(first, key, "本地模拟打款凭证-001"); w.Code != http.StatusOK {
		t.Fatalf("first mark paid=%d %s", w.Code, w.Body.String())
	}
	if w := call(first, key, "本地模拟打款凭证-001"); w.Code != http.StatusOK {
		t.Fatalf("same command replay=%d %s", w.Code, w.Body.String())
	}
	if w := call(first, key, "本地模拟篡改凭证-002"); w.Code != http.StatusConflict {
		t.Fatalf("changed reference=%d %s", w.Code, w.Body.String())
	}

	second := createRefund("中文模拟退款：幂等键冲突")
	if w := call(second, key, "本地模拟打款凭证-001"); w.Code != http.StatusConflict {
		t.Fatalf("cross refund key reuse=%d %s", w.Code, w.Body.String())
	}
	var accountRow struct {
		Available float64 `gorm:"column:available_amount"`
		State     string
		Version   uint
	}
	if err := db.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id = ?", depositAcceptanceMerchantID).Take(&accountRow).Error; err != nil {
		t.Fatal(err)
	}
	if accountRow.Available != 400 || accountRow.State != "shortfall" || accountRow.Version != 2 {
		t.Fatalf("account after replay/conflict = %+v", accountRow)
	}
	var paidCount, ledgerCount int64
	db.Table("qixi_crm_a_merchant_deposit_refund").Where("merchant_id = ? AND status = 'paid'", depositAcceptanceMerchantID).Count(&paidCount)
	db.Table("qixi_crm_a_merchant_deposit_ledger").Where("merchant_id = ? AND entry_type = 'refund_paid'", depositAcceptanceMerchantID).Count(&ledgerCount)
	if paidCount != 1 || ledgerCount != 1 {
		t.Fatalf("paid=%d ledger=%d", paidCount, ledgerCount)
	}
	var secondStatus string
	if err := db.Table("qixi_crm_a_merchant_deposit_refund").Where("id = ?", second).Pluck("status", &secondStatus).Error; err != nil || secondStatus != "approved" {
		t.Fatalf("second refund status=%q err=%v", secondStatus, err)
	}
	if err := db.Table("qixi_crm_a_merchant_deposit_refund").Create(map[string]any{"merchant_id": depositAcceptanceMerchantID, "amount": 500, "status": "approved", "reason": "中文模拟退款：余额不足登记"}).Error; err != nil {
		t.Fatal(err)
	}
	var insufficientRefundID uint
	if err := db.Table("qixi_crm_a_merchant_deposit_refund").Where("merchant_id = ? AND reason = ?", depositAcceptanceMerchantID, "中文模拟退款：余额不足登记").Pluck("id", &insufficientRefundID).Error; err != nil || insufficientRefundID == 0 {
		t.Fatalf("insufficient refund id=%d err=%v", insufficientRefundID, err)
	}
	if w := call(insufficientRefundID, "deposit-paid-insufficient-987680001", "本地模拟余额不足凭证"); w.Code != http.StatusConflict {
		t.Fatalf("insufficient mark paid=%d %s", w.Code, w.Body.String())
	}
	if err := db.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id = ?", depositAcceptanceMerchantID).Take(&accountRow).Error; err != nil || accountRow.Available != 400 || accountRow.Version != 2 {
		t.Fatalf("insufficient mark paid mutated account=%+v err=%v", accountRow, err)
	}
}

func TestDeductIntegrationIdempotencyFingerprintAndCentBoundary(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_DEPOSIT_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_DEPOSIT_TEST_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		db.Table("qixi_crm_a_merchant_deposit_ledger").Where("merchant_id = ?", depositDeductAcceptanceMerchantID).Delete(nil)
		db.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id = ?", depositDeductAcceptanceMerchantID).Delete(nil)
	}
	cleanup()
	defer cleanup()
	if err := db.Table("qixi_crm_a_merchant_deposit_account").Create(map[string]any{
		"merchant_id":      depositDeductAcceptanceMerchantID,
		"required_amount":  500,
		"available_amount": 500,
		"state":            "funded",
		"version":          1,
	}).Error; err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	h := NewHandler(db)
	r := gin.New()
	r.POST("/deposits/:merchant_id/deduct", func(c *gin.Context) {
		c.Set(middleware.CtxAdminID, uint(987680013))
		h.Deduct(c)
	})
	call := func(amount float64, reason string) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(gin.H{"amount": amount, "reason": reason, "idempotency_key": "deposit-deduct-987680004"})
		req := httptest.NewRequest(http.MethodPost, "/deposits/"+strconv.FormatUint(uint64(depositDeductAcceptanceMerchantID), 10)+"/deduct", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	if w := call(100, "中文模拟违规扣减"); w.Code != http.StatusOK {
		t.Fatalf("first deduct=%d %s", w.Code, w.Body.String())
	}
	if w := call(100, "中文模拟违规扣减"); w.Code != http.StatusOK {
		t.Fatalf("same command replay=%d %s", w.Code, w.Body.String())
	}
	if w := call(50, "中文模拟违规扣减"); w.Code != http.StatusConflict {
		t.Fatalf("changed amount=%d %s", w.Code, w.Body.String())
	}
	if w := call(100, "中文模拟篡改原因"); w.Code != http.StatusConflict {
		t.Fatalf("changed reason=%d %s", w.Code, w.Body.String())
	}
	if w := call(0.001, "中文模拟小数精度"); w.Code != http.StatusBadRequest {
		t.Fatalf("fractional cent=%d %s", w.Code, w.Body.String())
	}
	var accountRow struct {
		Available float64 `gorm:"column:available_amount"`
		Version   uint
	}
	if err := db.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id = ?", depositDeductAcceptanceMerchantID).Take(&accountRow).Error; err != nil {
		t.Fatal(err)
	}
	if accountRow.Available != 400 || accountRow.Version != 2 {
		t.Fatalf("account after idempotency conflicts=%+v", accountRow)
	}
	var ledgerCount int64
	db.Table("qixi_crm_a_merchant_deposit_ledger").Where("merchant_id = ? AND entry_type = 'deduct'", depositDeductAcceptanceMerchantID).Count(&ledgerCount)
	if ledgerCount != 1 {
		t.Fatalf("ledger=%d", ledgerCount)
	}
	insufficientRaw, _ := json.Marshal(gin.H{"amount": 401, "reason": "中文模拟余额不足扣减", "idempotency_key": "deposit-deduct-insufficient-987680004"})
	insufficientReq := httptest.NewRequest(http.MethodPost, "/deposits/"+strconv.FormatUint(uint64(depositDeductAcceptanceMerchantID), 10)+"/deduct", bytes.NewReader(insufficientRaw))
	insufficientReq.Header.Set("Content-Type", "application/json")
	insufficientRecorder := httptest.NewRecorder()
	r.ServeHTTP(insufficientRecorder, insufficientReq)
	if insufficientRecorder.Code != http.StatusConflict {
		t.Fatalf("insufficient deduct=%d %s", insufficientRecorder.Code, insufficientRecorder.Body.String())
	}
	var unchanged struct {
		Available float64 `gorm:"column:available_amount"`
		Version   uint
	}
	if err := db.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id = ?", depositDeductAcceptanceMerchantID).Take(&unchanged).Error; err != nil || unchanged.Available != 400 || unchanged.Version != 2 {
		t.Fatalf("insufficient deduction mutated account=%+v err=%v", unchanged, err)
	}
}

func TestReviewIntegrationConcurrencyAndImmutableLedger(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_DEPOSIT_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_DEPOSIT_TEST_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		db.Table("qixi_crm_a_merchant_deposit_ledger").Where("merchant_id = ?", depositReviewAcceptanceMerchantID).Delete(nil)
		db.Table("qixi_crm_a_merchant_deposit_refund").Where("merchant_id = ?", depositReviewAcceptanceMerchantID).Delete(nil)
		db.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id = ?", depositReviewAcceptanceMerchantID).Delete(nil)
	}
	cleanup()
	defer cleanup()
	if err := db.Table("qixi_crm_a_merchant_deposit_account").Create(map[string]any{"merchant_id": depositReviewAcceptanceMerchantID, "required_amount": 500, "available_amount": 500, "state": "funded", "version": 1}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Table("qixi_crm_a_merchant_deposit_refund").Create(map[string]any{"merchant_id": depositReviewAcceptanceMerchantID, "amount": 500, "status": "applied", "reason": "中文模拟退款审核并发"}).Error; err != nil {
		t.Fatal(err)
	}
	var refundID uint
	if err := db.Table("qixi_crm_a_merchant_deposit_refund").Where("merchant_id = ?", depositReviewAcceptanceMerchantID).Pluck("id", &refundID).Error; err != nil || refundID == 0 {
		t.Fatalf("refund id=%d err=%v", refundID, err)
	}
	gin.SetMode(gin.TestMode)
	h := NewHandler(db)
	r := gin.New()
	r.POST("/refunds/:id/approve", func(c *gin.Context) { c.Set(middleware.CtxAdminID, uint(987680013)); h.Approve(c) })
	r.POST("/refunds/:id/reject", func(c *gin.Context) { c.Set(middleware.CtxAdminID, uint(987680014)); h.Reject(c) })
	call := func(path, note string) int {
		raw, _ := json.Marshal(gin.H{"note": note})
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}
	statuses := make(chan int, 6)
	for index := 0; index < 6; index++ {
		go func(index int) {
			action := "approve"
			if index%2 == 1 {
				action = "reject"
			}
			statuses <- call("/refunds/"+strconv.FormatUint(uint64(refundID), 10)+"/"+action, "中文模拟审核并发说明")
		}(index)
	}
	success, conflict := 0, 0
	for index := 0; index < 6; index++ {
		switch code := <-statuses; code {
		case http.StatusOK:
			success++
		case http.StatusConflict:
			conflict++
		default:
			t.Fatalf("concurrent review=%d", code)
		}
	}
	if success != 1 || conflict != 5 {
		t.Fatalf("success=%d conflict=%d", success, conflict)
	}
	var refund struct{ Status string }
	if err := db.Table("qixi_crm_a_merchant_deposit_refund").Where("id = ?", refundID).Take(&refund).Error; err != nil {
		t.Fatal(err)
	}
	var accountRow struct {
		State   string
		Version uint
	}
	if err := db.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id = ?", depositReviewAcceptanceMerchantID).Take(&accountRow).Error; err != nil {
		t.Fatal(err)
	}
	var ledgerCount int64
	if err := db.Table("qixi_crm_a_merchant_deposit_ledger").Where("merchant_id = ? AND entry_type IN ?", depositReviewAcceptanceMerchantID, []string{"refund_approved", "refund_rejected"}).Count(&ledgerCount).Error; err != nil {
		t.Fatal(err)
	}
	if ledgerCount != 1 {
		t.Fatalf("review ledger=%d", ledgerCount)
	}
	if refund.Status == "approved" && (accountRow.State != "refund_pending" || accountRow.Version != 2) {
		t.Fatalf("approved account=%+v", accountRow)
	}
	if refund.Status == "rejected" && (accountRow.State != "funded" || accountRow.Version != 1) {
		t.Fatalf("rejected account=%+v", accountRow)
	}
	if code := call("/refunds/"+strconv.FormatUint(uint64(refundID), 10)+"/approve", ""); code != http.StatusBadRequest {
		t.Fatalf("blank review note=%d", code)
	}
}
