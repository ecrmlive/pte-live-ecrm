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
	if err := db.WithContext(ctx).Where("id = ? AND user_id = ?", groupOrderID, userID).First(&group).Error; err != nil {
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
		for _, item := range rows {
			items[item.OrderID] = append(items[item.OrderID], item)
		}
	}
	return GroupDetail{Group: group, Orders: orders, Items: items}, nil
}

func ListGroups(ctx context.Context, db *gorm.DB, userID uint64, page, limit int) ([]groupRow, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	query := db.WithContext(ctx).Model(&groupRow{}).Where("user_id = ?", userID)
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
	return order, items, nil
}
