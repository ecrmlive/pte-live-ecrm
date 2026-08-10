package feedback

import (
	command "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/feedbackmoderation"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	db, admin *gorm.DB
	commands  *command.Client
}

func New(db, admin *gorm.DB, c *command.Client) *Handler { return &Handler{db, admin, c} }
func (h *Handler) Register(r gin.IRoutes) {
	p := middleware.RequireAdminRoles("platform", "customer_service")
	read := middleware.RequireAdminMenu(h.admin, "user.feedback.read")
	write := middleware.RequireAdminMenu(h.admin, "user.feedback.manage")
	categoryRead := middleware.RequireAdminMenu(h.admin, "user.feedback.category.read")
	categoryManage := middleware.RequireAdminMenu(h.admin, "user.feedback.category.manage")
	opsOK := middleware.RequireAdminRoles("platform", "operations")
	r.GET("/user-feedback", p, read, h.List)
	r.POST("/user-feedback/:id/reply", p, write, h.Reply)
	r.POST("/user-feedback/:id/close", p, write, h.Close)
	r.DELETE("/user-feedback/:id", p, write, h.Delete)
	r.GET("/user-feedback/categories", opsOK, categoryRead, h.ListCategories)
	r.POST("/user-feedback/categories", opsOK, categoryManage, h.CreateCategory)
	r.PUT("/user-feedback/categories/:id", opsOK, categoryManage, h.UpdateCategory)
	r.PUT("/user-feedback/categories/:id/status", opsOK, categoryManage, h.SetCategoryStatus)
	r.DELETE("/user-feedback/categories/:id", opsOK, categoryManage, h.DeleteCategory)
}

type row struct {
	ID        uint64    `gorm:"column:id" json:"id"`
	UserID    uint64    `gorm:"column:user_id" json:"user_id"`
	Type      string    `gorm:"column:type" json:"type"`
	Content   string    `gorm:"column:content" json:"content"`
	Status    string    `gorm:"column:status" json:"status"`
	Reply     string    `gorm:"column:reply" json:"reply"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
type input struct {
	Reply          string `json:"reply"`
	IdempotencyKey string `json:"idempotency_key"`
}
type categoryRow struct {
	ID        uint64    `gorm:"column:id" json:"id"`
	PID       uint64    `gorm:"column:pid" json:"pid"`
	Name      string    `gorm:"column:name" json:"name"`
	Sort      int       `gorm:"column:sort" json:"sort"`
	Status    int       `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
type categoryInput struct {
	Name           string `json:"name"`
	PID            uint64 `json:"pid"`
	Sort           int    `json:"sort"`
	Status         int    `json:"status"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := queryfilter.Page(c)
	status := strings.TrimSpace(c.Query("status"))
	if status != "" {
		if status != "pending" && status != "replied" && status != "closed" {
			response.Fail(c, 400, "反馈状态错误")
			return
		}
	}
	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_feedback").Where("deleted_at IS NULL")
	if status != "" {
		q = q.Where("status=?", status)
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("content LIKE ? OR CAST(user_id AS CHAR) LIKE ?", like, like)
	}
	q = queryfilter.ApplyCreatedAtRange(q, c, "created_at")
	var total int64
	if q.Count(&total).Error != nil {
		response.Fail(c, 500, "反馈查询失败")
		return
	}
	var rows []row
	if q.Order("id DESC").Offset((page-1)*limit).Limit(limit).Scan(&rows).Error != nil {
		response.Fail(c, 500, "反馈查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}
func (h *Handler) Reply(c *gin.Context)  { h.write(c, "reply") }
func (h *Handler) Close(c *gin.Context)  { h.write(c, "close") }
func (h *Handler) Delete(c *gin.Context) { h.write(c, "delete") }

func (h *Handler) ListCategories(c *gin.Context) {
	var rows []categoryRow
	if e := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_feedback_category").Where("deleted_at IS NULL").Order("sort ASC,id ASC").Scan(&rows).Error; e != nil {
		response.Fail(c, 500, "反馈分类查询失败")
		return
	}
	byParent := map[uint64][]gin.H{}
	for _, row := range rows {
		item := gin.H{
			"id":         row.ID,
			"pid":        row.PID,
			"name":       row.Name,
			"sort":       row.Sort,
			"status":     row.Status,
			"created_at": row.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": row.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		byParent[row.PID] = append(byParent[row.PID], item)
	}
	var build func(uint64) []gin.H
	build = func(parent uint64) []gin.H {
		items := byParent[parent]
		for i := range items {
			id := items[i]["id"].(uint64)
			if children := build(id); len(children) > 0 {
				items[i]["children"] = children
			}
		}
		return items
	}
	response.OK(c, gin.H{"list": build(0)})
}
func (h *Handler) CreateCategory(c *gin.Context)    { h.writeCategory(c, "category_create") }
func (h *Handler) UpdateCategory(c *gin.Context)    { h.writeCategory(c, "category_update") }
func (h *Handler) SetCategoryStatus(c *gin.Context) { h.writeCategory(c, "category_status") }
func (h *Handler) DeleteCategory(c *gin.Context)    { h.writeCategory(c, "category_delete") }
func (h *Handler) writeCategory(c *gin.Context, action string) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in categoryInput
	if c.ShouldBindJSON(&in) != nil || (action != "category_create" && id == 0) {
		response.Fail(c, 400, "反馈分类参数错误")
		return
	}
	in.Name, in.IdempotencyKey = strings.TrimSpace(in.Name), strings.TrimSpace(in.IdempotencyKey)
	if action == "category_delete" {
		in.Status, in.Sort, in.PID = 0, 0, 0
	}
	out, e := h.commands.Dispatch(c.Request.Context(), command.Command{
		CategoryID: id, Action: action, Name: in.Name, PID: in.PID, Sort: in.Sort, Status: in.Status,
		OperatorID: uint64(middleware.AdminID(c)), IdempotencyKey: in.IdempotencyKey,
	})
	if e != nil {
		if strings.Contains(e.Error(), "参数错误") {
			response.Fail(c, 400, "反馈分类参数错误")
			return
		}
		response.Fail(c, 503, "反馈分类命令服务不可用")
		return
	}
	if out.Code == "" {
		response.OK(c, gin.H{"category_id": out.CategoryID})
		return
	}
	if out.Code == "not_found" {
		response.Fail(c, 404, "反馈分类不存在")
		return
	}
	if out.Code == "invalid_parent" {
		response.Fail(c, 400, "上级分类不合法，仅支持两级分类")
		return
	}
	if out.Code == "has_child" {
		response.Fail(c, http.StatusConflict, "该分类存在子集，请先处理子集")
		return
	}
	response.Fail(c, http.StatusConflict, "反馈分类已变化或幂等键冲突")
}
func (h *Handler) write(c *gin.Context, action string) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in input
	if id == 0 || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, 400, "反馈参数错误")
		return
	}
	out, e := h.commands.Dispatch(c.Request.Context(), command.Command{FeedbackID: id, Action: action, Reply: strings.TrimSpace(in.Reply), OperatorID: uint64(middleware.AdminID(c)), IdempotencyKey: strings.TrimSpace(in.IdempotencyKey)})
	if e != nil {
		response.Fail(c, 503, "反馈命令服务不可用")
		return
	}
	if out.Code == "" {
		response.OK(c, gin.H{"feedback_id": out.FeedbackID, "status": out.Status})
		return
	}
	if out.Code == "not_found" {
		response.Fail(c, 404, "反馈不存在")
		return
	}
	response.Fail(c, http.StatusConflict, "反馈状态已变化或幂等键冲突")
}
func page(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if p < 1 {
		p = 1
	}
	if l < 1 || l > 100 {
		l = 20
	}
	return p, l
}
