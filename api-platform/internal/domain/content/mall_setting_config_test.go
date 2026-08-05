package content

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

type mallSettingStore struct{ cache *Cache }

func (s *mallSettingStore) ListNotices(context.Context, bool, int, int) ([]Notice, int64, error) {
	return nil, 0, nil
}
func (s *mallSettingStore) GetNotice(context.Context, uint) (*Notice, error) {
	return nil, gorm.ErrRecordNotFound
}
func (s *mallSettingStore) CreateNotice(context.Context, *Notice) error  { return nil }
func (s *mallSettingStore) UpdateNotice(context.Context, *Notice) error  { return nil }
func (s *mallSettingStore) SoftDeleteNotice(context.Context, uint) error { return nil }
func (s *mallSettingStore) GetCache(_ context.Context, key string) (*Cache, error) {
	if s.cache == nil || s.cache.Key != key {
		return nil, gorm.ErrRecordNotFound
	}
	return s.cache, nil
}
func (s *mallSettingStore) UpsertCache(_ context.Context, row *Cache) error {
	s.cache = row
	return nil
}

func TestShopConfigRejectsUnknownKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveShopConfig(context.Background(), `{"site_name":"七禧商城","secret_key":"not-allowed"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("unknown shop key error = %v, want ErrBadParam", err)
	}
}

func TestShopConfigCanonicalizesChineseValues(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	got, err := svc.SaveShopConfig(context.Background(), `{"site_name":"七禧商城","site_url":"https://example.test","order_auto_cancel_minutes":15,"order_auto_receive_days":10,"enabled":true,"remark":"中文商城设置"}`)
	if err != nil {
		t.Fatalf("save shop config: %v", err)
	}
	want := `{"site_name":"七禧商城","site_url":"https://example.test","order_auto_cancel_minutes":15,"order_auto_receive_days":10,"enabled":true,"remark":"中文商城设置"}`
	if got != want || store.cache == nil || store.cache.Result != want {
		t.Fatalf("canonical shop config = %q, stored = %#v", got, store.cache)
	}
}

func TestShopConfigNeverEchoesLegacyMalformed(t *testing.T) {
	store := &mallSettingStore{cache: &Cache{Key: shopConfigKey, Result: `{"site_name":"旧值","unknown":true}`}}
	svc := NewService(store)
	safe, err := svc.GetShopConfig(context.Background())
	if err != nil {
		t.Fatalf("get shop config: %v", err)
	}
	if safe == store.cache.Result || safe != marshalShopConfig(defaultShopConfig()) {
		t.Fatalf("legacy malformed shop config must not be returned, got %q", safe)
	}
}

func TestPayConfigRejectsSensitiveKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	cases := []string{
		`{"wechat_enabled":true,"api_key":"not-allowed"}`,
		`{"alipay_enabled":true,"merchant_cert":"not-allowed"}`,
		`{"balance_enabled":true,"access_token":"not-allowed"}`,
		`{"wechat_enabled":true,"pay_password":"not-allowed"}`,
		`{"wechat_enabled":true,"app_secret":"not-allowed"}`,
	}
	for _, raw := range cases {
		if _, err := svc.SavePayConfig(context.Background(), raw); !errors.Is(err, ErrBadParam) {
			t.Fatalf("sensitive pay config %q error = %v, want ErrBadParam", raw, err)
		}
	}
}

func TestPayConfigCanonicalizesAndNeverEchoesLegacySecret(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	got, err := svc.SavePayConfig(context.Background(), `{"wechat_enabled":true,"alipay_enabled":false,"balance_enabled":true,"remark":"中文支付开关"}`)
	if err != nil {
		t.Fatalf("save pay config: %v", err)
	}
	want := `{"wechat_enabled":true,"alipay_enabled":false,"balance_enabled":true,"remark":"中文支付开关"}`
	if got != want || store.cache == nil || store.cache.Result != want {
		t.Fatalf("canonical pay config = %q, stored = %#v", got, store.cache)
	}
	store.cache = &Cache{Key: payConfigKey, Result: `{"wechat_enabled":true,"mch_secret":"must-not-leak"}`}
	safe, err := svc.GetPayConfig(context.Background())
	if err != nil {
		t.Fatalf("get pay config: %v", err)
	}
	if safe == store.cache.Result || safe != marshalPayConfig(defaultPayConfig()) {
		t.Fatalf("legacy secret pay config must not be returned, got %q", safe)
	}
}

func TestWechatAppConfigRejectsSensitiveKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveWechatAppConfig(context.Background(), `{"app_name":"七禧商城","enabled":true,"app_secret":"not-allowed"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("sensitive wechat app key error = %v, want ErrBadParam", err)
	}
}

func TestWechatAppConfigCanonicalizesChineseValues(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	got, err := svc.SaveWechatAppConfig(context.Background(), `{"app_name":"七禧商城公众号","enabled":true,"remark":"中文公众号开关"}`)
	if err != nil {
		t.Fatalf("save wechat app config: %v", err)
	}
	want := `{"app_name":"七禧商城公众号","enabled":true,"remark":"中文公众号开关"}`
	if got != want || store.cache == nil || store.cache.Result != want {
		t.Fatalf("canonical wechat app config = %q, stored = %#v", got, store.cache)
	}
}
