package comment

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestCommentViewKeepsChineseFields(t *testing.T) {
	got := (row{
		ID: 18, OrderItemID: 208, ProductID: 1001, Score: 5, Content: "尺码合适，面料柔软", Reply: "感谢您的认可", Status: "published", Title: "轻奢羊绒针织衫",
		CreatedAt: time.Date(2026, 8, 3, 10, 20, 0, 0, time.Local),
	}).view()
	if got["comment_id"] != uint64(18) || got["order_item_id"] != uint64(208) || got["title"] != "轻奢羊绒针织衫" || got["content"] != "尺码合适，面料柔软" || got["reply_content"] != "感谢您的认可" || got["create_time"] != "2026-08-03 10:20:00" {
		t.Fatalf("unexpected comment view: %#v", got)
	}
}

func TestCommentIDRejectsInvalidOrZero(t *testing.T) {
	for _, raw := range []string{"", "0", "-1", "comment"} {
		if _, ok := id(raw); ok {
			t.Fatalf("%q should not be a valid comment ID", raw)
		}
	}
	if value, ok := id("1001"); !ok || value != 1001 {
		t.Fatalf("id() = %d, %v; want 1001, true", value, ok)
	}
}

func TestCommentRoutesRejectInvalidInputBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(nil)
	handler.Register(router)
	handler.RegisterPublic(router)

	invalidCreate := httptest.NewRecorder()
	router.ServeHTTP(invalidCreate, httptest.NewRequest(http.MethodPost, "/comments", strings.NewReader(`{"order_item_id":0,"score":6,"content":""}`)))
	if invalidCreate.Code != http.StatusBadRequest {
		t.Fatalf("invalid create status = %d; want %d", invalidCreate.Code, http.StatusBadRequest)
	}

	emptyContent := httptest.NewRecorder()
	router.ServeHTTP(emptyContent, httptest.NewRequest(http.MethodPost, "/comments", strings.NewReader(`{"order_item_id":1,"score":5,"content":"   "}`)))
	if emptyContent.Code != http.StatusBadRequest {
		t.Fatalf("empty content status = %d; want %d", emptyContent.Code, http.StatusBadRequest)
	}

	invalidProduct := httptest.NewRecorder()
	router.ServeHTTP(invalidProduct, httptest.NewRequest(http.MethodGet, "/products/0/comments", nil))
	if invalidProduct.Code != http.StatusBadRequest {
		t.Fatalf("invalid product status = %d; want %d", invalidProduct.Code, http.StatusBadRequest)
	}
}

func TestCommentPaginationBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/comments/mine?page=0&limit=101", nil)
	current, limit := page(ctx)
	if current != 1 || limit != 20 {
		t.Fatalf("page() = %d, %d; want 1, 20", current, limit)
	}
}
