// Package smsclient is a provider-neutral HTTPS adapter for one-time SMS codes.
package smsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrNotConfigured = errors.New("短信网关未配置")
var ErrRejected = errors.New("短信网关拒绝请求")

type Config struct {
	Endpoint, Authorization, Template string
	Timeout                           time.Duration
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
