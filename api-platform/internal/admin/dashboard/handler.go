// Package dashboard exposes the unified admin console's native read model.
package dashboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/adminscope"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Summary is backed only by the current qixi_crm_a_/qixi_crm_b_ projections.
// It deliberately never falls back to the removed qixi_m_* compatibility tables.
type Summary struct {
	Scope                 string            `json:"scope"`
	StoreTotal            int64             `json:"store_total"`
	OnSaleProduct         int64             `json:"on_sale_product"`
	PaidOrder             int64             `json:"paid_order"`
	PendingRefund         int64             `json:"pending_refund"`
	PendingStoreAudit     int64             `json:"pending_store_audit"`
	PendingProductAudit   int64             `json:"pending_product_audit"`
	PendingDelivery       int64             `json:"pending_delivery"`
	PendingService        int64             `json:"pending_service"`
	PendingSpreadGift     int64             `json:"pending_spread_gift"`
	PendingWithdraw       int64             `json:"pending_withdraw"`
	PendingTransfer       int64             `json:"pending_transfer"`
	PendingCommunity      int64             `json:"pending_community"`
	PendingFeedback       int64             `json:"pending_feedback"`
	PendingIntegralShip   int64             `json:"pending_integral_ship"`
	NewUsers              Metric            `json:"new_users"`
	PageViews             Metric            `json:"page_views"`
	Visitors              Metric            `json:"visitors"`
	Stores                Metric            `json:"stores"`
	StoreCount            int64             `json:"store_count"`
	TodayOrderCount       int64             `json:"today_order_count"`
	TodayPayerCount       int64             `json:"today_payer_count"`
	TodayPaidAmount       float64           `json:"today_paid_amount"`
	TodayOrderHours       []HourAmountPoint `json:"today_order_hours"`
	OrderStats            OrderStatsBlock   `json:"order_stats"`
	UserTrend             []UserTrendPoint  `json:"user_trend"`
	DealFunnel            DealUserBlock     `json:"deal_funnel"`
	DealRatio             DealRatioBlock    `json:"deal_ratio"`
	StoreSalesRank        []StoreSalesRank  `json:"store_sales_rank"`
}

// Metric contains only persisted facts. Today/yesterday/month are calculated in
// MySQL so every card uses the same database calendar as the order ledger.
type Metric struct {
	Today     int64   `json:"today"`
	Yesterday int64   `json:"yesterday"`
	Month     int64   `json:"month"`
	WeekRatio float64 `json:"week_ratio"`
	LastWeek  int64   `json:"last_week"`
}

type StoreSalesRank struct {
	StoreID       uint64  `json:"store_id"`
	StoreName     string  `json:"store_name"`
	StoreImage    string  `json:"store_image"`
	FollowerCount int64   `json:"follower_count"`
	SaleCount     int64   `json:"sale_count"`
	SaleAmount    float64 `json:"sale_amount"`
}

// todayOrderStats is intentionally kept separate from Summary. GORM interprets
// Summary's nested Metric fields as model relations when scanning a partial
// SELECT, which turns a successful aggregation query into a 500 response.
type todayOrderStats struct {
	TodayOrderCount int64   `gorm:"column:today_order_count"`
	TodayPayerCount int64   `gorm:"column:today_payer_count"`
	TodayPaidAmount float64 `gorm:"column:today_paid_amount"`
}

type Handler struct {
	adminDB    *gorm.DB
	businessDB *gorm.DB
	merchantDB *gorm.DB
}

func NewHandler(adminDB, businessDB, merchantDB *gorm.DB) *Handler {
	return &Handler{adminDB: adminDB, businessDB: businessDB, merchantDB: merchantDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/dashboard/summary", h.GetSummary)
	r.GET("/dashboard/data-screen", h.GetDataScreen)
	h.RegisterConsole(r)
	h.RegisterProductAnalytics(r)
	h.RegisterOrderAnalytics(r)
	h.RegisterUserAnalytics(r)
}

func (h *Handler) GetSummary(c *gin.Context) {
	scope, err := h.resolveScope(c)
	if err != nil {
		if errors.Is(err, adminscope.ErrNotConfigured) {
			response.Fail(c, http.StatusForbidden, "未配置首页数据范围")
		} else {
			response.Fail(c, http.StatusInternalServerError, "加载首页数据范围失败")
		}
		return
	}
	var out Summary
	out.Scope = "store"
	if scope.Full {
		out.Scope = "all"
	}
	queries := []struct {
		db    *gorm.DB
		table string
		where string
		out   *int64
	}{
		{h.businessDB, "qixi_crm_b_store_view", "status = 1", &out.StoreTotal},
		{h.businessDB, "qixi_crm_b_product_view", "sale_status = 1", &out.OnSaleProduct},
		{h.businessDB, "qixi_crm_b_order", "status IN ('paid','fulfilling','shipped','completed')", &out.PaidOrder},
		{h.businessDB, "qixi_crm_b_order", "status = 'paid'", &out.PendingDelivery},
		{h.businessDB, "qixi_crm_b_customer_service_binding", "status = 'open'", &out.PendingService},
	}
	for _, query := range queries {
		q := query.db.WithContext(c.Request.Context()).Table(query.table).Where(query.where)
		if !scope.Full {
			q = q.Where("store_id IN ?", scope.StoreIDs)
		}
		if err := q.Count(query.out).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "加载经营概览失败")
			return
		}
	}
	refundQuery := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").Where("r.status IN ('applied','merchant_handling','awaiting_return','awaiting_receipt','platform_intervene')")
	if !scope.Full {
		refundQuery = refundQuery.Where("o.store_id IN ?", scope.StoreIDs)
	}
	if err := refundQuery.Count(&out.PendingRefund).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载经营概览失败")
		return
	}
	if scope.Full {
		if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_application").Where("status = 'pending'").Count(&out.PendingStoreAudit).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "加载经营概览失败")
			return
		}
	}
	productAuditQuery := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_product_review").Where("status = 'pending'")
	if !scope.Full {
		productAuditQuery = productAuditQuery.Where("store_id IN ?", scope.StoreIDs)
	}
	if err := productAuditQuery.Count(&out.PendingProductAudit).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载经营概览失败")
		return
	}
	if scope.Full {
		err = metricFor(h.businessDB, "qixi_crm_b_user", "created_at", "1 = 1", false, &out.NewUsers)
	} else {
		err = metricFor(h.businessDB, "qixi_crm_b_order", "created_at", "store_id IN ?", true, &out.NewUsers, scope.StoreIDs)
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载用户经营指标失败")
		return
	}
	metricArgs := []any{}
	metricWhere := "1 = 1"
	if !scope.Full {
		metricWhere, metricArgs = "store_id IN ?", []any{scope.StoreIDs}
	}
	if err := metricFor(h.businessDB, "qixi_crm_b_user_browse_history", "viewed_at", metricWhere, false, &out.PageViews, metricArgs...); err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载浏览经营指标失败")
		return
	}
	if err := metricFor(h.businessDB, "qixi_crm_b_user_browse_history", "viewed_at", metricWhere, true, &out.Visitors, metricArgs...); err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载访客经营指标失败")
		return
	}
	storeQuery := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_store_view").Where("status = 1")
	if !scope.Full {
		storeQuery = storeQuery.Where("store_id IN ?", scope.StoreIDs)
	}
	if err := storeQuery.Count(&out.StoreCount).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载店铺经营指标失败")
		return
	}
	var today todayOrderStats
	todayQuery := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order").
		Select(`
			COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN 1 ELSE 0 END), 0) AS today_order_count,
			COALESCE(COUNT(DISTINCT CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN user_id END), 0) AS today_payer_count,
			COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN pay_amount ELSE 0 END), 0) AS today_paid_amount`)
	if !scope.Full {
		todayQuery = todayQuery.Where("store_id IN ?", scope.StoreIDs)
	}
	if err := todayQuery.Scan(&today).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载当日订单统计失败")
		return
	}
	out.TodayOrderCount = today.TodayOrderCount
	out.TodayPayerCount = today.TodayPayerCount
	out.TodayPaidAmount = today.TodayPaidAmount
	h.fillExtraTodos(c, scope, &out)
	if err := metricFor(h.merchantDB, "qixi_crm_m_store", "created_at", "1 = 1", false, &out.Stores); err != nil {
		out.Stores = Metric{}
	}
	out.Stores.Month = out.StoreCount
	if hours, err := h.todayOrderHours(c, scope); err == nil {
		out.TodayOrderHours = hours
	} else {
		out.TodayOrderHours = []HourAmountPoint{}
	}
	if stats, err := h.orderStatsBlock(c, scope); err == nil {
		out.OrderStats = stats
	}
	if trend, err := h.userTrend(c, scope, "30d"); err == nil {
		out.UserTrend = trend
	} else {
		out.UserTrend = []UserTrendPoint{}
	}
	if funnel, ratio, err := h.dealBlocks(c, scope, "month"); err == nil {
		out.DealFunnel = funnel
		out.DealRatio = ratio
	}
	if err := h.storeSalesRank(c, scope, &out.StoreSalesRank); err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载店铺销售排行失败")
		return
	}
	// JSON must remain an array for an empty ranking. A null value crashes the
	// Vben empty-state branch (`store_sales_rank.length`) and turns a valid
	// zero-data dashboard into a client-side render error.
	if out.StoreSalesRank == nil {
		out.StoreSalesRank = []StoreSalesRank{}
	}
	if out.TodayOrderHours == nil {
		out.TodayOrderHours = []HourAmountPoint{}
	}
	if out.UserTrend == nil {
		out.UserTrend = []UserTrendPoint{}
	}
	response.OK(c, out)
}

func (h *Handler) fillExtraTodos(c *gin.Context, scope dashboardScope, out *Summary) {
	safeCount(h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_withdrawal_application").
		Where("status IN ('applied','reviewing')"), &out.PendingWithdraw)
	safeCount(h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_settlement_view").
		Where("status = 'withdraw_applied'"), &out.PendingTransfer)
	safeCount(h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_social_post").
		Where("is_del = 0 AND status = 0"), &out.PendingCommunity)
	safeCount(h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_feedback").
		Where("deleted_at IS NULL AND status = 'pending'"), &out.PendingFeedback)
	giftQ := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product AS p").
		Joins("LEFT JOIN qixi_crm_m_product_recycle_bin AS rb ON rb.product_id = p.id").
		Where("rb.product_id IS NULL AND p.is_gift_bag = 1 AND p.status IN ?", []string{"pending_review", "draft"})
	if !scope.Full {
		giftQ = giftQ.Where("p.store_id IN ?", scope.StoreIDs)
	}
	safeCount(giftQ, &out.PendingSpreadGift)
	// 积分订单发货：当前 schema 无独立待审表，固定 0，结构占位。
	out.PendingIntegralShip = 0
}

func metricFor(db *gorm.DB, table, timeColumn, where string, distinctUser bool, out *Metric, args ...any) error {
	type rawMetric struct {
		Today     int64 `gorm:"column:today"`
		Yesterday int64 `gorm:"column:yesterday"`
		Month     int64 `gorm:"column:month"`
		LastWeek  int64 `gorm:"column:last_week"`
	}
	var raw rawMetric
	var selectSQL string
	if !distinctUser {
		selectSQL = fmt.Sprintf(`
			COALESCE(SUM(CASE WHEN DATE(%[1]s) = CURDATE() THEN 1 ELSE 0 END), 0) AS today,
			COALESCE(SUM(CASE WHEN DATE(%[1]s) = DATE_SUB(CURDATE(), INTERVAL 1 DAY) THEN 1 ELSE 0 END), 0) AS yesterday,
			COALESCE(SUM(CASE WHEN DATE_FORMAT(%[1]s, '%%Y-%%m') = DATE_FORMAT(CURDATE(), '%%Y-%%m') THEN 1 ELSE 0 END), 0) AS month,
			COALESCE(SUM(CASE WHEN DATE(%[1]s) = DATE_SUB(CURDATE(), INTERVAL 7 DAY) THEN 1 ELSE 0 END), 0) AS last_week`, timeColumn)
	} else {
		selectSQL = fmt.Sprintf(`
			COALESCE(COUNT(DISTINCT CASE WHEN DATE(%[1]s) = CURDATE() THEN user_id END), 0) AS today,
			COALESCE(COUNT(DISTINCT CASE WHEN DATE(%[1]s) = DATE_SUB(CURDATE(), INTERVAL 1 DAY) THEN user_id END), 0) AS yesterday,
			COALESCE(COUNT(DISTINCT CASE WHEN DATE_FORMAT(%[1]s, '%%Y-%%m') = DATE_FORMAT(CURDATE(), '%%Y-%%m') THEN user_id END), 0) AS month,
			COALESCE(COUNT(DISTINCT CASE WHEN DATE(%[1]s) = DATE_SUB(CURDATE(), INTERVAL 7 DAY) THEN user_id END), 0) AS last_week`, timeColumn)
	}
	if err := db.Table(table).Where(where, args...).Select(selectSQL).Scan(&raw).Error; err != nil {
		return err
	}
	out.Today = raw.Today
	out.Yesterday = raw.Yesterday
	out.Month = raw.Month
	out.LastWeek = raw.LastWeek
	out.WeekRatio = growthRatio(raw.Today, raw.LastWeek)
	return nil
}

func (h *Handler) storeSalesRank(c *gin.Context, scope dashboardScope, out *[]StoreSalesRank) error {
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
		Where("o.status IN ('paid','fulfilling','shipped','completed') AND o.paid_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')")
	if !scope.Full {
		q = q.Where("o.store_id IN ?", scope.StoreIDs)
	}
	return q.
		Group("o.store_id").
		Order("sale_amount DESC, sale_count DESC, o.store_id ASC").
		Limit(10).
		Scan(out).Error
}

type dashboardScope struct {
	Full     bool
	StoreIDs []uint64
}

func (h *Handler) resolveScope(c *gin.Context) (dashboardScope, error) {
	claims := middleware.ClaimsFrom(c)
	if claims == nil {
		return dashboardScope{}, adminscope.ErrNotConfigured
	}
	if hasRole(claims.Roles, "platform") {
		return dashboardScope{Full: true}, nil
	}
	if hasRole(claims.Roles, "merchant") || hasRole(claims.Roles, "region") {
		merchantScope, err := adminscope.ResolveMerchantScope(c.Request.Context(), h.adminDB, claims)
		if err != nil {
			return dashboardScope{}, err
		}
		merchantIDs := append([]uint64{}, merchantScope.MerchantIDs...)
		if len(merchantScope.RegionIDs) > 0 {
			var regional []struct{ ID uint64 }
			if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_merchant").Select("id").Where("region_id IN ?", merchantScope.RegionIDs).Find(&regional).Error; err != nil {
				return dashboardScope{}, err
			}
			for _, row := range regional {
				merchantIDs = appendUnique(merchantIDs, row.ID)
			}
		}
		if len(merchantIDs) == 0 {
			return dashboardScope{}, adminscope.ErrNotConfigured
		}
		var stores []struct{ ID uint64 }
		if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_store").Select("id").Where("merchant_id IN ?", merchantIDs).Find(&stores).Error; err != nil {
			return dashboardScope{}, err
		}
		storeIDs := make([]uint64, 0, len(stores))
		for _, row := range stores {
			storeIDs = appendUnique(storeIDs, row.ID)
		}
		return dashboardScope{StoreIDs: storeIDs}, nil
	}
	if !hasRole(claims.Roles, "customer_service") {
		return dashboardScope{}, adminscope.ErrNotConfigured
	}
	var rows []struct {
		ScopeValue json.RawMessage `gorm:"column:scope_value"`
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_data_scope").Select("scope_value").Where("admin_user_id = ? AND scope_type = ?", claims.AdminID, "service_queue").Find(&rows).Error; err != nil {
		return dashboardScope{}, err
	}
	storeIDs := []uint64{}
	for _, row := range rows {
		var value struct {
			StoreIDs []uint64 `json:"store_ids"`
		}
		if json.Unmarshal(row.ScopeValue, &value) != nil {
			continue
		}
		for _, id := range value.StoreIDs {
			storeIDs = appendUnique(storeIDs, id)
		}
	}
	if len(storeIDs) == 0 {
		return dashboardScope{}, adminscope.ErrNotConfigured
	}
	return dashboardScope{StoreIDs: storeIDs}, nil
}

func appendUnique(values []uint64, value uint64) []uint64 {
	if value == 0 {
		return values
	}
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func hasRole(roles []string, expected string) bool {
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
}
