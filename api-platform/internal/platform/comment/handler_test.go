package comment

import "testing"

func TestCommentFiltersAreConstrained(t *testing.T) {
	for _, status := range []string{"pending", "published", "hidden"} {
		if !validStatus(status) {
			t.Fatalf("status %q must be allowed", status)
		}
	}
	if validStatus("deleted") {
		t.Fatal("unknown comment status must fail")
	}
	if value, ok := uintParam("8801"); !ok || value != 8801 {
		t.Fatal("valid product id must parse")
	}
	if _, ok := uintParam("0"); !ok {
		t.Fatal("zero must remain distinguishable for caller rejection")
	}
}
