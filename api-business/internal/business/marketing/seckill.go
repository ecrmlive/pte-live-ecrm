package marketing

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

// SeckillHandler is the public, read-only seckill projection.  Promotion
// configuration is written by the management domains and copied into business
// storage; the C-end never reads the merchant promotion tables directly.
type SeckillHandler struct{ db *gorm.DB }

func NewSeckillHandler(db *gorm.DB) *SeckillHandler { return &SeckillHandler{db: db} }

func (h *SeckillHandler) Register(r gin.IRoutes) {
	r.GET("/seckill/times", h.Times)
	r.GET("/seckill/actives", h.List)
}

type seckillActivityView struct {
	ActivityID uint64          `gorm:"column:activity_id"`
	StoreID    uint64          `gorm:"column:store_id"`
	Name       string          `gorm:"column:name"`
	Rules      json.RawMessage `gorm:"column:rules"`
	StartsAt   *time.Time      `gorm:"column:starts_at"`
	EndsAt     *time.Time      `gorm:"column:ends_at"`
}

type seckillRules struct {
	ProductID    uint64   `json:"product_id"`
	SeckillPrice float64  `json:"seckill_price"`
	TimeSlots    []string `json:"time_slots"`
}

type seckillProductView struct {
	ProductID    uint64  `gorm:"column:product_id"`
	MerchantID   uint64  `gorm:"column:merchant_id"`
	MerchantName string  `gorm:"column:merchant_name"`
	StoreName    string  `gorm:"column:store_name"`
	Title        string  `gorm:"column:title"`
	CoverURL     string  `gorm:"column:cover_url"`
	Price        float64 `gorm:"column:price"`
	Sales        int     `gorm:"column:sales"`
}

func (h *SeckillHandler) Times(c *gin.Context) {
	items, err := h.load(c)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "秒杀场次查询失败")
		return
	}
	seen := make(map[string]bool)
	labels := make([]string, 0, 4)
	for _, item := range items {
		for _, label := range item.rules.TimeSlots {
			label = strings.TrimSpace(label)
			if label != "" && !seen[label] {
				seen[label] = true
				labels = append(labels, label)
			}
		}
	}
	sort.Strings(labels)
	list := make([]gin.H, 0, len(labels))
	for _, label := range labels {
		list = append(list, gin.H{"label": label, "active": slotActive(label, time.Now())})
	}
	response.OK(c, gin.H{"list": list})
}

func (h *SeckillHandler) List(c *gin.Context) {
	items, err := h.load(c)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "秒杀活动查询失败")
		return
	}
	rows := make([]gin.H, 0, len(items))
	now := time.Now()
	for _, item := range items {
		product := item.product
		rows = append(rows, gin.H{
			"seckill_active_id": item.activity.ActivityID,
			"name":              item.activity.Name,
			"product_id":        product.ProductID,
			"seckill_price":     item.rules.SeckillPrice,
			"price":             product.Price,
			"image":             product.CoverURL,
			"store_name":        product.Title,
			"shop_name":         product.StoreName,
			"mer_name":          product.MerchantName,
			"sales":             product.Sales,
			"time_slots":        item.rules.TimeSlots,
			"in_window":         activityActive(item.activity, item.rules, now),
			"start_day":         formatDay(item.activity.StartsAt),
			"end_day":           formatDay(item.activity.EndsAt),
		})
	}
	response.OK(c, gin.H{"list": rows, "total": len(rows), "page": 1, "limit": len(rows)})
}

type seckillItem struct {
	activity seckillActivityView
	rules    seckillRules
	product  seckillProductView
}

func (h *SeckillHandler) load(c *gin.Context) ([]seckillItem, error) {
	activities := make([]seckillActivityView, 0)
	query := h.db.WithContext(c.Request.Context()).
		Table("qixi_crm_b_marketing_activity_view").
		Where("activity_type = ? AND status = 1", "seckill").
		Where("(starts_at IS NULL OR starts_at <= ?) AND (ends_at IS NULL OR ends_at >= ?)", time.Now(), time.Now()).
		Order("starts_at ASC,activity_id DESC")
	if err := query.Find(&activities).Error; err != nil {
		return nil, err
	}

	rulesByActivity := make(map[uint64]seckillRules, len(activities))
	productIDs := make([]uint64, 0, len(activities))
	for _, activity := range activities {
		var rules seckillRules
		if err := json.Unmarshal(activity.Rules, &rules); err != nil || rules.ProductID == 0 || rules.SeckillPrice <= 0 {
			continue
		}
		rulesByActivity[activity.ActivityID] = rules
		productIDs = append(productIDs, rules.ProductID)
	}
	if len(productIDs) == 0 {
		return []seckillItem{}, nil
	}
	products := make([]seckillProductView, 0, len(productIDs))
	if err := h.db.WithContext(c.Request.Context()).
		Table("qixi_crm_b_product_view").
		Where("product_id IN ? AND sale_status = 1", productIDs).
		Find(&products).Error; err != nil {
		return nil, err
	}
	productsByID := make(map[uint64]seckillProductView, len(products))
	for _, product := range products {
		productsByID[product.ProductID] = product
	}
	items := make([]seckillItem, 0, len(activities))
	for _, activity := range activities {
		rules, ok := rulesByActivity[activity.ActivityID]
		if !ok {
			continue
		}
		product, ok := productsByID[rules.ProductID]
		if !ok {
			continue
		}
		items = append(items, seckillItem{activity: activity, rules: rules, product: product})
	}
	return items, nil
}

func activityActive(activity seckillActivityView, rules seckillRules, now time.Time) bool {
	if activity.StartsAt != nil && now.Before(*activity.StartsAt) || activity.EndsAt != nil && now.After(*activity.EndsAt) {
		return false
	}
	for _, slot := range rules.TimeSlots {
		if slotActive(slot, now) {
			return true
		}
	}
	return false
}

func slotActive(slot string, now time.Time) bool {
	parsed, err := time.Parse("15:04", strings.TrimSpace(slot))
	if err != nil {
		return false
	}
	return now.Hour() >= parsed.Hour() && now.Hour() < parsed.Hour()+4
}

func formatDay(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("2006-01-02")
}
