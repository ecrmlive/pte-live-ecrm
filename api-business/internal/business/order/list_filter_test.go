package order

import "testing"

func TestNormalizeGroupPayStatus(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
		bad   bool
	}{
		{name: "全部", input: "", want: ""},
		{name: "全部别名", input: "all", want: ""},
		{name: "待付款", input: "pending", want: "pending"},
		{name: "已支付", input: "paid", want: "paid"},
		{name: "已取消", input: "closed", want: "closed"},
		{name: "拒绝未知状态", input: "completed", bad: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NormalizeGroupPayStatus(tc.input)
			if tc.bad {
				if err != ErrOrderListStatus {
					t.Fatalf("NormalizeGroupPayStatus(%q) error=%v, want %v", tc.input, err, ErrOrderListStatus)
				}
				return
			}
			if err != nil || got != tc.want {
				t.Fatalf("NormalizeGroupPayStatus(%q)=(%q,%v), want (%q,nil)", tc.input, got, err, tc.want)
			}
		})
	}
}

func TestNormalizeOrderListKeyword(t *testing.T) {
	keyword, err := NormalizeOrderListKeyword("  中文商品  ")
	if err != nil || keyword != "中文商品" {
		t.Fatalf("NormalizeOrderListKeyword()=(%q,%v), want (中文商品,nil)", keyword, err)
	}
	if _, err := NormalizeOrderListKeyword(string(make([]rune, 65))); err != ErrOrderListKeyword {
		t.Fatalf("long keyword error=%v, want %v", err, ErrOrderListKeyword)
	}
}

func TestNormalizeGroupFulfillmentStatus(t *testing.T) {
	if got, err := NormalizeGroupFulfillmentStatus("awaiting_receipt"); err != nil || got != "awaiting_receipt" {
		t.Fatalf("valid fulfillment filter=(%q,%v)", got, err)
	}
	if _, err := NormalizeGroupFulfillmentStatus("completed"); err != ErrOrderListFulfillmentStatus {
		t.Fatalf("invalid fulfillment filter error=%v, want %v", err, ErrOrderListFulfillmentStatus)
	}
}
