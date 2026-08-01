// Package merchantdiy materializes merchant DIY pages in qixi_crm_business.
// C-end APIs consume this projection and never connect to qixi_crm_merchant.
package merchantdiy

import (
	"context"
	"encoding/json"
	"log"

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
type Payload struct {
	PageID   uint64          `json:"page_id"`
	StoreID  uint64          `json:"store_id"`
	PageType string          `json:"page_type"`
	Name     string          `json:"name"`
	Document json.RawMessage `json:"document"`
	Status   string          `json:"status"`
	IsActive bool            `json:"is_active"`
}

func StartBusinessProjection(ctx context.Context, businessDB *gorm.DB, natsURL string) (*nats.Conn, error) {
	if natsURL == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("qixi_live_ecrm_api_business_diy_projection"))
	if err != nil {
		return nil, err
	}
	if _, err = nc.Subscribe(Subject, func(msg *nats.Msg) {
		if err := ApplyProjection(ctx, businessDB, msg.Data); err != nil {
			log.Printf("merchant diy projection failed: %v", err)
		}
	}); err != nil {
		nc.Close()
		return nil, err
	}
	return nc, nc.Flush()
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
	if payload.PageID == 0 || payload.StoreID == 0 {
		return nil
	}
	switch envelope.EventType {
	case Upserted:
		return businessDB.WithContext(ctx).Exec(`INSERT INTO qixi_crm_b_diy_page_view
      (source, page_id, store_id, page_type, name, document, status, is_active, updated_at)
      VALUES ('merchant', ?, ?, ?, ?, ?, ?, ?, NOW())
      ON DUPLICATE KEY UPDATE store_id=VALUES(store_id), page_type=VALUES(page_type), name=VALUES(name),
      document=VALUES(document), status=VALUES(status), is_active=VALUES(is_active), updated_at=NOW()`,
			payload.PageID, payload.StoreID, payload.PageType, payload.Name, payload.Document, payload.Status, payload.IsActive).Error
	case Deleted:
		return businessDB.WithContext(ctx).Exec("DELETE FROM qixi_crm_b_diy_page_view WHERE source = 'merchant' AND page_id = ? AND store_id = ?", payload.PageID, payload.StoreID).Error
	default:
		return nil
	}
}
