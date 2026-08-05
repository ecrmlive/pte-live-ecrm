package distribution

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCommissionMark(t *testing.T) {
	cases := map[string]string{"pending": "待结算", "available": "可结算", "settled": "已结算", "unknown": "佣金记录"}
	for status, want := range cases {
		if got := commissionMark(status); got != want {
			t.Fatalf("commissionMark(%q)=%q, want %q", status, got, want)
		}
	}
}

func TestBindRejectsEmptyPromoterBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil).Register(router)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/spread/bind", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestPageParamsBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=-4&limit=500", nil)
	page, limit := pageParams(ctx)
	if page != 1 || limit != 100 {
		t.Fatalf("page=%d limit=%d", page, limit)
	}
}
