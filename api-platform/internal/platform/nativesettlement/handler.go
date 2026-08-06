// Package nativesettlement exposes the platform-owned merchant settlement
// projection. It is intentionally read-only: settlement approval and payout
// remain merchant-owned high-risk workflows and must arrive through events.
package nativesettlement

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/adminscope"
	merchantsettlement "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/merchantsettlement"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	adminDB  *gorm.DB
	commands *merchantsettlement.Client
}

func NewHandler(adminDB *gorm.DB, commands *merchantsettlement.Client) *Handler {
	return &Handler{adminDB: adminDB, commands: commands}
}

func (h *Handler) Register(r gin.IRoutes) {
	platformOrRegion := middleware.RequireAdminRoles("platform", "region")
	readSettlement := middleware.RequireAdminMenu(h.adminDB, "accounts.merchant_settlement.read")
	r.GET("/finance/merchant-settlements", platformOrRegion, readSettlement, h.List)
	r.GET("/finance/merchant-settlements/summary", platformOrRegion, readSettlement, h.Summary)
	// 转账记录 = 结算打款链路只读投影，取代 setting_cache stub。
	r.GET("/finance/transfer-records", platformOrRegion, readSettlement, h.ListTransferRecords)
	platformOnly := middleware.RequireAdminRoles("platform")
	reviewSettlement := middleware.RequireAdminMenu(h.adminDB, "accounts.merchant_settlement.review")
	r.POST("/finance/merchant-settlements/:id/approve", platformOnly, reviewSettlement, h.approve)
	r.POST("/finance/merchant-settlements/:id/reject", platformOnly, reviewSettlement, h.reject)
	r.POST("/finance/merchant-settlements/:id/mark-paid", platformOnly, reviewSettlement, h.markPaid)
}

type settlementRow struct {
	SourceSettlementID uint64    `gorm:"column:source_settlement_id" json:"settlement_id"`
	MerchantID         uint64    `gorm:"column:merchant_id" json:"merchant_id"`
	StoreID            uint64    `gorm:"column:store_id" json:"store_id"`
	MerchantName       string    `gorm:"column:merchant_name" json:"merchant_name"`
	PeriodStart        time.Time `gorm:"column:period_start" json:"period_start"`
	PeriodEnd          time.Time `gorm:"column:period_end" json:"period_end"`
	Amount             float64   `gorm:"column:amount" json:"amount"`
	Status             string    `gorm:"column:status" json:"status"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type settlementSummary struct {
	Status string  `gorm:"column:status" json:"status"`
	Amount float64 `gorm:"column:amount" json:"amount"`
	Count  int64   `gorm:"column:count" json:"count"`
}

type commandInput struct {
	IdempotencyKey  string `json:"idempotency_key"`
	ReviewNote      string `json:"review_note"`
	PayoutReference string `json:"payout_reference"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := paging(c)
	q, err := h.scopedQuery(c)
	if err != nil {
		writeScopeFailure(c, err)
		return
	}
	if status, ok := settlementStatus(c.Query("status")); !ok {
		response.Fail(c, http.StatusBadRequest, "结算状态错误")
		return
	} else if status != "" {
		q = q.Where("status = ?", status)
	}
	if rawMerchantID := strings.TrimSpace(c.Query("merchant_id")); rawMerchantID != "" {
		merchantID, err := strconv.ParseUint(rawMerchantID, 10, 64)
		if err != nil || merchantID == 0 {
			response.Fail(c, http.StatusBadRequest, "商户 ID 错误")
			return
		}
		q = q.Where("merchant_id = ?", merchantID)
	}
	q = queryfilter.ApplyCreatedAtRange(q, c, "updated_at")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	rows := make([]settlementRow, 0)
	if err := q.Select("source_settlement_id,merchant_id,store_id,merchant_name,period_start,period_end,amount,status,updated_at").Order("updated_at DESC, source_settlement_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Summary(c *gin.Context) {
	q, err := h.scopedQuery(c)
	if err != nil {
		writeScopeFailure(c, err)
		return
	}
	rows := make([]settlementSummary, 0)
	err = q.
		Select("status, COALESCE(SUM(amount), 0) AS amount, COUNT(*) AS count").Group("status").Order("status ASC").Scan(&rows).Error
	if err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"list": rows})
}

// ListTransferRecords returns settlement payout pipeline rows only
// (approved / paid / rejected). It never returns bank accounts or secrets.
func (h *Handler) ListTransferRecords(c *gin.Context) {
	page, limit := paging(c)
	q, err := h.scopedQuery(c)
	if err != nil {
		writeScopeFailure(c, err)
		return
	}
	if status, ok := transferStatus(c.Query("status")); !ok {
		response.Fail(c, http.StatusBadRequest, "转账状态错误")
		return
	} else if status != "" {
		q = q.Where("status = ?", status)
	} else {
		q = q.Where("status IN ?", []string{"approved", "paid", "rejected"})
	}
	if rawMerchantID := strings.TrimSpace(c.Query("merchant_id")); rawMerchantID != "" {
		merchantID, err := strconv.ParseUint(rawMerchantID, 10, 64)
		if err != nil || merchantID == 0 {
			response.Fail(c, http.StatusBadRequest, "商户 ID 错误")
			return
		}
		q = q.Where("merchant_id = ?", merchantID)
	}
	q = queryfilter.ApplyCreatedAtRange(q, c, "updated_at")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	rows := make([]settlementRow, 0)
	if err := q.Select("source_settlement_id,merchant_id,store_id,merchant_name,period_start,period_end,amount,status,updated_at").Order("updated_at DESC, source_settlement_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) approve(c *gin.Context)  { h.command(c, "approve") }
func (h *Handler) reject(c *gin.Context)   { h.command(c, "reject") }
func (h *Handler) markPaid(c *gin.Context) { h.command(c, "mark_paid") }

func (h *Handler) command(c *gin.Context, action string) {
	settlementID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	var in commandInput
	if err != nil || settlementID == 0 || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "结算命令参数错误")
		return
	}
	in.IdempotencyKey, in.ReviewNote, in.PayoutReference = strings.TrimSpace(in.IdempotencyKey), strings.TrimSpace(in.ReviewNote), strings.TrimSpace(in.PayoutReference)
	command := merchantsettlement.Command{SettlementID: settlementID, Action: action, OperatorID: uint64(middleware.AdminID(c)), IdempotencyKey: in.IdempotencyKey, ReviewNote: in.ReviewNote, PayoutReference: in.PayoutReference}
	if !validSettlementCommand(command) {
		response.Fail(c, http.StatusBadRequest, "结算命令参数错误")
		return
	}
	out, err := h.commands.Dispatch(c.Request.Context(), command)
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "店铺结算命令服务不可用")
		return
	}
	switch out.Code {
	case "":
		response.OK(c, gin.H{"settlement_id": out.SettlementID, "status": out.Status, "message": "命令已完成，监管投影将通过事件异步刷新"})
	case "not_found":
		response.Fail(c, http.StatusNotFound, "结算单不存在")
	case "conflict":
		response.Fail(c, http.StatusConflict, "结算单状态已变化或幂等键冲突")
	case "invalid":
		response.Fail(c, http.StatusBadRequest, "结算命令参数错误")
	default:
		response.Fail(c, http.StatusServiceUnavailable, "店铺结算命令处理失败")
	}
}

func (h *Handler) scopedQuery(c *gin.Context) (*gorm.DB, error) {
	scope, err := adminscope.ResolveMerchantScope(c.Request.Context(), h.adminDB, middleware.ClaimsFrom(c))
	if err != nil {
		return nil, err
	}
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_settlement_view")
	regionIDs, err := settlementRegionScope(scope)
	if err != nil {
		return nil, err
	}
	if regionIDs == nil {
		return q, nil
	}
	return q.Where("region_id IN ?", regionIDs), nil
}

func settlementRegionScope(scope adminscope.MerchantScope) ([]uint64, error) {
	if scope.Full {
		return nil, nil
	}
	// 结算监管仅允许区域范围；兼有 merchant 角色也不能用 merchant_ids
	// 放大财务可见范围，店铺财务操作仍归店铺管理系统。
	if len(scope.RegionIDs) == 0 {
		return nil, adminscope.ErrNotConfigured
	}
	return scope.RegionIDs, nil
}

func settlementStatus(raw string) (string, bool) {
	switch strings.TrimSpace(raw) {
	case "", "bill_pending", "bill_frozen", "withdraw_applied", "approved", "paid", "rejected", "cancelled":
		return strings.TrimSpace(raw), true
	default:
		return "", false
	}
}

func transferStatus(raw string) (string, bool) {
	switch strings.TrimSpace(raw) {
	case "", "approved", "paid", "rejected":
		return strings.TrimSpace(raw), true
	default:
		return "", false
	}
}

func validSettlementCommand(in merchantsettlement.Command) bool {
	if in.SettlementID == 0 || in.OperatorID == 0 || !validIdempotencyKey(in.IdempotencyKey) {
		return false
	}
	switch in.Action {
	case "approve":
		return len([]rune(in.ReviewNote)) <= 500
	case "reject":
		return in.ReviewNote != "" && len([]rune(in.ReviewNote)) <= 500
	case "mark_paid":
		return len([]rune(in.PayoutReference)) >= 3 && len([]rune(in.PayoutReference)) <= 128
	default:
		return false
	}
}

func validIdempotencyKey(value string) bool {
	length := len([]rune(strings.TrimSpace(value)))
	return length >= 8 && length <= 128
}

func paging(c *gin.Context) (int, int) {
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
	response.Fail(c, http.StatusInternalServerError, "店铺结算监管投影查询失败")
}

func writeScopeFailure(c *gin.Context, err error) {
	if errors.Is(err, adminscope.ErrNotConfigured) {
		response.Fail(c, http.StatusForbidden, "未配置可监管的结算区域")
		return
	}
	fail(c)
}
