package merchantonboarding

import (
	"context"
	"testing"
)

func TestApplyRejectsIncompleteChineseOnboardingCommandBeforeDatabaseAccess(t *testing.T) {
	_, err := Apply(context.Background(), nil, []byte(`{"application_id":88,"region_id":10,"merchant_name":"七禧演示茶铺","contact_name":"测试商户王小明","contact_mobile":"13800000000","account":"demo_owner","password_hash":"not-a-bcrypt-hash"}`))
	if err == nil {
		t.Fatal("expected incomplete command to be rejected")
	}
}
