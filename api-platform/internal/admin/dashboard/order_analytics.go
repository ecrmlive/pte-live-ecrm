package dashboard

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// OrderTopCard mirrors CRMEB analytics/order/top rows.
// Count/Mom may be order counts or money amounts.
type OrderTopCard struct {
	Title     string  `json:"title"`
	Count     float64 `json:"count"`
	Mom       float64 `json:"mom"`
	Statistic float64 `json:"statistic"`
}

// OrderLinePoint mirrors CRMEB analytics/order/line_chart rows.
type OrderLinePoint struct {
	Xaxis       string  `json:"xaxis"`
	PayPrice    float64 `json:"pay_price"`
	OrderNum    int64   `json:"order_num"`
	RefundPrice float64 `json:"refund_price"`
	RefundNum   int64   `json:"refund_num"`
}

// OrderPieSlice mirrors CRMEB pie chart {name,value}.
type OrderPieSlice struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func (h *Handler) RegisterOrderAnalytics(r gin.IRoutes) {
	r.GET("/analytics/order/top", h.GetOrderTop)
	r.GET("/analytics/order/line_chart", h.GetOrderLineChart)
	r.GET("/analytics/order/pie_chart/:type", h.GetOrderPieChart)
}

func (h *Handler) GetOrderTop(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置订单统计范围")
		return
	}
	date := strings.TrimSpace(c.Query("date"))
	if date == "" {
		date = "lately7"
	}
	win, err := resolveAnalyticsWindow(date, time.Now())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "时间参数无效")
		return
	}

	payNum := h.countPaidOrders(c, scope, win.Start, win.EndExclusive)
	momPayNum := h.countPaidOrders(c, scope, win.MomStart, win.MomEndExcl)
	payMoney := h.sumPaidAmount(c, scope, win.Start, win.EndExclusive)
	momPayMoney := h.sumPaidAmount(c, scope, win.MomStart, win.MomEndExcl)
	couponMoney := h.sumCouponAmount(c, scope, win.Start, win.EndExclusive)
	momCoupon := h.sumCouponAmount(c, scope, win.MomStart, win.MomEndExcl)
	refundMoney := h.sumRefundAmount(c, scope, win.Start, win.EndExclusive, true)
	momRefundMoney := h.sumRefundAmount(c, scope, win.MomStart, win.MomEndExcl, true)
	refundNum := h.countRefundOrders(c, scope, win.Start, win.EndExclusive, false)
	momRefundNum := h.countRefundOrders(c, scope, win.MomStart, win.MomEndExcl, false)

	out := []OrderTopCard{
		{Title: "支付订单数", Count: float64(payNum), Mom: float64(momPayNum), Statistic: growthRate(payNum, momPayNum)},
		{Title: "订单实付金额", Count: payMoney, Mom: momPayMoney, Statistic: growthRateF(payMoney, momPayMoney)},
		{Title: "用券金额", Count: couponMoney, Mom: momCoupon, Statistic: growthRateF(couponMoney, momCoupon)},
		{Title: "退款金额", Count: refundMoney, Mom: momRefundMoney, Statistic: growthRateF(refundMoney, momRefundMoney)},
		{Title: "退款订单数", Count: float64(refundNum), Mom: float64(momRefundNum), Statistic: growthRate(refundNum, momRefundNum)},
	}
	response.OK(c, out)
}

func (h *Handler) GetOrderLineChart(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置订单统计范围")
		return
	}
	date := strings.TrimSpace(c.Query("date"))
	if date == "" {
		date = "lately7"
	}
	win, err := resolveAnalyticsWindow(date, time.Now())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "时间参数无效")
		return
	}

	out := make([]OrderLinePoint, 0, len(win.Buckets))
	index := map[string]int{}
	for i, b := range win.Buckets {
		index[b] = i
		out = append(out, OrderLinePoint{Xaxis: b})
	}

	type dayAgg struct {
		Bucket string  `gorm:"column:bucket"`
		Amount float64 `gorm:"column:amount"`
		Cnt    int64   `gorm:"column:cnt"`
	}

	var paidRows []dayAgg
	pq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select(fmt.Sprintf("DATE_FORMAT(created_at, '%s') AS bucket, COALESCE(SUM(pay_amount),0) AS amount, COUNT(*) AS cnt", win.SQLFormat)).
		Where("created_at >= ? AND created_at < ? AND status IN ?", win.Start, win.EndExclusive, paidOrderStatuses)
	if !scope.Full {
		pq = pq.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = pq.Group("bucket").Scan(&paidRows)
	for _, r := range paidRows {
		if i, ok := index[r.Bucket]; ok {
			out[i].PayPrice = r.Amount
			out[i].OrderNum = r.Cnt
		}
	}

	var refundRows []dayAgg
	rq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").
		Select(fmt.Sprintf("DATE_FORMAT(r.created_at, '%s') AS bucket, COALESCE(SUM(r.amount),0) AS amount, COUNT(*) AS cnt", win.SQLFormat)).
		Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").
		Where("r.created_at >= ? AND r.created_at < ? AND r.status <> ?", win.Start, win.EndExclusive, "cancelled")
	if !scope.Full {
		rq = rq.Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = rq.Group("bucket").Scan(&refundRows)
	for _, r := range refundRows {
		if i, ok := index[r.Bucket]; ok {
			out[i].RefundPrice = r.Amount
			out[i].RefundNum = r.Cnt
		}
	}

	response.OK(c, out)
}

func (h *Handler) GetOrderPieChart(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置订单统计范围")
		return
	}
	date := strings.TrimSpace(c.Query("date"))
	if date == "" {
		date = "lately7"
	}
	win, err := resolveAnalyticsWindow(date, time.Now())
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "时间参数无效")
		return
	}

	// CRMEB: type 真值 → 发货方式；假值 → 订单类型
	typeParam := strings.TrimSpace(c.Param("type"))
	asDelivery := typeParam != "" && typeParam != "0" && !strings.EqualFold(typeParam, "false")
	if asDelivery {
		response.OK(c, h.orderDeliveryPie(c, scope, win))
		return
	}
	response.OK(c, h.orderTypePie(c, scope, win))
}

func (h *Handler) orderTypePie(c *gin.Context, scope dashboardScope, win analyticsWindow) []OrderPieSlice {
	type row struct {
		ActivityType int8  `gorm:"column:activity_type"`
		Cnt          int64 `gorm:"column:cnt"`
	}
	var rows []row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select("activity_type, COUNT(*) AS cnt").
		Where("created_at >= ? AND created_at < ? AND status IN ?", win.Start, win.EndExclusive, paidOrderStatuses)
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = q.Group("activity_type").Scan(&rows)
	byType := map[int8]int64{}
	for _, r := range rows {
		byType[r.ActivityType] = r.Cnt
	}

	firstCnt := h.countFirstPaidOrders(c, scope, win.Start, win.EndExclusive)

	return []OrderPieSlice{
		{Name: "普通订单", Value: byType[0]},
		{Name: "秒杀订单", Value: byType[1]},
		{Name: "预售订单", Value: byType[2]},
		{Name: "砍价订单", Value: byType[3]},
		{Name: "拼团订单", Value: byType[4]},
		{Name: "积分订单", Value: byType[20]},
		{Name: "套餐订单", Value: byType[10]},
		{Name: "新人首单", Value: firstCnt},
	}
}

func (h *Handler) orderDeliveryPie(c *gin.Context, scope dashboardScope, win analyticsWindow) []OrderPieSlice {
	// Map qixi delivery_type → CRMEB 五类发货方式；无发货记录记为自动发货。
	labels := []struct {
		Name string
		Key  string
	}{
		{Name: "快递发货", Key: "express"},
		{Name: "配送订单", Key: "city"},
		{Name: "虚拟发货", Key: "service"},
		{Name: "核销订单", Key: "pickup"},
		{Name: "自动发货", Key: "auto"},
	}
	counts := map[string]int64{}

	type row struct {
		DeliveryType string `gorm:"column:delivery_type"`
		Cnt          int64  `gorm:"column:cnt"`
	}
	var rows []row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order AS o").
		Select("d.delivery_type, COUNT(*) AS cnt").
		Joins("JOIN qixi_crm_b_order_delivery AS d ON d.order_id = o.id").
		Where("o.created_at >= ? AND o.created_at < ? AND o.status IN ?", win.Start, win.EndExclusive, paidOrderStatuses)
	if !scope.Full {
		q = q.Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = q.Group("d.delivery_type").Scan(&rows)
	for _, r := range rows {
		counts[r.DeliveryType] = r.Cnt
	}

	var autoCnt int64
	aq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order AS o").
		Joins("LEFT JOIN qixi_crm_b_order_delivery AS d ON d.order_id = o.id").
		Where("o.created_at >= ? AND o.created_at < ? AND o.status IN ? AND d.id IS NULL", win.Start, win.EndExclusive, paidOrderStatuses)
	if !scope.Full {
		aq = aq.Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = aq.Count(&autoCnt)
	counts["auto"] = autoCnt

	out := make([]OrderPieSlice, 0, len(labels))
	for _, l := range labels {
		out = append(out, OrderPieSlice{Name: l.Name, Value: counts[l.Key]})
	}
	return out
}

func growthRateF(current, previous float64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	rate := (current - previous) / previous * 100
	return math.Round(rate*100) / 100
}

func (h *Handler) countPaidOrders(c *gin.Context, scope dashboardScope, start, endExcl time.Time) int64 {
	var n int64
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Where("created_at >= ? AND created_at < ? AND status IN ?", start, endExcl, paidOrderStatuses)
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = q.Count(&n)
	return n
}

func (h *Handler) sumPaidAmount(c *gin.Context, scope dashboardScope, start, endExcl time.Time) float64 {
	type row struct {
		Total float64 `gorm:"column:total"`
	}
	var r row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select("COALESCE(SUM(pay_amount),0) AS total").
		Where("created_at >= ? AND created_at < ? AND status IN ?", start, endExcl, paidOrderStatuses)
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = q.Scan(&r)
	return r.Total
}

func (h *Handler) sumCouponAmount(c *gin.Context, scope dashboardScope, start, endExcl time.Time) float64 {
	type row struct {
		Total float64 `gorm:"column:total"`
	}
	var r row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select("COALESCE(SUM(discount_amount),0) AS total").
		Where("created_at >= ? AND created_at < ? AND status IN ? AND discount_amount > 0", start, endExcl, paidOrderStatuses)
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = q.Scan(&r)
	return r.Total
}

func (h *Handler) sumRefundAmount(c *gin.Context, scope dashboardScope, start, endExcl time.Time, refundedOnly bool) float64 {
	type row struct {
		Total float64 `gorm:"column:total"`
	}
	var r row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").
		Select("COALESCE(SUM(r.amount),0) AS total").
		Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").
		Where("r.created_at >= ? AND r.created_at < ?", start, endExcl)
	if refundedOnly {
		q = q.Where("r.status = ?", "refunded")
	} else {
		q = q.Where("r.status <> ?", "cancelled")
	}
	if !scope.Full {
		q = q.Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = q.Scan(&r)
	return r.Total
}

func (h *Handler) countRefundOrders(c *gin.Context, scope dashboardScope, start, endExcl time.Time, refundedOnly bool) int64 {
	var n int64
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").
		Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").
		Where("r.created_at >= ? AND r.created_at < ?", start, endExcl)
	if refundedOnly {
		q = q.Where("r.status = ?", "refunded")
	} else {
		q = q.Where("r.status <> ?", "cancelled")
	}
	if !scope.Full {
		q = q.Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = q.Count(&n)
	return n
}

func (h *Handler) countFirstPaidOrders(c *gin.Context, scope dashboardScope, start, endExcl time.Time) int64 {
	var n int64
	// 用户历史首笔已支付订单落在统计窗口内，记为「新人首单」。
	sub := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select("user_id, MIN(id) AS first_id").
		Where("status IN ?", paidOrderStatuses).
		Group("user_id")
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order AS o").
		Joins("JOIN (?) AS f ON f.first_id = o.id", sub).
		Where("o.created_at >= ? AND o.created_at < ?", start, endExcl)
	if !scope.Full {
		q = q.Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = q.Count(&n)
	return n
}
