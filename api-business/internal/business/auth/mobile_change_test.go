package auth

import (
	"context"
	"errors"
	"testing"
)

func TestChangeMobileRejectsUnsafeInputBeforeDatabaseAccess(t *testing.T) {
	svc := NewService(nil)
	if err := svc.ChangeMobile(context.Background(), 1, "13800000000", "123456", "13800000000", "654321"); !errors.Is(err, ErrBadParam) {
		t.Fatalf("same mobile error = %v, want ErrBadParam", err)
	}
	if err := svc.ChangeMobile(context.Background(), 1, "13800000000", "12345", "13900000000", "654321"); !errors.Is(err, ErrBadParam) {
		t.Fatalf("short old code error = %v, want ErrBadParam", err)
	}
}
