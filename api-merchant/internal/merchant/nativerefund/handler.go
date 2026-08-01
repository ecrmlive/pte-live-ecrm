// Package nativerefund owns store after-sale decisions on qixi_crm_b_*.
//
// A store decision is deliberately separated from a payment-provider refund:
// approve moves the request to merchant_handling and only a verified provider
// callback may later move it to refunded.
package nativerefund

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	db       *gorm.DB
	identity *identity.Service
}

func NewHandler(db *gorm.DB, id *identity.Service) *Handler { return &Handler{db: db, identity: id} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/refunds", h.list)
	r.GET("/refunds/:id", h.get)
	r.POST("/refunds/:id/approve", middleware.RequireMerchantMenu(h.identity, identity.MerPermRefundApprove), h.approve)
	r.POST("/refunds/:id/reject", middleware.RequireMerchantMenu(h.identity, identity.MerPermRefundReject), h.reject)
}

type refund struct {
	ID         uint64    `gorm:"column:id"`
	OrderID    uint64    `gorm:"column:order_id"`
	RefundNo   string    `gorm:"column:refund_no"`
	Reason     string    `gorm:"column:reason"`
	Amount     float64   `gorm:"column:amount"`
	Status     string    `gorm:"column:status"`
	MerchantID uint64    `gorm:"column:merchant_id"`
	StoreID    uint64    `gorm:"column:store_id"`
	UserID     uint64    `gorm:"column:user_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := normalizePage(c)
	q := h.base(c)
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
	row, ok := h.load(c, parseID(c))
	if !ok {
		return
	}
	response.OK(c, view(row))
}

func (h *Handler) approve(c *gin.Context) {
	h.transition(c, "merchant_handling", "商户已受理退款申请", func(current string) bool { return current == "applied" })
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
	h.transition(c, "rejected", reason, func(current string) bool { return current == "applied" || current == "merchant_handling" })
}

func (h *Handler) transition(c *gin.Context, target, reason string, allowed func(string) bool) {
	id := parseID(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "售后 ID 错误")
		return
	}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var row refund
		if err := tx.Table("qixi_crm_b_refund AS r").
			Select("r.id,r.order_id,r.status,o.merchant_id,o.store_id").
			Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").
			Where("r.id = ? AND o.store_id = ?", id, middleware.StoreID(c)).
			Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&row).Error; err != nil {
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
		return tx.Table("qixi_crm_b_refund_event").Create(map[string]any{
			"refund_id": row.ID, "from_status": row.Status, "to_status": target,
			"actor_type": "merchant", "actor_id": middleware.AdminID(c), "reason": reason,
			"idempotency_key": "merchant:" + strconv.FormatUint(row.ID, 10) + ":" + target,
		}).Error
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

func (h *Handler) load(c *gin.Context, id uint64) (refund, bool) {
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "售后 ID 错误")
		return refund{}, false
	}
	var row refund
	if err := h.base(c).Where("r.id = ?", id).Scan(&row).Error; err != nil {
		fail(c)
		return refund{}, false
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "售后单不存在")
		return refund{}, false
	}
	return row, true
}

func (h *Handler) base(c *gin.Context) *gorm.DB {
	return h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").
		Select("r.id,r.order_id,r.refund_no,r.reason,r.amount,r.status,r.created_at,r.updated_at,o.merchant_id,o.store_id,o.user_id").
		Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").Where("o.store_id = ?", middleware.StoreID(c))
}

func view(row refund) gin.H {
	return gin.H{"refund_order_id": row.ID, "refund_order_sn": row.RefundNo, "order_id": row.OrderID, "mer_id": row.MerchantID, "store_id": row.StoreID, "uid": row.UserID, "refund_type": 1, "refund_message": row.Reason, "refund_price": row.Amount, "status": legacyStatus(row.Status), "status_code": row.Status, "create_time": row.CreatedAt.Format("2006-01-02 15:04:05"), "status_time": row.UpdatedAt.Format("2006-01-02 15:04:05")}
}

func normalizePage(c *gin.Context) (int, int) {
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
func parseID(c *gin.Context) uint64 { id, _ := strconv.ParseUint(c.Param("id"), 10, 64); return id }
func normalizeStatus(value string) string {
	switch value {
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
		return value
	}
}
func legacyStatus(value string) int {
	switch value {
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
