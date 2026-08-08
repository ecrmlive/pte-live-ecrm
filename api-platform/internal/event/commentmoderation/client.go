// Package commentmoderation sends platform moderation commands to the service
// that owns product-comment facts. The message never contains user profile or
// payment data.
package commentmoderation

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

const CommandSubject = "qixi.platform.product-comment-moderation-command.v1"

type Command struct {
	CommentID           uint64   `json:"comment_id"`
	Action              string   `json:"action"`
	OperatorID          uint64   `json:"operator_id"`
	IdempotencyKey      string   `json:"idempotency_key"`
	Note                string   `json:"note,omitempty"`
	ProductID           uint64   `json:"product_id,omitempty"`
	Score               int      `json:"score,omitempty"`
	Content             string   `json:"content,omitempty"`
	VirtualAuthorName   string   `json:"virtual_author_name,omitempty"`
	VirtualAuthorAvatar string   `json:"virtual_author_avatar,omitempty"`
	Sort                int      `json:"sort,omitempty"`
	Media               []string `json:"media,omitempty"`
	MediaSet            bool     `json:"media_set,omitempty"`
}
type Result struct {
	CommentID uint64 `json:"comment_id"`
	Status    string `json:"status,omitempty"`
	Code      string `json:"code,omitempty"`
}
type Client struct{ nc *nats.Conn }

func New(natsURL string) (*Client, error) {
	if strings.TrimSpace(natsURL) == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_platform_product_comment_moderation"))
	if err != nil {
		return nil, err
	}
	return &Client{nc: nc}, nil
}
func (c *Client) Close() {
	if c != nil && c.nc != nil {
		c.nc.Close()
	}
}
func (c *Client) Dispatch(ctx context.Context, command Command) (Result, error) {
	if c == nil || c.nc == nil {
		return Result{}, errors.New("评论审核命令服务不可用")
	}
	if !Valid(command) {
		return Result{}, errors.New("评论审核命令参数错误")
	}
	body, err := json.Marshal(command)
	if err != nil {
		return Result{}, err
	}
	requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	msg, err := c.nc.RequestWithContext(requestCtx, CommandSubject, body)
	if err != nil {
		return Result{}, errors.New("评论审核命令未确认")
	}
	var out Result
	if err = json.Unmarshal(msg.Data, &out); err != nil {
		return Result{}, err
	}
	if out.CommentID == 0 || (command.CommentID > 0 && out.CommentID != command.CommentID) {
		return Result{}, errors.New("评论审核命令结果无效")
	}
	return out, nil
}
func Valid(in Command) bool {
	key, note := strings.TrimSpace(in.IdempotencyKey), strings.TrimSpace(in.Note)
	if in.OperatorID == 0 || len([]rune(key)) < 8 || len([]rune(key)) > 128 || len([]rune(note)) > 500 {
		return false
	}
	switch in.Action {
	case "publish", "hide":
		return in.CommentID > 0
	case "create_virtual":
		return in.CommentID == 0 && in.ProductID > 0 && validVirtual(in)
	case "update_virtual":
		return in.CommentID > 0 && validVirtual(in)
	case "sort_virtual":
		return in.CommentID > 0 && in.Sort >= 0 && in.Sort <= 999999
	case "delete_virtual":
		return in.CommentID > 0
	default:
		return false
	}
}
func validVirtual(in Command) bool {
	avatar := strings.TrimSpace(in.VirtualAuthorAvatar)
	if in.Score < 1 || in.Score > 5 || len([]rune(strings.TrimSpace(in.Content))) == 0 || len([]rune(strings.TrimSpace(in.Content))) > 2000 || len([]rune(strings.TrimSpace(in.VirtualAuthorName))) == 0 || len([]rune(strings.TrimSpace(in.VirtualAuthorName))) > 64 || len([]rune(avatar)) > 1024 || in.Sort < 0 || in.Sort > 999999 || len(in.Media) > 9 {
		return false
	}
	for _, media := range in.Media {
		if len([]rune(strings.TrimSpace(media))) == 0 || len([]rune(strings.TrimSpace(media))) > 1024 {
			return false
		}
	}
	return true
}
