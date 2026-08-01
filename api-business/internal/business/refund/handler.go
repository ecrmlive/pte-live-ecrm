// Package refund owns the C-end after-sale read and request workflow on
// qixi_crm_b_*. It does not mark a refund as succeeded: that state is reserved
// for a verified payment-channel refund callback.
package refund

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }
func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/refund/apply", h.apply)
	r.GET("/refunds", h.list)
	r.GET("/refunds/:id", h.get)
	r.POST("/refunds/:id/cancel", h.cancel)
	r.POST("/refunds/:id/platform", h.requestPlatform)
}

type applyRequest struct {
	OrderID        uint64 `json:"order_id"`
	RefundType     int    `json:"refund_type"`
	RefundMessage  string `json:"refund_message"`
	IdempotencyKey string `json:"idempotency_key"`
}

type order struct {
	ID        uint64  `gorm:"column:id"`
	UserID    uint64  `gorm:"column:user_id"`
	Status    string  `gorm:"column:status"`
	PayAmount float64 `gorm:"column:pay_amount"`
}
type orderItem struct {
	ID       uint64  `gorm:"column:id"`
	Quantity int     `gorm:"column:quantity"`
	Price    float64 `gorm:"column:unit_price"`
}
type refund struct {
	ID         uint64    `gorm:"column:id"`
	OrderID    uint64    `gorm:"column:order_id"`
	RefundNo   string    `gorm:"column:refund_no"`
	Reason     string    `gorm:"column:reason"`
	Amount     float64   `gorm:"column:amount"`
	Status     string    `gorm:"column:status"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	MerchantID uint64    `gorm:"column:merchant_id"`
	UserID     uint64    `gorm:"column:user_id"`
}
type refundItem struct {
	ID          uint64  `gorm:"column:id"`
	RefundID    uint64  `gorm:"column:refund_id"`
	OrderItemID uint64  `gorm:"column:order_item_id"`
	Quantity    int     `gorm:"column:quantity"`
	Amount      float64 `gorm:"column:amount"`
}

func (h *Handler) apply(c *gin.Context) {
	var req applyRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.OrderID == 0 || req.RefundType != 1 || strings.TrimSpace(req.RefundMessage) == "" {
		bad(c, "售后申请参数不合法")
		return
	}
	uid := uint64(middleware.UID(c))
	req.RefundMessage = strings.TrimSpace(req.RefundMessage)
	if len(req.RefundMessage) > 500 {
		bad(c, "退款原因不能超过 500 字")
		return
	}
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		req.IdempotencyKey = key("apply", uid, req.OrderID, req.RefundMessage)
	}
	var created refund
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var o order
		if err := tx.Table("qixi_crm_b_order").Where("id = ? AND user_id = ?", req.OrderID, uid).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&o).Error; err != nil {
			return err
		}
		if o.ID == 0 {
			return errNotFound
		}
		if !canApply(o.Status) {
			return errBadStatus
		}
		var previous refund
		if err := tx.Table("qixi_crm_b_refund").Where("order_id = ? AND idempotency_key = ?", req.OrderID, req.IdempotencyKey).Scan(&previous).Error; err != nil {
			return err
		}
		if previous.ID != 0 {
			created = previous
			return nil
		}
		var active int64
		if err := tx.Table("qixi_crm_b_refund").Where("order_id = ? AND status IN ?", req.OrderID, []string{"applied", "merchant_handling", "platform_intervene", "refunding"}).Count(&active).Error; err != nil {
			return err
		}
		if active > 0 {
			return errActive
		}
		var items []orderItem
		if err := tx.Table("qixi_crm_b_order_item").Where("order_id = ?", req.OrderID).Find(&items).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return errBadStatus
		}
		amount := 0.0
		for _, item := range items {
			amount += item.Price * float64(item.Quantity)
		}
		if amount <= 0 || amount > o.PayAmount {
			return errBadStatus
		}
		created = refund{OrderID: req.OrderID, RefundNo: refundNo(req.OrderID), Reason: req.RefundMessage, Amount: amount, Status: "applied"}
		if err := tx.Table("qixi_crm_b_refund").Create(map[string]any{
			"order_id": created.OrderID, "refund_no": created.RefundNo, "reason": created.Reason, "amount": created.Amount,
			"order_status_before": o.Status, "status": created.Status, "idempotency_key": req.IdempotencyKey,
		}).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_refund").Where("order_id = ? AND idempotency_key = ?", req.OrderID, req.IdempotencyKey).Scan(&created).Error; err != nil {
			return err
		}
		for _, item := range items {
			if err := tx.Table("qixi_crm_b_aftersale_item").Create(map[string]any{"refund_id": created.ID, "order_item_id": item.ID, "quantity": item.Quantity, "amount": item.Price * float64(item.Quantity)}).Error; err != nil {
				return err
			}
		}
		if err := tx.Table("qixi_crm_b_order").Where("id = ?", req.OrderID).Update("status", "aftersale").Error; err != nil {
			return err
		}
		return event(tx, created.ID, "", "applied", "user", uid, req.RefundMessage, req.IdempotencyKey)
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	if err := h.base(c).Where("r.id = ?", created.ID).Scan(&created).Error; err != nil || created.ID == 0 {
		fail(c)
		return
	}
	response.OK(c, h.view(c, created))
}

func (h *Handler) list(c *gin.Context) {
	page, limit := page(c)
	uid := uint64(middleware.UID(c))
	q := h.base(c).Where("o.user_id = ?", uid)
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
		list = append(list, h.view(c, row))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) get(c *gin.Context) {
	id := id(c)
	if id == 0 {
		bad(c, "售后 ID 错误")
		return
	}
	var row refund
	if err := h.base(c).Where("r.id = ? AND o.user_id = ?", id, middleware.UID(c)).Scan(&row).Error; err != nil {
		fail(c)
		return
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "售后单不存在")
		return
	}
	response.OK(c, h.view(c, row))
}

func (h *Handler) cancel(c *gin.Context) {
	h.transition(c, "cancelled", "取消售后申请", func(current string) bool { return current == "applied" || current == "merchant_handling" })
}
func (h *Handler) requestPlatform(c *gin.Context) {
	h.transition(c, "platform_intervene", "用户申请平台介入", func(current string) bool { return current == "applied" || current == "merchant_handling" })
}

func (h *Handler) transition(c *gin.Context, target, reason string, allowed func(string) bool) {
	refundID := id(c)
	uid := uint64(middleware.UID(c))
	if refundID == 0 {
		bad(c, "售后 ID 错误")
		return
	}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var row refund
		if err := tx.Table("qixi_crm_b_refund AS r").Select("r.id,r.order_id,r.status,o.user_id").Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").Where("r.id = ? AND o.user_id = ?", refundID, uid).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID == 0 {
			return errNotFound
		}
		if row.Status == target {
			return nil // the client retried an already completed state transition
		}
		if !allowed(row.Status) {
			return errBadStatus
		}
		if err := tx.Table("qixi_crm_b_refund").Where("id = ? AND status = ?", row.ID, row.Status).Update("status", target).Error; err != nil {
			return err
		}
		if target == "cancelled" {
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
		return event(tx, row.ID, row.Status, target, "user", uid, reason, key("transition", uid, row.ID, target))
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) base(c *gin.Context) *gorm.DB {
	return h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").
		Select("r.id,r.order_id,r.refund_no,r.reason,r.amount,r.status,r.created_at,r.updated_at,o.merchant_id,o.user_id").
		Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id")
}
func (h *Handler) view(c *gin.Context, row refund) gin.H {
	items := []gin.H{}
	refundNum := 0
	var lines []refundItem
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_aftersale_item").Where("refund_id = ?", row.ID).Find(&lines).Error; err == nil {
		for _, line := range lines {
			refundNum += line.Quantity
			items = append(items, gin.H{"refund_product_id": line.ID, "order_product_id": line.OrderItemID, "refund_price": line.Amount, "refund_num": line.Quantity})
		}
	}
	return gin.H{"refund_order_id": row.ID, "refund_order_sn": row.RefundNo, "order_id": row.OrderID, "mer_id": row.MerchantID, "uid": row.UserID, "refund_type": 1, "refund_message": row.Reason, "refund_price": row.Amount, "refund_num": refundNum, "status": legacyStatus(row.Status), "status_code": row.Status, "create_time": row.CreatedAt.Format("2006-01-02 15:04:05"), "status_time": row.UpdatedAt.Format("2006-01-02 15:04:05"), "products": items}
}

func event(tx *gorm.DB, refundID uint64, from, to, actor string, actorID uint64, reason, idem string) error {
	return tx.Table("qixi_crm_b_refund_event").Create(map[string]any{"refund_id": refundID, "from_status": from, "to_status": to, "actor_type": actor, "actor_id": actorID, "reason": reason, "idempotency_key": idem}).Error
}
func canApply(status string) bool {
	return status == "paid" || status == "fulfilling" || status == "shipped"
}
func legacyStatus(status string) int {
	switch status {
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
func id(c *gin.Context) uint64 { v, _ := strconv.ParseUint(c.Param("id"), 10, 64); return v }
func key(parts ...any) string {
	h := sha256.New()
	for _, p := range parts {
		_, _ = fmt.Fprint(h, p, "|")
	}
	return hex.EncodeToString(h.Sum(nil))
}
func refundNo(orderID uint64) string { return fmt.Sprintf("R%d%d", time.Now().UnixNano(), orderID) }
func bad(c *gin.Context, msg string) { response.Fail(c, http.StatusBadRequest, msg) }
func fail(c *gin.Context)            { response.Fail(c, http.StatusInternalServerError, "售后服务异常") }

var errNotFound = fmt.Errorf("refund not found")
var errBadStatus = fmt.Errorf("refund status invalid")
var errActive = fmt.Errorf("active refund exists")

func writeErr(c *gin.Context, err error) {
	switch err {
	case errNotFound:
		response.Fail(c, http.StatusNotFound, "售后单不存在")
	case errBadStatus:
		bad(c, "当前状态不可执行该操作")
	case errActive:
		bad(c, "订单已有处理中售后申请")
	default:
		fail(c)
	}
}
