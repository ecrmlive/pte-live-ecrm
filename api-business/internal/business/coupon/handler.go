// Package coupon owns the C-end coupon projection. It deliberately uses the
// business read model and never reaches the retired qixi_m_* coupon tables.
package coupon

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/business/order"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

// RegisterPublic exposes the coupon catalogue without granting any stateful
// coupon action. A guest may inspect the centre; receiving and all user coupon
// state remain behind the JWT-protected routes below.
func (h *Handler) RegisterPublic(r gin.IRoutes) {
	r.GET("/coupons/center", h.center)
}

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/coupons/:id/receive", h.receive)
	r.GET("/coupons/newcomer", h.newcomer)
	r.GET("/coupons/mine", h.mine)
	r.GET("/coupons/usable", h.usable)
}

type templateRow struct {
	CouponID      uint64     `gorm:"column:coupon_id"`
	StoreID       uint64     `gorm:"column:store_id"`
	MerchantID    uint64     `gorm:"column:merchant_id"`
	Name          string     `gorm:"column:name"`
	DiscountType  string     `gorm:"column:discount_type"`
	DiscountValue float64    `gorm:"column:discount_value"`
	MinAmount     float64    `gorm:"column:min_amount"`
	StartsAt      *time.Time `gorm:"column:starts_at"`
	EndsAt        *time.Time `gorm:"column:ends_at"`
	Status        int8       `gorm:"column:status"`
	Received      int        `gorm:"column:received"`
}

type userCouponRow struct {
	ID            uint64     `gorm:"column:id"`
	CouponID      uint64     `gorm:"column:coupon_id"`
	StoreID       uint64     `gorm:"column:store_id"`
	MerchantID    uint64     `gorm:"column:merchant_id"`
	Name          string     `gorm:"column:name"`
	DiscountType  string     `gorm:"column:discount_type"`
	DiscountValue float64    `gorm:"column:discount_value"`
	MinAmount     float64    `gorm:"column:min_amount"`
	Status        string     `gorm:"column:status"`
	StartsAt      *time.Time `gorm:"column:starts_at"`
	EndsAt        *time.Time `gorm:"column:ends_at"`
}

func (h *Handler) center(c *gin.Context) {
	uid := uint64(middleware.UID(c))
	now := time.Now()
	merchantID, err := queryMerchant(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "商户 ID 错误")
		return
	}
	q := h.baseTemplates(c, uid).Where("c.status = ?", 1).Where("(c.starts_at IS NULL OR c.starts_at <= ?)", now).Where("(c.ends_at IS NULL OR c.ends_at >= ?)", now)
	if merchantID != 0 {
		q = q.Where("s.merchant_id = ?", merchantID)
	}
	var rows []templateRow
	if err := q.Order("c.store_id DESC,c.coupon_id DESC").Scan(&rows).Error; err != nil {
		internal(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, couponView(row))
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) receive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "优惠券 ID 错误")
		return
	}
	uid := uint64(middleware.UID(c))
	now := time.Now()
	err = h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Table("qixi_crm_b_coupon_template_view").Where("coupon_id = ? AND status = 1 AND (starts_at IS NULL OR starts_at <= ?) AND (ends_at IS NULL OR ends_at >= ?)", id, now, now).Count(&count).Error; err != nil {
			return err
		}
		if count != 1 {
			return errCouponUnavailable
		}
		return tx.Table("qixi_crm_b_coupon_user").Clauses(clause.OnConflict{DoNothing: true}).Create(map[string]any{"user_id": uid, "coupon_id": id, "source": "center", "status": "unused"}).Error
	})
	if err != nil {
		if errors.Is(err, errCouponUnavailable) {
			response.Fail(c, http.StatusNotFound, "优惠券不存在或已失效")
		} else {
			internal(c)
		}
		return
	}
	var row userCouponRow
	if err := h.baseUserCoupons(c, uid).Where("u.coupon_id = ?", id).Scan(&row).Error; err != nil || row.ID == 0 {
		internal(c)
		return
	}
	response.OK(c, userCouponView(row))
}

func (h *Handler) mine(c *gin.Context) {
	uid := uint64(middleware.UID(c))
	page, limit := pageParams(c)
	status := strings.TrimSpace(c.Query("status"))
	if !validMineStatus(status) {
		response.Fail(c, http.StatusBadRequest, "优惠券状态错误")
		return
	}
	q := h.baseUserCoupons(c, uid)
	switch status {
	case "", "all":
	case "0", "unused":
		q = q.Where("u.status = ? AND (c.ends_at IS NULL OR c.ends_at >= ?)", "unused", time.Now())
	case "1", "used":
		q = q.Where("u.status = ?", "used")
	case "2", "expired":
		q = q.Where("u.status = ? OR (u.status = ? AND c.ends_at IS NOT NULL AND c.ends_at < ?)", "expired", "unused", time.Now())
	case "locked":
		q = q.Where("u.status = ?", "locked")
	case "history":
		q = q.Where("u.status IN ? OR (u.status = ? AND c.ends_at IS NOT NULL AND c.ends_at < ?)", []string{"used", "locked", "expired"}, "unused", time.Now())
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		internal(c)
		return
	}
	var rows []userCouponRow
	if err := q.Order("u.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		internal(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, userCouponView(row))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) usable(c *gin.Context) {
	uid := uint64(middleware.UID(c))
	ids, err := queryIDs(c.Query("cart_ids"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "购物车参数错误")
		return
	}
	checkout, err := order.LoadCheckout(c.Request.Context(), h.db, uid, ids)
	if err != nil {
		writeCheckoutError(c, err)
		return
	}
	if len(checkout.Stores) != 1 {
		response.Fail(c, http.StatusBadRequest, "当前订单不可使用优惠券")
		return
	}
	now := time.Now()
	q := h.baseUserCoupons(c, uid).
		Where("u.status = ?", "unused").
		Where("c.status = ?", 1).
		Where("(c.starts_at IS NULL OR c.starts_at <= ?)", now).
		Where("(c.ends_at IS NULL OR c.ends_at >= ?)", now).
		Where("(c.store_id = 0 OR c.store_id = ?)", checkout.Stores[0].StoreID).
		Where("c.min_amount <= ?", float64(checkout.TotalCents)/100)
	var rows []userCouponRow
	if err := q.Order("c.store_id DESC,c.discount_value DESC,u.id DESC").Scan(&rows).Error; err != nil {
		internal(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, userCouponView(row))
	}
	response.OK(c, gin.H{"list": list, "total_price": float64(checkout.TotalCents) / 100})
}

func (h *Handler) baseTemplates(c *gin.Context, uid uint64) *gorm.DB {
	return h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_coupon_template_view AS c").
		Select("c.coupon_id,c.store_id,COALESCE(s.merchant_id,0) AS merchant_id,c.name,c.discount_type,c.discount_value,c.min_amount,c.starts_at,c.ends_at,c.status,CASE WHEN u.id IS NULL THEN 0 ELSE 1 END AS received").
		Joins("LEFT JOIN qixi_crm_b_store_view AS s ON s.store_id = c.store_id").
		Joins("LEFT JOIN qixi_crm_b_coupon_user AS u ON u.coupon_id = c.coupon_id AND u.user_id = ?", uid)
}
func (h *Handler) baseUserCoupons(c *gin.Context, uid uint64) *gorm.DB {
	return h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_coupon_user AS u").
		Select("u.id,u.coupon_id,c.store_id,COALESCE(s.merchant_id,0) AS merchant_id,c.name,c.discount_type,c.discount_value,c.min_amount,u.status,c.starts_at,c.ends_at").
		Joins("JOIN qixi_crm_b_coupon_template_view AS c ON c.coupon_id = u.coupon_id").
		Joins("LEFT JOIN qixi_crm_b_store_view AS s ON s.store_id = c.store_id").Where("u.user_id = ?", uid)
}
func couponView(row templateRow) gin.H {
	return gin.H{
		"coupon_id":      row.CouponID,
		"mer_id":         row.MerchantID,
		"title":          row.Name,
		"coupon_price":   row.DiscountValue,
		"discount_type":  row.DiscountType,
		"discount_value": row.DiscountValue,
		"use_min_price":  row.MinAmount,
		"coupon_time":    0,
		"remain_count":   999999,
		"is_limited":     0,
		"type":           map[bool]int{true: 1, false: 0}[row.StoreID != 0],
		"received":       row.Received == 1,
	}
}
func userCouponView(row userCouponRow) gin.H {
	status := row.Status
	if status == "unused" && row.EndsAt != nil && row.EndsAt.Before(time.Now()) {
		status = "expired"
	}
	return gin.H{
		"coupon_user_id": row.ID,
		"coupon_id":      row.CouponID,
		"mer_id":         row.MerchantID,
		"coupon_title":   row.Name,
		"coupon_price":   row.DiscountValue,
		"discount_type":  row.DiscountType,
		"discount_value": row.DiscountValue,
		"use_min_price":  row.MinAmount,
		"status":         map[string]int{"unused": 0, "locked": 1, "used": 2, "expired": -1}[status],
		"coupon_kind":    map[bool]int{true: 1, false: 0}[row.StoreID != 0],
		"starts_at":      row.StartsAt,
		"ends_at":        row.EndsAt,
	}
}
func validMineStatus(status string) bool {
	switch status {
	case "", "all", "0", "unused", "1", "used", "2", "expired", "locked", "history":
		return true
	default:
		return false
	}
}

func queryMerchant(c *gin.Context) (uint64, error) {
	raw := c.Query("mer_id")
	if raw == "" {
		return 0, nil
	}
	return strconv.ParseUint(raw, 10, 64)
}
func queryIDs(raw string) ([]uint64, error) {
	parts := strings.Split(raw, ",")
	ids := make([]uint64, 0, len(parts))
	seen := map[uint64]bool{}
	for _, part := range parts {
		id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil || id == 0 || seen[id] {
			return nil, errCouponUnavailable
		}
		seen[id] = true
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, errCouponUnavailable
	}
	return ids, nil
}
func pageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
func writeCheckoutError(c *gin.Context, err error) {
	response.Fail(c, http.StatusBadRequest, err.Error())
}
func internal(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "优惠券服务异常")
}

var errCouponUnavailable = errors.New("coupon unavailable")
