package content

import (
	"context"
	"errors"
	"testing"
)

func TestStorageConfigRejectsSensitiveKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveStorageConfig(context.Background(), `{"provider":"cos","secret_key":"not-allowed"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("sensitive storage key error = %v, want ErrBadParam", err)
	}
}

func TestAppStubConfigRejectsSensitiveKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveAppStubConfig(context.Background(), RoutineAppConfigKey, `{"name":"小程序","app_secret":"not-allowed"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("sensitive app stub key error = %v, want ErrBadParam", err)
	}
}

