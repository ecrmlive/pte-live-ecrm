package order

import "testing"

func TestCanArchiveGroup(t *testing.T) {
	cases := []struct {
		name   string
		group  groupRow
		orders []orderRow
		want   bool
	}{
		{name: "未支付取消订单", group: groupRow{PayStatus: "closed"}, want: true},
		{name: "已完成订单", group: groupRow{PayStatus: "paid"}, orders: []orderRow{{Status: "completed"}}, want: true},
		{name: "多个订单均已完成或取消", group: groupRow{PayStatus: "paid"}, orders: []orderRow{{Status: "completed"}, {Status: "cancelled"}}, want: true},
		{name: "履约中的订单", group: groupRow{PayStatus: "paid"}, orders: []orderRow{{Status: "shipped"}}, want: false},
		{name: "待支付订单", group: groupRow{PayStatus: "pending"}, want: false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := CanArchiveGroup(testCase.group, testCase.orders); got != testCase.want {
				t.Fatalf("CanArchiveGroup()=%v, want %v", got, testCase.want)
			}
		})
	}
}
