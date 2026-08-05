package order

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	assistOrderActivityType      = 3
	combinationOrderActivityType = 4
	presellOrderActivityType     = 2
)

type orderActivityRow struct {
	ID                uint64 `gorm:"column:id"`
	GroupOrderID      uint64 `gorm:"column:group_order_id"`
	ActivityType      int    `gorm:"column:activity_type"`
	ActivityID        uint64 `gorm:"column:activity_id"`
	RelatedActivityID uint64 `gorm:"column:related_activity_id"`
	Quantity          int    `gorm:"column:quantity"`
	Status            string `gorm:"column:status"`
}

func (orderActivityRow) TableName() string { return "qixi_crm_b_order_activity" }

// settleOrderActivity runs inside the business payment transaction. Every
// special-order transition is driven from the immutable order activity link,
// rather than from display text or legacy qixi_m_* order rows.
func settleOrderActivity(tx *gorm.DB, group groupRow) error {
	var activity orderActivityRow
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("group_order_id = ?", group.ID).First(&activity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if activity.Status == "paid" {
		return nil
	}
	if activity.Status != "reserved" {
		return ErrOrderNotPayable
	}

	switch activity.ActivityType {
	case presellOrderActivityType:
		if err := settlePresellActivity(tx, group, activity); err != nil {
			return err
		}
	case combinationOrderActivityType:
		if err := settleCombinationActivity(tx, group); err != nil {
			return err
		}
	case assistOrderActivityType:
		updated := tx.Table("qixi_crm_b_assist_set").
			Where("product_assist_set_id = ? AND status = ? AND is_del = 0", activity.ActivityID, 11).
			Update("status", 20)
		if updated.Error != nil {
			return updated.Error
		}
		if updated.RowsAffected != 1 {
			return ErrOrderNotPayable
		}
		if err := tx.Table("qixi_crm_b_assist").
			Where("product_assist_id = ? AND is_del = 0", activity.RelatedActivityID).
			UpdateColumn("pay_count", gorm.Expr("pay_count + 1")).Error; err != nil {
			return err
		}
	}
	return tx.Model(&orderActivityRow{}).Where("id = ? AND status = ?", activity.ID, "reserved").Update("status", "paid").Error
}

func settlePresellActivity(tx *gorm.DB, group groupRow, activity orderActivityRow) error {
	var order struct {
		ID         uint64 `gorm:"column:id"`
		MerchantID uint64 `gorm:"column:merchant_id"`
	}
	if err := tx.Table("qixi_crm_b_order").Select("id,merchant_id").Where("group_order_id = ?", group.ID).First(&order).Error; err != nil {
		return err
	}
	var p struct {
		PresellType int     `gorm:"column:presell_type"`
		FinalPrice  float64 `gorm:"column:final_price"`
		FinalStart  string  `gorm:"column:final_start_time"`
		FinalEnd    string  `gorm:"column:final_end_time"`
	}
	if err := tx.Table("qixi_crm_b_presell").Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_presell_id = ? AND is_del = 0", activity.ActivityID).First(&p).Error; err != nil {
		return err
	}
	if err := tx.Table("qixi_crm_b_presell").Where("product_presell_id = ?", activity.ActivityID).UpdateColumn("pay_count", gorm.Expr("pay_count + 1")).Error; err != nil {
		return err
	}
	if p.PresellType == 1 {
		return tx.Table("qixi_crm_b_order").Where("id = ? AND status = ?", order.ID, "paid").Update("status", "fulfilling").Error
	}
	if p.PresellType != 2 || p.FinalPrice <= 0 {
		return ErrOrderNotPayable
	}
	start, end := parsePresellWindow(p.FinalStart, p.FinalEnd)
	if !end.After(start) {
		return ErrOrderNotPayable
	}
	if err := tx.Table("qixi_crm_b_presell_order").Create(map[string]any{"presell_order_sn": fmt.Sprintf("PF%d", order.ID), "uid": group.UserID, "mer_id": order.MerchantID, "order_id": order.ID, "product_presell_id": activity.ActivityID, "final_start_time": start, "final_end_time": end, "paid": 0, "status": 1, "pay_type": 0, "pay_price": p.FinalPrice * float64(activity.Quantity)}).Error; err != nil {
		return err
	}
	return tx.Table("qixi_crm_b_order").Where("id = ? AND status = ?", order.ID, "paid").Update("status", "awaiting_final").Error
}

func parsePresellWindow(startS, endS string) (time.Time, time.Time) {
	now := time.Now()
	start := now
	end := now.AddDate(0, 0, 30)
	for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339} {
		if t, e := time.ParseInLocation(layout, startS, time.Local); e == nil {
			start = t
			break
		}
	}
	for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339} {
		if t, e := time.ParseInLocation(layout, endS, time.Local); e == nil {
			end = t
			break
		}
	}
	return start, end
}

func settleCombinationActivity(tx *gorm.DB, group groupRow) error {
	var member struct {
		ID       uint64 `gorm:"column:id"`
		BuyingID uint64 `gorm:"column:group_buying_id"`
		Status   int    `gorm:"column:status"`
	}
	err := tx.Table("qixi_crm_b_combination_member AS m").Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("m.id,m.group_buying_id,m.status").Joins("JOIN qixi_crm_b_order AS o ON o.id = m.order_id").
		Where("o.group_order_id = ? AND m.is_del = 0", group.ID).First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrOrderNotPayable
	}
	if err != nil {
		return err
	}
	if member.Status == 1 {
		return nil
	}
	var buying struct {
		ID             uint64 `gorm:"column:group_buying_id"`
		Status         int    `gorm:"column:status"`
		IsDel          int    `gorm:"column:is_del"`
		BuyingCountNum int    `gorm:"column:buying_count_num"`
		YetBuyingNum   int    `gorm:"column:yet_buying_num"`
		EndTime        int64  `gorm:"column:end_time"`
		ProductGroupID uint64 `gorm:"column:product_group_id"`
	}
	if err := tx.Table("qixi_crm_b_combination_buying").Clauses(clause.Locking{Strength: "UPDATE"}).Where("group_buying_id = ?", member.BuyingID).First(&buying).Error; err != nil {
		return err
	}
	if buying.Status != 0 || buying.IsDel != 0 || (buying.EndTime > 0 && buying.EndTime <= time.Now().Unix()) {
		return ErrOrderNotPayable
	}
	if err := tx.Table("qixi_crm_b_combination_member").Where("id = ? AND status = 0 AND is_del = 0", member.ID).Update("status", 1).Error; err != nil {
		return err
	}
	next := buying.YetBuyingNum + 1
	updates := map[string]any{"yet_buying_num": next}
	if next >= buying.BuyingCountNum {
		updates["status"] = 10
	}
	if err := tx.Table("qixi_crm_b_combination_buying").Where("group_buying_id = ? AND status = 0 AND is_del = 0", buying.ID).Updates(updates).Error; err != nil {
		return err
	}
	if err := tx.Table("qixi_crm_b_combination_group").Where("product_group_id = ? AND is_del = 0", buying.ProductGroupID).UpdateColumn("pay_count", gorm.Expr("pay_count + 1")).Error; err != nil {
		return err
	}
	if next < buying.BuyingCountNum {
		return nil
	}
	if err := tx.Table("qixi_crm_b_combination_group").Where("product_group_id = ? AND is_del = 0", buying.ProductGroupID).UpdateColumn("success_num", gorm.Expr("success_num + 1")).Error; err != nil {
		return err
	}
	return tx.Table("qixi_crm_b_order AS o").Joins("JOIN qixi_crm_b_combination_member AS m ON m.order_id = o.id").Where("m.group_buying_id = ? AND m.status = 1 AND m.is_del = 0 AND o.status = ?", buying.ID, "paid").Update("o.status", "fulfilling").Error
}

func releaseCombinationActivity(tx *gorm.DB, group groupRow) error {
	var member struct {
		ID       uint64 `gorm:"column:id"`
		BuyingID uint64 `gorm:"column:group_buying_id"`
		Status   int    `gorm:"column:status"`
	}
	err := tx.Table("qixi_crm_b_combination_member AS m").Clauses(clause.Locking{Strength: "UPDATE"}).Select("m.id,m.group_buying_id,m.status").Joins("JOIN qixi_crm_b_order AS o ON o.id = m.order_id").Where("o.group_order_id = ? AND m.is_del = 0", group.ID).First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if member.Status != 0 {
		return nil
	}
	if err := tx.Table("qixi_crm_b_combination_member").Where("id = ? AND status = 0 AND is_del = 0", member.ID).Update("is_del", 1).Error; err != nil {
		return err
	}
	var left int64
	if err := tx.Table("qixi_crm_b_combination_member").Where("group_buying_id = ? AND is_del = 0", member.BuyingID).Count(&left).Error; err != nil {
		return err
	}
	if left == 0 {
		if err := tx.Table("qixi_crm_b_combination_buying").Where("group_buying_id = ? AND status = 0", member.BuyingID).Updates(map[string]any{"status": -1, "is_del": 1}).Error; err != nil {
			return err
		}
	}
	return nil
}

// releaseOrderActivity reverses a pending-only reservation during cancellation.
// Paid orders deliberately never enter this path.
func releaseOrderActivity(tx *gorm.DB, group groupRow) error {
	var activity orderActivityRow
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("group_order_id = ?", group.ID).First(&activity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if activity.Status != "reserved" {
		return nil
	}
	switch activity.ActivityType {
	case presellOrderActivityType:
		if err := tx.Table("qixi_crm_b_presell").Where("product_presell_id = ? AND is_del = 0", activity.ActivityID).UpdateColumn("stock", gorm.Expr("stock + ?", activity.Quantity)).Error; err != nil {
			return err
		}
	case combinationOrderActivityType:
		if err := releaseCombinationActivity(tx, group); err != nil {
			return err
		}
	case assistOrderActivityType:
		if err := tx.Table("qixi_crm_b_assist").
			Where("product_assist_id = ? AND is_del = 0", activity.RelatedActivityID).
			UpdateColumn("stock", gorm.Expr("stock + ?", activity.Quantity)).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_assist_set").
			Where("product_assist_set_id = ? AND status = ? AND is_del = 0", activity.ActivityID, 11).
			Update("status", 10).Error; err != nil {
			return err
		}
	}
	return tx.Model(&orderActivityRow{}).Where("id = ? AND status = ?", activity.ID, "reserved").Update("status", "released").Error
}
