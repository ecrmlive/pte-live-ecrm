package funding

import "testing"

func TestSVIPPlanValidation(t *testing.T) {
	days := uint(30)
	if !validSVIPPlan(svipPlan{Price: 29, PlanType: "period", DurationDays: &days}) {
		t.Fatal("paid periodic plan should be accepted")
	}
	if validSVIPPlan(svipPlan{Price: 0, PlanType: "trial", DurationDays: &days}) {
		t.Fatal("zero-price plan must not enter an H5 payment flow")
	}
	if validSVIPPlan(svipPlan{Price: 29, PlanType: "period"}) {
		t.Fatal("periodic plan without duration must be rejected")
	}
	if !validSVIPPlan(svipPlan{Price: 299, PlanType: "lifetime"}) {
		t.Fatal("paid lifetime plan should be accepted")
	}
}

func TestFundingInputGuards(t *testing.T) {
	if !validKey("recharge-20260804-123456") {
		t.Fatal("generated idempotency key should be accepted")
	}
	if validKey("短") || validKey("contains space 123") {
		t.Fatal("unsafe idempotency keys must be rejected")
	}
	if roundedAmount(12.345) != 12.35 || cents(12.345) != 1235 {
		t.Fatal("money conversion must round to cents consistently")
	}
}
