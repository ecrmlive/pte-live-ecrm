package merchantdeposit

import (
	"errors"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db} }
func (h *Handler) Register(r gin.IRoutes) {
	p := middleware.RequireAdminRoles("platform")
	m := middleware.RequireAdminMenu(h.db, "merchant.deposit.review")
	r.GET("/merchant-deposits", p, m, h.List)
	r.GET("/merchant-deposit-refunds", p, m, h.ListRefunds)
	r.POST("/merchant-deposits/:merchant_id/deduct", p, m, h.Deduct)
	r.POST("/merchant-deposit-refunds/:id/approve", p, m, h.Approve)
	r.POST("/merchant-deposit-refunds/:id/reject", p, m, h.Reject)
	r.POST("/merchant-deposit-refunds/:id/mark-paid", p, m, h.MarkPaid)
}

type account struct {
	MerchantID uint    `gorm:"column:merchant_id" json:"merchant_id"`
	Required   float64 `gorm:"column:required_amount" json:"required_amount"`
	Available  float64 `gorm:"column:available_amount" json:"available_amount"`
	State      string  `gorm:"column:state" json:"state"`
}

func (h *Handler) List(c *gin.Context) {
	var rows []account
	if e := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_deposit_account").Order("merchant_id").Find(&rows).Error; e != nil {
		fail(c, e)
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func (h *Handler) ListRefunds(c *gin.Context) {
	var rows []struct {
		ID              uint      `json:"id"`
		MerchantID      uint      `json:"merchant_id"`
		Amount          float64   `json:"amount"`
		Status          string    `json:"status"`
		Reason          string    `json:"reason"`
		ReviewNote      string    `json:"review_note"`
		PayoutReference *string   `json:"payout_reference"`
		CreatedAt       time.Time `json:"created_at"`
	}
	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_deposit_refund").Order("id DESC")
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		q = q.Where("status=?", status)
	}
	if e := q.Find(&rows).Error; e != nil {
		fail(c, e)
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func (h *Handler) Deduct(c *gin.Context) {
	mid, e := strconv.ParseUint(c.Param("merchant_id"), 10, 64)
	var q struct {
		Amount         float64 `json:"amount"`
		Reason         string  `json:"reason"`
		IdempotencyKey string  `json:"idempotency_key"`
	}
	if e != nil || mid == 0 || c.ShouldBindJSON(&q) != nil || !validDepositAmount(q.Amount) || strings.TrimSpace(q.Reason) == "" || strings.TrimSpace(q.IdempotencyKey) == "" {
		response.Fail(c, 400, "扣减参数错误")
		return
	}
	e = h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var a account
		if e := tx.Table("qixi_crm_a_merchant_deposit_account").Clauses(clause.Locking{Strength: "UPDATE"}).Where("merchant_id=?", mid).Take(&a).Error; e != nil {
			return e
		}
		var existing struct {
			Amount float64
			Reason string
		}
		err := tx.Table("qixi_crm_a_merchant_deposit_ledger").Select("amount, reason").Where("merchant_id=? AND idempotency_key=?", mid, strings.TrimSpace(q.IdempotencyKey)).Take(&existing).Error
		if err == nil {
			if sameDepositAmount(existing.Amount, q.Amount) && existing.Reason == strings.TrimSpace(q.Reason) {
				return nil
			}
			return errors.New("保证金扣减幂等键冲突")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if !depositCanDeduct(a.State) {
			return errors.New("当前保证金不可扣减")
		}
		if a.Available < q.Amount {
			return errors.New("保证金余额不足")
		}
		bal := a.Available - q.Amount
		state := depositStateAfterBalance(a.Required, bal)
		if e := tx.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id=?", mid).Updates(map[string]any{"available_amount": bal, "state": state, "version": gorm.Expr("version + 1")}).Error; e != nil {
			return e
		}
		return tx.Table("qixi_crm_a_merchant_deposit_ledger").Create(map[string]any{"merchant_id": mid, "entry_type": "deduct", "amount": q.Amount, "balance_after": bal, "reason": strings.TrimSpace(q.Reason), "idempotency_key": strings.TrimSpace(q.IdempotencyKey), "operator_admin_id": middleware.AdminID(c)}).Error
	})
	if e != nil {
		if depositConflict(e) {
			response.Fail(c, http.StatusConflict, e.Error())
			return
		}
		fail(c, e)
		return
	}
	response.OK(c, gin.H{"ok": true})
}
func (h *Handler) Approve(c *gin.Context) { h.review(c, "approved") }
func (h *Handler) Reject(c *gin.Context)  { h.review(c, "rejected") }
func (h *Handler) review(c *gin.Context, status string) {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	var q struct {
		Note string `json:"note"`
	}
	if e != nil || id == 0 || c.ShouldBindJSON(&q) != nil {
		response.Fail(c, 400, "退款审核参数错误")
		return
	}
	q.Note = strings.TrimSpace(q.Note)
	if q.Note == "" || len([]rune(q.Note)) > 500 {
		response.Fail(c, 400, "退款审核参数错误")
		return
	}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var refund struct {
			ID         uint
			MerchantID uint `gorm:"column:merchant_id"`
			Status     string
		}
		if err := tx.Table("qixi_crm_a_merchant_deposit_refund").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).Take(&refund).Error; err != nil {
			return err
		}
		if refund.Status != "applied" {
			return errors.New("退款状态已变化")
		}
		var acc account
		if err := tx.Table("qixi_crm_a_merchant_deposit_account").Clauses(clause.Locking{Strength: "UPDATE"}).Where("merchant_id=?", refund.MerchantID).Take(&acc).Error; err != nil {
			return err
		}
		if status == "approved" {
			if !depositCanApproveRefund(acc.State) {
				return errors.New("当前保证金不可进入退款流程")
			}
			if err := tx.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id=?", refund.MerchantID).Updates(map[string]any{"state": "refund_pending", "version": gorm.Expr("version + 1")}).Error; err != nil {
				return err
			}
		}
		res := tx.Table("qixi_crm_a_merchant_deposit_refund").Where("id=? AND status='applied'", id).Updates(map[string]any{"status": status, "review_note": q.Note, "reviewed_by": middleware.AdminID(c), "reviewed_at": time.Now()})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errors.New("退款状态已变化")
		}
		reason := "保证金退款审核通过：" + q.Note
		if status == "rejected" {
			reason = "保证金退款审核驳回：" + q.Note
		}
		return tx.Table("qixi_crm_a_merchant_deposit_ledger").Create(map[string]any{
			"merchant_id":       refund.MerchantID,
			"entry_type":        "refund_" + status,
			"amount":            0,
			"balance_after":     acc.Available,
			"reason":            reason,
			"idempotency_key":   "refund-review-" + strconv.FormatUint(uint64(refund.ID), 10) + "-" + status,
			"operator_admin_id": middleware.AdminID(c),
		}).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "保证金退款申请不存在")
		} else if depositConflict(err) {
			response.Fail(c, http.StatusConflict, err.Error())
		} else {
			fail(c, err)
		}
		return
	}
	response.OK(c, gin.H{"ok": true, "status": status})
}
func (h *Handler) MarkPaid(c *gin.Context) {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	var q struct {
		IdempotencyKey  string `json:"idempotency_key"`
		PayoutReference string `json:"payout_reference"`
	}
	if e != nil || id == 0 || c.ShouldBindJSON(&q) != nil || !validPayoutRegistrationInput(q.IdempotencyKey, q.PayoutReference) {
		response.Fail(c, 400, "打款登记参数错误")
		return
	}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var row struct {
			ID                   uint
			MerchantID           uint `gorm:"column:merchant_id"`
			Status               string
			PayoutIdempotencyKey *string `gorm:"column:payout_idempotency_key"`
			PayoutReference      string  `gorm:"column:payout_reference"`
			Amount               float64
		}
		if err := tx.Table("qixi_crm_a_merchant_deposit_refund").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).Take(&row).Error; err != nil {
			return err
		}
		if row.PayoutIdempotencyKey != nil && *row.PayoutIdempotencyKey == strings.TrimSpace(q.IdempotencyKey) && row.Status == "paid" {
			if row.PayoutReference == strings.TrimSpace(q.PayoutReference) {
				return nil
			}
			return errors.New("退款状态已变化或幂等键冲突")
		}
		if row.PayoutIdempotencyKey != nil || row.Status != "approved" {
			return errors.New("退款状态已变化或幂等键冲突")
		}
		var acc account
		if err := tx.Table("qixi_crm_a_merchant_deposit_account").Clauses(clause.Locking{Strength: "UPDATE"}).Where("merchant_id = ?", row.MerchantID).Take(&acc).Error; err != nil {
			return err
		}
		if acc.Available < row.Amount {
			return errors.New("保证金余额不足，不能登记退款打款")
		}
		var duplicate int64
		if err := tx.Table("qixi_crm_a_merchant_deposit_refund").Where("merchant_id=? AND payout_idempotency_key=?", row.MerchantID, strings.TrimSpace(q.IdempotencyKey)).Count(&duplicate).Error; err != nil {
			return err
		}
		if duplicate > 0 {
			return errors.New("退款状态已变化或幂等键冲突")
		}
		balance := acc.Available - row.Amount
		state := depositStateAfterRefundBalance(acc.Required, balance)
		if err := tx.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id = ?", row.MerchantID).Updates(map[string]any{"available_amount": balance, "state": state, "version": gorm.Expr("version + 1")}).Error; err != nil {
			return err
		}
		res := tx.Table("qixi_crm_a_merchant_deposit_refund").Where("id=? AND status='approved' AND payout_idempotency_key IS NULL", id).Updates(map[string]any{"status": "paid", "payout_idempotency_key": strings.TrimSpace(q.IdempotencyKey), "payout_reference": strings.TrimSpace(q.PayoutReference), "paid_by": middleware.AdminID(c), "paid_at": time.Now()})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errors.New("退款状态已变化或幂等键冲突")
		}
		return tx.Table("qixi_crm_a_merchant_deposit_ledger").Create(map[string]any{"merchant_id": row.MerchantID, "entry_type": "refund_paid", "amount": row.Amount, "balance_after": balance, "reason": "保证金退款已登记打款", "idempotency_key": "refund-paid-" + strconv.FormatUint(uint64(row.ID), 10), "operator_admin_id": middleware.AdminID(c)}).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "保证金退款申请不存在")
		} else if depositConflict(err) {
			response.Fail(c, http.StatusConflict, err.Error())
		} else {
			fail(c, err)
		}
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func depositCanDeduct(state string) bool {
	return state == "funded" || state == "shortfall"
}

func depositCanApproveRefund(state string) bool {
	return state == "funded" || state == "shortfall"
}

func depositStateAfterBalance(required, balance float64) string {
	if balance < required {
		return "shortfall"
	}
	return "funded"
}

func depositStateAfterRefundBalance(required, balance float64) string {
	if balance == 0 {
		return "refunded"
	}
	return depositStateAfterBalance(required, balance)
}

func validPayoutRegistrationInput(key, reference string) bool {
	key = strings.TrimSpace(key)
	reference = strings.TrimSpace(reference)
	return len([]rune(key)) >= 8 && len([]rune(key)) <= 128 && len([]rune(reference)) > 0 && len([]rune(reference)) <= 128
}

func validDepositAmount(amount float64) bool {
	if math.IsNaN(amount) || math.IsInf(amount, 0) || amount <= 0 || amount > 9999999999.99 {
		return false
	}
	return sameDepositAmount(amount, math.Round(amount*100)/100)
}

func sameDepositAmount(left, right float64) bool {
	return math.Abs(left-right) < 0.000001
}

func fail(c *gin.Context, e error) {
	if errors.Is(e, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "保证金账户不存在")
		return
	}
	if strings.Contains(e.Error(), "保证金") {
		response.Fail(c, 400, e.Error())
		return
	}
	response.Fail(c, 500, "保证金操作失败")
}

func depositConflict(err error) bool {
	message := err.Error()
	return strings.Contains(message, "幂等") || strings.Contains(message, "状态") || strings.Contains(message, "退款流程") || strings.Contains(message, "余额不足")
}
