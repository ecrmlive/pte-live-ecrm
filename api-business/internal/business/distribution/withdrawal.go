package distribution

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	errNotPromoter            = errors.New("not promoter")
	errInsufficientCommission = errors.New("insufficient commission")
)

type withdrawalInput struct {
	Amount         float64 `json:"amount"`
	Channel        string  `json:"channel"`
	AccountName    string  `json:"account_name"`
	AccountNo      string  `json:"account_no"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type withdrawalRow struct {
	ID             uint64    `gorm:"column:id"`
	WithdrawalNo   string    `gorm:"column:withdrawal_no"`
	Amount         float64   `gorm:"column:amount"`
	Channel        string    `gorm:"column:channel"`
	Status         string    `gorm:"column:status"`
	ReviewNote     string    `gorm:"column:review_note"`
	IdempotencyKey string    `gorm:"column:idempotency_key"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

func (h *Handler) withdrawals(c *gin.Context) {
	page, limit := pageParams(c)
	uid := uint64(middleware.UID(c))
	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_withdrawal_application").Where("user_id = ?", uid)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	rows := make([]withdrawalRow, 0)
	if err := q.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		fail(c)
		return
	}
	available, reserved, err := withdrawalAmounts(c, h.db, uid, false)
	if err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, gin.H{"id": row.ID, "withdrawal_no": row.WithdrawalNo, "amount": row.Amount, "channel": withdrawalChannelText(row.Channel), "status": row.Status, "status_text": withdrawalStatusText(row.Status), "review_note": row.ReviewNote, "created_at": row.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit, "available_commission": available, "reserved_commission": reserved, "withdrawable_commission": roundMoney(math.Max(0, available-reserved))})
}

func (h *Handler) applyWithdrawal(c *gin.Context) {
	var input withdrawalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, "提现信息格式错误")
		return
	}
	input.Channel = strings.ToLower(strings.TrimSpace(input.Channel))
	input.AccountName = strings.TrimSpace(input.AccountName)
	input.AccountNo = strings.TrimSpace(input.AccountNo)
	input.IdempotencyKey = strings.TrimSpace(input.IdempotencyKey)
	input.Amount = roundMoney(input.Amount)
	if !validWithdrawalInput(input) {
		response.Fail(c, http.StatusBadRequest, "提现信息填写不完整或金额不正确")
		return
	}

	uid := uint64(middleware.UID(c))
	var created withdrawalRow
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var promoterCount int64
		if err := tx.Table("qixi_crm_b_distribution_promoter").Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND status = 1", uid).Count(&promoterCount).Error; err != nil {
			return err
		}
		if promoterCount != 1 {
			return errNotPromoter
		}
		err := tx.Table("qixi_crm_b_withdrawal_application").Where("user_id = ? AND idempotency_key = ?", uid, input.IdempotencyKey).First(&created).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		available, reserved, err := withdrawalAmounts(c, tx, uid, true)
		if err != nil {
			return err
		}
		if roundMoney(available-reserved) < input.Amount {
			return errInsufficientCommission
		}
		snapshot, err := json.Marshal(gin.H{"account_name": input.AccountName, "account_no": input.AccountNo, "channel": input.Channel})
		if err != nil {
			return err
		}
		created = withdrawalRow{WithdrawalNo: withdrawalNo(uid), Amount: input.Amount, Channel: input.Channel, Status: "applied", IdempotencyKey: input.IdempotencyKey}
		if err := tx.Table("qixi_crm_b_withdrawal_application").Create(map[string]any{"withdrawal_no": created.WithdrawalNo, "user_id": uid, "amount": created.Amount, "channel": created.Channel, "account_snapshot": string(snapshot), "status": created.Status, "idempotency_key": created.IdempotencyKey}).Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_b_withdrawal_application").Where("user_id = ? AND idempotency_key = ?", uid, input.IdempotencyKey).First(&created).Error
	})
	if err != nil {
		switch {
		case errors.Is(err, errNotPromoter):
			response.Fail(c, http.StatusBadRequest, "当前账户未开通推广资格")
		case errors.Is(err, errInsufficientCommission):
			response.Fail(c, http.StatusBadRequest, "可提现佣金不足")
		default:
			fail(c)
		}
		return
	}
	response.OK(c, gin.H{"id": created.ID, "withdrawal_no": created.WithdrawalNo, "amount": created.Amount, "status": created.Status, "status_text": withdrawalStatusText(created.Status)})
}

func withdrawalAmounts(c *gin.Context, db *gorm.DB, uid uint64, lock bool) (float64, float64, error) {
	ledger := db.WithContext(c.Request.Context()).Table("qixi_crm_b_commission_ledger").Where("user_id = ? AND status = 'available'", uid)
	applications := db.WithContext(c.Request.Context()).Table("qixi_crm_b_withdrawal_application").Where("user_id = ? AND status IN ('applied','reviewing','approved','paying')", uid)
	if !lock {
		var available, reserved float64
		if err := ledger.Select("COALESCE(SUM(amount), 0)").Scan(&available).Error; err != nil {
			return 0, 0, err
		}
		if err := applications.Select("COALESCE(SUM(amount), 0)").Scan(&reserved).Error; err != nil {
			return 0, 0, err
		}
		return roundMoney(available), roundMoney(reserved), nil
	}
	ledger = ledger.Clauses(clause.Locking{Strength: "UPDATE"})
	applications = applications.Clauses(clause.Locking{Strength: "UPDATE"})
	var commissionRows []commissionRow
	if err := ledger.Select("amount").Find(&commissionRows).Error; err != nil {
		return 0, 0, err
	}
	var applicationsRows []withdrawalRow
	if err := applications.Select("amount").Find(&applicationsRows).Error; err != nil {
		return 0, 0, err
	}
	var available, reserved float64
	for _, row := range commissionRows {
		available += row.Amount
	}
	for _, row := range applicationsRows {
		reserved += row.Amount
	}
	return roundMoney(available), roundMoney(reserved), nil
}

func validWithdrawalInput(input withdrawalInput) bool {
	if input.Amount < 0.01 || input.Amount > 1000000 || len(input.IdempotencyKey) < 12 || len(input.IdempotencyKey) > 128 || utf8.RuneCountInString(input.AccountName) < 2 || utf8.RuneCountInString(input.AccountName) > 64 {
		return false
	}
	if input.Channel == "wechat" {
		return input.AccountNo == ""
	}
	return input.Channel == "bank" && len(input.AccountNo) >= 6 && len(input.AccountNo) <= 64
}

func roundMoney(amount float64) float64 { return math.Round(amount*100) / 100 }

func withdrawalNo(uid uint64) string {
	return fmt.Sprintf("WD%s%06d", fmt.Sprint(time.Now().UTC().UnixNano()), uid%1000000)
}

func withdrawalChannelText(channel string) string {
	if channel == "bank" {
		return "银行卡"
	}
	return "微信零钱"
}

func withdrawalStatusText(status string) string {
	switch status {
	case "applied":
		return "已申请，待审核"
	case "reviewing":
		return "审核中"
	case "approved":
		return "审核通过"
	case "paying":
		return "打款中"
	case "paid":
		return "已打款"
	case "rejected":
		return "已驳回"
	default:
		return "状态未知"
	}
}
