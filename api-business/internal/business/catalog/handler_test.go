package catalog

import "testing"

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
