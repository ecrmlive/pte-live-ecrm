package contentview

import "testing"

func TestAgreementLabel(t *testing.T) {
	label, ok := agreementLabel("sys_user_agree")
	if !ok || label != "用户协议" {
		t.Fatalf("agreementLabel(sys_user_agree) = (%q, %v)", label, ok)
	}
	if _, ok := agreementLabel("not-a-contract"); ok {
		t.Fatal("unknown agreement key must be rejected")
	}
}
