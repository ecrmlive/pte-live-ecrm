package merchantintegralpolicy

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const Subject = "qixi.merchant.integral-policy.v1"
const Upserted = "merchant.integral_policy.upsert"

type payload struct {
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

func StartBusinessProjection(ctx context.Context, db *gorm.DB, natsURL string) (*nats.Conn, error) {
	if natsURL == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_business_integral_policy"))
	if err != nil {
		return nil, err
	}
	_, err = nc.Subscribe(Subject, func(msg *nats.Msg) {
		if err := ApplyProjection(ctx, db, msg.Data); err != nil {
			log.Printf("merchant integral policy projection failed: %v", err)
		}
	})
	if err != nil {
		nc.Close()
		return nil, err
	}
	return nc, nc.Flush()
}
func valid(p payload) bool {
	return p.StoreID > 0 && p.PointsPerYuan > 0 && p.MaxDeductionBps > 0 && p.MaxDeductionBps <= 10000
}

func ApplyProjection(ctx context.Context, db *gorm.DB, raw []byte) error {
	var e envelope
	if err := json.Unmarshal(raw, &e); err != nil {
		return err
	}
	if e.EventType != Upserted {
		return nil
	}
	var p payload
	if err := json.Unmarshal(e.Payload, &p); err != nil {
		return err
	}
	if !valid(p) {
		return nil
	}
	return db.WithContext(ctx).Table("qixi_crm_b_store_view").Where("store_id = ?", p.StoreID).Updates(map[string]any{"integral_enabled": p.Enabled, "integral_points_per_yuan": p.PointsPerYuan, "integral_max_deduction_bps": p.MaxDeductionBps}).Error
}
