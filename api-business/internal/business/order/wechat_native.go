package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api-business/internal/paymentconfig"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/wechatpayv3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrPaymentProcessing = errors.New("支付单正在创建，请稍后刷新支付页面")

type wechatPaymentIntent struct {
	Status       string  `json:"status"`
	Channel      string  `json:"channel"`
	PaymentMode  string  `json:"payment_mode"`
	GroupOrderID uint64  `json:"group_order_id"`
	OutTradeNo   string  `json:"out_trade_no"`
	PayPrice     float64 `json:"pay_price"`
	CodeURL      string  `json:"code_url,omitempty"`
	ExpiresAt    string  `json:"expires_at,omitempty"`
}

type prepayPayload struct {
	CodeURL   string `json:"code_url"`
	ExpiresAt string `json:"expires_at"`
}

func (h *Handler) createWechatNativePay(ctx context.Context, userID, groupOrderID uint64) (wechatPaymentIntent, error) {
	config, group, err := h.wechatConfigForOrder(ctx, userID, groupOrderID)
	if err != nil {
		return wechatPaymentIntent{}, err
	}
	if group.PayStatus == "paid" {
		return wechatPaymentIntent{Status: "paid", Channel: "wechat", PaymentMode: "native", GroupOrderID: group.ID, OutTradeNo: group.OrderNo, PayPrice: group.PayAmount}, nil
	}

	existing, createPrepay, err := h.beginWechatPayment(ctx, userID, group)
	if err != nil {
		return wechatPaymentIntent{}, err
	}
	if !createPrepay {
		return existing, nil
	}

	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	response, err := h.wechatClient.NativePrepay(ctx, config, wechatpayv3.NativeRequest{
		Description: "七禧商城订单",
		OutTradeNo:  group.OrderNo,
		Attach:      fmt.Sprintf("group_order:%d", group.ID),
		AmountCents: decimalCents(group.PayAmount),
		ExpireAt:    expiresAt,
	})
	if err != nil {
		_ = h.db.WithContext(ctx).Model(&paymentRow{}).Where("group_order_id = ? AND channel = ? AND status = ?", group.ID, "wechat", "processing").Update("status", "failed").Error
		return wechatPaymentIntent{}, err
	}
	payload, err := json.Marshal(prepayPayload{CodeURL: response.CodeURL, ExpiresAt: expiresAt.Format(time.RFC3339)})
	if err != nil {
		return wechatPaymentIntent{}, err
	}
	if err := h.db.WithContext(ctx).Model(&paymentRow{}).Where("group_order_id = ? AND channel = ? AND status = ?", group.ID, "wechat", "processing").Updates(map[string]any{"status": "created", "provider_payload": payload}).Error; err != nil {
		return wechatPaymentIntent{}, err
	}
	return wechatPaymentIntent{Status: "pending", Channel: "wechat", PaymentMode: "native", GroupOrderID: group.ID, OutTradeNo: group.OrderNo, PayPrice: group.PayAmount, CodeURL: response.CodeURL, ExpiresAt: expiresAt.Format(time.RFC3339)}, nil
}

func (h *Handler) beginWechatPayment(ctx context.Context, userID uint64, group groupRow) (wechatPaymentIntent, bool, error) {
	var intent wechatPaymentIntent
	createPrepay := false
	err := h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var locked groupRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ?", group.ID, userID).First(&locked).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOrderOwnership
			}
			return err
		}
		intent = wechatPaymentIntent{Status: "pending", Channel: "wechat", PaymentMode: "native", GroupOrderID: locked.ID, OutTradeNo: locked.OrderNo, PayPrice: locked.PayAmount}
		if locked.PayStatus == "paid" {
			intent.Status = "paid"
			return nil
		}
		if locked.PayStatus != "pending" {
			return ErrOrderNotPayable
		}

		var payment paymentRow
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("group_order_id = ? AND channel = ?", locked.ID, "wechat").First(&payment).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			payment = paymentRow{GroupOrderID: locked.ID, Channel: "wechat", TransactionNo: locked.OrderNo, Amount: locked.PayAmount, Status: "processing"}
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
			createPrepay = true
			return nil
		}
		if err != nil {
			return err
		}
		switch payment.Status {
		case "succeeded":
			intent.Status = "paid"
			return nil
		case "created":
			var payload prepayPayload
			if len(payment.ProviderPayload) != 0 && json.Unmarshal(payment.ProviderPayload, &payload) == nil && strings.TrimSpace(payload.CodeURL) != "" {
				intent.CodeURL, intent.ExpiresAt = payload.CodeURL, payload.ExpiresAt
				return nil
			}
			return ErrPaymentProcessing
		case "processing":
			return ErrPaymentProcessing
		case "failed":
			if err := tx.Model(&paymentRow{}).Where("id = ?", payment.ID).Updates(map[string]any{"status": "processing", "provider_payload": nil}).Error; err != nil {
				return err
			}
			createPrepay = true
			return nil
		default:
			return ErrOrderNotPayable
		}
	})
	return intent, createPrepay, err
}

func (h *Handler) wechatConfigForOrder(ctx context.Context, userID, groupOrderID uint64) (wechatpayv3.Config, groupRow, error) {
	if err := assertPaymentChannelAvailable(ctx, h.db, h.configs, h.platformConfigs, userID, groupOrderID, "wechat"); err != nil {
		return wechatpayv3.Config{}, groupRow{}, err
	}
	var group groupRow
	if err := h.db.WithContext(ctx).Where("id = ? AND user_id = ?", groupOrderID, userID).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return wechatpayv3.Config{}, groupRow{}, ErrOrderOwnership
		}
		return wechatpayv3.Config{}, groupRow{}, err
	}
	subject, err := paymentSubject(ctx, h.db, userID, groupOrderID)
	if err != nil {
		return wechatpayv3.Config{}, groupRow{}, err
	}
	platformOwned := subject.MerchantID == 0 && subject.StoreID == 0
	values, err := loadPaymentValues(ctx, h.configs, h.platformConfigs, subject, platformOwned)
	if err != nil {
		return wechatpayv3.Config{}, groupRow{}, err
	}
	return wechatConfig(values, platformOwned), group, nil
}

func wechatConfig(values paymentconfig.Values, platformOwned bool) wechatpayv3.Config {
	if platformOwned {
		return wechatpayv3.Config{AppID: values["wechat_app_id"], MchID: values["wechat_mch_id"], MerchantSerialNo: values["wechat_serial_no"], MerchantPrivateKey: values["wechat_private_key"], MerchantCertPEM: values["wechat_merchant_cert"], APIv3Key: values["wechat_api_v3_key"], NotifyURL: values["wechat_notify_url"], PublicKeyID: values["wechat_public_key_id"], PublicKeyPEM: values["wechat_public_key"], PlatformCertSerial: values["wechat_platform_cert_serial"], PlatformCertPEM: values["wechat_platform_cert"]}
	}
	return wechatpayv3.Config{AppID: values["app_id"], MchID: values["mch_id"], MerchantSerialNo: values["serial_no"], MerchantPrivateKey: values["private_key"], APIv3Key: values["api_v3_key"], NotifyURL: values["notify_url"]}
}

func decimalCents(amount float64) int64 { return int64(math.Round(amount * 100)) }
