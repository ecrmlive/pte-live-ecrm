// Package ordertimeout closes business-database orders that were never paid.
//
// The job service owns the schedule; this package deliberately operates only
// on qixi_crm_b_* tables and has no dependency on legacy mall tables.
package ordertimeout

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	presellActivityType     = 2
	assistActivityType      = 3
	combinationActivityType = 4
	reservationActivityType = 30
)

var errNoExpiredPendingOrder = errors.New("no expired pending order")

type groupRow struct {
	ID           uint64    `gorm:"column:id"`
	UserID       uint64    `gorm:"column:user_id"`
	PayStatus    string    `gorm:"column:pay_status"`
	ActivityType int       `gorm:"column:activity_type"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (groupRow) TableName() string { return "qixi_crm_b_group_order" }

type activityRow struct {
	ID                uint64 `gorm:"column:id"`
	GroupOrderID      uint64 `gorm:"column:group_order_id"`
	ActivityType      int    `gorm:"column:activity_type"`
	ActivityID        uint64 `gorm:"column:activity_id"`
	RelatedActivityID uint64 `gorm:"column:related_activity_id"`
	Quantity          int    `gorm:"column:quantity"`
	Status            string `gorm:"column:status"`
}

func (activityRow) TableName() string { return "qixi_crm_b_order_activity" }

// ExpireUnpaid closes no more than batch unpaid groups created before now-TTL.
// SKIP LOCKED lets multiple job instances safely share the same queue.
func ExpireUnpaid(ctx context.Context, db *gorm.DB, now time.Time, ttl time.Duration, batch int) (int, error) {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	if batch <= 0 {
		batch = 50
	}
	cutoff := now.UTC().Add(-ttl)
	expired := 0
	for expired < batch {
		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var group groupRow
			err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
				Where("pay_status = ? AND created_at <= ?", "pending", cutoff).
				Order("id ASC").Limit(1).First(&group).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errNoExpiredPendingOrder
			}
			if err != nil {
				return err
			}
			return closePendingGroup(tx, group)
		})
		if errors.Is(err, errNoExpiredPendingOrder) {
			return expired, nil
		}
		if err != nil {
			return expired, err
		}
		expired++
	}
	return expired, nil
}

// closePendingGroup shares the state contract used by api-business cancellation:
// one transaction closes the order and releases every pending-only reservation.
func closePendingGroup(tx *gorm.DB, group groupRow) error {
	closed := tx.Model(&groupRow{}).Where("id = ? AND pay_status = ?", group.ID, "pending").Update("pay_status", "closed")
	if closed.Error != nil {
		return closed.Error
	}
	if closed.RowsAffected != 1 {
		return nil // another worker or payment callback has already advanced it
	}
	if err := tx.Table("qixi_crm_b_order").Where("group_order_id = ? AND status = ?", group.ID, "pending_pay").Update("status", "cancelled").Error; err != nil {
		return err
	}
	if err := releaseActivity(tx, group); err != nil {
		return err
	}
	if group.ActivityType == reservationActivityType {
		if err := releaseReservation(tx, group.ID); err != nil {
			return err
		}
	}
	if err := tx.Table("qixi_crm_b_payment_transaction").Where("group_order_id = ? AND status IN ?", group.ID, []string{"created", "processing", "failed"}).Update("status", "closed").Error; err != nil {
		return err
	}
	return tx.Table("qixi_crm_b_coupon_user").
		Where("user_id = ? AND used_order_id = ? AND status = ?", group.UserID, group.ID, "locked").
		Updates(map[string]any{"status": "unused", "used_order_id": nil}).Error
}

func releaseReservation(tx *gorm.DB, groupID uint64) error {
	var bookings []struct {
		ID     uint64 `gorm:"column:id"`
		SlotID uint64 `gorm:"column:slot_id"`
	}
	if err := tx.Table("qixi_crm_b_reservation_booking AS b").Select("b.id,b.slot_id").
		Joins("JOIN qixi_crm_b_order AS o ON o.id = b.order_id").
		Where("o.group_order_id = ? AND b.status = 1", groupID).Find(&bookings).Error; err != nil {
		return err
	}
	if len(bookings) == 0 {
		return nil
	}
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
	return nil
}

func releaseActivity(tx *gorm.DB, group groupRow) error {
	var activity activityRow
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("group_order_id = ?", group.ID).First(&activity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil || activity.Status != "reserved" {
		return err
	}
	switch activity.ActivityType {
	case presellActivityType:
		if err := tx.Table("qixi_crm_b_presell").Where("product_presell_id = ? AND is_del = 0", activity.ActivityID).UpdateColumn("stock", gorm.Expr("stock + ?", activity.Quantity)).Error; err != nil {
			return err
		}
	case assistActivityType:
		if err := tx.Table("qixi_crm_b_assist").Where("product_assist_id = ? AND is_del = 0", activity.RelatedActivityID).UpdateColumn("stock", gorm.Expr("stock + ?", activity.Quantity)).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_assist_set").Where("product_assist_set_id = ? AND status = ? AND is_del = 0", activity.ActivityID, 11).Update("status", 10).Error; err != nil {
			return err
		}
	case combinationActivityType:
		if err := releaseCombination(tx, group.ID); err != nil {
			return err
		}
	}
	return tx.Model(&activityRow{}).Where("id = ? AND status = ?", activity.ID, "reserved").Update("status", "released").Error
}

func releaseCombination(tx *gorm.DB, groupID uint64) error {
	var member struct {
		ID       uint64 `gorm:"column:id"`
		BuyingID uint64 `gorm:"column:group_buying_id"`
		Status   int    `gorm:"column:status"`
	}
	err := tx.Table("qixi_crm_b_combination_member AS m").Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("m.id,m.group_buying_id,m.status").Joins("JOIN qixi_crm_b_order AS o ON o.id = m.order_id").
		Where("o.group_order_id = ? AND m.is_del = 0", groupID).First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil || member.Status != 0 {
		return err
	}
	if err := tx.Table("qixi_crm_b_combination_member").Where("id = ? AND status = 0 AND is_del = 0", member.ID).Update("is_del", 1).Error; err != nil {
		return err
	}
	var left int64
	if err := tx.Table("qixi_crm_b_combination_member").Where("group_buying_id = ? AND is_del = 0", member.BuyingID).Count(&left).Error; err != nil {
		return err
	}
	if left == 0 {
		return tx.Table("qixi_crm_b_combination_buying").Where("group_buying_id = ? AND status = 0", member.BuyingID).Updates(map[string]any{"status": -1, "is_del": 1}).Error
	}
	return nil
}
