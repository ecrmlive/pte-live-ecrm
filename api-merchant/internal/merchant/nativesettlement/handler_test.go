package nativesettlement

import "testing"

func TestSettlementInputRejectsWeakIdempotencyAndUnknownStatus(t *testing.T) {
	for _, value := range []string{"settlement-demo-001", "七禧结算幂等键-0001"} {
		if !validIdempotencyKey(value) {
			t.Fatalf("idempotency key %q rejected", value)
		}
	}
	if validIdempotencyKey("短") || validIdempotencyKey("") {
		t.Fatal("weak idempotency keys must be rejected")
	}
	if _, ok := settlementStatus("未知状态"); ok {
		t.Fatal("unknown settlement status must be rejected")
	}
}
