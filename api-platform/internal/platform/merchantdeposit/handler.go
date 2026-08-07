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
	r.GET("/merchant-deposits/:merchant_id/ledgers", p, m, h.ListLedgers)
	r.GET("/merchant-deposit-refunds", p, m, h.ListRefunds)
	r.POST("/merchant-deposits/:merchant_id/deduct", p, m, h.Deduct)
	r.POST("/merchant-deposits/:merchant_id/fund-offline", p, m, h.FundOffline)
	r.POST("/merchant-deposit-refunds/:id/approve", p, m, h.Approve)
	r.POST("/merchant-deposit-refunds/:id/reject", p, m, h.Reject)
	r.POST("/merchant-deposit-refunds/:id/mark", p, m, h.MarkRefundNote)
	r.POST("/merchant-deposit-refunds/:id/mark-paid", p, m, h.MarkPaid)
}

type account struct {
	MerchantID uint    `gorm:"column:merchant_id" json:"merchant_id"`
	Required   float64 `gorm:"column:required_amount" json:"required_amount"`
	Available  float64 `gorm:"column:available_amount" json:"available_amount"`
	State      string  `gorm:"column:state" json:"state"`
}

type accountRow struct {
	MerchantID     uint       `gorm:"column:merchant_id" json:"merchant_id"`
	MerchantName   string     `gorm:"column:merchant_name" json:"merchant_name"`
	OwnerName      string     `gorm:"column:owner_name" json:"owner_name"`
	TypeName       string     `gorm:"column:type_name" json:"type_name"`
	CategoryName   string     `gorm:"column:category_name" json:"category_name"`
	RequiredAmount float64    `gorm:"column:required_amount" json:"required_amount"`
	AvailableAmount float64   `gorm:"column:available_amount" json:"available_amount"`
	PayableAmount  float64    `gorm:"column:payable_amount" json:"payable_amount"`
	State          string     `gorm:"column:state" json:"state"`
	Mark           string     `gorm:"column:mark" json:"mark"`
	IsTrader       int        `gorm:"column:is_trader" json:"is_trader"`
	TypeID         uint       `gorm:"column:type_id" json:"type_id"`
	CategoryID     uint       `gorm:"column:category_id" json:"category_id"`
	PaidAt         *time.Time `gorm:"column:paid_at" json:"paid_at"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := pageLimit(c)
	tab := strings.TrimSpace(c.Query("tab"))
	if tab == "" {
		if status := strings.TrimSpace(c.Query("status")); status != "" {
			switch status {
			case "pending":
				tab = "pending"
			case "funded", "shortfall":
				// shortfall 归属缴存列表（对齐 CRMEB is_margin=10）
				tab = "funded"
			default:
				tab = status
			}
		} else {
			tab = "funded"
		}
	}
	q := h.db.WithContext(c.Request.Context()).
		Table("qixi_crm_a_merchant_deposit_account AS a").
		Joins("LEFT JOIN qixi_crm_a_merchant_view AS m ON m.merchant_id = a.merchant_id").
		Joins("LEFT JOIN qixi_crm_a_merchant_type AS t ON t.id = m.type_id").
		Joins("LEFT JOIN qixi_crm_a_merchant_category AS c ON c.id = m.category_id")
	switch tab {
	case "pending":
		// CRMEB 待缴：尚未缴纳（is_margin=1）。shortfall 仍留在缴存列表，可继续扣费。
		q = q.Where("a.state = ?", "pending")
	case "funded":
		// CRMEB 缴存：已缴纳过保证金（is_margin=10），含扣费后不足额度的 shortfall。
		q = q.Where("a.state IN ?", []string{"funded", "shortfall"})
	default:
		if tab != "" && tab != "all" {
			q = q.Where("a.state = ?", tab)
		}
	}
	if mid := strings.TrimSpace(c.Query("merchant_id")); mid != "" {
		id, err := strconv.ParseUint(mid, 10, 64)
		if err != nil || id == 0 {
			response.Fail(c, 400, "商户 ID 参数错误")
			return
		}
		q = q.Where("a.merchant_id = ?", id)
	}
	if typeID := strings.TrimSpace(c.Query("type_id")); typeID != "" {
		id, err := strconv.ParseUint(typeID, 10, 64)
		if err != nil || id == 0 {
			response.Fail(c, 400, "店铺类型参数错误")
			return
		}
		q = q.Where("m.type_id = ?", id)
	}
	if categoryID := strings.TrimSpace(c.Query("category_id")); categoryID != "" {
		id, err := strconv.ParseUint(categoryID, 10, 64)
		if err != nil || id == 0 {
			response.Fail(c, 400, "店铺分类参数错误")
			return
		}
		q = q.Where("m.category_id = ?", id)
	}
	if trader := strings.TrimSpace(c.Query("is_trader")); trader != "" {
		v, err := strconv.Atoi(trader)
		if err != nil || (v != 0 && v != 1) {
			response.Fail(c, 400, "店铺类别参数错误")
			return
		}
		q = q.Where("m.is_trader = ?", v)
	}
	if kw := strings.TrimSpace(c.Query("keyword")); kw != "" {
		like := "%" + kw + "%"
		q = q.Where(
			"CAST(a.merchant_id AS CHAR) LIKE ? OR COALESCE(m.merchant_name,'') LIKE ? OR COALESCE(m.contact_mobile,'') LIKE ? OR COALESCE(m.owner_name,'') LIKE ? OR COALESCE(m.contact_name,'') LIKE ?",
			like, like, like, like, like,
		)
	}
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		q = q.Where("a.updated_at >= ?", from+" 00:00:00")
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		q = q.Where("a.updated_at <= ?", to+" 23:59:59")
	}
	var total int64
	if e := q.Count(&total).Error; e != nil {
		fail(c, e)
		return
	}
	var rows []accountRow
	selectSQL := `a.merchant_id, a.required_amount, a.available_amount, a.state,
		GREATEST(a.required_amount - a.available_amount, 0) AS payable_amount,
		COALESCE(m.merchant_name,'') AS merchant_name,
		COALESCE(NULLIF(m.owner_name,''), NULLIF(m.contact_name,''), '') AS owner_name,
		COALESCE(t.name,'') AS type_name,
		COALESCE(c.name,'') AS category_name,
		COALESCE(m.mark,'') AS mark,
		COALESCE(m.is_trader,0) AS is_trader,
		COALESCE(m.type_id,0) AS type_id,
		COALESCE(m.category_id,0) AS category_id,
		(SELECT MAX(l.created_at) FROM qixi_crm_a_merchant_deposit_ledger l WHERE l.merchant_id = a.merchant_id AND l.entry_type = 'fund') AS paid_at`
	if e := q.Select(selectSQL).Order("a.merchant_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; e != nil {
		fail(c, e)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit, "tab": tab})
}

// ListLedgers returns immutable deposit operation records for one merchant (CRMEB「操作记录」).
func (h *Handler) ListLedgers(c *gin.Context) {
	mid, err := strconv.ParseUint(c.Param("merchant_id"), 10, 64)
	if err != nil || mid == 0 {
		response.Fail(c, 400, "商户 ID 参数错误")
		return
	}
	page, limit := pageLimit(c)
	base := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_deposit_ledger AS l").Where("l.merchant_id = ?", mid)
	var total int64
	if e := base.Count(&total).Error; e != nil {
		response.Fail(c, 500, "操作记录加载失败")
		return
	}
	var rows []struct {
		ID              uint      `json:"id"`
		MerchantID      uint      `gorm:"column:merchant_id" json:"merchant_id"`
		EntryType       string    `gorm:"column:entry_type" json:"entry_type"`
		Amount          float64   `json:"amount"`
		BalanceAfter    float64   `gorm:"column:balance_after" json:"balance_after"`
		Reason          string    `json:"reason"`
		OperatorAdminID uint      `gorm:"column:operator_admin_id" json:"operator_admin_id"`
		OperatorName    string    `gorm:"column:operator_name" json:"operator_name"`
		CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	}
	if e := h.db.WithContext(c.Request.Context()).
		Table("qixi_crm_a_merchant_deposit_ledger AS l").
		Select(`l.id, l.merchant_id, l.entry_type, l.amount, l.balance_after, l.reason, l.operator_admin_id,
			COALESCE(NULLIF(u.display_name,''), NULLIF(u.username,''), '') AS operator_name, l.created_at`).
		Joins("LEFT JOIN qixi_crm_a_admin_user AS u ON u.id = l.operator_admin_id AND u.deleted_at IS NULL").
		Where("l.merchant_id = ?", mid).
		Order("l.id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&rows).Error; e != nil {
		response.Fail(c, 500, "操作记录加载失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) ListRefunds(c *gin.Context) {
	page, limit := pageLimit(c)
	var rows []struct {
		ID               uint      `json:"id"`
		MerchantID       uint      `gorm:"column:merchant_id" json:"merchant_id"`
		MerchantName     string    `gorm:"column:merchant_name" json:"merchant_name"`
		OwnerName        string    `gorm:"column:owner_name" json:"owner_name"`
		RequiredAmount   float64   `gorm:"column:required_amount" json:"required_amount"`
		AvailableAmount  float64   `gorm:"column:available_amount" json:"available_amount"`
		Amount           float64   `json:"amount"`
		Status           string    `json:"status"`
		Reason           string    `json:"reason"`
		ReviewNote       string    `gorm:"column:review_note" json:"review_note"`
		RefundMethod     string    `gorm:"column:refund_method" json:"refund_method"`
		PayoutReference  *string   `json:"payout_reference"`
		CreatedAt        time.Time `json:"created_at"`
	}
	q := h.db.WithContext(c.Request.Context()).
		Table("qixi_crm_a_merchant_deposit_refund AS r").
		Joins("LEFT JOIN qixi_crm_a_merchant_view AS m ON m.merchant_id = r.merchant_id").
		Joins("LEFT JOIN qixi_crm_a_merchant_deposit_account AS a ON a.merchant_id = r.merchant_id")
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		q = q.Where("r.status=?", status)
	}
	if mid := strings.TrimSpace(c.Query("merchant_id")); mid != "" {
		id, err := strconv.ParseUint(mid, 10, 64)
		if err != nil || id == 0 {
			response.Fail(c, 400, "商户 ID 参数错误")
			return
		}
		q = q.Where("r.merchant_id = ?", id)
	}
	if typeID := strings.TrimSpace(c.Query("type_id")); typeID != "" {
		id, err := strconv.ParseUint(typeID, 10, 64)
		if err != nil || id == 0 {
			response.Fail(c, 400, "店铺类型参数错误")
			return
		}
		q = q.Where("m.type_id = ?", id)
	}
	if categoryID := strings.TrimSpace(c.Query("category_id")); categoryID != "" {
		id, err := strconv.ParseUint(categoryID, 10, 64)
		if err != nil || id == 0 {
			response.Fail(c, 400, "店铺分类参数错误")
			return
		}
		q = q.Where("m.category_id = ?", id)
	}
	if trader := strings.TrimSpace(c.Query("is_trader")); trader != "" {
		v, err := strconv.Atoi(trader)
		if err != nil || (v != 0 && v != 1) {
			response.Fail(c, 400, "店铺类别参数错误")
			return
		}
		q = q.Where("m.is_trader = ?", v)
	}
	if kw := strings.TrimSpace(c.Query("keyword")); kw != "" {
		like := "%" + kw + "%"
		q = q.Where(
			"CAST(r.merchant_id AS CHAR) LIKE ? OR COALESCE(m.merchant_name,'') LIKE ? OR COALESCE(m.contact_mobile,'') LIKE ? OR r.reason LIKE ? OR COALESCE(r.review_note,'') LIKE ?",
			like, like, like, like, like,
		)
	}
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		q = q.Where("r.created_at >= ?", from+" 00:00:00")
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		q = q.Where("r.created_at <= ?", to+" 23:59:59")
	}
	var total int64
	if e := q.Count(&total).Error; e != nil {
		fail(c, e)
		return
	}
	selectSQL := `r.id, r.merchant_id, r.amount, r.status, r.reason, r.review_note, r.payout_reference, r.created_at,
		COALESCE(m.merchant_name,'') AS merchant_name,
		COALESCE(NULLIF(m.owner_name,''), NULLIF(m.contact_name,''), '') AS owner_name,
		COALESCE(a.required_amount,0) AS required_amount,
		COALESCE(a.available_amount,0) AS available_amount,
		'线下' AS refund_method`
	if e := q.Select(selectSQL).Order("r.id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; e != nil {
		fail(c, e)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func pageLimit(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
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

// FundOffline registers an offline margin top-up (CRMEB「线下付款」).
func (h *Handler) FundOffline(c *gin.Context) {
	mid, e := strconv.ParseUint(c.Param("merchant_id"), 10, 64)
	var q struct {
		Amount         float64 `json:"amount"`
		Mark           string  `json:"mark"`
		IdempotencyKey string  `json:"idempotency_key"`
	}
	if e != nil || mid == 0 || c.ShouldBindJSON(&q) != nil || strings.TrimSpace(q.IdempotencyKey) == "" {
		response.Fail(c, 400, "线下付款参数错误")
		return
	}
	q.Mark = strings.TrimSpace(q.Mark)
	if len([]rune(q.Mark)) > 500 {
		response.Fail(c, 400, "备注不能超过 500 个字符")
		return
	}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var a account
		if e := tx.Table("qixi_crm_a_merchant_deposit_account").Clauses(clause.Locking{Strength: "UPDATE"}).Where("merchant_id=?", mid).Take(&a).Error; e != nil {
			return e
		}
		if !depositCanFundOffline(a.State) {
			return errors.New("当前保证金不可线下付款")
		}
		payable := a.Required - a.Available
		if payable < 0 {
			payable = 0
		}
		amount := q.Amount
		if amount == 0 {
			amount = payable
		}
		if !validDepositAmount(amount) {
			return errors.New("线下付款金额错误")
		}
		if amount > payable+0.000001 {
			return errors.New("线下付款金额超过待缴金额")
		}
		var existing struct {
			Amount float64
			Reason string
		}
		err := tx.Table("qixi_crm_a_merchant_deposit_ledger").Select("amount, reason").Where("merchant_id=? AND idempotency_key=?", mid, strings.TrimSpace(q.IdempotencyKey)).Take(&existing).Error
		if err == nil {
			if sameDepositAmount(existing.Amount, amount) {
				return nil
			}
			return errors.New("保证金线下付款幂等键冲突")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		bal := a.Available + amount
		state := depositStateAfterBalance(a.Required, bal)
		if e := tx.Table("qixi_crm_a_merchant_deposit_account").Where("merchant_id=?", mid).Updates(map[string]any{
			"available_amount": bal,
			"state":            state,
			"version":          gorm.Expr("version + 1"),
		}).Error; e != nil {
			return e
		}
		reason := "线下补缴保证金"
		if q.Mark != "" {
			reason = q.Mark
			_ = tx.Table("qixi_crm_a_merchant_view").Where("merchant_id=?", mid).Update("mark", q.Mark).Error
		}
		return tx.Table("qixi_crm_a_merchant_deposit_ledger").Create(map[string]any{
			"merchant_id":       mid,
			"entry_type":        "fund",
			"amount":            amount,
			"balance_after":     bal,
			"reason":            reason,
			"idempotency_key":   strings.TrimSpace(q.IdempotencyKey),
			"operator_admin_id": middleware.AdminID(c),
		}).Error
	})
	if err != nil {
		if depositConflict(err) {
			response.Fail(c, http.StatusConflict, err.Error())
			return
		}
		fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) MarkRefundNote(c *gin.Context) {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	var q struct {
		Note string `json:"note"`
	}
	if e != nil || id == 0 || c.ShouldBindJSON(&q) != nil {
		response.Fail(c, 400, "备注参数错误")
		return
	}
	q.Note = strings.TrimSpace(q.Note)
	if len([]rune(q.Note)) > 500 {
		response.Fail(c, 400, "备注不能超过 500 个字符")
		return
	}
	res := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_deposit_refund").Where("id=?", id).Update("review_note", q.Note)
	if res.Error != nil {
		fail(c, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "保证金退款申请不存在")
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

func depositCanFundOffline(state string) bool {
	return state == "pending" || state == "shortfall"
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
	return strings.Contains(message, "幂等") || strings.Contains(message, "状态") || strings.Contains(message, "退款流程") || strings.Contains(message, "余额不足") || strings.Contains(message, "线下付款") || strings.Contains(message, "待缴金额")
}
