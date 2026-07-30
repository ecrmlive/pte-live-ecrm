package order

import (
	"errors"
	"testing"
)

func TestUniqueCartIDsAndCents(t *testing.T) {
	ids, err := uniqueCartIDs([]uint64{8, 3})
	if err != nil || len(ids) != 2 {
		t.Fatalf("ids=%v err=%v", ids, err)
	}
	if _, err := uniqueCartIDs([]uint64{8, 8}); !errors.Is(err, ErrCartOwnership) {
		t.Fatalf("err=%v", err)
	}
	value, err := cents(19.99)
	if err != nil || value != 1999 {
		t.Fatalf("cents=%d err=%v", value, err)
	}
}
