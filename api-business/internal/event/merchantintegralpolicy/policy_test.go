package merchantintegralpolicy

import "testing"

func TestProjectionPolicyValidation(t *testing.T) {
	validPolicy := payload{StoreID: 1, Enabled: true, PointsPerYuan: 100, MaxDeductionBps: 2000}
	if !valid(validPolicy) {
		t.Fatal("expected valid merchant integral policy")
	}
	for _, value := range []payload{
		{PointsPerYuan: 100, MaxDeductionBps: 2000},
		{StoreID: 1, MaxDeductionBps: 2000},
		{StoreID: 1, PointsPerYuan: 100},
		{StoreID: 1, PointsPerYuan: 100, MaxDeductionBps: 10001},
	} {
		if valid(value) {
			t.Fatalf("unexpected valid policy: %+v", value)
		}
	}
}
