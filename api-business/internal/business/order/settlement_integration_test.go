package order

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestConfirmReceiptIntegrationEnqueuesFactAnchoredSettlementAccrual proves
// the user-facing completion endpoint changes the order and records its
// merchant-ledger command in one business transaction. It needs an isolated
// MySQL schema because the contract is intentionally database-backed.
func TestConfirmReceiptIntegrationEnqueuesFactAnchoredSettlementAccrual(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_SETTLEMENT_BUSINESS_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_SETTLEMENT_BUSINESS_TEST_DSN to run settlement receipt integration test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated business database: %v", err)
	}
	const orderID uint64 = 987677901
	const userID uint64 = 987677902
	const storeID uint64 = 987677903
	const merchantID uint64 = 987677904
	const orderNo = "SETTLEMENT-RECEIPT-987677901"
	ctx := context.Background()
	cleanup := func() {
		_ = db.WithContext(ctx).Table("qixi_crm_b_settlement_command_outbox").Where("order_id = ?", orderID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_order").Where("id = ? AND order_no = ?", orderID, orderNo).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	if err := db.WithContext(ctx).Table("qixi_crm_b_order").Create(map[string]any{
		"id":                     orderID,
		"group_order_id":         orderID,
		"order_no":               orderNo,
		"merchant_id":            merchantID,
		"merchant_name_snapshot": "七禧结算事实验收茶铺",
		"store_id":               storeID,
		"store_name_snapshot":    "七禧结算事实验收店",
		"user_id":                userID,
		"total_amount":           128.5,
		"pay_amount":             128.5,
		"total_quantity":         1,
		"recipient_snapshot":     `{}`,
		"status":                 "shipped",
	}).Error; err != nil {
		t.Fatalf("seed Chinese shipped order: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.CtxUID, uint(userID))
		c.Next()
	})
	NewHandler(db, nil, true).Register(router)
	request := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/order/987677901/confirm-receipt", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		return resp
	}
	if resp := request(); resp.Code != http.StatusOK {
		t.Fatalf("first confirm receipt status=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp := request(); resp.Code != http.StatusOK {
		t.Fatalf("receipt replay status=%d body=%s", resp.Code, resp.Body.String())
	}
	var state struct {
		Status string `gorm:"column:status"`
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_order").Select("status").Where("id = ?", orderID).Scan(&state).Error; err != nil || state.Status != "completed" {
		t.Fatalf("order status=%q err=%v, want completed", state.Status, err)
	}
	var commands []struct {
		Action         string  `gorm:"column:action"`
		StoreID        uint64  `gorm:"column:store_id"`
		MerchantID     uint64  `gorm:"column:merchant_id"`
		Amount         float64 `gorm:"column:amount"`
		IdempotencyKey string  `gorm:"column:idempotency_key"`
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_settlement_command_outbox").Where("order_id = ?", orderID).Find(&commands).Error; err != nil {
		t.Fatalf("load settlement command: %v", err)
	}
	if len(commands) != 1 || commands[0].Action != "accrue" || commands[0].StoreID != storeID || commands[0].MerchantID != merchantID || commands[0].Amount != 128.5 || commands[0].IdempotencyKey != "settlement:accrue:987677901" {
		t.Fatalf("unexpected receipt settlement command: %+v", commands)
	}
}
