package order

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestSandboxRefundCallbackIntegrationRecordsServerDerivedReversal exercises
// the only local provider simulator. The request supplies a refund ID only;
// amount, store, merchant, SKU and quantity are all derived from locked facts.
// Production never registers this route.
func TestSandboxRefundCallbackIntegrationRecordsServerDerivedReversal(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_SETTLEMENT_BUSINESS_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_SETTLEMENT_BUSINESS_TEST_DSN to run refund callback integration test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated business database: %v", err)
	}
	const groupID uint64 = 987677921
	const orderID uint64 = 987677922
	const refundID uint64 = 987677923
	const orderItemID uint64 = 987677924
	const transactionID uint64 = 987677925
	const merchantID uint64 = 987677926
	const storeID uint64 = 987677927
	const skuID uint64 = 987677928
	const groupNo = "SETTLEMENT-REFUND-GROUP-987677921"
	const orderNo = "SETTLEMENT-REFUND-987677922"
	const refundNo = "SETTLEMENT-REFUND-987677923"
	ctx := context.Background()
	cleanup := func() {
		_ = db.WithContext(ctx).Table("qixi_crm_b_refund_event").Where("refund_id = ?", refundID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_refund_callback").Where("provider_event_id = ?", "mock-refund:987677923").Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_refund_transaction").Where("refund_id = ?", refundID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_aftersale_item").Where("refund_id = ?", refundID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_stock_command_outbox").Where("order_id = ?", orderID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_settlement_command_outbox").Where("order_id = ?", orderID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_refund").Where("id = ? AND refund_no = ?", refundID, refundNo).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_payment_transaction").Where("id = ?", transactionID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_order_item").Where("id = ?", orderItemID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_order").Where("id = ? AND order_no = ?", orderID, orderNo).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_group_order").Where("id = ? AND order_no = ?", groupID, groupNo).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	if err := db.WithContext(ctx).Table("qixi_crm_b_group_order").Create(map[string]any{
		"id": groupID, "order_no": groupNo, "user_id": 987677929, "total_amount": 66.5, "pay_amount": 66.5,
		"total_quantity": 1, "recipient_snapshot": `{}`, "pay_channel": "mock", "pay_status": "paid", "idempotency_key": "fixture:refund-group:987677921",
	}).Error; err != nil {
		t.Fatalf("seed paid group order: %v", err)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_order").Create(map[string]any{
		"id": orderID, "group_order_id": groupID, "order_no": orderNo, "merchant_id": merchantID,
		"merchant_name_snapshot": "七禧退款结算验收茶铺", "store_id": storeID, "store_name_snapshot": "七禧退款结算验收店",
		"user_id": 987677929, "total_amount": 66.5, "pay_amount": 66.5, "total_quantity": 1, "recipient_snapshot": `{}`, "status": "aftersale",
	}).Error; err != nil {
		t.Fatalf("seed Chinese aftersale order: %v", err)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_order_item").Create(map[string]any{
		"id": orderItemID, "order_id": orderID, "product_id": 987677930, "merchant_sku_id": skuID,
		"sku_key": "模拟退款规格", "title_snapshot": "七禧退款结算验收商品", "spec_snapshot": `{"规格":"标准"}`, "unit_price": 66.5, "quantity": 1,
	}).Error; err != nil {
		t.Fatalf("seed order item: %v", err)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_refund").Create(map[string]any{
		"id": refundID, "order_id": orderID, "refund_no": refundNo, "reason": "中文模拟退款验收", "amount": 66.5,
		"refund_type": "money_only", "order_status_before": "completed", "status": "refunding", "idempotency_key": "fixture:refund:987677923",
	}).Error; err != nil {
		t.Fatalf("seed refund: %v", err)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_refund_transaction").Create(map[string]any{
		"id": transactionID, "refund_id": refundID, "channel": "mock", "provider_refund_no": "fixture-provider-refund-987677923",
		"amount": 66.5, "status": "created", "idempotency_key": "fixture:refund-transaction:987677923",
	}).Error; err != nil {
		t.Fatalf("seed refund transaction: %v", err)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_payment_transaction").Create(map[string]any{
		"id": transactionID, "group_order_id": groupID, "channel": "mock", "transaction_no": "fixture-payment-987677921", "amount": 66.5, "status": "succeeded",
	}).Error; err != nil {
		t.Fatalf("seed mock payment: %v", err)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_aftersale_item").Create(map[string]any{"refund_id": refundID, "order_item_id": orderItemID, "quantity": 1, "amount": 66.5}).Error; err != nil {
		t.Fatalf("seed aftersale item: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewCallbackHandler(db, nil, true).Register(router)
	request := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/refund/mock", strings.NewReader(`{"refund_id":987677923}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		return resp
	}
	if resp := request(); resp.Code != http.StatusOK {
		t.Fatalf("first sandbox refund callback status=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp := request(); resp.Code != http.StatusOK {
		t.Fatalf("refund callback replay status=%d body=%s", resp.Code, resp.Body.String())
	}
	var final struct {
		RefundStatus  string `gorm:"column:refund_status"`
		OrderStatus   string `gorm:"column:order_status"`
		GroupStatus   string `gorm:"column:group_status"`
		PaymentStatus string `gorm:"column:payment_status"`
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_refund AS r").Select("r.status AS refund_status,o.status AS order_status,g.pay_status AS group_status,p.status AS payment_status").Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").Joins("JOIN qixi_crm_b_group_order AS g ON g.id = o.group_order_id").Joins("JOIN qixi_crm_b_payment_transaction AS p ON p.group_order_id = g.id AND p.channel = 'mock'").Where("r.id = ?", refundID).Scan(&final).Error; err != nil {
		t.Fatalf("load refund terminal state: %v", err)
	}
	if final.RefundStatus != "refunded" || final.OrderStatus != "cancelled" || final.GroupStatus != "refunded" || final.PaymentStatus != "refunded" {
		t.Fatalf("unexpected refund terminal state: %+v", final)
	}
	var settlement []struct {
		Action         string  `gorm:"column:action"`
		RefundID       uint64  `gorm:"column:refund_id"`
		StoreID        uint64  `gorm:"column:store_id"`
		MerchantID     uint64  `gorm:"column:merchant_id"`
		Amount         float64 `gorm:"column:amount"`
		IdempotencyKey string  `gorm:"column:idempotency_key"`
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_settlement_command_outbox").Where("order_id = ?", orderID).Find(&settlement).Error; err != nil {
		t.Fatalf("load settlement reversal command: %v", err)
	}
	if len(settlement) != 1 || settlement[0].Action != "reverse" || settlement[0].RefundID != refundID || settlement[0].StoreID != storeID || settlement[0].MerchantID != merchantID || settlement[0].Amount != 66.5 || settlement[0].IdempotencyKey != "settlement:reverse:987677923" {
		t.Fatalf("unexpected server-derived settlement reversal: %+v", settlement)
	}
	var stock []struct {
		Action         string `gorm:"column:action"`
		MerchantSKUID  uint64 `gorm:"column:merchant_sku_id"`
		Quantity       int    `gorm:"column:quantity"`
		IdempotencyKey string `gorm:"column:idempotency_key"`
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_stock_command_outbox").Where("order_id = ?", orderID).Find(&stock).Error; err != nil {
		t.Fatalf("load stock restock command: %v", err)
	}
	if len(stock) != 1 || stock[0].Action != "restock" || stock[0].MerchantSKUID != skuID || stock[0].Quantity != 1 || stock[0].IdempotencyKey != "stock:restock:987677922:987677928" {
		t.Fatalf("unexpected server-derived stock restock: %+v", stock)
	}
	var callbacks, events int64
	if err := db.WithContext(ctx).Table("qixi_crm_b_refund_callback").Where("channel = ? AND provider_event_id = ?", "mock", "mock-refund:987677923").Count(&callbacks).Error; err != nil {
		t.Fatalf("count refund callback: %v", err)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_refund_event").Where("refund_id = ?", refundID).Count(&events).Error; err != nil {
		t.Fatalf("count refund event: %v", err)
	}
	if callbacks != 1 || events != 1 {
		t.Fatalf("refund replay created duplicate facts: callbacks=%d events=%d", callbacks, events)
	}
}
