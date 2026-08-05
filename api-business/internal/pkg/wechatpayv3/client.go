// Package wechatpayv3 implements the small, auditable subset of the WeChat
// Pay v3 protocol that the business service needs.  It deliberately keeps
// merchant credentials server-side and never accepts an amount from a client.
package wechatpayv3

import (
	"bytes"
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DefaultBaseURL = "https://api.mch.weixin.qq.com"

var (
	ErrInvalidConfig    = errors.New("微信支付配置不完整")
	ErrInvalidSignature = errors.New("微信支付回调签名无效")
	ErrInvalidCallback  = errors.New("微信支付回调内容无效")
)

type Config struct {
	AppID              string
	H5AppID            string
	H5SiteURL          string
	MchID              string
	MerchantSerialNo   string
	MerchantPrivateKey string
	MerchantCertPEM    string
	APIv3Key           string
	NotifyURL          string
	// RefundNotifyURL is deliberately separate from payment notifications so a
	// refund callback cannot be parsed as a payment-success notification.
	RefundNotifyURL string
	// PublicKeyID/PublicKeyPEM are the platform callback verification key.
	// The fields also support the platform-certificate PEM/serial fallback.
	PublicKeyID        string
	PublicKeyPEM       string
	PlatformCertSerial string
	PlatformCertPEM    string
}

func (c Config) ValidForNative() bool {
	return strings.TrimSpace(c.AppID) != "" && strings.TrimSpace(c.MchID) != "" &&
		c.resolvedMerchantSerialNo() != "" && strings.TrimSpace(c.MerchantPrivateKey) != "" &&
		strings.TrimSpace(c.NotifyURL) != ""
}

// ValidForJSAPI has the same merchant-signing requirements as Native. The
// caller must additionally provide the authenticated mini-program openid.
func (c Config) ValidForJSAPI() bool { return c.ValidForNative() }

func (c Config) ValidForH5() bool {
	return c.ValidForNative() && strings.TrimSpace(c.H5AppID) != "" && strings.HasPrefix(strings.TrimSpace(c.H5SiteURL), "https://")
}

func (c Config) ValidForRefund() bool {
	return strings.TrimSpace(c.MchID) != "" && c.resolvedMerchantSerialNo() != "" &&
		strings.TrimSpace(c.MerchantPrivateKey) != "" && strings.HasPrefix(strings.TrimSpace(c.RefundNotifyURL), "https://")
}

func (c Config) ValidForCallback() bool {
	return len([]byte(c.APIv3Key)) == 32 && c.callbackPublicKeyPEM() != ""
}

type NativeRequest struct {
	Description string
	OutTradeNo  string
	AmountCents int64
	Attach      string
	ExpireAt    time.Time
}

type NativeResponse struct {
	CodeURL string
}

type JSAPIRequest struct {
	Description string
	OutTradeNo  string
	AmountCents int64
	OpenID      string
	Attach      string
	ExpireAt    time.Time
}

// JSAPIResponse is passed to uni.requestPayment without exposing any merchant
// credential. Field names intentionally match the Mini Program API.
type JSAPIResponse struct {
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Now        func() time.Time
	Nonce      func() (string, error)
}

func (c *Client) NativePrepay(ctx context.Context, config Config, request NativeRequest) (NativeResponse, error) {
	if !config.ValidForNative() || request.AmountCents <= 0 || strings.TrimSpace(request.OutTradeNo) == "" {
		return NativeResponse{}, ErrInvalidConfig
	}
	privateKey, err := parseRSAPrivateKey(config.MerchantPrivateKey)
	if err != nil {
		return NativeResponse{}, ErrInvalidConfig
	}
	now := c.now()
	expires := request.ExpireAt.UTC()
	if expires.IsZero() || !expires.After(now) {
		expires = now.Add(15 * time.Minute)
	}
	body, err := json.Marshal(struct {
		AppID       string `json:"appid"`
		MchID       string `json:"mchid"`
		Description string `json:"description"`
		OutTradeNo  string `json:"out_trade_no"`
		Attach      string `json:"attach,omitempty"`
		NotifyURL   string `json:"notify_url"`
		TimeExpire  string `json:"time_expire"`
		Amount      struct {
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		} `json:"amount"`
	}{
		AppID:       strings.TrimSpace(config.AppID),
		MchID:       strings.TrimSpace(config.MchID),
		Description: truncate(strings.TrimSpace(request.Description), 127),
		OutTradeNo:  strings.TrimSpace(request.OutTradeNo),
		Attach:      truncate(strings.TrimSpace(request.Attach), 127),
		NotifyURL:   strings.TrimSpace(config.NotifyURL),
		TimeExpire:  expires.Format(time.RFC3339),
		Amount: struct {
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		}{Total: request.AmountCents, Currency: "CNY"},
	})
	if err != nil {
		return NativeResponse{}, err
	}
	const path = "/v3/pay/transactions/native"
	nonce, err := c.nonce()
	if err != nil {
		return NativeResponse{}, err
	}
	timestamp := fmt.Sprintf("%d", now.Unix())
	signature, err := sign(privateKey, strings.Join([]string{http.MethodPost, path, timestamp, nonce, string(body), ""}, "\n"))
	if err != nil {
		return NativeResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL()+path, bytes.NewReader(body))
	if err != nil {
		return NativeResponse{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",timestamp="%s",serial_no="%s",signature="%s"`, config.MchID, nonce, timestamp, config.resolvedMerchantSerialNo(), signature))

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return NativeResponse{}, fmt.Errorf("微信支付下单请求失败: %w", err)
	}
	defer resp.Body.Close()
	limited := io.LimitReader(resp.Body, 1<<20)
	data, err := io.ReadAll(limited)
	if err != nil {
		return NativeResponse{}, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var failure struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		_ = json.Unmarshal(data, &failure)
		return NativeResponse{}, fmt.Errorf("微信支付下单失败: %s", providerMessage(failure.Code, failure.Message, resp.StatusCode))
	}
	var output struct {
		CodeURL string `json:"code_url"`
	}
	if err := json.Unmarshal(data, &output); err != nil || strings.TrimSpace(output.CodeURL) == "" {
		return NativeResponse{}, ErrInvalidCallback
	}
	return NativeResponse{CodeURL: output.CodeURL}, nil
}

// JSAPIPrepay creates a Mini Program payment and signs the client invocation
// parameters. openid is sourced from a verified mini-program login, never from
// a request body.
func (c *Client) JSAPIPrepay(ctx context.Context, config Config, request JSAPIRequest) (JSAPIResponse, error) {
	if !config.ValidForJSAPI() || request.AmountCents <= 0 || strings.TrimSpace(request.OutTradeNo) == "" || strings.TrimSpace(request.OpenID) == "" {
		return JSAPIResponse{}, ErrInvalidConfig
	}
	privateKey, err := parseRSAPrivateKey(config.MerchantPrivateKey)
	if err != nil {
		return JSAPIResponse{}, ErrInvalidConfig
	}
	now := c.now()
	expires := request.ExpireAt.UTC()
	if expires.IsZero() || !expires.After(now) {
		expires = now.Add(15 * time.Minute)
	}
	body, err := json.Marshal(struct {
		AppID       string `json:"appid"`
		MchID       string `json:"mchid"`
		Description string `json:"description"`
		OutTradeNo  string `json:"out_trade_no"`
		Attach      string `json:"attach,omitempty"`
		NotifyURL   string `json:"notify_url"`
		TimeExpire  string `json:"time_expire"`
		Amount      struct {
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		} `json:"amount"`
		Payer struct {
			OpenID string `json:"openid"`
		} `json:"payer"`
	}{
		AppID:       strings.TrimSpace(config.AppID),
		MchID:       strings.TrimSpace(config.MchID),
		Description: truncate(strings.TrimSpace(request.Description), 127),
		OutTradeNo:  strings.TrimSpace(request.OutTradeNo),
		Attach:      truncate(strings.TrimSpace(request.Attach), 127),
		NotifyURL:   strings.TrimSpace(config.NotifyURL),
		TimeExpire:  expires.Format(time.RFC3339),
		Amount: struct {
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		}{Total: request.AmountCents, Currency: "CNY"},
		Payer: struct {
			OpenID string `json:"openid"`
		}{OpenID: strings.TrimSpace(request.OpenID)},
	})
	if err != nil {
		return JSAPIResponse{}, err
	}
	const path = "/v3/pay/transactions/jsapi"
	nonce, err := c.nonce()
	if err != nil {
		return JSAPIResponse{}, err
	}
	timestamp := fmt.Sprintf("%d", now.Unix())
	signature, err := sign(privateKey, strings.Join([]string{http.MethodPost, path, timestamp, nonce, string(body), ""}, "\n"))
	if err != nil {
		return JSAPIResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL()+path, bytes.NewReader(body))
	if err != nil {
		return JSAPIResponse{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",timestamp="%s",serial_no="%s",signature="%s"`, config.MchID, nonce, timestamp, config.resolvedMerchantSerialNo(), signature))

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return JSAPIResponse{}, fmt.Errorf("微信支付下单请求失败: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return JSAPIResponse{}, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var failure struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		_ = json.Unmarshal(data, &failure)
		return JSAPIResponse{}, fmt.Errorf("微信支付下单失败: %s", providerMessage(failure.Code, failure.Message, resp.StatusCode))
	}
	var output struct {
		PrepayID string `json:"prepay_id"`
	}
	if err := json.Unmarshal(data, &output); err != nil || strings.TrimSpace(output.PrepayID) == "" {
		return JSAPIResponse{}, ErrInvalidCallback
	}
	payNonce, err := c.nonce()
	if err != nil {
		return JSAPIResponse{}, err
	}
	payTimestamp := fmt.Sprintf("%d", c.now().Unix())
	payPackage := "prepay_id=" + strings.TrimSpace(output.PrepayID)
	paySign, err := sign(privateKey, strings.Join([]string{strings.TrimSpace(config.AppID), payTimestamp, payNonce, payPackage, ""}, "\n"))
	if err != nil {
		return JSAPIResponse{}, err
	}
	return JSAPIResponse{AppID: strings.TrimSpace(config.AppID), TimeStamp: payTimestamp, NonceStr: payNonce, Package: payPackage, SignType: "RSA", PaySign: paySign}, nil
}

type CallbackTransaction struct {
	EventID       string
	OutTradeNo    string
	TransactionID string
	MchID         string
	AppID         string
	TradeState    string
	AmountCents   int64
	SuccessTime   time.Time
	RawPayload    json.RawMessage
}

// VerifyAndDecryptCallback verifies the WeChat signature before attempting to
// decrypt its AES-GCM resource. The caller must additionally match the result
// to its own payment transaction and order amount.
func VerifyAndDecryptCallback(config Config, header http.Header, body []byte, now time.Time) (CallbackTransaction, error) {
	if !config.ValidForCallback() || len(body) == 0 || len(body) > 1<<20 {
		return CallbackTransaction{}, ErrInvalidConfig
	}
	timestamp := strings.TrimSpace(header.Get("Wechatpay-Timestamp"))
	nonce := strings.TrimSpace(header.Get("Wechatpay-Nonce"))
	signature := strings.TrimSpace(header.Get("Wechatpay-Signature"))
	serial := strings.TrimSpace(header.Get("Wechatpay-Serial"))
	if timestamp == "" || nonce == "" || signature == "" || serial == "" {
		return CallbackTransaction{}, ErrInvalidSignature
	}
	if expected := strings.TrimSpace(config.PublicKeyID); expected != "" && serial != expected {
		return CallbackTransaction{}, ErrInvalidSignature
	}
	if expected := strings.TrimSpace(config.PlatformCertSerial); expected != "" && config.PublicKeyID == "" && serial != expected {
		return CallbackTransaction{}, ErrInvalidSignature
	}
	seconds, err := parseUnix(timestamp)
	if err != nil || now.UTC().Sub(time.Unix(seconds, 0).UTC()) > 5*time.Minute || time.Unix(seconds, 0).UTC().Sub(now.UTC()) > 5*time.Minute {
		return CallbackTransaction{}, ErrInvalidSignature
	}
	if err := VerifyCallbackSignature(config, header, body, now); err != nil {
		return CallbackTransaction{}, ErrInvalidSignature
	}
	var envelope struct {
		ID        string `json:"id"`
		EventType string `json:"event_type"`
		Resource  struct {
			Algorithm      string `json:"algorithm"`
			Ciphertext     string `json:"ciphertext"`
			AssociatedData string `json:"associated_data"`
			Nonce          string `json:"nonce"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil || envelope.EventType != "TRANSACTION.SUCCESS" || envelope.ID == "" || envelope.Resource.Algorithm != "AEAD_AES_256_GCM" {
		return CallbackTransaction{}, ErrInvalidCallback
	}
	plain, err := decryptResource(config.APIv3Key, envelope.Resource.AssociatedData, envelope.Resource.Nonce, envelope.Resource.Ciphertext)
	if err != nil {
		return CallbackTransaction{}, ErrInvalidCallback
	}
	var transaction struct {
		OutTradeNo    string `json:"out_trade_no"`
		TransactionID string `json:"transaction_id"`
		MchID         string `json:"mchid"`
		AppID         string `json:"appid"`
		TradeState    string `json:"trade_state"`
		SuccessTime   string `json:"success_time"`
		Amount        struct {
			Total int64 `json:"total"`
		} `json:"amount"`
	}
	if err := json.Unmarshal(plain, &transaction); err != nil || transaction.OutTradeNo == "" || transaction.TransactionID == "" || transaction.Amount.Total <= 0 {
		return CallbackTransaction{}, ErrInvalidCallback
	}
	paidAt, err := time.Parse(time.RFC3339, transaction.SuccessTime)
	if err != nil {
		return CallbackTransaction{}, ErrInvalidCallback
	}
	return CallbackTransaction{EventID: envelope.ID, OutTradeNo: transaction.OutTradeNo, TransactionID: transaction.TransactionID, MchID: transaction.MchID, AppID: transaction.AppID, TradeState: transaction.TradeState, AmountCents: transaction.Amount.Total, SuccessTime: paidAt.UTC(), RawPayload: append(json.RawMessage(nil), plain...)}, nil
}

// VerifyCallbackSignature verifies only the signed callback envelope. It is
// separated from decryption because a multi-merchant callback must select the
// APIv3 key after the common WeChat platform signature has been authenticated.
func VerifyCallbackSignature(config Config, header http.Header, body []byte, now time.Time) error {
	if config.callbackPublicKeyPEM() == "" || len(body) == 0 || len(body) > 1<<20 {
		return ErrInvalidConfig
	}
	timestamp := strings.TrimSpace(header.Get("Wechatpay-Timestamp"))
	nonce := strings.TrimSpace(header.Get("Wechatpay-Nonce"))
	signature := strings.TrimSpace(header.Get("Wechatpay-Signature"))
	serial := strings.TrimSpace(header.Get("Wechatpay-Serial"))
	if timestamp == "" || nonce == "" || signature == "" || serial == "" {
		return ErrInvalidSignature
	}
	if expected := strings.TrimSpace(config.PublicKeyID); expected != "" && serial != expected {
		return ErrInvalidSignature
	}
	if expected := strings.TrimSpace(config.PlatformCertSerial); expected != "" && config.PublicKeyID == "" && serial != expected {
		return ErrInvalidSignature
	}
	seconds, err := parseUnix(timestamp)
	if err != nil || now.UTC().Sub(time.Unix(seconds, 0).UTC()) > 5*time.Minute || time.Unix(seconds, 0).UTC().Sub(now.UTC()) > 5*time.Minute {
		return ErrInvalidSignature
	}
	key, err := parseRSAPublicKey(config.callbackPublicKeyPEM())
	if err != nil {
		return ErrInvalidConfig
	}
	signed := timestamp + "\n" + nonce + "\n" + string(body) + "\n"
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil || rsa.VerifyPKCS1v15(key, crypto.SHA256, hash([]byte(signed)), signatureBytes) != nil {
		return ErrInvalidSignature
	}
	return nil
}

func (c *Client) baseURL() string {
	if base := strings.TrimRight(strings.TrimSpace(c.BaseURL), "/"); base != "" {
		return base
	}
	return DefaultBaseURL
}
func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: 15 * time.Second}
}
func (c *Client) now() time.Time {
	if c.Now != nil {
		return c.Now().UTC()
	}
	return time.Now().UTC()
}
func (c *Client) nonce() (string, error) {
	if c.Nonce != nil {
		return c.Nonce()
	}
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func parseRSAPrivateKey(input string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(input))
	if block == nil {
		return nil, errors.New("invalid PEM")
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
func parseRSAPublicKey(input string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(input))
	if block == nil {
		return nil, errors.New("invalid PEM")
	}
	if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
		if key, ok := cert.PublicKey.(*rsa.PublicKey); ok {
			return key, nil
		}
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return rsaKey, nil
}
func (c Config) callbackPublicKeyPEM() string {
	if strings.TrimSpace(c.PublicKeyPEM) != "" {
		return c.PublicKeyPEM
	}
	return c.PlatformCertPEM
}

// ResolvedMerchantSerialNo returns the configured serial number, or derives
// it from the merchant certificate. The latter keeps bootstrap key data
// complete when a certificate was supplied but its serial was not duplicated.
func (c Config) resolvedMerchantSerialNo() string {
	if serial := strings.ToUpper(strings.TrimSpace(c.MerchantSerialNo)); serial != "" {
		return serial
	}
	block, _ := pem.Decode([]byte(c.MerchantCertPEM))
	if block == nil {
		return ""
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil || certificate.SerialNumber == nil {
		return ""
	}
	return strings.ToUpper(certificate.SerialNumber.Text(16))
}
func sign(key *rsa.PrivateKey, data string) (string, error) {
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash([]byte(data)))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}
func hash(data []byte) []byte { sum := sha256.Sum256(data); return sum[:] }
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}
func providerMessage(code, message string, status int) string {
	if strings.TrimSpace(message) != "" {
		return truncate(message, 240)
	}
	if strings.TrimSpace(code) != "" {
		return truncate(code, 120)
	}
	return fmt.Sprintf("HTTP %d", status)
}
func parseUnix(raw string) (int64, error) {
	return strconv.ParseInt(raw, 10, 64)
}

func decryptResource(apiV3Key, associatedData, nonce, ciphertext string) ([]byte, error) {
	if len([]byte(apiV3Key)) != 32 {
		return nil, ErrInvalidConfig
	}
	block, err := aes.NewCipher([]byte(apiV3Key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len([]byte(nonce)) != gcm.NonceSize() {
		return nil, ErrInvalidCallback
	}
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, []byte(nonce), data, []byte(associatedData))
}
