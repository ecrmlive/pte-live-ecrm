package catalog

import "testing"

func TestBrandFilterKeepsChineseBrandAndRejectsUnsafeInput(t *testing.T) {
	if got, err := brandFilter(" 云锦织造 "); err != nil || got != "云锦织造" {
		t.Fatalf("brandFilter() = %q, %v", got, err)
	}
	for _, raw := range []string{"品牌\n换行", "\x00品牌"} {
		if _, err := brandFilter(raw); err == nil {
			t.Fatalf("%q should be rejected", raw)
		}
	}
}

func TestResponseProductExposesBrand(t *testing.T) {
	item := responseProduct(productView{ProductID: 1001, Title: "中文模拟针织衫", BrandName: "云锦织造"})
	if item["brand_name"] != "云锦织造" {
		t.Fatalf("brand_name = %v", item["brand_name"])
	}
}
