// Package platformcity materializes platform-owned administrative areas for C-end reads.
package platformcity

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const Subject = "qixi.platform.city.v1"
const (
	Upserted = "platform.city.upsert"
	Deleted  = "platform.city.deleted"
)

type envelope struct {
	EventID   uint64          `json:"event_id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}
type payload struct {
	CityID   uint64 `json:"city_id"`
	ParentID uint64 `json:"parent_id"`
	Name     string `json:"name"`
	Level    int8   `json:"level"`
	IsShow   int8   `json:"is_show"`
}

func Start(ctx context.Context, db *gorm.DB, url string) (*nats.Conn, error) {
	if url == "" {
		return nil, nil
	}
	nc, err := nats.Connect(url, nats.Name("pte_live_ecrm_api_business_platform_city_projection"))
	if err != nil {
		return nil, err
	}
	if _, err = nc.Subscribe(Subject, func(msg *nats.Msg) {
		if err := Apply(ctx, db, msg.Data); err != nil {
			log.Printf("platform city projection failed: %v", err)
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
	if data.CityID == 0 {
		return nil
	}
	switch env.EventType {
	case Upserted:
		return db.WithContext(ctx).Exec("INSERT INTO qixi_crm_b_city_view (city_id,parent_id,name,level,is_show,updated_at) VALUES (?,?,?,?,?,NOW()) ON DUPLICATE KEY UPDATE parent_id=VALUES(parent_id),name=VALUES(name),level=VALUES(level),is_show=VALUES(is_show),updated_at=NOW()", data.CityID, data.ParentID, data.Name, data.Level, data.IsShow).Error
	case Deleted:
		return db.WithContext(ctx).Exec("DELETE FROM qixi_crm_b_city_view WHERE city_id = ?", data.CityID).Error
	}
	return nil
}
