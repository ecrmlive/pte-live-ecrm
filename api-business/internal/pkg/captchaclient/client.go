// Package captchaclient implements the server-to-server pte-tools-captcha contract.
// It deliberately never exposes HMAC credentials to a browser or mobile client.
package captchaclient

import (
	"bytes"
	"context"
	"crypto/hmac"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/config"
)

var ErrUnavailable = errors.New("验证码服务未配置")

// Client only communicates with pte-tools-captcha over the private service network.
type Client struct {
	baseURL string
	appID   string
	secret  []byte
	http    *http.Client
}

func New(cfg config.CaptchaConfig) (*Client, error) {
	if !cfg.Enabled {
		return nil, ErrUnavailable
	}
	secret := cfg.Secret()
	if strings.TrimSpace(cfg.BaseURL) == "" || strings.TrimSpace(cfg.ApplicationID) == "" || len(secret) < 32 {
		return nil, ErrUnavailable
	}
	return &Client{
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		appID:   cfg.ApplicationID,
		secret:  []byte(secret),
		http:    &http.Client{Timeout: cfg.Timeout()},
	}, nil
}

func (c *Client) Create(ctx context.Context, payload any) (json.RawMessage, error) {
	return c.post(ctx, "/api/v1/challenges", payload)
}

func (c *Client) Verify(ctx context.Context, challengeID string, payload any) (json.RawMessage, error) {
	challengeID = strings.TrimSpace(challengeID)
	if challengeID == "" {
		return nil, fmt.Errorf("challenge_id is required")
	}
	return c.post(ctx, "/api/v1/challenges/"+challengeID+"/verify", payload)
}

func (c *Client) Validate(ctx context.Context, verificationToken, action string) error {
	body, err := c.post(ctx, "/api/v1/tokens/validate", map[string]string{
		"verification_token": verificationToken,
		"action":             action,
	})
	if err != nil {
		return err
	}
	var result struct {
		Valid bool `json:"valid"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析验证码校验响应: %w", err)
	}
	if !result.Valid {
		return errors.New("验证码校验失败")
	}
	return nil
}

func (c *Client) post(ctx context.Context, path string, payload any) (json.RawMessage, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("编码验证码请求: %w", err)
	}
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce, err := randomNonce()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建验证码请求: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Captcha-App-ID", c.appID)
	req.Header.Set("X-Captcha-Timestamp", timestamp)
	req.Header.Set("X-Captcha-Nonce", nonce)
	req.Header.Set("X-Captcha-Signature", sign(c.secret, http.MethodPost, path, body, timestamp, nonce))

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("验证码服务不可用: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("读取验证码响应: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var failure struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(data, &failure) == nil && failure.Error != "" {
			return nil, errors.New(failure.Error)
		}
		return nil, fmt.Errorf("验证码服务返回 HTTP %d", resp.StatusCode)
	}
	return json.RawMessage(data), nil
}

func randomNonce() (string, error) {
	b := make([]byte, 24)
	if _, err := cryptorand.Read(b); err != nil {
		return "", fmt.Errorf("生成验证码请求随机数: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func sign(secret []byte, method, path string, body []byte, timestamp, nonce string) string {
	digest := sha256.Sum256(body)
	canonical := strings.ToUpper(method) + "\n" + path + "\n" + hex.EncodeToString(digest[:]) + "\n" + timestamp + "\n" + nonce
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(canonical))
	return hex.EncodeToString(mac.Sum(nil))
}
