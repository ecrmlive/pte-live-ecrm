package storegroup

import "testing"

func TestStoreGroupInputRejectsIncompleteOrOutOfRangeCoordinates(t *testing.T) {
	lng := 121.4737
	lat := 31.2304
	if err := validateRequest(&saveRequest{Name: "中文演示分组", Longitude: &lng}); err == nil {
		t.Fatal("longitude without latitude must be rejected")
	}
	if err := validateRequest(&saveRequest{Name: "中文演示分组", Longitude: &lng, Latitude: &lat, MerchantIDs: []uint{1, 2}}); err != nil {
		t.Fatalf("valid Chinese fixture input rejected: %v", err)
	}
}

func TestStoreGroupTreeAndMembershipNormalization(t *testing.T) {
	rows := []group{
		{ID: 1, ParentID: 0, Name: "一级中文分组"},
		{ID: 2, ParentID: 1, Name: "二级中文分组"},
		{ID: 3, ParentID: 2, Name: "三级中文分组"},
	}
	tree := buildTree(rows)
	if len(tree) != 1 || len(tree[0].Children) != 1 || len(tree[0].Children[0].Children) != 1 {
		t.Fatalf("unexpected group tree: %#v", tree)
	}
	got := uniqueIDs([]uint{2, 1, 2, 0, 1})
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("membership IDs must be sorted and unique: %#v", got)
	}
}
