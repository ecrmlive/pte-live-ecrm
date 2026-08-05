package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrCancellationConfirmation = errors.New("注销确认已失效，请重新发起")

const accountCancellationConfirmationTTL = 10 * time.Minute

// CancellationBlockers is deliberately limited to business facts that make
// account cancellation consequential. Historical rows are retained either way.
type CancellationBlockers struct {
	Balance          float64 `json:"balance"`
	Points           int64   `json:"points"`
	Commission       float64 `json:"commission"`
	ActiveOrderCount int64   `json:"active_order_count"`
	OpenRefundCount  int64   `json:"open_refund_count"`
}

func (b CancellationBlockers) RequiresConfirmation() bool {
	return b.Balance > 0 || b.Points > 0 || b.Commission > 0 || b.ActiveOrderCount > 0 || b.OpenRefundCount > 0
}

func (b CancellationBlockers) Messages() []string {
	messages := make([]string, 0, 5)
	if b.Balance > 0 {
		messages = append(messages, "账户仍有余额")
	}
	if b.Points > 0 {
		messages = append(messages, "账户仍有积分")
	}
	if b.Commission > 0 {
		messages = append(messages, "账户仍有佣金")
	}
	if b.ActiveOrderCount > 0 {
		messages = append(messages, "存在未完成订单")
	}
	if b.OpenRefundCount > 0 {
		messages = append(messages, "存在处理中的售后")
	}
	return messages
}

type accountCancellationConfirmation struct {
	UserID    uint64    `gorm:"column:user_id;primaryKey"`
	TokenHash string    `gorm:"column:token_hash"`
	Blockers  string    `gorm:"column:blockers"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (accountCancellationConfirmation) TableName() string {
	return "qixi_crm_b_user_cancellation_confirmation"
}

type accountCancellationAudit struct {
	ID                   uint64    `gorm:"column:id;primaryKey"`
	UserID               uint64    `gorm:"column:user_id"`
	ConfirmationRequired bool      `gorm:"column:confirmation_required"`
	Blockers             string    `gorm:"column:blockers"`
	CancelledAt          time.Time `gorm:"column:cancelled_at"`
}

func (accountCancellationAudit) TableName() string { return "qixi_crm_b_user_cancellation_audit" }

type CancellationResult struct {
	Cancelled            bool     `json:"cancelled"`
	RequiresConfirmation bool     `json:"requires_confirmation"`
	ConfirmationToken    string   `json:"confirmation_token,omitempty"`
	Blockers             []string `json:"blockers,omitempty"`
	ExpiresAt            string   `json:"expires_at,omitempty"`
}

func accountCancellationToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func accountCancellationTokenHash(token string) string {
	digest := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(digest[:])
}

func (s *Service) cancellationBlockers(ctx context.Context, tx *gorm.DB, userID uint64) (CancellationBlockers, error) {
	var blockers CancellationBlockers
	if err := tx.WithContext(ctx).Table("qixi_crm_b_member_account").
		Select("COALESCE(balance, 0) AS balance, COALESCE(points, 0) AS points, COALESCE(commission, 0) AS commission").
		Where("user_id = ?", userID).Scan(&blockers).Error; err != nil {
		return CancellationBlockers{}, err
	}
	if err := tx.WithContext(ctx).Table("qixi_crm_b_order").
		Where("user_id = ? AND status NOT IN ?", userID, []string{"completed", "cancelled"}).Count(&blockers.ActiveOrderCount).Error; err != nil {
		return CancellationBlockers{}, err
	}
	if err := tx.WithContext(ctx).Table("qixi_crm_b_refund AS r").
		Joins("INNER JOIN qixi_crm_b_order AS o ON o.id = r.order_id").
		Where("o.user_id = ? AND r.status NOT IN ?", userID, []string{"refunded", "rejected", "cancelled"}).Count(&blockers.OpenRefundCount).Error; err != nil {
		return CancellationBlockers{}, err
	}
	return blockers, nil
}

// CancelAccount logically deactivates an account. It never deletes the user,
// orders, assets or aftersale records, because those are legal and financial
// facts. A fresh confirmation is required whenever mutable obligations remain.
func (s *Service) CancelAccount(ctx context.Context, userID uint64, confirmationToken string) (CancellationResult, error) {
	var result CancellationResult
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
		blockers, err := s.cancellationBlockers(ctx, tx, userID)
		if err != nil {
			return err
		}
		blockerJSON, err := json.Marshal(blockers)
		if err != nil {
			return err
		}
		if blockers.RequiresConfirmation() {
			result.Blockers = blockers.Messages()
			if strings.TrimSpace(confirmationToken) == "" {
				token, err := accountCancellationToken()
				if err != nil {
					return err
				}
				expiresAt := time.Now().Add(accountCancellationConfirmationTTL)
				row := accountCancellationConfirmation{UserID: userID, TokenHash: accountCancellationTokenHash(token), Blockers: string(blockerJSON), ExpiresAt: expiresAt}
				if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "user_id"}}, DoUpdates: clause.AssignmentColumns([]string{"token_hash", "blockers", "expires_at", "updated_at"})}).Create(&row).Error; err != nil {
					return err
				}
				result.RequiresConfirmation = true
				result.ConfirmationToken = token
				result.ExpiresAt = expiresAt.UTC().Format(time.RFC3339)
				return nil
			}
			var confirmation accountCancellationConfirmation
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&confirmation).Error; err != nil {
				return ErrCancellationConfirmation
			}
			if confirmation.ExpiresAt.Before(time.Now()) || subtle.ConstantTimeCompare([]byte(confirmation.TokenHash), []byte(accountCancellationTokenHash(confirmationToken))) != 1 {
				return ErrCancellationConfirmation
			}
		}
		if err := tx.Create(&accountCancellationAudit{UserID: userID, ConfirmationRequired: blockers.RequiresConfirmation(), Blockers: string(blockerJSON), CancelledAt: time.Now()}).Error; err != nil {
			return err
		}
		updated := tx.Model(&User{}).Where("id = ? AND status = ? AND auth_version = ?", userID, 1, user.AuthVersion).Updates(map[string]any{"status": 0, "auth_version": user.AuthVersion + 1})
		if updated.Error != nil {
			return updated.Error
		}
		if updated.RowsAffected != 1 {
			return ErrAccountDisabled
		}
		result.Cancelled = true
		return nil
	})
	return result, err
}
