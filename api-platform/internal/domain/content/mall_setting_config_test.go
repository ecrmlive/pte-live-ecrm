package content

import (
	"context"
	"errors"
	"strings"
	"testing"

	"gorm.io/gorm"
)

type mallSettingStore struct{ cache *Cache }

func (s *mallSettingStore) ListNotices(context.Context, bool, NoticeListFilter) ([]Notice, int64, error) {
	return nil, 0, nil
}
func (s *mallSettingStore) GetNotice(context.Context, uint) (*Notice, error) {
	return nil, gorm.ErrRecordNotFound
}
func (s *mallSettingStore) CreateNotice(context.Context, *Notice, []uint) error { return nil }
func (s *mallSettingStore) UpdateNotice(context.Context, *Notice, []uint) error { return nil }
func (s *mallSettingStore) ListNoticeScopes(context.Context, []uint) ([]NoticeScope, error) {
	return nil, nil
}
func (s *mallSettingStore) UpdateNoticeStatus(context.Context, uint, int8) error { return nil }
func (s *mallSettingStore) SoftDeleteNotice(context.Context, uint) error         { return nil }
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
	got, err := svc.SaveShopConfig(context.Background(), `{"auto_parse_clipboard":true,"arrival_notice_enabled":true,"product_comment_enabled":true,"auto_positive_review_enabled":true,"default_copy_times":8,"order_auto_cancel_minutes":15,"order_auto_receive_days":10,"after_sale_days":1,"merchant_refund_auto_days":1,"refund_reasons":["测试","不要了"],"platform_rights_enabled":true,"platform_rights_days":1,"merge_payment_enabled":true,"merchant_apply_enabled":true,"merchant_qualification_required":true,"merchant_margin_badge_enabled":false,"merchant_margin_badge_image":"/attachment/badge.png","merchant_category_limit":5,"mall_show_stores":true,"mall_recommend_enabled":true,"mall_recommend_distance_enabled":true,"mall_recommend_sort":"star","live_stream_auto_approve":false,"live_product_auto_approve":false,"hot_ranking_enabled":true,"hot_ranking_category_level":2,"hot_ranking_refresh_hours":24,"mall_search_mode":"fuzzy","product_ranking_period":"month","product_ranking_metric":"sales_amount","shop_ranking_period":"month","shop_ranking_metric":"product_count","dashboard_display_name":"数据大屏"}`)
	if err != nil {
		t.Fatalf("save shop config: %v", err)
	}
	if !strings.Contains(got, `"dashboard_display_name":"数据大屏"`) || !strings.Contains(got, `"refund_reasons":["测试","不要了"]`) || store.cache == nil || store.cache.Result != got {
		t.Fatalf("canonical shop config = %q, stored = %#v", got, store.cache)
	}
}

func TestShopConfigRejectsInvalidEnumsAndSensitiveKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	for _, raw := range []string{
		`{"mall_search_mode":"anything"}`,
		`{"dashboard_display_name":"数据大屏","secret_key":"not-allowed"}`,
	} {
		if _, err := svc.SaveShopConfig(context.Background(), raw); !errors.Is(err, ErrBadParam) {
			t.Fatalf("shop config %q error = %v, want ErrBadParam", raw, err)
		}
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

func TestMarginConfigRejectsUnknownKeys(t *testing.T) {
	svc := NewService(&mallSettingStore{})
	_, err := svc.SaveMarginConfig(context.Background(), `{"margin_remind_switch":true,"secret_key":"not-allowed"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("unknown margin key error = %v, want ErrBadParam", err)
	}
}

func TestMarginConfigCanonicalizesAndDefaults(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	got, err := svc.SaveMarginConfig(context.Background(), `{"margin_remind_switch":true,"margin_remind_day":30}`)
	if err != nil {
		t.Fatalf("save margin config: %v", err)
	}
	want := `{"margin_remind_switch":true,"margin_remind_day":30}`
	if got != want || store.cache == nil || store.cache.Result != want {
		t.Fatalf("canonical margin config = %q, stored = %#v", got, store.cache)
	}
	store.cache = &Cache{Key: marginConfigKey, Result: `{"margin_remind_switch":true,"unknown":1}`}
	safe, err := svc.GetMarginConfig(context.Background())
	if err != nil {
		t.Fatalf("get margin config: %v", err)
	}
	if safe == store.cache.Result || safe != marshalMarginConfig(defaultMarginConfig()) {
		t.Fatalf("legacy malformed margin config must not be returned, got %q", safe)
	}
}
