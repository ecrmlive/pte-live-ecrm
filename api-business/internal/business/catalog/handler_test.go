package catalog

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProductSortWhitelist(t *testing.T) {
	cases := []struct {
		name      string
		sort      string
		order     string
		wantOrder string
	}{
		{name: "default ignores input order", sort: "unexpected", order: "asc", wantOrder: "sales DESC,updated_at DESC,product_id DESC"},
		{name: "sales ascending", sort: "sales", order: "asc", wantOrder: "sales ASC,updated_at DESC,product_id DESC"},
		{name: "price descending", sort: "price", order: "desc", wantOrder: "price DESC,updated_at DESC,product_id DESC"},
		{name: "unsafe direction is rejected", sort: "price", order: "desc;DROP TABLE", wantOrder: "price DESC,updated_at DESC,product_id DESC"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if got := productOrder(test.sort, test.order); got != test.wantOrder {
				t.Fatalf("productOrder() = %q, want %q", got, test.wantOrder)
			}
		})
	}
}

func TestDescendantCategoryIDs(t *testing.T) {
	rows := []categoryView{
		{CategoryID: 1, ParentID: 0},
		{CategoryID: 2, ParentID: 1},
		{CategoryID: 3, ParentID: 2},
		{CategoryID: 4, ParentID: 0},
	}

	got := descendantCategoryIDs(rows, 1)
	want := []uint64{1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("descendantCategoryIDs length = %d, want %d; got %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("descendantCategoryIDs[%d] = %d, want %d", i, got[i], want[i])
		}
	}

	unknown := descendantCategoryIDs(rows, 999)
	if len(unknown) != 1 || unknown[0] != 999 {
		t.Fatalf("unknown category fallback = %#v, want []uint64{999}", unknown)
	}
}

func TestPriceRange(t *testing.T) {
	cases := []struct {
		name    string
		query   string
		wantMin *float64
		wantMax *float64
		wantErr bool
	}{
		{name: "empty", query: ""},
		{name: "minimum only", query: "min_price=12.5", wantMin: float64Ptr(12.5)},
		{name: "inclusive range", query: "min_price=12.5&max_price=25", wantMin: float64Ptr(12.5), wantMax: float64Ptr(25)},
		{name: "negative rejected", query: "min_price=-1", wantErr: true},
		{name: "reversed rejected", query: "min_price=30&max_price=20", wantErr: true},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/catalog/products?"+test.query, nil)
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = request
			min, max, err := priceRange(ctx)
			if (err != nil) != test.wantErr {
				t.Fatalf("priceRange() error = %v, wantErr %v", err, test.wantErr)
			}
			if test.wantErr {
				return
			}
			assertFloatPtr(t, "min", min, test.wantMin)
			assertFloatPtr(t, "max", max, test.wantMax)
		})
	}
}

func float64Ptr(value float64) *float64 { return &value }

func assertFloatPtr(t *testing.T, label string, got, want *float64) {
	t.Helper()
	if got == nil || want == nil {
		if got != want {
			t.Fatalf("%s = %v, want %v", label, got, want)
		}
		return
	}
	if *got != *want {
		t.Fatalf("%s = %v, want %v", label, *got, *want)
	}
}
