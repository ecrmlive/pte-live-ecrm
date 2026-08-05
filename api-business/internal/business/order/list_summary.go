package order

import (
	"context"

	"gorm.io/gorm"
)

type groupListSummary struct {
	FulfillmentStatus  string
	HasUncommentedItem bool
}

// GroupFulfillmentStatus returns a conservative, user-facing summary for a
// grouped payment. A group can contain multiple merchant orders; active work
// always wins over a completed child so the user does not miss fulfillment.
func GroupFulfillmentStatus(group groupRow, orders []orderRow) string {
	if group.PayStatus != "paid" {
		return group.PayStatus
	}
	if len(orders) == 0 {
		return "paid"
	}
	seen := make(map[string]bool, len(orders))
	for _, order := range orders {
		seen[order.Status] = true
	}
	if seen["awaiting_final"] || seen["final_timeout"] {
		if seen["awaiting_final"] {
			return "awaiting_final"
		}
		return "final_timeout"
	}
	if seen["paid"] || seen["fulfilling"] {
		return "fulfilling"
	}
	if seen["shipped"] {
		return "shipped"
	}
	if seen["aftersale"] {
		return "aftersale"
	}
	if seen["completed"] {
		return "completed"
	}
	if seen["cancelled"] {
		return "cancelled"
	}
	return "paid"
}

func loadGroupListSummaries(ctx context.Context, db *gorm.DB, userID uint64, groups []groupRow) (map[uint64]groupListSummary, error) {
	result := make(map[uint64]groupListSummary, len(groups))
	if len(groups) == 0 {
		return result, nil
	}
	groupIDs := make([]uint64, 0, len(groups))
	groupByID := make(map[uint64]groupRow, len(groups))
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID)
		groupByID[group.ID] = group
	}
	var orders []orderRow
	if err := db.WithContext(ctx).Where("group_order_id IN ? AND user_id = ?", groupIDs, userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	ordersByGroup := make(map[uint64][]orderRow, len(groups))
	completedOrderIDs := make([]uint64, 0)
	for _, order := range orders {
		ordersByGroup[order.GroupOrderID] = append(ordersByGroup[order.GroupOrderID], order)
		if order.Status == "completed" {
			completedOrderIDs = append(completedOrderIDs, order.ID)
		}
	}
	uncommentedByGroup := make(map[uint64]bool)
	if len(completedOrderIDs) > 0 {
		var rows []struct {
			GroupOrderID uint64  `gorm:"column:group_order_id"`
			CommentID    *uint64 `gorm:"column:comment_id"`
		}
		if err := db.WithContext(ctx).Table("qixi_crm_b_order_item AS oi").
			Select("o.group_order_id,pc.id AS comment_id").
			Joins("JOIN qixi_crm_b_order AS o ON o.id = oi.order_id").
			Joins("LEFT JOIN qixi_crm_b_product_comment AS pc ON pc.order_item_id = oi.id AND pc.user_id = ?", userID).
			Where("oi.order_id IN ?", completedOrderIDs).Find(&rows).Error; err != nil {
			return nil, err
		}
		for _, row := range rows {
			if row.CommentID == nil {
				uncommentedByGroup[row.GroupOrderID] = true
			}
		}
	}
	for _, group := range groups {
		result[group.ID] = groupListSummary{
			FulfillmentStatus:  GroupFulfillmentStatus(groupByID[group.ID], ordersByGroup[group.ID]),
			HasUncommentedItem: uncommentedByGroup[group.ID],
		}
	}
	return result, nil
}
