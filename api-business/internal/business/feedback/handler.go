// Package feedback owns C-end user feedback submissions and private progress queries.
package feedback

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/feedback", h.create)
	r.GET("/feedback/categories", h.categories)
	r.GET("/feedback/list", h.list)
	r.GET("/feedback/detail/:id", h.detail)
}

type createRequest struct {
	CategoryID uint64 `json:"category_id"`
	Type       string `json:"type"`
	Content    string `json:"content"`
}

type row struct {
	ID         uint64    `gorm:"column:id"`
	UserID     uint64    `gorm:"column:user_id"`
	CategoryID uint64    `gorm:"column:category_id"`
	Type       string    `gorm:"column:type"`
	Content    string    `gorm:"column:content"`
	Status     string    `gorm:"column:status"`
	Reply      string    `gorm:"column:reply"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}
type categoryRow struct {
	ID   uint64 `gorm:"column:id"`
	Name string `gorm:"column:name"`
	Sort int    `gorm:"column:sort"`
}

func (h *Handler) create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		bad(c, "反馈参数错误")
		return
	}
	req.Type = strings.TrimSpace(req.Type)
	req.Content = strings.TrimSpace(req.Content)
	if req.Type == "" {
		req.Type = "功能建议"
	}
	if len([]rune(req.Type)) > 32 || req.Content == "" || len([]rune(req.Content)) > 1000 {
		bad(c, "反馈内容不合法")
		return
	}
	if req.CategoryID > 0 {
		var category categoryRow
		err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_feedback_category").Where("id=? AND status=1 AND deleted_at IS NULL", req.CategoryID).Take(&category).Error
		if err == gorm.ErrRecordNotFound {
			bad(c, "反馈分类不可用")
			return
		}
		if err != nil {
			internal(c)
			return
		}
		req.Type = category.Name
	}
	created := row{UserID: uint64(middleware.UID(c)), CategoryID: req.CategoryID, Type: req.Type, Content: req.Content, Status: "pending"}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_feedback").Create(&created).Error; err != nil {
		internal(c)
		return
	}
	response.OK(c, created.view())
}

func (h *Handler) categories(c *gin.Context) {
	var rows []categoryRow
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_feedback_category").Where("status=1 AND deleted_at IS NULL").Order("sort ASC,id ASC").Find(&rows).Error; err != nil {
		internal(c)
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pagination(c)
	query := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_feedback").Where("user_id = ? AND deleted_at IS NULL", middleware.UID(c))
	var total int64
	if err := query.Count(&total).Error; err != nil {
		internal(c)
		return
	}
	rows := make([]row, 0)
	if err := query.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		internal(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, item := range rows {
		list = append(list, item.view())
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) detail(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		bad(c, "反馈 ID 错误")
		return
	}
	var item row
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_feedback").Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, middleware.UID(c)).Take(&item).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, http.StatusNotFound, "反馈不存在")
		return
	}
	if err != nil {
		internal(c)
		return
	}
	response.OK(c, item.view())
}

func (r row) view() gin.H {
	return gin.H{"feedback_id": r.ID, "category_id": r.CategoryID, "type": r.Type, "content": r.Content, "status": r.Status, "reply": r.Reply, "create_time": r.CreatedAt.Format("2006-01-02 15:04:05"), "update_time": r.UpdatedAt.Format("2006-01-02 15:04:05")}
}

func positiveID(raw string) (uint64, bool) {
	id, err := strconv.ParseUint(raw, 10, 64)
	return id, err == nil && id > 0
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
	if limit > 50 {
		limit = 50
	}
	return page, limit
}
func bad(c *gin.Context, message string) { response.Fail(c, http.StatusBadRequest, message) }
func internal(c *gin.Context)            { response.Fail(c, http.StatusInternalServerError, "反馈服务异常") }
