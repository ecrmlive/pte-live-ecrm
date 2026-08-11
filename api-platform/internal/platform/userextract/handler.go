// Package userextract exposes platform user brokerage withdrawal (提现管理)
// list, detail, audit and CSV export. Aligns with CRMEB UserExtract.
// Data lives in qixi_crm_b_user_extract (+ qixi_crm_b_user for nickname).
package userextract

import (
	"bytes"
	"encoding/csv"
	"fmt"
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
	menuRead   = "accounts.withdraw.read"
	menuReview = "accounts.withdraw.review"
	menuExport = "accounts.withdraw.export"
	exportLimit = 5000
	tableName   = "qixi_crm_b_user_extract"
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
	read := middleware.RequireAdminMenu(h.adminDB, menuRead)
	review := middleware.RequireAdminMenu(h.adminDB, menuReview)
	export := middleware.RequireAdminMenu(h.adminDB, menuExport)
	r.GET("/finance/user-extracts", access, read, h.List)
	r.GET("/finance/user-extracts/:id", access, read, h.Detail)
	r.POST("/finance/user-extracts/:id/status", access, review, h.SwitchStatus)
	r.POST("/finance/user-extracts/export", access, export, h.Export)
}

type listFilter struct {
	DateFrom     string
	DateTo       string
	Status       *int
	ExtractType  *int
	UserType     string
	UserKeyword  string
	AccountKeyword string
}

type row struct {
	ExtractID    uint64     `gorm:"column:extract_id" json:"extract_id"`
	UID          uint64     `gorm:"column:uid" json:"uid"`
	Nickname     string     `gorm:"column:nickname" json:"nickname"`
	ExtractSN    string     `gorm:"column:extract_sn" json:"extract_sn"`
	RealName     string     `gorm:"column:real_name" json:"real_name"`
	ExtractType  int        `gorm:"column:extract_type" json:"extract_type"`
	ExtractTypeName string  `gorm:"-" json:"extract_type_name"`
	BankCode     string     `gorm:"column:bank_code" json:"bank_code"`
	BankAddress  string     `gorm:"column:bank_address" json:"bank_address"`
	BankName     string     `gorm:"column:bank_name" json:"bank_name"`
	AlipayCode   string     `gorm:"column:alipay_code" json:"alipay_code"`
	Wechat       string     `gorm:"column:wechat" json:"wechat"`
	ExtractPic   string     `gorm:"column:extract_pic" json:"extract_pic"`
	ExtractPrice float64    `gorm:"column:extract_price" json:"extract_price"`
	Balance      float64    `gorm:"column:balance" json:"balance"`
	Mark         string     `gorm:"column:mark" json:"mark"`
	AdminID      uint64     `gorm:"column:admin_id" json:"admin_id"`
	FailMsg      string     `gorm:"column:fail_msg" json:"fail_msg"`
	Status       int        `gorm:"column:status" json:"status"`
	StatusName   string     `gorm:"-" json:"status_name"`
	Account      string     `gorm:"-" json:"account"`
	StatusTime   *time.Time `gorm:"column:status_time" json:"status_time"`
	CreateTime   *time.Time `gorm:"column:create_time" json:"create_time"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := queryfilter.Page(c)
	f := parseFilter(c)
	var total int64
	if err := h.baseQuery(c, f).Count(&total).Error; err != nil {
		fail(c, "提现列表查询失败")
		return
	}
	rows := make([]row, 0)
	err := h.baseQuery(c, f).Select(`e.extract_id,e.uid,COALESCE(u.nickname,'') AS nickname,e.extract_sn,e.real_name,
e.extract_type,e.bank_code,e.bank_address,e.bank_name,e.alipay_code,e.wechat,e.extract_pic,
e.extract_price,e.balance,e.mark,e.admin_id,e.fail_msg,e.status,e.status_time,e.create_time`).
		Order("e.extract_id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		fail(c, "提现列表查询失败")
		return
	}
	for i := range rows {
		decorate(&rows[i])
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Detail(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var item row
	err := h.businessDB.WithContext(c.Request.Context()).Table(tableName+" AS e").
		Joins("LEFT JOIN qixi_crm_b_user u ON u.id = e.uid").
		Select(`e.extract_id,e.uid,COALESCE(u.nickname,'') AS nickname,e.extract_sn,e.real_name,
e.extract_type,e.bank_code,e.bank_address,e.bank_name,e.alipay_code,e.wechat,e.extract_pic,
e.extract_price,e.balance,e.mark,e.admin_id,e.fail_msg,e.status,e.status_time,e.create_time`).
		Where("e.extract_id = ?", id).
		Take(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, "提现记录不存在")
			return
		}
		fail(c, "提现详情查询失败")
		return
	}
	decorate(&item)
	response.OK(c, item)
}

func (h *Handler) SwitchStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req struct {
		Status  int    `json:"status"`
		FailMsg string `json:"fail_msg"`
		Mark    string `json:"mark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "审核参数错误")
		return
	}
	if req.Status != 1 && req.Status != -1 {
		fail(c, "审核状态错误")
		return
	}
	req.FailMsg = strings.TrimSpace(req.FailMsg)
	req.Mark = strings.TrimSpace(req.Mark)
	if req.Status == -1 {
		if req.FailMsg == "" || len([]rune(req.FailMsg)) > 200 {
			fail(c, "请填写拒绝原因（不超过 200 字）")
			return
		}
	} else {
		req.FailMsg = ""
	}
	if len([]rune(req.Mark)) > 500 {
		fail(c, "备注不能超过 500 字")
		return
	}

	var current row
	err := h.businessDB.WithContext(c.Request.Context()).Table(tableName).
		Select("extract_id,status").
		Where("extract_id = ?", id).
		Take(&current).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, "提现记录不存在")
			return
		}
		fail(c, "提现审核失败")
		return
	}
	if current.Status != 0 {
		fail(c, "数据不存在或状态错误")
		return
	}

	now := time.Now()
	result := h.businessDB.WithContext(c.Request.Context()).Table(tableName).
		Where("extract_id = ? AND status = ?", id, 0).
		Updates(map[string]any{
			"status":      req.Status,
			"fail_msg":    req.FailMsg,
			"mark":        req.Mark,
			"admin_id":    middleware.AdminID(c),
			"status_time": now,
		})
	if result.Error != nil {
		fail(c, "提现审核失败")
		return
	}
	if result.RowsAffected == 0 {
		fail(c, "数据不存在或状态错误")
		return
	}
	response.OK(c, gin.H{"ok": true, "status": req.Status})
}

func (h *Handler) Export(c *gin.Context) {
	f := parseFilter(c)
	var body struct {
		DateFrom       string `json:"date_from"`
		DateTo         string `json:"date_to"`
		Status         *int   `json:"status"`
		ExtractType    *int   `json:"extract_type"`
		UserType       string `json:"user_type"`
		UserKeyword    string `json:"user_keyword"`
		AccountKeyword string `json:"account_keyword"`
	}
	if err := c.ShouldBindJSON(&body); err == nil {
		if body.DateFrom != "" {
			f.DateFrom = strings.TrimSpace(body.DateFrom)
		}
		if body.DateTo != "" {
			f.DateTo = strings.TrimSpace(body.DateTo)
		}
		if body.UserType != "" {
			f.UserType = strings.TrimSpace(body.UserType)
		}
		if body.UserKeyword != "" {
			f.UserKeyword = strings.TrimSpace(body.UserKeyword)
		}
		if body.AccountKeyword != "" {
			f.AccountKeyword = strings.TrimSpace(body.AccountKeyword)
		}
		if body.Status != nil {
			f.Status = body.Status
		}
		if body.ExtractType != nil {
			f.ExtractType = body.ExtractType
		}
	}
	q := h.baseQuery(c, f)
	rows := make([]row, 0)
	err := q.Select(`e.extract_id,e.uid,COALESCE(u.nickname,'') AS nickname,e.extract_sn,e.real_name,
e.extract_type,e.bank_code,e.bank_address,e.bank_name,e.alipay_code,e.wechat,e.extract_pic,
e.extract_price,e.balance,e.mark,e.admin_id,e.fail_msg,e.status,e.status_time,e.create_time`).
		Order("e.extract_id DESC").
		Limit(exportLimit + 1).
		Scan(&rows).Error
	if err != nil {
		fail(c, "提现导出失败")
		return
	}
	truncated := len(rows) > exportLimit
	if truncated {
		rows = rows[:exportLimit]
	}
	for i := range rows {
		decorate(&rows[i])
	}

	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"序号", "用户信息", "用户UID", "户名", "提现金额", "提现方式", "银行名称", "账号", "提现状态", "拒绝原因", "添加时间"})
	for i, item := range rows {
		_ = w.Write([]string{
			strconv.Itoa(i + 1),
			item.Nickname,
			strconv.FormatUint(item.UID, 10),
			dash(item.RealName),
			fmt.Sprintf("%.2f", item.ExtractPrice),
			item.ExtractTypeName,
			dash(item.BankName),
			dash(item.Account),
			item.StatusName,
			dash(item.FailMsg),
			formatTime(item.CreateTime),
		})
	}
	w.Flush()
	response.OK(c, gin.H{
		"content":   buf.String(),
		"file_name": fmt.Sprintf("提现管理_%s.csv", time.Now().Format("20060102150405")),
		"row_count": len(rows),
		"truncated": truncated,
	})
}

func (h *Handler) baseQuery(c *gin.Context, f listFilter) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).Table(tableName + " AS e").
		Joins("LEFT JOIN qixi_crm_b_user u ON u.id = e.uid")
	q = applyDateRange(q, f.DateFrom, f.DateTo, "e.create_time")
	if f.Status != nil {
		q = q.Where("e.status = ?", *f.Status)
	}
	if f.ExtractType != nil {
		q = q.Where("e.extract_type = ?", *f.ExtractType)
	}
	if kw := strings.TrimSpace(f.UserKeyword); kw != "" {
		switch strings.TrimSpace(f.UserType) {
		case "uid":
			if id, err := strconv.ParseUint(kw, 10, 64); err == nil {
				q = q.Where("e.uid = ?", id)
			} else {
				q = q.Where("1 = 0")
			}
		case "phone":
			q = q.Where("u.mobile LIKE ?", "%"+kw+"%")
		case "real_name":
			q = q.Where("e.real_name LIKE ?", "%"+kw+"%")
		default:
			q = q.Where("u.nickname LIKE ?", "%"+kw+"%")
		}
	}
	if ak := strings.TrimSpace(f.AccountKeyword); ak != "" {
		like := "%" + ak + "%"
		q = q.Where("(e.bank_code LIKE ? OR e.alipay_code LIKE ? OR e.wechat LIKE ?)", like, like, like)
	}
	return q
}

func parseFilter(c *gin.Context) listFilter {
	f := listFilter{
		DateFrom:       strings.TrimSpace(c.Query("date_from")),
		DateTo:         strings.TrimSpace(c.Query("date_to")),
		UserType:       strings.TrimSpace(c.Query("user_type")),
		UserKeyword:    strings.TrimSpace(c.Query("user_keyword")),
		AccountKeyword: strings.TrimSpace(c.Query("account_keyword")),
	}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			f.Status = &v
		}
	}
	if raw := strings.TrimSpace(c.Query("extract_type")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			f.ExtractType = &v
		}
	}
	return f
}

func applyDateRange(q *gorm.DB, from, to, column string) *gorm.DB {
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

func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		fail(c, "提现 ID 无效")
		return 0, false
	}
	return id, true
}

func decorate(item *row) {
	item.ExtractTypeName = typeName(item.ExtractType)
	item.StatusName = statusName(item.Status)
	item.Account = accountOf(item)
	if item.Nickname == "" {
		item.Nickname = "用户"
	}
}

func accountOf(item *row) string {
	switch item.ExtractType {
	case 0:
		return item.BankCode
	case 1, 3:
		return item.Wechat
	case 2:
		return item.AlipayCode
	default:
		return ""
	}
}

func typeName(t int) string {
	switch t {
	case 0:
		return "银行卡"
	case 1:
		return "微信"
	case 2:
		return "支付宝"
	case 3:
		return "微信零钱"
	case 4:
		return "余额"
	default:
		return "未知"
	}
}

func statusName(status int) string {
	switch status {
	case 1:
		return "已通过"
	case 0:
		return "审核中"
	case -1:
		return "已拒绝"
	default:
		return "未知"
	}
}

func dash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func fail(c *gin.Context, msg string) {
	response.Fail(c, http.StatusBadRequest, msg)
}
