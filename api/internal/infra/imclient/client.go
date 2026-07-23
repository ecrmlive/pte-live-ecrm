package imclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client 调用 pte-live-im mall S2S 集成。
type Client struct {
	baseURL string
	token   string
	appID   int
	http    *http.Client
}

func New(baseURL, token string, appID int) *Client {
	if appID <= 0 {
		appID = 30001
	}
	return &Client{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		token:   strings.TrimSpace(token),
		appID:   appID,
		http:    &http.Client{Timeout: 8 * time.Second},
	}
}

func (c *Client) Enabled() bool {
	return c != nil && c.baseURL != "" && c.token != ""
}

func (c *Client) AppID() int { return c.appID }

type UserSigResult struct {
	UserSig    string `json:"user_sig"`
	WsURL      string `json:"ws_url"`
	AppID      string `json:"app_id"`
	SDKAppID   string `json:"sdk_app_id"`
	UserID     string `json:"user_id"`
	Identifier string `json:"identifier"`
	ExpireAt   int64  `json:"expire_at"`
}

type apiEnvelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (c *Client) IssueUserSig(ctx context.Context, userID int64, identifier, userType, deviceID, platform string) (*UserSigResult, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("im client disabled")
	}
	if identifier == "" {
		identifier = fmt.Sprintf("%d", userID)
	}
	if userType == "" {
		userType = "user"
	}
	body := map[string]any{
		"user_id":    fmt.Sprintf("%d", userID),
		"identifier": identifier,
		"user_type":  userType,
		"device_id":  deviceID,
		"platform":   platform,
		"scene":      "chat",
		"expire":     86400,
	}
	var out UserSigResult
	if err := c.post(ctx, "/api/v1/integrations/mall/usersig", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type OpenSingleResult struct {
	ConversationID uint64 `json:"conversation_id"`
	ID             uint64 `json:"id"`
}

func (c *Client) OpenSingle(ctx context.Context, userID, peerUserID int64) (uint64, error) {
	if !c.Enabled() {
		return 0, fmt.Errorf("im client disabled")
	}
	body := map[string]any{
		"user_id":      userID,
		"peer_user_id": peerUserID,
	}
	var out OpenSingleResult
	if err := c.post(ctx, "/api/v1/integrations/mall/conversation/open-single", body, &out); err != nil {
		return 0, err
	}
	id := out.ConversationID
	if id == 0 {
		id = out.ID
	}
	if id == 0 {
		return 0, fmt.Errorf("empty conversation_id")
	}
	return id, nil
}

func (c *Client) post(ctx context.Context, path string, body any, dest any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Pte-Mall-Integration-Token", c.token)
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("im mall auth failed: %s", strings.TrimSpace(string(data)))
	}
	var env apiEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return fmt.Errorf("im response decode: %w body=%s", err, truncate(string(data), 200))
	}
	// pte-live-im Success = code 1
	if env.Code != 1 {
		if env.Msg == "" {
			env.Msg = "im api error"
		}
		return fmt.Errorf("%s", env.Msg)
	}
	if dest == nil || len(env.Data) == 0 || string(env.Data) == "null" {
		return nil
	}
	return json.Unmarshal(env.Data, dest)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
