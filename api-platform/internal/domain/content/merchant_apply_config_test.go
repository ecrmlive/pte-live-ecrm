package content

import (
	"context"
	"errors"
	"testing"
)

func TestMerchantApplyConfigRejectsUnknownKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveMerchantApplyConfig(context.Background(), `{"background_image":"","secret_key":"x"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("unknown key error = %v, want ErrBadParam", err)
	}
}

func TestMerchantApplyConfigCanonicalizesFields(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	raw := `{"background_image":"https://example.test/bg.png","form_fields":[{"id":"f1","type":"text","title":"营业面积","content_type":"number","default_value":"","placeholder":"请输入","required":true}]}`
	got, err := svc.SaveMerchantApplyConfig(context.Background(), raw)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	want := `{"background_image":"https://example.test/bg.png","form_fields":[{"id":"f1","type":"text","title":"营业面积","content_type":"number","default_value":"","placeholder":"请输入","required":true}]}`
	if got != want || store.cache == nil || store.cache.Result != want {
		t.Fatalf("canonical = %q, stored = %#v", got, store.cache)
	}
}

func TestMerchantApplyConfigNeverEchoesLegacyMalformed(t *testing.T) {
	store := &mallSettingStore{cache: &Cache{Key: merchantApplyConfigKey, Result: `{"background_image":"x","unknown":1}`}}
	svc := NewService(store)
	safe, err := svc.GetMerchantApplyConfig(context.Background())
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if safe == store.cache.Result || safe != marshalMerchantApplyConfig(defaultMerchantApplyConfig()) {
		t.Fatalf("legacy malformed must not be returned, got %q", safe)
	}
}

func TestMerchantApplyConfigRejectsBadFieldType(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveMerchantApplyConfig(context.Background(), `{"background_image":"","form_fields":[{"id":"f1","type":"hack","title":"x","content_type":"text","default_value":"","placeholder":"","required":false}]}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("bad type error = %v, want ErrBadParam", err)
	}
}

func TestMerchantApplyConfigAllowsTextarea(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	raw := `{"background_image":"","form_fields":[{"id":"ta1","type":"textarea","title":"备注说明","content_type":"text","default_value":"","placeholder":"请填写","required":false}]}`
	got, err := svc.SaveMerchantApplyConfig(context.Background(), raw)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	want := `{"background_image":"","form_fields":[{"id":"ta1","type":"textarea","title":"备注说明","content_type":"text","default_value":"","placeholder":"请填写","required":false}]}`
	if got != want {
		t.Fatalf("canonical = %q\nwant       = %q", got, want)
	}
}

func TestMerchantApplyConfigCanonicalizesExtendedProps(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	raw := `{"background_image":"","form_fields":[{"id":"img1","type":"image","title":"门头照","content_type":"text","default_value":"","placeholder":"","required":true,"max_upload":3},{"id":"city1","type":"city","title":"所在城市","content_type":"text","default_value":"","placeholder":"请选择","required":false,"city_level":"province_city"},{"id":"tr1","type":"timerange","title":"营业时间","content_type":"text","default_value":"","placeholder":"请选择","required":false,"default_visible":"show","default_mode":"specify","specify_value":"09:00 ~ 18:00"}]}`
	got, err := svc.SaveMerchantApplyConfig(context.Background(), raw)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	want := `{"background_image":"","form_fields":[{"id":"img1","type":"image","title":"门头照","content_type":"text","default_value":"","placeholder":"","required":true,"max_upload":3},{"id":"city1","type":"city","title":"所在城市","content_type":"text","default_value":"","placeholder":"请选择","required":false,"city_level":"province_city"},{"id":"tr1","type":"timerange","title":"营业时间","content_type":"text","default_value":"","placeholder":"请选择","required":false,"default_visible":"show","default_mode":"specify","specify_value":"09:00 ~ 18:00"}]}`
	if got != want {
		t.Fatalf("canonical = %q\nwant       = %q", got, want)
	}
}
