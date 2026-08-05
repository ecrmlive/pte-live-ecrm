package nativesettlement

import (
	"errors"
	"reflect"
	"testing"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/adminscope"
	merchantsettlement "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/merchantsettlement"
)

func TestSettlementStatusAllowsOnlyProjectionStates(t *testing.T) {
	for _, raw := range []string{"", "bill_pending", "bill_frozen", "withdraw_applied", "approved", "paid", "rejected", "cancelled"} {
		if _, ok := settlementStatus(raw); !ok {
			t.Fatalf("status %q should be allowed", raw)
		}
	}
	if _, ok := settlementStatus("已打款"); ok {
		t.Fatal("arbitrary Chinese status must be rejected")
	}
}

func TestTransferStatusAllowsOnlyPayoutPipeline(t *testing.T) {
	for _, raw := range []string{"", "approved", "paid", "rejected"} {
		if _, ok := transferStatus(raw); !ok {
			t.Fatalf("transfer status %q should be allowed", raw)
		}
	}
	for _, raw := range []string{"bill_pending", "withdraw_applied", "已打款"} {
		if _, ok := transferStatus(raw); ok {
			t.Fatalf("transfer status %q must be rejected", raw)
		}
	}
}

func TestSettlementRegionScopeNeverUsesMerchantIDs(t *testing.T) {
	cases := []struct {
		name  string
		scope adminscope.MerchantScope
		want  []uint64
		err   error
	}{
		{name: "platform has full projection", scope: adminscope.MerchantScope{Full: true}, want: nil},
		{name: "region only sees configured regions", scope: adminscope.MerchantScope{RegionIDs: []uint64{10, 20}}, want: []uint64{10, 20}},
		{name: "merchant ids cannot grant settlement read", scope: adminscope.MerchantScope{MerchantIDs: []uint64{1}}, err: adminscope.ErrNotConfigured},
		{name: "missing region scope fails closed", scope: adminscope.MerchantScope{}, err: adminscope.ErrNotConfigured},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := settlementRegionScope(tc.scope)
			if !errors.Is(err, tc.err) || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("settlement scope = %#v, %v; want %#v, %v", got, err, tc.want, tc.err)
			}
		})
	}
}

func TestSettlementCommandValidation(t *testing.T) {
	base := merchantsettlement.Command{SettlementID: 9001, OperatorID: 1, IdempotencyKey: "settlement-review-9001"}
	base.Action = "approve"
	if !validSettlementCommand(base) {
		t.Fatal("approve command should accept a valid idempotency key")
	}
	base.Action, base.ReviewNote = "reject", "结算周期数据不完整，请补充后再申请"
	if !validSettlementCommand(base) {
		t.Fatal("reject command should require and accept Chinese review note")
	}
	base.Action, base.PayoutReference = "mark_paid", "凭证-20260803-0001"
	if !validSettlementCommand(base) {
		t.Fatal("mark-paid command should accept non-sensitive payout reference")
	}
	base.IdempotencyKey = "短"
	if validSettlementCommand(base) {
		t.Fatal("weak command idempotency key must be rejected")
	}
}
