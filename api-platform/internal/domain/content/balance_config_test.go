package content

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestBalanceConfigRejectsUnknownKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveBalanceConfig(context.Background(), `{"balance_func_status":1,"secret_key":"x"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("unknown key error = %v, want ErrBadParam", err)
	}
}

func TestBalanceConfigValidatesRanges(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveBalanceConfig(context.Background(), `{"balance_func_status":2,"recharge_switch":1,"store_user_min_recharge":1,"recharge_attention":""}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("invalid balance status should fail, err=%v", err)
	}
	_, err = svc.SaveBalanceConfig(context.Background(), `{"balance_func_status":1,"recharge_switch":3,"store_user_min_recharge":1,"recharge_attention":""}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("invalid recharge switch should fail, err=%v", err)
	}
	_, err = svc.SaveBalanceConfig(context.Background(), `{"balance_func_status":1,"recharge_switch":1,"store_user_min_recharge":-1,"recharge_attention":""}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("negative min recharge should fail, err=%v", err)
	}
}

func TestBalanceConfigCanonicalizesAndDefaults(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	payload := `{"balance_func_status":1,"recharge_switch":1,"store_user_min_recharge":1,"recharge_attention":"1、注意A\n2、注意B\n3、注意C"}`
	got, err := svc.SaveBalanceConfig(context.Background(), payload)
	if err != nil {
		t.Fatalf("save balance config: %v", err)
	}
	var parsed balanceConfig
	if err := json.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("unmarshal saved: %v", err)
	}
	if parsed.BalanceFuncStatus != 1 ||
		parsed.RechargeSwitch != 1 ||
		parsed.StoreUserMinRecharge < 0.999 || parsed.StoreUserMinRecharge > 1.001 ||
		parsed.RechargeAttention != "1、注意A\n2、注意B\n3、注意C" {
		t.Fatalf("unexpected saved config: %#v", parsed)
	}
	if store.cache == nil || store.cache.Result != got {
		t.Fatalf("store mismatch: %#v", store.cache)
	}

	safe, err := svc.GetBalanceConfig(context.Background())
	if err != nil {
		t.Fatalf("get after save: %v", err)
	}
	if safe != got {
		t.Fatalf("get = %q, want %q", safe, got)
	}

	store.cache = &Cache{Key: BalanceConfigKey, Result: `{"balance_func_status":1,"unknown":1}`}
	fallbackRaw, err := svc.GetBalanceConfig(context.Background())
	if err != nil {
		t.Fatalf("get malformed: %v", err)
	}
	fallback, _ := json.Marshal(defaultBalanceConfig())
	if fallbackRaw != string(fallback) {
		t.Fatalf("legacy malformed must fallback, got %q", fallbackRaw)
	}
}
