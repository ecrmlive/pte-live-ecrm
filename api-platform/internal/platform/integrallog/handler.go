// Package integrallog exposes platform integral (points) bill list, summary and CSV export.
// Data lives in qixi_crm_b_user_bill (+ member_account / user_sign for title stats).
package integrallog

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	// RequireAdminMenu 仅认 kind=button；page 码 marketing.integral.log 只用于导航。
	menuIntegralLogRead = "marketing.integral.log.read"
	exportLimit         = 5000
)

type Handler struct {
	businessDB *gorm.DB
	adminDB    *gorm.DB
}

func NewHandler(businessDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	read := middleware.RequireAdminMenu(h.adminDB, menuIntegralLogRead)
	r.GET("/integral/logs", access, read, h.List)
	r.GET("/integral/logs/title", access, read, h.Title)
	r.POST("/integral/logs/export", access, read, h.Export)
}

type billRow struct {
	BillID     uint64    `gorm:"column:bill_id" json:"bill_id"`
	UID        uint64    `gorm:"column:uid" json:"uid"`
	Nickname   string    `gorm:"column:nickname" json:"nickname"`
	Title      string    `gorm:"column:title" json:"title"`
	PM         int8      `gorm:"column:pm" json:"pm"`
	Number     float64   `gorm:"column:number" json:"number"`
	Balance    float64   `gorm:"column:balance" json:"balance"`
	Mark       string    `gorm:"column:mark" json:"mark"`
	Category   string    `gorm:"column:category" json:"category"`
	Type       string    `gorm:"column:type" json:"type"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

type titleStat struct {
	TotalIntegral  float64 `json:"total_integral"`
	SignCount      int64   `json:"sign_count"`
	SignIntegral   float64 `json:"sign_integral"`
	UsedIntegral   float64 `json:"used_integral"`
	OrderIntegral  float64 `json:"order_integral"`
	FreezeIntegral float64 `json:"freeze_integral"`
}

type listFilter struct {
	Keyword  string
	DateFrom string
	DateTo   string
}

func (h *Handler) List(c *gin.Context) {
	page, limit := queryfilter.Page(c)
	f := parseListFilter(c)
	q := h.baseQuery(c, f)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "积分日志查询失败")
		return
	}
	rows := make([]billRow, 0)
	if err := q.Select("a.bill_id,a.uid,COALESCE(b.nickname,'') AS nickname,a.title,a.pm,a.number,a.balance,a.mark,a.category,a.type,a.create_time").
		Order("a.create_time DESC, a.bill_id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&rows).Error; err != nil {
		fail(c, "积分日志查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Title(c *gin.Context) {
	ctx := c.Request.Context()
	out := titleStat{}

	// 总积分：启用用户会员账户积分合计
	if err := h.businessDB.WithContext(ctx).
		Table("qixi_crm_b_member_account AS a").
		Joins("INNER JOIN qixi_crm_b_user AS u ON u.id = a.user_id").
		Where("u.status = ?", 1).
		Select("COALESCE(SUM(a.points),0)").
		Scan(&out.TotalIntegral).Error; err != nil {
		fail(c, "积分统计失败")
		return
	}

	if err := h.businessDB.WithContext(ctx).Table("qixi_crm_b_user_sign").Count(&out.SignCount).Error; err != nil {
		fail(c, "积分统计失败")
		return
	}

	if err := h.businessDB.WithContext(ctx).Table("qixi_crm_b_user_bill").
		Where("type = ?", "sign_integral").
		Select("COALESCE(SUM(number),0)").Scan(&out.SignIntegral).Error; err != nil {
		fail(c, "积分统计失败")
		return
	}

	var deduction, merRefund float64
	if err := h.businessDB.WithContext(ctx).Table("qixi_crm_b_user_bill").
		Where("category = ? AND type = ?", "integral", "deduction").
		Select("COALESCE(SUM(number),0)").Scan(&deduction).Error; err != nil {
		fail(c, "积分统计失败")
		return
	}
	if err := h.businessDB.WithContext(ctx).Table("qixi_crm_b_user_bill").
		Where("category = ? AND type = ?", "mer_integral", "refund").
		Select("COALESCE(SUM(number),0)").Scan(&merRefund).Error; err != nil {
		fail(c, "积分统计失败")
		return
	}
	out.UsedIntegral = deduction - merRefund

	var lockSum, refundLockSum float64
	if err := h.businessDB.WithContext(ctx).Table("qixi_crm_b_user_bill").
		Where("category = ? AND type = ?", "integral", "lock").
		Select("COALESCE(SUM(number),0)").Scan(&lockSum).Error; err != nil {
		fail(c, "积分统计失败")
		return
	}
	if err := h.businessDB.WithContext(ctx).Table("qixi_crm_b_user_bill").
		Where("category = ? AND type = ?", "integral", "refund_lock").
		Select("COALESCE(SUM(number),0)").Scan(&refundLockSum).Error; err != nil {
		fail(c, "积分统计失败")
		return
	}
	out.OrderIntegral = lockSum - refundLockSum
	if out.OrderIntegral < 0 {
		out.OrderIntegral = 0
	}

	out.FreezeIntegral = h.freezeIntegral(c)

	response.OK(c, out)
}

type exportInput struct {
	Keyword  string `json:"keyword"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (h *Handler) Export(c *gin.Context) {
	var in exportInput
	_ = c.ShouldBindJSON(&in)
	f := listFilter{
		Keyword:  strings.TrimSpace(in.Keyword),
		DateFrom: strings.TrimSpace(in.DateFrom),
		DateTo:   strings.TrimSpace(in.DateTo),
	}
	if f.Keyword == "" {
		f.Keyword = strings.TrimSpace(c.Query("keyword"))
	}
	if f.DateFrom == "" {
		f.DateFrom = strings.TrimSpace(c.Query("date_from"))
	}
	if f.DateTo == "" {
		f.DateTo = strings.TrimSpace(c.Query("date_to"))
	}

	q := h.baseQuery(c, f)
	rows := make([]billRow, 0)
	if err := q.Select("a.bill_id,a.uid,COALESCE(b.nickname,'') AS nickname,a.title,a.pm,a.number,a.balance,a.mark,a.category,a.type,a.create_time").
		Order("a.create_time DESC, a.bill_id DESC").
		Limit(exportLimit).
		Scan(&rows).Error; err != nil {
		fail(c, "积分日志导出失败")
		return
	}
	content, err := billsCSV(rows)
	if err != nil {
		fail(c, "积分日志导出失败")
		return
	}
	response.OK(c, gin.H{
		"file_name":  "积分日志_" + time.Now().Format("20060102150405") + ".csv",
		"content":    content,
		"row_count":  len(rows),
		"truncated": len(rows) == exportLimit,
	})
}

func (h *Handler) freezeIntegral(c *gin.Context) float64 {
	ctx := c.Request.Context()
	type lockRow struct {
		LinkID string  `gorm:"column:link_id"`
		Number float64 `gorm:"column:number"`
	}
	locks := make([]lockRow, 0)
	if err := h.businessDB.WithContext(ctx).Table("qixi_crm_b_user_bill").
		Select("link_id,number").
		Where("category = ? AND type = ? AND status = ?", "integral", "lock", 0).
		Scan(&locks).Error; err != nil || len(locks) == 0 {
		return 0
	}
	linkIDs := make([]string, 0, len(locks))
	var sum float64
	for _, row := range locks {
		sum += row.Number
		if row.LinkID != "" {
			linkIDs = append(linkIDs, row.LinkID)
		}
	}
	if len(linkIDs) == 0 {
		return sum
	}
	var refunded float64
	_ = h.businessDB.WithContext(ctx).Table("qixi_crm_b_user_bill").
		Where("category = ? AND type = ? AND link_id IN ?", "integral", "refund_lock", linkIDs).
		Select("COALESCE(SUM(number),0)").Scan(&refunded).Error
	frozen := sum - refunded
	if frozen < 0 {
		return 0
	}
	return frozen
}

func (h *Handler) baseQuery(c *gin.Context, f listFilter) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_user_bill AS a").
		Joins("LEFT JOIN qixi_crm_b_user AS b ON b.id = a.uid").
		Where("a.category IN ?", []string{"integral", "mer_integral"})
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		q = q.Where("CAST(a.uid AS CHAR) LIKE ? OR b.nickname LIKE ? OR a.title LIKE ?", like, like, like)
	}
	q = applyCreateTimeRange(q, f.DateFrom, f.DateTo, "a.create_time")
	return q
}

func parseListFilter(c *gin.Context) listFilter {
	return listFilter{
		Keyword:  strings.TrimSpace(c.Query("keyword")),
		DateFrom: strings.TrimSpace(c.Query("date_from")),
		DateTo:   strings.TrimSpace(c.Query("date_to")),
	}
}

func applyCreateTimeRange(q *gorm.DB, from, to, column string) *gorm.DB {
	if from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where(column+" >= ?", t)
		}
	}
	if to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where(column+" < ?", t.AddDate(0, 0, 1))
		}
	}
	return q
}

func billsCSV(rows []billRow) (string, error) {
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	if err := w.Write([]string{"ID", "用户昵称", "积分标题", "积分变动", "当前积分额度", "备注", "添加时间"}); err != nil {
		return "", err
	}
	for _, row := range rows {
		change := formatSignedNumber(row.PM, row.Number)
		if err := w.Write([]string{
			strconv.FormatUint(row.BillID, 10),
			csvSafe(row.Nickname),
			csvSafe(row.Title),
			change,
			formatPlain(row.Balance),
			csvSafe(row.Mark),
			row.CreateTime.Format("2006-01-02 15:04:05"),
		}); err != nil {
			return "", err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func formatSignedNumber(pm int8, number float64) string {
	n := formatPlain(number)
	if pm == 1 {
		return "+" + n
	}
	return "-" + n
}

func formatPlain(v float64) string {
	if v == float64(int64(v)) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func csvSafe(s string) string {
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '+', '-', '@':
		return "'" + s
	default:
		return s
	}
}

func fail(c *gin.Context, msg string) {
	response.Fail(c, http.StatusInternalServerError, msg)
}
