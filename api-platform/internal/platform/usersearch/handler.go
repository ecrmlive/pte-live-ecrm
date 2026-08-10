// Package usersearch exposes the platform-only supervision plane for C-end
// search history. List joins nickname/avatar/user type for CRMEB parity;
// clear/export stay audited and never return mobile or device identifiers.
package usersearch

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ business, admin *gorm.DB }

func New(business, admin *gorm.DB) *Handler { return &Handler{business: business, admin: admin} }

func (h *Handler) Register(r gin.IRoutes) {
	read := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.search_record.read")}
	clear := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.search_record.clear")}
	export := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.search_record.export")}
	r.GET("/user-search-records", append(read, h.List)...)
	r.POST("/user-search-records/clear", append(clear, h.Clear)...)
	r.POST("/user-search-records/export", append(export, h.Export)...)
}

type row struct {
	ID        uint64 `gorm:"column:id" json:"id"`
	UserID    uint64 `gorm:"column:user_id" json:"user_id"`
	Nickname  string `gorm:"column:nickname" json:"nickname"`
	AvatarURL string `gorm:"column:avatar_url" json:"avatar_url"`
	UserType  string `gorm:"column:user_type" json:"user_type"`
	Keyword   string `gorm:"column:keyword" json:"keyword"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}

type query struct {
	Keyword   string `form:"keyword" json:"keyword"`
	Nickname  string `form:"nickname" json:"nickname"`
	UserType  string `form:"user_type" json:"user_type"`
	StartDate string `form:"start_date" json:"start_date"`
	EndDate   string `form:"end_date" json:"end_date"`
}

type clearInput struct {
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

type exportInput struct {
	query
	Reason string `json:"reason"`
}

type clearAudit struct {
	UserID         uint64 `gorm:"column:user_id"`
	Reason         string `gorm:"column:reason"`
	IdempotencyKey string `gorm:"column:idempotency_key"`
	OperatorAdmin  uint64 `gorm:"column:operator_admin_id"`
	ClearedCount   int64  `gorm:"column:cleared_count"`
}

var errClearConflict = errors.New("search record clear idempotency conflict")

func (h *Handler) List(c *gin.Context) {
	filter, ok := parseQuery(c, query{
		Keyword:   c.Query("keyword"),
		Nickname:  c.Query("nickname"),
		UserType:  c.Query("user_type"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	})
	if !ok {
		return
	}
	page, limit := pagination(c)
	q := applyQuery(h.business.WithContext(c.Request.Context()), filter)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "搜索记录查询失败")
		return
	}
	rows := make([]row, 0)
	if err := applyQuery(h.business.WithContext(c.Request.Context()), filter).
		Select(`r.id,r.user_id,COALESCE(u.nickname,'') AS nickname,COALESCE(p.avatar_url,'') AS avatar_url,
COALESCE(p.source_channel,'') AS user_type,r.keyword,r.created_at`).
		Order("r.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "搜索记录查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Clear(c *gin.Context) {
	var in clearInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "清理搜索记录参数错误")
		return
	}
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if in.Reason == "" {
		in.Reason = "一键清空搜索记录"
	}
	if len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len(in.IdempotencyKey) < 8 || len(in.IdempotencyKey) > 128 {
		response.Fail(c, http.StatusBadRequest, "清理搜索记录参数错误")
		return
	}
	operator := middleware.AdminID(c)
	var out clearAudit
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// user_id=0 表示一键清空全部可见搜索记录。
		insert := tx.Exec(`INSERT INTO qixi_crm_b_user_search_record_clear_audit (user_id,reason,idempotency_key,operator_admin_id,cleared_count)
VALUES (0,?,?,?,0) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id)`, in.Reason, in.IdempotencyKey, operator)
		if insert.Error != nil {
			return insert.Error
		}
		if err := tx.Table("qixi_crm_b_user_search_record_clear_audit").Where("operator_admin_id=? AND idempotency_key=?", operator, in.IdempotencyKey).Take(&out).Error; err != nil {
			return err
		}
		if out.UserID != 0 || out.Reason != in.Reason {
			return errClearConflict
		}
		if insert.RowsAffected == 0 {
			return nil
		}
		cleared := tx.Table("qixi_crm_b_user_search_record").Where("deleted_at IS NULL").Update("deleted_at", time.Now())
		if cleared.Error != nil {
			return cleared.Error
		}
		out.ClearedCount = cleared.RowsAffected
		return tx.Table("qixi_crm_b_user_search_record_clear_audit").Where("operator_admin_id=? AND idempotency_key=?", operator, in.IdempotencyKey).Update("cleared_count", cleared.RowsAffected).Error
	})
	if err != nil {
		if errors.Is(err, errClearConflict) {
			response.Fail(c, http.StatusConflict, "幂等键已用于不同清理请求")
			return
		}
		fail(c, "搜索记录清理失败")
		return
	}
	response.OK(c, gin.H{"cleared_count": out.ClearedCount})
}

func (h *Handler) Export(c *gin.Context) {
	var in exportInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "导出搜索记录参数错误")
		return
	}
	filter, ok := parseQuery(c, in.query)
	if !ok {
		return
	}
	in.Reason = strings.TrimSpace(in.Reason)
	if in.Reason == "" {
		in.Reason = "平台后台导出搜索记录"
	}
	if len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 {
		response.Fail(c, http.StatusBadRequest, "导出原因错误")
		return
	}
	rows := make([]row, 0)
	if err := applyQuery(h.business.WithContext(c.Request.Context()), filter).
		Select(`r.id,r.user_id,COALESCE(u.nickname,'') AS nickname,COALESCE(p.avatar_url,'') AS avatar_url,
COALESCE(p.source_channel,'') AS user_type,r.keyword,r.created_at`).
		Order("r.id DESC").Limit(5000).Scan(&rows).Error; err != nil {
		fail(c, "搜索记录导出查询失败")
		return
	}
	content, err := makeCSV(rows)
	if err != nil {
		fail(c, "搜索记录导出生成失败")
		return
	}
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_user_search_record_export_audit").Create(map[string]any{"query_fingerprint": fingerprint(filter), "row_count": len(rows), "reason": in.Reason, "operator_admin_id": middleware.AdminID(c)}).Error; err != nil {
		fail(c, "搜索记录导出审计写入失败")
		return
	}
	response.OK(c, gin.H{"file_name": "用户搜索记录导出_" + time.Now().Format("20060102150405") + ".csv", "content": content, "row_count": len(rows), "truncated": len(rows) == 5000})
}

func parseQuery(c *gin.Context, in query) (query, bool) {
	in.Keyword = strings.TrimSpace(in.Keyword)
	in.Nickname = strings.TrimSpace(in.Nickname)
	in.UserType = strings.TrimSpace(in.UserType)
	in.StartDate = strings.TrimSpace(in.StartDate)
	in.EndDate = strings.TrimSpace(in.EndDate)
	if len([]rune(in.Keyword)) > 128 || len([]rune(in.Nickname)) > 64 || (in.UserType != "" && !validUserType(in.UserType)) || !validDateRange(in.StartDate, in.EndDate) {
		response.Fail(c, http.StatusBadRequest, "搜索记录筛选参数错误")
		return query{}, false
	}
	return in, true
}

func applyQuery(db *gorm.DB, in query) *gorm.DB {
	q := db.Table("qixi_crm_b_user_search_record AS r").
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id = r.user_id").
		Joins("LEFT JOIN qixi_crm_b_user_profile AS p ON p.user_id = r.user_id").
		Where("r.deleted_at IS NULL")
	if in.Keyword != "" {
		q = q.Where("r.keyword LIKE ?", "%"+in.Keyword+"%")
	}
	if in.Nickname != "" {
		q = q.Where("u.nickname LIKE ?", "%"+in.Nickname+"%")
	}
	switch in.UserType {
	case "wechat", "mini_program", "h5", "pc":
		q = q.Where("p.source_channel=?", in.UserType)
	case "app":
		q = q.Where("p.source_channel IN ?", []string{"ios", "android", "harmony"})
	}
	if in.StartDate != "" {
		q = q.Where("r.created_at>=?", in.StartDate+" 00:00:00")
	}
	if in.EndDate != "" {
		q = q.Where("r.created_at<?", nextDate(in.EndDate))
	}
	return q
}

func validUserType(value string) bool {
	switch value {
	case "wechat", "mini_program", "h5", "app", "pc":
		return true
	default:
		return false
	}
}

func validDateRange(start, end string) bool {
	if start == "" && end == "" {
		return true
	}
	var s, e time.Time
	var err error
	if start != "" {
		if s, err = time.Parse("2006-01-02", start); err != nil {
			return false
		}
	}
	if end != "" {
		if e, err = time.Parse("2006-01-02", end); err != nil || (!s.IsZero() && e.Before(s)) {
			return false
		}
	}
	return true
}

func nextDate(value string) string {
	t, _ := time.Parse("2006-01-02", value)
	return t.AddDate(0, 0, 1).Format("2006-01-02 15:04:05")
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

func userTypeLabel(value string) string {
	switch value {
	case "wechat":
		return "微信用户"
	case "mini_program":
		return "小程序用户"
	case "h5":
		return "H5用户"
	case "pc":
		return "PC用户"
	case "ios", "android", "harmony":
		return "APP用户"
	default:
		return value
	}
}

func makeCSV(rows []row) (string, error) {
	var out bytes.Buffer
	out.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&out)
	if err := w.Write([]string{"用户ID", "昵称", "用户类型", "搜索词", "搜索时间"}); err != nil {
		return "", err
	}
	for _, item := range rows {
		if err := w.Write([]string{
			strconv.FormatUint(item.UserID, 10),
			csvCell(item.Nickname),
			userTypeLabel(item.UserType),
			csvCell(item.Keyword),
			item.CreatedAt,
		}); err != nil {
			return "", err
		}
	}
	w.Flush()
	return out.String(), w.Error()
}

func csvCell(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "=") || strings.HasPrefix(value, "+") || strings.HasPrefix(value, "-") || strings.HasPrefix(value, "@") {
		return "'" + value
	}
	return value
}

func fingerprint(in query) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{in.Keyword, in.Nickname, in.UserType, in.StartDate, in.EndDate}, "|")))
	return hex.EncodeToString(sum[:])
}

func fail(c *gin.Context, message string) { response.Fail(c, http.StatusInternalServerError, message) }
