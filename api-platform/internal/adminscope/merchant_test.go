package adminscope

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestDecodeMerchantScopeWithChineseFixture(t *testing.T) {
	direct, err := json.Marshal(directMerchantScope{MerchantIDs: []uint64{2001, 2002}})
	if err != nil {
		t.Fatal(err)
	}
	regions, err := json.Marshal([]uint64{10, 20})
	if err != nil {
		t.Fatal(err)
	}
	scope := decodeMerchantScope([]scopeRow{
		{ScopeType: "merchant", ScopeValue: direct},
		{ScopeType: "region", ScopeValue: regions},
	}, true, true)
	if !reflect.DeepEqual(scope.MerchantIDs, []uint64{2001, 2002}) || !reflect.DeepEqual(scope.RegionIDs, []uint64{10, 20}) {
		t.Fatalf("商户运营李航的范围 = %#v", scope)
	}
}
