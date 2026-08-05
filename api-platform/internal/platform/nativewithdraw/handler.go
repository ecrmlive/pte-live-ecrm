// Package nativewithdraw serves business-owned user withdrawal applications.
// It deliberately does not read legacy qixi_m_admin_* financial tables.
package nativewithdraw

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	businessDB *gorm.DB
	adminDB    *gorm.DB
}

func NewHandler(businessDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/finance/withdraws", h.List)
	r.GET("/finance/withdraws/:id", h.Get)
	platformOnly := middleware.RequireAdminRoles("platform")
	reviewWithdraw := middleware.RequireAdminMenu(h.adminDB, "accounts.withdraw.review")
	r.POST("/finance/withdraws/:id/approve", platformOnly, reviewWithdraw, h.Approve)
	r.POST("/finance/withdraws/:id/reject", platformOnly, reviewWithdraw, h.Reject)
	r.POST("/finance/withdraws/:id/mark-paid", platformOnly, reviewWithdraw, h.MarkPaid)
}

type withdraw struct {
	ID, UserID                    uint64
	WithdrawalNo, Channel, Status string
	Amount                        float64
	ReviewNote                    string
	CreatedAt                     time.Time
	PaidAt                        *time.Time
	PayoutReference               *string
	PayoutIdempotencyKey          *string
}

var (
	errPayoutConflict = errors.New("withdrawal payout state conflict")
	errPayoutNotFound = errors.New("withdrawal payout not found")
)

func (h *Handler) List(c *gin.Context) {
	if !isPlatform(c) {
		response.Fail(c, http.StatusForbidden, "用户提现仅允许平台角色监管")
		return
	}
	page, limit := page(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_withdrawal_application")
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		statuses, ok := statuses(raw)
		if !ok {
			response.Fail(c, http.StatusBadRequest, "提现状态错误")
			return
		}
		q = q.Where("status IN ?", statuses)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	var rows []withdraw
	if err := q.Select("id,user_id,withdrawal_no,channel,status,amount,review_note,payout_reference,payout_idempotency_key,created_at,paid_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, view(row))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Get(c *gin.Context) {
	if !isPlatform(c) {
		response.Fail(c, http.StatusForbidden, "用户提现仅允许平台角色监管")
		return
	}
	id := id(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "提现 ID 错误")
		return
	}
	var row withdraw
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_withdrawal_application").Select("id,user_id,withdrawal_no,channel,status,amount,review_note,payout_reference,payout_idempotency_key,created_at,paid_at").Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, http.StatusNotFound, "提现申请不存在")
		return
	}
	if err != nil {
		fail(c)
		return
	}
	response.OK(c, view(row))
}

func (h *Handler) Approve(c *gin.Context) {
	h.transition(c, "approved", "", []string{"applied", "reviewing"})
}

func (h *Handler) Reject(c *gin.Context) {
	var req struct {
		Refusal string `json:"refusal"`
	}
	if c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Refusal) == "" || len([]rune(strings.TrimSpace(req.Refusal))) > 500 {
		response.Fail(c, http.StatusBadRequest, "拒绝原因不能为空且不能超过 500 字")
		return
	}
	h.transition(c, "rejected", strings.TrimSpace(req.Refusal), []string{"applied", "reviewing"})
}

func (h *Handler) MarkPaid(c *gin.Context) {
	withdrawalID := id(c)
	var req struct {
		IdempotencyKey  string `json:"idempotency_key"`
		PayoutReference string `json:"payout_reference"`
	}
	if withdrawalID == 0 || c.ShouldBindJSON(&req) != nil {
		response.Fail(c, http.StatusBadRequest, "打款登记参数错误")
		return
	}
	req.IdempotencyKey, req.PayoutReference = strings.TrimSpace(req.IdempotencyKey), strings.TrimSpace(req.PayoutReference)
	if !validPayoutInput(req.IdempotencyKey, req.PayoutReference) {
		response.Fail(c, http.StatusBadRequest, "打款凭证编号或幂等键错误")
		return
	}
	var out withdraw
	err := h.businessDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var row withdraw
		if err := tx.Table("qixi_crm_b_withdrawal_application").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", withdrawalID).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID == 0 {
			return errPayoutNotFound
		}
		if row.PayoutIdempotencyKey != nil && *row.PayoutIdempotencyKey == req.IdempotencyKey {
			if row.PayoutReference == nil || *row.PayoutReference != req.PayoutReference {
				return errPayoutConflict
			}
			out = row
			return nil
		}
		if row.PayoutIdempotencyKey != nil || row.Status != "approved" {
			return errPayoutConflict
		}
		var existing withdraw
		if err := tx.Table("qixi_crm_b_withdrawal_application").Where("user_id = ? AND payout_idempotency_key = ?", row.UserID, req.IdempotencyKey).Scan(&existing).Error; err != nil {
			return err
		}
		if existing.ID != 0 {
			return errPayoutConflict
		}
		now := time.Now()
		result := tx.Table("qixi_crm_b_withdrawal_application").Where("id = ? AND status = 'approved' AND payout_idempotency_key IS NULL", row.ID).Updates(map[string]any{"status": "paid", "payout_idempotency_key": req.IdempotencyKey, "payout_reference": req.PayoutReference, "paid_by": middleware.AdminID(c), "paid_at": now})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errPayoutConflict
		}
		row.Status, row.PayoutIdempotencyKey, row.PayoutReference, row.PaidAt = "paid", &req.IdempotencyKey, &req.PayoutReference, &now
		out = row
		return nil
	})
	switch {
	case err == nil:
		response.OK(c, gin.H{"financial_id": out.ID, "withdrawal_status": out.Status, "financial_status": 1, "payout_reference": out.PayoutReference, "paid_at": out.PaidAt})
	case errors.Is(err, errPayoutNotFound):
		response.Fail(c, http.StatusNotFound, "提现申请不存在")
	case errors.Is(err, errPayoutConflict):
		response.Fail(c, http.StatusConflict, "提现状态已变化或幂等键冲突")
	default:
		fail(c)
	}
}

func (h *Handler) transition(c *gin.Context, target, note string, from []string) {
	id := id(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "提现 ID 错误")
		return
	}
	updates := map[string]any{"status": target, "review_note": note, "reviewed_by": middleware.AdminID(c), "reviewed_at": time.Now()}
	result := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_withdrawal_application").Where("id = ? AND status IN ?", id, from).Updates(updates)
	if result.Error != nil {
		fail(c)
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusConflict, "提现申请状态已变化，请刷新后重试")
		return
	}
	response.OK(c, gin.H{"ok": true, "status": target})
}

func view(row withdraw) gin.H {
	status := 0
	if row.Status == "approved" || row.Status == "paying" || row.Status == "paid" {
		status = 1
	}
	if row.Status == "rejected" {
		status = -1
	}
	channel := map[string]int{"bank": 1, "wechat": 2}[row.Channel]
	return gin.H{"financial_id": row.ID, "financial_sn": row.WithdrawalNo, "mer_id": row.UserID, "user_id": row.UserID, "extract_money": row.Amount, "financial_type": channel, "financial_account": "已脱敏收款账户", "financial_status": map[bool]int{true: 1, false: 0}[row.Status == "paid"], "status": status, "withdrawal_status": row.Status, "refusal": row.ReviewNote, "mark": "用户提现申请", "create_time": row.CreatedAt, "paid_at": row.PaidAt, "payout_reference": row.PayoutReference}
}

func statuses(raw string) ([]string, bool) {
	switch raw {
	case "0":
		return []string{"applied", "reviewing"}, true
	case "1":
		return []string{"approved", "paying", "paid"}, true
	case "-1":
		return []string{"rejected"}, true
	default:
		return nil, false
	}
}
func page(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if p < 1 {
		p = 1
	}
	if l < 1 || l > 100 {
		l = 20
	}
	return p, l
}
func id(c *gin.Context) uint64 { value, _ := strconv.ParseUint(c.Param("id"), 10, 64); return value }
func validPayoutInput(idempotencyKey, payoutReference string) bool {
	return len([]rune(idempotencyKey)) >= 8 && len([]rune(idempotencyKey)) <= 128 && len([]rune(payoutReference)) >= 3 && len([]rune(payoutReference)) <= 128
}
func isPlatform(c *gin.Context) bool {
	claims := middleware.ClaimsFrom(c)
	if claims == nil {
		return false
	}
	for _, role := range claims.Roles {
		if role == "platform" {
			return true
		}
	}
	return false
}
func fail(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "提现查询或审核失败")
}
