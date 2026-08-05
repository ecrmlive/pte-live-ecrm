package distribution

import "testing"

func TestValidWithdrawalInput(t *testing.T) {
	goodWechat := withdrawalInput{Amount: 0.01, Channel: "wechat", AccountName: "张三", IdempotencyKey: "withdrawal-123"}
	if !validWithdrawalInput(goodWechat) {
		t.Fatal("wechat withdrawal should be valid")
	}
	goodBank := withdrawalInput{Amount: 18.88, Channel: "bank", AccountName: "李四", AccountNo: "622202123456", IdempotencyKey: "withdrawal-456"}
	if !validWithdrawalInput(goodBank) {
		t.Fatal("bank withdrawal should be valid")
	}
	for _, input := range []withdrawalInput{
		{Amount: 0, Channel: "wechat", AccountName: "张三", IdempotencyKey: "withdrawal-123"},
		{Amount: 1, Channel: "bank", AccountName: "李四", AccountNo: "123", IdempotencyKey: "withdrawal-456"},
		{Amount: 1, Channel: "wechat", AccountName: "王", IdempotencyKey: "withdrawal-789"},
	} {
		if validWithdrawalInput(input) {
			t.Fatalf("input should be invalid: %#v", input)
		}
	}
}

func TestWithdrawalStatusText(t *testing.T) {
	if got := withdrawalStatusText("rejected"); got != "已驳回" {
		t.Fatalf("status text = %q", got)
	}
	if got := withdrawalChannelText("bank"); got != "银行卡" {
		t.Fatalf("channel text = %q", got)
	}
	if got := roundMoney(1.235); got != 1.24 {
		t.Fatalf("round money = %v", got)
	}
}
