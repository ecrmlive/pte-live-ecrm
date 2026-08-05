// Package refundprocessor executes server-created provider refund requests.
// It never accepts a browser request and never marks a refund successful: a
// verified provider callback owns the terminal state.
package refundprocessor

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/business/order"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/cloudconfig"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/paymentconfig"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/wechatpayv3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Processor struct {
	db              *gorm.DB
	configs         *paymentconfig.Store
	platformConfigs *cloudconfig.Service
	client          *wechatpayv3.Client
}

func New(db *gorm.DB, configs *paymentconfig.Store, platformConfigs *cloudconfig.Service) *Processor {
	return &Processor{db: db, configs: configs, platformConfigs: platformConfigs, client: &wechatpayv3.Client{}}
}

func (p *Processor) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := p.ProcessPending(ctx, 20); err != nil {
					log.Printf("refund processor: %v", err)
				}
			}
		}
	}()
}

type pending struct {
	RefundID uint64 `gorm:"column:refund_id"`
}

const processingRecoveryAfter = 15 * time.Minute

// ProcessPending is exported for deterministic job tests. It processes only
// server-created WeChat transactions; balance/Alipay/mock have dedicated
// settlement adapters and cannot be silently routed through WeChat.
func (p *Processor) ProcessPending(ctx context.Context, limit int) (int, error) {
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}
	var rows []pending
	recoveryBefore := time.Now().UTC().Add(-processingRecoveryAfter)
	if err := p.db.WithContext(ctx).Table("qixi_crm_b_refund_transaction AS t").Select("t.refund_id").Joins("JOIN qixi_crm_b_refund AS r ON r.id = t.refund_id").Where("r.status = ? AND t.channel = ? AND (t.status IN ? OR (t.status = ? AND t.updated_at < ?))", "refunding", "wechat", []string{"created", "failed"}, "processing", recoveryBefore).Order("t.id ASC").Limit(limit).Scan(&rows).Error; err != nil {
		return 0, err
	}
	processed := 0
	for _, row := range rows {
		if err := p.processWechat(ctx, row.RefundID); err != nil {
			log.Printf("refund %d execution deferred: %v", row.RefundID, err)
			continue
		}
		processed++
	}
	return processed, nil
}

type execution struct {
	RefundID, OrderID, MerchantID, StoreID uint64
	RefundNo, Reason                       string
	RefundAmount, PaymentAmount            float64
	TransactionNo                          string
}

func (p *Processor) processWechat(ctx context.Context, refundID uint64) error {
	var work execution
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var refund struct {
			ID, OrderID uint64
			RefundNo    string
			Reason      string
			Amount      float64
			Status      string
		}
		if err := tx.Table("qixi_crm_b_refund").Where("id = ?", refundID).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&refund).Error; err != nil {
			return err
		}
		if refund.ID == 0 || refund.Status != "refunding" {
			return nil
		}
		var transaction struct {
			ID        uint64
			Status    string
			UpdatedAt time.Time `gorm:"column:updated_at"`
		}
		if err := tx.Table("qixi_crm_b_refund_transaction").Where("refund_id = ? AND channel = ?", refund.ID, "wechat").Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&transaction).Error; err != nil {
			return err
		}
		if transaction.ID == 0 || !recoverable(transaction.Status, transaction.UpdatedAt, time.Now().UTC()) {
			return nil
		}
		var payment struct {
			TransactionNo string  `gorm:"column:transaction_no"`
			Amount        float64 `gorm:"column:amount"`
			Status        string  `gorm:"column:status"`
			MerchantID    uint64  `gorm:"column:merchant_id"`
			StoreID       uint64  `gorm:"column:store_id"`
		}
		if err := tx.Table("qixi_crm_b_order AS o").Select("p.transaction_no,p.amount,p.status,o.merchant_id,o.store_id").Joins("JOIN qixi_crm_b_payment_transaction AS p ON p.group_order_id = o.group_order_id AND p.channel = 'wechat'").Where("o.id = ?", refund.OrderID).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&payment).Error; err != nil {
			return err
		}
		if payment.TransactionNo == "" || payment.Status != "succeeded" || refund.Amount <= 0 || refund.Amount > payment.Amount {
			return gorm.ErrRecordNotFound
		}
		now := time.Now().UTC()
		if result := tx.Table("qixi_crm_b_refund_transaction").Where("id = ? AND (status IN ? OR (status = ? AND updated_at < ?))", transaction.ID, []string{"created", "failed"}, "processing", now.Add(-processingRecoveryAfter)).Updates(map[string]any{"status": "processing", "updated_at": now}); result.Error != nil {
			return result.Error
		} else if result.RowsAffected != 1 {
			return nil
		}
		work = execution{RefundID: refund.ID, OrderID: refund.OrderID, MerchantID: payment.MerchantID, StoreID: payment.StoreID, RefundNo: refund.RefundNo, Reason: refund.Reason, RefundAmount: refund.Amount, PaymentAmount: payment.Amount, TransactionNo: payment.TransactionNo}
		return nil
	})
	if err != nil || work.RefundID == 0 {
		return err
	}
	values, err := p.paymentValues(ctx, work)
	if err != nil {
		_ = p.markFailed(ctx, work.RefundID)
		return err
	}
	config := order.WechatConfig(values, work.MerchantID == 0 && work.StoreID == 0)
	response, err := p.client.Refund(ctx, config, wechatpayv3.RefundRequest{OutTradeNo: work.TransactionNo, OutRefundNo: work.RefundNo, Reason: work.Reason, TotalCents: refundCents(work.PaymentAmount), RefundCents: refundCents(work.RefundAmount)})
	if err != nil {
		_ = p.markFailed(ctx, work.RefundID)
		return err
	}
	return p.db.WithContext(ctx).Table("qixi_crm_b_refund_transaction").Where("refund_id = ? AND channel = ? AND status = ?", work.RefundID, "wechat", "processing").Updates(map[string]any{"provider_refund_no": response.ProviderRefundNo, "status": "processing", "updated_at": time.Now().UTC()}).Error
}

func refundCents(amount float64) int64 { return int64(math.Round(amount * 100)) }

func recoverable(status string, updatedAt, now time.Time) bool {
	return status == "created" || status == "failed" || (status == "processing" && updatedAt.Before(now.Add(-processingRecoveryAfter)))
}

func (p *Processor) paymentValues(ctx context.Context, work execution) (paymentconfig.Values, error) {
	if work.MerchantID != 0 || work.StoreID != 0 {
		return p.configs.LoadStore(ctx, uint(work.StoreID))
	}
	values, err := p.platformConfigs.Values(ctx, "payment")
	return paymentconfig.Values(values), err
}

func (p *Processor) markFailed(ctx context.Context, refundID uint64) error {
	return p.db.WithContext(ctx).Table("qixi_crm_b_refund_transaction").Where("refund_id = ? AND channel = ? AND status = ?", refundID, "wechat", "processing").Updates(map[string]any{"status": "failed", "updated_at": time.Now().UTC()}).Error
}
