package seckill

import (
	"context"
	"errors"
	"testing"
)

func TestUpdateRejectsInvalidStatusBeforeDataRead(t *testing.T) {
	invalid := int8(2)
	_, err := NewService(nil).Update(context.Background(), 0, 1, ActiveInput{Status: &invalid})
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("Update invalid status error = %v, want ErrBadParam", err)
	}
}

func TestValidActivityDates(t *testing.T) {
	if !validActivityDates("2026-08-01", "2026-08-31") {
		t.Fatal("same-month activity dates should be valid")
	}
	if validActivityDates("2026-08-31", "2026-08-01") || validActivityDates("2026年08月01日", "2026-08-31") {
		t.Fatal("reversed or non-ISO dates must be rejected")
	}
}
