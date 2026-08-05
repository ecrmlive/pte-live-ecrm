package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// grantOnboardingCoupons is deliberately invoked from the same transaction as
// account creation. A unique (user_id, coupon_id) index makes retries safe and
// prevents a second benefit from concurrent registration requests.
func grantOnboardingCoupons(tx *gorm.DB, userID uint64) error {
	var policy struct {
		Enabled       int8 `gorm:"column:enabled"`
		CouponEnabled int8 `gorm:"column:coupon_enabled"`
	}
	err := tx.Table("qixi_crm_b_onboarding_policy").Where("id = ?", 1).First(&policy).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if policy.Enabled == 0 || policy.CouponEnabled == 0 {
		return nil
	}
	now := time.Now()
	var couponIDs []uint64
	if err := tx.Table("qixi_crm_b_onboarding_coupon AS oc").
		Joins("JOIN qixi_crm_b_coupon_template_view AS c ON c.coupon_id = oc.coupon_id").
		Where("oc.enabled = 1 AND c.status = 1").
		Where("(c.starts_at IS NULL OR c.starts_at <= ?)", now).
		Where("(c.ends_at IS NULL OR c.ends_at >= ?)", now).
		Order("oc.sort ASC, oc.coupon_id ASC").Pluck("oc.coupon_id", &couponIDs).Error; err != nil {
		return err
	}
	for _, couponID := range couponIDs {
		if err := tx.Table("qixi_crm_b_coupon_user").Clauses(clause.OnConflict{DoNothing: true}).
			Create(map[string]any{"user_id": userID, "coupon_id": couponID, "source": "newcomer", "status": "unused"}).Error; err != nil {
			return err
		}
	}
	return nil
}
