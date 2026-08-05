package merchantledger

import "testing"

func TestSettlementLedgerCommandValidation(t *testing.T) {
	if !valid(command{Action: "accrue", OrderID: 1, StoreID: 2, MerchantID: 3, Amount: 88.5, IdempotencyKey: "settlement:accrue:1"}) {
		t.Fatal("expected accrual command to be valid")
	}
	if !valid(command{Action: "reverse", OrderID: 1, RefundID: 9, StoreID: 2, MerchantID: 3, Amount: 88.5, IdempotencyKey: "settlement:reverse:9"}) {
		t.Fatal("expected reversal command to be valid")
	}
	if valid(command{Action: "reverse", OrderID: 1, StoreID: 2, MerchantID: 3, Amount: 88.5, IdempotencyKey: "settlement:reverse:0"}) {
		t.Fatal("unanchored reversal must be rejected")
	}
	if valid(command{Action: "accrue", OrderID: 1, StoreID: 2, MerchantID: 3, Amount: 88.5, IdempotencyKey: "settlement:accrue:2"}) {
		t.Fatal("accrual idempotency key must be anchored to its order")
	}
	if valid(command{Action: "reverse", OrderID: 1, RefundID: 9, StoreID: 2, MerchantID: 3, Amount: 88.5, IdempotencyKey: "settlement:reverse:10"}) {
		t.Fatal("reversal idempotency key must be anchored to its refund")
	}
}
