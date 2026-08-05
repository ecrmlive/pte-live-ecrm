package combination

import (
	"context"
	"errors"
	"testing"
)

func TestUpdateRejectsInvalidStatusBeforeDataRead(t *testing.T) {
	invalid := 2
	_, err := NewService(nil).Update(context.Background(), 0, 1, SaveInput{Status: &invalid})
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("Update invalid status error = %v, want ErrBadParam", err)
	}
}

func TestValidActivityRange(t *testing.T) {
	start, end := parseRange("2026-08-01 09:00:00", "2026-08-31 21:00:00")
	if !validActivityRange("2026-08-01 09:00:00", "2026-08-31 21:00:00", start, end) {
		t.Fatal("valid activity range must be accepted")
	}
	start, end = parseRange("2026-08-31", "2026-08-01")
	if validActivityRange("2026-08-31", "2026-08-01", start, end) {
		t.Fatal("reversed activity range must be rejected")
	}
	if validActivityTime("2026年08月01日") {
		t.Fatal("non-ISO date must be rejected")
	}
}
