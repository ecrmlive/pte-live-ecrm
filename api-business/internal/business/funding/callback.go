package funding

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/wechatpayv3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MarkWechatPaid consumes a transaction that has already passed WeChat v3
// signature, decryption, AppID and merchant checks in the shared callback
// adapter. handled=false means this is not a funding payment and lets the
// caller continue with the normal merchandise-order callback.
func MarkWechatPaid(ctx context.Context, db *gorm.DB, notify wechatpayv3.CallbackTransaction) (handled bool, err error) {
	var exists paymentRow
	err = db.WithContext(ctx).Where("out_trade_no = ? AND channel = ?", notify.OutTradeNo, "wechat").Take(&exists).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var payment paymentRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", exists.ID).Take(&payment).Error; err != nil {
			return err
		}
		if payment.Channel != "wechat" || payment.OutTradeNo != notify.OutTradeNo || cents(payment.Amount) != notify.AmountCents {
			return ErrOrderNotPayable
		}
		var callback struct {
			ID uint64 `gorm:"column:id"`
		}
		err := tx.Table("qixi_crm_b_funding_payment_callback").Where("channel = ? AND provider_event_id = ?", "wechat", notify.EventID).Take(&callback).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		payload, err := json.Marshal(json.RawMessage(notify.RawPayload))
		if err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_funding_payment_callback").Create(map[string]any{"channel": "wechat", "provider_event_id": notify.EventID, "out_trade_no": payment.OutTradeNo, "payload": payload, "verified": true, "processed_at": notify.SuccessTime}).Error; err != nil {
			return err
		}
		if payment.Status == "succeeded" {
			return nil
		}
		if payment.Status != "created" && payment.Status != "processing" {
			return ErrOrderNotPayable
		}
		switch payment.OrderType {
		case "recharge":
			if err := settleRecharge(tx, payment, notify.SuccessTime); err != nil {
				return err
			}
		case "svip":
			if err := settleSVIP(tx, payment, notify.SuccessTime); err != nil {
				return err
			}
		default:
			return ErrOrderNotPayable
		}
		return tx.Model(&paymentRow{}).Where("id = ? AND status IN ?", payment.ID, []string{"created", "processing"}).Updates(map[string]any{"status": "succeeded", "provider_transaction_no": notify.TransactionID, "callback_idempotency_key": notify.EventID, "paid_at": notify.SuccessTime}).Error
	})
}

func settleRecharge(tx *gorm.DB, payment paymentRow, paidAt time.Time) error {
	var order rechargeOrder
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ?", payment.FundingOrderID, payment.UserID).Take(&order).Error; err != nil {
		return err
	}
	if order.RechargeNo != payment.OutTradeNo || order.Status == "closed" || cents(order.Amount) != cents(payment.Amount) {
		return ErrOrderNotPayable
	}
	if order.Status == "paid" {
		return nil
	}
	credit := roundedAmount(order.Amount + order.BonusAmount)
	if err := tx.Exec("INSERT INTO qixi_crm_b_member_account (user_id, balance) VALUES (?, ?) ON DUPLICATE KEY UPDATE balance = balance + VALUES(balance)", order.UserID, credit).Error; err != nil {
		return err
	}
	key := "recharge:" + order.RechargeNo
	if err := tx.Exec("INSERT INTO qixi_crm_b_asset_ledger (user_id, asset_type, amount, reference_type, reference_id, idempotency_key) VALUES (?, 'balance', ?, 'recharge', ?, ?)", order.UserID, credit, order.RechargeNo, key).Error; err != nil {
		return err
	}
	result := tx.Model(&rechargeOrder{}).Where("id = ? AND status = ?", order.ID, "pending").Updates(map[string]any{"status": "paid", "paid_at": paidAt})
	if result.Error != nil || result.RowsAffected != 1 {
		if result.Error != nil {
			return result.Error
		}
		return ErrOrderNotPayable
	}
	return nil
}

func settleSVIP(tx *gorm.DB, payment paymentRow, paidAt time.Time) error {
	var order svipOrder
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ?", payment.FundingOrderID, payment.UserID).Take(&order).Error; err != nil {
		return err
	}
	if order.OrderNo != payment.OutTradeNo || order.Status == "closed" || cents(order.Amount) != cents(payment.Amount) || !validSVIPOrder(order) {
		return ErrOrderNotPayable
	}
	if order.Status == "paid" {
		return nil
	}
	// The upsert creates a row only in this transaction, then FOR UPDATE makes
	// concurrent periodic purchases serialize before calculating the expiry.
	if err := tx.Exec("INSERT INTO qixi_crm_b_user_svip (user_id, status, expires_at) VALUES (?, 'period', NULL) ON DUPLICATE KEY UPDATE user_id = VALUES(user_id)", order.UserID).Error; err != nil {
		return err
	}
	var current struct {
		Status    string     `gorm:"column:status"`
		ExpiresAt *time.Time `gorm:"column:expires_at"`
	}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("qixi_crm_b_user_svip").Where("user_id = ?", order.UserID).Take(&current).Error; err != nil {
		return err
	}
	if current.Status == "lifetime" {
		return ErrLifetimeActive
	}
	if order.PlanType == "trial" {
		var used int64
		if err := tx.Model(&svipOrder{}).Where("user_id = ? AND plan_type = ? AND status = ? AND id <> ?", order.UserID, "trial", "paid", order.ID).Count(&used).Error; err != nil {
			return err
		}
		if used != 0 {
			return ErrTrialAlreadyUsed
		}
	}
	status := order.PlanType
	var expiresAt *time.Time
	if order.PlanType == "lifetime" {
		status = "lifetime"
	} else {
		start := paidAt.UTC()
		if current.ExpiresAt != nil && current.ExpiresAt.After(start) {
			start = *current.ExpiresAt
		}
		next := start.AddDate(0, 0, int(*order.DurationDays))
		expiresAt = &next
	}
	if err := tx.Table("qixi_crm_b_user_svip").Where("user_id = ?", order.UserID).Updates(map[string]any{"status": status, "expires_at": expiresAt}).Error; err != nil {
		return err
	}
	result := tx.Model(&svipOrder{}).Where("id = ? AND status = ?", order.ID, "pending").Updates(map[string]any{"status": "paid", "paid_at": paidAt})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrOrderNotPayable
	}
	return nil
}

func validSVIPOrder(order svipOrder) bool {
	if order.PlanType == "lifetime" {
		return order.Amount > 0
	}
	return order.DurationDays != nil && *order.DurationDays > 0 && order.Amount > 0 && (order.PlanType == "trial" || order.PlanType == "period")
}
