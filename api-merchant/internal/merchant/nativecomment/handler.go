// Package nativecomment exposes store-scoped read of product comments.
package nativecomment

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
	r.GET("/product/comments", h.list)
	r.POST("/product/comments/:id/reply", h.reply)
}

type row struct {
	ID           uint64    `gorm:"column:id"`
	ProductID    uint64    `gorm:"column:product_id"`
	UserID       uint64    `gorm:"column:user_id"`
	Score        int       `gorm:"column:score"`
	Content      string    `gorm:"column:content"`
	ReplyContent string    `gorm:"column:reply_content"`
	Source       string    `gorm:"column:source"`
	Status       string    `gorm:"column:status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	ProductTitle string    `gorm:"column:product_title"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pagination(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_comment AS pc").
		Select("pc.id,pc.product_id,pc.user_id,pc.score,pc.content,pc.reply_content,pc.source,pc.status,pc.created_at,COALESCE(p.title,'') AS product_title").
		Joins("LEFT JOIN qixi_crm_b_product_view AS p ON p.product_id=pc.product_id").
		Where("pc.store_id = ? AND pc.deleted_at IS NULL", middleware.StoreID(c))
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		q = q.Where("pc.status = ?", status)
	}
	if productID, err := strconv.ParseUint(c.Query("product_id"), 10, 64); err == nil && productID > 0 {
		q = q.Where("pc.product_id = ?", productID)
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("(pc.content LIKE ? OR p.title LIKE ?)", like, like)
	}
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where("pc.created_at >= ?", t)
		}
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where("pc.created_at < ?", t.AddDate(0, 0, 1))
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品评论失败")
		return
	}
	var rows []row
	if err := q.Order("pc.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品评论失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, x := range rows {
		items = append(items, gin.H{
			"id": x.ID, "product_id": x.ProductID, "product_title": x.ProductTitle, "user_id": x.UserID,
			"score": x.Score, "content": x.Content, "reply_content": x.ReplyContent, "source": x.Source,
			"status": x.Status, "created_at": x.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func (h *Handler) reply(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		ReplyContent string `json:"reply_content"`
	}
	if id == 0 || c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.ReplyContent) == "" {
		response.Fail(c, http.StatusBadRequest, "回复内容不能为空")
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_comment").
		Where("id = ? AND store_id = ? AND deleted_at IS NULL", id, middleware.StoreID(c)).
		Updates(map[string]any{"reply_content": strings.TrimSpace(req.ReplyContent), "replied_at": time.Now()})
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "保存商家回复失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "评论不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
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
