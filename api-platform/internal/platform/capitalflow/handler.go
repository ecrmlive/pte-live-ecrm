// Package capitalflow exposes platform capital flow list (资金流水).
// Aligns with CRMEB FinancialRecord::lst for platform console.
package capitalflow

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
	menuRead    = "accounts.capital_flow.read"
	menuExport  = "accounts.capital_flow.export"
	exportLimit = 5000
	tableName   = "qixi_crm_b_financial_record"
)

// Platform list financial types (CRMEB FinancialRecord::lst).
var platformTypes = []string{
	"order", "sys_accoubts", "brokerage_one", "brokerage_two",
	"refund_brokerage_one", "refund_brokerage_two", "refund_order",
	"order_platform_coupon", "order_svip_coupon", "svip",
}

// payTypeFilter mirrors StoreOrderRepository::PAY_TYPE_FILTEER.
var payTypeFilter = map[string][]int{
	"0": {0},
	"1": {1, 2, 3, 6, 8},
	"2": {4, 5, 9},
	"3": {7},
}

var typeNames = map[string]string{
	"sys_accoubts":          "财务对账",
	"refund_order":          "退款订单",
	"brokerage_one":         "一级分佣",
	"brokerage_two":         "二级分佣",
	"refund_brokerage_one":  "返还一级分佣",
	"refund_brokerage_two":  "返还二级分佣",
	"order":                 "订单支付",
	"svip":                  "支付会员费",
	"order_platform_coupon": "平台优惠券",
	"order_svip_coupon":     "会员优惠券",
}

var payTypeNames = map[int]string{
	0: "余额支付",
	1: "微信支付",
	2: "小程序",
	3: "微信支付",
	4: "支付宝",
	5: "支付宝扫码",
	6: "微信扫码",
	7: "线下支付",
	8: "微信付款码",
	9: "支付宝付款码",
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
	r.GET("/finance/capital-flows", access, read, h.List)
	r.GET("/finance/capital-flows/:id", access, read, h.Detail)
	r.POST("/finance/capital-flows/export", access, export, h.Export)
}

type listFilter struct {
	DateFrom    string
	DateTo      string
	OrderSN     string
	PayTypeKey  string
	UserType    string
	UserKeyword string
}

type row struct {
	FinancialRecordID uint64     `gorm:"column:financial_record_id" json:"financial_record_id"`
	FinancialRecordSN string     `gorm:"column:financial_record_sn" json:"financial_record_sn"`
	OrderID           uint64     `gorm:"column:order_id" json:"order_id"`
	OrderSN           string     `gorm:"column:order_sn" json:"order_sn"`
	UserInfo          string     `gorm:"column:user_info" json:"user_info"`
	UserID            uint64     `gorm:"column:user_id" json:"user_id"`
	FinancialType     string     `gorm:"column:financial_type" json:"financial_type"`
	FinancialTypeName string     `gorm:"-" json:"financial_type_name"`
	FinancialPM       int        `gorm:"column:financial_pm" json:"financial_pm"`
	Number            float64    `gorm:"column:number" json:"number"`
	SignedNumber      float64    `gorm:"-" json:"signed_number"`
	Type              int        `gorm:"column:type" json:"type"`
	MerID             uint64     `gorm:"column:mer_id" json:"mer_id"`
	PayType           int        `gorm:"column:pay_type" json:"pay_type"`
	PayTypeName       string     `gorm:"-" json:"pay_type_name"`
	TransactionID     string     `gorm:"-" json:"transaction_id"`
	CreateTime        *time.Time `gorm:"column:create_time" json:"create_time"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := queryfilter.Page(c)
	f := parseFilter(c)
	var total int64
	if err := h.baseQuery(c, f).Count(&total).Error; err != nil {
		fail(c, "资金流水查询失败")
		return
	}
	rows := make([]row, 0)
	err := h.baseQuery(c, f).Select(`f.financial_record_id,f.financial_record_sn,f.order_id,f.order_sn,
f.user_info,f.user_id,f.financial_type,f.financial_pm,f.number,f.type,f.mer_id,f.pay_type,f.create_time`).
		Order("f.create_time DESC, f.financial_record_id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		fail(c, "资金流水查询失败")
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
	err := h.businessDB.WithContext(c.Request.Context()).Table(tableName+" AS f").
		Where("f.financial_record_id = ?", id).
		Where("f.type IN ?", []int{1, 2}).
		Select(`f.financial_record_id,f.financial_record_sn,f.order_id,f.order_sn,
f.user_info,f.user_id,f.financial_type,f.financial_pm,f.number,f.type,f.mer_id,f.pay_type,f.create_time`).
		Take(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "资金流水不存在")
			return
		}
		fail(c, "资金流水查询失败")
		return
	}
	decorate(&item)
	response.OK(c, item)
}

func (h *Handler) Export(c *gin.Context) {
	f := parseFilter(c)
	var body struct {
		DateFrom    string `json:"date_from"`
		DateTo      string `json:"date_to"`
		OrderSN     string `json:"order_sn"`
		PayTypeKey  string `json:"pay_type"`
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
		if body.OrderSN != "" {
			f.OrderSN = strings.TrimSpace(body.OrderSN)
		}
		if body.PayTypeKey != "" {
			f.PayTypeKey = strings.TrimSpace(body.PayTypeKey)
		}
		if body.UserType != "" {
			f.UserType = strings.TrimSpace(body.UserType)
		}
		if body.UserKeyword != "" {
			f.UserKeyword = strings.TrimSpace(body.UserKeyword)
		}
	}
	rows := make([]row, 0)
	err := h.baseQuery(c, f).Select(`f.financial_record_id,f.financial_record_sn,f.order_id,f.order_sn,
f.user_info,f.user_id,f.financial_type,f.financial_pm,f.number,f.type,f.mer_id,f.pay_type,f.create_time`).
		Order("f.create_time DESC, f.financial_record_id DESC").
		Limit(exportLimit + 1).
		Scan(&rows).Error
	if err != nil {
		fail(c, "资金流水导出失败")
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
	_ = w.Write([]string{"订单号", "交易流水号", "第三方交易单号", "交易时间", "对方信息", "交易类型", "支付方式", "收支金额"})
	for _, item := range rows {
		_ = w.Write([]string{
			dash(item.OrderSN),
			item.FinancialRecordSN,
			"-",
			formatTime(item.CreateTime),
			dash(item.UserInfo),
			item.FinancialTypeName,
			item.PayTypeName,
			fmt.Sprintf("%.2f", item.SignedNumber),
		})
	}
	w.Flush()
	response.OK(c, gin.H{
		"content":   buf.String(),
		"file_name": fmt.Sprintf("资金流水_%s.csv", time.Now().Format("20060102150405")),
		"row_count": len(rows),
		"truncated": truncated,
	})
}

func (h *Handler) baseQuery(c *gin.Context, f listFilter) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).Table(tableName+" AS f").
		Joins("LEFT JOIN qixi_crm_b_user u ON u.id = f.user_id").
		Where("f.financial_type IN ?", platformTypes).
		Where("f.type IN ?", []int{1, 2})
	if f.DateFrom != "" {
		q = q.Where("f.create_time >= ?", f.DateFrom+" 00:00:00")
	}
	if f.DateTo != "" {
		q = q.Where("f.create_time <= ?", f.DateTo+" 23:59:59")
	}
	if sn := strings.TrimSpace(f.OrderSN); sn != "" {
		q = q.Where("f.order_sn LIKE ?", "%"+sn+"%")
	}
	if ids, ok := payTypeFilter[f.PayTypeKey]; ok {
		q = q.Where("f.pay_type IN ?", ids)
	}
	if kw := strings.TrimSpace(f.UserKeyword); kw != "" {
		switch f.UserType {
		case "uid":
			if uid, err := strconv.ParseUint(kw, 10, 64); err == nil && uid > 0 {
				q = q.Where("f.user_id = ?", uid)
			} else {
				q = q.Where("1 = 0")
			}
		case "phone":
			q = q.Where("u.mobile LIKE ?", "%"+kw+"%")
		default:
			q = q.Where("(u.nickname LIKE ? OR f.user_info LIKE ?)", "%"+kw+"%", "%"+kw+"%")
		}
	}
	return q
}

func parseFilter(c *gin.Context) listFilter {
	f := listFilter{
		DateFrom:    strings.TrimSpace(c.Query("date_from")),
		DateTo:      strings.TrimSpace(c.Query("date_to")),
		OrderSN:     strings.TrimSpace(c.Query("order_sn")),
		PayTypeKey:  strings.TrimSpace(c.Query("pay_type")),
		UserType:    strings.TrimSpace(c.Query("user_type")),
		UserKeyword: strings.TrimSpace(c.Query("user_keyword")),
	}
	if f.UserType == "" {
		f.UserType = "nickname"
	}
	return f
}

func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "资金流水参数错误")
		return 0, false
	}
	return id, true
}

func decorate(r *row) {
	if name, ok := typeNames[r.FinancialType]; ok {
		r.FinancialTypeName = name
	} else if r.FinancialType != "" {
		r.FinancialTypeName = r.FinancialType
	} else {
		r.FinancialTypeName = "-"
	}
	if name, ok := payTypeNames[r.PayType]; ok {
		r.PayTypeName = name
	} else {
		r.PayTypeName = "-"
	}
	if r.FinancialPM == 1 {
		r.SignedNumber = r.Number
	} else {
		r.SignedNumber = -r.Number
	}
	r.TransactionID = "-"
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func dash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}

func fail(c *gin.Context, msg string) {
	response.Fail(c, http.StatusInternalServerError, msg)
}
