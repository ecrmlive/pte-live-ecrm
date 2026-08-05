package refundprocessor

import (
	"testing"
	"time"
)

func TestRefundCentsRoundsDecimalCurrency(t *testing.T) {
	if got := refundCents(19.99); got != 1999 {
		t.Fatalf("cents=%d, want 1999", got)
	}
	if got := refundCents(0.005); got != 1 {
		t.Fatalf("half cent=%d, want 1", got)
	}
}

func TestRecoverableOnlyReclaimsStaleProcessingRefund(t *testing.T) {
	now := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	if !recoverable("created", now, now) || !recoverable("failed", now, now) {
		t.Fatal("created and failed transactions must be retryable")
	}
	if recoverable("processing", now.Add(-processingRecoveryAfter+time.Second), now) {
		t.Fatal("fresh processing transaction must wait for provider callback")
	}
	if !recoverable("processing", now.Add(-processingRecoveryAfter-time.Second), now) {
		t.Fatal("stale processing transaction must be recovered with the same provider idempotency key")
	}
	if recoverable("succeeded", now.Add(-24*time.Hour), now) {
		t.Fatal("terminal transaction must never be retried")
	}
}
