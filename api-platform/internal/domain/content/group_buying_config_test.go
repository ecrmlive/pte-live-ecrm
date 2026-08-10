package content

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestGroupBuyingConfigRejectsUnknownKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveGroupBuyingConfig(context.Background(), `{"ficti_status":1,"group_buying_rate":30,"secret_key":"x"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("unknown key error = %v, want ErrBadParam", err)
	}
}

func TestGroupBuyingConfigValidatesRate(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveGroupBuyingConfig(context.Background(), `{"ficti_status":1,"group_buying_rate":101}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("rate > 100 should fail, err=%v", err)
	}
	_, err = svc.SaveGroupBuyingConfig(context.Background(), `{"ficti_status":1,"group_buying_rate":-1}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("rate < 0 should fail, err=%v", err)
	}
	_, err = svc.SaveGroupBuyingConfig(context.Background(), `{"ficti_status":2,"group_buying_rate":30}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("invalid ficti_status should fail, err=%v", err)
	}
}

func TestGroupBuyingConfigCanonicalizesAndDefaults(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	got, err := svc.SaveGroupBuyingConfig(context.Background(), `{"ficti_status":1,"group_buying_rate":30}`)
	if err != nil {
		t.Fatalf("save group buying config: %v", err)
	}
	want := `{"ficti_status":1,"group_buying_rate":30}`
	if got != want || store.cache == nil || store.cache.Result != want {
		t.Fatalf("canonical group buying config = %q, stored = %#v", got, store.cache)
	}

	store.cache = &Cache{Key: GroupBuyingConfigKey, Result: `{"ficti_status":1,"unknown":1}`}
	safe, err := svc.GetGroupBuyingConfig(context.Background())
	if err != nil {
		t.Fatalf("get group buying config: %v", err)
	}
	fallback, _ := json.Marshal(defaultGroupBuyingConfig())
	if safe == store.cache.Result || safe != string(fallback) {
		t.Fatalf("legacy malformed group buying config must not be returned, got %q", safe)
	}
}
