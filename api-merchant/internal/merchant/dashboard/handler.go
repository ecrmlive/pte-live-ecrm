// Package dashboard serves the store console's current read-model overview.
package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/response"
	"gorm.io/gorm"
)

type Summary struct {
	ProductTotal     int64   `json:"product_total"`
	PaidOrder        int64   `json:"paid_order"`
	AvailableBalance float64 `json:"available_balance"`
	PendingRefund    int64   `json:"pending_refund"`
}

type Handler struct {
	merchantDB *gorm.DB
	businessDB *gorm.DB
}

func NewHandler(merchantDB, businessDB *gorm.DB) *Handler {
	return &Handler{merchantDB: merchantDB, businessDB: businessDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/dashboard/summary", h.GetSummary)
}

func (h *Handler) GetSummary(c *gin.Context) {
	storeID := middleware.StoreID(c)
	if storeID == 0 {
		response.Fail(c, http.StatusForbidden, "缺少店铺上下文")
		return
	}
	var out Summary
	queries := []struct {
		db    *gorm.DB
		table string
		where string
		out   *int64
	}{
		{h.merchantDB, "qixi_crm_m_product", "store_id = ?", &out.ProductTotal},
		{h.businessDB, "qixi_crm_b_order", "store_id = ? AND status IN ('paid','fulfilling','shipped','completed')", &out.PaidOrder},
		{h.businessDB, "qixi_crm_b_refund AS r JOIN qixi_crm_b_order AS o ON o.id = r.order_id", "o.store_id = ? AND r.status IN ('applied','merchant_handling','platform_intervene')", &out.PendingRefund},
	}
	for _, query := range queries {
		if err := query.db.WithContext(c.Request.Context()).Table(query.table).Where(query.where, storeID).Count(query.out).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "加载店铺经营概览失败")
			return
		}
	}
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_finance_ledger").
		Select("COALESCE(SUM(amount), 0)").Where("store_id = ?", storeID).Scan(&out.AvailableBalance).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载店铺经营概览失败")
		return
	}
	response.OK(c, out)
}
