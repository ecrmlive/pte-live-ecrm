package merchanttype

import "testing"

func TestMerchantTypeGuaranteeRuleAndMenuNormalization(t *testing.T) {
	if err := validate(&saveReq{Name: "中文演示类型", Description: "中文说明", IsMargin: true, Margin: 0}); err == nil {
		t.Fatal("margin-enabled merchant type must reject zero margin")
	}
	req := &saveReq{Name: "中文演示类型", Description: "中文说明", IsMargin: false, Margin: 500}
	if err := validate(req); err != nil || req.Margin != 0 {
		t.Fatalf("non-margin type must normalize margin: err=%v margin=%v", err, req.Margin)
	}
	got := unique([]string{"merchant.catalog", " merchant.catalog ", "merchant.order"})
	if len(got) != 2 {
		t.Fatalf("menu codes must be de-duplicated: %#v", got)
	}
}
