package recharge

import "testing"

func TestValidPlan(t *testing.T) {
	if !valid(input{Name: "虚构演示充值", Amount: 100, BonusAmount: 8, Status: 1, Sort: 10}) {
		t.Fatal("valid recharge plan rejected")
	}
	if valid(input{Name: "", Amount: 100, Status: 1}) || valid(input{Name: "x", Amount: 0, Status: 1}) {
		t.Fatal("invalid recharge plan accepted")
	}
}
