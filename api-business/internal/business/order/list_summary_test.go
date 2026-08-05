package order

import "testing"

func TestGroupFulfillmentStatus(t *testing.T) {
	cases := []struct {
		name   string
		group  groupRow
		orders []orderRow
		want   string
	}{
		{name: "待付款", group: groupRow{PayStatus: "pending"}, want: "pending"},
		{name: "已取消", group: groupRow{PayStatus: "closed"}, want: "closed"},
		{name: "待发货", group: groupRow{PayStatus: "paid"}, orders: []orderRow{{Status: "paid"}}, want: "fulfilling"},
		{name: "待收货", group: groupRow{PayStatus: "paid"}, orders: []orderRow{{Status: "shipped"}}, want: "shipped"},
		{name: "售后中", group: groupRow{PayStatus: "paid"}, orders: []orderRow{{Status: "aftersale"}}, want: "aftersale"},
		{name: "已完成", group: groupRow{PayStatus: "paid"}, orders: []orderRow{{Status: "completed"}}, want: "completed"},
		{name: "混合订单优先提示待发货", group: groupRow{PayStatus: "paid"}, orders: []orderRow{{Status: "shipped"}, {Status: "paid"}}, want: "fulfilling"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := GroupFulfillmentStatus(tc.group, tc.orders); got != tc.want {
				t.Fatalf("GroupFulfillmentStatus()=%q, want %q", got, tc.want)
			}
		})
	}
}
