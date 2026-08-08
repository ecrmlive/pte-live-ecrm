package nativeorder

import "testing"

func TestPayChannelFromType(t *testing.T) {
	if payChannelFromType("1") != "wechat" {
		t.Fatal("wechat")
	}
	if payChannelFromType("balance") != "balance" {
		t.Fatal("balance")
	}
	if payChannelFromType("x") != "" {
		t.Fatal("unknown")
	}
}

func TestMaskPhone(t *testing.T) {
	if got := maskPhone("13012348558"); got != "130****8558" {
		t.Fatalf("got %s", got)
	}
	if got := maskPhone(""); got != "--" {
		t.Fatalf("empty %s", got)
	}
}

func TestOrderStatusLabel(t *testing.T) {
	label, code := orderStatusLabel(order{Status: "pending_pay"})
	if label != "待付款" || code != -2 {
		t.Fatalf("%s %d", label, code)
	}
	label, code = orderStatusLabel(order{Status: "completed", PendingComment: 1})
	if label != "待评价" || code != 2 {
		t.Fatalf("%s %d", label, code)
	}
	label, code = orderStatusLabel(order{HasRefunded: 1, Status: "completed"})
	if label != "已退款" || code != -1 {
		t.Fatalf("%s %d", label, code)
	}
}

func TestActivityAndDeliveryLabels(t *testing.T) {
	if activityTypeLabel(0) != "普通" {
		t.Fatal("activity")
	}
	if deliveryTypeLabel("express") != "快递发货" {
		t.Fatal("delivery")
	}
}

func TestVerifierDisplayName(t *testing.T) {
	if got := verifierDisplayName(0, accountNameRow{}); got != "管理员核销" {
		t.Fatalf("admin got %s", got)
	}
	if got := verifierDisplayName(9, accountNameRow{DisplayName: "店员甲"}); got != "店员甲" {
		t.Fatalf("display got %s", got)
	}
	if got := verifierDisplayName(9527, accountNameRow{}); got != "9527" {
		t.Fatalf("id got %s", got)
	}
}

func TestVerifyPayTypeLabel(t *testing.T) {
	if got := verifyPayTypeLabel("balance", 1); got != "余额支付" {
		t.Fatalf("got %s", got)
	}
	if got := verifyPayTypeLabel("wechat", 1); got != "微信支付" {
		t.Fatalf("got %s", got)
	}
}

func TestVerifyRecordRowAsOrder(t *testing.T) {
	row := verifyRecordRow{
		ID: 961004, OrderNo: "PSTAT-O-0004", MerchantID: 1, StoreID: 1,
		PayAmount: 180, Status: "completed", PayChannel: "balance",
		VerifiedByAccountID: 2, VerifyStatus: "used",
	}
	base := row.asOrder()
	if base.ID != 961004 || base.OrderNo != "PSTAT-O-0004" || base.MerchantID != 1 || base.PayAmount != 180 {
		t.Fatalf("asOrder lost fields: %+v", base)
	}
}
