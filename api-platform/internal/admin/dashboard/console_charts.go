package dashboard

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HourAmountPoint is one hour bucket for the「当日订单金额」dual line.
type HourAmountPoint struct {
	Hour            string  `json:"hour"`
	TodayAmount     float64 `json:"today_amount"`
	YesterdayAmount float64 `json:"yesterday_amount"`
}

// SparkStat is one cell in「数据统计」2×2 grid.
type SparkStat struct {
	Value    int64     `json:"value"`
	Ratio    float64   `json:"ratio"`
	Spark    []float64 `json:"spark"`
	SparkLabel string  `json:"spark_label"`
}

// OrderStatsBlock holds the four「数据统计」cells.
type OrderStatsBlock struct {
	TodayOrderCount SparkStat `json:"today_order_count"`
	TodayPayerCount SparkStat `json:"today_payer_count"`
	MonthOrderCount SparkStat `json:"month_order_count"`
	MonthPayerCount SparkStat `json:"month_payer_count"`
}

// UserTrendPoint is one day on the「用户数据」chart.
type UserTrendPoint struct {
	Day          string `json:"day"`
	NewUsers     int64  `json:"new_users"`
	VisitUsers   int64  `json:"visit_users"`
	TotalUsers   int64  `json:"total_users"`
}

// DealUserBlock is「成交用户」grid + funnel.
type DealUserBlock struct {
	VisitUsers       int64   `json:"visit_users"`
	OrderUsers       int64   `json:"order_users"`
	OrderAmount      float64 `json:"order_amount"`
	PayUsers         int64   `json:"pay_users"`
	PayAmount        float64 `json:"pay_amount"`
	AvgOrderAmount   float64 `json:"avg_order_amount"`
	VisitOrderRate   float64 `json:"visit_order_rate"`
	OrderPayRate     float64 `json:"order_pay_rate"`
}

// DealRatioBlock is「成交用户占比」donut.
type DealRatioBlock struct {
	NewUsers       int64   `json:"new_users"`
	OldUsers       int64   `json:"old_users"`
	NewAmount      float64 `json:"new_amount"`
	OldAmount      float64 `json:"old_amount"`
}

func (h *Handler) RegisterConsole(r gin.IRoutes) {
	r.GET("/dashboard/merchant-top", h.GetMerchantTop)
	r.GET("/dashboard/user-trend", h.GetUserTrend)
	r.GET("/dashboard/deal", h.GetDeal)
}

func (h *Handler) GetMerchantTop(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置首页数据范围")
		return
	}
	period := normalizePeriod(c.Query("period"), "month")
	out := []StoreSalesRank{}
	if err := h.storeSalesRankPeriod(c, scope, period, &out); err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载店铺销售排行失败")
		return
	}
	if out == nil {
		out = []StoreSalesRank{}
	}
	response.OK(c, gin.H{"list": out, "period": period})
}

func (h *Handler) GetUserTrend(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置首页数据范围")
		return
	}
	period := normalizePeriod(c.Query("period"), "30d")
	if period == "year" {
		period = "30d"
	}
	points, err := h.userTrend(c, scope, period)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载用户趋势失败")
		return
	}
	response.OK(c, gin.H{"list": points, "period": period})
}

func (h *Handler) GetDeal(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置首页数据范围")
		return
	}
	period := normalizePeriod(c.Query("period"), "month")
	funnel, ratio, err := h.dealBlocks(c, scope, period)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载成交用户数据失败")
		return
	}
	response.OK(c, gin.H{"funnel": funnel, "ratio": ratio, "period": period})
}

func normalizePeriod(raw, fallback string) string {
	switch strings.TrimSpace(raw) {
	case "7d", "lately7", "week":
		return "7d"
	case "30d", "lately30":
		return "30d"
	case "month", "this_month":
		return "month"
	case "year", "this_year":
		return "year"
	default:
		return fallback
	}
}

func periodCond(col, period string) string {
	switch period {
	case "7d":
		return col + " >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)"
	case "30d":
		return col + " >= DATE_SUB(CURDATE(), INTERVAL 29 DAY)"
	case "year":
		return "YEAR(" + col + ") = YEAR(CURDATE())"
	default: // month
		return col + " >= DATE_FORMAT(CURDATE(), '%Y-%m-01')"
	}
}

func (h *Handler) storeSalesRankPeriod(c *gin.Context, scope dashboardScope, period string, out *[]StoreSalesRank) error {
	q := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order AS o").
		Select(`
			o.store_id,
			MAX(o.store_name_snapshot) AS store_name,
			COALESCE(MAX(f.follower_count), 0) AS follower_count,
			COALESCE(SUM(o.total_quantity), 0) AS sale_count,
			COALESCE(SUM(o.pay_amount), 0) AS sale_amount`).
		Joins(`LEFT JOIN (
			SELECT store_id, COUNT(DISTINCT user_id) AS follower_count
			FROM qixi_crm_b_user_follow_store
			GROUP BY store_id
		) AS f ON f.store_id = o.store_id`).
		Where("o.status IN ('paid','fulfilling','shipped','completed') AND " + periodCond("o.paid_at", period))
	if !scope.Full {
		q = q.Where("o.store_id IN ?", scope.StoreIDs)
	}
	return q.
		Group("o.store_id").
		Order("sale_amount DESC, sale_count DESC, o.store_id ASC").
		Limit(10).
		Scan(out).Error
}

func (h *Handler) todayOrderHours(c *gin.Context, scope dashboardScope) ([]HourAmountPoint, error) {
	type row struct {
		HourBucket      int     `gorm:"column:hour_bucket"`
		TodayAmount     float64 `gorm:"column:today_amount"`
		YesterdayAmount float64 `gorm:"column:yesterday_amount"`
	}
	var rows []row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select(`
			HOUR(paid_at) AS hour_bucket,
			COALESCE(SUM(CASE WHEN DATE(paid_at) = CURDATE() THEN pay_amount ELSE 0 END), 0) AS today_amount,
			COALESCE(SUM(CASE WHEN DATE(paid_at) = DATE_SUB(CURDATE(), INTERVAL 1 DAY) THEN pay_amount ELSE 0 END), 0) AS yesterday_amount`).
		Where("status IN ('paid','fulfilling','shipped','completed') AND paid_at >= DATE_SUB(CURDATE(), INTERVAL 1 DAY)")
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	if err := q.Group("HOUR(paid_at)").Scan(&rows).Error; err != nil {
		return nil, err
	}
	m := map[int]row{}
	for _, r := range rows {
		m[r.HourBucket] = r
	}
	out := make([]HourAmountPoint, 0, 24)
	for hour := 0; hour < 24; hour++ {
		r := m[hour]
		out = append(out, HourAmountPoint{
			Hour:            fmt.Sprintf("%02d:00", hour),
			TodayAmount:     r.TodayAmount,
			YesterdayAmount: r.YesterdayAmount,
		})
	}
	return out, nil
}

func (h *Handler) orderStatsBlock(c *gin.Context, scope dashboardScope) (OrderStatsBlock, error) {
	var out OrderStatsBlock
	type agg struct {
		TodayOrder      int64   `gorm:"column:today_order"`
		YesterdayOrder  int64   `gorm:"column:yesterday_order"`
		TodayPayer      int64   `gorm:"column:today_payer"`
		YesterdayPayer  int64   `gorm:"column:yesterday_payer"`
		MonthOrder      int64   `gorm:"column:month_order"`
		LastMonthOrder  int64   `gorm:"column:last_month_order"`
		MonthPayer      int64   `gorm:"column:month_payer"`
		LastMonthPayer  int64   `gorm:"column:last_month_payer"`
	}
	var a agg
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").Select(`
		COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN 1 ELSE 0 END), 0) AS today_order,
		COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = DATE_SUB(CURDATE(), INTERVAL 1 DAY) THEN 1 ELSE 0 END), 0) AS yesterday_order,
		COALESCE(COUNT(DISTINCT CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN user_id END), 0) AS today_payer,
		COALESCE(COUNT(DISTINCT CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = DATE_SUB(CURDATE(), INTERVAL 1 DAY) THEN user_id END), 0) AS yesterday_payer,
		COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND paid_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END), 0) AS month_order,
		COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND paid_at >= DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-01') AND paid_at < DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END), 0) AS last_month_order,
		COALESCE(COUNT(DISTINCT CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND paid_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN user_id END), 0) AS month_payer,
		COALESCE(COUNT(DISTINCT CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND paid_at >= DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-01') AND paid_at < DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN user_id END), 0) AS last_month_payer`)
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	if err := q.Scan(&a).Error; err != nil {
		return out, err
	}
	hourSpark, err := h.hourSparkline(c, scope, "order")
	if err != nil {
		return out, err
	}
	hourPayerSpark, err := h.hourSparkline(c, scope, "payer")
	if err != nil {
		return out, err
	}
	dayOrderSpark, err := h.monthDaySparkline(c, scope, "order")
	if err != nil {
		return out, err
	}
	dayPayerSpark, err := h.monthDaySparkline(c, scope, "payer")
	if err != nil {
		return out, err
	}
	out.TodayOrderCount = SparkStat{Value: a.TodayOrder, Ratio: growthRatio(a.TodayOrder, a.YesterdayOrder), Spark: hourSpark, SparkLabel: "今天"}
	out.TodayPayerCount = SparkStat{Value: a.TodayPayer, Ratio: growthRatio(a.TodayPayer, a.YesterdayPayer), Spark: hourPayerSpark, SparkLabel: "今天"}
	out.MonthOrderCount = SparkStat{Value: a.MonthOrder, Ratio: growthRatio(a.MonthOrder, a.LastMonthOrder), Spark: dayOrderSpark, SparkLabel: "本月"}
	out.MonthPayerCount = SparkStat{Value: a.MonthPayer, Ratio: growthRatio(a.MonthPayer, a.LastMonthPayer), Spark: dayPayerSpark, SparkLabel: "本月"}
	return out, nil
}

func growthRatio(now, prev int64) float64 {
	if prev == 0 {
		if now == 0 {
			return 0
		}
		return 100
	}
	return float64(now-prev) / float64(abs64(prev)) * 100
}

func abs64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func (h *Handler) hourSparkline(c *gin.Context, scope dashboardScope, kind string) ([]float64, error) {
	type row struct {
		HourBucket int     `gorm:"column:hour_bucket"`
		Value      float64 `gorm:"column:value"`
	}
	selectExpr := "HOUR(paid_at) AS hour_bucket, COUNT(*) AS value"
	if kind == "payer" {
		selectExpr = "HOUR(paid_at) AS hour_bucket, COUNT(DISTINCT user_id) AS value"
	}
	var rows []row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select(selectExpr).
		Where("status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE()")
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	if err := q.Group("HOUR(paid_at)").Scan(&rows).Error; err != nil {
		return nil, err
	}
	m := map[int]float64{}
	for _, r := range rows {
		m[r.HourBucket] = r.Value
	}
	out := make([]float64, 24)
	for i := 0; i < 24; i++ {
		out[i] = m[i]
	}
	return out, nil
}

func (h *Handler) monthDaySparkline(c *gin.Context, scope dashboardScope, kind string) ([]float64, error) {
	type row struct {
		DayBucket int     `gorm:"column:day_bucket"`
		Value     float64 `gorm:"column:value"`
	}
	selectExpr := "DAY(paid_at) AS day_bucket, COUNT(*) AS value"
	if kind == "payer" {
		selectExpr = "DAY(paid_at) AS day_bucket, COUNT(DISTINCT user_id) AS value"
	}
	var rows []row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select(selectExpr).
		Where("status IN ('paid','fulfilling','shipped','completed') AND paid_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')")
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	if err := q.Group("DAY(paid_at)").Scan(&rows).Error; err != nil {
		return nil, err
	}
	now := time.Now()
	days := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()
	m := map[int]float64{}
	for _, r := range rows {
		m[r.DayBucket] = r.Value
	}
	out := make([]float64, days)
	for i := 1; i <= days; i++ {
		out[i-1] = m[i]
	}
	return out, nil
}

func (h *Handler) userTrend(c *gin.Context, scope dashboardScope, period string) ([]UserTrendPoint, error) {
	days := periodDays(period)
	start := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -(days - 1))
	if period == "month" {
		now := time.Now()
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		days = now.Day()
	}

	type dayCount struct {
		Day   string `gorm:"column:day_key"`
		Count int64  `gorm:"column:cnt"`
	}

	newMap := map[string]int64{}
	if scope.Full {
		var rows []dayCount
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user").
			Select("DATE_FORMAT(created_at, '%Y-%m-%d') AS day_key, COUNT(*) AS cnt").
			Where("created_at >= ?", start).
			Group("DATE_FORMAT(created_at, '%Y-%m-%d')").
			Scan(&rows)
		for _, r := range rows {
			newMap[r.Day] = r.Count
		}
	} else {
		var rows []dayCount
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
			Select("DATE_FORMAT(created_at, '%Y-%m-%d') AS day_key, COUNT(DISTINCT user_id) AS cnt").
			Where("store_id IN ? AND created_at >= ?", scope.StoreIDs, start).
			Group("DATE_FORMAT(created_at, '%Y-%m-%d')").
			Scan(&rows)
		for _, r := range rows {
			newMap[r.Day] = r.Count
		}
	}

	visitMap := map[string]int64{}
	{
		var rows []dayCount
		q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_browse_history").
			Select("DATE_FORMAT(viewed_at, '%Y-%m-%d') AS day_key, COUNT(DISTINCT user_id) AS cnt").
			Where("viewed_at >= ?", start)
		if !scope.Full {
			q = q.Where("store_id IN ?", scope.StoreIDs)
		}
		_ = q.Group("DATE_FORMAT(viewed_at, '%Y-%m-%d')").Scan(&rows)
		for _, r := range rows {
			visitMap[r.Day] = r.Count
		}
	}

	var baseTotal int64
	if scope.Full {
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user").
			Where("created_at < ?", start).Count(&baseTotal)
	}

	dailyNew := make([]int64, days)
	out := make([]UserTrendPoint, 0, days)
	running := baseTotal
	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i)
		key := d.Format("2006-01-02")
		n := newMap[key]
		dailyNew[i] = n
		running += n
		out = append(out, UserTrendPoint{
			Day:        d.Format("01-02"),
			NewUsers:   n,
			VisitUsers: visitMap[key],
			TotalUsers: running,
		})
	}
	return out, nil
}

func periodDays(period string) int {
	switch period {
	case "7d":
		return 7
	case "month":
		return time.Now().Day()
	default:
		return 30
	}
}

func (h *Handler) dealBlocks(c *gin.Context, scope dashboardScope, period string) (DealUserBlock, DealRatioBlock, error) {
	var funnel DealUserBlock
	var ratio DealRatioBlock

	visitQ := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_browse_history").
		Where(periodCond("viewed_at", period))
	if !scope.Full {
		visitQ = visitQ.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = visitQ.Distinct("user_id").Count(&funnel.VisitUsers)

	paidCond := periodCond("paid_at", period)
	createdCond := periodCond("created_at", period)
	type orderAgg struct {
		OrderUsers  int64   `gorm:"column:order_users"`
		OrderAmount float64 `gorm:"column:order_amount"`
		PayUsers    int64   `gorm:"column:pay_users"`
		PayAmount   float64 `gorm:"column:pay_amount"`
	}
	var oa orderAgg
	oq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").Select(`
		COALESCE(COUNT(DISTINCT CASE WHEN ` + createdCond + ` THEN user_id END), 0) AS order_users,
		COALESCE(SUM(CASE WHEN ` + createdCond + ` THEN pay_amount ELSE 0 END), 0) AS order_amount,
		COALESCE(COUNT(DISTINCT CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND ` + paidCond + ` THEN user_id END), 0) AS pay_users,
		COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND ` + paidCond + ` THEN pay_amount ELSE 0 END), 0) AS pay_amount`)
	if !scope.Full {
		oq = oq.Where("store_id IN ?", scope.StoreIDs)
	}
	if err := oq.Scan(&oa).Error; err != nil {
		return funnel, ratio, err
	}
	funnel.OrderUsers = oa.OrderUsers
	funnel.OrderAmount = oa.OrderAmount
	funnel.PayUsers = oa.PayUsers
	funnel.PayAmount = oa.PayAmount
	if oa.PayUsers > 0 {
		funnel.AvgOrderAmount = oa.PayAmount / float64(oa.PayUsers)
	}
	if funnel.VisitUsers > 0 {
		funnel.VisitOrderRate = float64(funnel.OrderUsers) / float64(funnel.VisitUsers) * 100
	}
	if funnel.OrderUsers > 0 {
		funnel.OrderPayRate = float64(funnel.PayUsers) / float64(funnel.OrderUsers) * 100
	}

	userCond := periodCond("u.created_at", period)
	type ratioRow struct {
		NewUsers  int64   `gorm:"column:new_users"`
		OldUsers  int64   `gorm:"column:old_users"`
		NewAmount float64 `gorm:"column:new_amount"`
		OldAmount float64 `gorm:"column:old_amount"`
	}
	var rr ratioRow
	rq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order AS o").
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id = o.user_id").
		Select(`
			COALESCE(COUNT(DISTINCT CASE WHEN ` + userCond + ` THEN o.user_id END), 0) AS new_users,
			COALESCE(COUNT(DISTINCT CASE WHEN NOT (` + userCond + `) OR u.id IS NULL THEN o.user_id END), 0) AS old_users,
			COALESCE(SUM(CASE WHEN ` + userCond + ` THEN o.pay_amount ELSE 0 END), 0) AS new_amount,
			COALESCE(SUM(CASE WHEN NOT (` + userCond + `) OR u.id IS NULL THEN o.pay_amount ELSE 0 END), 0) AS old_amount`).
		Where("o.status IN ('paid','fulfilling','shipped','completed') AND " + periodCond("o.paid_at", period))
	if !scope.Full {
		rq = rq.Where("o.store_id IN ?", scope.StoreIDs)
	}
	if err := rq.Scan(&rr).Error; err != nil {
		return funnel, ratio, nil
	}
	ratio.NewUsers = rr.NewUsers
	ratio.OldUsers = rr.OldUsers
	ratio.NewAmount = rr.NewAmount
	ratio.OldAmount = rr.OldAmount
	return funnel, ratio, nil
}

// safeCount ignores missing-table errors so optional todo buckets stay at 0.
func safeCount(db *gorm.DB, out *int64) {
	if err := db.Count(out).Error; err != nil {
		*out = 0
	}
}
