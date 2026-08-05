package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

var ErrSMSRateLimited = errors.New("验证码发送过于频繁")
var ErrSMSInvalid = errors.New("验证码错误或已失效")

const smsCodeTTL = 5 * time.Minute
const smsResendInterval = time.Minute

type SMSCode struct {
	ID         uint64     `gorm:"column:id;primaryKey"`
	Mobile     string     `gorm:"column:mobile"`
	Purpose    string     `gorm:"column:purpose"`
	CodeHash   string     `gorm:"column:code_hash"`
	ExpiresAt  time.Time  `gorm:"column:expires_at"`
	ConsumedAt *time.Time `gorm:"column:consumed_at"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
}

func (SMSCode) TableName() string { return "qixi_crm_b_auth_sms_code" }

func validSMSPurpose(value string) bool {
	return value == "login" || value == "binding" || value == "reset_password" || value == "change_mobile"
}
func validMobile(value string) bool {
	if len(value) != 11 || value[0] != '1' || value[1] < '3' || value[1] > '9' {
		return false
	}
	for _, c := range value {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
func smsHash(mobile, purpose, code string) string {
	sum := sha256.Sum256([]byte(mobile + ":" + purpose + ":" + code))
	return hex.EncodeToString(sum[:])
}

// CreateSMSCode returns the plaintext only to the server-side gateway caller.
func (s *Service) CreateSMSCode(ctx context.Context, mobile, purpose string) (string, error) {
	mobile = strings.TrimSpace(mobile)
	purpose = strings.TrimSpace(purpose)
	if !validMobile(mobile) || !validSMSPurpose(purpose) {
		return "", ErrBadParam
	}
	now := time.Now()
	var last SMSCode
	err := s.db.WithContext(ctx).Where("mobile=? AND purpose=?", mobile, purpose).Order("id DESC").First(&last).Error
	if err == nil && now.Sub(last.CreatedAt) < smsResendInterval {
		return "", ErrSMSRateLimited
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	var raw [4]byte
	if _, err = rand.Read(raw[:]); err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", (uint32(raw[0])<<24|uint32(raw[1])<<16|uint32(raw[2])<<8|uint32(raw[3]))%1000000)
	row := SMSCode{Mobile: mobile, Purpose: purpose, CodeHash: smsHash(mobile, purpose, code), ExpiresAt: now.Add(smsCodeTTL), CreatedAt: now}
	if err = s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return "", err
	}
	return code, nil
}

// DiscardSMSCode removes a newly issued code when the gateway did not accept it.
// The plaintext code stays server-side and is never persisted or returned.
func (s *Service) DiscardSMSCode(ctx context.Context, mobile, purpose, code string) error {
	mobile = strings.TrimSpace(mobile)
	purpose = strings.TrimSpace(purpose)
	code = strings.TrimSpace(code)
	if !validMobile(mobile) || !validSMSPurpose(purpose) || len(code) != 6 {
		return ErrBadParam
	}
	return s.db.WithContext(ctx).Where("mobile=? AND purpose=? AND code_hash=? AND consumed_at IS NULL", mobile, purpose, smsHash(mobile, purpose, code)).Delete(&SMSCode{}).Error
}

func (s *Service) ConsumeSMSCode(ctx context.Context, mobile, purpose, code string) error {
	mobile = strings.TrimSpace(mobile)
	purpose = strings.TrimSpace(purpose)
	code = strings.TrimSpace(code)
	if !validMobile(mobile) || !validSMSPurpose(purpose) || len(code) != 6 {
		return ErrSMSInvalid
	}
	now := time.Now()
	result := s.db.WithContext(ctx).Model(&SMSCode{}).Where("mobile=? AND purpose=? AND code_hash=? AND consumed_at IS NULL AND expires_at>?", mobile, purpose, smsHash(mobile, purpose, code), now).Update("consumed_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrSMSInvalid
	}
	return nil
}
