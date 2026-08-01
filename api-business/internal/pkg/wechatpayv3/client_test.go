package wechatpayv3

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNativePrepaySignsAndReturnsCodeURL(t *testing.T) {
	privateKey, privatePEM, _ := testKeys(t)
	now := time.Date(2026, 8, 1, 8, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/pay/transactions/native" || r.Method != http.MethodPost {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "WECHATPAY2-SHA256-RSA2048 ") {
			t.Fatal("missing wechat authorization")
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["out_trade_no"] != "G202608010001" || body["notify_url"] != "https://callback.example/pay/wechat" {
			t.Fatalf("unexpected body %#v", body)
		}
		if privateKey.PublicKey.N.BitLen() < 2048 {
			t.Fatal("weak test key")
		}
		_, _ = w.Write([]byte(`{"code_url":"weixin://wxpay/bizpayurl?pr=test"}`))
	}))
	defer server.Close()
	client := &Client{BaseURL: server.URL, Now: func() time.Time { return now }, Nonce: func() (string, error) { return "nonce-for-test", nil }}
	result, err := client.NativePrepay(t.Context(), Config{AppID: "wx-test", MchID: "1900000001", MerchantSerialNo: "merchant-serial", MerchantPrivateKey: privatePEM, NotifyURL: "https://callback.example/pay/wechat"}, NativeRequest{Description: "七禧商城订单", OutTradeNo: "G202608010001", AmountCents: 199, ExpireAt: now.Add(15 * time.Minute)})
	if err != nil {
		t.Fatal(err)
	}
	if result.CodeURL != "weixin://wxpay/bizpayurl?pr=test" {
		t.Fatalf("code url=%q", result.CodeURL)
	}
}

func TestVerifyAndDecryptCallback(t *testing.T) {
	privateKey, _, publicPEM := testKeys(t)
	now := time.Date(2026, 8, 1, 8, 0, 0, 0, time.UTC)
	apiV3Key := "0123456789abcdef0123456789abcdef"
	plain := []byte(`{"out_trade_no":"G202608010001","transaction_id":"42000000000001","mchid":"1900000001","appid":"wx-test","trade_state":"SUCCESS","success_time":"2026-08-01T08:00:00+00:00","amount":{"total":199}}`)
	resource := encryptedResource(t, apiV3Key, plain)
	body, err := json.Marshal(map[string]any{"id": "event-1", "event_type": "TRANSACTION.SUCCESS", "resource": resource})
	if err != nil {
		t.Fatal(err)
	}
	timestamp := "1785571200"
	nonce := "callback-nce"
	sig, err := sign(privateKey, timestamp+"\n"+nonce+"\n"+string(body)+"\n")
	if err != nil {
		t.Fatal(err)
	}
	headers := http.Header{"Wechatpay-Timestamp": []string{timestamp}, "Wechatpay-Nonce": []string{nonce}, "Wechatpay-Signature": []string{sig}, "Wechatpay-Serial": []string{"pub-key-id"}}
	result, err := VerifyAndDecryptCallback(Config{APIv3Key: apiV3Key, PublicKeyID: "pub-key-id", PublicKeyPEM: publicPEM}, headers, body, now)
	if err != nil {
		t.Fatal(err)
	}
	if result.OutTradeNo != "G202608010001" || result.TransactionID != "42000000000001" || result.AmountCents != 199 {
		t.Fatalf("unexpected transaction %#v", result)
	}
	if _, err := VerifyAndDecryptCallback(Config{APIv3Key: apiV3Key, PublicKeyID: "wrong", PublicKeyPEM: publicPEM}, headers, body, now); err == nil {
		t.Fatal("expected serial mismatch")
	}
}

func testKeys(t *testing.T) (*rsa.PrivateKey, string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	privateDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	publicDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	return key, string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateDER})), string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER}))
}

func encryptedResource(t *testing.T, apiV3Key string, plain []byte) map[string]string {
	t.Helper()
	block, err := aes.NewCipher([]byte(apiV3Key))
	if err != nil {
		t.Fatal(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatal(err)
	}
	nonce := "123456789012"
	associatedData := "transaction"
	return map[string]string{"algorithm": "AEAD_AES_256_GCM", "associated_data": associatedData, "nonce": nonce, "ciphertext": base64.StdEncoding.EncodeToString(gcm.Seal(nil, []byte(nonce), plain, []byte(associatedData)))}
}
