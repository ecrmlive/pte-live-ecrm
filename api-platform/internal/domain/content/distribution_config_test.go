package content

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestDistributionConfigRejectsUnknownKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveDistributionConfig(context.Background(), `{"extension_status":true,"secret_key":"not-allowed"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("unknown key error = %v, want ErrBadParam", err)
	}
}

func TestDistributionConfigValidatesRates(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveDistributionConfig(context.Background(), `{
		"extension_status":true,"extension_self":false,"extension_limit":false,
		"extension_limit_day":15,"promoter_type":0,"promoter_low_money":0,"extension_pop":0,
		"extension_one_rate":0.2,"extension_two_rate":0.3,"user_extract_min":1,
		"lock_brokerage_timer":0,"sys_extension_type":0,"withdraw_type":["1"],
		"extract_switch":1,"transfer_scene_id":0,"max_bag_number":10
	}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("one < two should fail, err=%v", err)
	}

	_, err = svc.SaveDistributionConfig(context.Background(), `{
		"extension_status":true,"extension_self":false,"extension_limit":false,
		"extension_limit_day":15,"promoter_type":0,"promoter_low_money":0,"extension_pop":0,
		"extension_one_rate":0.6,"extension_two_rate":0.5,"user_extract_min":1,
		"lock_brokerage_timer":0,"sys_extension_type":0,"withdraw_type":["1"],
		"extract_switch":1,"transfer_scene_id":0,"max_bag_number":10
	}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("sum > 1 should fail, err=%v", err)
	}
}

func TestDistributionConfigCanonicalizesAndDefaults(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	got, err := svc.SaveDistributionConfig(context.Background(), `{
		"extension_status":true,"extension_self":true,"extension_limit":true,
		"extension_limit_day":30,"promoter_type":2,"promoter_low_money":0,"extension_pop":1,
		"extension_one_rate":0.15,"extension_two_rate":0.05,"user_extract_min":10.5,
		"lock_brokerage_timer":7,"sys_extension_type":0,"withdraw_type":["1","0","9"],
		"extract_switch":1,"transfer_scene_id":1001,"max_bag_number":5
	}`)
	if err != nil {
		t.Fatalf("save distribution config: %v", err)
	}
	want := `{"extension_status":true,"extension_self":true,"extension_limit":true,"extension_limit_day":30,"promoter_type":2,"promoter_low_money":0,"extension_pop":1,"extension_one_rate":0.15,"extension_two_rate":0.05,"user_extract_min":10.5,"lock_brokerage_timer":7,"sys_extension_type":0,"withdraw_type":["1","0"],"extract_switch":1,"transfer_scene_id":1001,"max_bag_number":5}`
	if got != want || store.cache == nil || store.cache.Result != want {
		t.Fatalf("canonical distribution config = %q, stored = %#v", got, store.cache)
	}

	store.cache = &Cache{Key: DistributionConfigKey, Result: `{"extension_status":true,"unknown":1}`}
	safe, err := svc.GetDistributionConfig(context.Background())
	if err != nil {
		t.Fatalf("get distribution config: %v", err)
	}
	fallback, _ := json.Marshal(defaultDistributionConfig())
	if safe == store.cache.Result || safe != string(fallback) {
		t.Fatalf("legacy malformed distribution config must not be returned, got %q", safe)
	}
}
