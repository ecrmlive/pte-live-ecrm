package merchantledger

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestSettlementLedgerIntegrationKeepsFrozenBillImmutable verifies the
// merchant-owned NATS consumer writes immutable entries, rejects a changed
// replay, and posts a refund adjustment to the current open bill instead of
// mutating a historical frozen bill.
func TestSettlementLedgerIntegrationKeepsFrozenBillImmutable(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_SETTLEMENT_MERCHANT_TEST_DSN"))
	natsURL := strings.TrimSpace(os.Getenv("ECRM_SETTLEMENT_NATS_URL"))
	if dsn == "" || natsURL == "" {
		t.Skip("set ECRM_SETTLEMENT_MERCHANT_TEST_DSN and ECRM_SETTLEMENT_NATS_URL to run settlement ledger integration test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated merchant database: %v", err)
	}
	const merchantID uint64 = 987677904
	const storeID uint64 = 987677903
	const accrueOrderID uint64 = 987677901
	const reverseOrderID uint64 = 987677905
	const refundID uint64 = 987677906
	ctx := context.Background()
	t.Cleanup(func() {
		var billIDs []uint64
		_ = db.WithContext(ctx).Table("qixi_crm_m_settlement_bill").Where("merchant_id = ?", merchantID).Pluck("id", &billIDs).Error
		if len(billIDs) > 0 {
			_ = db.WithContext(ctx).Table("qixi_crm_m_outbox").Where("aggregate_type = ? AND aggregate_id IN ?", "settlement_bill", billIDs).Delete(nil).Error
		}
		_ = db.WithContext(ctx).Table("qixi_crm_m_settlement_entry").Where("merchant_id = ?", merchantID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_m_settlement_bill").Where("merchant_id = ?", merchantID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_m_store").Where("id = ?", storeID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_m_merchant").Where("id = ?", merchantID).Delete(nil).Error
	})
	if err := db.WithContext(ctx).Table("qixi_crm_m_merchant").Create(map[string]any{"id": merchantID, "name": "七禧结算账本验收茶铺", "region_id": 987677900, "status": 1}).Error; err != nil {
		t.Fatalf("seed Chinese merchant: %v", err)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_m_store").Create(map[string]any{"id": storeID, "merchant_id": merchantID, "app_id": "qixi.settlement.facts.987677903", "name": "七禧结算账本验收店", "status": 1}).Error; err != nil {
		t.Fatalf("seed Chinese store: %v", err)
	}
	subscriber, err := StartCommandSubscriber(ctx, db, natsURL)
	if err != nil || subscriber == nil {
		t.Fatalf("start settlement ledger subscriber: %v", err)
	}
	t.Cleanup(subscriber.Close)
	publisher, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_settlement_ledger_acceptance"))
	if err != nil {
		t.Fatalf("connect isolated NATS: %v", err)
	}
	t.Cleanup(publisher.Close)
	request := func(in command) commandResult {
		wire, err := json.Marshal(in)
		if err != nil {
			t.Fatalf("marshal settlement command: %v", err)
		}
		msg, err := publisher.Request(CommandSubject, wire, 2*time.Second)
		if err != nil {
			t.Fatalf("NATS settlement request: %v", err)
		}
		var out commandResult
		if err := json.Unmarshal(msg.Data, &out); err != nil {
			t.Fatalf("decode settlement reply: %v", err)
		}
		return out
	}
	accrual := command{Action: "accrue", OrderID: accrueOrderID, StoreID: storeID, MerchantID: merchantID, Amount: 128.5, IdempotencyKey: "settlement:accrue:987677901"}
	if out := request(accrual); out.Status != "accepted" || out.Code != "" {
		t.Fatalf("accrual result=%+v, want accepted", out)
	}
	if out := request(accrual); out.Status != "accepted" || out.Code != "" {
		t.Fatalf("exact accrual replay result=%+v, want accepted", out)
	}
	accrual.Amount = 128.6
	if out := request(accrual); out.Code != "conflict" {
		t.Fatalf("changed accrual replay result=%+v, want conflict", out)
	}
	// Age the generated pending bill, then exercise the real freezer. This
	// makes the refund assertion independent from the wall-clock month while
	// preserving the production rule that only elapsed bills can freeze.
	historicalStart := time.Now().UTC().AddDate(0, -1, 0)
	historicalStart = time.Date(historicalStart.Year(), historicalStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	historicalEnd := historicalStart.AddDate(0, 1, 0).Add(-time.Second)
	if err := db.WithContext(ctx).Table("qixi_crm_m_settlement_bill").Where("merchant_id = ? AND status = ?", merchantID, "bill_pending").Updates(map[string]any{"period_start": historicalStart, "period_end": historicalEnd}).Error; err != nil {
		t.Fatalf("age pending bill for freezer acceptance: %v", err)
	}
	freezeExpired(ctx, db, time.Now().UTC())
	var frozen struct {
		Amount float64 `gorm:"column:amount"`
		Status string  `gorm:"column:status"`
	}
	if err := db.WithContext(ctx).Table("qixi_crm_m_settlement_bill").Select("amount,status").Where("merchant_id = ? AND status = ?", merchantID, "bill_frozen").Scan(&frozen).Error; err != nil || frozen.Status != "bill_frozen" || frozen.Amount != 128.5 {
		t.Fatalf("frozen bill=%+v err=%v, want accrued historical bill", frozen, err)
	}
	reversal := command{Action: "reverse", OrderID: reverseOrderID, RefundID: refundID, StoreID: storeID, MerchantID: merchantID, Amount: 128.5, IdempotencyKey: "settlement:reverse:987677906"}
	if out := request(reversal); out.Status != "accepted" || out.Code != "" {
		t.Fatalf("reversal result=%+v, want accepted", out)
	}
	var entries []struct {
		EntryType string  `gorm:"column:entry_type"`
		Amount    float64 `gorm:"column:amount"`
	}
	if err := db.WithContext(ctx).Table("qixi_crm_m_settlement_entry").Select("entry_type,amount").Where("merchant_id = ?", merchantID).Order("id ASC").Find(&entries).Error; err != nil {
		t.Fatalf("load immutable settlement entries: %v", err)
	}
	if len(entries) != 2 || entries[0].EntryType != "order_accrual" || entries[0].Amount != 128.5 || entries[1].EntryType != "refund_reversal" || entries[1].Amount != -128.5 {
		t.Fatalf("unexpected immutable entries: %+v", entries)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_m_settlement_bill").Select("amount,status").Where("merchant_id = ? AND status = ?", merchantID, "bill_frozen").Scan(&frozen).Error; err != nil || frozen.Status != "bill_frozen" || frozen.Amount != 128.5 {
		t.Fatalf("historical frozen bill=%+v err=%v, want unchanged", frozen, err)
	}
	var current struct {
		Amount float64 `gorm:"column:amount"`
		Status string  `gorm:"column:status"`
	}
	if err := db.WithContext(ctx).Table("qixi_crm_m_settlement_bill").Select("amount,status").Where("merchant_id = ? AND status = ?", merchantID, "bill_pending").Order("id ASC").Limit(1).Scan(&current).Error; err != nil || current.Status != "bill_pending" || current.Amount != -128.5 {
		t.Fatalf("current adjustment bill=%+v err=%v, want negative refund adjustment", current, err)
	}
}
