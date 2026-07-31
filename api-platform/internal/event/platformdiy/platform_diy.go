package platformdiy

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"time"
)

const Subject = "qixi.platform.diy-page.v1"
const (
	Upserted = "platform.diy_page.upsert"
	Deleted  = "platform.diy_page.deleted"
)

type envelope struct {
	EventID   uint64          `json:"event_id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}
type row struct {
	ID        uint64          `gorm:"column:id"`
	EventType string          `gorm:"column:event_type"`
	Payload   json.RawMessage `gorm:"column:payload"`
}

func (row) TableName() string { return "qixi_crm_a_outbox" }
func Start(ctx context.Context, db *gorm.DB, url string) {
	if url == "" {
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
				nc, err := nats.Connect(url, nats.Timeout(2*time.Second))
				if err == nil {
					_ = dispatch(ctx, db, nc)
					nc.Close()
				}
			}
		}
	}()
}
func dispatch(ctx context.Context, db *gorm.DB, nc *nats.Conn) error {
	var rows []row
	if err := db.WithContext(ctx).Where("status = ? AND event_type IN ?", "pending", []string{Upserted, Deleted}).Order("id ASC").Limit(50).Find(&rows).Error; err != nil {
		return err
	}
	for _, item := range rows {
		wire, err := json.Marshal(envelope{EventID: item.ID, EventType: item.EventType, Payload: item.Payload})
		if err != nil {
			return err
		}
		if err = nc.Publish(Subject, wire); err != nil {
			return err
		}
		if err = nc.FlushTimeout(2 * time.Second); err != nil {
			return err
		}
		if err = db.WithContext(ctx).Table("qixi_crm_a_outbox").Where("id = ? AND status = ?", item.ID, "pending").Updates(map[string]any{"status": "published", "published_at": time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}
