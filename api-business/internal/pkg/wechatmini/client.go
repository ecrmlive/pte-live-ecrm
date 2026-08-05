// Package wechatmini contains the server-side subset of the WeChat Mini
// Program login protocol. AppSecret is never returned to a caller or logged.
package wechatmini

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultCode2SessionEndpoint = "https://api.weixin.qq.com/sns/jscode2session"

var (
	ErrNotConfigured = errors.New("微信小程序登录未配置")
	ErrInvalidCode   = errors.New("微信登录凭证无效")
	ErrProvider      = errors.New("微信登录服务暂不可用")
)

type Config struct {
	AppID     string
	AppSecret string
}

type Session struct {
	OpenID  string
	UnionID string
}

type Client struct {
	HTTPClient *http.Client
	Endpoint   string
}

func (c *Client) Code2Session(ctx context.Context, config Config, code string) (Session, error) {
	if strings.TrimSpace(config.AppID) == "" || strings.TrimSpace(config.AppSecret) == "" {
		return Session{}, ErrNotConfigured
	}
	if strings.TrimSpace(code) == "" {
		return Session{}, ErrInvalidCode
	}
	endpoint := strings.TrimSpace(c.Endpoint)
	if endpoint == "" {
		endpoint = defaultCode2SessionEndpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return Session{}, ErrProvider
	}
	query := u.Query()
	query.Set("appid", config.AppID)
	query.Set("secret", config.AppSecret)
	query.Set("js_code", code)
	query.Set("grant_type", "authorization_code")
	u.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return Session{}, ErrProvider
	}
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return Session{}, ErrProvider
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil || resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return Session{}, ErrProvider
	}
	var payload struct {
		OpenID  string `json:"openid"`
		UnionID string `json:"unionid"`
		ErrCode int    `json:"errcode"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return Session{}, ErrProvider
	}
	if payload.ErrCode != 0 || strings.TrimSpace(payload.OpenID) == "" {
		return Session{}, ErrInvalidCode
	}
	return Session{OpenID: strings.TrimSpace(payload.OpenID), UnionID: strings.TrimSpace(payload.UnionID)}, nil
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: 8 * time.Second}
}
