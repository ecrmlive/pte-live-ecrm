package merchantledger

import (
	"context"
	"time"

	merchantsettlement "github.com/crmlive/pte-live-ecrm/api-merchant/internal/event/merchantsettlement"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StartPendingBillFreezer closes elapsed calendar-month bills. The process is
// intentionally merchant-local and writes the normal settlement projection
// outbox in the same transaction; platform never freezes merchant facts.
func StartPendingBillFreezer(ctx context.Context, merchantDB *gorm.DB) {
	if merchantDB == nil {
		return
	}
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		freezeExpired(ctx, merchantDB, time.Now().UTC())
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-ticker.C:
				freezeExpired(ctx, merchantDB, now.UTC())
			}
		}
	}()
}

func freezeExpired(ctx context.Context, db *gorm.DB, now time.Time) {
	var ids []uint64
	if err := db.WithContext(ctx).Table("qixi_crm_m_settlement_bill").Where("status = ? AND period_end < ?", "bill_pending", now).Order("id ASC").Limit(100).Pluck("id", &ids).Error; err != nil {
		return
	}
	for _, id := range ids {
		_ = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var row settlementBill
			if err := tx.Table("qixi_crm_m_settlement_bill").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).Scan(&row).Error; err != nil || row.ID == 0 || row.Status != "bill_pending" || !row.PeriodEnd.Before(now) {
				return err
			}
			if result := tx.Table("qixi_crm_m_settlement_bill").Where("id = ? AND status = ?", row.ID, "bill_pending").Updates(map[string]any{"status": "bill_frozen", "version": gorm.Expr("version + 1")}); result.Error != nil || result.RowsAffected != 1 {
				return result.Error
			}
			if err := tx.Table("qixi_crm_m_settlement_bill AS b").Select("b.id,b.store_id,b.merchant_id,m.name AS merchant_name,m.region_id,b.period_start,b.period_end,b.amount,b.status,b.updated_at").Joins("JOIN qixi_crm_m_merchant AS m ON m.id = b.merchant_id").Where("b.id = ?", row.ID).Scan(&row).Error; err != nil {
				return err
			}
			return merchantsettlement.Enqueue(ctx, tx, merchantsettlement.Payload{SettlementID: row.ID, MerchantID: row.MerchantID, StoreID: row.StoreID, MerchantName: row.MerchantName, RegionID: row.RegionID, PeriodStart: row.PeriodStart, PeriodEnd: row.PeriodEnd, Amount: row.Amount, Status: row.Status, UpdatedAt: row.UpdatedAt})
		})
	}
}
