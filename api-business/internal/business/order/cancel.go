package order

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const reservationOrderActivityType = 30

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
		return closePendingGroup(tx, group)
	})
}

// closePendingGroup is the only business transition from pending to closed.
// It is used by both a user cancellation and the scheduled unpaid-order sweep.
func closePendingGroup(tx *gorm.DB, group groupRow) error {
	if group.PayStatus != "pending" {
		return ErrOrderNotCancellable
	}
	closed := tx.Model(&groupRow{}).Where("id = ? AND pay_status = ?", group.ID, "pending").Update("pay_status", "closed")
	if closed.Error != nil {
		return closed.Error
	}
	if closed.RowsAffected != 1 {
		return ErrOrderNotCancellable
	}
	if err := tx.Model(&orderRow{}).Where("group_order_id = ? AND status = ?", group.ID, "pending_pay").Update("status", "cancelled").Error; err != nil {
		return err
	}
	if err := enqueueStockActionForGroup(tx, "release", group.ID); err != nil {
		return err
	}
	if err := releaseOrderActivity(tx, group); err != nil {
		return err
	}
	if group.ActivityType == reservationOrderActivityType {
		var bookings []struct {
			ID     uint64 `gorm:"column:id"`
			SlotID uint64 `gorm:"column:slot_id"`
		}
		if err := tx.Table("qixi_crm_b_reservation_booking AS b").Select("b.id,b.slot_id").Joins("JOIN qixi_crm_b_order AS o ON o.id = b.order_id").Where("o.group_order_id = ? AND b.status = 1", group.ID).Find(&bookings).Error; err != nil {
			return err
		}
		if len(bookings) > 0 {
			ids := make([]uint64, 0, len(bookings))
			for _, booking := range bookings {
				ids = append(ids, booking.ID)
			}
			if err := tx.Table("qixi_crm_b_reservation_booking").Where("id IN ? AND status = 1", ids).Update("status", 0).Error; err != nil {
				return err
			}
			for _, booking := range bookings {
				if err := tx.Table("qixi_crm_b_reservation_slot").Where("attr_reservation_id = ? AND use_num > 0", booking.SlotID).UpdateColumn("use_num", gorm.Expr("use_num - 1")).Error; err != nil {
					return err
				}
			}
		}
	}
	if err := tx.Model(&paymentRow{}).Where("group_order_id = ? AND status IN ?", group.ID, []string{"created", "processing", "failed"}).Update("status", "closed").Error; err != nil {
		return err
	}
	if err := restoreIntegral(tx, group); err != nil {
		return err
	}
	return tx.Table("qixi_crm_b_coupon_user").
		Where("user_id = ? AND used_order_id = ? AND status = ?", group.UserID, group.ID, "locked").
		Updates(map[string]any{"status": "unused", "used_order_id": nil}).Error
}
