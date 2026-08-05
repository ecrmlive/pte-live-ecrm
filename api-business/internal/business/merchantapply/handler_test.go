package merchantapply

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestApplicationPageParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		query string
		page  int
	}{
		{query: "?page=-2", page: 1},
		{query: "?page=3", page: 3},
	} {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodGet, "/merchant-applications/mine"+tc.query, nil)
		page, limit := applicationPageParams(ctx)
		if page != tc.page || limit != 20 {
			t.Fatalf("query %s: page=%d limit=%d", tc.query, page, limit)
		}
	}
}
