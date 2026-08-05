// Package merchantledger owns merchant-side settlement accrual and reversal
// mutations. It accepts only idempotent NATS commands from the business order
// service and never accepts a browser payload as a finance fact.
package merchantledger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	merchantsettlement "github.com/crmlive/pte-live-ecrm/api-merchant/internal/event/merchantsettlement"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const CommandSubject = "qixi.business.merchant-settlement-command.v1"

type command struct {
	Action         string  `json:"action"`
	OrderID        uint64  `json:"order_id"`
	RefundID       uint64  `json:"refund_id,omitempty"`
	StoreID        uint64  `json:"store_id"`
	MerchantID     uint64  `json:"merchant_id"`
	Amount         float64 `json:"amount"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type commandResult struct {
	OrderID uint64 `json:"order_id"`
	Status  string `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
}

type settlementEntry struct {
	ID             uint64    `gorm:"column:id"`
	StoreID        uint64    `gorm:"column:store_id"`
	MerchantID     uint64    `gorm:"column:merchant_id"`
	OrderID        uint64    `gorm:"column:order_id"`
	RefundID       uint64    `gorm:"column:refund_id"`
	EntryType      string    `gorm:"column:entry_type"`
	Amount         float64   `gorm:"column:amount"`
	IdempotencyKey string    `gorm:"column:idempotency_key"`
	OccurredAt     time.Time `gorm:"column:occurred_at"`
}

func (settlementEntry) TableName() string { return "qixi_crm_m_settlement_entry" }

type settlementBill struct {
	ID           uint64    `gorm:"column:id"`
	StoreID      uint64    `gorm:"column:store_id"`
	MerchantID   uint64    `gorm:"column:merchant_id"`
	MerchantName string    `gorm:"column:merchant_name"`
	RegionID     *uint64   `gorm:"column:region_id"`
	PeriodStart  time.Time `gorm:"column:period_start"`
	PeriodEnd    time.Time `gorm:"column:period_end"`
	Amount       float64   `gorm:"column:amount"`
	Status       string    `gorm:"column:status"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

var (
	errInvalidCommand  = errors.New("invalid settlement ledger command")
	errStoreMismatch   = errors.New("merchant store mismatch")
	errCommandConflict = errors.New("settlement ledger command conflict")
)

// StartCommandSubscriber serializes settlement facts through the merchant
// database. The returned connection is held by main until process shutdown.
func StartCommandSubscriber(ctx context.Context, merchantDB *gorm.DB, natsURL string) (*nats.Conn, error) {
	if strings.TrimSpace(natsURL) == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_merchant_settlement_ledger"))
	if err != nil {
		return nil, err
	}
	if _, err = nc.QueueSubscribe(CommandSubject, "pte_live_ecrm_merchant_settlement_ledger", func(msg *nats.Msg) {
		out, applyErr := ApplyCommand(ctx, merchantDB, msg.Data)
		if applyErr != nil && out.Code == "" {
			out.Code = "failed"
		}
		wire, _ := json.Marshal(out)
		_ = msg.Respond(wire)
	}); err != nil {
		nc.Close()
		return nil, err
	}
	return nc, nc.Flush()
}

// ApplyCommand records the immutable entry and adjusts only the current open
// bill in one merchant transaction. A refund after a frozen period becomes a
// negative current-period adjustment, preserving historical bills intact.
func ApplyCommand(ctx context.Context, db *gorm.DB, raw []byte) (commandResult, error) {
	var in command
	if err := json.Unmarshal(raw, &in); err != nil {
		return commandResult{Code: "invalid"}, err
	}
	in.Action, in.IdempotencyKey = strings.TrimSpace(in.Action), strings.TrimSpace(in.IdempotencyKey)
	out := commandResult{OrderID: in.OrderID}
	if !valid(in) {
		out.Code = "invalid"
		return out, errInvalidCommand
	}
	if db == nil {
		out.Code = "failed"
		return out, errors.New("merchant database unavailable")
	}
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var store struct {
			ID uint64 `gorm:"column:id"`
		}
		if err := tx.Table("qixi_crm_m_store").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND merchant_id = ? AND status = ?", in.StoreID, in.MerchantID, 1).Scan(&store).Error; err != nil {
			return err
		}
		if store.ID == 0 {
			return errStoreMismatch
		}
		var existing settlementEntry
		if err := tx.Model(&settlementEntry{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("idempotency_key = ?", in.IdempotencyKey).Scan(&existing).Error; err != nil {
			return err
		}
		if existing.ID != 0 {
			if sameEntry(existing, in) {
				out.Status = "accepted"
				return nil
			}
			return errCommandConflict
		}
		now := time.Now().UTC()
		entryType, signedAmount := "order_accrual", in.Amount
		if in.Action == "reverse" {
			entryType, signedAmount = "refund_reversal", -in.Amount
		}
		if err := tx.Create(&settlementEntry{StoreID: in.StoreID, MerchantID: in.MerchantID, OrderID: in.OrderID, RefundID: in.RefundID, EntryType: entryType, Amount: signedAmount, IdempotencyKey: in.IdempotencyKey, OccurredAt: now}).Error; err != nil {
			return err
		}
		bill, err := appendOpenBill(tx, in.StoreID, in.MerchantID, signedAmount, now)
		if err != nil {
			return err
		}
		if err := merchantsettlement.Enqueue(ctx, tx, merchantsettlement.Payload{SettlementID: bill.ID, MerchantID: bill.MerchantID, StoreID: bill.StoreID, MerchantName: bill.MerchantName, RegionID: bill.RegionID, PeriodStart: bill.PeriodStart, PeriodEnd: bill.PeriodEnd, Amount: bill.Amount, Status: bill.Status, UpdatedAt: bill.UpdatedAt}); err != nil {
			return err
		}
		out.Status = "accepted"
		return nil
	})
	if err == nil {
		return out, nil
	}
	switch {
	case errors.Is(err, errInvalidCommand):
		out.Code = "invalid"
	case errors.Is(err, errStoreMismatch):
		out.Code = "not_found"
	case errors.Is(err, errCommandConflict):
		out.Code = "conflict"
	default:
		out.Code = "failed"
	}
	return out, err
}

func valid(in command) bool {
	if in.OrderID == 0 || in.StoreID == 0 || in.MerchantID == 0 || in.Amount <= 0 || math.IsNaN(in.Amount) || math.IsInf(in.Amount, 0) || len(in.IdempotencyKey) < 8 || len(in.IdempotencyKey) > 128 {
		return false
	}
	switch in.Action {
	case "accrue":
		return in.RefundID == 0 && in.IdempotencyKey == fmt.Sprintf("settlement:accrue:%d", in.OrderID)
	case "reverse":
		return in.RefundID != 0 && in.IdempotencyKey == fmt.Sprintf("settlement:reverse:%d", in.RefundID)
	default:
		return false
	}
}

func sameEntry(row settlementEntry, in command) bool {
	wantType, wantAmount := "order_accrual", in.Amount
	if in.Action == "reverse" {
		wantType, wantAmount = "refund_reversal", -in.Amount
	}
	return row.StoreID == in.StoreID && row.MerchantID == in.MerchantID && row.OrderID == in.OrderID && row.RefundID == in.RefundID && row.EntryType == wantType && row.Amount == wantAmount
}

func appendOpenBill(tx *gorm.DB, storeID, merchantID uint64, signedAmount float64, now time.Time) (settlementBill, error) {
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	var bill settlementBill
	if err := tx.Table("qixi_crm_m_settlement_bill").Clauses(clause.Locking{Strength: "UPDATE"}).Where("store_id = ? AND merchant_id = ? AND period_start = ? AND period_end = ? AND status = ?", storeID, merchantID, start, end, "bill_pending").Scan(&bill).Error; err != nil {
		return settlementBill{}, err
	}
	if bill.ID == 0 {
		if err := tx.Table("qixi_crm_m_settlement_bill").Create(map[string]any{"store_id": storeID, "merchant_id": merchantID, "period_start": start, "period_end": end, "amount": signedAmount, "status": "bill_pending"}).Error; err != nil {
			return settlementBill{}, err
		}
	} else if err := tx.Table("qixi_crm_m_settlement_bill").Where("id = ? AND status = ?", bill.ID, "bill_pending").Updates(map[string]any{"amount": gorm.Expr("amount + ?", signedAmount), "version": gorm.Expr("version + 1")}).Error; err != nil {
		return settlementBill{}, err
	}
	if err := tx.Table("qixi_crm_m_settlement_bill AS b").Select("b.id,b.store_id,b.merchant_id,m.name AS merchant_name,m.region_id,b.period_start,b.period_end,b.amount,b.status,b.updated_at").Joins("JOIN qixi_crm_m_merchant AS m ON m.id = b.merchant_id").Where("b.store_id = ? AND b.merchant_id = ? AND b.period_start = ? AND b.period_end = ? AND b.status = ?", storeID, merchantID, start, end, "bill_pending").Scan(&bill).Error; err != nil {
		return settlementBill{}, err
	}
	if bill.ID == 0 {
		return settlementBill{}, errCommandConflict
	}
	return bill, nil
}
