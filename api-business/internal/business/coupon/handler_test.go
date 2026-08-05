package coupon

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestMineRejectsUnknownStatusBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/coupons/mine", NewHandler(nil).mine)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/coupons/mine?status=unexpected", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestUserCouponViewPreservesUsedAndExpiredState(t *testing.T) {
	used := userCouponView(userCouponRow{ID: 1, Status: "used"})
	if used["status"] != 2 {
		t.Fatalf("used coupon status=%v, want 2", used["status"])
	}
	expiredAt := time.Now().Add(-time.Minute)
	expired := userCouponView(userCouponRow{ID: 2, Status: "unused", EndsAt: &expiredAt})
	if expired["status"] != -1 {
		t.Fatalf("expired coupon status=%v, want -1", expired["status"])
	}
}

func TestValidMineStatus(t *testing.T) {
	for _, status := range []string{"", "all", "unused", "used", "expired", "locked", "history", "0", "1", "2"} {
		if !validMineStatus(status) {
			t.Fatalf("status %q should be valid", status)
		}
	}
	if validMineStatus("unexpected") {
		t.Fatal("unknown status should be invalid")
	}
}
