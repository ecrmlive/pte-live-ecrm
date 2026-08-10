package feedbackmoderation

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nats-io/nats.go"
	"strings"
	"time"
)

const Subject = "qixi.platform.feedback-moderation-command.v1"

type Command struct {
	FeedbackID     uint64 `json:"feedback_id"`
	CategoryID     uint64 `json:"category_id,omitempty"`
	Action         string `json:"action"`
	Reply          string `json:"reply"`
	Name           string `json:"name,omitempty"`
	PID            uint64 `json:"pid,omitempty"`
	Sort           int    `json:"sort,omitempty"`
	Status         int    `json:"status,omitempty"`
	OperatorID     uint64 `json:"operator_id"`
	IdempotencyKey string `json:"idempotency_key"`
}
type Result struct {
	FeedbackID uint64 `json:"feedback_id"`
	CategoryID uint64 `json:"category_id,omitempty"`
	Status     string `json:"status,omitempty"`
	Code       string `json:"code,omitempty"`
}
type Client struct{ nc *nats.Conn }

func New(url string) (*Client, error) {
	if strings.TrimSpace(url) == "" {
		return nil, nil
	}
	nc, e := nats.Connect(url, nats.Name("pte_live_ecrm_api_platform_feedback_moderation"))
	if e != nil {
		return nil, e
	}
	return &Client{nc}, nil
}
func (c *Client) Close() {
	if c != nil && c.nc != nil {
		c.nc.Close()
	}
}
func Valid(x Command) bool {
	key, reply, name := strings.TrimSpace(x.IdempotencyKey), strings.TrimSpace(x.Reply), strings.TrimSpace(x.Name)
	if x.OperatorID == 0 || len([]rune(key)) < 8 || len([]rune(key)) > 128 || len([]rune(reply)) > 1000 || x.Sort < 0 || x.Sort > 9999 {
		return false
	}
	if (x.Action == "reply" && reply != "") || x.Action == "close" || x.Action == "delete" {
		return x.FeedbackID > 0
	}
	if x.Action == "category_create" {
		return x.CategoryID == 0 && len([]rune(name)) > 0 && len([]rune(name)) <= 32 && (x.Status == 0 || x.Status == 1)
	}
	if x.Action == "category_update" {
		return x.CategoryID > 0 && len([]rune(name)) > 0 && len([]rune(name)) <= 32 && (x.Status == 0 || x.Status == 1)
	}
	if x.Action == "category_status" {
		return x.CategoryID > 0 && (x.Status == 0 || x.Status == 1)
	}
	return x.Action == "category_delete" && x.CategoryID > 0
}
func (c *Client) Dispatch(ctx context.Context, x Command) (Result, error) {
	if c == nil || c.nc == nil {
		return Result{}, errors.New("反馈命令服务不可用")
	}
	if !Valid(x) {
		return Result{}, errors.New("反馈命令参数错误")
	}
	b, e := json.Marshal(x)
	if e != nil {
		return Result{}, e
	}
	t, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	m, e := c.nc.RequestWithContext(t, Subject, b)
	if e != nil {
		return Result{}, errors.New("反馈命令未确认")
	}
	var out Result
	if e = json.Unmarshal(m.Data, &out); e != nil {
		return Result{}, errors.New("反馈命令结果无效")
	}
	if out.Code != "" {
		return out, nil
	}
	if (x.FeedbackID > 0 && out.FeedbackID != x.FeedbackID) || (x.CategoryID > 0 && out.CategoryID != x.CategoryID) || (x.Action == "category_create" && out.CategoryID == 0) {
		return Result{}, errors.New("反馈命令结果无效")
	}
	return out, nil
}
