package order

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// issueVerificationsForPaidGroup creates unused pickup/service verification
// codes in the same payment transaction. Idempotent: existing rows are kept.
func issueVerificationsForPaidGroup(tx *gorm.DB, groupOrderID uint64) error {
	if tx == nil || groupOrderID == 0 {
		return gorm.ErrInvalidData
	}
	var orders []struct {
		ID           uint64 `gorm:"column:id"`
		DeliveryType string `gorm:"column:delivery_type"`
	}
	err := tx.Table("qixi_crm_b_order AS o").
		Select("o.id, d.delivery_type").
		Joins("JOIN qixi_crm_b_order_delivery AS d ON d.order_id = o.id").
		Where("o.group_order_id = ? AND o.status = ? AND d.delivery_type IN ?", groupOrderID, "paid", []string{"pickup", "service"}).
		Find(&orders).Error
	if err != nil {
		return err
	}
	for _, order := range orders {
		var count int64
		if err := tx.Table("qixi_crm_b_order_verification").Where("order_id = ?", order.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		code := generateVerifyCode()
		if err := tx.Exec(
			`INSERT INTO qixi_crm_b_order_verification (order_id, verify_code,verify_code_hash,status) VALUES (?,?,?,'unused')`,
			order.ID, code, hashVerifyCode(code),
		).Error; err != nil {
			if isDuplicateKey(err) {
				// Rare hash collision or concurrent pay replay; retry once with a new code.
				code = generateVerifyCode()
				if err := tx.Exec(
					`INSERT INTO qixi_crm_b_order_verification (order_id,verify_code,verify_code_hash,status) VALUES (?,?,?,'unused')`,
					order.ID, code, hashVerifyCode(code),
				).Error; err != nil {
					return err
				}
				continue
			}
			return err
		}
	}
	return nil
}

func hashVerifyCode(code string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(code)))
	return hex.EncodeToString(sum[:])
}

func generateVerifyCode() string {
	buf := make([]byte, 5)
	_, _ = rand.Read(buf)
	n := uint64(buf[0])<<32 | uint64(buf[1])<<24 | uint64(buf[2])<<16 | uint64(buf[3])<<8 | uint64(buf[4])
	return fmt.Sprintf("%010d", n%10000000000)
}

func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "uk_verify_code_hash")
}

func normalizeDeliveryType(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "express":
		return "express", nil
	case "pickup", "1", "take":
		return "pickup", nil
	case "city", "local":
		return "city", nil
	case "service":
		return "service", nil
	default:
		return "", ErrDeliveryType
	}
}

var ErrDeliveryType = fmt.Errorf("配送方式不合法")
