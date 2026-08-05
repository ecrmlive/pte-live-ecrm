// Package points exposes the platform supervision plane for the business
// points-mall view. It never recalculates or mutates an existing points order:
// payment snapshots the points amount when the order is created.
package points

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const activityTypePoints = 20

type Handler struct{ businessDB, adminDB *gorm.DB }

func NewHandler(businessDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	manage := middleware.RequireAdminMenu(h.adminDB, "marketing.points.manage")
	r.GET("/points/products", access, manage, h.ListProducts)
	r.GET("/points/products/summary", access, manage, h.ProductSummary)
	r.PUT("/points/products/:id", access, manage, h.UpdateProduct)
	r.GET("/points/orders", access, manage, h.ListOrders)
}

type productRow struct {
	ProductID      uint64  `gorm:"column:product_id" json:"product_id"`
	MerchantID     uint64  `gorm:"column:merchant_id" json:"merchant_id"`
	StoreID        uint64  `gorm:"column:store_id" json:"store_id"`
	MerchantName   string  `gorm:"column:merchant_name" json:"merchant_name"`
	StoreName      string  `gorm:"column:store_name" json:"store_name"`
	Title          string  `gorm:"column:title" json:"title"`
	OriginalPrice  float64 `gorm:"column:original_price" json:"original_price"`
	PointsRequired int64   `gorm:"column:points_required" json:"points_required"`
	Stock          int     `gorm:"column:stock" json:"stock"`
	SaleStatus     int     `gorm:"column:sale_status" json:"sale_status"`
	Version        uint64  `gorm:"column:version" json:"version"`
}
type productSummary struct {
	Total  int64 `json:"total"`
	OnSale int64 `json:"on_sale"`
	Stock  int64 `json:"stock"`
}
type productUpdateInput struct {
	PointsRequired *int64 `json:"points_required"`
	Stock          *int   `json:"stock"`
	SaleStatus     *int   `json:"sale_status"`
	Version        uint64 `json:"version"`
}
type orderRow struct {
	ID            uint64 `gorm:"column:id" json:"id"`
	OrderNo       string `gorm:"column:order_no" json:"order_no"`
	UserID        uint64 `gorm:"column:user_id" json:"user_id"`
	PayStatus     string `gorm:"column:pay_status" json:"pay_status"`
	PointsAmount  int64  `gorm:"column:points_amount" json:"points_amount"`
	TotalQuantity int    `gorm:"column:total_quantity" json:"total_quantity"`
	CreatedAt     string `gorm:"column:created_at" json:"created_at"`
}

func (h *Handler) ListProducts(c *gin.Context) {
	page, limit := pagination(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view")
	if id, err := optionalPositive(c.Query("merchant_id")); err != nil {
		response.Fail(c, http.StatusBadRequest, "商户 ID 参数错误")
		return
	} else if id != 0 {
		q = q.Where("merchant_id = ?", id)
	}
	if status := strings.TrimSpace(c.Query("sale_status")); status != "" {
		v, err := strconv.Atoi(status)
		if err != nil || (v != 0 && v != 1) {
			response.Fail(c, http.StatusBadRequest, "上架状态参数错误")
			return
		}
		q = q.Where("sale_status = ?", v)
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("title LIKE ?", "%"+keyword+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "积分商品查询失败")
		return
	}
	rows := make([]productRow, 0)
	if err := q.Select("product_id,merchant_id,store_id,merchant_name,store_name,title,original_price,points_required,stock,sale_status,version").Order("updated_at DESC,product_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "积分商品查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) ProductSummary(c *gin.Context) {
	var out productSummary
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").Select("COUNT(*) AS total, COALESCE(SUM(sale_status = 1),0) AS on_sale, COALESCE(SUM(stock),0) AS stock").Scan(&out).Error; err != nil {
		fail(c, "积分商品统计失败")
		return
	}
	response.OK(c, out)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in productUpdateInput
	if id == 0 || c.ShouldBindJSON(&in) != nil || !validUpdate(&in) {
		response.Fail(c, http.StatusBadRequest, "积分商品更新参数错误")
		return
	}
	changes := map[string]any{}
	if in.PointsRequired != nil {
		changes["points_required"] = *in.PointsRequired
	}
	if in.Stock != nil {
		changes["stock"] = *in.Stock
	}
	if in.SaleStatus != nil {
		changes["sale_status"] = *in.SaleStatus
	}
	changes["version"] = gorm.Expr("version + 1")
	result := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").Where("product_id = ? AND version = ?", id, in.Version).Updates(changes)
	if result.Error != nil {
		fail(c, "积分商品更新失败")
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusConflict, "积分商品已变更，请刷新后重试")
		return
	}
	var row productRow
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").Where("product_id = ?", id).Scan(&row).Error; err != nil || row.ProductID == 0 {
		fail(c, "积分商品读取失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListOrders(c *gin.Context) {
	page, limit := pagination(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_group_order").Where("activity_type = ?", activityTypePoints)
	if raw := strings.TrimSpace(c.Query("pay_status")); raw != "" {
		if raw != "pending" && raw != "paid" && raw != "closed" {
			response.Fail(c, http.StatusBadRequest, "积分订单状态参数错误")
			return
		}
		q = q.Where("pay_status = ?", raw)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "积分订单查询失败")
		return
	}
	rows := make([]orderRow, 0)
	if err := q.Select("id,order_no,user_id,pay_status,points_amount,total_quantity,created_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "积分订单查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func validUpdate(in *productUpdateInput) bool {
	return in != nil && in.Version > 0 && ((in.PointsRequired != nil && *in.PointsRequired > 0) || (in.Stock != nil && *in.Stock >= 0) || (in.SaleStatus != nil && (*in.SaleStatus == 0 || *in.SaleStatus == 1)))
}
func optionalPositive(raw string) (uint64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, strconv.ErrSyntax
	}
	return id, nil
}
func pagination(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if p < 1 {
		p = 1
	}
	if l < 1 || l > 100 {
		l = 20
	}
	return p, l
}
func fail(c *gin.Context, message string) { response.Fail(c, http.StatusInternalServerError, message) }
