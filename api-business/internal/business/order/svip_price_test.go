package order

import "testing"

func TestResolveSVIPUnitCents(t *testing.T) {
	tests := []struct {
		name       string
		list       int64
		product    int8
		active     bool
		priceType  int8
		configured float64
		want       int64
	}{
		{name: "nine discount", list: 36900, active: true, priceType: 1, want: 33210},
		{name: "fixed member price", list: 46900, active: true, priceType: 2, configured: 429, want: 42900},
		{name: "inactive", list: 36900, active: false, priceType: 1, want: 36900},
		{name: "activity never stacks", list: 36900, product: 1, active: true, priceType: 1, want: 36900},
		{name: "higher fixed price ignored", list: 36900, active: true, priceType: 2, configured: 399, want: 36900},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolveSVIPUnitCents(tt.list, tt.product, tt.active, tt.priceType, tt.configured); got != tt.want {
				t.Fatalf("ResolveSVIPUnitCents() = %d, want %d", got, tt.want)
			}
		})
	}
}
