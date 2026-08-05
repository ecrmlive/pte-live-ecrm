package merchantledger

import "testing"

func TestSettlementCommandRequiresFactAnchoredIdempotencyKey(t *testing.T) {
	if !valid(Command{Action: "accrue", OrderID: 71, StoreID: 8, MerchantID: 9, Amount: 128.5, IdempotencyKey: "settlement:accrue:71"}) {
		t.Fatal("expected fact-anchored accrual command to be valid")
	}
	if !valid(Command{Action: "reverse", OrderID: 71, RefundID: 17, StoreID: 8, MerchantID: 9, Amount: 128.5, IdempotencyKey: "settlement:reverse:17"}) {
		t.Fatal("expected fact-anchored reversal command to be valid")
	}
	if valid(Command{Action: "accrue", OrderID: 71, StoreID: 8, MerchantID: 9, Amount: 128.5, IdempotencyKey: "settlement:accrue:72"}) {
		t.Fatal("accrual command accepted a mismatched idempotency key")
	}
	if valid(Command{Action: "reverse", OrderID: 71, RefundID: 17, StoreID: 8, MerchantID: 9, Amount: 128.5, IdempotencyKey: "settlement:reverse:18"}) {
		t.Fatal("reversal command accepted a mismatched idempotency key")
	}
}
