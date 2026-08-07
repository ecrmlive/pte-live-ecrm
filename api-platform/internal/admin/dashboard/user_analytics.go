package dashboard

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// UserTopCard mirrors CRMEB analytics/user/top rows.
type UserTopCard struct {
	Title     string  `json:"title"`
	Count     int64   `json:"count"`
	Mom       int64   `json:"mom"`
	Statistic float64 `json:"statistic"`
}

// UserLinePoint mirrors CRMEB analytics/user/line_chart rows.
type UserLinePoint struct {
	Xaxis string `json:"xaxis"`
	Count int64  `json:"count"`
}

// UserDealPoint mirrors CRMEB analytics/user/pie_chart (grouped bar: old/new).
type UserDealPoint struct {
	Xaxis string `json:"xaxis"`
	Old   int64  `json:"old"`
	New   int64  `json:"new"`
}

func (h *Handler) RegisterUserAnalytics(r gin.IRoutes) {
	r.GET("/analytics/user/top", h.GetUserTop)
	r.GET("/analytics/user/line_chart", h.GetUserLineChart)
	r.GET("/analytics/user/pie_chart", h.GetUserDealChart)
}

func (h *Handler) GetUserTop(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置用户统计范围")
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

	totalUsers := h.countUsers(c, time.Time{}, time.Time{})
	newUsers := h.countUsers(c, win.Start, win.EndExclusive)
	momNew := h.countUsers(c, win.MomStart, win.MomEndExcl)
	orderUsers := h.countPaidOrderUsers(c, scope, win.Start, win.EndExclusive)
	momOrder := h.countPaidOrderUsers(c, scope, win.MomStart, win.MomEndExcl)
	active := h.countActiveUsers(c, scope, win.Start, win.EndExclusive)
	momActive := h.countActiveUsers(c, scope, win.MomStart, win.MomEndExcl)
	svipTotal := h.countActiveSvip(c)
	newSvip := h.countNewSvip(c, win.Start, win.EndExclusive)
	momSvip := h.countNewSvip(c, win.MomStart, win.MomEndExcl)

	out := []UserTopCard{
		{Title: "用户数量", Count: totalUsers, Mom: 0, Statistic: 0},
		{Title: "新增用户", Count: newUsers, Mom: momNew, Statistic: growthRate(newUsers, momNew)},
		{Title: "下单用户", Count: orderUsers, Mom: momOrder, Statistic: growthRate(orderUsers, momOrder)},
		{Title: "活跃用户", Count: active, Mom: momActive, Statistic: growthRate(active, momActive)},
		{Title: "付费会员", Count: svipTotal, Mom: 0, Statistic: 0},
		{Title: "新增付费会员", Count: newSvip, Mom: momSvip, Statistic: growthRate(newSvip, momSvip)},
	}
	response.OK(c, out)
}

func (h *Handler) GetUserLineChart(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置用户统计范围")
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
	// CRMEB: 0 新增用户 / 1 活跃用户 / 2 新增付费会员
	typeParam := strings.TrimSpace(c.DefaultQuery("type", "0"))

	out := make([]UserLinePoint, 0, len(win.Buckets))
	index := map[string]int{}
	for i, b := range win.Buckets {
		index[b] = i
		out = append(out, UserLinePoint{Xaxis: b})
	}

	type dayCount struct {
		Bucket string `gorm:"column:bucket"`
		Count  int64  `gorm:"column:cnt"`
	}
	fill := func(rows []dayCount) {
		for _, r := range rows {
			if i, ok := index[r.Bucket]; ok {
				out[i].Count = r.Count
			}
		}
	}

	var rows []dayCount
	switch typeParam {
	case "1":
		q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_browse_history").
			Select(fmt.Sprintf("DATE_FORMAT(viewed_at, '%s') AS bucket, COUNT(DISTINCT user_id) AS cnt", win.SQLFormat)).
			Where("viewed_at >= ? AND viewed_at < ?", win.Start, win.EndExclusive)
		if !scope.Full {
			q = q.Where("store_id IN ?", scope.StoreIDs)
		}
		_ = q.Group("bucket").Scan(&rows)
	case "2":
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_order").
			Select(fmt.Sprintf("DATE_FORMAT(COALESCE(paid_at, created_at), '%s') AS bucket, COUNT(*) AS cnt", win.SQLFormat)).
			Where("status = 'paid' AND COALESCE(paid_at, created_at) >= ? AND COALESCE(paid_at, created_at) < ?", win.Start, win.EndExclusive).
			Group("bucket").Scan(&rows)
	default:
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user").
			Select(fmt.Sprintf("DATE_FORMAT(created_at, '%s') AS bucket, COUNT(*) AS cnt", win.SQLFormat)).
			Where("created_at >= ? AND created_at < ?", win.Start, win.EndExclusive).
			Group("bucket").Scan(&rows)
	}
	fill(rows)
	response.OK(c, out)
}

func (h *Handler) GetUserDealChart(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置用户统计范围")
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

	out := make([]UserDealPoint, 0, len(win.Buckets))
	index := map[string]int{}
	for i, b := range win.Buckets {
		index[b] = i
		out = append(out, UserDealPoint{Xaxis: b})
	}

	type dayCount struct {
		Bucket string `gorm:"column:bucket"`
		Count  int64  `gorm:"column:cnt"`
	}

	// 新用户：窗口内首次支付落在该桶；老用户：该桶有支付但首次支付早于该桶。
	var newRows []dayCount
	nq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_group_order AS g").
		Select(fmt.Sprintf("DATE_FORMAT(g.paid_at, '%s') AS bucket, COUNT(DISTINCT g.user_id) AS cnt", win.SQLFormat)).
		Joins(`JOIN (
			SELECT user_id, MIN(paid_at) AS first_paid
			FROM qixi_crm_b_group_order
			WHERE pay_status = 'paid' AND paid_at IS NOT NULL
			GROUP BY user_id
		) AS f ON f.user_id = g.user_id AND f.first_paid = g.paid_at`).
		Where("g.pay_status = 'paid' AND g.paid_at >= ? AND g.paid_at < ?", win.Start, win.EndExclusive)
	if !scope.Full {
		nq = nq.Joins("JOIN qixi_crm_b_order AS o ON o.group_order_id = g.id").
			Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = nq.Group("bucket").Scan(&newRows)
	for _, r := range newRows {
		if i, ok := index[r.Bucket]; ok {
			out[i].New = r.Count
		}
	}

	var oldRows []dayCount
	oq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_group_order AS g").
		Select(fmt.Sprintf("DATE_FORMAT(g.paid_at, '%s') AS bucket, COUNT(DISTINCT g.user_id) AS cnt", win.SQLFormat)).
		Joins(`JOIN (
			SELECT user_id, MIN(paid_at) AS first_paid
			FROM qixi_crm_b_group_order
			WHERE pay_status = 'paid' AND paid_at IS NOT NULL
			GROUP BY user_id
		) AS f ON f.user_id = g.user_id AND DATE(f.first_paid) < DATE(g.paid_at)`).
		Where("g.pay_status = 'paid' AND g.paid_at >= ? AND g.paid_at < ?", win.Start, win.EndExclusive)
	if !scope.Full {
		oq = oq.Joins("JOIN qixi_crm_b_order AS o ON o.group_order_id = g.id").
			Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = oq.Group("bucket").Scan(&oldRows)
	for _, r := range oldRows {
		if i, ok := index[r.Bucket]; ok {
			out[i].Old = r.Count
		}
	}

	response.OK(c, out)
}

func (h *Handler) countUsers(c *gin.Context, start, endExcl time.Time) int64 {
	var n int64
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user")
	if !start.IsZero() {
		q = q.Where("created_at >= ? AND created_at < ?", start, endExcl)
	}
	_ = q.Count(&n)
	return n
}

func (h *Handler) countPaidOrderUsers(c *gin.Context, scope dashboardScope, start, endExcl time.Time) int64 {
	var n int64
	if scope.Full {
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_group_order").
			Where("pay_status = 'paid' AND paid_at >= ? AND paid_at < ?", start, endExcl).
			Distinct("user_id").Count(&n)
		return n
	}
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Where("store_id IN ? AND status IN ? AND paid_at >= ? AND paid_at < ?",
			scope.StoreIDs, paidOrderStatuses, start, endExcl).
		Distinct("user_id").Count(&n)
	return n
}

func (h *Handler) countActiveUsers(c *gin.Context, scope dashboardScope, start, endExcl time.Time) int64 {
	var n int64
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_browse_history").
		Where("viewed_at >= ? AND viewed_at < ?", start, endExcl)
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = q.Distinct("user_id").Count(&n)
	return n
}

func (h *Handler) countActiveSvip(c *gin.Context) int64 {
	var n int64
	now := time.Now().In(shanghaiLoc())
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_svip").
		Where("status = 'lifetime' OR (expires_at IS NOT NULL AND expires_at > ?)", now).
		Count(&n)
	return n
}

func (h *Handler) countNewSvip(c *gin.Context, start, endExcl time.Time) int64 {
	var n int64
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_order").
		Where("status = 'paid' AND COALESCE(paid_at, created_at) >= ? AND COALESCE(paid_at, created_at) < ?", start, endExcl).
		Count(&n)
	return n
}
