package catalog

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResponseSKUsKeepsChineseSpecsAndStableLabel(t *testing.T) {
	items := responseSKUs([]skuView{{
		MerchantSKUID: 61003,
		SKUKey:        "61003",
		SpecSnapshot:  json.RawMessage(`{"颜色":"晨雾灰","尺码":"40"}`),
		Price:         369,
		Stock:         24,
	}})
	if len(items) != 1 {
		t.Fatalf("SKU response count = %d, want 1", len(items))
	}
	if got := items[0]["spec_text"]; got != "尺码：40；颜色：晨雾灰" {
		t.Fatalf("Chinese SKU label = %q", got)
	}
	if got := items[0]["price"]; got != "369.00" {
		t.Fatalf("SKU price = %q", got)
	}
	if got := items[0]["stock"]; got != 24 {
		t.Fatalf("SKU stock = %v", got)
	}
}

func TestResponseSKUsUsesDefaultLabelForEmptySnapshot(t *testing.T) {
	items := responseSKUs([]skuView{{SKUKey: "61008", SpecSnapshot: json.RawMessage(`{}`)}})
	if got := items[0]["spec_text"]; got != "默认规格" {
		t.Fatalf("empty snapshot label = %q", got)
	}
}

func TestSellableSKUConditionGuardsCatalogRead(t *testing.T) {
	if !strings.Contains(sellableSKUCondition, "qixi_crm_b_product_sku_view") || !strings.Contains(sellableSKUCondition, "ps.sale_status = 1") {
		t.Fatalf("sellable SKU condition lost projection guard: %s", sellableSKUCondition)
	}
}
func TestResponseProductExposesSVIPConfiguration(t *testing.T) {
	item := responseProduct(productView{ProductID: 1002, Title: "七禧会员通勤包", Price: 469, SVIPPriceType: 2, SVIPPrice: 429})
	if got := item["svip_price_type"]; got != int8(2) {
		t.Fatalf("svip_price_type = %v, want 2", got)
	}
	if got := item["svip_price"]; got != "429.00" {
		t.Fatalf("svip_price = %v, want 429.00", got)
	}
}
