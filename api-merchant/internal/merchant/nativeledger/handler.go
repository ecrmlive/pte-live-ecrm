// Package nativeledger exposes store finance ledger read models.
package nativeledger

import (
	"net/http"
	"strconv"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	merchantDB *gorm.DB
}

func NewHandler(merchantDB *gorm.DB) *Handler { return &Handler{merchantDB: merchantDB} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/finance/ledger", h.listLedger)
	r.GET("/finance/statements", h.listStatements)
}

type ledgerRow struct {
	ID            uint64    `gorm:"column:id"`
	EntryType     string    `gorm:"column:entry_type"`
	Amount        float64   `gorm:"column:amount"`
	ReferenceType string    `gorm:"column:reference_type"`
	ReferenceID   string    `gorm:"column:reference_id"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (h *Handler) listLedger(c *gin.Context) {
	page, limit := pagination(c)
	q := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_finance_ledger").
		Where("store_id = ?", middleware.StoreID(c))
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询资金流水失败")
		return
	}
	var rows []ledgerRow
	if err := q.Order("created_at DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询资金流水失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{
			"id": row.ID, "entry_type": row.EntryType, "amount": row.Amount,
			"reference_type": row.ReferenceType, "reference_id": row.ReferenceID,
			"created_at": row.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

type statementRow struct {
	ID          uint64    `gorm:"column:id"`
	PeriodStart time.Time `gorm:"column:period_start"`
	PeriodEnd   time.Time `gorm:"column:period_end"`
	Amount      float64   `gorm:"column:amount"`
	Status      string    `gorm:"column:status"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (h *Handler) listStatements(c *gin.Context) {
	page, limit := pagination(c)
	q := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_settlement_bill").
		Where("store_id = ?", middleware.StoreID(c))
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询对账单失败")
		return
	}
	var rows []statementRow
	if err := q.Order("period_end DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询对账单失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{
			"statement_id": row.ID, "period_start": row.PeriodStart.Format("2006-01-02"),
			"period_end": row.PeriodEnd.Format("2006-01-02"), "amount": row.Amount,
			"status": row.Status, "updated_at": row.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func pagination(c *gin.Context) (int, int) {
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
