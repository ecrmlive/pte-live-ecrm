package assist

import "testing"

func TestNormalizePageBounds(t *testing.T) {
	for _, tc := range []struct {
		page, limit int
		wantPage    int
		wantLimit   int
	}{
		{page: 0, limit: 101, wantPage: 1, wantLimit: 20},
		{page: 3, limit: 12, wantPage: 3, wantLimit: 12},
	} {
		page, limit := normalize(tc.page, tc.limit)
		if page != tc.wantPage || limit != tc.wantLimit {
			t.Fatalf("normalize(%d,%d) = %d,%d, want %d,%d", tc.page, tc.limit, page, limit, tc.wantPage, tc.wantLimit)
		}
	}
}
