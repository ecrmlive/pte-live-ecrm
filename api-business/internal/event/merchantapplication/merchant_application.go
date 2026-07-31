// Package merchantapplication publishes C-end merchant application outbox
// events. It uses the business outbox so a successful HTTP response is never
// lost merely because NATS is temporarily unavailable.
package merchantapplication

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const Subject = "qixi.business.merchant-application.v1"
const Submitted = "business.merchant_application.submitted"

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

func StartOutboxDispatcher(ctx context.Context, businessDB *gorm.DB, natsURL string) {
	if natsURL == "" {
		log.Printf("merchant application outbox disabled: nats.url is empty")
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
					log.Printf("merchant application outbox: %v", err)
				}
				nc.Close()
			}
		}
	}()
}

func dispatchPending(ctx context.Context, db *gorm.DB, nc *nats.Conn, limit int) error {
	var rows []outboxRow
	if err := db.WithContext(ctx).Table("qixi_crm_b_outbox").Where("status = ? AND event_type = ?", "pending", Submitted).Order("id ASC").Limit(limit).Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		wire, err := json.Marshal(Envelope{EventID: row.ID, EventType: row.EventType, Payload: row.Payload})
		if err != nil {
			return err
		}
		if err := nc.Publish(Subject, wire); err != nil {
			return err
		}
		if err := nc.FlushTimeout(2 * time.Second); err != nil {
			return err
		}
		if err := db.WithContext(ctx).Table("qixi_crm_b_outbox").Where("id = ? AND status = ?", row.ID, "pending").Updates(map[string]any{"status": "published", "published_at": time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}
