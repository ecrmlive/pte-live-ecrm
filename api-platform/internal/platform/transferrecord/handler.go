// Package transferrecord exposes platform merchant transfer (提现转账) list,
// title stats, detail, audit, mark, payout voucher and CSV export.
// Data lives in qixi_crm_a_financial (+ qixi_crm_a_merchant_view for shop fields).
package transferrecord

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
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
	menuTransferRead   = "accounts.transfer.read"
	menuTransferManage = "accounts.transfer.manage"
	menuTransferExport = "accounts.transfer.export"
	exportLimit        = 5000
	extractMinimumLine = 0.0 // 对齐 CRMEB extract_minimum_line，本地默认 0
	typeExtract        = 0
)

type Handler struct {
	adminDB *gorm.DB
}

func NewHandler(adminDB *gorm.DB) *Handler {
	return &Handler{adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	read := middleware.RequireAdminMenu(h.adminDB, menuTransferRead)
	manage := middleware.RequireAdminMenu(h.adminDB, menuTransferManage)
	export := middleware.RequireAdminMenu(h.adminDB, menuTransferExport)
	r.GET("/finance/transfers/title", access, read, h.Title)
	r.GET("/finance/transfers", access, read, h.List)
	r.GET("/finance/transfers/:id", access, read, h.Detail)
	r.POST("/finance/transfers/:id/status", access, manage, h.SwitchStatus)
	r.POST("/finance/transfers/:id/mark", access, manage, h.Mark)
	r.POST("/finance/transfers/:id/pay", access, manage, h.Pay)
	r.POST("/finance/transfers/export", access, export, h.Export)
}

type listFilter struct {
	DateFrom        string
	DateTo          string
	Status          *int
	MerName         string
	IsTrader        *int
	FinancialType   *int
	FinancialStatus *int
	AdminKeyword    string
}

type titleStat struct {
	PayableAmount       float64 `json:"payable_amount"`
	WithdrawableAmount  float64 `json:"withdrawable_amount"`
	ApplyingMerchantCnt int64   `json:"applying_merchant_count"`
	ApplyingAmount      float64 `json:"applying_amount"`
	PendingAmount       float64 `json:"pending_amount"`
	FreezeAmount        float64 `json:"freeze_amount"`
}

type transferRow struct {
	FinancialID      uint64     `gorm:"column:financial_id" json:"financial_id"`
	FinancialSN      string     `gorm:"column:financial_sn" json:"financial_sn"`
	MerID            uint64     `gorm:"column:mer_id" json:"mer_id"`
	MerName          string     `gorm:"column:mer_name" json:"mer_name"`
	IsTrader         int        `gorm:"column:is_trader" json:"is_trader"`
	MerMoney         float64    `gorm:"column:mer_money" json:"mer_money"`
	ExtractMoney     float64    `gorm:"column:extract_money" json:"extract_money"`
	FinancialType    int        `gorm:"column:financial_type" json:"financial_type"`
	FinancialAccount string     `gorm:"column:financial_account" json:"financial_account"`
	FinancialStatus  int        `gorm:"column:financial_status" json:"financial_status"`
	Status           int        `gorm:"column:status" json:"status"`
	Refusal          string     `gorm:"column:refusal" json:"refusal"`
	Image            string     `gorm:"column:image" json:"image"`
	AdminID          *uint64    `gorm:"column:admin_id" json:"admin_id"`
	AdminName        string     `gorm:"column:admin_name" json:"admin_name"`
	Mark             string     `gorm:"column:mark" json:"mark"`
	AdminMark        string     `gorm:"column:admin_mark" json:"admin_mark"`
	CreateTime       *time.Time `gorm:"column:create_time" json:"create_time"`
	StatusTime       *time.Time `gorm:"column:status_time" json:"status_time"`
	UpdateTime       *time.Time `gorm:"column:update_time" json:"update_time"`
}

type accountInfo struct {
	Name       string `json:"name"`
	Bank       string `json:"bank"`
	BankCode   string `json:"bank_code"`
	Wechat     string `json:"wechat"`
	WechatCode string `json:"wechat_code"`
	Alipay     string `json:"alipay"`
	AlipayCode string `json:"alipay_code"`
}

func (h *Handler) Title(c *gin.Context) {
	ctx := c.Request.Context()
	out := titleStat{}

	_ = h.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_view").
		Select("COALESCE(SUM(mer_money),0)").Scan(&out.PayableAmount).Error

	line := extractMinimumLine
	var withdrawable float64
	_ = h.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_view").
		Where("mer_money > ?", line).
		Select(fmt.Sprintf("COALESCE(SUM(mer_money - %v),0)", line)).
		Scan(&withdrawable).Error
	out.WithdrawableAmount = round2(withdrawable)

	_ = h.adminDB.WithContext(ctx).Table("qixi_crm_a_financial").
		Where("is_del = 0 AND type = ? AND financial_status = 0 AND status > -1", typeExtract).
		Distinct("mer_id").Count(&out.ApplyingMerchantCnt).Error

	_ = h.adminDB.WithContext(ctx).Table("qixi_crm_a_financial").
		Where("is_del = 0 AND type = ? AND financial_status = 0 AND status = 1", typeExtract).
		Select("COALESCE(SUM(extract_money),0)").Scan(&out.ApplyingAmount).Error

	_ = h.adminDB.WithContext(ctx).Table("qixi_crm_a_financial").
		Where("is_del = 0 AND type = ? AND financial_status = 0 AND status = 0", typeExtract).
		Select("COALESCE(SUM(extract_money),0)").Scan(&out.PendingAmount).Error

	_ = h.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_view").
		Select("COALESCE(SUM(freeze_money),0)").Scan(&out.FreezeAmount).Error

	response.OK(c, out)
}

func (h *Handler) List(c *gin.Context) {
	page, limit := queryfilter.Page(c)
	f := parseFilter(c)
	q := h.baseQuery(c, f)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "转账记录查询失败")
		return
	}
	rows := make([]transferRow, 0)
	err := q.Select(`
f.financial_id,f.financial_sn,f.mer_id,
COALESCE(m.merchant_name,'') AS mer_name,
COALESCE(m.is_trader,0) AS is_trader,
f.mer_money,f.extract_money,f.financial_type,f.financial_account,
f.financial_status,f.status,f.refusal,f.image,f.admin_id,
COALESCE(u.display_name,'') AS admin_name,
f.mark,f.admin_mark,f.create_time,f.status_time,f.update_time`).
		Order("f.financial_id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&rows).Error
	if err != nil {
		fail(c, "转账记录查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Detail(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	row, err := h.loadRow(c, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "转账记录不存在")
			return
		}
		fail(c, "转账记录详情加载失败")
		return
	}
	response.OK(c, row)
}

type statusInput struct {
	Status  int    `json:"status"`
	Refusal string `json:"refusal"`
}

func (h *Handler) SwitchStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in statusInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if in.Status != 1 && in.Status != -1 {
		response.Fail(c, http.StatusBadRequest, "审核状态错误")
		return
	}
	refusal := strings.TrimSpace(in.Refusal)
	if in.Status == -1 && refusal == "" {
		response.Fail(c, http.StatusBadRequest, "请输入拒绝理由")
		return
	}
	if in.Status == 1 {
		refusal = ""
	}
	now := time.Now()
	adminID := middleware.AdminID(c)
	res := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_financial").
		Where("financial_id = ? AND is_del = 0 AND type = ? AND status = 0", id, typeExtract).
		Updates(map[string]any{
			"status":      in.Status,
			"refusal":     refusal,
			"admin_id":    adminID,
			"status_time": now,
		})
	if res.Error != nil {
		fail(c, "审核失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusBadRequest, "当前状态不允许审核")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

type markInput struct {
	AdminMark string `json:"admin_mark"`
}

func (h *Handler) Mark(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in markInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	mark := strings.TrimSpace(in.AdminMark)
	if mark == "" {
		response.Fail(c, http.StatusBadRequest, "请输入备注")
		return
	}
	res := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_financial").
		Where("financial_id = ? AND is_del = 0 AND type = ?", id, typeExtract).
		Update("admin_mark", mark)
	if res.Error != nil {
		fail(c, "备注保存失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "转账记录不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

type payInput struct {
	Image string `json:"image"`
}

func (h *Handler) Pay(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in payInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	image := strings.TrimSpace(in.Image)
	if image == "" {
		response.Fail(c, http.StatusBadRequest, "请上传转账凭证")
		return
	}
	now := time.Now()
	adminID := middleware.AdminID(c)
	res := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_financial").
		Where("financial_id = ? AND is_del = 0 AND type = ? AND status = 1 AND financial_status = 0", id, typeExtract).
		Updates(map[string]any{
			"image":            image,
			"financial_status": 1,
			"admin_id":         adminID,
			"update_time":      now,
		})
	if res.Error != nil {
		fail(c, "登记转账失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusBadRequest, "仅已审核且未到账的申请可登记转账")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Export(c *gin.Context) {
	f := parseFilter(c)
	var body struct {
		DateFrom        string `json:"date_from"`
		DateTo          string `json:"date_to"`
		Status          *int   `json:"status"`
		MerName         string `json:"mer_name"`
		IsTrader        *int   `json:"is_trader"`
		FinancialType   *int   `json:"financial_type"`
		FinancialStatus *int   `json:"financial_status"`
		AdminKeyword    string `json:"admin_keyword"`
	}
	if err := c.ShouldBindJSON(&body); err == nil {
		if body.DateFrom != "" {
			f.DateFrom = strings.TrimSpace(body.DateFrom)
		}
		if body.DateTo != "" {
			f.DateTo = strings.TrimSpace(body.DateTo)
		}
		if body.MerName != "" {
			f.MerName = strings.TrimSpace(body.MerName)
		}
		if body.AdminKeyword != "" {
			f.AdminKeyword = strings.TrimSpace(body.AdminKeyword)
		}
		if body.Status != nil {
			f.Status = body.Status
		}
		if body.IsTrader != nil {
			f.IsTrader = body.IsTrader
		}
		if body.FinancialType != nil {
			f.FinancialType = body.FinancialType
		}
		if body.FinancialStatus != nil {
			f.FinancialStatus = body.FinancialStatus
		}
	}
	q := h.baseQuery(c, f)
	rows := make([]transferRow, 0)
	err := q.Select(`
f.financial_id,f.financial_sn,f.mer_id,
COALESCE(m.merchant_name,'') AS mer_name,
COALESCE(m.is_trader,0) AS is_trader,
f.mer_money,f.extract_money,f.financial_type,f.financial_account,
f.financial_status,f.status,f.refusal,f.image,f.admin_id,
COALESCE(u.display_name,'') AS admin_name,
f.mark,f.admin_mark,f.create_time,f.status_time,f.update_time`).
		Order("f.financial_id DESC").
		Limit(exportLimit + 1).
		Scan(&rows).Error
	if err != nil {
		fail(c, "转账记录导出失败")
		return
	}
	truncated := len(rows) > exportLimit
	if truncated {
		rows = rows[:exportLimit]
	}

	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"商户名称", "申请时间", "转账金额", "到账状态", "审核状态", "拒绝理由", "商户余额", "转账信息", "平台管理员"})
	for _, row := range rows {
		_ = w.Write([]string{
			row.MerName,
			formatTime(row.CreateTime),
			fmt.Sprintf("%.2f", row.ExtractMoney),
			arrivalLabel(row.FinancialStatus),
			auditLabel(row.Status),
			row.Refusal,
			fmt.Sprintf("%.2f", row.MerMoney),
			formatAccount(row.FinancialType, row.FinancialAccount),
			row.AdminName,
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		fail(c, "转账记录导出失败")
		return
	}
	response.OK(c, gin.H{
		"content":   buf.String(),
		"file_name": fmt.Sprintf("转账记录_%s.csv", time.Now().Format("20060102150405")),
		"row_count": len(rows),
		"truncated": truncated,
	})
}

func (h *Handler) baseQuery(c *gin.Context, f listFilter) *gorm.DB {
	q := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_financial AS f").
		Joins("LEFT JOIN qixi_crm_a_merchant_view AS m ON m.merchant_id = f.mer_id").
		Joins("LEFT JOIN qixi_crm_a_admin_user AS u ON u.id = f.admin_id").
		Where("f.is_del = 0 AND f.type = ?", typeExtract)
	q = applyCreateTimeRange(q, f.DateFrom, f.DateTo, "f.create_time")
	if f.Status != nil {
		q = q.Where("f.status = ?", *f.Status)
	}
	if name := strings.TrimSpace(f.MerName); name != "" {
		like := "%" + name + "%"
		q = q.Where("m.merchant_name LIKE ?", like)
	}
	if f.IsTrader != nil {
		q = q.Where("m.is_trader = ?", *f.IsTrader)
	}
	if f.FinancialType != nil {
		q = q.Where("f.financial_type = ?", *f.FinancialType)
	}
	if f.FinancialStatus != nil {
		q = q.Where("f.financial_status = ?", *f.FinancialStatus)
	}
	if kw := strings.TrimSpace(f.AdminKeyword); kw != "" {
		like := "%" + kw + "%"
		if id, err := strconv.ParseUint(kw, 10, 64); err == nil {
			q = q.Where("(f.admin_id = ? OR u.display_name LIKE ? OR u.username LIKE ?)", id, like, like)
		} else {
			q = q.Where("(u.display_name LIKE ? OR u.username LIKE ?)", like, like)
		}
	}
	return q
}

func (h *Handler) loadRow(c *gin.Context, id uint64) (*transferRow, error) {
	var row transferRow
	err := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_financial AS f").
		Joins("LEFT JOIN qixi_crm_a_merchant_view AS m ON m.merchant_id = f.mer_id").
		Joins("LEFT JOIN qixi_crm_a_admin_user AS u ON u.id = f.admin_id").
		Select(`
f.financial_id,f.financial_sn,f.mer_id,
COALESCE(m.merchant_name,'') AS mer_name,
COALESCE(m.is_trader,0) AS is_trader,
f.mer_money,f.extract_money,f.financial_type,f.financial_account,
f.financial_status,f.status,f.refusal,f.image,f.admin_id,
COALESCE(u.display_name,'') AS admin_name,
f.mark,f.admin_mark,f.create_time,f.status_time,f.update_time`).
		Where("f.financial_id = ? AND f.is_del = 0 AND f.type = ?", id, typeExtract).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func parseFilter(c *gin.Context) listFilter {
	f := listFilter{
		DateFrom:     strings.TrimSpace(c.Query("date_from")),
		DateTo:       strings.TrimSpace(c.Query("date_to")),
		MerName:      strings.TrimSpace(c.Query("mer_name")),
		AdminKeyword: strings.TrimSpace(c.Query("admin_keyword")),
	}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && (v == 0 || v == 1 || v == -1) {
			f.Status = &v
		}
	}
	if raw := strings.TrimSpace(c.Query("is_trader")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && (v == 0 || v == 1) {
			f.IsTrader = &v
		}
	}
	if raw := strings.TrimSpace(c.Query("financial_type")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v >= 1 && v <= 3 {
			f.FinancialType = &v
		}
	}
	if raw := strings.TrimSpace(c.Query("financial_status")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && (v == 0 || v == 1) {
			f.FinancialStatus = &v
		}
	}
	return f
}

func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "ID 错误")
		return 0, false
	}
	return id, true
}

func applyCreateTimeRange(q *gorm.DB, from, to, column string) *gorm.DB {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	if from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where(column+" >= ?", t)
		} else if t, err := time.ParseInLocation("2006-01-02 15:04:05", from, time.Local); err == nil {
			q = q.Where(column+" >= ?", t)
		}
	}
	if to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where(column+" < ?", t.Add(24*time.Hour))
		} else if t, err := time.ParseInLocation("2006-01-02 15:04:05", to, time.Local); err == nil {
			q = q.Where(column+" <= ?", t)
		}
	}
	return q
}

func formatAccount(financialType int, raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	var acc accountInfo
	if err := json.Unmarshal([]byte(raw), &acc); err != nil {
		return raw
	}
	switch financialType {
	case 1:
		return fmt.Sprintf("姓名：%s；银行：%s；卡号：%s", acc.Name, acc.Bank, acc.BankCode)
	case 2:
		return fmt.Sprintf("姓名：%s；微信号：%s", acc.Name, acc.Wechat)
	case 3:
		return fmt.Sprintf("姓名：%s；支付宝：%s", acc.Name, acc.Alipay)
	default:
		return raw
	}
}

func auditLabel(status int) string {
	switch status {
	case 0:
		return "待审核"
	case 1:
		return "审核通过"
	case -1:
		return "审核未通过"
	default:
		return "未知"
	}
}

func arrivalLabel(status int) string {
	if status == 1 {
		return "已到账"
	}
	return "未到账"
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func round2(v float64) float64 {
	x := fmt.Sprintf("%.2f", v)
	n, _ := strconv.ParseFloat(x, 64)
	return n
}

func fail(c *gin.Context, msg string) {
	response.Fail(c, http.StatusInternalServerError, msg)
}
