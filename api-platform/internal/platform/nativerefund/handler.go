// Package nativerefund serves platform and region after-sale supervision from
// qixi_crm_b_* and maps region assignments to current qixi_crm_m_ merchants.
package nativerefund

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	businessDB, merchantDB *gorm.DB
	identity               *identity.Service
}

func NewHandler(businessDB, merchantDB *gorm.DB, id *identity.Service) *Handler {
	return &Handler{businessDB: businessDB, merchantDB: merchantDB, identity: id}
}
func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/refunds", h.list)
	r.GET("/refunds/:id", h.get)
	r.POST("/refunds/:id/approve", middleware.RequirePlatformMenu(h.identity, identity.PlatPermRefundApprove), h.approve)
	r.POST("/refunds/:id/reject", middleware.RequirePlatformMenu(h.identity, identity.PlatPermRefundReject), h.reject)
}

type refund struct {
	ID, OrderID, MerchantID, StoreID, UserID uint64
	RefundNo, Reason, Status                 string
	Amount                                   float64
	CreatedAt, UpdatedAt                     time.Time
}
type merchant struct {
	ID       uint64 `gorm:"column:id"`
	RegionID uint64 `gorm:"column:region_id"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := page(c)
	ids, ok := h.scope(c)
	if !ok {
		return
	}
	if ids != nil && len(ids) == 0 {
		response.OK(c, gin.H{"list": []gin.H{}, "total": 0, "page": page, "limit": limit})
		return
	}
	q := h.base(c, ids)
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		q = q.Where("r.status = ?", normalizeStatus(status))
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	var rows []refund
	if err := q.Order("r.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, view(row))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) get(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "售后 ID 错误")
		return
	}
	ids, ok := h.scope(c)
	if !ok {
		return
	}
	if ids != nil && len(ids) == 0 {
		response.Fail(c, http.StatusNotFound, "售后单不存在")
		return
	}
	var row refund
	if err := h.base(c, ids).Where("r.id = ?", id).Scan(&row).Error; err != nil {
		fail(c)
		return
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "售后单不存在")
		return
	}
	response.OK(c, view(row))
}

// approve records platform acceptance. It never claims the money has returned;
// only a verified WeChat/Alipay refund callback may set refunded.
func (h *Handler) approve(c *gin.Context) {
	h.transition(c, "refunding", "平台已受理退款申请", func(s string) bool { return s == "applied" || s == "merchant_handling" || s == "platform_intervene" })
}
func (h *Handler) reject(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
		Remark string `json:"remark"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = strings.TrimSpace(req.Remark)
	}
	if reason == "" || len([]rune(reason)) > 500 {
		response.Fail(c, http.StatusBadRequest, "驳回原因不能为空且不能超过 500 字")
		return
	}
	h.transition(c, "rejected", reason, func(s string) bool { return s == "applied" || s == "merchant_handling" || s == "platform_intervene" })
}

func (h *Handler) transition(c *gin.Context, target, reason string, allowed func(string) bool) {
	id := parseID(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "售后 ID 错误")
		return
	}
	ids, ok := h.scope(c)
	if !ok {
		return
	}
	if ids != nil && len(ids) == 0 {
		response.Fail(c, http.StatusNotFound, "售后单不存在")
		return
	}
	err := h.businessDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		q := tx.Table("qixi_crm_b_refund AS r").Select("r.id,r.order_id,r.status,o.merchant_id").Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").Where("r.id = ?", id)
		if ids != nil {
			q = q.Where("o.merchant_id IN ?", ids)
		}
		var row refund
		if err := q.Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID == 0 {
			return gorm.ErrRecordNotFound
		}
		if row.Status == target {
			return nil
		}
		if !allowed(row.Status) {
			return errStatus
		}
		if err := tx.Table("qixi_crm_b_refund").Where("id = ? AND status = ?", row.ID, row.Status).Update("status", target).Error; err != nil {
			return err
		}
		if target == "rejected" {
			var prior struct {
				Status string `gorm:"column:order_status_before"`
			}
			if err := tx.Table("qixi_crm_b_refund").Select("order_status_before").Where("id = ?", row.ID).Scan(&prior).Error; err != nil {
				return err
			}
			if err := tx.Table("qixi_crm_b_order").Where("id = ? AND status = 'aftersale'", row.OrderID).Update("status", prior.Status).Error; err != nil {
				return err
			}
		}
		return tx.Table("qixi_crm_b_refund_event").Create(map[string]any{"refund_id": row.ID, "from_status": row.Status, "to_status": target, "actor_type": "platform", "actor_id": middleware.AdminID(c), "reason": reason, "idempotency_key": "platform:" + strconv.FormatUint(row.ID, 10) + ":" + target}).Error
	})
	switch err {
	case nil:
		response.OK(c, gin.H{"ok": true})
	case gorm.ErrRecordNotFound:
		response.Fail(c, http.StatusNotFound, "售后单不存在")
	case errStatus:
		response.Fail(c, http.StatusConflict, "当前售后状态不可执行该操作")
	default:
		fail(c)
	}
}

func (h *Handler) base(c *gin.Context, merchantIDs []uint64) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").Select("r.id,r.order_id,r.refund_no,r.reason,r.amount,r.status,r.created_at,r.updated_at,o.merchant_id,o.store_id,o.user_id").Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id")
	if merchantIDs != nil {
		q = q.Where("o.merchant_id IN ?", merchantIDs)
	}
	return q
}

func (h *Handler) scope(c *gin.Context) ([]uint64, bool) {
	regionIDs, err := h.identity.PlatformRegionScope(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		fail(c)
		return nil, false
	}
	if regionIDs == nil {
		return nil, true
	}
	if len(regionIDs) == 0 {
		return []uint64{}, true
	}
	var rows []merchant
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_merchant").Select("id,region_id").Where("region_id IN ?", regionIDs).Find(&rows).Error; err != nil {
		fail(c)
		return nil, false
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	return ids, true
}

func view(row refund) gin.H {
	return gin.H{"refund_order_id": row.ID, "refund_order_sn": row.RefundNo, "order_id": row.OrderID, "mer_id": row.MerchantID, "store_id": row.StoreID, "uid": row.UserID, "refund_type": 1, "refund_message": row.Reason, "refund_price": row.Amount, "status": legacyStatus(row.Status), "status_code": row.Status, "create_time": row.CreatedAt.Format("2006-01-02 15:04:05"), "status_time": row.UpdatedAt.Format("2006-01-02 15:04:05")}
}
func page(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if p < 1 {
		p = 1
	}
	if l < 1 {
		l = 20
	}
	if l > 100 {
		l = 100
	}
	return p, l
}
func parseID(c *gin.Context) uint64 { id, _ := strconv.ParseUint(c.Param("id"), 10, 64); return id }
func normalizeStatus(v string) string {
	switch v {
	case "0":
		return "applied"
	case "3":
		return "refunded"
	case "4":
		return "platform_intervene"
	case "-1":
		return "rejected"
	case "-2":
		return "cancelled"
	default:
		return v
	}
}
func legacyStatus(v string) int {
	switch v {
	case "platform_intervene":
		return 4
	case "refunded":
		return 3
	case "rejected":
		return -1
	case "cancelled":
		return -2
	default:
		return 0
	}
}
func fail(c *gin.Context) { response.Fail(c, http.StatusInternalServerError, "售后服务异常") }

var errStatus = &statusError{}

type statusError struct{}

func (*statusError) Error() string { return "invalid refund status" }
