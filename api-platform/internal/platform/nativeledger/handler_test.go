package nativeledger

import "testing"

func TestAssetTypeAllowsOnlyBusinessAssetKinds(t *testing.T) {
	for _, value := range []string{"", "balance", "points", "commission"} {
		if got, ok := assetType(value); !ok || got != value {
			t.Fatalf("assetType(%q) = (%q, %t)", value, got, ok)
		}
	}
	if _, ok := assetType("银行卡余额"); ok {
		t.Fatal("arbitrary asset type must be rejected")
	}
}
