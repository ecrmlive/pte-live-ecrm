package nativerefund

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestReturnRefundStatusCompatibility(t *testing.T) {
	if got := normalizeStatus("1"); got != "awaiting_return" {
		t.Fatalf("status 1 = %q", got)
	}
	if got := normalizeStatus("2"); got != "awaiting_receipt" {
		t.Fatalf("status 2 = %q", got)
	}
	if got := legacyStatus("awaiting_return"); got != 1 {
		t.Fatalf("awaiting_return = %d", got)
	}
	if got := legacyStatus("awaiting_receipt"); got != 2 {
		t.Fatalf("awaiting_receipt = %d", got)
	}
}

func TestRefundExportCSVUsesChineseHeadersAndNeutralisesFormula(t *testing.T) {
	content, err := refundCSV([]refund{{
		OrderID: 123, MerchantID: 8, StoreID: 9, RefundNo: "=恶意公式", Amount: 19.9,
		RefundType: "return_and_refund", Status: "awaiting_receipt",
		CreatedAt: time.Date(2026, 8, 4, 12, 0, 0, 0, time.Local), UpdatedAt: time.Date(2026, 8, 4, 12, 1, 0, 0, time.Local),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(content, "\ufeff退款单号,订单ID") {
		t.Fatalf("CSV header=%q", content)
	}
	if !strings.Contains(content, "'=恶意公式") || !strings.Contains(content, "退货退款") || !strings.Contains(content, "待收货") {
		t.Fatalf("CSV did not retain protected Chinese row: %q", content)
	}
}

func TestRefundExportStatusAndFingerprint(t *testing.T) {
	if got, ok := exportStatus(-1); !ok || got != "rejected" {
		t.Fatalf("status -1 = %q, %v", got, ok)
	}
	if _, ok := exportStatus(99); ok {
		t.Fatal("unexpected valid unknown status")
	}
	if a, b := refundExportFingerprint("applied", []uint64{9, 2}), refundExportFingerprint("applied", []uint64{2, 9}); a != b {
		t.Fatalf("scope fingerprint must ignore merchant ID order: %q != %q", a, b)
	}
}

func TestRefundExportStatusAcceptsCurrentStateCodeWithoutLegacyMislabel(t *testing.T) {
	currentBytes, _ := json.Marshal("refunding")
	current := json.RawMessage(currentBytes)
	if got, ok := parseExportStatus(&current); !ok || got != "refunding" {
		t.Fatalf("current refunding status = %q, %v", got, ok)
	}
	legacyBytes, _ := json.Marshal(3)
	legacy := json.RawMessage(legacyBytes)
	if got, ok := parseExportStatus(&legacy); !ok || got != "refunded" {
		t.Fatalf("legacy status 3 = %q, %v", got, ok)
	}
	invalidBytes, _ := json.Marshal("not-a-refund-state")
	invalid := json.RawMessage(invalidBytes)
	if _, ok := parseExportStatus(&invalid); ok {
		t.Fatal("unknown current status was accepted")
	}
}
