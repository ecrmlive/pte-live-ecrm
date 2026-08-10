package recharge

import "testing"

func TestValidPlan(t *testing.T) {
	if !valid(input{Name: "虚构演示充值", Amount: 100, BonusAmount: 8, Status: 1, Sort: 10}, false) {
		t.Fatal("valid recharge plan rejected")
	}
	if !valid(input{Name: "", Amount: 100, BonusAmount: 0, Status: 1, Sort: 0}, false) {
		t.Fatal("nameless recharge plan should be accepted")
	}
	if valid(input{Name: "", Amount: 0, Status: 1}, false) || valid(input{Name: "x", Amount: 100.001, Status: 1}, false) {
		t.Fatal("invalid recharge plan accepted")
	}
	if planName(input{Amount: 10}) != "10元充值" {
		t.Fatalf("unexpected plan name: %s", planName(input{Amount: 10}))
	}
}
