// Package nativeledger exposes platform-supervised user asset ledgers from the
// business schema. It never reads the retired qixi_m_admin_financial_record
// tables and intentionally returns no user identity or payment credentials.
package nativeledger

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
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
	platformOnly := middleware.RequireAdminRoles("platform")
	readLedger := middleware.RequireAdminMenu(h.adminDB, "accounts.user_assets.read")
	r.GET("/finance/user-assets", platformOnly, readLedger, h.List)
	r.GET("/finance/user-assets/summary", platformOnly, readLedger, h.Summary)
}

type ledgerRow struct {
	ID            uint64    `gorm:"column:id" json:"id"`
	UserID        uint64    `gorm:"column:user_id" json:"user_id"`
	AssetType     string    `gorm:"column:asset_type" json:"asset_type"`
	Amount        float64   `gorm:"column:amount" json:"amount"`
	ReferenceType string    `gorm:"column:reference_type" json:"reference_type"`
	ReferenceID   string    `gorm:"column:reference_id" json:"reference_id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

type ledgerSummary struct {
	AssetType string  `gorm:"column:asset_type" json:"asset_type"`
	Income    float64 `gorm:"column:income" json:"income"`
	Expense   float64 `gorm:"column:expense" json:"expense"`
	Count     int64   `gorm:"column:count" json:"count"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := pageParams(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_asset_ledger")
	if assetType, ok := assetType(c.Query("asset_type")); !ok {
		response.Fail(c, http.StatusBadRequest, "资产类型错误")
		return
	} else if assetType != "" {
		q = q.Where("asset_type = ?", assetType)
	}
	if rawUserID := strings.TrimSpace(c.Query("user_id")); rawUserID != "" {
		userID, err := strconv.ParseUint(rawUserID, 10, 64)
		if err != nil || userID == 0 {
			response.Fail(c, http.StatusBadRequest, "用户 ID 错误")
			return
		}
		q = q.Where("user_id = ?", userID)
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("reference_type LIKE ? OR reference_id LIKE ?", like, like)
	}
	q = queryfilter.ApplyCreatedAtRange(q, c, "created_at")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		writeFailure(c)
		return
	}
	rows := make([]ledgerRow, 0)
	if err := q.Select("id,user_id,asset_type,amount,reference_type,reference_id,created_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		writeFailure(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Summary(c *gin.Context) {
	rows := make([]ledgerSummary, 0)
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_asset_ledger").
		Select("asset_type, COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS income, COALESCE(SUM(CASE WHEN amount < 0 THEN -amount ELSE 0 END), 0) AS expense, COUNT(*) AS count").
		Group("asset_type").Order("asset_type ASC").Scan(&rows).Error
	if err != nil {
		writeFailure(c)
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func assetType(raw string) (string, bool) {
	switch strings.TrimSpace(raw) {
	case "", "balance", "points", "commission":
		return strings.TrimSpace(raw), true
	default:
		return "", false
	}
}

func pageParams(c *gin.Context) (int, int) {
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

func writeFailure(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "用户资产流水查询失败")
}
