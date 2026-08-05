package nativeorder

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterIncludesProxyRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(nil, nil).Register(r)

	req := httptest.NewRequest(http.MethodPost, "/orders/proxy", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatalf("POST /orders/proxy must be registered")
	}
}

func TestRoundMoney(t *testing.T) {
	if got := roundMoney(19.999); got != 20 {
		t.Fatalf("roundMoney(19.999)=%v", got)
	}
	if got := roundMoney(10.123); got != 10.12 {
		t.Fatalf("roundMoney(10.123)=%v", got)
	}
}

func TestIsDuplicateKey(t *testing.T) {
	if !isDuplicateKey(fmtError("Error 1062: Duplicate entry for key 'uk_user_idempotency'")) {
		t.Fatal("expected duplicate detection")
	}
	if isDuplicateKey(fmtError("connection refused")) {
		t.Fatal("did not expect duplicate")
	}
}
