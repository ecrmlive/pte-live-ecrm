// Package merchantstock owns merchant-side inventory mutations requested by
// the business order saga. It never accepts a browser request and never reads
// the business database: NATS carries an idempotent command only.
package merchantstock

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const CommandSubject = "qixi.business.merchant-stock-command.v1"

type command struct {
	Action         string    `json:"action"`
	OrderID        uint64    `json:"order_id"`
	StoreID        uint64    `json:"store_id"`
	MerchantSKUID  uint64    `json:"merchant_sku_id"`
	Quantity       int       `json:"quantity"`
	ExpiresAt      time.Time `json:"expires_at,omitempty"`
	IdempotencyKey string    `json:"idempotency_key"`
}

type commandResult struct {
	OrderID       uint64 `json:"order_id"`
	MerchantSKUID uint64 `json:"merchant_sku_id"`
	Status        string `json:"status,omitempty"`
	Code          string `json:"code,omitempty"`
}

type stockReservation struct {
	ID        uint64    `gorm:"column:id"`
	SKUId     uint64    `gorm:"column:sku_id"`
	OrderID   uint64    `gorm:"column:order_id"`
	Quantity  int       `gorm:"column:quantity"`
	Status    string    `gorm:"column:status"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
}

func (stockReservation) TableName() string { return "qixi_crm_m_stock_reservation" }

type productSKU struct {
	ID            uint64 `gorm:"column:id"`
	Stock         int    `gorm:"column:stock"`
	Status        int8   `gorm:"column:status"`
	ProductStatus string `gorm:"column:product_status"`
}

func (productSKU) TableName() string { return "qixi_crm_m_product_sku" }

type stockLedger struct {
	SKUId           uint64 `gorm:"column:sku_id"`
	ChangeQuantity  int    `gorm:"column:change_quantity"`
	BalanceQuantity int    `gorm:"column:balance_quantity"`
	ReasonType      string `gorm:"column:reason_type"`
	ReferenceType   string `gorm:"column:reference_type"`
	ReferenceID     string `gorm:"column:reference_id"`
	IdempotencyKey  string `gorm:"column:idempotency_key"`
}

func (stockLedger) TableName() string { return "qixi_crm_m_stock_ledger" }

var (
	errStockCommandConflict = errors.New("stock command state conflict")
	errStockNotFound        = errors.New("merchant sku not found")
	errStockInsufficient    = errors.New("merchant sku stock insufficient")
)

// StartCommandSubscriber owns all merchant inventory changes from order
// workflows. A NATS reply only acknowledges the local merchant transaction;
// business outbox delivery is responsible for retrying commands.
func StartCommandSubscriber(ctx context.Context, merchantDB *gorm.DB, natsURL string) (*nats.Conn, error) {
	if strings.TrimSpace(natsURL) == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_merchant_stock_command"))
	if err != nil {
		return nil, err
	}
	if _, err = nc.QueueSubscribe(CommandSubject, "pte_live_ecrm_merchant_stock_command", func(msg *nats.Msg) {
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

// ApplyCommand holds the SKU/reservation row locks in the merchant database.
// Reserve deducts sellable stock, confirm preserves that deduction, release
// restores an unpaid reservation, and restock restores a confirmed refund.
func ApplyCommand(ctx context.Context, db *gorm.DB, raw []byte) (commandResult, error) {
	var in command
	if err := json.Unmarshal(raw, &in); err != nil {
		return commandResult{Code: "invalid"}, err
	}
	in.Action, in.IdempotencyKey = strings.TrimSpace(in.Action), strings.TrimSpace(in.IdempotencyKey)
	result := commandResult{OrderID: in.OrderID, MerchantSKUID: in.MerchantSKUID}
	if !validCommand(in) {
		result.Code = "invalid"
		return result, errors.New("invalid stock command")
	}
	if db == nil {
		result.Code = "failed"
		return result, errors.New("merchant database unavailable")
	}
	var reservation stockReservation
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var sku productSKU
		if err := tx.Table("qixi_crm_m_product_sku AS s").Select("s.id,s.stock,s.status,p.status AS product_status").Joins("JOIN qixi_crm_m_product AS p ON p.id = s.product_id").Clauses(clause.Locking{Strength: "UPDATE"}).Where("s.id = ? AND p.store_id = ?", in.MerchantSKUID, in.StoreID).Scan(&sku).Error; err != nil {
			return err
		}
		if sku.ID == 0 {
			return errStockNotFound
		}
		if sku.Status != 1 {
			return errStockNotFound
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND sku_id = ?", in.OrderID, in.MerchantSKUID).First(&reservation).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) || in.Action != "reserve" {
				return err
			}
			reservation = stockReservation{}
		}
		switch in.Action {
		case "reserve":
			if sku.ProductStatus != "on_sale" {
				return errStockNotFound
			}
			if reservation.ID != 0 {
				if reservation.Quantity == in.Quantity && reservation.Status == "reserved" {
					return nil
				}
				return errStockCommandConflict
			}
			if sku.Stock < in.Quantity {
				return errStockInsufficient
			}
			balance := sku.Stock - in.Quantity
			if err := tx.Model(&productSKU{}).Where("id = ? AND stock >= ?", sku.ID, in.Quantity).Update("stock", balance).Error; err != nil {
				return err
			}
			reservation = stockReservation{SKUId: sku.ID, OrderID: in.OrderID, Quantity: in.Quantity, Status: "reserved", ExpiresAt: in.ExpiresAt.UTC()}
			if err := tx.Create(&reservation).Error; err != nil {
				return err
			}
			return writeLedger(tx, sku.ID, -in.Quantity, balance, "order_reserve", in)
		case "confirm":
			if reservation.ID == 0 || reservation.Quantity != in.Quantity {
				return errStockCommandConflict
			}
			if reservation.Status == "confirmed" {
				return nil
			}
			if reservation.Status != "reserved" {
				return errStockCommandConflict
			}
			return tx.Model(&stockReservation{}).Where("id = ? AND status = ?", reservation.ID, "reserved").Update("status", "confirmed").Error
		case "release", "restock":
			if reservation.ID == 0 || reservation.Quantity != in.Quantity {
				return errStockCommandConflict
			}
			if reservation.Status == "released" {
				return nil
			}
			if in.Action == "release" && reservation.Status != "reserved" {
				return errStockCommandConflict
			}
			if in.Action == "restock" && reservation.Status != "confirmed" {
				return errStockCommandConflict
			}
			balance := sku.Stock + in.Quantity
			if err := tx.Model(&productSKU{}).Where("id = ?", sku.ID).Update("stock", balance).Error; err != nil {
				return err
			}
			if err := tx.Model(&stockReservation{}).Where("id = ? AND status = ?", reservation.ID, reservation.Status).Update("status", "released").Error; err != nil {
				return err
			}
			return writeLedger(tx, sku.ID, in.Quantity, balance, "order_"+in.Action, in)
		default:
			return errors.New("invalid stock action")
		}
	})
	if err == nil {
		if reservation.Status == "" {
			result.Status = "reserved"
		} else if in.Action == "confirm" {
			result.Status = "confirmed"
		} else if in.Action == "release" || in.Action == "restock" {
			result.Status = "released"
		} else {
			result.Status = reservation.Status
		}
		return result, nil
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound), errors.Is(err, errStockNotFound):
		result.Code = "not_found"
	case errors.Is(err, errStockInsufficient):
		result.Code = "insufficient"
	case errors.Is(err, errStockCommandConflict):
		result.Code = "conflict"
	default:
		result.Code = "failed"
	}
	return result, err
}

func writeLedger(tx *gorm.DB, skuID uint64, change, balance int, reason string, in command) error {
	return tx.Create(&stockLedger{SKUId: skuID, ChangeQuantity: change, BalanceQuantity: balance, ReasonType: reason, ReferenceType: "order", ReferenceID: strconv.FormatUint(in.OrderID, 10), IdempotencyKey: in.IdempotencyKey}).Error
}

func validCommand(in command) bool {
	if in.OrderID == 0 || in.StoreID == 0 || in.MerchantSKUID == 0 || in.Quantity < 1 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 {
		return false
	}
	switch in.Action {
	case "reserve":
		return !in.ExpiresAt.IsZero() && in.ExpiresAt.After(time.Now().UTC().Add(time.Minute)) && in.ExpiresAt.Before(time.Now().UTC().Add(25*time.Hour))
	case "confirm", "release", "restock":
		return in.ExpiresAt.IsZero()
	default:
		return false
	}
}
