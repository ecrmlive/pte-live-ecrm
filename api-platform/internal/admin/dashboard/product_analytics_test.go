package dashboard

import (
	"testing"
	"time"
)

func TestGrowthRate(t *testing.T) {
	cases := []struct {
		cur, prev int64
		want      float64
	}{
		{0, 0, 0},
		{10, 0, 100},
		{110, 100, 10},
		{2, 3, -33.33},
		{0, 5, -100},
	}
	for _, c := range cases {
		got := growthRate(c.cur, c.prev)
		if got != c.want {
			t.Fatalf("growthRate(%d,%d)=%v want %v", c.cur, c.prev, got, c.want)
		}
	}
}

func TestResolveAnalyticsWindowLately7(t *testing.T) {
	now := time.Date(2026, 8, 6, 15, 0, 0, 0, shanghaiLoc())
	win, err := resolveAnalyticsWindow("lately7", now)
	if err != nil {
		t.Fatal(err)
	}
	if len(win.Buckets) != 7 {
		t.Fatalf("buckets=%d want 7: %v", len(win.Buckets), win.Buckets)
	}
	if win.Buckets[0] != "07-31" || win.Buckets[6] != "08-06" {
		t.Fatalf("buckets=%v", win.Buckets)
	}
	if win.SQLFormat != "%m-%d" {
		t.Fatalf("format=%s", win.SQLFormat)
	}
}

func TestResolveAnalyticsWindowCustom(t *testing.T) {
	now := time.Date(2026, 8, 6, 15, 0, 0, 0, shanghaiLoc())
	win, err := resolveAnalyticsWindow("2026/07/31-2026/08/06", now)
	if err != nil {
		t.Fatal(err)
	}
	if len(win.Buckets) != 7 {
		t.Fatalf("buckets=%v", win.Buckets)
	}
}
