package assistorder

import "testing"

func TestNormalizeQuantityRejectsInvalidSet(t *testing.T) {
	in := input{CartNum: 1}
	if err := normalizeQuantity(&in); err == nil {
		t.Fatal("expected invalid set error")
	}
}

func TestNormalizeQuantityDefaultsSingleItem(t *testing.T) {
	in := input{ProductAssistSetID: 88}
	if err := normalizeQuantity(&in); err != nil {
		t.Fatalf("normalize quantity: %v", err)
	}
	if in.CartNum != 1 {
		t.Fatalf("cart num=%d, want 1", in.CartNum)
	}
}

func TestNormalizeQuantityRejectsMoreThanOne(t *testing.T) {
	in := input{ProductAssistSetID: 88, CartNum: 2}
	if err := normalizeQuantity(&in); err == nil {
		t.Fatal("expected single-item constraint")
	}
}
