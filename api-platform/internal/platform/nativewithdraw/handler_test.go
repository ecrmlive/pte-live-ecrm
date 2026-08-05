package nativewithdraw

import "testing"

func TestWithdrawalStatusMappings(t *testing.T) {
	cases := []struct {
		filter string
		want   []string
	}{
		{filter: "0", want: []string{"applied", "reviewing"}},
		{filter: "1", want: []string{"approved", "paying", "paid"}},
		{filter: "-1", want: []string{"rejected"}},
	}

	for _, tc := range cases {
		t.Run(tc.filter, func(t *testing.T) {
			got, ok := statuses(tc.filter)
			if !ok || len(got) != len(tc.want) {
				t.Fatalf("filter %q = (%v, %t), want (%v, true)", tc.filter, got, ok, tc.want)
			}
			for index := range tc.want {
				if got[index] != tc.want[index] {
					t.Fatalf("filter %q status[%d] = %q, want %q", tc.filter, index, got[index], tc.want[index])
				}
			}
		})
	}

	if _, ok := statuses("已完成"); ok {
		t.Fatal("unknown status must not be accepted")
	}
}

func TestViewMasksFinancialAccount(t *testing.T) {
	got := view(withdraw{
		ID:           8,
		UserID:       10086,
		WithdrawalNo: "TX-模拟-20260803-001",
		Channel:      "bank",
		Status:       "approved",
		Amount:       88.5,
	})

	if got["financial_account"] != "已脱敏收款账户" {
		t.Fatalf("financial account = %v, want masked value", got["financial_account"])
	}
	if got["status"] != 1 || got["financial_status"] != 0 {
		t.Fatalf("approved withdrawal state = (%v, %v), want (1, 0)", got["status"], got["financial_status"])
	}
	if got["user_id"] != uint64(10086) {
		t.Fatalf("user id = %v, want 10086", got["user_id"])
	}
}

func TestPayoutInputRejectsWeakIdempotencyOrReference(t *testing.T) {
	if !validPayoutInput("withdraw-paid-demo-001", "凭证-20260803-0001") {
		t.Fatal("valid Chinese payout reference rejected")
	}
	if validPayoutInput("短", "凭证-20260803-0001") || validPayoutInput("withdraw-paid-demo-001", "短") {
		t.Fatal("weak payout input accepted")
	}
}
