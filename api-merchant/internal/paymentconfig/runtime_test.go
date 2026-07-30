package paymentconfig

import "testing"

func TestChannelReady(t *testing.T) {
	wechat := Values{"wechat_enabled": "true", "wechat_app_id": "wxid", "wechat_mch_id": "mch", "wechat_api_v3_key": "key", "wechat_private_key": "private", "wechat_notify_url": "https://example.test/wechat"}
	if !ChannelReady(wechat, "wechat") {
		t.Fatal("expected complete wechat config to be ready")
	}
	delete(wechat, "wechat_private_key")
	if ChannelReady(wechat, "wechat") {
		t.Fatal("expected incomplete wechat config to be unavailable")
	}
	alipay := Values{"alipay_enabled": "true", "alipay_app_id": "app", "alipay_private_key": "private", "alipay_public_key": "public", "alipay_notify_url": "https://example.test/alipay"}
	if !ChannelReady(alipay, "alipay") {
		t.Fatal("expected complete alipay config to be ready")
	}
	merchantWechat := Values{"enabled": "true", "app_id": "wxid", "mch_id": "mch", "api_v3_key": "key", "serial_no": "serial", "private_key": "private", "notify_url": "https://merchant.example.test/wechat"}
	if !StoreChannelReady(merchantWechat, "wechat") {
		t.Fatal("expected complete merchant wechat config to be ready")
	}
	merchantWechat["enabled"] = "false"
	if StoreChannelReady(merchantWechat, "wechat") {
		t.Fatal("expected disabled merchant wechat config to be unavailable")
	}
}
