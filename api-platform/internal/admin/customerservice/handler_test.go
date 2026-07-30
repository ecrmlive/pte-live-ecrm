package customerservice

import "testing"

func TestPositiveInt(t *testing.T) {
	if got := positiveInt("0", 20); got != 20 {
		t.Fatalf("zero got %d", got)
	}
	if got := positiveInt("101", 20); got != 101 {
		t.Fatalf("positive got %d", got)
	}
}

func TestHasRole(t *testing.T) {
	if !hasRole([]string{"customer_service"}, "customer_service") {
		t.Fatal("customer service role must be recognized")
	}
	if hasRole([]string{"operations"}, "customer_service") {
		t.Fatal("operations must not get customer service queue access")
	}
}
