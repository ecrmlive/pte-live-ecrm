package order

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrInsufficientPoints = errors.New("账户积分不足")

// ResolveIntegral derives the policy from the validated checkout snapshot; a
// browser never supplies an exchange rate or a discount amount.
func ResolveIntegral(ctx context.Context, db *gorm.DB, userID uint64, checkout Checkout, couponDiscountCents int64, requested, lock bool) (IntegralQuote, error) {
	if len(checkout.Stores) != 1 || len(checkout.Stores[0].Lines) == 0 || checkout.Stores[0].Lines[0].ProductType != 0 || !requested {
		return IntegralQuote{}, nil
	}
	payable := checkout.TotalCents - couponDiscountCents
	if payable <= 0 {
		return IntegralQuote{}, nil
	}
	var account struct {
		Points int64 `gorm:"column:points"`
	}
	query := db.WithContext(ctx).Table("qixi_crm_b_member_account").Select("points").Where("user_id = ?", userID)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := query.Take(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return IntegralQuote{}, nil
		}
		return IntegralQuote{}, err
	}
	return QuoteIntegral(checkout.Stores[0].IntegralPolicy, account.Points, payable, true), nil
}

func deductIntegral(tx *gorm.DB, userID uint64, group groupRow, quote IntegralQuote) error {
	if quote.PointsUsed == 0 {
		return nil
	}
	result := tx.Table("qixi_crm_b_member_account").Where("user_id = ? AND points >= ?", userID, quote.PointsUsed).Update("points", gorm.Expr("points - ?", quote.PointsUsed))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrInsufficientPoints
	}
	return tx.Exec("INSERT INTO qixi_crm_b_asset_ledger (user_id, asset_type, amount, reference_type, reference_id, idempotency_key) VALUES (?, 'points', ?, 'order_integral_deduction', ?, ?)", userID, -quote.PointsUsed, group.OrderNo, fmt.Sprintf("order-integral:%d", group.ID)).Error
}

func restoreIntegral(tx *gorm.DB, group groupRow) error {
	if group.PointsAmount <= 0 {
		return nil
	}
	if err := tx.Table("qixi_crm_b_member_account").Where("user_id = ?", group.UserID).Update("points", gorm.Expr("points + ?", group.PointsAmount)).Error; err != nil {
		return err
	}
	return tx.Exec("INSERT INTO qixi_crm_b_asset_ledger (user_id, asset_type, amount, reference_type, reference_id, idempotency_key) VALUES (?, 'points', ?, 'order_integral_restore', ?, ?)", group.UserID, group.PointsAmount, group.OrderNo, fmt.Sprintf("order-integral-cancel:%d", group.ID)).Error
}
