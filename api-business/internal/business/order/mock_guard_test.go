package order

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProductionHandlerRejectsMockPayment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(nil, nil, false)
	r := gin.New()
	r.POST("/order/pay/:id", h.Pay)

	req := httptest.NewRequest(http.MethodPost, "/order/pay/1", strings.NewReader(`{"pay_type":"mock"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("mock payment status = %d, want %d", resp.Code, http.StatusBadRequest)
	}
}

func TestProductionCallbackDoesNotExposeMockPayment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCallbackHandler(nil, nil, false)
	r := gin.New()
	r.POST("/pay/mock", h.Mock)

	req := httptest.NewRequest(http.MethodPost, "/pay/mock", strings.NewReader(`{"group_order_id":1,"uid":1}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("mock callback status = %d, want %d", resp.Code, http.StatusNotFound)
	}
}

func TestRefundCallbackRouteRejectsUnsignedPayloadBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewCallbackHandler(nil, nil, false).Register(r)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/refund/wechat", strings.NewReader(`{"id":"forged","event_type":"REFUND.SUCCESS"}`))
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest && resp.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d, want rejected callback", resp.Code)
	}
}

func TestProductionCallbackDoesNotExposeMockRefund(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewCallbackHandler(nil, nil, false).Register(r)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/refund/mock", strings.NewReader(`{"refund_id":1}`)))
	if resp.Code != http.StatusNotFound {
		t.Fatalf("mock refund route status=%d, want %d", resp.Code, http.StatusNotFound)
	}
}

func TestSandboxMockRefundRejectsMissingIDBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewCallbackHandler(nil, nil, true).Register(r)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/refund/mock", strings.NewReader(`{}`)))
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("sandbox mock refund status=%d, want %d", resp.Code, http.StatusBadRequest)
	}
}
