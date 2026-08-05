package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	merchantstock "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantstock"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrInsufficientBalance = errors.New("账户余额不足")

// PayBalance is a server-authoritative wallet payment. It locks the group
// order before the account, conditionally deducts the business-owned balance,
// and writes an immutable debit with the same transaction as fulfillment.
func PayBalance(ctx context.Context, db *gorm.DB, userID, groupOrderID uint64) (CreatedOrder, error) {
	var result CreatedOrder
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var group groupRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ?", groupOrderID, userID).First(&group).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOrderOwnership
			}
			return err
		}
		result = CreatedOrder{GroupOrderID: group.ID, GroupOrderNo: group.OrderNo, PayCents: int64(group.PayAmount * 100), TotalQuantity: group.TotalQuantity}
		if groupPayStatus(group) == "paid" {
			return nil
		}
		if group.ActivityType == pointsOrderActivityType || group.PayAmount <= 0 || groupPayStatus(group) != "pending" {
			return ErrOrderNotPayable
		}
		if err := merchantstock.ReservationsReady(tx, group.ID); err != nil {
			return err
		}
		var account struct {
			Balance float64 `gorm:"column:balance"`
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("qixi_crm_b_member_account").Select("balance").Where("user_id = ?", userID).Take(&account).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInsufficientBalance
			}
			return err
		}
		if decimalCents(account.Balance) < decimalCents(group.PayAmount) {
			return ErrInsufficientBalance
		}
		var payment paymentRow
		if err := tx.Where("group_order_id = ? AND channel = ?", group.ID, "balance").FirstOrCreate(&payment, &paymentRow{GroupOrderID: group.ID, Channel: "balance", TransactionNo: orderNo("B"), Amount: group.PayAmount, Status: "created"}).Error; err != nil {
			return err
		}
		if decimalCents(payment.Amount) != decimalCents(group.PayAmount) || payment.Status == "succeeded" {
			return ErrOrderNotPayable
		}
		debit := tx.Table("qixi_crm_b_member_account").Where("user_id = ? AND balance >= ?", userID, group.PayAmount).Update("balance", gorm.Expr("balance - ?", group.PayAmount))
		if debit.Error != nil {
			return debit.Error
		}
		if debit.RowsAffected != 1 {
			return ErrInsufficientBalance
		}
		ledgerKey := fmt.Sprintf("balance-order:%d", group.ID)
		if err := tx.Exec("INSERT INTO qixi_crm_b_asset_ledger (user_id, asset_type, amount, reference_type, reference_id, idempotency_key) VALUES (?, 'balance', ?, 'order_payment', ?, ?)", userID, -group.PayAmount, group.OrderNo, ledgerKey).Error; err != nil {
			return err
		}
		now := time.Now().UTC()
		updated := tx.Model(&groupRow{}).Where("id = ? AND pay_status = ?", group.ID, "pending").Updates(map[string]any{"pay_status": "paid", "pay_channel": "balance", "paid_at": now})
		if updated.Error != nil {
			return updated.Error
		}
		if updated.RowsAffected != 1 {
			return ErrOrderNotPayable
		}
		if err := tx.Model(&orderRow{}).Where("group_order_id = ? AND status = ?", group.ID, "pending_pay").Updates(map[string]any{"status": "paid", "paid_at": now}).Error; err != nil {
			return err
		}
		if err := enqueueStockActionForGroup(tx, "confirm", group.ID); err != nil {
			return err
		}
		if err := settleOrderActivity(tx, group); err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_coupon_user").Where("user_id = ? AND used_order_id = ? AND status = ?", userID, group.ID, "locked").Update("status", "used").Error; err != nil {
			return err
		}
		return tx.Model(&paymentRow{}).Where("id = ? AND status = ?", payment.ID, "created").Updates(map[string]any{"status": "succeeded", "provider_transaction_no": fmt.Sprintf("balance-%d", group.ID), "paid_at": now}).Error
	})
	return result, err
}
