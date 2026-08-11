package invoice

import "testing"

func TestInvoiceMasksAndLabels(t *testing.T) {
	if got := maskTax("91310000DEMO12345X"); got != "**************345X" {
		t.Fatalf("tax mask=%q", got)
	}
	if got := maskEmail("finance.invalid"); got != "" {
		t.Fatalf("invalid email=%q", got)
	}
	if got := maskEmail("demo@invoice.invalid"); got != "d***@invoice.invalid" {
		t.Fatalf("email mask=%q", got)
	}
	if invoiceStatusLabel("issued") != "已开" ||
		invoiceStatusLabel("requested") != "未开" ||
		invoiceStatusLabel("rejected") != "未开" {
		t.Fatal("status label mismatch")
	}
	if invoiceTypeLabel(2) != "专用发票" || titleTypeLabel("enterprise") != "企业" {
		t.Fatal("type label mismatch")
	}
	if detailTitle(1, "personal") != "个人电子普通发票" {
		t.Fatal("personal detail title mismatch")
	}
	if detailTitle(1, "enterprise") != "企业电子普通发票" {
		t.Fatal("enterprise detail title mismatch")
	}
	if detailTitle(2, "enterprise") != "企业专用纸质发票" {
		t.Fatal("special detail title mismatch")
	}
}
