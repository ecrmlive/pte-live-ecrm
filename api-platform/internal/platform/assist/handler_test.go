package assist

import "testing"

func TestValidVisibility(t *testing.T) {
	if !validVisibility(&visibilityInput{IsShow: intPointer(1)}) {
		t.Fatal("show=1 must be accepted")
	}
	if !validVisibility(&visibilityInput{IsShow: intPointer(0)}) {
		t.Fatal("show=0 must be accepted")
	}
	if validVisibility(&visibilityInput{}) || validVisibility(&visibilityInput{IsShow: intPointer(2)}) {
		t.Fatal("only an explicit 0 or 1 visibility value is allowed")
	}
}

func TestOptionalMerID(t *testing.T) {
	if id, err := optionalMerID(""); err != nil || id != nil {
		t.Fatalf("empty merchant id = %v, %v", id, err)
	}
	if id, err := optionalMerID("2"); err != nil || id == nil || *id != 2 {
		t.Fatalf("valid merchant id = %v, %v", id, err)
	}
	if _, err := optionalMerID("0"); err == nil {
		t.Fatal("zero merchant id must be rejected")
	}
}

func intPointer(value int) *int { return &value }
