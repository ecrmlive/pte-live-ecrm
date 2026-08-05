// Package usersearch exposes the platform-only supervision plane for C-end
// search history. Search terms can be sensitive, so it deliberately returns
// no account, mobile, address or other user profile fields.
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
	Keyword   string `gorm:"column:keyword" json:"keyword"`
	Source    string `gorm:"column:source" json:"source"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}

type query struct {
	UserID    uint64 `form:"user_id" json:"user_id"`
	Keyword   string `form:"keyword" json:"keyword"`
	Source    string `form:"source" json:"source"`
	StartDate string `form:"start_date" json:"start_date"`
	EndDate   string `form:"end_date" json:"end_date"`
}

type clearInput struct {
	UserID         uint64 `json:"user_id"`
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
	userID, validID := queryUserID(c)
	if !validID {
		response.Fail(c, http.StatusBadRequest, "用户 ID 参数错误")
		return
	}
	filter, ok := parseQuery(c, query{UserID: userID, Keyword: c.Query("keyword"), Source: c.Query("source"), StartDate: c.Query("start_date"), EndDate: c.Query("end_date")})
	if !ok {
		return
	}
	page, limit := pagination(c)
	q := applyQuery(h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_user_search_record"), filter)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "搜索记录查询失败")
		return
	}
	rows := make([]row, 0)
	if err := q.Select("id,user_id,keyword,source,created_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
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
	if in.UserID == 0 || len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len(in.IdempotencyKey) < 8 || len(in.IdempotencyKey) > 128 {
		response.Fail(c, http.StatusBadRequest, "清理搜索记录参数错误")
		return
	}
	operator := middleware.AdminID(c)
	var out clearAudit
	err := h.business.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		insert := tx.Exec(`INSERT INTO qixi_crm_b_user_search_record_clear_audit (user_id,reason,idempotency_key,operator_admin_id,cleared_count)
VALUES (?,?,?,?,0) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id)`, in.UserID, in.Reason, in.IdempotencyKey, operator)
		if insert.Error != nil {
			return insert.Error
		}
		if err := tx.Table("qixi_crm_b_user_search_record_clear_audit").Where("operator_admin_id=? AND idempotency_key=?", operator, in.IdempotencyKey).Take(&out).Error; err != nil {
			return err
		}
		if out.UserID != in.UserID || out.Reason != in.Reason {
			return errClearConflict
		}
		if insert.RowsAffected == 0 {
			return nil
		}
		cleared := tx.Table("qixi_crm_b_user_search_record").Where("user_id=? AND deleted_at IS NULL", in.UserID).Update("deleted_at", time.Now())
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
	response.OK(c, gin.H{"user_id": in.UserID, "cleared_count": out.ClearedCount})
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
	if len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 {
		response.Fail(c, http.StatusBadRequest, "导出原因错误")
		return
	}
	rows := make([]row, 0)
	if err := applyQuery(h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_user_search_record"), filter).Select("id,user_id,keyword,source,created_at").Order("id DESC").Limit(5000).Scan(&rows).Error; err != nil {
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
	in.Keyword, in.Source, in.StartDate, in.EndDate = strings.TrimSpace(in.Keyword), strings.TrimSpace(in.Source), strings.TrimSpace(in.StartDate), strings.TrimSpace(in.EndDate)
	if len([]rune(in.Keyword)) > 128 || (in.Source != "" && !validSource(in.Source)) || !validDateRange(in.StartDate, in.EndDate) {
		response.Fail(c, http.StatusBadRequest, "搜索记录筛选参数错误")
		return query{}, false
	}
	return in, true
}

func applyQuery(q *gorm.DB, in query) *gorm.DB {
	q = q.Where("deleted_at IS NULL")
	if in.UserID != 0 {
		q = q.Where("user_id=?", in.UserID)
	}
	if in.Keyword != "" {
		q = q.Where("keyword LIKE ?", "%"+in.Keyword+"%")
	}
	if in.Source != "" {
		q = q.Where("source=?", in.Source)
	}
	if in.StartDate != "" {
		q = q.Where("created_at>=?", in.StartDate+" 00:00:00")
	}
	if in.EndDate != "" {
		q = q.Where("created_at<?", nextDate(in.EndDate))
	}
	return q
}

func validSource(value string) bool { return value == "pc" || value == "h5" || value == "mini" }
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
func queryUserID(c *gin.Context) (uint64, bool) {
	raw := strings.TrimSpace(c.Query("user_id"))
	if raw == "" {
		return 0, true
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	return value, err == nil && value > 0
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
func makeCSV(rows []row) (string, error) {
	var out bytes.Buffer
	out.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&out)
	if err := w.Write([]string{"记录ID", "用户ID", "搜索关键词", "来源", "搜索时间"}); err != nil {
		return "", err
	}
	for _, item := range rows {
		if err := w.Write([]string{strconv.FormatUint(item.ID, 10), strconv.FormatUint(item.UserID, 10), csvCell(item.Keyword), item.Source, item.CreatedAt}); err != nil {
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
	sum := sha256.Sum256([]byte(strings.Join([]string{strconv.FormatUint(in.UserID, 10), in.Keyword, in.Source, in.StartDate, in.EndDate}, "|")))
	return hex.EncodeToString(sum[:])
}
func fail(c *gin.Context, message string) { response.Fail(c, http.StatusInternalServerError, message) }
