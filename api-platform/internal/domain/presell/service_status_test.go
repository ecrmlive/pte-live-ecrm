package presell

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestUpdateRejectsInvalidStatusBeforeDataRead(t *testing.T) {
	invalid := 2
	_, err := NewService(nil).Update(context.Background(), 0, 1, SaveInput{Status: &invalid})
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("Update invalid status error = %v, want ErrBadParam", err)
	}
}

func TestValidPresellTimeRanges(t *testing.T) {
	start, _ := parseTime("2026-08-10 10:00:00")
	end, _ := parseTime("2026-08-11 10:00:00")
	if !validActivityRange("2026-08-10 10:00:00", "2026-08-11 10:00:00", start, end) {
		t.Fatal("expected valid activity range")
	}
	if validActivityRange("not-a-date", "", time.Now(), time.Now().Add(time.Hour)) {
		t.Fatal("invalid activity time must be rejected")
	}
	if validActivityRange("2026-08-11 10:00:00", "2026-08-10 10:00:00", start, end) {
		t.Fatal("reversed activity range must be rejected")
	}
	if validFinalRange("2026-08-12 10:00:00", "2026-08-11 10:00:00") {
		t.Fatal("reversed final payment range must be rejected")
	}
}
