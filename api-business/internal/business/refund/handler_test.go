package refund

import "testing"

func TestUserRefundStateMapping(t *testing.T) {
	cases := map[string]int{
		"applied": 0, "merchant_handling": 0, "platform_intervene": 4,
		"refunding": 0, "refunded": 3, "rejected": -1, "cancelled": -2,
	}
	for state, want := range cases {
		if got := legacyStatus(state); got != want {
			t.Fatalf("legacyStatus(%q) = %d, want %d", state, got, want)
		}
	}
}

func TestRefundApplyOrderStates(t *testing.T) {
	for _, state := range []string{"paid", "fulfilling", "shipped"} {
		if !canApply(state) {
			t.Fatalf("expected %q to allow after-sale apply", state)
		}
	}
	for _, state := range []string{"pending_pay", "completed", "cancelled", "aftersale"} {
		if canApply(state) {
			t.Fatalf("expected %q to reject after-sale apply", state)
		}
	}
}
