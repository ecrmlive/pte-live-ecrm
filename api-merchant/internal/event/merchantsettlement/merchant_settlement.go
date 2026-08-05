// Package merchantsettlement publishes non-sensitive merchant settlement
// projection events. Settlement facts remain in qixi_crm_m_settlement_bill;
// platform receives only a read-model payload through NATS.
package merchantsettlement

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const (
	Subject  = "qixi.merchant.settlement.v1"
	Upserted = "merchant.settlement.upsert"
)

type Payload struct {
	SettlementID uint64    `json:"settlement_id"`
	MerchantID   uint64    `json:"merchant_id"`
	StoreID      uint64    `json:"store_id"`
	MerchantName string    `json:"merchant_name"`
	RegionID     *uint64   `json:"region_id,omitempty"`
	PeriodStart  time.Time `json:"period_start"`
	PeriodEnd    time.Time `json:"period_end"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Envelope struct {
	EventID   uint64          `json:"event_id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

type outboxRow struct {
	ID        uint64          `gorm:"column:id"`
	EventType string          `gorm:"column:event_type"`
	Payload   json.RawMessage `gorm:"column:payload"`
}

func (outboxRow) TableName() string { return "qixi_crm_m_outbox" }

func ValidStatus(status string) bool {
	switch status {
	case "bill_pending", "bill_frozen", "withdraw_applied", "approved", "paid", "rejected":
		return true
	default:
		return false
	}
}

func ValidPayload(p Payload) bool {
	return p.SettlementID != 0 && p.MerchantID != 0 && p.StoreID != 0 && strings.TrimSpace(p.MerchantName) != "" &&
		!math.IsNaN(p.Amount) && !math.IsInf(p.Amount, 0) && !p.PeriodStart.IsZero() && !p.PeriodEnd.IsZero() && !p.PeriodEnd.Before(p.PeriodStart) &&
		!p.UpdatedAt.IsZero() && ValidStatus(p.Status)
}

// Enqueue must be called in the same merchant DB transaction as the status
// change. It contains no bank account, contact number, token or payout data.
func Enqueue(ctx context.Context, tx *gorm.DB, payload Payload) error {
	if tx == nil || !ValidPayload(payload) {
		return gorm.ErrInvalidData
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return tx.WithContext(ctx).Table("qixi_crm_m_outbox").Create(map[string]any{
		"event_type": Upserted, "aggregate_type": "settlement_bill", "aggregate_id": payload.SettlementID, "payload": body,
	}).Error
}

// StartOutboxDispatcher leaves events pending until NATS confirms acceptance.
// Platform projection is an idempotent upsert, so duplicate delivery is safe.
func StartOutboxDispatcher(ctx context.Context, merchantDB *gorm.DB, natsURL string) {
	if merchantDB == nil || strings.TrimSpace(natsURL) == "" {
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
				if err == nil {
					_ = dispatchPending(ctx, merchantDB, nc, 50)
					nc.Close()
				}
			}
		}
	}()
}

func dispatchPending(ctx context.Context, merchantDB *gorm.DB, nc *nats.Conn, limit int) error {
	var rows []outboxRow
	if err := merchantDB.WithContext(ctx).Where("status = ? AND event_type = ?", "pending", Upserted).Order("id ASC").Limit(limit).Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		wire, err := json.Marshal(Envelope{EventID: row.ID, EventType: row.EventType, Payload: row.Payload})
		if err != nil {
			return err
		}
		if err = nc.Publish(Subject, wire); err != nil {
			return err
		}
		if err = nc.FlushTimeout(2 * time.Second); err != nil {
			return err
		}
		if err = merchantDB.WithContext(ctx).Table("qixi_crm_m_outbox").Where("id = ? AND status = ?", row.ID, "pending").Updates(map[string]any{"status": "published", "published_at": time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}
