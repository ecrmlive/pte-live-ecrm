// Package nativeoperationlog lists store console audit entries.
package nativeoperationlog

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
	merchantDB *gorm.DB
}

func NewHandler(merchantDB *gorm.DB) *Handler { return &Handler{merchantDB: merchantDB} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/setting/operation-logs", h.list)
}

type logRow struct {
	ID           uint64    `gorm:"column:id"`
	AccountID    uint64    `gorm:"column:account_id"`
	Action       string    `gorm:"column:action"`
	ResourceType string    `gorm:"column:resource_type"`
	ResourceID   string    `gorm:"column:resource_id"`
	RequestID    string    `gorm:"column:request_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pagination(c)
	q := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_operation_log").
		Where("store_id = ?", middleware.StoreID(c))
	if action := strings.TrimSpace(c.Query("action")); action != "" {
		q = q.Where("action = ?", action)
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("(action LIKE ? OR resource_type LIKE ? OR resource_id LIKE ? OR request_id LIKE ?)", like, like, like, like)
	}
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where("created_at >= ?", t)
		}
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where("created_at < ?", t.AddDate(0, 0, 1))
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询操作日志失败")
		return
	}
	var rows []logRow
	if err := q.Order("created_at DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询操作日志失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{
			"id": row.ID, "account_id": row.AccountID, "action": row.Action,
			"resource_type": row.ResourceType, "resource_id": row.ResourceID,
			"request_id": row.RequestID, "created_at": row.CreatedAt.Format("2006-01-02 15:04:05"),
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
