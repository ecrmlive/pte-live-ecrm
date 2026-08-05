package wechatpayv3

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestH5PrepaySendsH5SceneAndReturnsMWEBURL(t *testing.T) {
	_, privatePEM, _ := testKeys(t)
	now := time.Date(2026, 8, 1, 8, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/pay/transactions/h5" || r.Method != http.MethodPost {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		var body struct {
			AppID     string `json:"appid"`
			SceneInfo struct {
				PayerClientIP string `json:"payer_client_ip"`
				H5Info        struct {
					Type   string `json:"type"`
					AppURL string `json:"app_url"`
				} `json:"h5_info"`
			} `json:"scene_info"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.AppID != "wx-h5-test" || body.SceneInfo.PayerClientIP != "203.0.113.10" || body.SceneInfo.H5Info.Type != "Wap" || body.SceneInfo.H5Info.AppURL != "https://mall.example.test" {
			t.Fatalf("unexpected h5 request %#v", body)
		}
		_, _ = w.Write([]byte(`{"h5_url":"https://wx.tenpay.com/cgi-bin/mmpayweb-bin/checkmweb?prepay_id=test"}`))
	}))
	defer server.Close()
	client := &Client{BaseURL: server.URL, Now: func() time.Time { return now }, Nonce: func() (string, error) { return "nonce-for-test", nil }}
	result, err := client.H5Prepay(t.Context(), Config{AppID: "wx-mini-test", H5AppID: "wx-h5-test", H5SiteURL: "https://mall.example.test", MchID: "1900000001", MerchantSerialNo: "merchant-serial", MerchantPrivateKey: privatePEM, NotifyURL: "https://callback.example/pay/wechat"}, H5Request{Description: "CRM Live 订单", OutTradeNo: "G202608010003", AmountCents: 199, PayerClientIP: "203.0.113.10", ExpireAt: now.Add(15 * time.Minute)})
	if err != nil {
		t.Fatal(err)
	}
	if result.MWEBURL != "https://wx.tenpay.com/cgi-bin/mmpayweb-bin/checkmweb?prepay_id=test" {
		t.Fatalf("mweb url=%q", result.MWEBURL)
	}
}
