package merchantintegralpolicy

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const Subject = "qixi.merchant.integral-policy.v1"
const Upserted = "merchant.integral_policy.upsert"

type Payload struct {
	StoreID         uint64 `json:"store_id"`
	Enabled         bool   `json:"enabled"`
	PointsPerYuan   int64  `json:"points_per_yuan"`
	MaxDeductionBps int64  `json:"max_deduction_bps"`
}
type envelope struct {
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

func Enqueue(tx *gorm.DB, payload Payload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return tx.Table("qixi_crm_m_outbox").Create(map[string]any{"event_type": Upserted, "aggregate_type": "integral_policy", "aggregate_id": payload.StoreID, "payload": body}).Error
}
func StartMerchantOutboxDispatcher(ctx context.Context, db *gorm.DB, natsURL string) {
	if natsURL == "" {
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
					_ = dispatch(ctx, db, nc)
					nc.Close()
				}
			}
		}
	}()
}
func dispatch(ctx context.Context, db *gorm.DB, nc *nats.Conn) error {
	var rows []outboxRow
	if err := db.WithContext(ctx).Where("status = ? AND event_type = ?", "pending", Upserted).Order("id ASC").Limit(50).Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		wire, err := json.Marshal(envelope{EventID: row.ID, EventType: row.EventType, Payload: row.Payload})
		if err != nil {
			return err
		}
		if err = nc.Publish(Subject, wire); err != nil {
			return err
		}
		if err = nc.FlushTimeout(2 * time.Second); err != nil {
			return err
		}
		if err = db.WithContext(ctx).Table("qixi_crm_m_outbox").Where("id = ? AND status = ?", row.ID, "pending").Updates(map[string]any{"status": "published", "published_at": time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}
