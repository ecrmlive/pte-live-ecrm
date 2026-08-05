package presellorder

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNormalizeDefaultsToOne(t *testing.T) {
	in := input{ProductPresellID: 701}
	if err := normalize(&in); err != nil {
		t.Fatalf("normalize: %v", err)
	}
	if in.CartNum != 1 {
		t.Fatalf("cart num=%d, want 1", in.CartNum)
	}
}

func TestNormalizeRejectsBadInput(t *testing.T) {
	for _, in := range []input{{}, {ProductPresellID: 1, CartNum: -1}, {ProductPresellID: 1, CartNum: 100}} {
		if err := normalize(&in); err == nil {
			t.Fatalf("expected rejection: %+v", in)
		}
	}
}

func TestFinalPageParamsBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		query               string
		wantPage, wantLimit int
	}{
		{query: "?page=-3&limit=101", wantPage: 1, wantLimit: 20},
		{query: "?page=2&limit=15", wantPage: 2, wantLimit: 15},
	} {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodGet, "/presell/finals"+tc.query, nil)
		page, limit := finalPageParams(ctx)
		if page != tc.wantPage || limit != tc.wantLimit {
			t.Fatalf("query %s: page=%d limit=%d, want %d,%d", tc.query, page, limit, tc.wantPage, tc.wantLimit)
		}
	}
}
