package integralorder

import (
	"testing"
	"time"
)

func TestStatusLabelOf(t *testing.T) {
	label, code := statusLabelOf(orderRow{Status: "paid"})
	if label != "待发货" || code != 0 {
		t.Fatalf("paid => %s/%d", label, code)
	}
	label, code = statusLabelOf(orderRow{Status: "shipped"})
	if label != "待收货" || code != 1 {
		t.Fatalf("shipped => %s/%d", label, code)
	}
	archived := orderRow{Status: "fulfilling"}
	now := archived.CreatedAt
	archived.UserArchivedAt = &now
	label, code = statusLabelOf(archived)
	if label != "已删除" || code != -10 {
		t.Fatalf("user deleted => %s/%d", label, code)
	}
}

func TestCanDeliver(t *testing.T) {
	if !canDeliver(orderRow{Status: "paid"}) {
		t.Fatal("paid should deliver")
	}
	if canDeliver(orderRow{Status: "shipped"}) {
		t.Fatal("shipped should not deliver")
	}
	if canDeliver(orderRow{Status: "paid", IsSystemDel: 1}) {
		t.Fatal("system deleted should not deliver")
	}
	now := time.Now()
	if canDeliver(orderRow{Status: "paid", UserArchivedAt: &now}) {
		t.Fatal("user deleted should not deliver")
	}
}

func TestSpecText(t *testing.T) {
	got := specText([]byte(`{"规格":"礼盒装"}`))
	if got == "" || !contains(got, "礼盒装") {
		t.Fatalf("specText=%q", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(len(s) > 0 && (func() bool {
			for i := 0; i+len(sub) <= len(s); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})()))
}
