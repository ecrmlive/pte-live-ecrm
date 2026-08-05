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

const pointsOrderActivityType = 20

var (
	ErrOrderOwnership     = errors.New("订单不存在或无权访问")
	ErrPayChannel         = errors.New("暂不支持该支付方式")
	ErrOrderNotPayable    = errors.New("当前订单不可支付")
	ErrOrderNotReceivable = errors.New("订单当前状态不可确认收货")
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
		if groupPayStatus(group) == "paid" {
			return nil
		}
		if group.ActivityType == pointsOrderActivityType {
			return ErrOrderNotPayable
		}
		if groupPayStatus(group) != "pending" {
			return ErrOrderNotPayable
		}
		if err := merchantstock.ReservationsReady(tx, group.ID); err != nil {
			return err
		}
		if err := tx.Where("group_order_id = ? AND channel = ?", group.ID, "mock").FirstOrCreate(&paymentRow{GroupOrderID: group.ID, Channel: "mock", TransactionNo: orderNo("P"), Amount: group.PayAmount, Status: "created"}).Error; err != nil {
			return err
		}
		now := time.Now().UTC()
		updated := tx.Model(&groupRow{}).Where("id = ? AND pay_status = ?", group.ID, "pending").Updates(map[string]any{"pay_status": "paid", "pay_channel": "mock", "paid_at": now})
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
		if err := issueVerificationsForPaidGroup(tx, group.ID); err != nil {
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
	ID                    uint64     `gorm:"column:id"`
	GroupOrderID          uint64     `gorm:"column:group_order_id"`
	Channel               string     `gorm:"column:channel"`
	TransactionNo         string     `gorm:"column:transaction_no"`
	Amount                float64    `gorm:"column:amount"`
	Status                string     `gorm:"column:status"`
	ProviderTransactionNo string     `gorm:"column:provider_transaction_no"`
	ProviderPayload       []byte     `gorm:"column:provider_payload"`
	PaidAt                *time.Time `gorm:"column:paid_at"`
}

func (paymentRow) TableName() string { return "qixi_crm_b_payment_transaction" }
