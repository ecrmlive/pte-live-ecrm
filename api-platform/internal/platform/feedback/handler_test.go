package feedback

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFeedbackPaginationBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/?page=0&limit=101", nil)
	p, l := page(c)
	if p != 1 || l != 20 {
		t.Fatalf("page=%d limit=%d", p, l)
	}
}

func TestFeedbackListRejectsInvalidStatusBeforeDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := New(nil, nil, nil)
	router.GET("/feedback", h.List)
	result := httptest.NewRecorder()
	router.ServeHTTP(result, httptest.NewRequest(http.MethodGet, "/feedback?status=unknown", nil))
	if result.Code != http.StatusBadRequest {
		t.Fatalf("status=%d want=%d", result.Code, http.StatusBadRequest)
	}
}
