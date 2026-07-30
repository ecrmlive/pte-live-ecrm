package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrOrderOwnership = errors.New("订单不存在或无权访问")
	ErrPayChannel     = errors.New("暂不支持该支付方式")
)

// PayMock is deliberately limited to local/test mock payment. Real channels
// must enter through verified provider callbacks before this state transition.
func PayMock(ctx context.Context, db *gorm.DB, userID, groupOrderID uint64) (CreatedOrder, error) {
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
		if err := tx.Where("group_order_id = ? AND channel = ?", group.ID, "mock").FirstOrCreate(&paymentRow{GroupOrderID: group.ID, Channel: "mock", TransactionNo: orderNo("P"), Amount: group.PayAmount, Status: "created"}).Error; err != nil {
			return err
		}
		if groupPayStatus(group) == "paid" {
			return nil
		}
		now := time.Now().UTC()
		if err := tx.Model(&groupRow{}).Where("id = ? AND pay_status = ?", group.ID, "pending").Updates(map[string]any{"pay_status": "paid", "pay_channel": "mock", "paid_at": now}).Error; err != nil {
			return err
		}
		if err := tx.Model(&orderRow{}).Where("group_order_id = ? AND status = ?", group.ID, "pending_pay").Updates(map[string]any{"status": "paid", "paid_at": now}).Error; err != nil {
			return err
		}
		return tx.Model(&paymentRow{}).Where("group_order_id = ? AND channel = ?", group.ID, "mock").Updates(map[string]any{"status": "succeeded", "paid_at": now, "provider_transaction_no": fmt.Sprintf("mock-%d", group.ID)}).Error
	})
	return result, err
}

// Group rows use an explicit status column in the schema. This helper keeps
// compatibility with the small row used by Create until read models are added.
func groupPayStatus(group groupRow) string { return group.PayStatus }

type paymentRow struct {
	GroupOrderID          uint64     `gorm:"column:group_order_id"`
	Channel               string     `gorm:"column:channel"`
	TransactionNo         string     `gorm:"column:transaction_no"`
	Amount                float64    `gorm:"column:amount"`
	Status                string     `gorm:"column:status"`
	ProviderTransactionNo string     `gorm:"column:provider_transaction_no"`
	PaidAt                *time.Time `gorm:"column:paid_at"`
}

func (paymentRow) TableName() string { return "qixi_crm_b_payment_transaction" }
