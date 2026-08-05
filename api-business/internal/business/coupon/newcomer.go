package coupon

import (
	"errors"
	"net/http"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// newcomer mirrors CRMEB's new_people read flow. Benefits are created during
// account registration; this endpoint is only an authenticated presentation of
// the active policy and the coupons already credited to the current user.
func (h *Handler) newcomer(c *gin.Context) {
	var policy struct {
		Enabled       int8   `gorm:"column:enabled"`
		CouponEnabled int8   `gorm:"column:coupon_enabled"`
		Title         string `gorm:"column:title"`
		Description   string `gorm:"column:description"`
	}
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_onboarding_policy").Where("id = ?", 1).First(&policy).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || policy.Enabled == 0 || policy.CouponEnabled == 0 {
		response.OK(c, gin.H{"enabled": false, "title": "新人礼", "description": "当前暂无新人权益", "coupon": []gin.H{}})
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "新人权益服务异常")
		return
	}
	now := time.Now()
	var rows []templateRow
	query := h.baseTemplates(c, uint64(middleware.UID(c))).
		Joins("JOIN qixi_crm_b_onboarding_coupon AS oc ON oc.coupon_id = c.coupon_id").
		Where("oc.enabled = 1 AND c.status = 1").
		Where("(c.starts_at IS NULL OR c.starts_at <= ?)", now).
		Where("(c.ends_at IS NULL OR c.ends_at >= ?)", now)
	if err := query.Order("oc.sort ASC, c.coupon_id ASC").Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "新人权益服务异常")
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, couponView(row))
	}
	response.OK(c, gin.H{"enabled": true, "title": policy.Title, "description": policy.Description, "coupon": list})
}
