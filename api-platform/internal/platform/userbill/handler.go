// Package userbill exposes platform user balance bill list (资金记录).
// Aligns with CRMEB UserBill (category now_money / svip_pay).
package userbill

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
	menuRead   = "accounts.user_assets.read"
	menuExport = "accounts.user_assets.export"
	exportLimit = 5000
	tableName  = "qixi_crm_b_user_bill"
)

// TYPE_INFO mirrors CRMEB UserBillRepository::TYPE_INFO for now_money + svip_pay.
var typeOptions = []struct {
	Type  string
	Title string
}{
	{Type: "now_money/brokerage", Title: "佣金转入余额"},
	{Type: "now_money/pay_product", Title: "购买商品"},
	{Type: "now_money/presell", Title: "支付预售尾款"},
	{Type: "now_money/recharge", Title: "余额充值"},
	{Type: "now_money/sys_dec_money", Title: "系统减少余额"},
	{Type: "now_money/sys_inc_money", Title: "系统增加余额"},
	{Type: "svip_pay/svip_pay", Title: "付费会员支付"},
}

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
	export := middleware.RequireAdminMenu(h.adminDB, menuExport)
	r.GET("/finance/user-bills", access, read, h.List)
	r.GET("/finance/user-bills/types", access, read, h.Types)
	r.POST("/finance/user-bills/export", access, export, h.Export)
}

type listFilter struct {
	DateFrom    string
	DateTo      string
	BillType    string
	UserType    string
	UserKeyword string
}

type row struct {
	BillID     uint64     `gorm:"column:bill_id" json:"bill_id"`
	UID        uint64     `gorm:"column:uid" json:"uid"`
	Nickname   string     `gorm:"column:nickname" json:"nickname"`
	PM         int        `gorm:"column:pm" json:"pm"`
	Title      string     `gorm:"column:title" json:"title"`
	Category   string     `gorm:"column:category" json:"category"`
	Type       string     `gorm:"column:type" json:"type"`
	Number     float64    `gorm:"column:number" json:"number"`
	Balance    float64    `gorm:"column:balance" json:"balance"`
	Mark       string     `gorm:"column:mark" json:"mark"`
	CreateTime *time.Time `gorm:"column:create_time" json:"create_time"`
}

func (h *Handler) Types(c *gin.Context) {
	list := make([]gin.H, 0, len(typeOptions))
	for _, opt := range typeOptions {
		list = append(list, gin.H{"type": opt.Type, "title": opt.Title})
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) List(c *gin.Context) {
	page, limit := queryfilter.Page(c)
	f := parseFilter(c)
	var total int64
	if err := h.baseQuery(c, f).Count(&total).Error; err != nil {
		fail(c, "资金记录查询失败")
		return
	}
	rows := make([]row, 0)
	err := h.baseQuery(c, f).Select(`a.bill_id,a.uid,COALESCE(u.nickname,'') AS nickname,a.pm,a.title,
a.category,a.type,a.number,a.balance,a.mark,a.create_time`).
		Order("a.create_time DESC, a.bill_id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		fail(c, "资金记录查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Export(c *gin.Context) {
	f := parseFilter(c)
	var body struct {
		DateFrom    string `json:"date_from"`
		DateTo      string `json:"date_to"`
		BillType    string `json:"type"`
		UserType    string `json:"user_type"`
		UserKeyword string `json:"user_keyword"`
	}
	if err := c.ShouldBindJSON(&body); err == nil {
		if body.DateFrom != "" {
			f.DateFrom = strings.TrimSpace(body.DateFrom)
		}
		if body.DateTo != "" {
			f.DateTo = strings.TrimSpace(body.DateTo)
		}
		if body.BillType != "" {
			f.BillType = strings.TrimSpace(body.BillType)
		}
		if body.UserType != "" {
			f.UserType = strings.TrimSpace(body.UserType)
		}
		if body.UserKeyword != "" {
			f.UserKeyword = strings.TrimSpace(body.UserKeyword)
		}
	}
	rows := make([]row, 0)
	err := h.baseQuery(c, f).Select(`a.bill_id,a.uid,COALESCE(u.nickname,'') AS nickname,a.pm,a.title,
a.category,a.type,a.number,a.balance,a.mark,a.create_time`).
		Order("a.create_time DESC, a.bill_id DESC").
		Limit(exportLimit + 1).
		Scan(&rows).Error
	if err != nil {
		fail(c, "资金记录导出失败")
		return
	}
	truncated := len(rows) > exportLimit
	if truncated {
		rows = rows[:exportLimit]
	}
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"会员ID", "昵称", "金额", "明细类型", "备注", "创建时间"})
	for _, item := range rows {
		_ = w.Write([]string{
			strconv.FormatUint(item.UID, 10),
			item.Nickname,
			fmt.Sprintf("%.2f", item.Number),
			item.Title,
			item.Mark,
			formatTime(item.CreateTime),
		})
	}
	w.Flush()
	response.OK(c, gin.H{
		"content":   buf.String(),
		"file_name": fmt.Sprintf("资金记录_%s.csv", time.Now().Format("20060102150405")),
		"row_count": len(rows),
		"truncated": truncated,
	})
}

func (h *Handler) baseQuery(c *gin.Context, f listFilter) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).Table(tableName+" AS a").
		Joins("LEFT JOIN qixi_crm_b_user u ON u.id = a.uid").
		Where("a.category IN ?", []string{"now_money", "svip_pay"}).
		Where("a.category <> ?", "sys_brokerage")
	if f.DateFrom != "" {
		q = q.Where("a.create_time >= ?", f.DateFrom+" 00:00:00")
	}
	if f.DateTo != "" {
		q = q.Where("a.create_time <= ?", f.DateTo+" 23:59:59")
	}
	if t := strings.TrimSpace(f.BillType); t != "" {
		parts := strings.SplitN(t, "/", 2)
		if len(parts) == 2 {
			q = q.Where("a.category = ? AND a.type = ?", parts[0], parts[1])
		} else {
			q = q.Where("a.type = ?", t)
		}
	}
	if kw := strings.TrimSpace(f.UserKeyword); kw != "" {
		switch f.UserType {
		case "uid":
			if uid, err := strconv.ParseUint(kw, 10, 64); err == nil && uid > 0 {
				q = q.Where("a.uid = ?", uid)
			} else {
				q = q.Where("1 = 0")
			}
		case "phone":
			q = q.Where("u.mobile LIKE ?", "%"+kw+"%")
		default:
			q = q.Where("u.nickname LIKE ?", "%"+kw+"%")
		}
	}
	return q
}

func parseFilter(c *gin.Context) listFilter {
	f := listFilter{
		DateFrom:    strings.TrimSpace(c.Query("date_from")),
		DateTo:      strings.TrimSpace(c.Query("date_to")),
		BillType:    strings.TrimSpace(c.Query("type")),
		UserType:    strings.TrimSpace(c.Query("user_type")),
		UserKeyword: strings.TrimSpace(c.Query("user_keyword")),
	}
	if f.UserType == "" {
		f.UserType = "nickname"
	}
	return f
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func fail(c *gin.Context, msg string) {
	response.Fail(c, http.StatusInternalServerError, msg)
}
