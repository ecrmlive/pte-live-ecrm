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

func TestUserSetupConfigCanonicalizesChineseValues(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	got, err := svc.SaveUserSetupConfig(context.Background(), `{"register_enabled":true,"mobile_required":true,"invite_required":false,"remark":"中文注册设置"}`)
	if err != nil {
		t.Fatalf("save user setup config: %v", err)
	}
	want := `{"register_enabled":true,"mobile_required":true,"invite_required":false,"remark":"中文注册设置"}`
	if got != want || store.cache == nil || store.cache.Result != want {
		t.Fatalf("canonical user setup config = %q, stored = %#v", got, store.cache)
	}
}
