package order

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type GroupDetail struct {
	Group  groupRow
	Orders []orderRow
	Items  map[uint64][]orderItemRow
}

func GetGroup(ctx context.Context, db *gorm.DB, userID, groupOrderID uint64) (GroupDetail, error) {
	var group groupRow
	if err := db.WithContext(ctx).Where("id = ? AND user_id = ? AND user_archived_at IS NULL", groupOrderID, userID).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GroupDetail{}, ErrOrderOwnership
		}
		return GroupDetail{}, err
	}
	var orders []orderRow
	if err := db.WithContext(ctx).Where("group_order_id = ? AND user_id = ?", group.ID, userID).Order("id ASC").Find(&orders).Error; err != nil {
		return GroupDetail{}, err
	}
	ids := make([]uint64, 0, len(orders))
	for _, order := range orders {
		ids = append(ids, order.ID)
	}
	items := map[uint64][]orderItemRow{}
	if len(ids) > 0 {
		var rows []orderItemRow
		if err := db.WithContext(ctx).Where("order_id IN ?", ids).Order("id ASC").Find(&rows).Error; err != nil {
			return GroupDetail{}, err
		}
		if err := markCommented(ctx, db, userID, rows); err != nil {
			return GroupDetail{}, err
		}
		for _, item := range rows {
			items[item.OrderID] = append(items[item.OrderID], item)
		}
	}
	return GroupDetail{Group: group, Orders: orders, Items: items}, nil
}

func markCommented(ctx context.Context, db *gorm.DB, userID uint64, items []orderItemRow) error {
	if len(items) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	var rows []struct {
		OrderItemID uint64 `gorm:"column:order_item_id"`
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_product_comment").Select("order_item_id").Where("user_id = ? AND order_item_id IN ?", userID, ids).Find(&rows).Error; err != nil {
		return err
	}
	commented := make(map[uint64]bool, len(rows))
	for _, row := range rows {
		commented[row.OrderItemID] = true
	}
	for index := range items {
		items[index].Commented = commented[items[index].ID]
	}
	return nil
}

func ListGroups(ctx context.Context, db *gorm.DB, userID uint64, payStatus, fulfillmentStatus, keyword string, page, limit int) ([]groupRow, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	query := db.WithContext(ctx).Model(&groupRow{}).Where("user_id = ? AND user_archived_at IS NULL", userID)
	if payStatus != "" {
		query = query.Where("pay_status = ?", payStatus)
	}
	switch fulfillmentStatus {
	case "awaiting_fulfillment":
		query = query.Where("pay_status = ? AND EXISTS (SELECT 1 FROM qixi_crm_b_order AS o WHERE o.group_order_id = qixi_crm_b_group_order.id AND o.status IN ?)", "paid", []string{"paid", "fulfilling"})
	case "awaiting_receipt":
		query = query.Where("pay_status = ? AND EXISTS (SELECT 1 FROM qixi_crm_b_order AS o WHERE o.group_order_id = qixi_crm_b_group_order.id AND o.status = ?)", "paid", "shipped")
	case "awaiting_comment":
		query = query.Where("pay_status = ? AND EXISTS (SELECT 1 FROM qixi_crm_b_order AS o JOIN qixi_crm_b_order_item AS oi ON oi.order_id = o.id LEFT JOIN qixi_crm_b_product_comment AS pc ON pc.order_item_id = oi.id AND pc.user_id = ? WHERE o.group_order_id = qixi_crm_b_group_order.id AND o.status = ? AND pc.id IS NULL)", "paid", userID, "completed")
	}
	if keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("(order_no LIKE ? OR EXISTS (SELECT 1 FROM qixi_crm_b_order AS o JOIN qixi_crm_b_order_item AS oi ON oi.order_id = o.id WHERE o.group_order_id = qixi_crm_b_group_order.id AND (o.order_no LIKE ? OR oi.title_snapshot LIKE ?)))", pattern, pattern, pattern)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []groupRow
	err := query.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func GetOrder(ctx context.Context, db *gorm.DB, userID, orderID uint64) (orderRow, []orderItemRow, error) {
	var order orderRow
	if err := db.WithContext(ctx).Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return order, nil, ErrOrderOwnership
		}
		return order, nil, err
	}
	var items []orderItemRow
	if err := db.WithContext(ctx).Where("order_id = ?", order.ID).Order("id ASC").Find(&items).Error; err != nil {
		return order, nil, err
	}
	if err := markCommented(ctx, db, userID, items); err != nil {
		return order, nil, err
	}
	return order, items, nil
}
