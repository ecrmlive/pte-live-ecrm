package presell

import "testing"

func TestNormalizePageBounds(t *testing.T) {
	for _, tc := range []struct {
		page, limit int
		wantPage    int
		wantLimit   int
	}{
		{page: -1, limit: 101, wantPage: 1, wantLimit: 20},
		{page: 2, limit: 15, wantPage: 2, wantLimit: 15},
	} {
		page, limit := normalize(tc.page, tc.limit)
		if page != tc.wantPage || limit != tc.wantLimit {
			t.Fatalf("normalize(%d,%d) = %d,%d, want %d,%d", tc.page, tc.limit, page, limit, tc.wantPage, tc.wantLimit)
		}
	}
}
