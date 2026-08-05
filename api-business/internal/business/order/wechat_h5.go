package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/wechatpayv3"
)

var ErrH5PaymentConfig = errors.New("微信 H5 支付尚未配置")

// createWechatH5Pay uses WeChat MWEB for H5 tokens. It never falls back to a
// Native QR link: a phone browser cannot reliably scan its own payment code.
func (h *Handler) createWechatH5Pay(ctx context.Context, userID, groupOrderID uint64, payerClientIP string) (wechatPaymentIntent, error) {
	config, group, err := h.wechatConfigForOrder(ctx, userID, groupOrderID)
	if err != nil {
		return wechatPaymentIntent{}, err
	}
	if !config.ValidForH5() {
		return wechatPaymentIntent{}, ErrH5PaymentConfig
	}
	if group.PayStatus == "paid" {
		return wechatPaymentIntent{Status: "paid", Channel: "wechat", PaymentMode: "mweb", GroupOrderID: group.ID, OutTradeNo: group.OrderNo, PayPrice: group.PayAmount}, nil
	}
	if strings.TrimSpace(payerClientIP) == "" {
		return wechatPaymentIntent{}, wechatpayv3.ErrInvalidConfig
	}
	existing, createPrepay, err := h.beginWechatPayment(ctx, userID, group, "mweb")
	if err != nil {
		return wechatPaymentIntent{}, err
	}
	if !createPrepay {
		return existing, nil
	}
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	result, err := h.wechatClient.H5Prepay(ctx, config, wechatpayv3.H5Request{
		Description:   "CRM Live 商城订单",
		OutTradeNo:    group.OrderNo,
		Attach:        fmt.Sprintf("group_order:%d", group.ID),
		AmountCents:   decimalCents(group.PayAmount),
		PayerClientIP: payerClientIP,
		ExpireAt:      expiresAt,
	})
	if err != nil {
		_ = h.db.WithContext(ctx).Model(&paymentRow{}).Where("group_order_id = ? AND channel = ? AND status = ?", group.ID, "wechat", "processing").Update("status", "failed").Error
		return wechatPaymentIntent{}, err
	}
	payload, err := json.Marshal(prepayPayload{PaymentMode: "mweb", MWEBURL: result.MWEBURL, ExpiresAt: expiresAt.Format(time.RFC3339)})
	if err != nil {
		return wechatPaymentIntent{}, err
	}
	if err := h.db.WithContext(ctx).Model(&paymentRow{}).Where("group_order_id = ? AND channel = ? AND status = ?", group.ID, "wechat", "processing").Updates(map[string]any{"status": "created", "provider_payload": payload}).Error; err != nil {
		return wechatPaymentIntent{}, err
	}
	return wechatPaymentIntent{Status: "pending", Channel: "wechat", PaymentMode: "mweb", GroupOrderID: group.ID, OutTradeNo: group.OrderNo, PayPrice: group.PayAmount, MWEBURL: result.MWEBURL, ExpiresAt: expiresAt.Format(time.RFC3339)}, nil
}
