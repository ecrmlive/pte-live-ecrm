// Package nativestorecustomer lists buyers who ordered from the current store.
package nativestorecustomer

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	businessDB *gorm.DB
}

func NewHandler(businessDB *gorm.DB) *Handler { return &Handler{businessDB: businessDB} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/store-customers", h.list)
}

type customerRow struct {
	UserID       uint64    `gorm:"column:user_id"`
	Nickname     string    `gorm:"column:nickname"`
	Mobile       string    `gorm:"column:mobile"`
	Status       int8      `gorm:"column:status"`
	OrderCount   int64     `gorm:"column:order_count"`
	LastOrderAt  time.Time `gorm:"column:last_order_at"`
	TotalPay     float64   `gorm:"column:total_pay"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pagination(c)
	storeID := middleware.StoreID(c)
	sub := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order AS o").
		Select("o.user_id, COUNT(*) AS order_count, MAX(o.created_at) AS last_order_at, COALESCE(SUM(o.pay_amount), 0) AS total_pay").
		Where("o.store_id = ? AND o.status <> 'pending_pay' AND o.status <> 'cancelled'", storeID).
		Group("o.user_id")
	q := h.businessDB.WithContext(c.Request.Context()).Table("(?) AS stats", sub).
		Select("stats.user_id, u.nickname, u.mobile, u.status, stats.order_count, stats.last_order_at, stats.total_pay").
		Joins("JOIN qixi_crm_b_user AS u ON u.id = stats.user_id")
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		if id, err := strconv.ParseUint(keyword, 10, 64); err == nil && id > 0 {
			q = q.Where("stats.user_id = ?", id)
		} else {
			q = q.Where("u.nickname LIKE ?", "%"+keyword+"%")
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询店铺客户失败")
		return
	}
	var rows []customerRow
	if err := q.Order("stats.last_order_at DESC, stats.user_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询店铺客户失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{
			"user_id": row.UserID, "nickname": row.Nickname, "mobile": maskMobile(row.Mobile),
			"status": row.Status, "order_count": row.OrderCount, "total_pay": row.TotalPay,
			"last_order_at": row.LastOrderAt.Format("2006-01-02 15:04:05"),
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

func maskMobile(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) < 7 {
		return raw
	}
	return raw[:3] + "****" + raw[len(raw)-4:]
}
