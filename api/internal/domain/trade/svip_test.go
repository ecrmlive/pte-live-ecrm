package trade

import (
	"testing"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/cart"
)

func TestUserSvipActiveAt(t *testing.T) {
	now := time.Date(2026, 7, 22, 12, 0, 0, 0, time.Local)
	endOk := now.Add(24 * time.Hour)
	endPast := now.Add(-time.Hour)

	cases := []struct {
		name   string
		isSvip int8
		end    *time.Time
		want   bool
	}{
		{"off", 0, nil, false},
		{"trial", 1, nil, true},
		{"forever", 3, nil, true},
		{"dated_ok", 2, &endOk, true},
		{"dated_past", 2, &endPast, false},
		{"dated_nil", 2, nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := UserSvipActiveAt(tc.isSvip, tc.end, now); got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestResolveSvipLinePrice(t *testing.T) {
	if got := resolveSvipLinePrice(cart.Cart{SvipPriceType: 1, Price: 100}); got != 90 {
		t.Fatalf("ratio price got %v", got)
	}
	if got := resolveSvipLinePrice(cart.Cart{SvipPriceType: 2, SvipPrice: 19.9, Price: 100}); got != 19.9 {
		t.Fatalf("custom price got %v", got)
	}
	if got := resolveSvipLinePrice(cart.Cart{SvipPriceType: 0, SvipPrice: 10, Price: 100}); got != 0 {
		t.Fatalf("off got %v", got)
	}
}
