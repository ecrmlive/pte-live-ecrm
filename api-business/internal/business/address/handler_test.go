package address

import "testing"

func TestMergeAddressSupportsPartialUpdate(t *testing.T) {
	item := row{Recipient: "李四", Mobile: "13900000000", Detail: "原详细地址", IsDefault: 0}
	city := "深圳市"
	defaultFlag := int8(1)
	if err := merge(&item, request{City: &city, IsDefault: &defaultFlag}, false); err != nil {
		t.Fatalf("merge() error = %v", err)
	}
	if item.City != "深圳市" || item.Recipient != "李四" || item.IsDefault != 1 {
		t.Fatalf("merge() changed data incorrectly: %#v", item)
	}
}

func TestMergeAddressRequiresFieldsOnCreate(t *testing.T) {
	if err := merge(&row{}, request{}, true); err == nil {
		t.Fatal("merge() must reject an incomplete new address")
	}
}
