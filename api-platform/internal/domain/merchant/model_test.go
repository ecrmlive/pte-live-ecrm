package merchant

import "testing"

func TestPlatformMerchantModelsUseAdminProjectionTables(t *testing.T) {
	if got := (Merchant{}).TableName(); got != "qixi_crm_a_merchant_view" {
		t.Fatalf("merchant table=%q", got)
	}
	if got := (Intention{}).TableName(); got != "qixi_crm_a_merchant_application" {
		t.Fatalf("intention table=%q", got)
	}
}
