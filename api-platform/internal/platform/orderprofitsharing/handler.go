// Package orderprofitsharing exposes platform order profit-sharing list,
// detail, retry and CSV export (对齐 CRMEB OrderProfitsharing).
package orderprofitsharing

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
	menuRead   = "accounts.profitsharing.read"
	menuManage = "accounts.profitsharing.manage"
	menuExport = "accounts.profitsharing.export"
	exportLimit = 5000
	tableName   = "qixi_crm_b_store_order_profitsharing"
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
	manage := middleware.RequireAdminMenu(h.adminDB, menuManage)
	export := middleware.RequireAdminMenu(h.adminDB, menuExport)
	r.GET("/finance/order-profitsharings", access, read, h.List)
	r.GET("/finance/order-profitsharings/:id", access, read, h.Detail)
	r.POST("/finance/order-profitsharings/:id/again", access, manage, h.Again)
	r.POST("/finance/order-profitsharings/export", access, export, h.Export)
}

type listFilter struct {
	DateFrom       string
	DateTo         string
	ProfitDateFrom string
	ProfitDateTo   string
	Status         *int
	Type           string
	MerID          *uint64
	MerName        string
}

type row struct {
	ProfitsharingID      uint64     `gorm:"column:profitsharing_id" json:"profitsharing_id"`
	ProfitsharingSN      string     `gorm:"column:profitsharing_sn" json:"profitsharing_sn"`
	OrderID              uint64     `gorm:"column:order_id" json:"order_id"`
	OrderSN              string     `gorm:"column:order_sn" json:"order_sn"`
	MerID                uint64     `gorm:"column:mer_id" json:"mer_id"`
	MerName              string     `gorm:"column:mer_name" json:"mer_name"`
	TransactionID        string     `gorm:"column:transaction_id" json:"transaction_id"`
	ProfitsharingPrice   float64    `gorm:"column:profitsharing_price" json:"profitsharing_price"`
	ProfitsharingRefund  float64    `gorm:"column:profitsharing_refund" json:"profitsharing_refund"`
	ProfitsharingMerPrice float64   `gorm:"column:profitsharing_mer_price" json:"profitsharing_mer_price"`
	Type                 string     `gorm:"column:type" json:"type"`
	TypeName             string     `gorm:"-" json:"type_name"`
	Status               int        `gorm:"column:status" json:"status"`
	StatusName           string     `gorm:"-" json:"status_name"`
	ErrorMsg             string     `gorm:"column:error_msg" json:"error_msg"`
	ProfitsharingTime    *time.Time `gorm:"column:profitsharing_time" json:"profitsharing_time"`
	CreateTime           *time.Time `gorm:"column:create_time" json:"create_time"`
	IsCombine            int        `gorm:"column:is_combine" json:"is_combine"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := queryfilter.Page(c)
	f := parseFilter(c)
	q := h.baseQuery(c, f)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "分账列表查询失败")
		return
	}
	rows := make([]row, 0)
	err := q.Order("profitsharing_id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		fail(c, "分账列表查询失败")
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
	err := h.businessDB.WithContext(c.Request.Context()).Table(tableName).
		Where("profitsharing_id = ?", id).
		Take(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, "分账单不存在")
			return
		}
		fail(c, "分账详情查询失败")
		return
	}
	decorate(&item)
	response.OK(c, item)
}

func (h *Handler) Again(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var item row
	err := h.businessDB.WithContext(c.Request.Context()).Table(tableName).
		Where("profitsharing_id = ?", id).
		Take(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fail(c, "分账单不存在")
			return
		}
		fail(c, "分账单查询失败")
		return
	}
	if item.Status != -2 {
		fail(c, "分账单状态有误，不能重新分账")
		return
	}
	now := time.Now()
	err = h.businessDB.WithContext(c.Request.Context()).Table(tableName).
		Where("profitsharing_id = ? AND status = ?", id, -2).
		Updates(map[string]any{
			"status":             1,
			"error_msg":          "",
			"profitsharing_time": now,
		}).Error
	if err != nil {
		fail(c, "重新分账失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Export(c *gin.Context) {
	f := parseFilter(c)
	var body struct {
		DateFrom       string `json:"date_from"`
		DateTo         string `json:"date_to"`
		ProfitDateFrom string `json:"profit_date_from"`
		ProfitDateTo   string `json:"profit_date_to"`
		Status         *int   `json:"status"`
		Type           string `json:"type"`
		MerID          *uint64 `json:"mer_id"`
		MerName        string `json:"mer_name"`
	}
	if err := c.ShouldBindJSON(&body); err == nil {
		if body.DateFrom != "" {
			f.DateFrom = strings.TrimSpace(body.DateFrom)
		}
		if body.DateTo != "" {
			f.DateTo = strings.TrimSpace(body.DateTo)
		}
		if body.ProfitDateFrom != "" {
			f.ProfitDateFrom = strings.TrimSpace(body.ProfitDateFrom)
		}
		if body.ProfitDateTo != "" {
			f.ProfitDateTo = strings.TrimSpace(body.ProfitDateTo)
		}
		if body.Type != "" {
			f.Type = strings.TrimSpace(body.Type)
		}
		if body.MerName != "" {
			f.MerName = strings.TrimSpace(body.MerName)
		}
		if body.Status != nil {
			f.Status = body.Status
		}
		if body.MerID != nil {
			f.MerID = body.MerID
		}
	}
	q := h.baseQuery(c, f)
	rows := make([]row, 0)
	err := q.Order("profitsharing_id DESC").Limit(exportLimit + 1).Scan(&rows).Error
	if err != nil {
		fail(c, "分账导出失败")
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
	_ = w.Write([]string{"分账ID", "分账单号", "订单编号", "店铺名称", "订单金额", "账单类型", "状态", "创建时间", "分账时间", "失败原因"})
	for _, item := range rows {
		_ = w.Write([]string{
			strconv.FormatUint(item.ProfitsharingID, 10),
			item.ProfitsharingSN,
			item.OrderSN,
			item.MerName,
			fmt.Sprintf("%.2f", item.ProfitsharingPrice),
			item.TypeName,
			item.StatusName,
			formatTime(item.CreateTime),
			formatTime(item.ProfitsharingTime),
			item.ErrorMsg,
		})
	}
	w.Flush()
	response.OK(c, gin.H{
		"content":   buf.String(),
		"file_name": fmt.Sprintf("分账管理_%s.csv", time.Now().Format("20060102150405")),
		"row_count": len(rows),
		"truncated": truncated,
	})
}

func (h *Handler) baseQuery(c *gin.Context, f listFilter) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).Table(tableName)
	q = applyDateRange(q, f.DateFrom, f.DateTo, "create_time")
	q = applyDateRange(q, f.ProfitDateFrom, f.ProfitDateTo, "profitsharing_time")
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.MerID != nil && *f.MerID > 0 {
		q = q.Where("mer_id = ?", *f.MerID)
	}
	if name := strings.TrimSpace(f.MerName); name != "" {
		q = q.Where("mer_name LIKE ?", "%"+name+"%")
	}
	return q
}

func parseFilter(c *gin.Context) listFilter {
	f := listFilter{
		DateFrom:       strings.TrimSpace(c.Query("date_from")),
		DateTo:         strings.TrimSpace(c.Query("date_to")),
		ProfitDateFrom: strings.TrimSpace(c.Query("profit_date_from")),
		ProfitDateTo:   strings.TrimSpace(c.Query("profit_date_to")),
		Type:           strings.TrimSpace(c.Query("type")),
		MerName:        strings.TrimSpace(c.Query("mer_name")),
	}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			f.Status = &v
		}
	}
	if raw := strings.TrimSpace(c.Query("mer_id")); raw != "" {
		if v, err := strconv.ParseUint(raw, 10, 64); err == nil {
			f.MerID = &v
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
		fail(c, "分账 ID 无效")
		return 0, false
	}
	return id, true
}

func decorate(item *row) {
	item.TypeName = typeName(item.Type)
	item.StatusName = statusName(item.Status)
}

func typeName(t string) string {
	switch t {
	case "presell":
		return "尾款支付"
	default:
		return "订单支付"
	}
}

func statusName(status int) string {
	switch status {
	case 2:
		return "分账中"
	case 1:
		return "已分账"
	case 0:
		return "待分账"
	case -1:
		return "已退款"
	case -2:
		return "分账失败"
	default:
		return "未知"
	}
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
