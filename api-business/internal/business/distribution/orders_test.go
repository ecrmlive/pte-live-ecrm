package distribution

import "testing"

func TestSpreadOrderKeyword(t *testing.T) {
	cases := []struct {
		input    string
		orderID  uint64
		filtered bool
		valid    bool
	}{
		{"", 0, false, true},
		{" 20301 ", 20301, true, true},
		{"订单20301", 0, false, false},
		{"-20301", 0, false, false},
	}
	for _, item := range cases {
		got, filtered, err := spreadOrderKeyword(item.input)
		if (err == nil) != item.valid || got != item.orderID || filtered != item.filtered {
			t.Fatalf("spreadOrderKeyword(%q)=(%d,%t,%v)", item.input, got, filtered, err)
		}
	}
}

func TestCommissionOrderStatus(t *testing.T) {
	if commissionOrderStatus("available") != "可提现" || commissionOrderStatus("voided") != "佣金记录" {
		t.Fatal("unexpected commission status label")
	}
}
