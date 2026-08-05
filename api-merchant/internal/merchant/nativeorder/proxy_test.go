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

func TestRegisterIncludesVerifyRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(nil, nil).Register(r)

	req := httptest.NewRequest(http.MethodPost, "/orders/1/verify", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Fatalf("POST /orders/:id/verify must be registered")
	}
}

func TestHashVerifyCodeStable(t *testing.T) {
	a := hashVerifyCode(" ABCD1234 ")
	b := hashVerifyCode("ABCD1234")
	if a == "" || a != b {
		t.Fatalf("hashVerifyCode unstable: %q vs %q", a, b)
	}
	if len(a) != 64 {
		t.Fatalf("hash length=%d, want 64", len(a))
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
