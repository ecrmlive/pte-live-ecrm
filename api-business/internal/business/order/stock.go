package order

import (
	merchantledger "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantledger"
	merchantstock "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantstock"
	"gorm.io/gorm"
)

// enqueueStockActionForGroup derives stock commands exclusively from immutable
// order-item facts. Callers never accept SKU IDs or quantities from requests.
func enqueueStockActionForGroup(tx *gorm.DB, action string, groupOrderID uint64) error {
	var lines []struct {
		OrderID       uint64 `gorm:"column:order_id"`
		StoreID       uint64 `gorm:"column:store_id"`
		MerchantSKUID uint64 `gorm:"column:merchant_sku_id"`
		Quantity      int    `gorm:"column:quantity"`
	}
	if err := tx.Table("qixi_crm_b_order_item AS i").Select("i.order_id,o.store_id,i.merchant_sku_id,i.quantity").Joins("JOIN qixi_crm_b_order AS o ON o.id = i.order_id").Where("o.group_order_id = ?", groupOrderID).Find(&lines).Error; err != nil {
		return err
	}
	for _, line := range lines {
		if err := merchantstock.EnqueueAction(tx, action, line.OrderID, line.StoreID, line.MerchantSKUID, line.Quantity, nil); err != nil {
			return err
		}
	}
	return nil
}

// enqueueStockRestockForRefund derives returned quantity from the immutable
// after-sale item rows. It never uses browser-provided SKU or quantity data.
func enqueueStockRestockForRefund(tx *gorm.DB, refundID uint64) error {
	var lines []struct {
		OrderID       uint64 `gorm:"column:order_id"`
		StoreID       uint64 `gorm:"column:store_id"`
		MerchantSKUID uint64 `gorm:"column:merchant_sku_id"`
		Quantity      int    `gorm:"column:quantity"`
	}
	if err := tx.Table("qixi_crm_b_aftersale_item AS a").Select("i.order_id,o.store_id,i.merchant_sku_id,a.quantity").Joins("JOIN qixi_crm_b_order_item AS i ON i.id = a.order_item_id").Joins("JOIN qixi_crm_b_order AS o ON o.id = i.order_id").Where("a.refund_id = ?", refundID).Find(&lines).Error; err != nil {
		return err
	}
	if len(lines) == 0 {
		return gorm.ErrRecordNotFound
	}
	for _, line := range lines {
		if err := merchantstock.EnqueueAction(tx, "restock", line.OrderID, line.StoreID, line.MerchantSKUID, line.Quantity, nil); err != nil {
			return err
		}
	}
	return nil
}

// enqueueSettlementReversalForRefund obtains store, merchant and amount from
// locked business facts. The refund callback never accepts any of these money
// fields from a browser or payment-provider payload.
func enqueueSettlementReversalForRefund(tx *gorm.DB, refundID uint64) error {
	var row struct {
		OrderID    uint64  `gorm:"column:order_id"`
		StoreID    uint64  `gorm:"column:store_id"`
		MerchantID uint64  `gorm:"column:merchant_id"`
		Amount     float64 `gorm:"column:amount"`
	}
	if err := tx.Table("qixi_crm_b_refund AS r").Select("r.order_id,o.store_id,o.merchant_id,r.amount").Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").Where("r.id = ?", refundID).Scan(&row).Error; err != nil {
		return err
	}
	if row.OrderID == 0 {
		return gorm.ErrRecordNotFound
	}
	return merchantledger.EnqueueReversal(tx, row.OrderID, refundID, row.StoreID, row.MerchantID, row.Amount)
}

func enqueueSettlementAccrual(tx *gorm.DB, order orderRow) error {
	return merchantledger.EnqueueAccrual(tx, order.ID, order.StoreID, order.MerchantID, order.PayAmount)
}
