// Package merchantstock reliably delivers business order inventory commands to
// the merchant-owned stock command subscriber. It owns no merchant database
// connection and never mutates inventory itself.
package merchantstock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const CommandSubject = "qixi.business.merchant-stock-command.v1"

var (
	ErrReservationsPending = errors.New("库存预留处理中")
	ErrReservationsFailed  = errors.New("库存预留失败")
)

type Command struct {
	Action         string    `json:"action"`
	OrderID        uint64    `json:"order_id"`
	StoreID        uint64    `json:"store_id"`
	MerchantSKUID  uint64    `json:"merchant_sku_id"`
	Quantity       int       `json:"quantity"`
	ExpiresAt      time.Time `json:"expires_at,omitempty"`
	IdempotencyKey string    `json:"idempotency_key"`
}

type Result struct {
	OrderID       uint64 `json:"order_id"`
	MerchantSKUID uint64 `json:"merchant_sku_id"`
	Status        string `json:"status"`
	Code          string `json:"code"`
}

type outboxRow struct {
	ID             uint64     `gorm:"column:id"`
	Action         string     `gorm:"column:action"`
	OrderID        uint64     `gorm:"column:order_id"`
	StoreID        uint64     `gorm:"column:store_id"`
	MerchantSKUID  uint64     `gorm:"column:merchant_sku_id"`
	Quantity       int        `gorm:"column:quantity"`
	ExpiresAt      *time.Time `gorm:"column:expires_at"`
	IdempotencyKey string     `gorm:"column:idempotency_key"`
}

func (outboxRow) TableName() string { return "qixi_crm_b_stock_command_outbox" }

// EnqueueReserve must be called in the same business transaction that creates
// the order line. The deterministic key makes repeat submission safe.
func EnqueueReserve(tx *gorm.DB, orderID, storeID, merchantSKUID uint64, quantity int, expiresAt time.Time) error {
	if expiresAt.IsZero() {
		return gorm.ErrInvalidData
	}
	return EnqueueAction(tx, "reserve", orderID, storeID, merchantSKUID, quantity, &expiresAt)
}

// EnqueueAction appends a deterministic per-line saga command. It is used for
// payment confirmation and cancellation release in the same business
// transaction as their order state changes.
func EnqueueAction(tx *gorm.DB, action string, orderID, storeID, merchantSKUID uint64, quantity int, expiresAt *time.Time) error {
	if tx == nil || orderID == 0 || storeID == 0 || merchantSKUID == 0 || quantity < 1 {
		return gorm.ErrInvalidData
	}
	if action != "reserve" && action != "confirm" && action != "release" && action != "restock" {
		return gorm.ErrInvalidData
	}
	if action == "reserve" && (expiresAt == nil || expiresAt.IsZero()) {
		return gorm.ErrInvalidData
	}
	key := fmt.Sprintf("stock:%s:%d:%d", action, orderID, merchantSKUID)
	return tx.Create(&outboxRow{Action: action, OrderID: orderID, StoreID: storeID, MerchantSKUID: merchantSKUID, Quantity: quantity, ExpiresAt: expiresAt, IdempotencyKey: key}).Error
}

// ReservationsReady is called inside the payment transaction. It compares
// every order line with an accepted reserve command so a missing outbox row
// cannot accidentally be treated as available inventory.
func ReservationsReady(tx *gorm.DB, groupOrderID uint64) error {
	if tx == nil || groupOrderID == 0 {
		return gorm.ErrInvalidData
	}
	var total, accepted, failed int64
	base := tx.Table("qixi_crm_b_order AS o").Joins("JOIN qixi_crm_b_order_item AS i ON i.order_id = o.id").Where("o.group_order_id = ?", groupOrderID)
	if err := base.Count(&total).Error; err != nil {
		return err
	}
	if total == 0 {
		return ErrReservationsFailed
	}
	if err := base.Joins("JOIN qixi_crm_b_stock_command_outbox AS s ON s.order_id = o.id AND s.merchant_sku_id = i.merchant_sku_id").Where("s.action = ? AND s.status = ?", "reserve", "accepted").Count(&accepted).Error; err != nil {
		return err
	}
	if err := base.Joins("JOIN qixi_crm_b_stock_command_outbox AS s ON s.order_id = o.id AND s.merchant_sku_id = i.merchant_sku_id").Where("s.action = ? AND s.status = ?", "reserve", "failed").Count(&failed).Error; err != nil {
		return err
	}
	if failed > 0 {
		return ErrReservationsFailed
	}
	if accepted != total {
		return ErrReservationsPending
	}
	return nil
}

// StartOutboxDispatcher retries transport errors forever. Business state is
// not advanced until merchant replies accepted; terminal inventory conflicts
// are recorded as failed and must be surfaced by the order state machine.
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
					log.Printf("merchant stock command dispatch: %v", err)
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
		cmd := Command{Action: row.Action, OrderID: row.OrderID, StoreID: row.StoreID, MerchantSKUID: row.MerchantSKUID, Quantity: row.Quantity, IdempotencyKey: row.IdempotencyKey}
		if row.ExpiresAt != nil {
			cmd.ExpiresAt = row.ExpiresAt.UTC()
		}
		body, err := json.Marshal(cmd)
		if err != nil {
			return err
		}
		reply, err := nc.Request(CommandSubject, body, 2*time.Second)
		if err != nil {
			_ = businessDB.WithContext(ctx).Table("qixi_crm_b_stock_command_outbox").Where("id = ? AND status = ?", row.ID, "pending").Updates(map[string]any{"attempts": gorm.Expr("attempts + 1"), "last_error": "库存命令投递暂不可用"}).Error
			continue
		}
		var result Result
		if err := json.Unmarshal(reply.Data, &result); err != nil {
			return err
		}
		updates := map[string]any{"attempts": gorm.Expr("attempts + 1"), "processed_at": time.Now().UTC()}
		if result.Code == "" {
			updates["status"] = "accepted"
			updates["last_error"] = ""
		} else {
			updates["status"] = "failed"
			updates["last_error"] = result.Code
		}
		if err := businessDB.WithContext(ctx).Table("qixi_crm_b_stock_command_outbox").Where("id = ? AND status = ?", row.ID, "pending").Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}
