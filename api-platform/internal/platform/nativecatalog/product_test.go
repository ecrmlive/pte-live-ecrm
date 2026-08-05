package nativecatalog

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMergeMerchantIDsPreservesDirectAndRegionalScope(t *testing.T) {
	got := mergeMerchantIDs([]uint64{2001, 2002, 0}, []uint64{2002, 3001, 0})
	want := []uint64{2001, 2002, 3001}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("merchant scope = %v, want %v", got, want)
	}
}

func TestProductStatusMapping(t *testing.T) {
	if statusName("1") != "on_sale" || statusName("-1") != "rejected" || statusName("unexpected") != "pending_review" {
		t.Fatal("product list status filters must map only accepted values")
	}
	if statusCode("on_sale") != 1 || statusCode("rejected") != -1 || showCode("pending_review") != 0 {
		t.Fatal("product audit response status mapping is invalid")
	}
}

func TestProductAuditRequiresSellableSKU(t *testing.T) {
	if errMissingSellableSKU == nil || errMissingSellableSKU.Error() == "" {
		t.Fatal("product audit must reject products without a sellable SKU")
	}
}

func TestProductResponseKeepsProductTitleSeparateFromStoreName(t *testing.T) {
	payload, err := json.Marshal(productResponse{Title: "中文青瓷茶具", StoreName: "七禧茶铺旗舰店"})
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		StoreName string `json:"store_name"`
		Title     string `json:"title"`
	}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Title != "中文青瓷茶具" || decoded.StoreName != "七禧茶铺旗舰店" {
		t.Fatalf("product title/store contract lost: %s", payload)
	}
}
