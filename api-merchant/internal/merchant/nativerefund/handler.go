// Package nativerefund owns store after-sale decisions on qixi_crm_b_*.
//
// A store decision is deliberately separated from a payment-provider refund:
// approve moves the request to merchant_handling and only a verified provider
// callback may later move it to refunded.
package nativerefund

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	db, merchantDB *gorm.DB
}

func NewHandler(db, merchantDB *gorm.DB) *Handler { return &Handler{db: db, merchantDB: merchantDB} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/refunds", h.list)
	r.GET("/refunds/export", middleware.RequireStorePermission(h.merchantDB, "refund.export"), h.export)
	r.GET("/refunds/:id/events", middleware.RequireStorePermission(h.merchantDB, "refund.log"), h.events)
	r.GET("/refunds/:id/express", middleware.RequireStorePermission(h.merchantDB, "refund.express"), h.express)
	r.GET("/refunds/:id", h.get)
	r.POST("/refunds/:id/remark", middleware.RequireStorePermission(h.merchantDB, "refund.remark"), h.remark)
	r.DELETE("/refunds/:id", middleware.RequireStorePermission(h.merchantDB, "refund.delete"), h.remove)
	r.POST("/refunds/:id/approve", middleware.RequireStorePermission(h.merchantDB, "refund.approve"), h.approve)
	r.POST("/refunds/:id/confirm-return", middleware.RequireStorePermission(h.merchantDB, "refund.approve"), h.confirmReturn)
	r.POST("/refunds/:id/reject", middleware.RequireStorePermission(h.merchantDB, "refund.reject"), h.reject)
}

type refund struct {
	ID         uint64    `gorm:"column:id"`
	OrderID    uint64    `gorm:"column:order_id"`
	RefundNo   string    `gorm:"column:refund_no"`
	Reason     string    `gorm:"column:reason"`
	Amount     float64   `gorm:"column:amount"`
	RefundType string    `gorm:"column:refund_type"`
	Status     string    `gorm:"column:status"`
	MerchantID uint64    `gorm:"column:merchant_id"`
	StoreID    uint64    `gorm:"column:store_id"`
	UserID     uint64    `gorm:"column:user_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}
type returnShipment struct {
	CarrierName string    `gorm:"column:carrier_name"`
	TrackingNo  string    `gorm:"column:tracking_no"`
	Remark      string    `gorm:"column:remark"`
	SubmittedAt time.Time `gorm:"column:submitted_at"`
}
type refundEvent struct {
	ID                              uint64 `gorm:"column:id"`
	FromStatus, ToStatus, ActorType string
	ActorID                         uint64    `gorm:"column:actor_id"`
	Reason                          string    `gorm:"column:reason"`
	CreatedAt                       time.Time `gorm:"column:created_at"`
}
type refundRemark struct {
	ID             uint64    `gorm:"column:id"`
	RefundID       uint64    `gorm:"column:refund_id"`
	StoreID        uint64    `gorm:"column:store_id"`
	AccountID      uint64    `gorm:"column:account_id"`
	Action         string    `gorm:"column:action"`
	Note           string    `gorm:"column:note"`
	IdempotencyKey string    `gorm:"column:idempotency_key"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}
type hiddenRefund struct {
	RefundID, StoreID uint64
	Reason            string
	IdempotencyKey    string
}

func (h *Handler) list(c *gin.Context) {
	page, limit := normalizePage(c)
	q := h.base(c)
	hidden, err := h.hiddenIDs(c)
	if err != nil {
		fail(c)
		return
	}
	if len(hidden) > 0 {
		q = q.Where("r.id NOT IN ?", hidden)
	}
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
		list = append(list, h.view(c, row))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) get(c *gin.Context) {
	row, ok := h.load(c, parseID(c))
	if !ok {
		return
	}
	response.OK(c, h.view(c, row))
}

// events exposes only the current store's immutable refund transitions.
func (h *Handler) events(c *gin.Context) {
	row, ok := h.load(c, parseID(c))
	if !ok {
		return
	}
	page, limit := normalizePage(c)
	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_refund_event").Where("refund_id = ?", row.ID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	var rows []refundEvent
	if err := q.Select("id,from_status,to_status,actor_type,actor_id,reason,created_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, item := range rows {
		list = append(list, gin.H{"id": item.ID, "from_status": item.FromStatus, "to_status": item.ToStatus, "actor_type": item.ActorType, "actor_id": item.ActorID, "reason": item.Reason, "created_at": item.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

// express returns the buyer's registered return shipment after store-scoped
// refund lookup. It is a tracking snapshot, not a third-party logistics call.
func (h *Handler) express(c *gin.Context) {
	row, ok := h.load(c, parseID(c))
	if !ok {
		return
	}
	var shipment returnShipment
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_refund_return_shipment").Where("refund_id = ?", row.ID).Scan(&shipment).Error; err != nil {
		fail(c)
		return
	}
	if shipment.TrackingNo == "" {
		response.Fail(c, http.StatusNotFound, "用户尚未登记退货物流")
		return
	}
	response.OK(c, gin.H{"carrier_name": shipment.CarrierName, "tracking_no": shipment.TrackingNo, "remark": shipment.Remark, "submitted_at": shipment.SubmittedAt.Format("2006-01-02 15:04:05")})
}

// export produces a bounded CSV without the buyer's ID, reason or shipment
// information, preventing an export permission from becoming a privacy bypass.
func (h *Handler) export(c *gin.Context) {
	status := strings.TrimSpace(c.Query("status"))
	if status != "" && !validExportStatus(status) {
		response.Fail(c, http.StatusBadRequest, "退款状态筛选错误")
		return
	}
	q := h.base(c)
	hidden, err := h.hiddenIDs(c)
	if err != nil {
		fail(c)
		return
	}
	if len(hidden) > 0 {
		q = q.Where("r.id NOT IN ?", hidden)
	}
	if status != "" {
		q = q.Where("r.status = ?", normalizeStatus(status))
	}
	var rows []refund
	if err := q.Order("r.id DESC").Limit(5000).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "退款导出查询失败")
		return
	}
	content, err := refundCSV(rows)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "退款导出生成失败")
		return
	}
	response.OK(c, gin.H{"file_name": "店铺退款导出_" + time.Now().Format("20060102150405") + ".csv", "content": content, "row_count": len(rows), "truncated": len(rows) == 5000})
}

// remove hides a refund only from the current store's console. It does not
// delete or alter business facts, payment transactions, platform supervision
// or callback processing.
func (h *Handler) remove(c *gin.Context) {
	var in struct {
		Reason         string `json:"reason"`
		IdempotencyKey string `json:"idempotency_key"`
	}
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "删除参数错误")
		return
	}
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	if len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 || len([]rune(in.IdempotencyKey)) < 8 || len([]rune(in.IdempotencyKey)) > 128 {
		response.Fail(c, http.StatusBadRequest, "删除原因或幂等键错误")
		return
	}
	id := parseID(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "售后 ID 错误")
		return
	}
	var row refund
	if err := h.base(c).Where("r.id = ?", id).Scan(&row).Error; err != nil {
		fail(c)
		return
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "售后单不存在")
		return
	}
	storeID, accountID := uint64(middleware.StoreID(c)), uint64(middleware.AdminID(c))
	var replayed bool
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var prior hiddenRefund
		err := tx.Table("qixi_crm_m_refund_hidden").Where("store_id = ? AND idempotency_key = ?", storeID, in.IdempotencyKey).Take(&prior).Error
		if err == nil {
			if prior.RefundID != row.ID || prior.Reason != in.Reason {
				return errIdempotency
			}
			replayed = true
			return nil
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		return tx.Table("qixi_crm_m_refund_hidden").Create(map[string]any{"refund_id": row.ID, "store_id": storeID, "deleted_by_account_id": accountID, "reason": in.Reason, "idempotency_key": in.IdempotencyKey}).Error
	})
	if err == errIdempotency {
		response.Fail(c, http.StatusConflict, "幂等键已用于不同删除请求")
		return
	}
	if err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"ok": true, "replayed": replayed})
}

// remark writes an immutable shop-side operation note. It never mutates the
// business refund reason or status, and the request key makes retry safe.
func (h *Handler) remark(c *gin.Context) {
	var in struct {
		Note           string `json:"note"`
		IdempotencyKey string `json:"idempotency_key"`
	}
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "备注参数错误")
		return
	}
	in.Note, in.IdempotencyKey = strings.TrimSpace(in.Note), strings.TrimSpace(in.IdempotencyKey)
	if !validRemark(in.Note, in.IdempotencyKey) {
		response.Fail(c, http.StatusBadRequest, "备注需为 1 至 500 个字符，幂等键需为 8 至 128 个字符")
		return
	}
	row, ok := h.load(c, parseID(c))
	if !ok {
		return
	}
	storeID, accountID := uint64(middleware.StoreID(c)), uint64(middleware.AdminID(c))
	var replayed bool
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var prior refundRemark
		err := tx.Table("qixi_crm_m_aftersale_action").Where("refund_id = ? AND store_id = ? AND account_id = ? AND action = 'remark' AND idempotency_key = ?", row.ID, storeID, accountID, in.IdempotencyKey).Take(&prior).Error
		if err == nil {
			if prior.Note != in.Note {
				return errIdempotency
			}
			replayed = true
			return nil
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		return tx.Table("qixi_crm_m_aftersale_action").Create(map[string]any{"refund_id": row.ID, "store_id": storeID, "account_id": accountID, "action": "remark", "note": in.Note, "attachments": gorm.Expr("JSON_ARRAY()"), "idempotency_key": in.IdempotencyKey}).Error
	})
	if err == errIdempotency {
		response.Fail(c, http.StatusConflict, "幂等键已用于不同备注")
		return
	}
	if err != nil {
		// A concurrent same-key insert is safely resolved by reading the
		// immutable row back instead of writing a second remark.
		var prior refundRemark
		lookup := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_aftersale_action").Where("refund_id = ? AND store_id = ? AND account_id = ? AND action = 'remark' AND idempotency_key = ?", row.ID, storeID, accountID, in.IdempotencyKey).Take(&prior).Error
		if lookup == nil && prior.Note == in.Note {
			response.OK(c, gin.H{"ok": true, "replayed": true})
			return
		}
		fail(c)
		return
	}
	response.OK(c, gin.H{"ok": true, "replayed": replayed})
}

func (h *Handler) approve(c *gin.Context) {
	row, ok := h.load(c, parseID(c))
	if !ok {
		return
	}
	if row.RefundType == "return_and_refund" {
		h.transition(c, "awaiting_return", "商户已同意退货退款，等待用户寄回商品", func(current string) bool { return current == "applied" })
		return
	}
	h.transition(c, "refunding", "商户已同意仅退款，等待支付渠道退款回调", func(current string) bool { return current == "applied" })
}

// confirmReturn is only available after the buyer has registered return
// logistics. It intentionally moves to refunding rather than refunded: money
// state stays owned by the verified payment callback.
func (h *Handler) confirmReturn(c *gin.Context) {
	row, ok := h.load(c, parseID(c))
	if !ok {
		return
	}
	if row.RefundType != "return_and_refund" {
		response.Fail(c, http.StatusConflict, "仅退货退款可确认收货")
		return
	}
	h.transition(c, "refunding", "商户已确认收到退回商品，等待支付渠道退款回调", func(current string) bool { return current == "awaiting_receipt" })
}

func (h *Handler) reject(c *gin.Context) {
	var req struct {
		Reason      string `json:"reason"`
		Remark      string `json:"remark"`
		FailMessage string `json:"fail_message"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = strings.TrimSpace(req.Remark)
	}
	if reason == "" {
		reason = strings.TrimSpace(req.FailMessage)
	}
	if reason == "" || len([]rune(reason)) > 500 {
		response.Fail(c, http.StatusBadRequest, "驳回原因不能为空且不能超过 500 字")
		return
	}
	h.transition(c, "rejected", reason, func(current string) bool {
		return current == "applied" || current == "awaiting_return" || current == "merchant_handling"
	})
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
			Select("r.id,r.order_id,r.refund_no,r.amount,r.status,r.refund_type,o.merchant_id,o.store_id").
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
		if target == "refunding" {
			if err := ensureRefundTransaction(tx, row); err != nil {
				return err
			}
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

// ensureRefundTransaction binds an approved after-sale request to the
// server-owned successful payment record. It is called under the refund row
// lock, so duplicate operator clicks cannot create duplicate provider refunds.
func ensureRefundTransaction(tx *gorm.DB, row refund) error {
	var payment struct {
		ID      uint64  `gorm:"column:id"`
		Channel string  `gorm:"column:channel"`
		Amount  float64 `gorm:"column:amount"`
		Status  string  `gorm:"column:status"`
	}
	if err := tx.Table("qixi_crm_b_order AS o").Select("p.id,p.channel,p.amount,p.status").Joins("JOIN qixi_crm_b_payment_transaction AS p ON p.group_order_id = o.group_order_id").Where("o.id = ?", row.OrderID).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&payment).Error; err != nil {
		return err
	}
	if payment.ID == 0 || payment.Status != "succeeded" || (payment.Channel != "wechat" && payment.Channel != "alipay" && payment.Channel != "balance" && payment.Channel != "mock") || row.Amount <= 0 || row.Amount > payment.Amount {
		return errStatus
	}
	key := payment.Channel + ":" + row.RefundNo
	var count int64
	if err := tx.Table("qixi_crm_b_refund_transaction").Where("refund_id = ? AND idempotency_key = ?", row.ID, key).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return nil
	}
	return tx.Table("qixi_crm_b_refund_transaction").Create(map[string]any{"refund_id": row.ID, "channel": payment.Channel, "provider_refund_no": row.RefundNo, "amount": row.Amount, "status": "created", "idempotency_key": key}).Error
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
	hidden, err := h.isHidden(c, row.ID)
	if err != nil {
		fail(c)
		return refund{}, false
	}
	if hidden {
		response.Fail(c, http.StatusNotFound, "售后单不存在")
		return refund{}, false
	}
	return row, true
}

func (h *Handler) hiddenIDs(c *gin.Context) ([]uint64, error) {
	var rows []struct {
		RefundID uint64 `gorm:"column:refund_id"`
	}
	err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_refund_hidden").Select("refund_id").Where("store_id = ?", middleware.StoreID(c)).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]uint64, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.RefundID)
	}
	return out, nil
}
func (h *Handler) isHidden(c *gin.Context, refundID uint64) (bool, error) {
	var count int64
	err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_refund_hidden").Where("refund_id = ? AND store_id = ?", refundID, middleware.StoreID(c)).Count(&count).Error
	return count > 0, err
}

func (h *Handler) base(c *gin.Context) *gorm.DB {
	return h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").
		Select("r.id,r.order_id,r.refund_no,r.reason,r.amount,r.refund_type,r.status,r.created_at,r.updated_at,o.merchant_id,o.store_id,o.user_id").
		Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").Where("o.store_id = ?", middleware.StoreID(c))
}

func view(row refund) gin.H {
	return gin.H{"refund_order_id": row.ID, "refund_order_sn": row.RefundNo, "order_id": row.OrderID, "mer_id": row.MerchantID, "store_id": row.StoreID, "uid": row.UserID, "refund_type": legacyRefundType(row.RefundType), "refund_message": row.Reason, "refund_price": row.Amount, "status": legacyStatus(row.Status), "status_code": row.Status, "create_time": row.CreatedAt.Format("2006-01-02 15:04:05"), "status_time": row.UpdatedAt.Format("2006-01-02 15:04:05")}
}
func (h *Handler) view(c *gin.Context, row refund) gin.H {
	out := view(row)
	var shipment returnShipment
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_refund_return_shipment").Where("refund_id = ?", row.ID).Scan(&shipment).Error; err == nil && shipment.TrackingNo != "" {
		out["return_shipment"] = gin.H{"carrier_name": shipment.CarrierName, "tracking_no": shipment.TrackingNo, "remark": shipment.Remark, "submitted_at": shipment.SubmittedAt.Format("2006-01-02 15:04:05")}
	}
	return out
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
	case "1":
		return "awaiting_return"
	case "2":
		return "awaiting_receipt"
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
func legacyRefundType(value string) int {
	if value == "return_and_refund" {
		return 2
	}
	return 1
}
func legacyStatus(value string) int {
	switch value {
	case "awaiting_return":
		return 1
	case "awaiting_receipt":
		return 2
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

func validExportStatus(value string) bool {
	switch value {
	case "-2", "-1", "0", "1", "2", "3", "4":
		return true
	default:
		return false
	}
}

func validRemark(note, key string) bool {
	return len([]rune(note)) >= 1 && len([]rune(note)) <= 500 && len([]rune(key)) >= 8 && len([]rune(key)) <= 128
}

func refundCSV(rows []refund) (string, error) {
	var output bytes.Buffer
	output.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&output)
	if err := w.Write([]string{"退款单号", "订单ID", "售后类型", "退款金额", "当前状态", "申请时间", "状态更新时间"}); err != nil {
		return "", err
	}
	for _, item := range rows {
		if err := w.Write([]string{csvCell(item.RefundNo), strconv.FormatUint(item.OrderID, 10), refundTypeName(item.RefundType), strconv.FormatFloat(item.Amount, 'f', 2, 64), refundStatusName(item.Status), item.CreatedAt.Format("2006-01-02 15:04:05"), item.UpdatedAt.Format("2006-01-02 15:04:05")}); err != nil {
			return "", err
		}
	}
	w.Flush()
	return output.String(), w.Error()
}

func csvCell(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "=") || strings.HasPrefix(value, "+") || strings.HasPrefix(value, "-") || strings.HasPrefix(value, "@") {
		return "'" + value
	}
	return value
}

func refundTypeName(value string) string {
	if value == "return_and_refund" {
		return "退货退款"
	}
	return "仅退款"
}

func refundStatusName(value string) string {
	switch value {
	case "awaiting_return":
		return "待退货"
	case "awaiting_receipt":
		return "待收货"
	case "refunding":
		return "退款中"
	case "refunded":
		return "已退款"
	case "platform_intervene":
		return "平台介入"
	case "rejected":
		return "审核拒绝"
	case "cancelled":
		return "用户已取消"
	default:
		return "待审核"
	}
}
func fail(c *gin.Context) { response.Fail(c, http.StatusInternalServerError, "售后服务异常") }

var errStatus = &statusError{}
var errIdempotency = &idempotencyError{}

type statusError struct{}

func (*statusError) Error() string { return "invalid refund status" }

type idempotencyError struct{}

func (*idempotencyError) Error() string { return "idempotency conflict" }
