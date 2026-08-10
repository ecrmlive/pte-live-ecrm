// Package statement exposes platform bill title, day/month list, detail and CSV download.
// Data lives in qixi_crm_b_financial_record (+ qixi_crm_b_user_bill for recharge stats).
package statement

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
	menuStatementRead     = "accounts.statement.read"
	menuStatementDownload = "accounts.statement.download"
	payTypeOffline        = 7
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
	read := middleware.RequireAdminMenu(h.adminDB, menuStatementRead)
	download := middleware.RequireAdminMenu(h.adminDB, menuStatementDownload)
	r.GET("/finance/statements/title", access, read, h.Title)
	r.GET("/finance/statements", access, read, h.List)
	r.GET("/finance/statements/detail", access, read, h.Detail)
	r.POST("/finance/statements/export", access, download, h.Export)
}

type listFilter struct {
	Type     int // 1 day, 2 month
	DateFrom string
	DateTo   string
}

type titleStat struct {
	OrderIncome      float64 `json:"order_income"`
	RefundExpense    float64 `json:"refund_expense"`
	BrokerageExpense float64 `json:"brokerage_expense"`
	PlatformCharge   float64 `json:"platform_charge"`
	RechargeAmount   float64 `json:"recharge_amount"`
	RechargeConsume  float64 `json:"recharge_consume"`
	MerchantCount    int64   `json:"merchant_count"`
	CouponAmount     float64 `json:"coupon_amount"`
}

type billRow struct {
	Date    string  `json:"date"`
	Income  float64 `json:"income"`
	Expend  float64 `json:"expend"`
	Offline float64 `json:"offline"`
	Charge  float64 `json:"charge"`
}

type detailLine struct {
	Label  string `json:"label"`
	Amount string `json:"amount"`
	Count  string `json:"count"`
}

type detailBlock struct {
	Title  string       `json:"title"`
	Number float64      `json:"number"`
	Count  string       `json:"count"`
	Data   []detailLine `json:"data"`
}

type detailResp struct {
	Date   string      `json:"date"`
	Income detailBlock `json:"income"`
	Bill   detailBlock `json:"bill"`
	Expend detailBlock `json:"expend"`
	Charge detailBlock `json:"charge"`
}

func (h *Handler) Title(c *gin.Context) {
	f := parseFilter(c)
	ctx := c.Request.Context()
	out := titleStat{}

	out.OrderIncome = h.sumTypes(c, f, nil, "order", "order_presell", "presell")
	out.RefundExpense = h.sumTypes(c, f, nil, "refund_order")

	brokerage := h.sumTypes(c, f, nil, "brokerage_one", "brokerage_two")
	refundBrokerage := h.sumTypes(c, f, nil, "refund_brokerage_one", "refund_brokerage_two")
	out.BrokerageExpense = round2(brokerage - refundBrokerage)

	charge := h.sumTypes(c, f, nil, "order_charge", "presell_charge")
	refundCharge := h.sumTypes(c, f, nil, "refund_charge")
	out.PlatformCharge = round2(charge - refundCharge)

	coupon := h.sumTypes(c, f, nil, "order_platform_coupon", "order_svip_coupon")
	refundCoupon := h.sumTypes(c, f, nil, "refund_platform_coupon", "refund_svip_coupon")
	out.CouponAmount = round2(coupon - refundCoupon)

	qBill := h.businessDB.WithContext(ctx).Table("qixi_crm_b_user_bill").
		Where("status = ? AND category = ? AND type IN ?", 1, "now_money", []string{"sys_inc_money", "recharge"})
	qBill = applyCreateTimeRange(qBill, f.DateFrom, f.DateTo, "create_time")
	_ = qBill.Select("COALESCE(SUM(number),0)").Scan(&out.RechargeAmount).Error

	qConsume := h.businessDB.WithContext(ctx).Table("qixi_crm_b_user_bill").
		Where("pm = ? AND status = ? AND category = ? AND type IN ?", 0, 1, "now_money", []string{"presell", "pay_product", "sys_dec_money"})
	qConsume = applyCreateTimeRange(qConsume, f.DateFrom, f.DateTo, "create_time")
	_ = qConsume.Select("COALESCE(SUM(number),0)").Scan(&out.RechargeConsume).Error

	qMer := h.financialBase(c, f)
	_ = qMer.Where("mer_id > 0").Distinct("mer_id").Count(&out.MerchantCount).Error

	response.OK(c, out)
}

func (h *Handler) List(c *gin.Context) {
	page, limit := queryfilter.Page(c)
	f := parseFilter(c)
	if f.Type != 2 {
		f.Type = 1
	}

	format := "%Y-%m-%d"
	if f.Type == 2 {
		format = "%Y-%m"
	}

	type periodRow struct {
		Time string `gorm:"column:time"`
	}
	base := h.financialBase(c, f)
	countSQL := fmt.Sprintf("COUNT(DISTINCT DATE_FORMAT(create_time,'%s'))", format)
	var total int64
	if err := base.Select(countSQL).Scan(&total).Error; err != nil {
		fail(c, "平台账单查询失败")
		return
	}

	periods := make([]periodRow, 0)
	if err := h.financialBase(c, f).
		Select(fmt.Sprintf("DATE_FORMAT(create_time,'%s') AS time", format)).
		Group("time").
		Order("time DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&periods).Error; err != nil {
		fail(c, "平台账单查询失败")
		return
	}

	list := make([]billRow, 0, len(periods))
	for _, p := range periods {
		income := h.periodIncome(c, f.Type, p.Time)
		expend := h.periodExpend(c, f.Type, p.Time)
		offline := h.periodOffline(c, f.Type, p.Time)
		list = append(list, billRow{
			Date:    p.Time,
			Income:  income,
			Expend:  expend,
			Offline: offline,
			Charge:  round2(income - expend),
		})
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Detail(c *gin.Context) {
	f := parseFilter(c)
	if f.Type != 2 {
		f.Type = 1
	}
	date := strings.TrimSpace(c.Query("date"))
	if date == "" {
		if f.Type == 1 {
			date = time.Now().Format("2006-01-02")
		} else {
			date = time.Now().Format("2006-01")
		}
	}
	date = normalizePeriod(f.Type, date)

	incomeNum, incomeOrder, incomeOrderCnt, incomeCoupon, incomeCouponCnt, incomeSvip, incomeSvipCnt := h.periodIncomeParts(c, f.Type, date)
	billNum, billCnt := h.periodRecharge(c, f.Type, date)
	expendNum, expendParts := h.periodExpendParts(c, f.Type, date)
	chargeNum := round2(incomeNum - expendNum)

	resp := detailResp{
		Date: date,
		Income: detailBlock{
			Title:  "订单收入总金额",
			Number: incomeNum,
			Count:  fmt.Sprintf("%d笔", incomeOrderCnt),
			Data: []detailLine{
				{Label: "订单支付", Amount: moneyYuan(incomeOrder), Count: fmt.Sprintf("%d笔", incomeOrderCnt)},
				{Label: "退回优惠券补贴", Amount: moneyYuan(incomeCoupon), Count: fmt.Sprintf("%d笔", incomeCouponCnt)},
				{Label: "退回会员优惠券补贴", Amount: moneyYuan(incomeSvip), Count: fmt.Sprintf("%d笔", incomeSvipCnt)},
			},
		},
		Bill: detailBlock{
			Title:  "充值金额",
			Number: billNum,
			Count:  fmt.Sprintf("%d笔", billCnt),
			Data:   []detailLine{},
		},
		Expend: detailBlock{
			Title:  "支出总金额",
			Number: expendNum,
			Count:  fmt.Sprintf("%d笔", expendParts.totalCount),
			Data: []detailLine{
				{Label: "应付商户金额", Amount: moneyYuan(expendParts.order), Count: fmt.Sprintf("%d笔", expendParts.orderCnt)},
				{Label: "商户线下已收", Amount: moneyYuan(expendParts.offline), Count: fmt.Sprintf("%d笔", expendParts.offlineCnt)},
				{Label: "佣金", Amount: moneyYuan(expendParts.brokerage), Count: fmt.Sprintf("%d笔", expendParts.brokerageCnt)},
				{Label: "返还手续费", Amount: moneyYuan(expendParts.charge), Count: fmt.Sprintf("%d笔", expendParts.chargeCnt)},
				{Label: "优惠券补贴", Amount: moneyYuan(expendParts.coupon), Count: fmt.Sprintf("%d笔", expendParts.couponCnt)},
				{Label: "会员优惠券补贴", Amount: moneyYuan(expendParts.svip), Count: fmt.Sprintf("%d笔", expendParts.svipCnt)},
			},
		},
		Charge: detailBlock{
			Title:  "平台应入账金额",
			Number: chargeNum,
			Count:  "",
			Data:   []detailLine{},
		},
	}
	response.OK(c, resp)
}

func (h *Handler) Export(c *gin.Context) {
	var in struct {
		Type int    `json:"type"`
		Date string `json:"date"`
	}
	_ = c.ShouldBindJSON(&in)
	if in.Type != 2 {
		in.Type = 1
	}
	date := strings.TrimSpace(in.Date)
	if date == "" {
		date = strings.TrimSpace(c.Query("date"))
	}
	if date == "" {
		fail(c, "请指定账单日期")
		return
	}
	date = normalizePeriod(in.Type, date)

	incomeNum, incomeOrder, incomeOrderCnt, incomeCoupon, incomeCouponCnt, incomeSvip, incomeSvipCnt := h.periodIncomeParts(c, in.Type, date)
	billNum, billCnt := h.periodRecharge(c, in.Type, date)
	expendNum, expendParts := h.periodExpendParts(c, in.Type, date)
	chargeNum := round2(incomeNum - expendNum)

	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"账期", date})
	_ = w.Write([]string{"项目", "金额(元)", "笔数"})
	_ = w.Write([]string{"订单收入总金额", formatPlain(incomeNum), strconv.FormatInt(incomeOrderCnt, 10)})
	_ = w.Write([]string{"  订单支付", formatPlain(incomeOrder), strconv.FormatInt(incomeOrderCnt, 10)})
	_ = w.Write([]string{"  退回优惠券补贴", formatPlain(incomeCoupon), strconv.FormatInt(incomeCouponCnt, 10)})
	_ = w.Write([]string{"  退回会员优惠券补贴", formatPlain(incomeSvip), strconv.FormatInt(incomeSvipCnt, 10)})
	_ = w.Write([]string{"充值金额", formatPlain(billNum), strconv.FormatInt(billCnt, 10)})
	_ = w.Write([]string{"支出总金额", formatPlain(expendNum), strconv.FormatInt(expendParts.totalCount, 10)})
	_ = w.Write([]string{"  应付商户金额", formatPlain(expendParts.order), strconv.FormatInt(expendParts.orderCnt, 10)})
	_ = w.Write([]string{"  商户线下已收", formatPlain(expendParts.offline), strconv.FormatInt(expendParts.offlineCnt, 10)})
	_ = w.Write([]string{"  佣金", formatPlain(expendParts.brokerage), strconv.FormatInt(expendParts.brokerageCnt, 10)})
	_ = w.Write([]string{"  返还手续费", formatPlain(expendParts.charge), strconv.FormatInt(expendParts.chargeCnt, 10)})
	_ = w.Write([]string{"  优惠券补贴", formatPlain(expendParts.coupon), strconv.FormatInt(expendParts.couponCnt, 10)})
	_ = w.Write([]string{"  会员优惠券补贴", formatPlain(expendParts.svip), strconv.FormatInt(expendParts.svipCnt, 10)})
	_ = w.Write([]string{"平台应入账金额", formatPlain(chargeNum), ""})
	w.Flush()
	if w.Error() != nil {
		fail(c, "平台账单导出失败")
		return
	}
	label := "日账单"
	if in.Type == 2 {
		label = "月账单"
	}
	response.OK(c, gin.H{
		"file_name": fmt.Sprintf("平台账单_%s_%s.csv", label, strings.ReplaceAll(date, "-", "")),
		"content":   buf.String(),
		"row_count": 14,
		"truncated": false,
	})
}

type expendParts struct {
	order, offline, brokerage, charge, coupon, svip                               float64
	orderCnt, offlineCnt, brokerageCnt, chargeCnt, couponCnt, svipCnt, totalCount int64
}

func (h *Handler) periodIncome(c *gin.Context, typ int, period string) float64 {
	n, _, _, _, _, _, _ := h.periodIncomeParts(c, typ, period)
	return n
}

func (h *Handler) periodIncomeParts(c *gin.Context, typ int, period string) (
	total, order float64, orderCnt int64,
	coupon float64, couponCnt int64,
	svip float64, svipCnt int64,
) {
	order, orderCnt = h.sumCountTypesPeriod(c, typ, period, "order", "order_presell", "presell")
	coupon, couponCnt = h.sumCountTypesPeriod(c, typ, period, "refund_platform_coupon")
	svip, svipCnt = h.sumCountTypesPeriod(c, typ, period, "refund_svip_coupon")
	total = round2(order + coupon + svip)
	return
}

func (h *Handler) periodExpend(c *gin.Context, typ int, period string) float64 {
	n, _ := h.periodExpendParts(c, typ, period)
	return n
}

func (h *Handler) periodExpendParts(c *gin.Context, typ int, period string) (float64, expendParts) {
	var p expendParts
	p.brokerage, p.brokerageCnt = h.sumCountTypesPeriod(c, typ, period, "brokerage_one", "brokerage_two")
	p.charge, p.chargeCnt = h.sumCountTypesPeriod(c, typ, period, "refund_charge")
	p.order, p.orderCnt = h.sumCountTypesPeriod(c, typ, period, "order_true", "presell_true")
	p.coupon, p.couponCnt = h.sumCountTypesPeriod(c, typ, period, "order_platform_coupon")
	p.svip, p.svipCnt = h.sumCountTypesPeriod(c, typ, period, "order_svip_coupon")
	p.offline, p.offlineCnt = h.sumCountOfflineRefund(c, typ, period)

	p.order = round2(p.order - p.offline)
	if p.orderCnt > p.offlineCnt {
		p.orderCnt -= p.offlineCnt
	} else {
		p.orderCnt = 0
	}
	number := round2(p.brokerage + p.order + p.coupon + p.svip + p.charge + p.offline)
	p.totalCount = p.brokerageCnt + p.orderCnt + p.chargeCnt + p.offlineCnt + p.couponCnt + p.svipCnt
	return number, p
}

func (h *Handler) periodOffline(c *gin.Context, typ int, period string) float64 {
	// 列表「店铺线下已收」：线下支付订单收入金额
	n, _ := h.sumCountTypesPeriodPay(c, typ, period, payTypeOffline, "order")
	return n
}

func (h *Handler) periodRecharge(c *gin.Context, typ int, period string) (float64, int64) {
	base := func() *gorm.DB {
		q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_bill").
			Where("pm = ? AND status = ? AND category = ? AND type IN ?", 1, 1, "now_money", []string{"sys_inc_money", "recharge"})
		return applyPeriod(q, typ, period, "create_time")
	}
	var number float64
	var count int64
	_ = base().Count(&count).Error
	_ = base().Select("COALESCE(SUM(number),0)").Scan(&number).Error
	return number, count
}

func (h *Handler) financialBase(c *gin.Context, f listFilter) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_financial_record")
	return applyCreateTimeRange(q, f.DateFrom, f.DateTo, "create_time")
}

func (h *Handler) sumTypes(c *gin.Context, f listFilter, extra func(*gorm.DB) *gorm.DB, types ...string) float64 {
	q := h.financialBase(c, f).Where("financial_type IN ?", types)
	if extra != nil {
		q = extra(q)
	}
	var n float64
	_ = q.Select("COALESCE(SUM(number),0)").Scan(&n).Error
	return n
}

func (h *Handler) sumCountTypesPeriod(c *gin.Context, typ int, period string, types ...string) (float64, int64) {
	base := func() *gorm.DB {
		q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_financial_record").
			Where("financial_type IN ?", types)
		return applyPeriod(q, typ, period, "create_time")
	}
	var number float64
	var count int64
	_ = base().Count(&count).Error
	_ = base().Select("COALESCE(SUM(number),0)").Scan(&number).Error
	return number, count
}

func (h *Handler) sumCountTypesPeriodPay(c *gin.Context, typ int, period string, payType int, types ...string) (float64, int64) {
	base := func() *gorm.DB {
		q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_financial_record").
			Where("financial_type IN ? AND pay_type = ?", types, payType)
		return applyPeriod(q, typ, period, "create_time")
	}
	var number float64
	var count int64
	_ = base().Count(&count).Error
	_ = base().Select("COALESCE(SUM(number),0)").Scan(&number).Error
	return number, count
}

func (h *Handler) sumCountOfflineRefund(c *gin.Context, typ int, period string) (float64, int64) {
	// 详情支出「商户线下已收」：线下订单应付商户金额（order_true + 线下支付）
	return h.sumCountTypesPeriodPay(c, typ, period, payTypeOffline, "order_true")
}

func parseFilter(c *gin.Context) listFilter {
	typ, _ := strconv.Atoi(strings.TrimSpace(c.Query("type")))
	return listFilter{
		Type:     typ,
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

func applyPeriod(q *gorm.DB, typ int, period, column string) *gorm.DB {
	period = normalizePeriod(typ, period)
	if typ == 2 {
		if t, err := time.ParseInLocation("2006-01", period, time.Local); err == nil {
			start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
			end := start.AddDate(0, 1, 0)
			return q.Where(column+" >= ? AND "+column+" < ?", start, end)
		}
		return q
	}
	if t, err := time.ParseInLocation("2006-01-02", period, time.Local); err == nil {
		return q.Where(column+" >= ? AND "+column+" < ?", t, t.AddDate(0, 0, 1))
	}
	return q
}

func normalizePeriod(typ int, period string) string {
	period = strings.TrimSpace(period)
	if typ == 2 {
		if len(period) >= 7 {
			return period[:7]
		}
		return period
	}
	if len(period) >= 10 {
		return period[:10]
	}
	return period
}

func round2(v float64) float64 {
	x, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", v), 64)
	return x
}

func moneyYuan(v float64) string {
	return formatPlain(v) + "元"
}

func formatPlain(v float64) string {
	return strconv.FormatFloat(round2(v), 'f', 2, 64)
}

func fail(c *gin.Context, msg string) {
	response.Fail(c, http.StatusInternalServerError, msg)
}
