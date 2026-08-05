package operationlog

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ admin *gorm.DB }

func New(admin *gorm.DB) *Handler { return &Handler{admin: admin} }

func (h *Handler) Register(r gin.IRoutes) {
	read := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "setting.operation_log.read")}
	r.GET("/operation-logs", append(read, h.List)...)
	r.GET("/login-logs", append(read, h.ListLoginLogs)...)
}

type row struct {
	ID           uint64    `gorm:"column:id" json:"id"`
	AdminUserID  uint64    `gorm:"column:admin_user_id" json:"admin_user_id"`
	RoleCode     string    `gorm:"column:role_code" json:"role_code"`
	Action       string    `gorm:"column:action" json:"action"`
	ResourceType string    `gorm:"column:resource_type" json:"resource_type"`
	ResourceID   string    `gorm:"column:resource_id" json:"resource_id"`
	RequestID    string    `gorm:"column:request_id" json:"request_id"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := pagination(c)
	adminID, ok := optionalID(c, "admin_user_id")
	if !ok {
		response.Fail(c, http.StatusBadRequest, "管理员 ID 参数错误")
		return
	}
	roleCode, action, resourceType := strings.TrimSpace(c.Query("role_code")), strings.TrimSpace(c.Query("action")), strings.TrimSpace(c.Query("resource_type"))
	if len(roleCode) > 32 || len(action) > 128 || len(resourceType) > 64 {
		response.Fail(c, http.StatusBadRequest, "操作日志筛选参数过长")
		return
	}
	start, end, ok := dateRange(c.Query("start_date"), c.Query("end_date"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "操作日期参数错误")
		return
	}
	db := h.admin.WithContext(c.Request.Context()).Table("qixi_crm_a_operation_log")
	if adminID != 0 {
		db = db.Where("admin_user_id=?", adminID)
	}
	if roleCode != "" {
		db = db.Where("role_code=?", roleCode)
	}
	if action != "" {
		db = db.Where("action LIKE ?", "%"+action+"%")
	}
	if resourceType != "" {
		db = db.Where("resource_type=?", resourceType)
	}
	if !start.IsZero() {
		db = db.Where("created_at>=?", start)
	}
	if !end.IsZero() {
		db = db.Where("created_at<?", end)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "操作日志查询失败")
		return
	}
	rows := make([]row, 0)
	if err := db.Select("id,admin_user_id,role_code,action,resource_type,resource_id,request_id,created_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "操作日志查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total})
}

type loginRow struct {
	ID          uint64    `gorm:"column:id" json:"id"`
	AdminUserID *uint64   `gorm:"column:admin_user_id" json:"admin_user_id"`
	Username    string    `gorm:"column:username" json:"username"`
	RoleCode    string    `gorm:"column:role_code" json:"role_code"`
	Success     bool      `gorm:"column:success" json:"success"`
	IP          string    `gorm:"column:ip" json:"ip"`
	UserAgent   string    `gorm:"column:user_agent" json:"user_agent"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (h *Handler) ListLoginLogs(c *gin.Context) {
	page, limit := pagination(c)
	username := strings.TrimSpace(c.Query("username"))
	if len(username) > 64 {
		response.Fail(c, http.StatusBadRequest, "账号参数过长")
		return
	}
	successRaw := strings.TrimSpace(c.Query("success"))
	if successRaw != "" && successRaw != "0" && successRaw != "1" {
		response.Fail(c, http.StatusBadRequest, "登录结果参数错误")
		return
	}
	start, end, ok := dateRange(c.Query("start_date"), c.Query("end_date"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "登录日期参数错误")
		return
	}
	db := h.admin.WithContext(c.Request.Context()).Table("qixi_crm_a_login_log")
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if successRaw != "" {
		db = db.Where("success=?", successRaw == "1")
	}
	if !start.IsZero() {
		db = db.Where("created_at>=?", start)
	}
	if !end.IsZero() {
		db = db.Where("created_at<?", end)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "登录日志查询失败")
		return
	}
	rows := make([]loginRow, 0)
	if err := db.Select("id,admin_user_id,username,role_code,success,ip,user_agent,created_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "登录日志查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total})
}

func pagination(c *gin.Context) (int, int) {
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

func optionalID(c *gin.Context, key string) (uint64, bool) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0, true
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	return id, err == nil && id > 0
}

func dateRange(startRaw, endRaw string) (time.Time, time.Time, bool) {
	startRaw, endRaw = strings.TrimSpace(startRaw), strings.TrimSpace(endRaw)
	var start, end time.Time
	var err error
	if startRaw != "" {
		start, err = time.ParseInLocation("2006-01-02", startRaw, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, false
		}
	}
	if endRaw != "" {
		end, err = time.ParseInLocation("2006-01-02", endRaw, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, false
		}
		end = end.AddDate(0, 0, 1)
	}
	if !start.IsZero() && !end.IsZero() && !start.Before(end) {
		return time.Time{}, time.Time{}, false
	}
	return start, end, true
}
