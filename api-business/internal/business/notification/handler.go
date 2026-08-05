package notification

import (
	"net/http"
	"strconv"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }
func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/notifications", h.list)
	r.POST("/notifications/:id/read", h.read)
	r.POST("/notifications/read-all", h.readAll)
}

type row struct {
	ID        uint64     `gorm:"column:id"`
	Title     string     `gorm:"column:title"`
	Body      string     `gorm:"column:body"`
	ReadAt    *time.Time `gorm:"column:read_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

func notificationPageParams(c *gin.Context) (int, int) {
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

func (h *Handler) list(c *gin.Context) {
	uid := middleware.UID(c)
	page, limit := notificationPageParams(c)
	base := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_notification").Where("user_id=?", uid)
	var total, unreadTotal int64
	if err := base.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_notification").Where("user_id=? AND read_at IS NULL", uid).Count(&unreadTotal).Error; err != nil {
		fail(c)
		return
	}
	var rows []row
	if err := base.Order("created_at DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		fail(c)
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, v := range rows {
		items = append(items, gin.H{"id": v.ID, "title": v.Title, "body": v.Body, "read": v.ReadAt != nil, "created_at": v.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	response.OK(c, gin.H{"list": items, "total": total, "unread_total": unreadTotal, "page": page, "limit": limit})
}
func (h *Handler) read(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "通知 ID 错误")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_notification").Where("id=? AND user_id=? AND read_at IS NULL", id, middleware.UID(c)).Update("read_at", time.Now()).Error; err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"id": id, "read": true})
}
func (h *Handler) readAll(c *gin.Context) {
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_notification").Where("user_id=? AND read_at IS NULL", middleware.UID(c)).Update("read_at", time.Now()).Error; err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"read": true})
}
func fail(c *gin.Context) { response.Fail(c, http.StatusInternalServerError, "通知服务异常") }
