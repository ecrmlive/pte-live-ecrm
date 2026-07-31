package order

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrOrderNotCancellable = errors.New("当前订单不可取消")

// CancelPending closes a not-yet-paid group order. Coupon locks belong to the
// same transaction so a failed cancellation cannot leave a user coupon stuck.
func CancelPending(ctx context.Context, db *gorm.DB, userID, groupOrderID uint64) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var group groupRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", groupOrderID, userID).First(&group).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOrderOwnership
			}
			return err
		}
		if group.PayStatus == "closed" {
			return nil
		}
		if group.PayStatus != "pending" {
			return ErrOrderNotCancellable
		}
		if err := tx.Model(&groupRow{}).Where("id = ? AND pay_status = ?", group.ID, "pending").
			Update("pay_status", "closed").Error; err != nil {
			return err
		}
		if err := tx.Model(&orderRow{}).Where("group_order_id = ? AND status = ?", group.ID, "pending_pay").
			Update("status", "cancelled").Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_b_coupon_user").
			Where("user_id = ? AND used_order_id = ? AND status = ?", userID, group.ID, "locked").
			Updates(map[string]any{"status": "unused", "used_order_id": nil}).Error
	})
}
