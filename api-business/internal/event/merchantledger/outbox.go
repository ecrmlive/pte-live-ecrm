// Package merchantledger reliably delivers order settlement facts to the
// merchant-owned ledger. The business service never writes merchant finance
// tables directly.
package merchantledger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const CommandSubject = "qixi.business.merchant-settlement-command.v1"

type Command struct {
	Action         string  `json:"action"`
	OrderID        uint64  `json:"order_id"`
	RefundID       uint64  `json:"refund_id,omitempty"`
	StoreID        uint64  `json:"store_id"`
	MerchantID     uint64  `json:"merchant_id"`
	Amount         float64 `json:"amount"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type Result struct {
	OrderID uint64 `json:"order_id"`
	Status  string `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
}

type outboxRow struct {
	ID             uint64  `gorm:"column:id"`
	Action         string  `gorm:"column:action"`
	OrderID        uint64  `gorm:"column:order_id"`
	RefundID       uint64  `gorm:"column:refund_id"`
	StoreID        uint64  `gorm:"column:store_id"`
	MerchantID     uint64  `gorm:"column:merchant_id"`
	Amount         float64 `gorm:"column:amount"`
	IdempotencyKey string  `gorm:"column:idempotency_key"`
}

func (outboxRow) TableName() string { return "qixi_crm_b_settlement_command_outbox" }

// EnqueueAccrual records a completed order's server-owned payable amount in
// the same transaction that changes the order to completed.
func EnqueueAccrual(tx *gorm.DB, orderID, storeID, merchantID uint64, amount float64) error {
	return enqueue(tx, Command{Action: "accrue", OrderID: orderID, StoreID: storeID, MerchantID: merchantID, Amount: amount, IdempotencyKey: fmt.Sprintf("settlement:accrue:%d", orderID)})
}

// EnqueueReversal records a money-success refund. refundID is deliberately in
// the key, so provider callback retries cannot create a second adjustment.
func EnqueueReversal(tx *gorm.DB, orderID, refundID, storeID, merchantID uint64, amount float64) error {
	return enqueue(tx, Command{Action: "reverse", OrderID: orderID, RefundID: refundID, StoreID: storeID, MerchantID: merchantID, Amount: amount, IdempotencyKey: fmt.Sprintf("settlement:reverse:%d", refundID)})
}

func enqueue(tx *gorm.DB, in Command) error {
	if tx == nil || !valid(in) {
		return gorm.ErrInvalidData
	}
	return tx.Create(&outboxRow{Action: in.Action, OrderID: in.OrderID, RefundID: in.RefundID, StoreID: in.StoreID, MerchantID: in.MerchantID, Amount: in.Amount, IdempotencyKey: in.IdempotencyKey}).Error
}

func valid(in Command) bool {
	if in.OrderID == 0 || in.StoreID == 0 || in.MerchantID == 0 || in.Amount <= 0 || math.IsNaN(in.Amount) || math.IsInf(in.Amount, 0) || len(strings.TrimSpace(in.IdempotencyKey)) < 8 {
		return false
	}
	switch in.Action {
	case "accrue":
		return in.RefundID == 0 && in.IdempotencyKey == fmt.Sprintf("settlement:accrue:%d", in.OrderID)
	case "reverse":
		return in.RefundID != 0 && in.IdempotencyKey == fmt.Sprintf("settlement:reverse:%d", in.RefundID)
	default:
		return false
	}
}

// StartOutboxDispatcher retries transport failures. A terminal merchant
// validation error is retained as failed for reconciliation; it never makes a
// business order appear settled.
func StartOutboxDispatcher(ctx context.Context, businessDB *gorm.DB, natsURL string) {
	if businessDB == nil || strings.TrimSpace(natsURL) == "" {
		return
	}
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				nc, err := nats.Connect(natsURL, nats.Timeout(2*time.Second))
				if err != nil {
					continue
				}
				if err := dispatchPending(ctx, businessDB, nc, 50); err != nil {
					log.Printf("merchant settlement command dispatch: %v", err)
				}
				nc.Close()
			}
		}
	}()
}

func dispatchPending(ctx context.Context, businessDB *gorm.DB, nc *nats.Conn, limit int) error {
	var rows []outboxRow
	if err := businessDB.WithContext(ctx).Where("status = ?", "pending").Order("id ASC").Limit(limit).Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		body, err := json.Marshal(Command{Action: row.Action, OrderID: row.OrderID, RefundID: row.RefundID, StoreID: row.StoreID, MerchantID: row.MerchantID, Amount: row.Amount, IdempotencyKey: row.IdempotencyKey})
		if err != nil {
			return err
		}
		reply, err := nc.Request(CommandSubject, body, 2*time.Second)
		if err != nil {
			_ = businessDB.WithContext(ctx).Table("qixi_crm_b_settlement_command_outbox").Where("id = ? AND status = ?", row.ID, "pending").Updates(map[string]any{"attempts": gorm.Expr("attempts + 1"), "last_error": "结算命令投递暂不可用"}).Error
			continue
		}
		var result Result
		if err := json.Unmarshal(reply.Data, &result); err != nil {
			return err
		}
		updates := map[string]any{"attempts": gorm.Expr("attempts + 1"), "processed_at": time.Now().UTC()}
		if result.Code == "" {
			updates["status"], updates["last_error"] = "accepted", ""
		} else {
			updates["status"], updates["last_error"] = "failed", result.Code
		}
		if err := businessDB.WithContext(ctx).Table("qixi_crm_b_settlement_command_outbox").Where("id = ? AND status = ?", row.ID, "pending").Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}
