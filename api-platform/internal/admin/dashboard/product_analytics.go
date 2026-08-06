package dashboard

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// ProductTopCard mirrors CRMEB analytics/product/top rows.
type ProductTopCard struct {
	Title     string  `json:"title"`
	Count     int64   `json:"count"`
	Mom       int64   `json:"mom"`
	Statistic float64 `json:"statistic"`
}

// ProductLinePoint mirrors CRMEB analytics/product/line_chart rows.
type ProductLinePoint struct {
	Xaxis    string `json:"xaxis"`
	Visit    int64  `json:"visit"`
	Relation int64  `json:"relation"`
	TotalNum int64  `json:"total_num"`
	PaidNum  int64  `json:"paid_num"`
}

// ProductPieSlice mirrors CRMEB pie chart {name,value}.
type ProductPieSlice struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func (h *Handler) RegisterProductAnalytics(r gin.IRoutes) {
	r.GET("/analytics/product/top", h.GetProductTop)
	r.GET("/analytics/product/line_chart", h.GetProductLineChart)
	r.GET("/analytics/product/pie_chart/:type", h.GetProductPieChart)
}

func (h *Handler) GetProductTop(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置商品统计范围")
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

	visit := h.countBrowse(c, scope, win.Start, win.EndExclusive)
	momVisit := h.countBrowse(c, scope, win.MomStart, win.MomEndExcl)
	relation := h.countFavorite(c, scope, win.Start, win.EndExclusive)
	momRelation := h.countFavorite(c, scope, win.MomStart, win.MomEndExcl)
	cart := h.sumCart(c, scope, win.Start, win.EndExclusive)
	momCart := h.sumCart(c, scope, win.MomStart, win.MomEndExcl)
	totalNum := h.sumOrderQty(c, scope, win.Start, win.EndExclusive, false)
	momTotal := h.sumOrderQty(c, scope, win.MomStart, win.MomEndExcl, false)
	paidNum := h.sumOrderQty(c, scope, win.Start, win.EndExclusive, true)
	momPaid := h.sumOrderQty(c, scope, win.MomStart, win.MomEndExcl, true)
	refundNum := h.countRefund(c, scope, win.Start, win.EndExclusive)
	momRefund := h.countRefund(c, scope, win.MomStart, win.MomEndExcl)

	out := []ProductTopCard{
		{Title: "浏览量", Count: visit, Mom: momVisit, Statistic: growthRate(visit, momVisit)},
		{Title: "收藏量", Count: relation, Mom: momRelation, Statistic: growthRate(relation, momRelation)},
		{Title: "加购数", Count: cart, Mom: momCart, Statistic: growthRate(cart, momCart)},
		{Title: "下单数", Count: totalNum, Mom: momTotal, Statistic: growthRate(totalNum, momTotal)},
		{Title: "支付数", Count: paidNum, Mom: momPaid, Statistic: growthRate(paidNum, momPaid)},
		{Title: "退款数", Count: refundNum, Mom: momRefund, Statistic: growthRate(refundNum, momRefund)},
	}
	response.OK(c, out)
}

func (h *Handler) GetProductLineChart(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置商品统计范围")
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

	out := make([]ProductLinePoint, 0, len(win.Buckets))
	index := map[string]int{}
	for i, b := range win.Buckets {
		index[b] = i
		out = append(out, ProductLinePoint{Xaxis: b})
	}

	type dayCount struct {
		Bucket string `gorm:"column:bucket"`
		Count  int64  `gorm:"column:cnt"`
	}

	fill := func(rows []dayCount, set func(p *ProductLinePoint, n int64)) {
		for _, r := range rows {
			if i, ok := index[r.Bucket]; ok {
				set(&out[i], r.Count)
			}
		}
	}

	var visitRows []dayCount
	vq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_browse_history").
		Select(fmt.Sprintf("DATE_FORMAT(viewed_at, '%s') AS bucket, COUNT(*) AS cnt", win.SQLFormat)).
		Where("viewed_at >= ? AND viewed_at < ?", win.Start, win.EndExclusive)
	if !scope.Full {
		vq = vq.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = vq.Group("bucket").Scan(&visitRows)
	fill(visitRows, func(p *ProductLinePoint, n int64) { p.Visit = n })

	var favRows []dayCount
	fq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_favorite AS f").
		Select(fmt.Sprintf("DATE_FORMAT(f.created_at, '%s') AS bucket, COUNT(*) AS cnt", win.SQLFormat)).
		Where("f.created_at >= ? AND f.created_at < ?", win.Start, win.EndExclusive)
	if !scope.Full {
		fq = fq.Joins("JOIN qixi_crm_b_product_view AS p ON p.product_id = f.product_id").
			Where("p.store_id IN ?", scope.StoreIDs)
	}
	_ = fq.Group("bucket").Scan(&favRows)
	fill(favRows, func(p *ProductLinePoint, n int64) { p.Relation = n })

	var orderRows []dayCount
	oq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select(fmt.Sprintf("DATE_FORMAT(created_at, '%s') AS bucket, COALESCE(SUM(total_quantity),0) AS cnt", win.SQLFormat)).
		Where("created_at >= ? AND created_at < ?", win.Start, win.EndExclusive)
	if !scope.Full {
		oq = oq.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = oq.Group("bucket").Scan(&orderRows)
	fill(orderRows, func(p *ProductLinePoint, n int64) { p.TotalNum = n })

	var paidRows []dayCount
	pq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select(fmt.Sprintf("DATE_FORMAT(created_at, '%s') AS bucket, COALESCE(SUM(total_quantity),0) AS cnt", win.SQLFormat)).
		Where("created_at >= ? AND created_at < ? AND status IN ?", win.Start, win.EndExclusive, paidOrderStatuses)
	if !scope.Full {
		pq = pq.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = pq.Group("bucket").Scan(&paidRows)
	fill(paidRows, func(p *ProductLinePoint, n int64) { p.PaidNum = n })

	response.OK(c, out)
}

func (h *Handler) GetProductPieChart(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置商品统计范围")
		return
	}
	// CRMEB: type 真值 → 分类；假值 → 商品类型
	typeParam := strings.TrimSpace(c.Param("type"))
	asCategory := typeParam != "" && typeParam != "0" && !strings.EqualFold(typeParam, "false")

	if asCategory {
		type row struct {
			Name  string `gorm:"column:name"`
			Count int64  `gorm:"column:cnt"`
		}
		var rows []row
		q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_view AS p").
			Select("COALESCE(NULLIF(c.name,''),'未分类') AS name, COUNT(*) AS cnt").
			Joins("LEFT JOIN qixi_crm_b_category_view AS c ON c.category_id = p.category_id").
			Where("p.sale_status = 1")
		if !scope.Full {
			q = q.Where("p.store_id IN ?", scope.StoreIDs)
		}
		_ = q.Group("COALESCE(NULLIF(c.name,''),'未分类')").Order("cnt DESC").Limit(12).Scan(&rows)
		out := make([]ProductPieSlice, 0, len(rows))
		for _, r := range rows {
			out = append(out, ProductPieSlice{Name: r.Name, Value: r.Count})
		}
		response.OK(c, out)
		return
	}

	names := []string{"普通商品", "虚拟商品", "网盘商品", "卡密商品", "预约商品"}
	type row struct {
		Type  int8  `gorm:"column:product_type"`
		Count int64 `gorm:"column:cnt"`
	}
	var rows []row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_view").
		Select("product_type, COUNT(*) AS cnt").
		Where("sale_status = 1")
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = q.Group("product_type").Order("product_type ASC").Scan(&rows)
	out := make([]ProductPieSlice, 0, len(rows))
	for _, r := range rows {
		name := fmt.Sprintf("类型%d", r.Type)
		if int(r.Type) >= 0 && int(r.Type) < len(names) {
			name = names[r.Type]
		}
		out = append(out, ProductPieSlice{Name: name, Value: r.Count})
	}
	response.OK(c, out)
}

var paidOrderStatuses = []string{"paid", "fulfilling", "shipped", "completed"}

type analyticsWindow struct {
	Start         time.Time
	End           time.Time
	EndExclusive  time.Time
	MomStart      time.Time
	MomEndExcl    time.Time
	Buckets       []string
	SQLFormat     string
}

func shanghaiLoc() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.FixedZone("CST", 8*3600)
	}
	return loc
}

func resolveAnalyticsWindow(date string, now time.Time) (analyticsWindow, error) {
	loc := shanghaiLoc()
	now = now.In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	var start, endInclusive time.Time
	var sqlFmt string
	var buckets []string

	switch date {
	case "lately7":
		start = today.AddDate(0, 0, -6)
		endInclusive = today
		sqlFmt = "%m-%d"
		for d := start; !d.After(endInclusive); d = d.AddDate(0, 0, 1) {
			buckets = append(buckets, d.Format("01-02"))
		}
	case "lately30":
		start = today.AddDate(0, 0, -29)
		endInclusive = today
		sqlFmt = "%m-%d"
		for d := start; !d.After(endInclusive); d = d.AddDate(0, 0, 1) {
			buckets = append(buckets, d.Format("01-02"))
		}
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		endInclusive = time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, loc)
		sqlFmt = "%d"
		days := endInclusive.Day()
		for day := 1; day <= days; day++ {
			buckets = append(buckets, fmt.Sprintf("%02d", day))
		}
	case "year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
		endInclusive = today
		sqlFmt = "%m"
		for m := time.January; m <= now.Month(); m++ {
			buckets = append(buckets, fmt.Sprintf("%02d", int(m)))
		}
	default:
		parts := strings.Split(date, "-")
		if len(parts) != 2 {
			return analyticsWindow{}, fmt.Errorf("invalid date %q", date)
		}
		s, err1 := parseSlashDate(parts[0], loc)
		e, err2 := parseSlashDate(parts[1], loc)
		if err1 != nil || err2 != nil || e.Before(s) {
			return analyticsWindow{}, fmt.Errorf("invalid date range %q", date)
		}
		start, endInclusive = s, e
		diffDays := int(e.Sub(s).Hours()/24) + 1
		if diffDays <= 30 {
			sqlFmt = "%m-%d"
			for d := s; !d.After(e); d = d.AddDate(0, 0, 1) {
				buckets = append(buckets, d.Format("01-02"))
			}
		} else if s.Year() == e.Year() {
			sqlFmt = "%m"
			for m := s.Month(); m <= e.Month(); m++ {
				buckets = append(buckets, fmt.Sprintf("%02d", int(m)))
			}
		} else {
			sqlFmt = "%Y"
			for y := s.Year(); y <= e.Year(); y++ {
				buckets = append(buckets, strconv.Itoa(y))
			}
		}
	}

	endExclusive := endInclusive.AddDate(0, 0, 1)
	duration := endExclusive.Sub(start)
	momEndExcl := start
	momStart := momEndExcl.Add(-duration)
	if date == "lately7" {
		momStart = today.AddDate(0, 0, -13)
		momEndExcl = today.AddDate(0, 0, -6)
	} else if date == "lately30" {
		momStart = today.AddDate(0, 0, -59)
		momEndExcl = today.AddDate(0, 0, -29)
	} else if date == "month" {
		momStart = start.AddDate(0, -1, 0)
		momEndExcl = start
	} else if date == "year" {
		momStart = start.AddDate(-1, 0, 0)
		momEndExcl = start
	}

	return analyticsWindow{
		Start:        start,
		End:          endInclusive,
		EndExclusive: endExclusive,
		MomStart:     momStart,
		MomEndExcl:   momEndExcl,
		Buckets:      buckets,
		SQLFormat:    sqlFmt,
	}, nil
}

func parseSlashDate(s string, loc *time.Location) (time.Time, error) {
	s = strings.TrimSpace(s)
	for _, layout := range []string{"2006/01/02", "2006-01-02", "2006/1/2"} {
		if t, err := time.ParseInLocation(layout, s, loc); err == nil {
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc), nil
		}
	}
	return time.Time{}, fmt.Errorf("bad date %q", s)
}

func growthRate(current, previous int64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	rate := float64(current-previous) / float64(previous) * 100
	return math.Round(rate*100) / 100
}

func (h *Handler) countBrowse(c *gin.Context, scope dashboardScope, start, endExcl time.Time) int64 {
	var n int64
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_browse_history").
		Where("viewed_at >= ? AND viewed_at < ?", start, endExcl)
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = q.Count(&n)
	return n
}

func (h *Handler) countFavorite(c *gin.Context, scope dashboardScope, start, endExcl time.Time) int64 {
	var n int64
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_favorite AS f").
		Where("f.created_at >= ? AND f.created_at < ?", start, endExcl)
	if !scope.Full {
		q = q.Joins("JOIN qixi_crm_b_product_view AS p ON p.product_id = f.product_id").
			Where("p.store_id IN ?", scope.StoreIDs)
	}
	_ = q.Count(&n)
	return n
}

func (h *Handler) sumCart(c *gin.Context, scope dashboardScope, start, endExcl time.Time) int64 {
	type row struct {
		Total int64 `gorm:"column:total"`
	}
	var r row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_cart").
		Select("COALESCE(SUM(quantity),0) AS total").
		Where("created_at >= ? AND created_at < ?", start, endExcl)
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = q.Scan(&r)
	return r.Total
}

func (h *Handler) sumOrderQty(c *gin.Context, scope dashboardScope, start, endExcl time.Time, paidOnly bool) int64 {
	type row struct {
		Total int64 `gorm:"column:total"`
	}
	var r row
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select("COALESCE(SUM(total_quantity),0) AS total").
		Where("created_at >= ? AND created_at < ?", start, endExcl)
	if paidOnly {
		q = q.Where("status IN ?", paidOrderStatuses)
	}
	if !scope.Full {
		q = q.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = q.Scan(&r)
	return r.Total
}

func (h *Handler) countRefund(c *gin.Context, scope dashboardScope, start, endExcl time.Time) int64 {
	var n int64
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").
		Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").
		Where("r.created_at >= ? AND r.created_at < ?", start, endExcl)
	if !scope.Full {
		q = q.Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = q.Count(&n)
	return n
}
