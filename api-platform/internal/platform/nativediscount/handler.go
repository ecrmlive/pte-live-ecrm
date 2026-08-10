// Package nativediscount exposes platform-supervised discount packages from
// qixi_crm_b_marketing_activity_view. Merchants own writes; platform may only
// toggle the published projection status for supervision.
package nativediscount

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	businessDB *gorm.DB
	adminDB    *gorm.DB
}

func NewHandler(businessDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	platformOrOps := middleware.RequireAdminRoles("platform", "operations")
	read := middleware.RequireAdminMenu(h.adminDB, "marketing.discounts.read")
	manage := middleware.RequireAdminMenu(h.adminDB, "marketing.discounts.manage")
	r.GET("/marketing/discounts", platformOrOps, read, h.list)
	r.GET("/marketing/discounts/:id", platformOrOps, read, h.detail)
	r.PUT("/marketing/discounts/:id/status", platformOrOps, manage, h.setStatus)
}

type viewRow struct {
	ActivityID uint64     `gorm:"column:activity_id"`
	StoreID    uint64     `gorm:"column:store_id"`
	Name       string     `gorm:"column:name"`
	Rules      string     `gorm:"column:rules"`
	Status     int        `gorm:"column:status"`
	Version    uint64     `gorm:"column:version"`
	StartsAt   *time.Time `gorm:"column:starts_at"`
	EndsAt     *time.Time `gorm:"column:ends_at"`
}

// CRMEB: type 0=固定套餐 1=搭配套餐；product.type 0=主商品 1=搭配商品
type discountProduct struct {
	ProductID uint64 `json:"product_id"`
	StoreName string `json:"store_name"`
	Image     string `json:"image"`
	Type      int    `json:"type"`
	Spec      string `json:"spec"`
}

type discountRules struct {
	PackagePrice float64           `json:"package_price"`
	PackageType  int               `json:"package_type"`
	Type         int               `json:"type"`
	IsLimit      int               `json:"is_limit"`
	LimitNum     int               `json:"limit_num"`
	IsTime       int               `json:"is_time"`
	ProductIDs   []uint64          `json:"product_ids"`
	Products     []discountProduct `json:"products"`
	FreeShipping bool              `json:"free_shipping"`
	Remark       string            `json:"remark"`
	CreateTime   string            `json:"create_time"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pageLimit(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_marketing_activity_view").
		Where("activity_type = ?", "discount")
	if raw := strings.TrimSpace(c.Query("store_id")); raw != "" {
		storeID, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || storeID == 0 {
			response.Fail(c, http.StatusBadRequest, "店铺 ID 错误")
			return
		}
		q = q.Where("store_id = ?", storeID)
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" && status != "all" {
		switch status {
		case "1", "active", "enabled":
			q = q.Where("status = 1")
		case "0", "inactive", "disabled":
			q = q.Where("status = 0")
		default:
			response.Fail(c, http.StatusBadRequest, "活动状态错误")
			return
		}
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	packageTypeFilter := strings.TrimSpace(c.Query("package_type"))
	if packageTypeFilter == "" {
		packageTypeFilter = strings.TrimSpace(c.Query("type"))
	}
	wantType := -1
	if packageTypeFilter != "" && packageTypeFilter != "all" {
		n, err := strconv.Atoi(packageTypeFilter)
		if err != nil || (n != 0 && n != 1) {
			response.Fail(c, http.StatusBadRequest, "套餐类型错误")
			return
		}
		wantType = n
	}

	// package_type 存于 rules JSON：先取匹配行再内存过滤分页（监管列表体量可控）
	rows := make([]viewRow, 0)
	if err := q.Select("activity_id,store_id,name,rules,status,version,starts_at,ends_at").
		Order("activity_id DESC").Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	storeNames := h.lookupStoreNames(c, rows)
	filtered := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		item := toItem(row, storeNames[row.StoreID])
		if wantType >= 0 {
			pt, _ := item["package_type"].(int)
			if pt != wantType {
				continue
			}
		}
		filtered = append(filtered, item)
	}
	total := int64(len(filtered))
	start := (page - 1) * limit
	list := []gin.H{}
	if start < len(filtered) {
		end := start + limit
		if end > len(filtered) {
			end = len(filtered)
		}
		list = filtered[start:end]
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "活动 ID 错误")
		return
	}
	var row viewRow
	err = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_marketing_activity_view").
		Select("activity_id,store_id,name,rules,status,version,starts_at,ends_at").
		Where("activity_id = ? AND activity_type = ?", id, "discount").Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, http.StatusNotFound, "优惠套餐不存在")
		return
	}
	if err != nil {
		fail(c)
		return
	}
	storeNames := h.lookupStoreNames(c, []viewRow{row})
	response.OK(c, toItem(row, storeNames[row.StoreID]))
}

func (h *Handler) setStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "活动 ID 错误")
		return
	}
	var in struct {
		Status *int `json:"status"`
	}
	if c.ShouldBindJSON(&in) != nil || in.Status == nil || (*in.Status != 0 && *in.Status != 1) {
		response.Fail(c, http.StatusBadRequest, "活动状态错误")
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Exec(
		`UPDATE qixi_crm_b_marketing_activity_view SET status=?, version=version+1
		 WHERE activity_id=? AND activity_type='discount'`,
		*in.Status, id,
	)
	if res.Error != nil {
		fail(c)
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "优惠套餐不存在")
		return
	}
	response.OK(c, gin.H{"activity_id": id, "status": *in.Status})
}

func (h *Handler) lookupStoreNames(c *gin.Context, rows []viewRow) map[uint64]string {
	out := map[uint64]string{}
	if h.adminDB == nil || len(rows) == 0 {
		return out
	}
	ids := make([]uint64, 0, len(rows))
	seen := map[uint64]struct{}{}
	for _, row := range rows {
		if row.StoreID == 0 {
			continue
		}
		if _, ok := seen[row.StoreID]; ok {
			continue
		}
		seen[row.StoreID] = struct{}{}
		ids = append(ids, row.StoreID)
	}
	if len(ids) == 0 {
		return out
	}
	type merRow struct {
		MerchantID   uint64 `gorm:"column:merchant_id"`
		MerchantName string `gorm:"column:merchant_name"`
	}
	list := make([]merRow, 0, len(ids))
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
		Select("merchant_id, merchant_name").
		Where("merchant_id IN ?", ids).Scan(&list)
	for _, m := range list {
		out[m.MerchantID] = m.MerchantName
	}
	return out
}

func toItem(row viewRow, storeName string) gin.H {
	rules := parseRules(row.Rules)
	packageType := rules.PackageType
	if packageType != 0 && packageType != 1 {
		packageType = rules.Type
	}
	if packageType != 0 && packageType != 1 {
		packageType = 0
	}
	isTime := rules.IsTime
	if hasJSONKey(row.Rules, "is_time") {
		isTime = rules.IsTime
	} else if row.StartsAt != nil || row.EndsAt != nil {
		// 旧夹具未写 is_time 但有起止时间：视为限时
		isTime = 1
	}

	products := rules.Products
	if len(products) == 0 && len(rules.ProductIDs) > 0 {
		products = make([]discountProduct, 0, len(rules.ProductIDs))
		for _, pid := range rules.ProductIDs {
			products = append(products, discountProduct{
				ProductID: pid,
				StoreName: "演示商品#" + strconv.FormatUint(pid, 10),
				Image:     "",
				Type:      0,
				Spec:      "| 0.00",
			})
		}
	}
	mainProducts := make([]gin.H, 0)
	comboProducts := make([]gin.H, 0)
	for _, p := range products {
		item := gin.H{
			"product_id": p.ProductID,
			"store_name": p.StoreName,
			"image":      p.Image,
			"type":       p.Type,
			"spec":       strings.TrimSpace(p.Spec),
		}
		if p.Type == 1 {
			comboProducts = append(comboProducts, item)
		} else {
			mainProducts = append(mainProducts, item)
		}
	}

	createTime := strings.TrimSpace(rules.CreateTime)
	if createTime == "" && row.StartsAt != nil {
		createTime = formatTime(row.StartsAt)
	}

	remainLabel := "不限量"
	remainNum := 0
	if rules.IsLimit == 1 {
		remainNum = rules.LimitNum
		remainLabel = strconv.Itoa(rules.LimitNum)
	}

	timeLabel := "不限时"
	startsAt := ""
	endsAt := ""
	if isTime == 1 {
		startsAt = formatTime(row.StartsAt)
		endsAt = formatTime(row.EndsAt)
		if startsAt != "" || endsAt != "" {
			timeLabel = strings.TrimSpace(startsAt + " ~ " + endsAt)
		}
	}

	if storeName == "" {
		storeName = "店铺#" + strconv.FormatUint(row.StoreID, 10)
	}

	packageTypeLabel := "固定套餐"
	if packageType == 1 {
		packageTypeLabel = "搭配套餐"
	}

	qtyLabel := "不限量"
	if rules.IsLimit == 1 {
		qtyLabel = strconv.Itoa(rules.LimitNum)
	}

	return gin.H{
		"activity_id":         row.ActivityID,
		"id":                  row.ActivityID,
		"store_id":            row.StoreID,
		"store_name":          storeName,
		"name":                row.Name,
		"package_price":       rules.PackagePrice,
		"package_type":        packageType,
		"package_type_label":  packageTypeLabel,
		"is_limit":            rules.IsLimit,
		"limit_num":           rules.LimitNum,
		"remain_num":          remainNum,
		"remain_label":        remainLabel,
		"qty_label":           qtyLabel,
		"is_time":             isTime,
		"time_label":          timeLabel,
		"product_ids":         rules.ProductIDs,
		"products":            productsToH(products),
		"main_products":       mainProducts,
		"combo_products":      comboProducts,
		"free_shipping":       rules.FreeShipping,
		"remark":              rules.Remark,
		"status":              row.Status,
		"status_label":        statusLabel(row.Status),
		"version":             row.Version,
		"starts_at":           startsAt,
		"ends_at":             endsAt,
		"create_time":         createTime,
		"created_at":          createTime,
	}
}

func productsToH(products []discountProduct) []gin.H {
	out := make([]gin.H, 0, len(products))
	for _, p := range products {
		out = append(out, gin.H{
			"product_id": p.ProductID,
			"store_name": p.StoreName,
			"image":      p.Image,
			"type":       p.Type,
			"spec":       strings.TrimSpace(p.Spec),
		})
	}
	return out
}

func parseRules(raw string) discountRules {
	rules := discountRules{}
	_ = json.Unmarshal([]byte(raw), &rules)
	return rules
}

func hasJSONKey(raw, key string) bool {
	var m map[string]json.RawMessage
	if json.Unmarshal([]byte(raw), &m) != nil {
		return false
	}
	_, ok := m[key]
	return ok
}

func statusLabel(status int) string {
	if status == 1 {
		return "上架"
	}
	return "下架"
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func pageLimit(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}

func fail(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "优惠套餐监管查询失败")
}
