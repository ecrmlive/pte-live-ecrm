package content

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestIntegralConfigRejectsUnknownKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveIntegralConfig(context.Background(), `{"integral_status":1,"secret_key":"x"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("unknown key error = %v, want ErrBadParam", err)
	}
}

func TestIntegralConfigValidatesRanges(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveIntegralConfig(context.Background(), `{"integral_status":2,"integral_money":0.1,"integral_order_rate":1,"integral_freeze":0,"integral_clear_time":24,"integral_user_give":50,"integral_community_give":10,"integral_community_give_limit":10,"rule":""}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("invalid status should fail, err=%v", err)
	}
	_, err = svc.SaveIntegralConfig(context.Background(), `{"integral_status":1,"integral_money":-0.1,"integral_order_rate":1,"integral_freeze":0,"integral_clear_time":24,"integral_user_give":50,"integral_community_give":10,"integral_community_give_limit":10,"rule":""}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("negative money should fail, err=%v", err)
	}
	_, err = svc.SaveIntegralConfig(context.Background(), `{"integral_status":1,"integral_money":0.1,"integral_order_rate":1,"integral_freeze":0,"integral_clear_time":24,"integral_user_give":50,"integral_community_give":10000,"integral_community_give_limit":10,"rule":""}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("community give > 9999 should fail, err=%v", err)
	}
}

func TestIntegralConfigCanonicalizesAndDefaults(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	got, err := svc.SaveIntegralConfig(context.Background(), `{"integral_status":1,"integral_money":0.1,"integral_order_rate":1,"integral_freeze":0,"integral_clear_time":24,"integral_user_give":50,"integral_community_give":10,"integral_community_give_limit":10,"rule":"<p>积分说明</p>"}`)
	if err != nil {
		t.Fatalf("save integral config: %v", err)
	}
	var parsed integralConfig
	if err := json.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("unmarshal saved: %v", err)
	}
	if parsed.IntegralStatus != 1 ||
		parsed.IntegralMoney < 0.099 || parsed.IntegralMoney > 0.101 ||
		parsed.IntegralOrderRate != 1 ||
		parsed.IntegralFreeze != 0 || parsed.IntegralClearTime != 24 || parsed.IntegralUserGive != 50 ||
		parsed.IntegralCommunityGive != 10 || parsed.IntegralCommunityGiveLimit != 10 ||
		parsed.Rule != "<p>积分说明</p>" {
		t.Fatalf("unexpected saved config: %#v", parsed)
	}
	if store.cache == nil || store.cache.Result != got {
		t.Fatalf("store mismatch: %#v", store.cache)
	}

	safe, err := svc.GetIntegralConfig(context.Background())
	if err != nil {
		t.Fatalf("get after save: %v", err)
	}
	if safe != got {
		t.Fatalf("get = %q, want %q", safe, got)
	}

	store.cache = &Cache{Key: IntegralConfigKey, Result: `{"integral_status":1,"unknown":1}`}
	fallbackRaw, err := svc.GetIntegralConfig(context.Background())
	if err != nil {
		t.Fatalf("get malformed: %v", err)
	}
	fallback, _ := json.Marshal(defaultIntegralConfig())
	if fallbackRaw != string(fallback) {
		t.Fatalf("legacy malformed must fallback, got %q", fallbackRaw)
	}
}
