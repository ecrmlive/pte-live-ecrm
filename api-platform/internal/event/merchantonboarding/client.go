// Package merchantonboarding sends the trusted command that provisions a
// merchant-owned store. Platform must never write qixi_crm_m_* tables.
package merchantonboarding

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/nats-io/nats.go"
)

const Subject = "qixi.platform.merchant-onboarding.v1"

type Request struct {
	ApplicationID uint   `json:"application_id"`
	RegionID      uint   `json:"region_id"`
	MerchantName  string `json:"merchant_name"`
	ContactName   string `json:"contact_name"`
	ContactMobile string `json:"contact_mobile"`
	Account       string `json:"account"`
	PasswordHash  string `json:"password_hash"`
}

type Result struct {
	MerchantID uint   `json:"merchant_id"`
	StoreID    uint   `json:"store_id"`
	Account    string `json:"account"`
	Error      string `json:"error,omitempty"`
}

type Client struct{ nc *nats.Conn }

func New(natsURL string) (*Client, error) {
	if strings.TrimSpace(natsURL) == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_platform_merchant_onboarding"))
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

func (c *Client) Provision(ctx context.Context, in Request) (Result, error) {
	if c == nil || c.nc == nil {
		return Result{}, errors.New("店铺开通命令服务不可用")
	}
	if in.ApplicationID == 0 || in.RegionID == 0 || strings.TrimSpace(in.MerchantName) == "" || strings.TrimSpace(in.Account) == "" || !strings.HasPrefix(in.PasswordHash, "$2") {
		return Result{}, errors.New("开通店铺参数不完整")
	}
	body, err := json.Marshal(in)
	if err != nil {
		return Result{}, err
	}
	msg, err := c.nc.RequestWithContext(ctx, Subject, body)
	if err != nil {
		return Result{}, errors.New("店铺开通命令未确认")
	}
	var out Result
	if err := json.Unmarshal(msg.Data, &out); err != nil {
		return Result{}, err
	}
	if out.Error != "" {
		return Result{}, errors.New(out.Error)
	}
	if out.MerchantID == 0 || out.StoreID == 0 {
		return Result{}, errors.New("店铺开通结果无效")
	}
	return out, nil
}
