package nativerefund

import (
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

func TestRefundExportCSVProtectsFormulaAndUsesUTF8Chinese(t *testing.T) {
	content, err := refundCSV([]refund{{RefundNo: "=公式", OrderID: 99, RefundType: "return_and_refund", Amount: 18.8, Status: "awaiting_receipt", CreatedAt: time.Date(2026, 8, 4, 12, 0, 0, 0, time.Local), UpdatedAt: time.Date(2026, 8, 4, 12, 1, 0, 0, time.Local)}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(content, "\ufeff退款单号") || !strings.Contains(content, "'=公式") || !strings.Contains(content, "退货退款") {
		t.Fatalf("unexpected CSV %q", content)
	}
}

func TestValidRemark(t *testing.T) {
	if !validRemark("虚构演示备注：等待仓库核验。", "refund-note-0001") {
		t.Fatal("expected valid Chinese remark")
	}
	if validRemark("", "refund-note-0001") || validRemark("说明", "short") {
		t.Fatal("accepted invalid remark command")
	}
}
