// Package dashboard exposes the unified admin console's native read model.
package dashboard

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/response"
	"gorm.io/gorm"
)

// Summary is backed only by the current qixi_crm_a_/qixi_crm_b_ projections.
// It deliberately never falls back to the removed qixi_m_* compatibility tables.
type Summary struct {
	StoreTotal          int64            `json:"store_total"`
	OnSaleProduct       int64            `json:"on_sale_product"`
	PaidOrder           int64            `json:"paid_order"`
	PendingRefund       int64            `json:"pending_refund"`
	PendingStoreAudit   int64            `json:"pending_store_audit"`
	PendingProductAudit int64            `json:"pending_product_audit"`
	PendingDelivery     int64            `json:"pending_delivery"`
	PendingService      int64            `json:"pending_service"`
	NewUsers            Metric           `json:"new_users"`
	PageViews           Metric           `json:"page_views"`
	Visitors            Metric           `json:"visitors"`
	StoreCount          int64            `json:"store_count"`
	TodayOrderCount     int64            `json:"today_order_count"`
	TodayPayerCount     int64            `json:"today_payer_count"`
	TodayPaidAmount     float64          `json:"today_paid_amount"`
	StoreSalesRank      []StoreSalesRank `json:"store_sales_rank"`
}

// Metric contains only persisted facts. Today/yesterday/month are calculated in
// MySQL so every card uses the same database calendar as the order ledger.
type Metric struct {
	Today     int64 `json:"today"`
	Yesterday int64 `json:"yesterday"`
	Month     int64 `json:"month"`
}

type StoreSalesRank struct {
	StoreID       uint64  `json:"store_id"`
	StoreName     string  `json:"store_name"`
	FollowerCount int64   `json:"follower_count"`
	SaleCount     int64   `json:"sale_count"`
	SaleAmount    float64 `json:"sale_amount"`
}

type Handler struct {
	adminDB    *gorm.DB
	businessDB *gorm.DB
}

func NewHandler(adminDB, businessDB *gorm.DB) *Handler {
	return &Handler{adminDB: adminDB, businessDB: businessDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/dashboard/summary", h.GetSummary)
}

func (h *Handler) GetSummary(c *gin.Context) {
	var out Summary
	queries := []struct {
		db    *gorm.DB
		table string
		where string
		out   *int64
	}{
		{h.businessDB, "qixi_crm_b_store_view", "status = 1", &out.StoreTotal},
		{h.businessDB, "qixi_crm_b_product_view", "sale_status = 1", &out.OnSaleProduct},
		{h.businessDB, "qixi_crm_b_order", "status IN ('paid','fulfilling','shipped','completed')", &out.PaidOrder},
		{h.businessDB, "qixi_crm_b_refund", "status IN ('applied','merchant_handling','platform_intervene')", &out.PendingRefund},
		{h.adminDB, "qixi_crm_a_merchant_application", "status = 'pending'", &out.PendingStoreAudit},
		{h.adminDB, "qixi_crm_a_product_review", "status = 'pending'", &out.PendingProductAudit},
		{h.businessDB, "qixi_crm_b_order", "status = 'paid'", &out.PendingDelivery},
		{h.businessDB, "qixi_crm_b_customer_service_binding", "status = 'open'", &out.PendingService},
	}
	for _, query := range queries {
		if err := query.db.WithContext(c.Request.Context()).Table(query.table).Where(query.where).Count(query.out).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "加载经营概览失败")
			return
		}
	}
	if err := metricFor(h.businessDB, "qixi_crm_b_user", "created_at", "1 = 1", false, &out.NewUsers); err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载用户经营指标失败")
		return
	}
	if err := metricFor(h.businessDB, "qixi_crm_b_user_browse_history", "viewed_at", "1 = 1", false, &out.PageViews); err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载浏览经营指标失败")
		return
	}
	if err := metricFor(h.businessDB, "qixi_crm_b_user_browse_history", "viewed_at", "1 = 1", true, &out.Visitors); err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载访客经营指标失败")
		return
	}
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_store_view").Where("status = 1").Count(&out.StoreCount).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载店铺经营指标失败")
		return
	}
	if err := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order").
		Select(`
			COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN 1 ELSE 0 END), 0) AS today_order_count,
			COALESCE(COUNT(DISTINCT CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN user_id END), 0) AS today_payer_count,
			COALESCE(SUM(CASE WHEN status IN ('paid','fulfilling','shipped','completed') AND DATE(paid_at) = CURDATE() THEN pay_amount ELSE 0 END), 0) AS today_paid_amount`).
		Scan(&out).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载当日订单统计失败")
		return
	}
	if err := h.storeSalesRank(c, &out.StoreSalesRank); err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载店铺销售排行失败")
		return
	}
	response.OK(c, out)
}

func metricFor(db *gorm.DB, table, timeColumn, where string, distinctUser bool, out *Metric) error {
	if !distinctUser {
		selectSQL := fmt.Sprintf(`
			COALESCE(SUM(CASE WHEN DATE(%[1]s) = CURDATE() THEN 1 ELSE 0 END), 0) AS today,
			COALESCE(SUM(CASE WHEN DATE(%[1]s) = DATE_SUB(CURDATE(), INTERVAL 1 DAY) THEN 1 ELSE 0 END), 0) AS yesterday,
			COALESCE(SUM(CASE WHEN DATE_FORMAT(%[1]s, '%%Y-%%m') = DATE_FORMAT(CURDATE(), '%%Y-%%m') THEN 1 ELSE 0 END), 0) AS month`, timeColumn)
		return db.Table(table).Where(where).Select(selectSQL).Scan(out).Error
	}
	selectSQL := fmt.Sprintf(`
		COALESCE(COUNT(DISTINCT CASE WHEN DATE(%[1]s) = CURDATE() THEN user_id END), 0) AS today,
		COALESCE(COUNT(DISTINCT CASE WHEN DATE(%[1]s) = DATE_SUB(CURDATE(), INTERVAL 1 DAY) THEN user_id END), 0) AS yesterday,
		COALESCE(COUNT(DISTINCT CASE WHEN DATE_FORMAT(%[1]s, '%%Y-%%m') = DATE_FORMAT(CURDATE(), '%%Y-%%m') THEN user_id END), 0) AS month`, timeColumn)
	return db.Table(table).Where(where).Select(selectSQL).Scan(out).Error
}

func (h *Handler) storeSalesRank(c *gin.Context, out *[]StoreSalesRank) error {
	return h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order AS o").
		Select(`
			o.store_id,
			MAX(o.store_name_snapshot) AS store_name,
			COUNT(DISTINCT f.user_id) AS follower_count,
			COALESCE(SUM(o.total_quantity), 0) AS sale_count,
			COALESCE(SUM(o.pay_amount), 0) AS sale_amount`).
		Joins("LEFT JOIN qixi_crm_b_user_follow_store AS f ON f.store_id = o.store_id").
		Where("o.status IN ('paid','fulfilling','shipped','completed') AND o.paid_at >= DATE_FORMAT(CURDATE(), '%Y-%m-01')").
		Group("o.store_id").
		Order("sale_amount DESC, sale_count DESC, o.store_id ASC").
		Limit(10).
		Scan(out).Error
}
