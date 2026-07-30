package order

import "testing"

func TestMoney(t *testing.T) {
	if got := money(1999); got != 19.99 {
		t.Fatalf("money=%v", got)
	}
}
