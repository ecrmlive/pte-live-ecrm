// Package nativerefund serves platform and region after-sale supervision from
// qixi_crm_b_* and maps region assignments to current qixi_crm_m_ merchants.
package nativerefund

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/adminscope"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	businessDB, merchantDB *gorm.DB
	adminDB                *gorm.DB
}

func NewHandler(businessDB, merchantDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, merchantDB: merchantDB, adminDB: adminDB}
}
func (h *Handler) Register(r gin.IRoutes) {
	platformOnly := middleware.RequireAdminRoles("platform")
	approve := middleware.RequireAdminMenu(h.adminDB, "order.refund.approve")
	reject := middleware.RequireAdminMenu(h.adminDB, "order.refund.reject")
	export := middleware.RequireAdminMenu(h.adminDB, "order.refund.export")
	r.GET("/refunds", h.list)
	r.GET("/refunds/tab-counts", h.tabCounts)
	r.POST("/refunds/export", platformOnly, export, h.export)
	// Events are platform-only read audit (detail tab「订单记录」); no money mutation.
	r.GET("/refunds/:id/events", platformOnly, h.events)
	r.GET("/refunds/:id", h.get)
	r.POST("/refunds/:id/approve", platformOnly, approve, h.approve)
	r.POST("/refunds/:id/reject", platformOnly, reject, h.reject)
}

type refund struct {
	ID, OrderID, MerchantID, StoreID, UserID uint64
	RefundNo, Reason, Status                 string
	RefundType                               string
	Amount                                   float64
	CreatedAt, UpdatedAt                     time.Time
}
type merchant struct {
	ID       uint64 `gorm:"column:id"`
	RegionID uint64 `gorm:"column:region_id"`
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

type refundExportInput struct {
	Reason string           `json:"reason"`
	Status *json.RawMessage `json:"status"`
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
	q := h.applyListFilters(c, h.base(c, ids), ids)
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
	response.OK(c, gin.H{"list": h.enrichRows(c, rows, false), "total": total, "page": page, "limit": limit})
}

func (h *Handler) tabCounts(c *gin.Context) {
	ids, ok := h.scope(c)
	if !ok {
		return
	}
	empty := gin.H{
		"all": 0, "applied": 0, "rejected": 0, "approved": 0,
		"awaiting_receipt": 0, "dispute": 0, "completed": 0,
	}
	if ids != nil && len(ids) == 0 {
		response.OK(c, empty)
		return
	}
	// Rebuild query without status/tab so each bucket shares the same non-status filters.
	values := c.Request.URL.Query()
	values.Del("tab_status")
	values.Del("status")
	origRaw := c.Request.URL.RawQuery
	c.Request.URL.RawQuery = values.Encode()
	q := h.applyListFilters(c, h.base(c, ids), ids)
	c.Request.URL.RawQuery = origRaw

	type statusCount struct {
		Status string `gorm:"column:status"`
		Cnt    int64  `gorm:"column:cnt"`
	}
	var rows []statusCount
	if err := q.Select("r.status AS status, COUNT(1) AS cnt").Group("r.status").Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	counts := map[string]int64{}
	var all int64
	for _, row := range rows {
		counts[row.Status] = row.Cnt
		all += row.Cnt
	}
	sum := func(keys ...string) int64 {
		var n int64
		for _, key := range keys {
			n += counts[key]
		}
		return n
	}
	response.OK(c, gin.H{
		"all":              all,
		"applied":          sum("applied", "merchant_handling"),
		"rejected":         sum("rejected"),
		"approved":         sum("awaiting_return", "refunding"),
		"awaiting_receipt": sum("awaiting_receipt"),
		"dispute":          sum("platform_intervene"),
		"completed":        sum("refunded"),
	})
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
	items := h.enrichRows(c, []refund{row}, true)
	if len(items) == 0 {
		response.Fail(c, http.StatusNotFound, "售后单不存在")
		return
	}
	response.OK(c, items[0])
}

// events returns the immutable state-transition trail after first checking the
// refund itself through the caller's data scope. This prevents an ID guessed
// from another merchant from disclosing its operating history.
func (h *Handler) events(c *gin.Context) {
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
	page, limit := page(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund_event").Where("refund_id = ?", id)
	if terminal := strings.TrimSpace(c.Query("terminal")); terminal != "" {
		q = q.Where("actor_type = ?", terminal)
	}
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		q = q.Where("created_at >= ?", from+" 00:00:00")
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		q = q.Where("created_at <= ?", to+" 23:59:59")
	}
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
	userIDs := make([]uint64, 0)
	for _, item := range rows {
		if item.ActorType == "user" || item.ActorType == "merchant" {
			userIDs = append(userIDs, item.ActorID)
		}
	}
	nicknames := map[uint64]string{}
	if len(userIDs) > 0 {
		var users []userRow
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user").
			Select("id,nickname").Where("id IN ?", userIDs).Find(&users)
		for _, u := range users {
			nicknames[u.ID] = u.Nickname
		}
	}
	list := make([]gin.H, 0, len(rows))
	for _, item := range rows {
		list = append(list, gin.H{
			"id": item.ID, "from_status": item.FromStatus, "to_status": item.ToStatus,
			"actor_type": item.ActorType, "actor_id": item.ActorID, "reason": item.Reason,
			"created_at": item.CreatedAt.Format("2006-01-02 15:04:05"),
			"order_sn":   row.RefundNo,
			"content":    eventContent(item.FromStatus, item.ToStatus, item.Reason),
			"role":       actorRole(item.ActorType),
			"operator":   operatorLabel(item.ActorType, item.ActorID, nicknames),
			"terminal":   item.ActorType,
			"operate_time": item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

// export creates a bounded audit CSV for a platform operator. It intentionally
// excludes user identity, refund reason and return-shipment details; the
// export cannot become an alternate personal-data extraction route.
func (h *Handler) export(c *gin.Context) {
	var in refundExportInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "退款导出参数错误")
		return
	}
	in.Reason = strings.TrimSpace(in.Reason)
	if len([]rune(in.Reason)) < 2 || len([]rune(in.Reason)) > 500 {
		response.Fail(c, http.StatusBadRequest, "导出原因需为 2 至 500 个字符")
		return
	}
	status, valid := parseExportStatus(in.Status)
	if !valid {
		response.Fail(c, http.StatusBadRequest, "退款状态筛选错误")
		return
	}
	ids, ok := h.scope(c)
	if !ok {
		return
	}
	q := h.base(c, ids)
	if status != "" {
		q = q.Where("r.status = ?", status)
	}
	rows := make([]refund, 0)
	if err := q.Order("r.id DESC").Limit(5000).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "退款导出查询失败")
		return
	}
	content, err := refundCSV(rows)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "退款导出生成失败")
		return
	}
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund_export_audit").Create(map[string]any{
		"query_fingerprint": refundExportFingerprint(status, ids), "row_count": len(rows), "reason": in.Reason, "operator_admin_id": middleware.AdminID(c),
	}).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "退款导出审计写入失败")
		return
	}
	response.OK(c, gin.H{"file_name": "退款监管导出_" + time.Now().Format("20060102150405") + ".csv", "content": content, "row_count": len(rows), "truncated": len(rows) == 5000})
}

// approve records platform acceptance. It never claims the money has returned;
// only a verified WeChat/Alipay refund callback may set refunded.
func (h *Handler) approve(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "售后 ID 错误")
		return
	}
	ids, ok := h.scope(c)
	if !ok {
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
	if row.RefundType == "return_and_refund" {
		h.transition(c, "awaiting_return", "平台已同意退货退款，等待用户寄回商品", func(s string) bool { return s == "applied" || s == "merchant_handling" || s == "platform_intervene" })
		return
	}
	h.transition(c, "refunding", "平台已同意仅退款，等待支付渠道退款回调", func(s string) bool { return s == "applied" || s == "merchant_handling" || s == "platform_intervene" })
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
		q := tx.Table("qixi_crm_b_refund AS r").Select("r.id,r.order_id,r.refund_no,r.amount,r.status,o.merchant_id").Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id").Where("r.id = ?", id)
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

// ensureRefundTransaction establishes the provider-refund idempotency record
// in the same business transaction as the platform approval. No browser value
// is used for the channel or money amount.
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

func (h *Handler) base(c *gin.Context, merchantIDs []uint64) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").Select("r.id,r.order_id,r.refund_no,r.reason,r.amount,r.refund_type,r.status,r.created_at,r.updated_at,o.merchant_id,o.store_id,o.user_id").Joins("JOIN qixi_crm_b_order AS o ON o.id = r.order_id")
	if merchantIDs != nil {
		q = q.Where("o.merchant_id IN ?", merchantIDs)
	}
	return q
}

func (h *Handler) scope(c *gin.Context) ([]uint64, bool) {
	scope, err := adminscope.ResolveMerchantScope(c.Request.Context(), h.adminDB, middleware.ClaimsFrom(c))
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置售后监管数据范围")
		return nil, false
	}
	if scope.Full {
		return nil, true
	}
	ids := append([]uint64{}, scope.MerchantIDs...)
	if len(scope.RegionIDs) == 0 {
		return ids, true
	}
	var rows []merchant
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_merchant").Select("id,region_id").Where("region_id IN ?", scope.RegionIDs).Find(&rows).Error; err != nil {
		fail(c)
		return nil, false
	}
	for _, row := range rows {
		if !containsID(ids, row.ID) {
			ids = append(ids, row.ID)
		}
	}
	return ids, true
}

func containsID(values []uint64, expected uint64) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func view(row refund) gin.H {
	return gin.H{"refund_order_id": row.ID, "refund_order_sn": row.RefundNo, "order_id": row.OrderID, "mer_id": row.MerchantID, "store_id": row.StoreID, "uid": row.UserID, "refund_type": legacyRefundType(row.RefundType), "refund_message": row.Reason, "refund_price": row.Amount, "status": legacyStatus(row.Status), "status_code": row.Status, "create_time": row.CreatedAt.Format("2006-01-02 15:04:05"), "status_time": row.UpdatedAt.Format("2006-01-02 15:04:05")}
}
func (h *Handler) view(c *gin.Context, row refund) gin.H {
	out := view(row)
	var shipment returnShipment
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund_return_shipment").Where("refund_id = ?", row.ID).Scan(&shipment).Error; err == nil && shipment.TrackingNo != "" {
		out["return_shipment"] = gin.H{"carrier_name": shipment.CarrierName, "tracking_no": shipment.TrackingNo, "remark": shipment.Remark, "submitted_at": shipment.SubmittedAt.Format("2006-01-02 15:04:05")}
	}
	return out
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
		return v
	}
}
func legacyRefundType(value string) int {
	if value == "return_and_refund" {
		return 2
	}
	return 1
}
func legacyStatus(v string) int {
	switch v {
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

func exportStatus(value int) (string, bool) {
	switch value {
	case -2, -1, 0, 1, 2, 3, 4:
		return normalizeStatus(strconv.Itoa(value)), true
	default:
		return "", false
	}
}

// parseExportStatus accepts legacy numeric status values and the current
// server-owned state code. Keeping both forms avoids forcing existing callers
// to mislabel `refunding` as the legacy "待审核" state.
func parseExportStatus(raw *json.RawMessage) (string, bool) {
	if raw == nil || len(bytes.TrimSpace(*raw)) == 0 || string(bytes.TrimSpace(*raw)) == "null" {
		return "", true
	}
	var text string
	if err := json.Unmarshal(*raw, &text); err == nil {
		text = normalizeStatus(strings.TrimSpace(text))
		return text, validRefundStatus(text)
	}
	var legacy int
	if err := json.Unmarshal(*raw, &legacy); err != nil {
		return "", false
	}
	return exportStatus(legacy)
}

func validRefundStatus(value string) bool {
	switch value {
	case "applied", "merchant_handling", "awaiting_return", "awaiting_receipt", "platform_intervene", "refunding", "refunded", "rejected", "cancelled":
		return true
	default:
		return false
	}
}

func refundCSV(rows []refund) (string, error) {
	var output bytes.Buffer
	output.Write([]byte{0xEF, 0xBB, 0xBF}) // Excel UTF-8 BOM; Chinese fixture data remains readable.
	w := csv.NewWriter(&output)
	if err := w.Write([]string{"退款单号", "订单ID", "商户ID", "店铺ID", "售后类型", "退款金额", "当前状态", "申请时间", "状态更新时间"}); err != nil {
		return "", err
	}
	for _, item := range rows {
		if err := w.Write([]string{
			refundCSVCell(item.RefundNo), strconv.FormatUint(item.OrderID, 10), strconv.FormatUint(item.MerchantID, 10), strconv.FormatUint(item.StoreID, 10),
			refundTypeName(item.RefundType), strconv.FormatFloat(item.Amount, 'f', 2, 64), refundStatusName(item.Status),
			item.CreatedAt.Format("2006-01-02 15:04:05"), item.UpdatedAt.Format("2006-01-02 15:04:05"),
		}); err != nil {
			return "", err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return output.String(), nil
}

func refundCSVCell(value string) string {
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

func refundExportFingerprint(status string, merchantIDs []uint64) string {
	parts := []string{"refund_export", status}
	if merchantIDs == nil {
		parts = append(parts, "full_scope")
	} else {
		ids := append([]uint64(nil), merchantIDs...)
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		for _, id := range ids {
			parts = append(parts, strconv.FormatUint(id, 10))
		}
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(sum[:])
}
func fail(c *gin.Context) { response.Fail(c, http.StatusInternalServerError, "售后服务异常") }

var errStatus = &statusError{}

type statusError struct{}

func (*statusError) Error() string { return "invalid refund status" }
