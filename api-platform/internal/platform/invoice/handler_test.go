package invoice

import "testing"

func TestInvoiceMasksAndStatus(t *testing.T) {
	if got := maskTax("91310000DEMO12345X"); got != "**************345X" {
		t.Fatalf("tax mask=%q", got)
	}
	if got := maskEmail("finance.invalid"); got != "" {
		t.Fatalf("invalid email=%q", got)
	}
	if got := maskEmail("demo@invoice.invalid"); got != "d***@invoice.invalid" {
		t.Fatalf("email mask=%q", got)
	}
	if !validStatus("issued") || validStatus("paid") {
		t.Fatal("invoice status validation mismatch")
	}
}
