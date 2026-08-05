package marketing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSeckillPageParamsBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		query string
		page  int
		limit int
	}{
		{query: "?page=-1&limit=101", page: 1, limit: 20},
		{query: "?page=2&limit=15", page: 2, limit: 15},
	} {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodGet, "/seckill/actives"+tc.query, nil)
		page, limit := seckillPageParams(ctx)
		if page != tc.page || limit != tc.limit {
			t.Fatalf("query %s = %d,%d, want %d,%d", tc.query, page, limit, tc.page, tc.limit)
		}
	}
}
