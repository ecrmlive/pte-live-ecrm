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
	ID             uint64    `gorm:"column:id" json:"id"`
	AdminUserID    uint64    `gorm:"column:admin_user_id" json:"admin_user_id"`
	AdminName      string    `gorm:"column:admin_name" json:"admin_name"`
	RoleCode       string    `gorm:"column:role_code" json:"-"`
	Action         string    `gorm:"column:action" json:"-"`
	ResourceType   string    `gorm:"column:resource_type" json:"-"`
	RequestMethod  string    `gorm:"column:request_method" json:"request_method"`
	RequestPath    string    `gorm:"column:request_path" json:"link"`
	RequestIP      string    `gorm:"column:request_ip" json:"ip"`
	PermissionName string    `gorm:"column:permission_name" json:"permission_name"`
	Request        string    `gorm:"-" json:"request"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := pagination(c)
	adminKeyword, requestMethod := strings.TrimSpace(c.Query("admin_keyword")), strings.ToUpper(strings.TrimSpace(c.Query("request_method")))
	if len(adminKeyword) > 128 || len(requestMethod) > 16 {
		response.Fail(c, http.StatusBadRequest, "操作日志筛选参数过长")
		return
	}
	if requestMethod != "" && !isRequestMethod(requestMethod) {
		response.Fail(c, http.StatusBadRequest, "请求方式参数错误")
		return
	}
	start, end, ok := dateRange(c.Query("start_date"), c.Query("end_date"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "操作日期参数错误")
		return
	}
	db := h.admin.WithContext(c.Request.Context()).
		Table("qixi_crm_a_operation_log AS l").
		Joins("LEFT JOIN qixi_crm_a_admin_user AS u ON u.id = l.admin_user_id")
	if adminKeyword != "" {
		like := "%" + adminKeyword + "%"
		db = db.Where("CAST(l.admin_user_id AS CHAR) LIKE ? OR u.username LIKE ? OR u.display_name LIKE ?", like, like, like)
	}
	if requestMethod != "" {
		db = db.Where("COALESCE(NULLIF(l.request_method, ''), SUBSTRING_INDEX(l.action, ' ', 1)) = ?", requestMethod)
	}
	if !start.IsZero() {
		db = db.Where("l.created_at>=?", start)
	}
	if !end.IsZero() {
		db = db.Where("l.created_at<?", end)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "操作日志查询失败")
		return
	}
	rows := make([]row, 0)
	if err := db.Select(`l.id,l.admin_user_id,l.role_code,l.action,l.resource_type,
      l.request_method,l.request_path,l.request_ip,l.permission_name,l.created_at,
      COALESCE(NULLIF(u.display_name,''), NULLIF(u.username,''), '管理员已删除') AS admin_name`).
		Order("l.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "操作日志查询失败")
		return
	}
	for i := range rows {
		rows[i].normalize()
	}
	response.OK(c, gin.H{"list": rows, "total": total})
}

func (r *row) normalize() {
	method, path := splitAction(r.Action)
	if r.RequestMethod == "" {
		r.RequestMethod = method
	}
	if r.RequestPath == "" {
		r.RequestPath = path
	}
	if r.RequestIP == "" {
		r.RequestIP = "—"
	}
	if r.PermissionName == "" {
		r.PermissionName = operationPermissionName(r.ResourceType)
	}
	r.Request = r.Action
	if r.Request == "" {
		r.Request = "—"
	}
}

func splitAction(action string) (string, string) {
	parts := strings.Fields(strings.TrimSpace(action))
	if len(parts) >= 2 && isRequestMethod(strings.ToUpper(parts[0])) {
		return strings.ToUpper(parts[0]), strings.Join(parts[1:], " ")
	}
	return "", ""
}

func isRequestMethod(value string) bool {
	switch value {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func operationPermissionName(resource string) string {
	switch resource {
	case "setting", "maintain":
		return "系统设置"
	case "customer-service":
		return "客服管理"
	case "product", "products":
		return "商品管理"
	case "merchants":
		return "商户管理"
	case "stores":
		return "店铺管理"
	case "users", "user-list":
		return "用户管理"
	default:
		return "平台管理"
	}
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
