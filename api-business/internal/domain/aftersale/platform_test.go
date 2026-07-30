package aftersale

import "testing"

func TestRefundAuditFromStatus(t *testing.T) {
	from, err := refundAuditFromStatus(1, StatusWait)
	if err != nil || from != StatusWait {
		t.Fatalf("merchant wait: from=%d err=%v", from, err)
	}
	if _, err := refundAuditFromStatus(1, StatusPlatform); err == nil {
		t.Fatal("merchant should not audit platform status")
	}
	from, err = refundAuditFromStatus(0, StatusPlatform)
	if err != nil || from != StatusPlatform {
		t.Fatalf("platform intervene: from=%d err=%v", from, err)
	}
	from, err = refundAuditFromStatus(0, StatusWait)
	if err != nil || from != StatusWait {
		t.Fatalf("platform wait: from=%d err=%v", from, err)
	}
	if _, err := refundAuditFromStatus(0, StatusRefunded); err == nil {
		t.Fatal("platform should reject refunded")
	}
}
