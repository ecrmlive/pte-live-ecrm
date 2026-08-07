package merchantdeposit

import "testing"

func TestDepositStateTransitionsFailClosed(t *testing.T) {
	for _, state := range []string{"funded", "shortfall"} {
		if !depositCanDeduct(state) || !depositCanApproveRefund(state) {
			t.Fatalf("state %q must allow the configured deposit transition", state)
		}
	}
	for _, state := range []string{"pending", "shortfall"} {
		if !depositCanFundOffline(state) {
			t.Fatalf("state %q must allow offline fund", state)
		}
	}
	for _, state := range []string{"not_required", "pending", "refund_pending", "refunded", "unknown"} {
		if depositCanDeduct(state) || depositCanApproveRefund(state) {
			t.Fatalf("state %q must fail closed", state)
		}
	}
	for _, state := range []string{"not_required", "funded", "refund_pending", "refunded", "unknown"} {
		if depositCanFundOffline(state) {
			t.Fatalf("state %q must not allow offline fund", state)
		}
	}
	if got := depositStateAfterBalance(500, 500); got != "funded" {
		t.Fatalf("full balance state = %q", got)
	}
	if got := depositStateAfterBalance(500, 499.99); got != "shortfall" {
		t.Fatalf("short balance state = %q", got)
	}
	if got := depositStateAfterRefundBalance(500, 0); got != "refunded" {
		t.Fatalf("zero balance refund state = %q", got)
	}
}

func TestPayoutRegistrationRequiresStableNonSensitiveReference(t *testing.T) {
	if !validPayoutRegistrationInput("deposit-paid-7301", "本地模拟凭证-7301") {
		t.Fatal("valid Chinese fixture input must pass")
	}
	for _, input := range [][2]string{{"短", "凭证"}, {"deposit-paid-7301", ""}} {
		if validPayoutRegistrationInput(input[0], input[1]) {
			t.Fatalf("invalid payout input must be rejected: %#v", input)
		}
	}
}

func TestDepositAmountUsesExactCentBoundary(t *testing.T) {
	for _, amount := range []float64{0.01, 100, 9999999999.99} {
		if !validDepositAmount(amount) {
			t.Fatalf("valid amount rejected: %v", amount)
		}
	}
	for _, amount := range []float64{0, -0.01, 0.001, 100.123, 10000000000} {
		if validDepositAmount(amount) {
			t.Fatalf("invalid amount accepted: %v", amount)
		}
	}
}
