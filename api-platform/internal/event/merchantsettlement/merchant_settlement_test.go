package merchantsettlement

import (
	"testing"
	"time"
)

func TestSettlementProjectionRejectsUnknownState(t *testing.T) {
	now := time.Date(2026, 8, 3, 10, 0, 0, 0, time.UTC)
	payload := Payload{SettlementID: 9001, MerchantID: 1, StoreID: 1, MerchantName: "七禧演示茶铺", PeriodStart: now.AddDate(0, -1, 0), PeriodEnd: now, Amount: 1280.50, Status: "withdraw_applied", UpdatedAt: now}
	if !validPayload(payload) {
		t.Fatal("valid settlement payload rejected")
	}
	payload.Amount = -12.5
	payload.Status = "bill_pending"
	if !validPayload(payload) {
		t.Fatal("negative current-period refund adjustment rejected")
	}
	payload.Amount = 1280.50
	payload.Status = "approved_by_anyone"
	if validPayload(payload) {
		t.Fatal("unknown settlement state accepted")
	}
}
