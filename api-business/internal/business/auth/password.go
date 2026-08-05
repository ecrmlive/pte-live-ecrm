package auth

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const minPasswordLength = 8

func validNewPassword(value string) bool {
	return len([]rune(value)) >= minPasswordLength && len([]rune(value)) <= 128
}

// ChangePassword verifies the credential attached to the current login
// channel, changes only that credential, then increments auth_version so all
// previously issued access and refresh tokens are invalidated.
func (s *Service) ChangePassword(ctx context.Context, userID uint64, channel LoginChannel, currentPassword, newPassword string) (*Profile, error) {
	currentPassword = strings.TrimSpace(currentPassword)
	if currentPassword == "" || !validNewPassword(newPassword) || currentPassword == newPassword {
		return nil, ErrBadParam
	}
	var profile Profile
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if user.Status != 1 {
			return ErrAccountDisabled
		}
		var identity Identity
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND channel = ?", userID, channel).First(&identity).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidCredentials
			}
			return err
		}
		if identity.CredentialHash == "" || bcrypt.CompareHashAndPassword([]byte(identity.CredentialHash), []byte(currentPassword)) != nil {
			return ErrInvalidCredentials
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err := tx.Model(&Identity{}).Where("id = ?", identity.ID).Update("credential_hash", string(hash)).Error; err != nil {
			return err
		}
		result := tx.Model(&User{}).Where("id = ? AND auth_version = ?", user.ID, user.AuthVersion).Update("auth_version", user.AuthVersion+1)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrInvalidCredentials
		}
		profile = Profile{ID: user.ID, UID: user.ID, Nickname: user.Nickname, Mobile: mobileText(user.Mobile), Channel: channel, Subject: identity.Subject, AuthVersion: user.AuthVersion + 1}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
