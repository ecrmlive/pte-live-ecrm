package wechatpayv3

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRefundSignsServerSideRequest(t *testing.T) {
	_, privatePEM, _ := testKeys(t)
	now := time.Date(2026, 8, 3, 8, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v3/refund/domestic/refunds" {
			t.Fatalf("request=%s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") == "" {
			t.Fatal("missing WeChat authorization")
		}
		var body struct {
			OutTradeNo  string                        `json:"out_trade_no"`
			OutRefundNo string                        `json:"out_refund_no"`
			NotifyURL   string                        `json:"notify_url"`
			Amount      struct{ Refund, Total int64 } `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.OutTradeNo != "G202608030001" || body.OutRefundNo != "R202608030001" || body.Amount.Refund != 199 || body.Amount.Total != 299 || body.NotifyURL != "https://callback.example/refund/wechat" {
			t.Fatalf("unexpected body %#v", body)
		}
		_, _ = w.Write([]byte(`{"out_refund_no":"R202608030001","refund_id":"5000000001","status":"PROCESSING"}`))
	}))
	defer server.Close()
	client := &Client{BaseURL: server.URL, Now: func() time.Time { return now }, Nonce: func() (string, error) { return "refund-nonce", nil }}
	result, err := client.Refund(t.Context(), Config{MchID: "1900000001", MerchantSerialNo: "merchant-serial", MerchantPrivateKey: privatePEM, RefundNotifyURL: "https://callback.example/refund/wechat"}, RefundRequest{OutTradeNo: "G202608030001", OutRefundNo: "R202608030001", Reason: "七禧演示退款", TotalCents: 299, RefundCents: 199})
	if err != nil {
		t.Fatal(err)
	}
	if result.ProviderRefundNo != "5000000001" || result.Status != "PROCESSING" {
		t.Fatalf("unexpected result %#v", result)
	}
}

func TestVerifyAndDecryptRefundCallback(t *testing.T) {
	privateKey, _, publicPEM := testKeys(t)
	now := time.Date(2026, 8, 3, 8, 0, 0, 0, time.UTC)
	apiV3Key := "0123456789abcdef0123456789abcdef"
	plain := []byte(`{"out_trade_no":"G202608030001","out_refund_no":"R202608030001","transaction_id":"42000000000001","refund_id":"5000000001","mchid":"1900000001","refund_status":"SUCCESS","success_time":"2026-08-03T08:00:00+00:00","amount":{"refund":199,"total":299,"currency":"CNY"}}`)
	resource := encryptedResource(t, apiV3Key, plain)
	body, err := json.Marshal(map[string]any{"id": "refund-event-1", "event_type": "REFUND.SUCCESS", "resource": resource})
	if err != nil {
		t.Fatal(err)
	}
	timestamp := "1785744000"
	nonce := "refund-callback-nonce"
	sig, err := sign(privateKey, timestamp+"\n"+nonce+"\n"+string(body)+"\n")
	if err != nil {
		t.Fatal(err)
	}
	headers := http.Header{"Wechatpay-Timestamp": []string{timestamp}, "Wechatpay-Nonce": []string{nonce}, "Wechatpay-Signature": []string{sig}, "Wechatpay-Serial": []string{"pub-key-id"}}
	result, err := VerifyAndDecryptRefundCallback(Config{APIv3Key: apiV3Key, PublicKeyID: "pub-key-id", PublicKeyPEM: publicPEM}, headers, body, now)
	if err != nil {
		t.Fatal(err)
	}
	if result.OutRefundNo != "R202608030001" || result.RefundCents != 199 || result.TotalCents != 299 {
		t.Fatalf("unexpected refund %#v", result)
	}
	wrongBody := append([]byte(nil), body...)
	wrongBody[len(wrongBody)-2] ^= 1
	if _, err := VerifyAndDecryptRefundCallback(Config{APIv3Key: apiV3Key, PublicKeyID: "pub-key-id", PublicKeyPEM: publicPEM}, headers, wrongBody, now); err == nil {
		t.Fatal("expected modified callback rejection")
	}
}
