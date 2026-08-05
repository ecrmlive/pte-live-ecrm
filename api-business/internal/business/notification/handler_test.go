package notification

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNotificationPageParamsBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=-1&limit=300", nil)
	page, limit := notificationPageParams(ctx)
	if page != 1 || limit != 100 {
		t.Fatalf("page=%d limit=%d", page, limit)
	}
}
