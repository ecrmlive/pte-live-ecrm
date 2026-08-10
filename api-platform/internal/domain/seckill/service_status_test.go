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
	if !validActivityDates("2026-08-08T00:00:00+08:00", "2026-11-07T00:00:00+08:00") {
		t.Fatal("RFC3339 DATE from MySQL parseTime should be accepted")
	}
	if validActivityDates("2026-08-31", "2026-08-01") || validActivityDates("2026年08月01日", "2026-08-31") {
		t.Fatal("reversed or non-ISO dates must be rejected")
	}
}

func TestNormalizeDay(t *testing.T) {
	if got := normalizeDay("2026-08-08T00:00:00+08:00"); got != "2026-08-08" {
		t.Fatalf("normalizeDay RFC3339 = %q, want 2026-08-08", got)
	}
	if got := normalizeDay("2026-08-08"); got != "2026-08-08" {
		t.Fatalf("normalizeDay plain = %q", got)
	}
}
