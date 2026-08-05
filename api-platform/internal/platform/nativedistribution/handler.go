// Package nativedistribution exposes the platform/operations read model for
// business-owned distribution data. It deliberately never adjusts commission
// balances, promoter relations, or withdrawal applications.
package nativedistribution

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	businessDB *gorm.DB
	adminDB    *gorm.DB
}

func NewHandler(businessDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	read := middleware.RequireAdminRoles("platform", "operations")
	readSpread := middleware.RequireAdminMenu(h.adminDB, "marketing.spread.read")
	r.GET("/distribution/promoters", read, readSpread, h.ListPromoters)
	r.GET("/distribution/commissions", read, readSpread, h.ListCommissions)
	r.GET("/distribution/summary", read, readSpread, h.Summary)
}

type promoterRow struct {
	UserID              uint64    `gorm:"column:user_id" json:"user_id"`
	Status              int8      `gorm:"column:status" json:"status"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"updated_at"`
	DirectUserCount     int64     `gorm:"column:direct_user_count" json:"direct_user_count"`
	PendingCommission   float64   `gorm:"column:pending_commission" json:"pending_commission"`
	AvailableCommission float64   `gorm:"column:available_commission" json:"available_commission"`
	SettledCommission   float64   `gorm:"column:settled_commission" json:"settled_commission"`
}

type commissionRow struct {
	ID          uint64     `gorm:"column:id" json:"commission_id"`
	UserID      uint64     `gorm:"column:user_id" json:"user_id"`
	OrderID     uint64     `gorm:"column:order_id" json:"order_id"`
	Amount      float64    `gorm:"column:amount" json:"amount"`
	Status      string     `gorm:"column:status" json:"status"`
	AvailableAt *time.Time `gorm:"column:available_at" json:"available_at"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
}

type distributionSummary struct {
	PromoterCount       int64   `gorm:"column:promoter_count" json:"promoter_count"`
	ActivePromoterCount int64   `gorm:"column:active_promoter_count" json:"active_promoter_count"`
	PendingCommission   float64 `gorm:"column:pending_commission" json:"pending_commission"`
	AvailableCommission float64 `gorm:"column:available_commission" json:"available_commission"`
	SettledCommission   float64 `gorm:"column:settled_commission" json:"settled_commission"`
}

func (h *Handler) ListPromoters(c *gin.Context) {
	page, limit := paging(c)
	q := h.promoterQuery(c)
	if status, ok := promoterStatus(c.Query("status")); !ok {
		response.Fail(c, http.StatusBadRequest, "推广员状态错误")
		return
	} else if status != nil {
		q = q.Where("p.status = ?", *status)
	}
	if userID, provided, ok := queryID(c, "user_id"); !ok {
		response.Fail(c, http.StatusBadRequest, "用户 ID 错误")
		return
	} else if provided {
		q = q.Where("p.user_id = ?", userID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		failure(c)
		return
	}
	rows := make([]promoterRow, 0)
	if err := q.Select("p.user_id,p.status,p.updated_at,COALESCE(rel.direct_user_count,0) AS direct_user_count,COALESCE(ledger.pending_commission,0) AS pending_commission,COALESCE(ledger.available_commission,0) AS available_commission,COALESCE(ledger.settled_commission,0) AS settled_commission").Order("p.updated_at DESC,p.user_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		failure(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) ListCommissions(c *gin.Context) {
	page, limit := paging(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_commission_ledger")
	if status, ok := commissionStatus(c.Query("status")); !ok {
		response.Fail(c, http.StatusBadRequest, "佣金状态错误")
		return
	} else if status != "" {
		q = q.Where("status = ?", status)
	}
	if userID, provided, ok := queryID(c, "user_id"); !ok {
		response.Fail(c, http.StatusBadRequest, "用户 ID 错误")
		return
	} else if provided {
		q = q.Where("user_id = ?", userID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		failure(c)
		return
	}
	rows := make([]commissionRow, 0)
	if err := q.Select("id,user_id,order_id,amount,status,available_at,created_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		failure(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Summary(c *gin.Context) {
	var out distributionSummary
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_promoter AS p").
		Select("COUNT(*) AS promoter_count,COALESCE(SUM(CASE WHEN p.status = 1 THEN 1 ELSE 0 END),0) AS active_promoter_count,COALESCE((SELECT SUM(amount) FROM qixi_crm_b_commission_ledger WHERE status = 'pending'),0) AS pending_commission,COALESCE((SELECT SUM(amount) FROM qixi_crm_b_commission_ledger WHERE status = 'available'),0) AS available_commission,COALESCE((SELECT SUM(amount) FROM qixi_crm_b_commission_ledger WHERE status = 'settled'),0) AS settled_commission").Scan(&out).Error
	if err != nil {
		failure(c)
		return
	}
	response.OK(c, out)
}

func (h *Handler) promoterQuery(c *gin.Context) *gorm.DB {
	return h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_promoter AS p").
		Joins("LEFT JOIN (SELECT user_id,COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END),0) AS pending_commission,COALESCE(SUM(CASE WHEN status = 'available' THEN amount ELSE 0 END),0) AS available_commission,COALESCE(SUM(CASE WHEN status = 'settled' THEN amount ELSE 0 END),0) AS settled_commission FROM qixi_crm_b_commission_ledger GROUP BY user_id) AS ledger ON ledger.user_id = p.user_id").
		Joins("LEFT JOIN (SELECT parent_user_id,COUNT(*) AS direct_user_count FROM qixi_crm_b_distribution_relation WHERE parent_user_id IS NOT NULL GROUP BY parent_user_id) AS rel ON rel.parent_user_id = p.user_id")
}

func promoterStatus(raw string) (*int8, bool) {
	switch strings.TrimSpace(raw) {
	case "":
		return nil, true
	case "0":
		value := int8(0)
		return &value, true
	case "1":
		value := int8(1)
		return &value, true
	default:
		return nil, false
	}
}

func commissionStatus(raw string) (string, bool) {
	switch strings.TrimSpace(raw) {
	case "", "pending", "available", "settled", "voided":
		return strings.TrimSpace(raw), true
	default:
		return "", false
	}
}

func queryID(c *gin.Context, field string) (uint64, bool, bool) {
	raw := strings.TrimSpace(c.Query(field))
	if raw == "" {
		return 0, false, true
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	return value, true, err == nil && value > 0
}

func paging(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}

func failure(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "分销监管数据查询失败")
}
