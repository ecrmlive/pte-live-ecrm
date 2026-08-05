package account

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestSignListRejectsMalformedMonthBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil).Register(router)
	for _, month := range []string{"2026-13", "2026-1", "x2026-08", "2026/08"} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/sign/lst?month="+month, nil)
		router.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("month %q status=%d, want %d", month, recorder.Code, http.StatusBadRequest)
		}
	}
}

func TestShanghaiTodayUsesISODate(t *testing.T) {
	got := shanghaiToday()
	if _, err := time.Parse("2006-01-02", got); err != nil {
		t.Fatalf("invalid date %q: %v", got, err)
	}
}

func TestLedgerPageParamsBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=-2&limit=101", nil)
	page, limit := ledgerPageParams(ctx)
	if page != 1 || limit != 100 {
		t.Fatalf("page=%d limit=%d", page, limit)
	}
	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=3&limit=0", nil)
	page, limit = ledgerPageParams(ctx)
	if page != 3 || limit != 20 {
		t.Fatalf("page=%d limit=%d", page, limit)
	}
}
