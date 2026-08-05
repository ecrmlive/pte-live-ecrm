package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ChangeMobile atomically verifies possession of the bound and replacement
// numbers, consumes both purpose-bound codes, then changes the mobile value.
// Codes are rolled back when the new number is already occupied or the user
// state changed, so a recoverable conflict does not waste a verification code.
func (s *Service) ChangeMobile(ctx context.Context, userID uint64, oldMobile, oldCode, newMobile, newCode string) error {
	oldMobile = strings.TrimSpace(oldMobile)
	newMobile = strings.TrimSpace(newMobile)
	oldCode = strings.TrimSpace(oldCode)
	newCode = strings.TrimSpace(newCode)
	if userID == 0 || !validMobile(oldMobile) || !validMobile(newMobile) || oldMobile == newMobile || len(oldCode) != 6 || len(newCode) != 6 {
		return ErrBadParam
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if user.Status != 1 || user.Mobile == nil || *user.Mobile != oldMobile {
			return ErrBadParam
		}

		now := time.Now()
		consume := func(mobile, code string) error {
			result := tx.Model(&SMSCode{}).Where("mobile=? AND purpose=? AND code_hash=? AND consumed_at IS NULL AND expires_at>?", mobile, "change_mobile", smsHash(mobile, "change_mobile", code), now).Update("consumed_at", now)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return ErrSMSInvalid
			}
			return nil
		}
		// A stable order prevents a two-session old/new-code lock inversion.
		if oldMobile < newMobile {
			if err := consume(oldMobile, oldCode); err != nil {
				return err
			}
			if err := consume(newMobile, newCode); err != nil {
				return err
			}
		} else {
			if err := consume(newMobile, newCode); err != nil {
				return err
			}
			if err := consume(oldMobile, oldCode); err != nil {
				return err
			}
		}

		var other User
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("mobile = ?", newMobile).First(&other).Error
		if err == nil && other.ID != user.ID {
			return ErrAccountExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		result := tx.Model(&User{}).Where("id = ? AND status = ? AND mobile = ?", user.ID, 1, oldMobile).Update("mobile", newMobile)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrNotFound
		}
		return nil
	})
}
