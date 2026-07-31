package order

import (
	"context"
	"errors"
	"math"
	"sort"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrCouponOwnership = errors.New("优惠券不存在、已失效或无权使用")
	ErrCouponConflict  = errors.New("优惠券不可同时使用")
	ErrCouponMinNotMet = errors.New("订单金额未达到优惠券使用门槛")
)

// CouponPricing is calculated entirely in integer cents. Allocations are
// persisted with the order so a coupon cannot be reused before payment.
type CouponPricing struct {
	DiscountCents int64
	Allocations   []CouponAllocation
}
type CouponAllocation struct {
	CouponUserID  uint64
	DiscountCents int64
}

type couponRow struct {
	UserCouponID  uint64     `gorm:"column:id"`
	CouponID      uint64     `gorm:"column:coupon_id"`
	StoreID       uint64     `gorm:"column:store_id"`
	DiscountType  string     `gorm:"column:discount_type"`
	DiscountValue float64    `gorm:"column:discount_value"`
	MinAmount     float64    `gorm:"column:min_amount"`
	StartsAt      *time.Time `gorm:"column:starts_at"`
	EndsAt        *time.Time `gorm:"column:ends_at"`
	Status        int8       `gorm:"column:status"`
}

func ResolveCoupons(ctx context.Context, db *gorm.DB, userID uint64, checkout Checkout, ids []uint64, lock bool) (CouponPricing, error) {
	ids, err := uniqueCouponIDs(ids)
	if err != nil {
		return CouponPricing{}, err
	}
	if len(ids) == 0 {
		return CouponPricing{}, nil
	}
	if len(checkout.Stores) != 1 {
		return CouponPricing{}, ErrCouponConflict
	}
	q := db.WithContext(ctx).Table("qixi_crm_b_coupon_user AS u").
		Select("u.id,u.coupon_id,c.store_id,c.discount_type,c.discount_value,c.min_amount,c.starts_at,c.ends_at,c.status").
		Joins("JOIN qixi_crm_b_coupon_template_view AS c ON c.coupon_id = u.coupon_id").
		Where("u.user_id = ? AND u.id IN ? AND u.status = ?", userID, ids, "unused")
	if lock {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var rows []couponRow
	if err := q.Find(&rows).Error; err != nil {
		return CouponPricing{}, err
	}
	if len(rows) != len(ids) {
		return CouponPricing{}, ErrCouponOwnership
	}
	storeID := checkout.Stores[0].StoreID
	now := time.Now()
	for _, row := range rows {
		if row.Status != 1 || (row.StartsAt != nil && row.StartsAt.After(now)) || (row.EndsAt != nil && row.EndsAt.Before(now)) {
			return CouponPricing{}, ErrCouponOwnership
		}
		if row.StoreID != 0 && row.StoreID != storeID {
			return CouponPricing{}, ErrCouponConflict
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].StoreID != 0 && rows[j].StoreID == 0 })
	storeCoupons, platformCoupons := 0, 0
	result := CouponPricing{Allocations: make([]CouponAllocation, 0, len(rows))}
	for _, row := range rows {
		if row.StoreID == 0 {
			platformCoupons++
		} else {
			storeCoupons++
		}
		if storeCoupons > 1 || platformCoupons > 1 {
			return CouponPricing{}, ErrCouponConflict
		}
		base := checkout.TotalCents - result.DiscountCents
		minimum := int64(math.Round(row.MinAmount * 100))
		if base < minimum {
			return CouponPricing{}, ErrCouponMinNotMet
		}
		value := int64(math.Round(row.DiscountValue * 100))
		if row.DiscountType == "rate" {
			// 90 means 9 折: the order is reduced by 10%.
			value = int64(math.Round(float64(base) * (100 - row.DiscountValue) / 100))
		}
		if value <= 0 {
			return CouponPricing{}, ErrCouponOwnership
		}
		if value > base {
			value = base
		}
		result.DiscountCents += value
		result.Allocations = append(result.Allocations, CouponAllocation{CouponUserID: row.UserCouponID, DiscountCents: value})
	}
	return result, nil
}

func uniqueCouponIDs(ids []uint64) ([]uint64, error) {
	seen := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			return nil, ErrCouponOwnership
		}
		if _, ok := seen[id]; ok {
			return nil, ErrCouponConflict
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out, nil
}

func lockCouponsToOrder(tx *gorm.DB, userID, groupOrderID uint64, pricing CouponPricing) error {
	if len(pricing.Allocations) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(pricing.Allocations))
	for _, item := range pricing.Allocations {
		ids = append(ids, item.CouponUserID)
	}
	result := tx.Table("qixi_crm_b_coupon_user").Where("user_id = ? AND id IN ? AND status = ?", userID, ids, "unused").Updates(map[string]any{"status": "locked", "used_order_id": groupOrderID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != int64(len(ids)) {
		return ErrCouponOwnership
	}
	return nil
}
