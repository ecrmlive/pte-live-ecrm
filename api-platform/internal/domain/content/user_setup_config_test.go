package content

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestUserSetupConfigRejectsUnknownKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveUserSetupConfig(context.Background(), `{"user_default_avatar":"","secret_key":"x"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("unknown key error = %v, want ErrBadParam", err)
	}
}

func TestUserSetupConfigValidatesGiftRules(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	base := defaultUserSetupConfig()
	base.NewcomerStatus = 1
	base.RegisterMoneyStatus = 1
	base.RegisterGiveMoney = 0
	raw, _ := json.Marshal(base)
	_, err := svc.SaveUserSetupConfig(context.Background(), string(raw))
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("money gift without amount should fail, err=%v", err)
	}

	base = defaultUserSetupConfig()
	base.NewcomerStatus = 1
	base.RegisterIntegralStatus = 1
	base.RegisterGiveIntegral = 0
	raw, _ = json.Marshal(base)
	_, err = svc.SaveUserSetupConfig(context.Background(), string(raw))
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("integral gift without value should fail, err=%v", err)
	}
}

func TestUserSetupConfigCanonicalizesAndDefaults(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	cfg := defaultUserSetupConfig()
	cfg.UserDefaultAvatar = "https://example.com/a.png"
	cfg.IsPhoneLogin = 1
	cfg.NewcomerStatus = 1
	cfg.RegisterMoneyStatus = 1
	cfg.RegisterGiveMoney = 1.5
	cfg.Fields[0].IsRequire = 1
	raw, _ := json.Marshal(cfg)
	got, err := svc.SaveUserSetupConfig(context.Background(), string(raw))
	if err != nil {
		t.Fatalf("save user setup: %v", err)
	}
	var parsed userSetupConfig
	if err := json.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if parsed.UserDefaultAvatar != "https://example.com/a.png" ||
		parsed.IsPhoneLogin != 1 ||
		parsed.RegisterGiveMoney < 1.49 || parsed.RegisterGiveMoney > 1.51 ||
		len(parsed.Fields) != 6 ||
		parsed.Fields[0].IsDefault != 1 {
		t.Fatalf("unexpected saved: %#v", parsed)
	}

	store.cache = &Cache{Key: UserSetupConfigKey, Result: `{"register_enabled":true}`}
	fallbackRaw, err := svc.GetUserSetupConfig(context.Background())
	if err != nil {
		t.Fatalf("get malformed: %v", err)
	}
	fallback, _ := json.Marshal(defaultUserSetupConfig())
	if fallbackRaw != string(fallback) {
		t.Fatalf("legacy malformed must fallback")
	}
}

func TestUserSetupConfigInjectsMissingDefaultField(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	cfg := defaultUserSetupConfig()
	cfg.Fields = cfg.Fields[1:] // drop real_name
	raw, _ := json.Marshal(cfg)
	got, err := svc.SaveUserSetupConfig(context.Background(), string(raw))
	if err != nil {
		t.Fatalf("missing default field should be injected, err=%v", err)
	}
	var parsed userSetupConfig
	if err := json.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	found := false
	for _, f := range parsed.Fields {
		if f.Field == "real_name" && f.IsDefault == 1 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("real_name should be injected: %#v", parsed.Fields)
	}
}
