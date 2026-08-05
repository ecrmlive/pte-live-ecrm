package order

import (
	"context"
	"errors"
	"math"
	"time"

	"gorm.io/gorm"
)

// ResolveSVIPUnitCents derives the final unit price from the server-owned
// product projection. Activity prices never stack with SVIP prices. A quoted
// SVIP price that is invalid or not cheaper is deliberately ignored.
func ResolveSVIPUnitCents(listCents int64, productType int8, active bool, priceType int8, configuredPrice float64) int64 {
	if !active || productType != 0 || listCents <= 0 {
		return listCents
	}
	memberCents := int64(0)
	switch priceType {
	case 1:
		// type=1 is the platform's fixed nine-discount rule.
		memberCents = int64(math.Round(float64(listCents) * 0.9))
	case 2:
		value, err := cents(configuredPrice)
		if err != nil {
			return listCents
		}
		memberCents = value
	default:
		return listCents
	}
	if memberCents <= 0 || memberCents >= listCents {
		return listCents
	}
	return memberCents
}

func activeSVIP(ctx context.Context, db *gorm.DB, userID uint64) (bool, error) {
	var row struct {
		Status    string     `gorm:"column:status"`
		ExpiresAt *time.Time `gorm:"column:expires_at"`
	}
	err := db.WithContext(ctx).Table("qixi_crm_b_user_svip").Select("status,expires_at").Where("user_id = ?", userID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return row.Status == "lifetime" || ((row.Status == "trial" || row.Status == "period") && row.ExpiresAt != nil && row.ExpiresAt.After(time.Now().UTC())), nil
}
