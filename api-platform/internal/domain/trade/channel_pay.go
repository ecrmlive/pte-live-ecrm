package trade

import (
	"context"
	"errors"
	"strings"

	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/pkg/paynotify"
	"gorm.io/gorm"
)

// PaymentSettings 第三方渠道开关（由 api-business 投影配置）。
type PaymentSettings struct {
	Sandbox      bool
	NotifySecret string
	Wechat       bool
	Alipay       bool
}

// PayIntent 异步渠道支付意图（沙箱下发 notify_token 供回调验签）。
type PayIntent struct {
	Status       string  `json:"status"`
	Channel      string  `json:"channel"`
	GroupOrderID uint    `json:"group_order_id"`
	OutTradeNo   string  `json:"out_trade_no"`
	PayPrice     float64 `json:"pay_price"`
	Sandbox      bool    `json:"sandbox"`
	NotifyToken  string  `json:"notify_token,omitempty"`
}

// ChannelNotifyInput 微信/支付宝沙箱回调入参。
type ChannelNotifyInput struct {
	GroupOrderID uint    `json:"group_order_id"`
	UID          uint    `json:"uid"`
	OutTradeNo   string  `json:"out_trade_no"`
	PayPrice     float64 `json:"pay_price"`
	NotifyToken  string  `json:"notify_token"`
}

func (s *Service) SetPayment(p PaymentSettings) { s.payment = p }

// IsChannelPayType wechat|alipay 走异步意图 + 回调。
func IsChannelPayType(payType string) bool {
	switch strings.ToLower(strings.TrimSpace(payType)) {
	case "wechat", "alipay":
		return true
	default:
		return false
	}
}

func (s *Service) channelEnabled(channel string) bool {
	switch strings.ToLower(strings.TrimSpace(channel)) {
	case "wechat":
		return s.payment.Wechat
	case "alipay":
		return s.payment.Alipay
	default:
		return false
	}
}

// CreatePayIntent 创建渠道支付意图；已支付则返回 paid 意图（无 token）。
func (s *Service) CreatePayIntent(ctx context.Context, uid, groupOrderID uint, payTypeStr string) (*PayIntent, error) {
	channel := strings.ToLower(strings.TrimSpace(payTypeStr))
	if !IsChannelPayType(channel) {
		return nil, ErrInvalidPayType
	}
	if !s.channelEnabled(channel) {
		return nil, ErrChannelDisabled
	}
	// 禁止把未实现的真实渠道伪装成已创建支付单；sandbox=true 才允许 HMAC 闭环。
	if !s.payment.Sandbox {
		return nil, ErrPaymentConfig
	}
	g, err := s.store.GetGroupOrder(ctx, groupOrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if g.UID != uid || g.IsDel == 1 {
		return nil, ErrForbidden
	}
	if g.ActivityType == ActivityTypePoints {
		return nil, ErrNotPointsProduct
	}
	if g.Paid == 1 {
		return &PayIntent{
			Status:       "paid",
			Channel:      channel,
			GroupOrderID: g.GroupOrderID,
			OutTradeNo:   g.GroupOrderSN,
			PayPrice:     g.PayPrice,
			Sandbox:      s.payment.Sandbox,
		}, nil
	}

	intent := &PayIntent{
		Status:       "pending",
		Channel:      channel,
		GroupOrderID: g.GroupOrderID,
		OutTradeNo:   g.GroupOrderSN,
		PayPrice:     g.PayPrice,
		Sandbox:      s.payment.Sandbox,
	}
	if s.payment.Sandbox {
		intent.NotifyToken = paynotify.MakeToken(
			s.payment.NotifySecret, channel, g.GroupOrderSN, uid, g.GroupOrderID, g.PayPrice,
		)
	}
	return intent, nil
}

// NotifyChannelPay 渠道回调：验签后走 PaySuccess（幂等）。
func (s *Service) NotifyChannelPay(ctx context.Context, channel string, in ChannelNotifyInput) (*GroupOrder, error) {
	channel = strings.ToLower(strings.TrimSpace(channel))
	if !IsChannelPayType(channel) {
		return nil, ErrInvalidPayType
	}
	if !s.channelEnabled(channel) {
		return nil, ErrChannelDisabled
	}
	if !s.payment.Sandbox {
		return nil, ErrPaymentConfig
	}
	if in.GroupOrderID == 0 || in.UID == 0 || strings.TrimSpace(in.NotifyToken) == "" {
		return nil, ErrBadNotify
	}
	g, err := s.store.GetGroupOrder(ctx, in.GroupOrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if g.UID != in.UID || g.IsDel == 1 {
		return nil, ErrForbidden
	}
	outTradeNo := strings.TrimSpace(in.OutTradeNo)
	if outTradeNo == "" {
		outTradeNo = g.GroupOrderSN
	}
	if outTradeNo != g.GroupOrderSN {
		return nil, ErrBadNotify
	}
	amount := in.PayPrice
	if amount <= 0 {
		amount = g.PayPrice
	}
	// 金额允许 0.01 浮点误差
	if absFloat(amount-g.PayPrice) > 0.01 {
		return nil, ErrBadNotify
	}
	if !paynotify.VerifyToken(s.payment.NotifySecret, channel, outTradeNo, in.UID, in.GroupOrderID, g.PayPrice, in.NotifyToken) {
		return nil, ErrBadNotify
	}
	return s.PaySuccess(ctx, in.UID, in.GroupOrderID, channel)
}

func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
