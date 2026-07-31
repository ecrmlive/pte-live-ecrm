// Package account exposes the C-end asset read model. It only reads the
// business-owned member account and immutable asset ledger.
package account

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler     { return &Handler{db: db} }
func (h *Handler) Register(r gin.IRoutes) { r.GET("/account/balance", h.Balance) }

type balanceSummary struct {
	Balance      float64 `json:"balance"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}
type ledgerRow struct {
	ID            uint64  `json:"id"`
	Amount        float64 `json:"amount"`
	ReferenceType string  `json:"reference_type"`
	ReferenceID   string  `json:"reference_id"`
	CreatedAt     string  `json:"created_at"`
}

func (h *Handler) Balance(c *gin.Context) {
	uid := middleware.UID(c)
	summary := balanceSummary{}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_member_account").Select("COALESCE(balance, 0) AS balance").Where("user_id = ?", uid).Scan(&summary).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询余额失败")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_asset_ledger").Select("COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END),0) AS total_income, COALESCE(SUM(CASE WHEN amount < 0 THEN -amount ELSE 0 END),0) AS total_expense").Where("user_id = ? AND asset_type = ?", uid, "balance").Scan(&summary).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询余额流水失败")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	const limit = 20
	var rows []ledgerRow
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_asset_ledger").Select("id, amount, reference_type, reference_id, created_at").Where("user_id = ? AND asset_type = ?", uid, "balance").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询余额流水失败")
		return
	}
	response.OK(c, gin.H{"summary": summary, "list": rows, "page": page, "limit": limit})
}
