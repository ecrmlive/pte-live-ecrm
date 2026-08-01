package platformdiy

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"log"
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
type payload struct {
	PageID   uint64          `json:"page_id"`
	PageType string          `json:"page_type"`
	Name     string          `json:"name"`
	Document json.RawMessage `json:"document"`
	Status   string          `json:"status"`
	IsActive bool            `json:"is_active"`
}

func Start(ctx context.Context, db *gorm.DB, url string) (*nats.Conn, error) {
	if url == "" {
		return nil, nil
	}
	nc, err := nats.Connect(url, nats.Name("qixi_live_ecrm_api_business_platform_diy_projection"))
	if err != nil {
		return nil, err
	}
	if _, err = nc.Subscribe(Subject, func(m *nats.Msg) {
		if err := Apply(ctx, db, m.Data); err != nil {
			log.Printf("platform diy projection failed: %v", err)
		}
	}); err != nil {
		nc.Close()
		return nil, err
	}
	return nc, nc.Flush()
}
func Apply(ctx context.Context, db *gorm.DB, raw []byte) error {
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return err
	}
	var data payload
	if err := json.Unmarshal(env.Payload, &data); err != nil {
		return err
	}
	if data.PageID == 0 {
		return nil
	}
	switch env.EventType {
	case Upserted:
		return db.WithContext(ctx).Exec(`INSERT INTO qixi_crm_b_diy_page_view (source,page_id,store_id,page_type,name,document,status,is_active,updated_at) VALUES ('platform',?,0,?,?,?, ?,?,NOW()) ON DUPLICATE KEY UPDATE page_type=VALUES(page_type),name=VALUES(name),document=VALUES(document),status=VALUES(status),is_active=VALUES(is_active),updated_at=NOW()`, data.PageID, data.PageType, data.Name, data.Document, data.Status, data.IsActive).Error
	case Deleted:
		return db.WithContext(ctx).Exec("DELETE FROM qixi_crm_b_diy_page_view WHERE source = 'platform' AND page_id = ?", data.PageID).Error
	default:
		return nil
	}
}
