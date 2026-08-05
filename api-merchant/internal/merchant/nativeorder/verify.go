package nativeorder

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	errVerifyCodeMismatch = errors.New("verify code mismatch")
	errVerifyNotEligible  = errors.New("order not verifiable")
)

type verificationRow struct {
	ID     uint64 `gorm:"column:id"`
	Status string `gorm:"column:status"`
	Hash   string `gorm:"column:verify_code_hash"`
}

func (h *Handler) verify(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "订单 ID 错误")
		return
	}
	var req struct {
		VerifyCode string `json:"verify_code"`
	}
	_ = c.ShouldBindJSON(&req)
	code := strings.TrimSpace(req.VerifyCode)
	adminID := uint64(middleware.AdminID(c))
	storeID := uint64(middleware.StoreID(c))

	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var row struct {
			ID         uint64  `gorm:"column:id"`
			MerchantID uint64  `gorm:"column:merchant_id"`
			StoreID    uint64  `gorm:"column:store_id"`
			PayAmount  float64 `gorm:"column:pay_amount"`
			Status     string  `gorm:"column:status"`
		}
		if err := tx.Table("qixi_crm_b_order").
			Select("id,merchant_id,store_id,pay_amount,status").
			Where("id = ? AND store_id = ?", id, storeID).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Take(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return gorm.ErrRecordNotFound
			}
			return err
		}
		if row.Status == "completed" {
			return nil
		}
		if row.Status != "paid" && row.Status != "fulfilling" {
			return errVerifyNotEligible
		}

		var ver verificationRow
		err := tx.Table("qixi_crm_b_order_verification").
			Select("id,status,verify_code_hash").
			Where("order_id = ? AND status = 'unused'", id).
			Order("id ASC").Limit(1).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Take(&ver).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		hasUnused := err == nil && ver.ID > 0
		if code != "" {
			wantHash := hashVerifyCode(code)
			if !hasUnused || ver.Hash != wantHash {
				return errVerifyCodeMismatch
			}
		}
		if !hasUnused {
			eligible, err := h.verifyEligibleWithoutCode(tx, id)
			if err != nil {
				return err
			}
			if !eligible {
				return errVerifyNotEligible
			}
			plain := randomVerifyCode()
			if err := tx.Exec(
				`INSERT INTO qixi_crm_b_order_verification (order_id,verify_code,verify_code_hash,status,verified_by_account_id,verified_at)
				 VALUES (?,?,?, 'used', ?, NOW())`,
				id, plain, hashVerifyCode(plain), adminID,
			).Error; err != nil {
				return err
			}
		} else {
			result := tx.Exec(
				`UPDATE qixi_crm_b_order_verification SET status='used', verify_code='', verified_by_account_id=?, verified_at=NOW()
				 WHERE id=? AND status='unused'`,
				adminID, ver.ID,
			)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return errVerifyNotEligible
			}
		}

		if err := tx.Table("qixi_crm_b_order").Where("id = ? AND store_id = ?", id, storeID).Update("status", "completed").Error; err != nil {
			return err
		}
		return enqueueSettlementAccrual(tx, row.ID, row.StoreID, row.MerchantID, row.PayAmount)
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "订单不存在")
		return
	}
	if errors.Is(err, errVerifyCodeMismatch) {
		response.Fail(c, http.StatusBadRequest, "核销码不匹配")
		return
	}
	if errors.Is(err, errVerifyNotEligible) {
		response.Fail(c, http.StatusConflict, "当前订单不可核销")
		return
	}
	if err != nil {
		fail(c, "订单核销失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) verifyEligibleWithoutCode(tx *gorm.DB, orderID uint64) (bool, error) {
	var deliveryType string
	err := tx.Table("qixi_crm_b_order_delivery").Select("delivery_type").Where("order_id = ?", orderID).Order("id DESC").Limit(1).Scan(&deliveryType).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return deliveryType == "pickup" || deliveryType == "service", nil
}

func applyVerifyTab(q *gorm.DB, tab string) *gorm.DB {
	switch strings.TrimSpace(tab) {
	case "pending":
		return q.Where(`o.status IN ('paid','fulfilling') AND NOT EXISTS (
			SELECT 1 FROM qixi_crm_b_order_verification v WHERE v.order_id = o.id AND v.status = 'used'
		) AND (
			EXISTS (SELECT 1 FROM qixi_crm_b_order_verification v WHERE v.order_id = o.id AND v.status = 'unused')
			OR EXISTS (SELECT 1 FROM qixi_crm_b_order_delivery d WHERE d.order_id = o.id AND d.delivery_type IN ('pickup','service'))
		)`)
	case "verified":
		return q.Where(`o.status = 'completed' AND EXISTS (
			SELECT 1 FROM qixi_crm_b_order_verification v WHERE v.order_id = o.id AND v.status = 'used'
		)`)
	default:
		return q
	}
}

func hashVerifyCode(code string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(code)))
	return hex.EncodeToString(sum[:])
}

func randomVerifyCode() string {
	buf := make([]byte, 6)
	_, _ = rand.Read(buf)
	n := uint64(buf[0])<<40 | uint64(buf[1])<<32 | uint64(buf[2])<<24 | uint64(buf[3])<<16 | uint64(buf[4])<<8 | uint64(buf[5])
	return strings.ToUpper(hex.EncodeToString([]byte{
		byte(n >> 32), byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n),
	}))[:10]
}

func enqueueSettlementAccrual(tx *gorm.DB, orderID, storeID, merchantID uint64, amount float64) error {
	if orderID == 0 || storeID == 0 || merchantID == 0 || amount <= 0 {
		return nil
	}
	key := "settlement:accrue:" + strconv.FormatUint(orderID, 10)
	return tx.Exec(`INSERT INTO qixi_crm_b_settlement_command_outbox
		(action,order_id,refund_id,store_id,merchant_id,amount,idempotency_key,status)
		VALUES ('accrue',?,0,?,?,?,?,'pending')
		ON DUPLICATE KEY UPDATE updated_at = updated_at`,
		orderID, storeID, merchantID, amount, key,
	).Error
}
