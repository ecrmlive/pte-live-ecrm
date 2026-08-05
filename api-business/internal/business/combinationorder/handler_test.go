package combinationorder

import "testing"

func TestNormalizeDefaultsToOne(t *testing.T) {
	in := input{ProductGroupID: 701}
	if err := normalize(&in); err != nil {
		t.Fatalf("normalize: %v", err)
	}
	if in.CartNum != 1 {
		t.Fatalf("cart num=%d, want 1", in.CartNum)
	}
}

func TestNormalizeRejectsMissingActivity(t *testing.T) {
	if err := normalize(&input{CartNum: 1}); err == nil {
		t.Fatal("expected missing activity rejection")
	}
}
