package merchantsettlement

import (
	"context"
	"testing"
)

func TestSettlementCommandValidationAndFailClosedDatabase(t *testing.T) {
	valid := command{SettlementID: 9001, Action: "reject", OperatorID: 101, IdempotencyKey: "settlement-reject-9001", ReviewNote: "结算周期数据不完整，请补充后再申请"}
	if !validCommand(valid) {
		t.Fatal("valid Chinese reject command rejected")
	}
	valid.Action, valid.ReviewNote, valid.PayoutReference = "mark_paid", "", "凭证-20260803-0001"
	if !validCommand(valid) {
		t.Fatal("valid non-sensitive payout reference rejected")
	}
	valid.IdempotencyKey = "短"
	if validCommand(valid) {
		t.Fatal("weak idempotency key accepted")
	}
	result, err := ApplyCommand(context.Background(), nil, []byte(`{"settlement_id":9001,"action":"approve","operator_id":101,"idempotency_key":"settlement-approve-9001"}`))
	if err == nil || result.Code != "failed" {
		t.Fatalf("unavailable merchant database must fail closed, got %#v, %v", result, err)
	}
}

func TestSettlementCommandReplayRequiresSameActionPayloadAndOperator(t *testing.T) {
	reviewKey, payoutKey, payoutReference := "settlement-review-9001", "settlement-paid-9001", "本地模拟凭证-001"
	approved := settlementBill{Status: "approved", ReviewedByAdminID: 101, ReviewIdempotencyKey: &reviewKey, ReviewNote: "审核通过", PayoutIdempotencyKey: &payoutKey, PayoutReference: &payoutReference}
	if !commandAlreadyApplied(approved, command{Action: "approve", OperatorID: 101, IdempotencyKey: reviewKey, ReviewNote: "审核通过"}) {
		t.Fatal("identical approval replay must be accepted")
	}
	if commandAlreadyApplied(approved, command{Action: "reject", OperatorID: 101, IdempotencyKey: reviewKey, ReviewNote: "审核通过"}) {
		t.Fatal("approval key must not replay as reject")
	}
	if commandAlreadyApplied(approved, command{Action: "approve", OperatorID: 102, IdempotencyKey: reviewKey, ReviewNote: "审核通过"}) {
		t.Fatal("approval key must not replay for another operator")
	}
	paid := settlementBill{Status: "paid", PayoutIdempotencyKey: &payoutKey, PayoutReference: &payoutReference}
	if !commandAlreadyApplied(paid, command{Action: "mark_paid", IdempotencyKey: payoutKey, PayoutReference: payoutReference}) {
		t.Fatal("identical payout replay must be accepted")
	}
	if commandAlreadyApplied(paid, command{Action: "mark_paid", IdempotencyKey: payoutKey, PayoutReference: "本地模拟篡改凭证-002"}) {
		t.Fatal("payout key must not replay with changed reference")
	}
}
