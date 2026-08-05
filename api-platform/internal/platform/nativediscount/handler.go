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

type discountRules struct {
	PackagePrice float64  `json:"package_price"`
	ProductIDs   []uint64 `json:"product_ids"`
	FreeShipping bool     `json:"free_shipping"`
	Remark       string   `json:"remark"`
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
	if status := strings.TrimSpace(c.Query("status")); status != "" {
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
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	rows := make([]viewRow, 0)
	if err := q.Select("activity_id,store_id,name,rules,status,version,starts_at,ends_at").
		Order("activity_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, toItem(row))
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
	response.OK(c, toItem(row))
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

func toItem(row viewRow) gin.H {
	rules := discountRules{}
	_ = json.Unmarshal([]byte(row.Rules), &rules)
	return gin.H{
		"activity_id":   row.ActivityID,
		"store_id":      row.StoreID,
		"name":          row.Name,
		"package_price": rules.PackagePrice,
		"product_ids":   rules.ProductIDs,
		"free_shipping": rules.FreeShipping,
		"remark":        rules.Remark,
		"status":        row.Status,
		"version":       row.Version,
		"starts_at":     formatTime(row.StartsAt),
		"ends_at":       formatTime(row.EndsAt),
	}
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
