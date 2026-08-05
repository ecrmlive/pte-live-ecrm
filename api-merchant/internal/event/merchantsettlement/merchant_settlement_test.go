package merchantsettlement

import (
	"testing"
	"time"
)

func TestSettlementEventPayloadRejectsSensitiveOrInvalidStateByShape(t *testing.T) {
	now := time.Date(2026, 8, 3, 10, 0, 0, 0, time.UTC)
	payload := Payload{SettlementID: 9001, MerchantID: 1, StoreID: 1, MerchantName: "七禧演示茶铺", PeriodStart: now.AddDate(0, -1, 0), PeriodEnd: now, Amount: 1280.50, Status: "withdraw_applied", UpdatedAt: now}
	if !ValidPayload(payload) {
		t.Fatal("valid Chinese settlement payload rejected")
	}
	payload.Amount = -12.5
	payload.Status = "bill_pending"
	if !ValidPayload(payload) {
		t.Fatal("negative current-period refund adjustment rejected")
	}
	payload.Amount = 1280.50
	payload.Status = "已打款"
	if ValidPayload(payload) {
		t.Fatal("unrecognised status must be rejected")
	}
}
