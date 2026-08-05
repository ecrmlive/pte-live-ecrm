package merchantsettlement

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const CommandSubject = "qixi.platform.merchant-settlement-command.v1"

type command struct {
	SettlementID    uint64 `json:"settlement_id"`
	Action          string `json:"action"`
	OperatorID      uint64 `json:"operator_id"`
	IdempotencyKey  string `json:"idempotency_key"`
	ReviewNote      string `json:"review_note,omitempty"`
	PayoutReference string `json:"payout_reference,omitempty"`
}

type commandResult struct {
	SettlementID uint64 `json:"settlement_id"`
	Status       string `json:"status,omitempty"`
	Code         string `json:"code,omitempty"`
}

type settlementBill struct {
	ID                   uint64    `gorm:"column:id"`
	StoreID              uint64    `gorm:"column:store_id"`
	MerchantID           uint64    `gorm:"column:merchant_id"`
	MerchantName         string    `gorm:"column:merchant_name"`
	RegionID             *uint64   `gorm:"column:region_id"`
	PeriodStart          time.Time `gorm:"column:period_start"`
	PeriodEnd            time.Time `gorm:"column:period_end"`
	Amount               float64   `gorm:"column:amount"`
	Status               string    `gorm:"column:status"`
	ReviewedByAdminID    uint64    `gorm:"column:reviewed_by_admin_id"`
	ReviewIdempotencyKey *string   `gorm:"column:review_idempotency_key"`
	PayoutIdempotencyKey *string   `gorm:"column:payout_idempotency_key"`
	ReviewNote           string    `gorm:"column:review_note"`
	PayoutReference      *string   `gorm:"column:payout_reference"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}

var (
	errCommandConflict = errors.New("settlement command status conflict")
	errCommandNotFound = errors.New("settlement command settlement not found")
)

// StartCommandSubscriber owns the merchant-side mutation path for platform
// review and payout commands. A reply only reports status; projection remains
// outbox/NATS driven.
func StartCommandSubscriber(ctx context.Context, merchantDB *gorm.DB, natsURL string) (*nats.Conn, error) {
	if strings.TrimSpace(natsURL) == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_merchant_settlement_command"))
	if err != nil {
		return nil, err
	}
	if _, err = nc.QueueSubscribe(CommandSubject, "pte_live_ecrm_merchant_settlement_command", func(msg *nats.Msg) {
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

func ApplyCommand(ctx context.Context, db *gorm.DB, raw []byte) (commandResult, error) {
	var in command
	if err := json.Unmarshal(raw, &in); err != nil {
		return commandResult{Code: "invalid"}, err
	}
	in.Action, in.IdempotencyKey = strings.TrimSpace(in.Action), strings.TrimSpace(in.IdempotencyKey)
	in.ReviewNote, in.PayoutReference = strings.TrimSpace(in.ReviewNote), strings.TrimSpace(in.PayoutReference)
	if !validCommand(in) {
		return commandResult{SettlementID: in.SettlementID, Code: "invalid"}, errors.New("invalid settlement command")
	}
	if db == nil {
		return commandResult{SettlementID: in.SettlementID, Code: "failed"}, errors.New("merchant database unavailable")
	}
	var out settlementBill
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row settlementBill
		if err := tx.Table("qixi_crm_m_settlement_bill").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", in.SettlementID).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID == 0 {
			return errCommandNotFound
		}
		if commandAlreadyApplied(row, in) {
			out = row
			return loadProjectionBill(tx, &out)
		}
		if commandKeyAlreadyUsed(row, in) || !commandAllowed(row.Status, in.Action) {
			return errCommandConflict
		}
		now := time.Now()
		updates := commandUpdates(in, now)
		where := "id = ? AND status = ?"
		args := []any{row.ID, row.Status}
		if in.Action == "mark_paid" {
			where += " AND payout_idempotency_key IS NULL"
		} else {
			where += " AND review_idempotency_key IS NULL"
		}
		result := tx.Table("qixi_crm_m_settlement_bill").Where(where, args...).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errCommandConflict
		}
		if err := loadProjectionBill(tx, &out, row.ID); err != nil {
			return err
		}
		return Enqueue(ctx, tx, Payload{SettlementID: out.ID, MerchantID: out.MerchantID, StoreID: out.StoreID, MerchantName: out.MerchantName, RegionID: out.RegionID, PeriodStart: out.PeriodStart, PeriodEnd: out.PeriodEnd, Amount: out.Amount, Status: out.Status, UpdatedAt: out.UpdatedAt})
	})
	result := commandResult{SettlementID: in.SettlementID}
	switch {
	case err == nil:
		result.Status = out.Status
	case errors.Is(err, errCommandNotFound):
		result.Code = "not_found"
	case errors.Is(err, errCommandConflict):
		result.Code = "conflict"
	default:
		result.Code = "failed"
	}
	return result, err
}

func loadProjectionBill(db *gorm.DB, target *settlementBill, ids ...uint64) error {
	q := db.Table("qixi_crm_m_settlement_bill AS b").Select("b.id,b.store_id,b.merchant_id,m.name AS merchant_name,m.region_id,b.period_start,b.period_end,b.amount,b.status,b.reviewed_by_admin_id,b.review_idempotency_key,b.payout_idempotency_key,b.review_note,b.payout_reference,b.updated_at").Joins("JOIN qixi_crm_m_merchant AS m ON m.id = b.merchant_id")
	if len(ids) == 0 {
		return q.Where("b.id = ?", target.ID).Scan(target).Error
	}
	return q.Where("b.id = ?", ids[0]).Scan(target).Error
}

func commandAlreadyApplied(row settlementBill, in command) bool {
	if in.Action == "mark_paid" {
		return row.Status == "paid" && row.PayoutIdempotencyKey != nil && *row.PayoutIdempotencyKey == in.IdempotencyKey && row.PayoutReference != nil && *row.PayoutReference == in.PayoutReference
	}
	status := "approved"
	if in.Action == "reject" {
		status = "rejected"
	}
	return row.Status == status && row.ReviewIdempotencyKey != nil && *row.ReviewIdempotencyKey == in.IdempotencyKey && row.ReviewedByAdminID == in.OperatorID && row.ReviewNote == in.ReviewNote
}

func commandKeyAlreadyUsed(row settlementBill, in command) bool {
	if in.Action == "mark_paid" {
		return row.PayoutIdempotencyKey != nil
	}
	return row.ReviewIdempotencyKey != nil
}

func commandAllowed(status, action string) bool {
	switch action {
	case "approve", "reject":
		return status == "withdraw_applied"
	case "mark_paid":
		return status == "approved"
	default:
		return false
	}
}

func commandUpdates(in command, now time.Time) map[string]any {
	switch in.Action {
	case "approve":
		return map[string]any{"status": "approved", "review_idempotency_key": in.IdempotencyKey, "reviewed_by_admin_id": in.OperatorID, "review_note": in.ReviewNote, "reviewed_at": now, "version": gorm.Expr("version + 1")}
	case "reject":
		return map[string]any{"status": "rejected", "review_idempotency_key": in.IdempotencyKey, "reviewed_by_admin_id": in.OperatorID, "review_note": in.ReviewNote, "reviewed_at": now, "version": gorm.Expr("version + 1")}
	default:
		return map[string]any{"status": "paid", "payout_idempotency_key": in.IdempotencyKey, "payout_reference": in.PayoutReference, "paid_at": now, "version": gorm.Expr("version + 1")}
	}
}

func validCommand(in command) bool {
	if in.SettlementID == 0 || in.OperatorID == 0 || !validCommandIdempotencyKey(in.IdempotencyKey) {
		return false
	}
	switch in.Action {
	case "approve":
		return len([]rune(in.ReviewNote)) <= 500
	case "reject":
		return in.ReviewNote != "" && len([]rune(in.ReviewNote)) <= 500
	case "mark_paid":
		return len([]rune(in.PayoutReference)) >= 3 && len([]rune(in.PayoutReference)) <= 128
	default:
		return false
	}
}

func validCommandIdempotencyKey(value string) bool {
	length := len([]rune(strings.TrimSpace(value)))
	return length >= 8 && length <= 128
}
