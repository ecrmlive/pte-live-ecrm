package refund

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUserRefundStateMapping(t *testing.T) {
	cases := map[string]int{
		"applied": 0, "merchant_handling": 0, "platform_intervene": 4,
		"awaiting_return": 1, "awaiting_receipt": 2, "refunding": 0, "refunded": 3, "rejected": -1, "cancelled": -2,
	}
	for state, want := range cases {
		if got := legacyStatus(state); got != want {
			t.Fatalf("legacyStatus(%q) = %d, want %d", state, got, want)
		}
	}
}

func TestRefundTypeMapping(t *testing.T) {
	if got := refundType(1); got != "money_only" {
		t.Fatalf("money-only type = %q", got)
	}
	if got := refundType(2); got != "return_and_refund" {
		t.Fatalf("return-and-refund type = %q", got)
	}
	if got := legacyRefundType("return_and_refund"); got != 2 {
		t.Fatalf("legacy return-and-refund type = %d", got)
	}
}

func TestFullOrderRefundAmountUsesServerPaidAmount(t *testing.T) {
	if got := fullOrderRefundAmount(88.6); got != 88.6 {
		t.Fatalf("refund amount=%v, want 88.6", got)
	}
	if got := fullOrderRefundAmount(0); got != 0 {
		t.Fatalf("zero paid amount=%v, want 0", got)
	}
}

func TestPageBounds(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "/refunds?page=-3&limit=999", nil)
	if gotPage, gotLimit := page(ctx); gotPage != 1 || gotLimit != 100 {
		t.Fatalf("page bounds = %d,%d, want 1,100", gotPage, gotLimit)
	}
	second, _ := gin.CreateTestContext(httptest.NewRecorder())
	second.Request = httptest.NewRequest(http.MethodGet, "/refunds?page=2&limit=15", nil)
	if gotPage, gotLimit := page(second); gotPage != 2 || gotLimit != 15 {
		t.Fatalf("page values = %d,%d, want 2,15", gotPage, gotLimit)
	}
}

func TestRefundApplyOrderStates(t *testing.T) {
	for _, state := range []string{"paid", "fulfilling", "shipped"} {
		if !canApply(state) {
			t.Fatalf("expected %q to allow after-sale apply", state)
		}
	}
	for _, state := range []string{"pending_pay", "completed", "cancelled", "aftersale"} {
		if canApply(state) {
			t.Fatalf("expected %q to reject after-sale apply", state)
		}
	}
}

func TestRefundReasonUsesUnicodeCharacterLimit(t *testing.T) {
	if exceedsRefundMessageLimit(strings.Repeat("中", 500)) {
		t.Fatal("500 Chinese characters must be accepted")
	}
	if !exceedsRefundMessageLimit(strings.Repeat("中", 501)) {
		t.Fatal("501 Chinese characters must exceed refund reason limit")
	}
}
