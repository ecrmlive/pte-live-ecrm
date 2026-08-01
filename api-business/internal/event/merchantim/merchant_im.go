// Package merchantim synchronizes the current merchant IM SDK AppId into the
// business read model. It never contacts pte-live-im or stores its secrets.
package merchantim

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const Subject = "qixi.merchant.im-sdk-app.v1"

const (
	Activated   = "merchant.im_sdk_app.activated"
	Deactivated = "merchant.im_sdk_app.deactivated"
)

type Payload struct {
	MerchantID   uint64 `json:"merchant_id"`
	SDKAppID     string `json:"sdk_app_id"`
	Status       string `json:"status"`
	IsActive     bool   `json:"is_active"`
	APIPublicURL string `json:"api_public_url"`
	WSPublicURL  string `json:"ws_public_url"`
	PTEProfileID string `json:"pte_profile_id"`
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

// StartMerchantOutboxDispatcher retries NATS connectivity indefinitely. Pending
// events stay in the merchant DB until NATS confirms acceptance; duplicate
// delivery is safe because the business projection is idempotent.
func StartMerchantOutboxDispatcher(ctx context.Context, merchantDB *gorm.DB, natsURL string) {
	if natsURL == "" {
		log.Printf("merchant IM outbox disabled: nats.url is empty")
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
				_ = dispatchPending(ctx, merchantDB, nc, 50)
				nc.Close()
			}
		}
	}()
}

func dispatchPending(ctx context.Context, merchantDB *gorm.DB, nc *nats.Conn, limit int) error {
	var rows []outboxRow
	if err := merchantDB.WithContext(ctx).Where("status = ? AND event_type IN ?", "pending", []string{Activated, Deactivated}).Order("id ASC").Limit(limit).Find(&rows).Error; err != nil {
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
		if err := merchantDB.WithContext(ctx).Table("qixi_crm_m_outbox").Where("id = ? AND status = ?", row.ID, "pending").Updates(map[string]any{"status": "published", "published_at": time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}

// StartBusinessProjection subscribes only to merchant outbox events and builds
// qixi_crm_b_merchant_im_sdk_app_view. No merchant database connection exists here.
func StartBusinessProjection(ctx context.Context, businessDB *gorm.DB, natsURL string) (*nats.Conn, error) {
	if natsURL == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_business_im_projection"))
	if err != nil {
		return nil, err
	}
	_, err = nc.Subscribe(Subject, func(msg *nats.Msg) {
		if err := ApplyProjection(ctx, businessDB, msg.Data); err != nil {
			log.Printf("merchant IM projection failed: %v", err)
		}
	})
	if err != nil {
		nc.Close()
		return nil, err
	}
	return nc, nc.FlushTimeout(2 * time.Second)
}

func ApplyProjection(ctx context.Context, businessDB *gorm.DB, raw []byte) error {
	var envelope Envelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return err
	}
	var payload Payload
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		return err
	}
	if payload.MerchantID == 0 || payload.SDKAppID == "" {
		return nil
	}
	switch envelope.EventType {
	case Activated:
		return businessDB.WithContext(ctx).Exec(`INSERT INTO qixi_crm_b_merchant_im_sdk_app_view
      (merchant_id, sdk_app_id, api_public_url, ws_public_url, pte_profile_id, updated_at)
      VALUES (?, ?, ?, ?, ?, NOW())
      ON DUPLICATE KEY UPDATE sdk_app_id=VALUES(sdk_app_id), api_public_url=VALUES(api_public_url),
      ws_public_url=VALUES(ws_public_url), pte_profile_id=VALUES(pte_profile_id), updated_at=NOW()`,
			payload.MerchantID, payload.SDKAppID, payload.APIPublicURL, payload.WSPublicURL, payload.PTEProfileID).Error
	case Deactivated:
		return businessDB.WithContext(ctx).Exec("DELETE FROM qixi_crm_b_merchant_im_sdk_app_view WHERE merchant_id = ? AND sdk_app_id = ?", payload.MerchantID, payload.SDKAppID).Error
	default:
		return nil
	}
}
