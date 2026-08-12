package content

import "testing"

func TestAgreementSettingsCatalogUsesRequiredProtocolTypes(t *testing.T) {
	want := []AgreeMeta{
		{Key: "sys_user_agree", Label: "用户协议"},
		{Key: "sys_userr_privacy", Label: "隐私政策"},
		{Key: "the_cancellation_prompt", Label: "注销提示"},
		{Key: "platform_rule", Label: "平台规则"},
		{Key: "sys_intention_agree", Label: "店铺入驻申请协议"},
		{Key: "circle_entry_agree", Label: "代理入驻申请协议"},
		{Key: "business_entry_agree", Label: "商户入驻申请协议"},
		{Key: "the_cancellation_msg", Label: "注销声明"},
		{Key: "sys_about_us", Label: "关于我们"},
		{Key: "sys_certificate", Label: "资质证照"},
	}

	got := AgreementSettingsCatalog()
	if len(got) != len(want) {
		t.Fatalf("settings agreement count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("settings agreement at %d = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestAgreeCatalogKeepsLegacyAgreementKeysAvailable(t *testing.T) {
	for _, key := range []string{"sys_svip", "sys_product_presell_agree", "sys_merchant_type"} {
		if _, ok := agreeMeta(key); !ok {
			t.Fatalf("legacy agreement key %q is no longer available", key)
		}
	}
}
