package order

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestValidIdempotencyKey(t *testing.T) {
	for _, tc := range []struct {
		key  string
		want bool
	}{
		{key: "", want: false},
		{key: "short-key", want: false},
		{key: "order-create-20260803-abc", want: true},
		{key: strings.Repeat("a", 129), want: false},
	} {
		if got := validIdempotencyKey(tc.key); got != tc.want {
			t.Fatalf("validIdempotencyKey(%q)=%v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestNormalizeOrderRemarkKeepsChineseAndBoundsLength(t *testing.T) {
	remark, err := normalizeOrderRemark("  请商家在工作日配送，感谢！  ")
	if err != nil || remark != "请商家在工作日配送，感谢！" {
		t.Fatalf("normalizeOrderRemark()=(%q,%v)", remark, err)
	}
	if _, err := normalizeOrderRemark(strings.Repeat("中", 201)); !errors.Is(err, ErrOrderRemark) {
		t.Fatalf("err=%v, want %v", err, ErrOrderRemark)
	}
}

func TestCreateRejectsOverlongRemarkBeforeDatabaseAccess(t *testing.T) {
	_, err := Create(context.Background(), nil, 7, CreateInput{AddressID: 1, IdempotencyKey: "order-create-20260803-abc", Remark: strings.Repeat("中", 201)})
	if !errors.Is(err, ErrOrderRemark) {
		t.Fatalf("err=%v, want %v", err, ErrOrderRemark)
	}
}

func TestCreateRejectsMissingIdempotencyKeyBeforeDatabaseAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil, nil, false).Register(router)
	request := httptest.NewRequest(http.MethodPost, "/v2/order/create", strings.NewReader(`{"cart_ids":[1],"address_id":1}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
