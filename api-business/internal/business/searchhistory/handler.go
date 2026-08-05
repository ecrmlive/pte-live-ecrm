// Package searchhistory implements private C-end product search history.
package searchhistory

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const maxItems = 20

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }
func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/search-history", h.list)
	r.POST("/search-history", h.record)
	r.DELETE("/search-history", h.clear)
	r.DELETE("/search-history/:id", h.remove)
}

type request struct {
	Keyword string `json:"keyword"`
}
type entry struct {
	ID        uint64    `gorm:"column:id"`
	Keyword   string    `gorm:"column:keyword"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (h *Handler) record(c *gin.Context) {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		bad(c, "搜索关键词错误")
		return
	}
	keyword, ok := normalize(req.Keyword)
	if !ok {
		bad(c, "搜索关键词长度应为 1 至 128 个字符")
		return
	}
	uid := uint64(middleware.UID(c))
	now := time.Now()
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_b_search_history").Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "user_id"}, {Name: "keyword"}}, DoUpdates: clause.Assignments(map[string]any{"created_at": now})}).Create(map[string]any{"user_id": uid, "keyword": keyword, "created_at": now}).Error; err != nil {
			return err
		}
		var stale []uint64
		if err := tx.Table("qixi_crm_b_search_history").Select("id").Where("user_id=?", uid).Order("created_at DESC,id DESC").Offset(maxItems).Find(&stale).Error; err != nil {
			return err
		}
		if len(stale) == 0 {
			return nil
		}
		return tx.Table("qixi_crm_b_search_history").Where("user_id=? AND id IN ?", uid, stale).Delete(&entry{}).Error
	})
	if err != nil {
		internal(c)
		return
	}
	response.OK(c, gin.H{"keyword": keyword, "recorded": true})
}

func (h *Handler) list(c *gin.Context) {
	uid := uint64(middleware.UID(c))
	var rows []entry
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_search_history").Where("user_id=?", uid).Order("created_at DESC,id DESC").Limit(maxItems).Find(&rows).Error; err != nil {
		internal(c)
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{"id": row.ID, "keyword": row.Keyword, "created_at": row.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	response.OK(c, gin.H{"list": items})
}

func (h *Handler) remove(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		bad(c, "搜索记录 ID 错误")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_search_history").Where("id=? AND user_id=?", id, middleware.UID(c)).Delete(&entry{}).Error; err != nil {
		internal(c)
		return
	}
	response.OK(c, gin.H{"id": id, "deleted": true})
}

func (h *Handler) clear(c *gin.Context) {
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_search_history").Where("user_id=?", middleware.UID(c)).Delete(&entry{}).Error; err != nil {
		internal(c)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

func normalize(raw string) (string, bool) {
	value := strings.TrimSpace(raw)
	return value, len([]rune(value)) > 0 && len([]rune(value)) <= 128
}
func bad(c *gin.Context, message string) { response.Fail(c, http.StatusBadRequest, message) }
func internal(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "搜索记录服务异常")
}
