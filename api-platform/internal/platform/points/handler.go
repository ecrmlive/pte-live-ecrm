// Package points exposes the platform supervision plane for the business
// points-mall view. It never recalculates or mutates an existing points order:
// payment snapshots the points amount when the order is created.
package points

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

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
	r.POST("/points/products", access, manage, h.CreateProduct)
	r.POST("/points/products/quick", access, manage, h.QuickAddProduct)
	r.GET("/points/products/:id", access, manage, h.GetProduct)
	r.PUT("/points/products/:id", access, manage, h.UpdateProduct)
	r.PUT("/points/products/:id/status", access, manage, h.SwitchStatus)
	r.POST("/points/products/:id/copy", access, manage, h.CopyProduct)
	r.DELETE("/points/products/:id", access, manage, h.DeleteProduct)
	r.GET("/points/products/:id/exchanges", access, manage, h.ListExchanges)
	r.GET("/points/orders", access, manage, h.ListOrders)
}

type productRow struct {
	ProductID       uint64  `gorm:"column:product_id" json:"product_id"`
	MerchantID      uint64  `gorm:"column:merchant_id" json:"merchant_id"`
	StoreID         uint64  `gorm:"column:store_id" json:"store_id"`
	MerchantName    string  `gorm:"column:merchant_name" json:"merchant_name"`
	StoreName       string  `gorm:"column:store_name" json:"store_name"`
	Title           string  `gorm:"column:title" json:"title"`
	CateID          uint64  `gorm:"column:cate_id" json:"cate_id"`
	CoverURL        string  `gorm:"column:cover_url" json:"cover_url"`
	OriginalPrice   float64 `gorm:"column:original_price" json:"original_price"`
	PointsRequired  int64   `gorm:"column:points_required" json:"points_required"`
	Stock           int     `gorm:"column:stock" json:"stock"`
	Sales           int     `gorm:"column:sales" json:"sales"`
	Sort            int     `gorm:"column:sort" json:"sort"`
	SaleStatus      int     `gorm:"column:sale_status" json:"sale_status"`
	SourceProductID uint64  `gorm:"column:source_product_id" json:"source_product_id"`
	Version         uint64  `gorm:"column:version" json:"version"`
	CreateTime      string  `gorm:"column:create_time" json:"create_time"`
	UpdatedAt       string  `gorm:"column:updated_at" json:"updated_at"`
}

type productSummary struct {
	Total  int64 `json:"total"`
	OnSale int64 `json:"on_sale"`
	Stock  int64 `json:"stock"`
}

type productSaveInput struct {
	Title           string   `json:"title"`
	CoverURL        string   `json:"cover_url"`
	CateID          *uint64  `json:"cate_id"`
	OriginalPrice   *float64 `json:"original_price"`
	PointsRequired  *int64   `json:"points_required"`
	Stock           *int     `json:"stock"`
	Sales           *int     `json:"sales"`
	Sort            *int     `json:"sort"`
	SaleStatus      *int     `json:"sale_status"`
	MerchantID      *uint64  `json:"merchant_id"`
	StoreID         *uint64  `json:"store_id"`
	MerchantName    string   `json:"merchant_name"`
	StoreName       string   `json:"store_name"`
	SourceProductID *uint64  `json:"source_product_id"`
	Version         uint64   `json:"version"`
}

type quickAddInput struct {
	SourceProductID uint64   `json:"source_product_id"`
	Title           string   `json:"title"`
	CoverURL        string   `json:"cover_url"`
	CateID          uint64   `json:"cate_id"`
	OriginalPrice   *float64 `json:"original_price"`
	PointsRequired  *int64   `json:"points_required"`
	Stock           *int     `json:"stock"`
	Sort            *int     `json:"sort"`
	SaleStatus      *int     `json:"sale_status"`
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

type exchangeRow struct {
	OrderID       uint64 `gorm:"column:order_id" json:"order_id"`
	OrderNo       string `gorm:"column:order_no" json:"order_no"`
	UserID        uint64 `gorm:"column:user_id" json:"user_id"`
	PayStatus     string `gorm:"column:pay_status" json:"pay_status"`
	PointsAmount  int64  `gorm:"column:points_amount" json:"points_amount"`
	Quantity      int    `gorm:"column:quantity" json:"quantity"`
	TitleSnapshot string `gorm:"column:title_snapshot" json:"title_snapshot"`
	CreatedAt     string `gorm:"column:created_at" json:"created_at"`
}

const productSelect = `product_id,merchant_id,store_id,merchant_name,store_name,title,cate_id,cover_url,original_price,points_required,stock,sales,sort,sale_status,source_product_id,version,create_time,updated_at`

func (h *Handler) ListProducts(c *gin.Context) {
	page, limit := pagination(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").Where("is_del = 0")
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
		if id, err := strconv.ParseUint(keyword, 10, 64); err == nil && id > 0 {
			q = q.Where("(product_id = ? OR title LIKE ?)", id, "%"+keyword+"%")
		} else {
			q = q.Where("title LIKE ?", "%"+keyword+"%")
		}
	}
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		q = q.Where("create_time >= ?", from+" 00:00:00")
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		q = q.Where("create_time <= ?", to+" 23:59:59")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "积分商品查询失败")
		return
	}
	rows := make([]productRow, 0)
	if err := q.Select(productSelect).Order("sort DESC, product_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "积分商品查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) ProductSummary(c *gin.Context) {
	var out productSummary
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").
		Where("is_del = 0").
		Select("COUNT(*) AS total, COALESCE(SUM(sale_status = 1),0) AS on_sale, COALESCE(SUM(stock),0) AS stock").
		Scan(&out).Error; err != nil {
		fail(c, "积分商品统计失败")
		return
	}
	response.OK(c, out)
}

func (h *Handler) GetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "积分商品参数错误")
		return
	}
	row, err := h.loadProduct(c, id)
	if err != nil || row.ProductID == 0 {
		response.Fail(c, http.StatusNotFound, "积分商品不存在")
		return
	}
	response.OK(c, row)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var in productSaveInput
	if c.ShouldBindJSON(&in) != nil || !validCreate(&in) {
		response.Fail(c, http.StatusBadRequest, "积分商品参数错误")
		return
	}
	id, err := h.nextProductID(c)
	if err != nil {
		fail(c, "积分商品添加失败")
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	row := map[string]any{
		"product_id":        id,
		"merchant_id":       derefU64(in.MerchantID, 0),
		"store_id":          derefU64(in.StoreID, 0),
		"merchant_name":     strings.TrimSpace(in.MerchantName),
		"store_name":        strings.TrimSpace(in.StoreName),
		"title":             strings.TrimSpace(in.Title),
		"cate_id":           derefU64(in.CateID, 0),
		"cover_url":         strings.TrimSpace(in.CoverURL),
		"original_price":    derefF64(in.OriginalPrice, 0),
		"points_required":   *in.PointsRequired,
		"stock":             *in.Stock,
		"sales":             derefInt(in.Sales, 0),
		"sort":              derefInt(in.Sort, 0),
		"sale_status":       derefInt(in.SaleStatus, 1),
		"source_product_id": derefU64(in.SourceProductID, 0),
		"is_del":            0,
		"version":           1,
		"create_time":       now,
		"updated_at":        now,
	}
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").Create(row).Error; err != nil {
		fail(c, "积分商品添加失败")
		return
	}
	out, _ := h.loadProduct(c, id)
	response.OK(c, out)
}

func (h *Handler) QuickAddProduct(c *gin.Context) {
	var in quickAddInput
	if c.ShouldBindJSON(&in) != nil || in.SourceProductID == 0 || in.PointsRequired == nil || *in.PointsRequired < 1 || in.Stock == nil || *in.Stock < 0 {
		response.Fail(c, http.StatusBadRequest, "快速添加参数错误")
		return
	}
	var src struct {
		ProductID    uint64  `gorm:"column:product_id"`
		MerchantID   uint64  `gorm:"column:merchant_id"`
		StoreID      uint64  `gorm:"column:store_id"`
		MerchantName string  `gorm:"column:merchant_name"`
		StoreName    string  `gorm:"column:store_name"`
		Title        string  `gorm:"column:title"`
		CoverURL     string  `gorm:"column:cover_url"`
		Price        float64 `gorm:"column:price"`
		Stock        int     `gorm:"column:stock"`
	}
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_view").
		Where("product_id = ?", in.SourceProductID).
		Select("product_id,merchant_id,store_id,merchant_name,store_name,title,cover_url,price,stock").
		Take(&src).Error; err != nil || src.ProductID == 0 {
		response.Fail(c, http.StatusNotFound, "来源商品不存在")
		return
	}
	title := strings.TrimSpace(in.Title)
	if title == "" {
		title = strings.TrimSpace(src.Title)
	}
	cover := strings.TrimSpace(in.CoverURL)
	if cover == "" {
		cover = strings.TrimSpace(src.CoverURL)
	}
	price := src.Price
	if in.OriginalPrice != nil {
		price = *in.OriginalPrice
	}
	stock := *in.Stock
	id, err := h.nextProductID(c)
	if err != nil {
		fail(c, "快速添加失败")
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	row := map[string]any{
		"product_id":        id,
		"merchant_id":       src.MerchantID,
		"store_id":          src.StoreID,
		"merchant_name":     src.MerchantName,
		"store_name":        src.StoreName,
		"title":             title,
		"cate_id":           in.CateID,
		"cover_url":         cover,
		"original_price":    price,
		"points_required":   *in.PointsRequired,
		"stock":             stock,
		"sales":             0,
		"sort":              derefInt(in.Sort, 0),
		"sale_status":       derefInt(in.SaleStatus, 1),
		"source_product_id": in.SourceProductID,
		"is_del":            0,
		"version":           1,
		"create_time":       now,
		"updated_at":        now,
	}
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").Create(row).Error; err != nil {
		fail(c, "快速添加失败")
		return
	}
	out, _ := h.loadProduct(c, id)
	response.OK(c, out)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in productSaveInput
	if id == 0 || c.ShouldBindJSON(&in) != nil || !validUpdate(&in) {
		response.Fail(c, http.StatusBadRequest, "积分商品更新参数错误")
		return
	}
	changes := map[string]any{}
	fullEdit := strings.TrimSpace(in.Title) != ""
	if fullEdit {
		changes["title"] = strings.TrimSpace(in.Title)
		changes["cover_url"] = strings.TrimSpace(in.CoverURL)
	}
	if in.CateID != nil {
		changes["cate_id"] = *in.CateID
	}
	if in.OriginalPrice != nil && *in.OriginalPrice >= 0 {
		changes["original_price"] = *in.OriginalPrice
	}
	if in.PointsRequired != nil {
		changes["points_required"] = *in.PointsRequired
	}
	if in.Stock != nil {
		changes["stock"] = *in.Stock
	}
	if in.Sales != nil && *in.Sales >= 0 {
		changes["sales"] = *in.Sales
	}
	if in.Sort != nil {
		changes["sort"] = *in.Sort
	}
	if in.SaleStatus != nil {
		changes["sale_status"] = *in.SaleStatus
	}
	if in.MerchantID != nil {
		changes["merchant_id"] = *in.MerchantID
	}
	if in.StoreID != nil {
		changes["store_id"] = *in.StoreID
	}
	if n := strings.TrimSpace(in.MerchantName); n != "" {
		changes["merchant_name"] = n
	}
	if n := strings.TrimSpace(in.StoreName); n != "" {
		changes["store_name"] = n
	}
	if len(changes) == 0 {
		response.Fail(c, http.StatusBadRequest, "积分商品更新参数错误")
		return
	}
	changes["version"] = gorm.Expr("version + 1")
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").Where("product_id = ? AND is_del = 0", id)
	if in.Version > 0 {
		q = q.Where("version = ?", in.Version)
	}
	result := q.Updates(changes)
	if result.Error != nil {
		fail(c, "积分商品更新失败")
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusConflict, "积分商品已变更，请刷新后重试")
		return
	}
	row, err := h.loadProduct(c, id)
	if err != nil || row.ProductID == 0 {
		fail(c, "积分商品读取失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) SwitchStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in struct {
		Status int `json:"status"`
	}
	if id == 0 || c.ShouldBindJSON(&in) != nil || (in.Status != 0 && in.Status != 1) {
		response.Fail(c, http.StatusBadRequest, "积分商品状态参数错误")
		return
	}
	result := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").
		Where("product_id = ? AND is_del = 0", id).
		Updates(map[string]any{"sale_status": in.Status, "version": gorm.Expr("version + 1")})
	if result.Error != nil {
		fail(c, "积分商品状态修改失败")
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "积分商品不存在")
		return
	}
	row, _ := h.loadProduct(c, id)
	response.OK(c, row)
}

func (h *Handler) CopyProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "积分商品参数错误")
		return
	}
	src, err := h.loadProduct(c, id)
	if err != nil || src.ProductID == 0 {
		response.Fail(c, http.StatusNotFound, "积分商品不存在")
		return
	}
	newID, err := h.nextProductID(c)
	if err != nil {
		fail(c, "积分商品复制失败")
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	row := map[string]any{
		"product_id":        newID,
		"merchant_id":       src.MerchantID,
		"store_id":          src.StoreID,
		"merchant_name":     src.MerchantName,
		"store_name":        src.StoreName,
		"title":             src.Title + "（复制）",
		"cate_id":           src.CateID,
		"cover_url":         src.CoverURL,
		"original_price":    src.OriginalPrice,
		"points_required":   src.PointsRequired,
		"stock":             src.Stock,
		"sales":             0,
		"sort":              src.Sort,
		"sale_status":       0,
		"source_product_id": src.SourceProductID,
		"is_del":            0,
		"version":           1,
		"create_time":       now,
		"updated_at":        now,
	}
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").Create(row).Error; err != nil {
		fail(c, "积分商品复制失败")
		return
	}
	out, _ := h.loadProduct(c, newID)
	response.OK(c, out)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "积分商品参数错误")
		return
	}
	result := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").
		Where("product_id = ? AND is_del = 0", id).
		Updates(map[string]any{"is_del": 1, "sale_status": 0, "version": gorm.Expr("version + 1")})
	if result.Error != nil {
		fail(c, "积分商品删除失败")
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "积分商品不存在")
		return
	}
	response.OK(c, gin.H{"product_id": id})
}

func (h *Handler) ListExchanges(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "积分商品参数错误")
		return
	}
	page, limit := pagination(c)
	base := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order_item AS oi").
		Joins("INNER JOIN qixi_crm_b_order AS o ON o.id = oi.order_id").
		Joins("INNER JOIN qixi_crm_b_group_order AS g ON g.id = o.group_order_id").
		Where("oi.product_id = ? AND g.activity_type = ?", id, activityTypePoints)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		fail(c, "兑换记录查询失败")
		return
	}
	rows := make([]exchangeRow, 0)
	if err := base.Select(`o.id AS order_id, o.order_no, o.user_id, g.pay_status, g.points_amount, oi.quantity, oi.title_snapshot, g.created_at`).
		Order("g.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "兑换记录查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
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
	if productID, err := optionalPositive(c.Query("product_id")); err != nil {
		response.Fail(c, http.StatusBadRequest, "商品 ID 参数错误")
		return
	} else if productID != 0 {
		q = q.Where(`EXISTS (
			SELECT 1 FROM qixi_crm_b_order o
			INNER JOIN qixi_crm_b_order_item oi ON oi.order_id = o.id
			WHERE o.group_order_id = qixi_crm_b_group_order.id AND oi.product_id = ?
		)`, productID)
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

func (h *Handler) loadProduct(c *gin.Context, id uint64) (productRow, error) {
	var row productRow
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_points_product_view").
		Select(productSelect).
		Where("product_id = ? AND is_del = 0", id).
		Scan(&row).Error
	return row, err
}

func (h *Handler) nextProductID(c *gin.Context) (uint64, error) {
	var maxID uint64
	if err := h.businessDB.WithContext(c.Request.Context()).
		Raw(`SELECT GREATEST(
			COALESCE((SELECT MAX(product_id) FROM qixi_crm_b_points_product_view), 0),
			COALESCE((SELECT MAX(product_id) FROM qixi_crm_b_product_view), 0)
		)`).Scan(&maxID).Error; err != nil {
		return 0, err
	}
	if maxID < 1400 {
		maxID = 1400
	}
	return maxID + 1, nil
}

func validCreate(in *productSaveInput) bool {
	if in == nil {
		return false
	}
	title := strings.TrimSpace(in.Title)
	if title == "" || utf8.RuneCountInString(title) > 100 {
		return false
	}
	if in.PointsRequired == nil || *in.PointsRequired < 1 {
		return false
	}
	if in.Stock == nil || *in.Stock < 0 {
		return false
	}
	if in.OriginalPrice != nil && *in.OriginalPrice < 0 {
		return false
	}
	if in.SaleStatus != nil && *in.SaleStatus != 0 && *in.SaleStatus != 1 {
		return false
	}
	return true
}

func validUpdate(in *productSaveInput) bool {
	if in == nil {
		return false
	}
	if in.PointsRequired != nil && *in.PointsRequired < 1 {
		return false
	}
	if in.Stock != nil && *in.Stock < 0 {
		return false
	}
	if in.OriginalPrice != nil && *in.OriginalPrice < 0 {
		return false
	}
	if in.SaleStatus != nil && *in.SaleStatus != 0 && *in.SaleStatus != 1 {
		return false
	}
	if in.Sales != nil && *in.Sales < 0 {
		return false
	}
	return in.Version > 0 ||
		strings.TrimSpace(in.Title) != "" ||
		in.CoverURL != "" ||
		in.CateID != nil ||
		in.OriginalPrice != nil ||
		in.PointsRequired != nil ||
		in.Stock != nil ||
		in.Sort != nil ||
		in.SaleStatus != nil ||
		in.Sales != nil
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

func derefU64(p *uint64, def uint64) uint64 {
	if p == nil {
		return def
	}
	return *p
}
func derefInt(p *int, def int) int {
	if p == nil {
		return def
	}
	return *p
}
func derefF64(p *float64, def float64) float64 {
	if p == nil {
		return def
	}
	return *p
}

func fail(c *gin.Context, message string) { response.Fail(c, http.StatusInternalServerError, message) }
