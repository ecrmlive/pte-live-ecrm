package nativedistribution

import "testing"

func TestDistributionFiltersFailClosed(t *testing.T) {
	for _, value := range []string{"", "0", "1"} {
		if _, ok := promoterStatus(value); !ok {
			t.Fatalf("promoter status %q should be allowed", value)
		}
	}
	if _, ok := promoterStatus("启用"); ok {
		t.Fatal("arbitrary promoter status must be rejected")
	}
	for _, value := range []string{"", "pending", "available", "settled", "voided"} {
		if _, ok := commissionStatus(value); !ok {
			t.Fatalf("commission status %q should be allowed", value)
		}
	}
	if _, ok := commissionStatus("已结算"); ok {
		t.Fatal("arbitrary commission status must be rejected")
	}
}
