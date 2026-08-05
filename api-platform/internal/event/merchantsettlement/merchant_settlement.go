// Package merchantsettlement projects merchant-owned settlement changes into
// the platform database. Platform never connects to qixi_crm_merchant.
package merchantsettlement

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const (
	Subject  = "qixi.merchant.settlement.v1"
	Upserted = "merchant.settlement.upsert"
)

type Envelope struct {
	EventID   uint64          `json:"event_id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

type Payload struct {
	SettlementID uint64    `json:"settlement_id"`
	MerchantID   uint64    `json:"merchant_id"`
	StoreID      uint64    `json:"store_id"`
	MerchantName string    `json:"merchant_name"`
	RegionID     *uint64   `json:"region_id,omitempty"`
	PeriodStart  time.Time `json:"period_start"`
	PeriodEnd    time.Time `json:"period_end"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func Start(ctx context.Context, adminDB *gorm.DB, natsURL string) (*nats.Conn, error) {
	if adminDB == nil || strings.TrimSpace(natsURL) == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_platform_merchant_settlement_projection"))
	if err != nil {
		return nil, err
	}
	if _, err = nc.Subscribe(Subject, func(msg *nats.Msg) {
		if err := Apply(ctx, adminDB, msg.Data); err != nil {
			log.Printf("merchant settlement projection failed: %v", err)
		}
	}); err != nil {
		nc.Close()
		return nil, err
	}
	return nc, nc.FlushTimeout(2 * time.Second)
}

func Apply(ctx context.Context, db *gorm.DB, raw []byte) error {
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return err
	}
	if env.EventType != Upserted {
		return nil
	}
	var payload Payload
	if err := json.Unmarshal(env.Payload, &payload); err != nil {
		return err
	}
	if !validPayload(payload) {
		return gorm.ErrInvalidData
	}
	return db.WithContext(ctx).Exec(`INSERT INTO qixi_crm_a_merchant_settlement_view
  (source_settlement_id,merchant_id,store_id,merchant_name,region_id,period_start,period_end,amount,status,updated_at)
  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
  ON DUPLICATE KEY UPDATE merchant_id=VALUES(merchant_id),store_id=VALUES(store_id),merchant_name=VALUES(merchant_name),region_id=VALUES(region_id),period_start=VALUES(period_start),period_end=VALUES(period_end),amount=VALUES(amount),status=VALUES(status),updated_at=VALUES(updated_at)`,
		payload.SettlementID, payload.MerchantID, payload.StoreID, payload.MerchantName, payload.RegionID, payload.PeriodStart, payload.PeriodEnd, payload.Amount, payload.Status, payload.UpdatedAt).Error
}

func validPayload(p Payload) bool {
	if p.SettlementID == 0 || p.MerchantID == 0 || p.StoreID == 0 || strings.TrimSpace(p.MerchantName) == "" || math.IsNaN(p.Amount) || math.IsInf(p.Amount, 0) || p.PeriodStart.IsZero() || p.PeriodEnd.Before(p.PeriodStart) || p.UpdatedAt.IsZero() {
		return false
	}
	switch p.Status {
	case "bill_pending", "bill_frozen", "withdraw_applied", "approved", "paid", "rejected", "cancelled":
		return true
	default:
		return false
	}
}
