package order

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestUniqueCouponIDsAllowsPlatformAndStorePair(t *testing.T) {
	ids, err := uniqueCouponIDs([]uint64{101, 202})
	if err != nil {
		t.Fatalf("distinct coupon pair err=%v", err)
	}
	if len(ids) != 2 || ids[0] != 101 || ids[1] != 202 {
		t.Fatalf("ids=%v", ids)
	}
}

func TestUniqueCouponIDsRejectsDuplicateOrZero(t *testing.T) {
	if _, err := uniqueCouponIDs([]uint64{101, 101}); !errors.Is(err, ErrCouponConflict) {
		t.Fatalf("duplicate err=%v", err)
	}
	if _, err := uniqueCouponIDs([]uint64{0}); !errors.Is(err, ErrCouponOwnership) {
		t.Fatalf("zero err=%v", err)
	}
}

func TestDeliverySummary(t *testing.T) {
	at := time.Date(2026, time.August, 3, 10, 30, 0, 0, time.Local)
	view := deliverySummary(deliveryRow{OrderID: 7, DeliveryType: "express", CarrierCode: "顺丰速运", TrackingNo: "SFDEMO20260803", Status: "shipped", DeliveredAt: &at})
	if view["delivery_name"] != "顺丰速运" || view["delivery_id"] != "SFDEMO20260803" || view["delivery_status"] != "shipped" {
		t.Fatalf("summary=%#v", view)
	}
}

func TestConfirmReceiptRejectsInvalidIDBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil, nil, false).Register(router)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/order/0/confirm-receipt", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
