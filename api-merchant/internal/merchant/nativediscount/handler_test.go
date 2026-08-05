package nativediscount

import "testing"

func TestNormalizeInputRequiresPackageAndProducts(t *testing.T) {
	_, _, _, _, _, err := normalizeInput(upsertInput{Name: "夏日套餐", PackagePrice: 0, ProductIDs: []uint64{1001}}, "draft")
	if err == nil {
		t.Fatal("zero package price must be rejected")
	}
	_, _, _, _, _, err = normalizeInput(upsertInput{Name: "夏日套餐", PackagePrice: 99, ProductIDs: nil}, "draft")
	if err == nil {
		t.Fatal("empty product list must be rejected")
	}
	name, raw, status, _, _, err := normalizeInput(upsertInput{
		Name: "夏日香氛套餐", PackagePrice: 199, ProductIDs: []uint64{1001, 1006}, FreeShipping: true, Remark: "中文演示套餐",
	}, "draft")
	if err != nil || name != "夏日香氛套餐" || status != "draft" || raw == "" {
		t.Fatalf("valid input rejected: name=%q status=%q raw=%q err=%v", name, status, raw, err)
	}
}

func TestValidStatus(t *testing.T) {
	for _, status := range []string{"draft", "pending", "active", "closed", "rejected"} {
		if !validStatus(status) {
			t.Fatalf("status %q should be allowed", status)
		}
	}
	if validStatus("enabled") {
		t.Fatal("unknown status must be rejected")
	}
}
