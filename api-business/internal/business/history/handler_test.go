package history

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPositiveID(t *testing.T) {
	for _, raw := range []string{"", "0", "-1", "bad"} {
		if _, ok := positiveID(raw); ok {
			t.Fatalf("%q should not be a valid ID", raw)
		}
	}
	if id, ok := positiveID("1001"); !ok || id != 1001 {
		t.Fatalf("positiveID() = %d, %v", id, ok)
	}
}

func TestHistoryViewKeepsChineseProductFields(t *testing.T) {
	got := (listRow{HistoryID: 7, ProductID: 1001, StoreID: 1, StoreName: "七禧服饰旗舰店", Title: "轻奢羊绒针织衫", Price: 299, Sales: 158}).view()
	if got["history_id"] != uint64(7) || got["title"] != "轻奢羊绒针织衫" || got["store_name"] != "七禧服饰旗舰店" || got["sales"] != 158 {
		t.Fatalf("unexpected history view: %#v", got)
	}
}

func TestPaginationBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/history?page=0&limit=99", nil)
	page, limit := pagination(ctx)
	if page != 1 || limit != 50 {
		t.Fatalf("pagination() = %d, %d; want 1, 50", page, limit)
	}
}
