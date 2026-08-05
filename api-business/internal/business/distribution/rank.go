package distribution

import (
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	rankMetricCommission = "commission"
	rankMetricPromoters  = "promoters"
)

type rankRow struct {
	UserID           uint64  `gorm:"column:user_id"`
	Nickname         string  `gorm:"column:nickname"`
	CommissionAmount float64 `gorm:"column:commission_amount"`
	ReferralCount    int64   `gorm:"column:referral_count"`
}

type directUserRow struct {
	UserID   uint64    `gorm:"column:user_id"`
	Nickname string    `gorm:"column:nickname"`
	BoundAt  time.Time `gorm:"column:bound_at"`
}

func normalizeRankMetric(value string) (string, bool) {
	switch strings.TrimSpace(value) {
	case "", rankMetricCommission:
		return rankMetricCommission, true
	case rankMetricPromoters:
		return rankMetricPromoters, true
	default:
		return "", false
	}
}

func maskNickname(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "推广用户"
	}
	first, _ := utf8.DecodeRuneInString(value)
	return string(first) + "**"
}

func (h *Handler) rank(c *gin.Context) {
	metric, ok := normalizeRankMetric(c.Query("metric"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "排行指标错误")
		return
	}
	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_promoter AS p").
		Joins("JOIN qixi_crm_b_user AS u ON u.id = p.user_id AND u.status = 1").
		Where("p.status = 1")
	rows := make([]rankRow, 0)
	if metric == rankMetricCommission {
		q = q.Select("p.user_id, u.nickname, COALESCE(SUM(CASE WHEN l.status <> 'voided' THEN l.amount ELSE 0 END), 0) AS commission_amount, 0 AS referral_count").
			Joins("LEFT JOIN qixi_crm_b_commission_ledger AS l ON l.user_id = p.user_id").
			Group("p.user_id, u.nickname").
			Order("commission_amount DESC, p.user_id ASC")
	} else {
		q = q.Select("p.user_id, u.nickname, 0 AS commission_amount, COUNT(r.user_id) AS referral_count").
			Joins("LEFT JOIN qixi_crm_b_distribution_relation AS r ON r.parent_user_id = p.user_id").
			Group("p.user_id, u.nickname").
			Order("referral_count DESC, p.user_id ASC")
	}
	if err := q.Limit(100).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for index, row := range rows {
		list = append(list, gin.H{
			"rank":              index + 1,
			"nickname":          maskNickname(row.Nickname),
			"commission_amount": row.CommissionAmount,
			"referral_count":    row.ReferralCount,
		})
	}
	response.OK(c, gin.H{"metric": metric, "list": list})
}

func (h *Handler) users(c *gin.Context) {
	page, limit := pageParams(c)
	uid := uint64(middleware.UID(c))
	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_relation AS r").
		Joins("JOIN qixi_crm_b_distribution_promoter AS p ON p.user_id = r.parent_user_id AND p.status = 1").
		Joins("JOIN qixi_crm_b_user AS u ON u.id = r.user_id AND u.status = 1").
		Where("r.parent_user_id = ?", uid)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	rows := make([]directUserRow, 0)
	if err := q.Select("r.user_id, u.nickname, r.bound_at").Order("r.bound_at DESC, r.user_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, gin.H{"nickname": maskNickname(row.Nickname), "bound_at": row.BoundAt.Format("2006-01-02 15:04:05")})
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}
