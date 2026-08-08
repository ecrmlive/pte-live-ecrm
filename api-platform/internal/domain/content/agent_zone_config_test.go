package content

import (
	"context"
	"errors"
	"testing"
)

func TestAgentZoneConfigRejectsUnknownKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveAgentZoneConfig(context.Background(), `{"one_agent_commission":1,"two_agent_commission":0,"three_agent_commission":0,"form_fields":[],"secret":1}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("unknown key error = %v, want ErrBadParam", err)
	}
}

func TestAgentZoneConfigRejectsInvalidOrder(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveAgentZoneConfig(context.Background(), `{"one_agent_commission":3,"two_agent_commission":5,"three_agent_commission":1,"form_fields":[]}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("order error = %v, want ErrBadParam", err)
	}
}

func TestAgentZoneConfigCanonicalizes(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	raw := `{"one_agent_commission":8,"two_agent_commission":5,"three_agent_commission":3,"form_fields":[{"id":"f1","type":"text","title":"备注","content_type":"text","default_value":"","placeholder":"请输入","required":false}]}`
	got, err := svc.SaveAgentZoneConfig(context.Background(), raw)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	want := `{"one_agent_commission":8,"two_agent_commission":5,"three_agent_commission":3,"form_fields":[{"id":"f1","type":"text","title":"备注","content_type":"text","default_value":"","placeholder":"请输入","required":false}]}`
	if got != want || store.cache == nil || store.cache.Key != agentZoneConfigKey || store.cache.Result != want {
		t.Fatalf("canonical = %q, stored = %#v", got, store.cache)
	}
}

func TestAgentZoneConfigDefaultsOnMalformed(t *testing.T) {
	store := &mallSettingStore{cache: &Cache{Key: agentZoneConfigKey, Result: `{"one_agent_commission":1,"bad":true}`}}
	svc := NewService(store)
	safe, err := svc.GetAgentZoneConfig(context.Background())
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if safe != marshalAgentZoneConfig(defaultAgentZoneConfig()) {
		t.Fatalf("legacy malformed must not be returned, got %q", safe)
	}
}

func TestAgentZoneConfigAllowsEqualRates(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	raw := `{"one_agent_commission":0,"two_agent_commission":0,"three_agent_commission":0,"form_fields":[]}`
	got, err := svc.SaveAgentZoneConfig(context.Background(), raw)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	want := `{"one_agent_commission":0,"two_agent_commission":0,"three_agent_commission":0,"form_fields":[]}`
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
