package dashboard

import (
	"fmt"
	"net/http"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// DataScreen is the platform「数据大屏」payload，字段对齐 CRMEB data_screen/*。
type DataScreen struct {
	Config               DataScreenConfig            `json:"config"`
	TodayPayCountNumber  DataScreenTodayNumbers      `json:"today_pay_count_number"`
	TodayPayNewOld       DataScreenNewOld            `json:"today_pay_new_old"`
	TodayPayNumber       DataScreenPaymentAmount     `json:"today_pay_number"`
	CityRanking          []DataScreenCityRank        `json:"city_ranking"`
	MonthPayCount        []DataScreenMonthPoint      `json:"month_pay_count"`
	TodayPayCount        []DataScreenHourPoint       `json:"today_pay_count"`
	TodayPayInfo         []DataScreenOrderInfo       `json:"today_pay_info"`
	TodayPayMerchantRank DataScreenMerchantRankBoard `json:"today_pay_merchant_rank"`
	PayProductRank       []DataScreenProductRank     `json:"pay_product_rank"`
}

type DataScreenConfig struct {
	Title string `json:"data_screen_title"`
}

type DataScreenTodayNumbers struct {
	VisitNum          int64 `json:"visit_num"`
	VisitUserNum      int64 `json:"visit_user_num"`
	TodayPayUserFirst int64 `json:"today_pay_user_first"`
	TodayPayNumber    int64 `json:"today_pay_number"`
}

type DataScreenNewOld struct {
	NewCount int64 `json:"new_count"`
	OldCount int64 `json:"old_count"`
}

// DataScreenPaymentAmount mirrors CRMEB today_pay_number.
type DataScreenPaymentAmount struct {
	Count   int64   `json:"count"`
	Number  float64 `json:"number"`
	OrderID uint64  `json:"order_id"`
	Paid    int     `json:"paid"`
}

type DataScreenCityRank struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
	Code  string `json:"code"`
}

type DataScreenMonthPoint struct {
	Day      string  `json:"day"`
	TotalSum float64 `json:"total_sum"`
}

type DataScreenHourPoint struct {
	Hours      string `json:"hours"`
	UserCount  int64  `json:"user_count"`
	OrderCount int64  `json:"order_count"`
}

type DataScreenOrderInfo struct {
	Number        float64           `json:"number"`
	PaymentMethod string            `json:"payment_method"`
	PayTime       string            `json:"paytime"`
	Product       DataScreenProduct `json:"product"`
	Store         DataScreenStore   `json:"store"`
}

type DataScreenProduct struct {
	Image       string `json:"image"`
	ProductName string `json:"product_name"`
}

type DataScreenStore struct {
	Image     string `json:"image"`
	StoreName string `json:"store_name"`
}

type DataScreenMerchantRankBoard struct {
	Data []DataScreenMerchantRank `json:"data"`
	Type string                   `json:"type"`
}

type DataScreenMerchantRank struct {
	Count  int64           `json:"count"`
	Number float64         `json:"number"`
	Store  DataScreenStore `json:"store"`
}

// DataScreenProductRank aligns with CRMEB pay_product_rank rows.
type DataScreenProductRank struct {
	Count   int64             `json:"count"`
	Number  float64           `json:"number"`
	Product DataScreenProduct `json:"product"`
}

func (h *Handler) GetDataScreen(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置数据大屏范围")
		return
	}
	out := DataScreen{
		Config:        DataScreenConfig{Title: "数据大屏"},
		CityRanking:   []DataScreenCityRank{},
		MonthPayCount: []DataScreenMonthPoint{},
		TodayPayCount: []DataScreenHourPoint{},
		TodayPayInfo:  []DataScreenOrderInfo{},
		TodayPayMerchantRank: DataScreenMerchantRankBoard{
			Data: []DataScreenMerchantRank{},
			Type: "元",
		},
		PayProductRank: []DataScreenProductRank{},
	}

	metricWhere := "1 = 1"
	metricArgs := []any{}
	if !scope.Full {
		metricWhere, metricArgs = "store_id IN ?", []any{scope.StoreIDs}
	}
	var pageViews, visitors Metric
	_ = metricFor(h.businessDB, "qixi_crm_b_user_browse_history", "viewed_at", metricWhere, false, &pageViews, metricArgs...)
	_ = metricFor(h.businessDB, "qixi_crm_b_user_browse_history", "viewed_at", metricWhere, true, &visitors, metricArgs...)
	out.TodayPayCountNumber.VisitNum = pageViews.Today
	out.TodayPayCountNumber.VisitUserNum = visitors.Today

	var today todayOrderStats
	todayQuery := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").Select(`
		COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN 1 ELSE 0 END), 0) AS today_order_count,
		COALESCE(COUNT(DISTINCT CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN user_id END), 0) AS today_payer_count,
		COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN pay_amount ELSE 0 END), 0) AS today_paid_amount`)
	if !scope.Full {
		todayQuery = todayQuery.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = todayQuery.Scan(&today).Error
	out.TodayPayCountNumber.TodayPayNumber = today.TodayOrderCount
	out.TodayPayNumber = DataScreenPaymentAmount{
		Count:  today.TodayOrderCount,
		Number: today.TodayPaidAmount,
		Paid:   1,
	}

	var newUsers Metric
	if scope.Full {
		_ = metricFor(h.businessDB, "qixi_crm_b_user", "created_at", "1 = 1", false, &newUsers)
	}
	out.TodayPayCountNumber.TodayPayUserFirst = newUsers.Today
	out.TodayPayNewOld.NewCount = newUsers.Today
	if today.TodayPayerCount > newUsers.Today {
		out.TodayPayNewOld.OldCount = today.TodayPayerCount - newUsers.Today
	}

	type monthRow struct {
		Day      string  `gorm:"column:day"`
		TotalSum float64 `gorm:"column:total_sum"`
	}
	var monthRows []monthRow
	mq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select("DATE_FORMAT(paid_at, '%d') AS day, COALESCE(SUM(pay_amount),0) AS total_sum").
		Where("status IN ('paid','fulfilling','shipped','completed') AND paid_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')")
	if !scope.Full {
		mq = mq.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = mq.Group("DATE_FORMAT(paid_at, '%d')").Order("day ASC").Scan(&monthRows).Error
	dayMap := map[string]float64{}
	for _, row := range monthRows {
		dayMap[row.Day] = row.TotalSum
	}
	now := time.Now()
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()
	for d := 1; d <= daysInMonth; d++ {
		key := fmt.Sprintf("%02d", d)
		out.MonthPayCount = append(out.MonthPayCount, DataScreenMonthPoint{Day: key, TotalSum: dayMap[key]})
	}

	type hourRow struct {
		HourBucket int   `gorm:"column:hour_bucket"`
		UserCount  int64 `gorm:"column:user_count"`
		OrderCount int64 `gorm:"column:order_count"`
	}
	var hourRows []hourRow
	hq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select("FLOOR(HOUR(paid_at)/2)*2 AS hour_bucket, COUNT(*) AS order_count, COUNT(DISTINCT user_id) AS user_count").
		Where("status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE()")
	if !scope.Full {
		hq = hq.Where("store_id IN ?", scope.StoreIDs)
	}
	_ = hq.Group("FLOOR(HOUR(paid_at)/2)*2").Scan(&hourRows).Error
	hourMap := map[int]hourRow{}
	for _, row := range hourRows {
		hourMap[row.HourBucket] = row
	}
	for hour := 0; hour < 24; hour += 2 {
		row := hourMap[hour]
		out.TodayPayCount = append(out.TodayPayCount, DataScreenHourPoint{
			Hours:      fmt.Sprintf("%02d~%02d", hour, hour+1),
			UserCount:  row.UserCount,
			OrderCount: row.OrderCount,
		})
	}

	type orderRow struct {
		PaymentMethod string    `gorm:"column:payment_method"`
		ProductImage  string    `gorm:"column:product_image"`
		ProductName   string    `gorm:"column:product_name"`
		StoreImage    string    `gorm:"column:store_image"`
		StoreName     string    `gorm:"column:store_name"`
		PayTime       time.Time `gorm:"column:paid_at"`
		Amount        float64   `gorm:"column:pay_amount"`
	}
	var orders []orderRow
	oq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order AS o").
		Select(`o.paid_at, o.pay_amount,
			COALESCE(NULLIF(TRIM(o.store_name_snapshot), ''), '—') AS store_name,
			COALESCE((
				SELECT NULLIF(TRIM(p.cover_url), '') FROM qixi_crm_b_product_view AS p
				WHERE p.store_id = o.store_id ORDER BY p.product_id ASC LIMIT 1
			), '') AS store_image,
			CASE COALESCE(g.pay_channel, 'mock')
				WHEN 'wechat' THEN '微信支付'
				WHEN 'alipay' THEN '支付宝'
				WHEN 'balance' THEN '余额支付'
				ELSE '线上支付'
			END AS payment_method,
			COALESCE(NULLIF(TRIM((
				SELECT i.title_snapshot FROM qixi_crm_b_order_item AS i
				WHERE i.order_id = o.id ORDER BY i.id ASC LIMIT 1
			)), ''), NULLIF(TRIM(o.store_name_snapshot), ''), '—') AS product_name,
			COALESCE((
				SELECT i.cover_url_snapshot FROM qixi_crm_b_order_item AS i
				WHERE i.order_id = o.id ORDER BY i.id ASC LIMIT 1
			), '') AS product_image`).
		Joins("LEFT JOIN qixi_crm_b_group_order AS g ON g.id = o.group_order_id").
		Where("o.status IN ('paid','fulfilling','shipped','completed') AND DATE(o.paid_at) = CURDATE()")
	if !scope.Full {
		oq = oq.Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = oq.Order("o.paid_at DESC").Limit(20).Scan(&orders).Error
	for _, row := range orders {
		out.TodayPayInfo = append(out.TodayPayInfo, DataScreenOrderInfo{
			Number:        row.Amount,
			PaymentMethod: row.PaymentMethod,
			PayTime:       row.PayTime.Format("2006-01-02 15:04:05"),
			Product:       DataScreenProduct{Image: row.ProductImage, ProductName: row.ProductName},
			Store:         DataScreenStore{Image: row.StoreImage, StoreName: row.StoreName},
		})
	}

	type storeRankRow struct {
		SaleAmount float64 `gorm:"column:sale_amount"`
		SaleCount  int64   `gorm:"column:sale_count"`
		StoreImage string  `gorm:"column:store_image"`
		StoreName  string  `gorm:"column:store_name"`
	}
	var ranks []storeRankRow
	rankQuery := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order AS o").
		Select(`
			MAX(o.store_name_snapshot) AS store_name,
			COALESCE((
				SELECT NULLIF(TRIM(p.cover_url), '') FROM qixi_crm_b_product_view AS p
				WHERE p.store_id = o.store_id ORDER BY p.product_id ASC LIMIT 1
			), '') AS store_image,
			COALESCE(SUM(o.total_quantity), 0) AS sale_count,
			COALESCE(SUM(o.pay_amount), 0) AS sale_amount`).
		Where("o.status IN ('paid','fulfilling','shipped','completed') AND DATE(o.paid_at) = CURDATE()")
	if !scope.Full {
		rankQuery = rankQuery.Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = rankQuery.
		Group("o.store_id").
		Order("sale_amount DESC, sale_count DESC, o.store_id ASC").
		Limit(20).
		Scan(&ranks).Error
	for _, row := range ranks {
		out.TodayPayMerchantRank.Data = append(out.TodayPayMerchantRank.Data, DataScreenMerchantRank{
			Count:  row.SaleCount,
			Number: row.SaleAmount,
			Store:  DataScreenStore{Image: row.StoreImage, StoreName: row.StoreName},
		})
	}

	type productRankRow struct {
		ProductName string  `gorm:"column:product_name"`
		StoreName   string  `gorm:"column:store_name"`
		Image       string  `gorm:"column:image"`
		SaleCount   int64   `gorm:"column:sale_count"`
		SaleAmount  float64 `gorm:"column:sale_amount"`
	}
	var productRanks []productRankRow
	pq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order_item AS i").
		Select(`
			i.product_id,
			MAX(COALESCE(NULLIF(TRIM(i.title_snapshot), ''), CONCAT('商品#', i.product_id))) AS product_name,
			MAX(COALESCE(NULLIF(TRIM(o.store_name_snapshot), ''), '—')) AS store_name,
			MAX(COALESCE(i.cover_url_snapshot, '')) AS image,
			COALESCE(SUM(i.quantity), 0) AS sale_count,
			COALESCE(SUM(i.unit_price * i.quantity), 0) AS sale_amount`).
		Joins("INNER JOIN qixi_crm_b_order AS o ON o.id = i.order_id").
		Where("o.status IN ('paid','fulfilling','shipped','completed') AND o.paid_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')")
	if !scope.Full {
		pq = pq.Where("o.store_id IN ?", scope.StoreIDs)
	}
	_ = pq.Group("i.product_id").
		Order("sale_amount DESC, sale_count DESC, i.product_id ASC").
		Limit(20).
		Scan(&productRanks).Error
	for _, row := range productRanks {
		out.PayProductRank = append(out.PayProductRank, DataScreenProductRank{
			Count:   row.SaleCount,
			Number:  row.SaleAmount,
			Product: DataScreenProduct{Image: row.Image, ProductName: row.ProductName},
		})
	}

	// 收货省从 recipient_snapshot JSON 提取；保留原始省级名称，以匹配 china.json。
	type cityRow struct {
		Name  string `gorm:"column:name"`
		Value int64  `gorm:"column:value"`
	}
	var cities []cityRow
	cq := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select("COALESCE(NULLIF(TRIM(JSON_UNQUOTE(JSON_EXTRACT(recipient_snapshot, '$.province'))), ''), '未知') AS name, COUNT(*) AS value").
		Where("status IN ('paid','fulfilling','shipped','completed') AND paid_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')")
	if !scope.Full {
		cq = cq.Where("store_id IN ?", scope.StoreIDs)
	}
	if err := cq.Group("name").Order("value DESC").Limit(40).Scan(&cities).Error; err == nil {
		for _, row := range cities {
			if row.Name == "未知" || row.Name == "null" {
				continue
			}
			out.CityRanking = append(out.CityRanking, DataScreenCityRank{
				Name:  row.Name,
				Value: row.Value,
			})
		}
	}

	response.OK(c, out)
}
