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
