// Package merchantdiy publishes non-sensitive store DIY page changes to the
// business read model. The document is public page content, never credentials.
package merchantdiy

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const Subject = "qixi.merchant.diy-page.v1"

const (
	Upserted = "merchant.diy_page.upsert"
	Deleted  = "merchant.diy_page.deleted"
)

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

// StartMerchantOutboxDispatcher leaves an event pending until NATS accepts it.
// The business projection is idempotent, so retry and duplicate delivery are safe.
func StartMerchantOutboxDispatcher(ctx context.Context, merchantDB *gorm.DB, natsURL string) {
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
					_ = dispatchPending(ctx, merchantDB, nc, 50)
					nc.Close()
				}
			}
		}
	}()
}

func dispatchPending(ctx context.Context, merchantDB *gorm.DB, nc *nats.Conn, limit int) error {
	var rows []outboxRow
	if err := merchantDB.WithContext(ctx).Where("status = ? AND event_type IN ?", "pending", []string{Upserted, Deleted}).Order("id ASC").Limit(limit).Find(&rows).Error; err != nil {
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

func PageAggregateID(id uint) string { return strconv.FormatUint(uint64(id), 10) }
