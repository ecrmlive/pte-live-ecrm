package merchantapplication

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const ReviewSubject = "qixi.platform.merchant-application.v1"
const Reviewed = "platform.merchant_application.reviewed"

type ReviewPayload struct {
	SourceApplicationID uint   `json:"source_application_id"`
	Status              string `json:"status"`
	ReviewNote          string `json:"review_note"`
}

type reviewOutboxRow struct {
	ID        uint64          `gorm:"column:id"`
	EventType string          `gorm:"column:event_type"`
	Payload   json.RawMessage `gorm:"column:payload"`
}

func (reviewOutboxRow) TableName() string { return "qixi_crm_a_outbox" }

// EnqueueReview writes the outward event in the same platform transaction as
// the audit decision. The business service receives the result only after it
// has been committed, never by directly reading the admin database.
func EnqueueReview(tx *gorm.DB, payload ReviewPayload) error {
	if payload.SourceApplicationID == 0 || (payload.Status != "approved" && payload.Status != "rejected") {
		return nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return tx.Table("qixi_crm_a_outbox").Create(map[string]any{
		"event_type":     Reviewed,
		"aggregate_type": "merchant_application",
		"aggregate_id":   payload.SourceApplicationID,
		"payload":        body,
		"status":         "pending",
	}).Error
}

func StartReviewOutboxDispatcher(ctx context.Context, db *gorm.DB, url string) {
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
					_ = dispatchReviews(ctx, db, nc)
					nc.Close()
				}
			}
		}
	}()
}

func dispatchReviews(ctx context.Context, db *gorm.DB, nc *nats.Conn) error {
	var rows []reviewOutboxRow
	if err := db.WithContext(ctx).Where("status = ? AND event_type = ?", "pending", Reviewed).Order("id ASC").Limit(50).Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		wire, err := json.Marshal(envelope{EventID: row.ID, EventType: row.EventType, Payload: row.Payload})
		if err != nil {
			return err
		}
		if err = nc.Publish(ReviewSubject, wire); err != nil {
			return err
		}
		if err = nc.FlushTimeout(2 * time.Second); err != nil {
			return err
		}
		if err = db.WithContext(ctx).Table("qixi_crm_a_outbox").Where("id = ? AND status = ?", row.ID, "pending").Updates(map[string]any{"status": "published", "published_at": time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}
