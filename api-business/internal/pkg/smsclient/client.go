// Package smsclient is a provider-neutral HTTPS adapter for one-time SMS codes.
package smsclient

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var ErrNotConfigured = errors.New("短信网关未配置")
var ErrRejected = errors.New("短信网关拒绝请求")

type Config struct {
	Endpoint, Authorization, Template string
	Timeout                           time.Duration
}

// TencentConfig contains the approved SMS application values. AppKey is kept
// encrypted in the platform database and is never returned by an API.
type TencentConfig struct {
	AppKey      string
	SDKAppID    string
	SignContent string
	TemplateID  string
}
type Client struct{ httpClient *http.Client }

func New() *Client { return &Client{httpClient: &http.Client{Timeout: 5 * time.Second}} }

// Send never logs or returns the code. Gateways must accept a JSON POST body
// and reply with any 2xx status only after accepting the delivery request.
func (c *Client) Send(ctx context.Context, cfg Config, phone, code string, expiresIn time.Duration) error {
	if strings.TrimSpace(cfg.Endpoint) == "" || strings.TrimSpace(cfg.Authorization) == "" || strings.TrimSpace(cfg.Template) == "" {
		return ErrNotConfigured
	}
	u, err := url.Parse(cfg.Endpoint)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return ErrNotConfigured
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}
	body, err := json.Marshal(map[string]any{"phone": phone, "template": cfg.Template, "code": code, "expires_in": int(expiresIn.Seconds())})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", cfg.Authorization)
	client := *c.httpClient
	client.Timeout = cfg.Timeout
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ErrRejected
	}
	return nil
}

// SendTencent calls Tencent Cloud's AppKey SMS endpoint. The request proof is
// calculated locally; neither key nor verification code is logged.
func (c *Client) SendTencent(ctx context.Context, cfg TencentConfig, phone, code string) error {
	if strings.TrimSpace(cfg.AppKey) == "" || strings.TrimSpace(cfg.SDKAppID) == "" || strings.TrimSpace(cfg.SignContent) == "" || strings.TrimSpace(cfg.TemplateID) == "" {
		return ErrNotConfigured
	}
	if _, err := strconv.ParseUint(cfg.SDKAppID, 10, 64); err != nil {
		return ErrNotConfigured
	}
	templateID, err := strconv.ParseUint(cfg.TemplateID, 10, 64)
	if err != nil {
		return ErrNotConfigured
	}
	random := strconv.FormatInt(time.Now().UnixNano(), 10)
	sentAt := time.Now().Unix()
	signSource := "appkey=" + cfg.AppKey + "&random=" + random + "&time=" + strconv.FormatInt(sentAt, 10) + "&mobile=" + phone
	signature := sha256.Sum256([]byte(signSource))
	body, err := json.Marshal(map[string]any{
		"ext":    "",
		"extend": "",
		"params": []string{code},
		"sig":    hex.EncodeToString(signature[:]),
		"sign":   cfg.SignContent,
		"tel":    map[string]string{"mobile": phone, "nationcode": "86"},
		"time":   sentAt,
		"tpl_id": templateID,
	})
	if err != nil {
		return err
	}
	endpoint := "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=" + url.QueryEscape(cfg.SDKAppID) + "&random=" + url.QueryEscape(random)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := *c.httpClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ErrRejected
	}
	var result struct {
		Result int `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.Result != 0 {
		return ErrRejected
	}
	return nil
}
