package merchantsettlement

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/nats-io/nats.go"
)

// CommandSubject is request/reply only. It carries no bank account, token or
// payout secret; the actual settlement write remains in api-merchant.
const CommandSubject = "qixi.platform.merchant-settlement-command.v1"

type Command struct {
	SettlementID    uint64 `json:"settlement_id"`
	Action          string `json:"action"`
	OperatorID      uint64 `json:"operator_id"`
	IdempotencyKey  string `json:"idempotency_key"`
	ReviewNote      string `json:"review_note,omitempty"`
	PayoutReference string `json:"payout_reference,omitempty"`
}

type CommandResult struct {
	SettlementID uint64 `json:"settlement_id"`
	Status       string `json:"status,omitempty"`
	Code         string `json:"code,omitempty"`
}

type Client struct{ nc *nats.Conn }

func NewCommandClient(natsURL string) (*Client, error) {
	if strings.TrimSpace(natsURL) == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_platform_merchant_settlement_command"))
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

func (c *Client) Dispatch(ctx context.Context, command Command) (CommandResult, error) {
	if c == nil || c.nc == nil {
		return CommandResult{}, errors.New("店铺结算命令服务不可用")
	}
	if !validCommand(command) {
		return CommandResult{}, errors.New("结算命令参数错误")
	}
	body, err := json.Marshal(command)
	if err != nil {
		return CommandResult{}, err
	}
	msg, err := c.nc.RequestWithContext(ctx, CommandSubject, body)
	if err != nil {
		return CommandResult{}, errors.New("店铺结算命令未确认")
	}
	var out CommandResult
	if err := json.Unmarshal(msg.Data, &out); err != nil {
		return CommandResult{}, err
	}
	if out.SettlementID != command.SettlementID {
		return CommandResult{}, errors.New("店铺结算命令结果无效")
	}
	return out, nil
}

func validCommand(in Command) bool {
	if in.SettlementID == 0 || in.OperatorID == 0 || !validIdempotencyKey(in.IdempotencyKey) {
		return false
	}
	in.ReviewNote = strings.TrimSpace(in.ReviewNote)
	in.PayoutReference = strings.TrimSpace(in.PayoutReference)
	switch in.Action {
	case "approve":
		return len([]rune(in.ReviewNote)) <= 500
	case "reject":
		return in.ReviewNote != "" && len([]rune(in.ReviewNote)) <= 500
	case "mark_paid":
		return len([]rune(in.PayoutReference)) >= 3 && len([]rune(in.PayoutReference)) <= 128
	default:
		return false
	}
}

func validIdempotencyKey(value string) bool {
	length := len([]rune(strings.TrimSpace(value)))
	return length >= 8 && length <= 128
}
