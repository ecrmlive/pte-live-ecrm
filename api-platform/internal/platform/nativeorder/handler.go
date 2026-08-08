// Package nativeorder implements the platform order supervision read model.
//
// It deliberately reads the new business and merchant databases instead of the
// retired qixi_m_* projection. Fulfilment, payment and after-sale mutations
// remain owned by their respective services; the platform order page is a
// supervision page and is read-only by design.
package nativeorder

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/adminscope"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	businessDB *gorm.DB
	merchantDB *gorm.DB
	adminDB    *gorm.DB
}

func NewHandler(businessDB, merchantDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, merchantDB: merchantDB, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/orders", h.list)
	r.GET("/orders/tab-counts", h.tabCounts)
	// Static paths must register before /orders/:id.
	r.GET("/orders/verify-records", h.verifyList)
	r.GET("/orders/verify-records/summary", h.verifySummary)
	r.GET("/orders/:id", h.get)
	r.GET("/orders/:id/logs", h.logs)
}

type order struct {
	ID               uint64          `gorm:"column:id"`
	GroupID          uint64          `gorm:"column:group_order_id"`
	GroupOrderNo     string          `gorm:"column:group_order_no"`
	OrderNo          string          `gorm:"column:order_no"`
	MerchantID       uint64          `gorm:"column:merchant_id"`
	StoreID          uint64          `gorm:"column:store_id"`
	UserID           uint64          `gorm:"column:user_id"`
	PayAmount        float64         `gorm:"column:pay_amount"`
	TotalAmount      float64         `gorm:"column:total_amount"`
	DiscountAmount   float64         `gorm:"column:discount_amount"`
	FreightAmount    float64         `gorm:"column:freight_amount"`
	PointsAmount     int64           `gorm:"column:points_amount"`
	TotalQuantity    int             `gorm:"column:total_quantity"`
	ActivityType     int             `gorm:"column:activity_type"`
	Recipient        json.RawMessage `gorm:"column:recipient_snapshot"`
	Remark           string          `gorm:"column:remark"`
	Status           string          `gorm:"column:status"`
	CreatedAt        time.Time       `gorm:"column:created_at"`
	PaidAt           *time.Time      `gorm:"column:paid_at"`
	PayChannel       string          `gorm:"column:pay_channel"`
	UserArchivedAt   *time.Time      `gorm:"column:user_archived_at"`
	HasRefunded      int             `gorm:"column:has_refunded"`
	PendingComment   int             `gorm:"column:pending_comment"`
}

type orderItem struct {
	ID        uint64          `gorm:"column:id"`
	OrderID   uint64          `gorm:"column:order_id"`
	ProductID uint64          `gorm:"column:product_id"`
	Title     string          `gorm:"column:title_snapshot"`
	CoverURL  string          `gorm:"column:cover_url_snapshot"`
	Spec      json.RawMessage `gorm:"column:spec_snapshot"`
	UnitPrice float64         `gorm:"column:unit_price"`
	Quantity  int             `gorm:"column:quantity"`
}

type delivery struct {
	OrderID      uint64 `gorm:"column:order_id"`
	DeliveryType string `gorm:"column:delivery_type"`
	CarrierCode  string `gorm:"column:carrier_code"`
	TrackingNo   string `gorm:"column:tracking_no"`
}

type merchant struct {
	ID       uint64 `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	RegionID uint64 `gorm:"column:region_id"`
}

type store struct {
	ID         uint64 `gorm:"column:id"`
	MerchantID uint64 `gorm:"column:merchant_id"`
	Name       string `gorm:"column:name"`
}

type merchantView struct {
	MerchantID uint64 `gorm:"column:merchant_id"`
	TypeID     uint64 `gorm:"column:type_id"`
	CategoryID uint64 `gorm:"column:category_id"`
	IsTrader   int    `gorm:"column:is_trader"`
	TypeName   string `gorm:"column:type_name"`
	CateName   string `gorm:"column:cate_name"`
}

type userRow struct {
	ID       uint64 `gorm:"column:id"`
	Nickname string `gorm:"column:nickname"`
	Mobile   string `gorm:"column:mobile"`
}

type recipient struct {
	Recipient string `json:"recipient"`
	Mobile    string `json:"mobile"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Detail    string `json:"detail"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := normalizePage(c)
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置订单监管数据范围")
		return
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.OK(c, gin.H{"list": []gin.H{}, "total": 0, "page": page, "limit": limit})
		return
	}

	countQ := applyOrderListFilters(h.scopeQuery(c, merchantIDs), c, h.adminDB)
	var total int64
	if err := countQ.Distinct("o.id").Count(&total).Error; err != nil {
		fail(c, "查询订单失败")
		return
	}
	q := applyOrderListFilters(h.base(c, merchantIDs), c, h.adminDB)
	var rows []order
	if err := q.Order("o.created_at DESC,o.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "查询订单失败")
		return
	}
	items, err := h.responses(c, rows, true)
	if err != nil {
		fail(c, "加载订单失败")
		return
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func (h *Handler) tabCounts(c *gin.Context) {
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置订单监管数据范围")
		return
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.OK(c, emptyTabCounts())
		return
	}
	counts := emptyTabCounts()
	tabs := []string{"all", "unpaid", "unshipped", "unreceived", "unevaluated", "completed", "refunded", "deleted"}
	for _, tab := range tabs {
		var n int64
		q := applyTabStatus(applyOrderListFiltersExceptTab(h.scopeQuery(c, merchantIDs), c, h.adminDB), tab)
		if err := q.Distinct("o.id").Count(&n).Error; err != nil {
			fail(c, "统计订单失败")
			return
		}
		counts[tab] = n
	}
	response.OK(c, counts)
}

func (h *Handler) get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "订单 ID 错误")
		return
	}
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置订单监管数据范围")
		return
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.Fail(c, http.StatusNotFound, "订单不存在")
		return
	}
	var row order
	if err := h.base(c, merchantIDs).Where("o.id = ?", id).Scan(&row).Error; err != nil {
		fail(c, "查询订单失败")
		return
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "订单不存在")
		return
	}
	items, err := h.responses(c, []order{row}, true)
	if err != nil || len(items) == 0 {
		fail(c, "加载订单失败")
		return
	}
	response.OK(c, items[0])
}

func (h *Handler) logs(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "订单 ID 错误")
		return
	}
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置订单监管数据范围")
		return
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.Fail(c, http.StatusNotFound, "订单不存在")
		return
	}
	var row order
	if err := h.base(c, merchantIDs).Where("o.id = ?", id).Scan(&row).Error; err != nil || row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "订单不存在")
		return
	}

	page, limit := normalizePage(c)
	terminal := strings.TrimSpace(c.Query("terminal"))
	dateFrom := strings.TrimSpace(c.Query("date_from"))
	dateTo := strings.TrimSpace(c.Query("date_to"))

	all := h.buildSyntheticLogs(c, row)
	filtered := make([]gin.H, 0, len(all))
	for _, item := range all {
		if terminal != "" && strVal(item["terminal"]) != terminal {
			continue
		}
		at := strVal(item["operate_time"])
		if dateFrom != "" && at < dateFrom+" 00:00:00" {
			continue
		}
		if dateTo != "" && at > dateTo+" 23:59:59" {
			continue
		}
		filtered = append(filtered, item)
	}
	total := len(filtered)
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}
	response.OK(c, gin.H{"list": filtered[start:end], "total": total, "page": page, "limit": limit})
}

func (h *Handler) buildSyntheticLogs(c *gin.Context, row order) []gin.H {
	var nickname string
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user").
		Select("nickname").Where("id = ?", row.UserID).Scan(&nickname).Error
	if nickname == "" {
		nickname = "用户"
	}

	out := []gin.H{}
	out = append(out, gin.H{
		"order_sn": row.OrderNo, "content": "订单生成", "role": "用户", "operator": nickname,
		"terminal": "user", "operate_time": row.CreatedAt.Format("2006-01-02 15:04:05"),
	})
	if row.PaidAt != nil {
		out = append(out, gin.H{
			"order_sn": row.OrderNo, "content": "订单支付", "role": "用户", "operator": nickname,
			"terminal": "user", "operate_time": row.PaidAt.Format("2006-01-02 15:04:05"),
		})
	}
	if row.Status == "cancelled" {
		at := row.CreatedAt
		if row.PaidAt != nil {
			at = *row.PaidAt
		}
		out = append(out, gin.H{
			"order_sn": row.OrderNo, "content": "取消订单", "role": "系统", "operator": "系统",
			"terminal": "system", "operate_time": at.Add(time.Minute).Format("2006-01-02 15:04:05"),
		})
	}
	if row.UserArchivedAt != nil {
		out = append(out, gin.H{
			"order_sn": row.OrderNo, "content": "用户删除订单", "role": "用户", "operator": nickname,
			"terminal": "user", "operate_time": row.UserArchivedAt.Format("2006-01-02 15:04:05"),
		})
	}

	type actionRow struct {
		Action    string    `gorm:"column:action"`
		AccountID uint64    `gorm:"column:account_id"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}
	var actions []actionRow
	_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_order_action").
		Select("action,account_id,created_at").Where("order_id = ?", row.ID).
		Order("created_at DESC").Find(&actions).Error
	for _, a := range actions {
		operator := "店铺账号"
		if a.AccountID > 0 {
			operator = "店铺账号#" + strconv.FormatUint(a.AccountID, 10)
		}
		out = append(out, gin.H{
			"order_sn": row.OrderNo, "content": actionLabel(a.Action), "role": "店铺", "operator": operator,
			"terminal": "merchant", "operate_time": a.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// newest first
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

func actionLabel(action string) string {
	switch strings.TrimSpace(action) {
	case "ship", "deliver":
		return "订单发货"
	case "verify":
		return "订单核销"
	case "remark":
		return "店铺备注"
	default:
		if action == "" {
			return "订单操作"
		}
		return action
	}
}

func (h *Handler) scopeQuery(c *gin.Context, merchantIDs []uint64) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order AS o").
		Joins("JOIN qixi_crm_b_group_order AS g ON g.id = o.group_order_id")
	if merchantIDs != nil {
		q = q.Where("o.merchant_id IN ?", merchantIDs)
	}
	return q
}

func (h *Handler) base(c *gin.Context, merchantIDs []uint64) *gorm.DB {
	return h.scopeQuery(c, merchantIDs).Select(`o.id,o.group_order_id,g.order_no AS group_order_no,o.order_no,o.merchant_id,o.store_id,o.user_id,
			o.pay_amount,o.total_amount,o.discount_amount,o.freight_amount,o.points_amount,o.total_quantity,o.activity_type,
			o.recipient_snapshot,o.remark,o.status,o.created_at,o.paid_at,g.pay_channel,g.user_archived_at,
			CASE WHEN EXISTS(SELECT 1 FROM qixi_crm_b_refund r WHERE r.order_id=o.id AND r.status='refunded') THEN 1 ELSE 0 END AS has_refunded,
			CASE WHEN o.status='completed' AND EXISTS(
				SELECT 1 FROM qixi_crm_b_order_item oi
				LEFT JOIN qixi_crm_b_product_comment pc ON pc.order_item_id=oi.id AND pc.deleted_at IS NULL
				WHERE oi.order_id=o.id AND pc.id IS NULL
			) THEN 1 ELSE 0 END AS pending_comment`)
}

func applyOrderListFilters(q *gorm.DB, c *gin.Context, adminDB *gorm.DB) *gorm.DB {
	q = applyOrderListFiltersExceptTab(q, c, adminDB)
	if tab := strings.TrimSpace(c.Query("tab_status")); tab != "" {
		q = applyTabStatus(q, tab)
	} else if statusStr := strings.TrimSpace(c.Query("status")); statusStr != "" {
		q = applyLegacyStatusFilter(q, statusStr)
	}
	return q
}

func applyOrderListFiltersExceptTab(q *gorm.DB, c *gin.Context, adminDB *gorm.DB) *gorm.DB {
	if paid := c.Query("paid"); paid == "0" {
		q = q.Where("o.status = 'pending_pay'")
	} else if paid == "1" {
		q = q.Where("o.status <> 'pending_pay'")
	}
	if orderSN := strings.TrimSpace(c.Query("order_sn")); orderSN != "" {
		q = q.Where("o.order_no LIKE ?", "%"+orderSN+"%")
	}
	if merID := strings.TrimSpace(c.Query("mer_id")); merID != "" {
		if id, err := strconv.ParseUint(merID, 10, 64); err == nil && id > 0 {
			q = q.Where("o.merchant_id = ?", id)
		}
	}
	if storeID := strings.TrimSpace(c.Query("store_id")); storeID != "" {
		if id, err := strconv.ParseUint(storeID, 10, 64); err == nil && id > 0 {
			q = q.Where("o.store_id = ?", id)
		}
	}
	if cateID := strings.TrimSpace(c.Query("mer_category_id")); cateID != "" {
		if id, err := strconv.ParseUint(cateID, 10, 64); err == nil && id > 0 && adminDB != nil {
			var ids []uint64
			_ = adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
				Select("merchant_id").Where("category_id = ?", id).Pluck("merchant_id", &ids)
			if len(ids) == 0 {
				q = q.Where("1 = 0")
			} else {
				q = q.Where("o.merchant_id IN ?", ids)
			}
		}
	}
	if typeID := strings.TrimSpace(c.Query("mer_type_id")); typeID != "" {
		if id, err := strconv.ParseUint(typeID, 10, 64); err == nil && id > 0 && adminDB != nil {
			var ids []uint64
			_ = adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
				Select("merchant_id").Where("type_id = ?", id).Pluck("merchant_id", &ids)
			if len(ids) == 0 {
				q = q.Where("1 = 0")
			} else {
				q = q.Where("o.merchant_id IN ?", ids)
			}
		}
	}
	if payType := strings.TrimSpace(c.Query("pay_type")); payType != "" {
		if channel := payChannelFromType(payType); channel != "" {
			q = q.Where("g.pay_channel = ?", channel)
		}
	}
	if activity := strings.TrimSpace(c.Query("activity_type")); activity != "" {
		if n, err := strconv.Atoi(activity); err == nil {
			q = q.Where("o.activity_type = ?", n)
		}
	}
	if productName := strings.TrimSpace(c.Query("product_name")); productName != "" {
		like := "%" + productName + "%"
		q = q.Where(`EXISTS (
			SELECT 1 FROM qixi_crm_b_order_item oi
			WHERE oi.order_id = o.id AND oi.title_snapshot LIKE ?
		)`, like)
	}
	if deliveryType := strings.TrimSpace(c.Query("delivery_type")); deliveryType != "" {
		q = q.Where(`EXISTS (
			SELECT 1 FROM qixi_crm_b_order_delivery d
			WHERE d.order_id = o.id AND d.delivery_type = ?
		)`, deliveryType)
	}
	if productType := strings.TrimSpace(c.Query("product_type")); productType != "" {
		// product_type 未落在订单快照上时忽略，避免伪造过滤；保留参数兼容前端。
		_ = productType
	}

	orderType := strings.TrimSpace(c.Query("order_search_type"))
	orderKW := strings.TrimSpace(c.Query("order_search_keyword"))
	if orderKW != "" {
		like := "%" + orderKW + "%"
		switch orderType {
		case "real_name":
			q = q.Where("JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.recipient')) LIKE ?", like)
		case "phone":
			q = q.Where("JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.mobile')) LIKE ?", like)
		default:
			q = q.Where("o.order_no LIKE ?", like)
		}
	}

	userType := strings.TrimSpace(c.Query("user_search_type"))
	userKW := strings.TrimSpace(c.Query("user_search_keyword"))
	if userKW != "" {
		like := "%" + userKW + "%"
		switch userType {
		case "uid":
			if id, err := strconv.ParseUint(userKW, 10, 64); err == nil {
				q = q.Where("o.user_id = ?", id)
			} else {
				q = q.Where("1 = 0")
			}
		case "phone":
			q = q.Where(`EXISTS (
				SELECT 1 FROM qixi_crm_b_user u
				WHERE u.id = o.user_id AND u.mobile LIKE ?
			)`, like)
		default:
			q = q.Where(`EXISTS (
				SELECT 1 FROM qixi_crm_b_user u
				WHERE u.id = o.user_id AND u.nickname LIKE ?
			)`, like)
		}
	}

	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where(
			"o.order_no LIKE ? OR o.recipient_snapshot LIKE ? OR CAST(o.id AS CHAR) = ?",
			like, like, keyword,
		)
	}
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where("o.created_at >= ?", t)
		}
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where("o.created_at < ?", t.AddDate(0, 0, 1))
		}
	}
	return q
}

// applyTabStatus maps CRMEB-style platform order tabs onto qixi_crm_b_* facts.
// Deleted = group_order.user_archived_at (user soft-hide); 全部 excludes deleted.
func applyTabStatus(q *gorm.DB, tab string) *gorm.DB {
	switch strings.TrimSpace(tab) {
	case "unpaid":
		return q.Where("g.user_archived_at IS NULL AND o.status = ?", "pending_pay")
	case "unshipped":
		return q.Where("g.user_archived_at IS NULL AND o.status IN ?", []string{"paid", "fulfilling", "awaiting_final"})
	case "unreceived":
		return q.Where("g.user_archived_at IS NULL AND o.status = ?", "shipped")
	case "unevaluated":
		return q.Where(`g.user_archived_at IS NULL AND o.status = 'completed' AND EXISTS(
			SELECT 1 FROM qixi_crm_b_order_item oi
			LEFT JOIN qixi_crm_b_product_comment pc ON pc.order_item_id=oi.id AND pc.deleted_at IS NULL
			WHERE oi.order_id=o.id AND pc.id IS NULL
		)`)
	case "completed":
		return q.Where(`g.user_archived_at IS NULL AND o.status = 'completed' AND NOT EXISTS(
			SELECT 1 FROM qixi_crm_b_order_item oi
			LEFT JOIN qixi_crm_b_product_comment pc ON pc.order_item_id=oi.id AND pc.deleted_at IS NULL
			WHERE oi.order_id=o.id AND pc.id IS NULL
		)`)
	case "refunded":
		return q.Where(`g.user_archived_at IS NULL AND EXISTS(
			SELECT 1 FROM qixi_crm_b_refund r WHERE r.order_id=o.id AND r.status='refunded'
		)`)
	case "deleted":
		return q.Where("g.user_archived_at IS NOT NULL")
	case "all", "":
		return q.Where("g.user_archived_at IS NULL")
	default:
		return q
	}
}

func applyLegacyStatusFilter(q *gorm.DB, statusStr string) *gorm.DB {
	status, err := strconv.Atoi(strings.TrimSpace(statusStr))
	if err != nil {
		return q
	}
	switch status {
	case -1:
		return q.Where("o.status = ?", "cancelled")
	case 0:
		return q.Where("o.status IN ?", []string{"paid", "fulfilling"})
	case 1:
		return q.Where("o.status = ?", "shipped")
	case 2:
		return applyTabStatus(q, "unevaluated")
	case 3:
		return q.Where("o.status = ?", "completed")
	default:
		return q
	}
}

func payChannelFromType(payType string) string {
	switch payType {
	case "0", "balance":
		return "balance"
	case "1", "wechat":
		return "wechat"
	case "2", "alipay":
		return "alipay"
	case "7", "mock":
		return "mock"
	default:
		return ""
	}
}

func emptyTabCounts() gin.H {
	return gin.H{
		"all": 0, "unpaid": 0, "unshipped": 0, "unreceived": 0,
		"unevaluated": 0, "completed": 0, "refunded": 0, "deleted": 0,
	}
}

// merchantScope maps unified-admin merchant and region assignments to the
// current merchant projection. nil means platform-wide supervision.
func (h *Handler) merchantScope(c *gin.Context) ([]uint64, error) {
	scope, err := adminscope.ResolveMerchantScope(c.Request.Context(), h.adminDB, middleware.ClaimsFrom(c))
	if err != nil {
		return nil, err
	}
	if scope.Full {
		return nil, nil
	}
	ids := append([]uint64{}, scope.MerchantIDs...)
	if len(scope.RegionIDs) == 0 {
		return ids, nil
	}
	var rows []merchant
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_merchant").
		Select("id,name,region_id").Where("region_id IN ?", scope.RegionIDs).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		if !containsID(ids, row.ID) {
			ids = append(ids, row.ID)
		}
	}
	return ids, nil
}

func containsID(values []uint64, expected uint64) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func (h *Handler) responses(c *gin.Context, rows []order, includeItems bool) ([]gin.H, error) {
	if len(rows) == 0 {
		return []gin.H{}, nil
	}
	orderIDs := make([]uint64, 0, len(rows))
	merchantIDs := make([]uint64, 0, len(rows))
	storeIDs := make([]uint64, 0, len(rows))
	userIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		orderIDs = append(orderIDs, row.ID)
		merchantIDs = append(merchantIDs, row.MerchantID)
		storeIDs = append(storeIDs, row.StoreID)
		userIDs = append(userIDs, row.UserID)
	}
	var merchants []merchant
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_merchant").Where("id IN ?", merchantIDs).Find(&merchants).Error; err != nil {
		return nil, err
	}
	merchantNames := make(map[uint64]string, len(merchants))
	for _, row := range merchants {
		merchantNames[row.ID] = row.Name
	}
	var stores []store
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_store").Where("id IN ?", storeIDs).Find(&stores).Error; err != nil {
		return nil, err
	}
	storeNames := make(map[uint64]string, len(stores))
	for _, row := range stores {
		storeNames[row.ID] = row.Name
	}

	views := map[uint64]merchantView{}
	var viewRows []merchantView
	_ = h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_merchant_view AS v").
		Select("v.merchant_id,v.type_id,v.category_id,v.is_trader,COALESCE(t.name,'') AS type_name,COALESCE(c.name,'') AS cate_name").
		Joins("LEFT JOIN qixi_crm_a_merchant_type t ON t.id = v.type_id").
		Joins("LEFT JOIN qixi_crm_a_merchant_category c ON c.id = v.category_id").
		Where("v.merchant_id IN ?", merchantIDs).
		Scan(&viewRows)
	for _, row := range viewRows {
		views[row.MerchantID] = row
	}

	users := map[uint64]userRow{}
	var userRows []userRow
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user").
		Select("id,nickname,mobile").Where("id IN ?", userIDs).Find(&userRows)
	for _, row := range userRows {
		users[row.ID] = row
	}

	var deliveries []delivery
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order_delivery").Where("order_id IN ?", orderIDs).Find(&deliveries).Error; err != nil {
		return nil, err
	}
	deliveriesByOrder := map[uint64]delivery{}
	for _, row := range deliveries {
		deliveriesByOrder[row.OrderID] = row
	}

	type commissionAgg struct {
		OrderID uint64  `gorm:"column:order_id"`
		Amount  float64 `gorm:"column:amount"`
	}
	commissions := map[uint64]float64{}
	var commissionRows []commissionAgg
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_commission_ledger").
		Select("order_id, SUM(amount) AS amount").Where("order_id IN ?", orderIDs).
		Group("order_id").Scan(&commissionRows)
	for _, row := range commissionRows {
		commissions[row.OrderID] = row.Amount
	}

	spreadNames := map[uint64]string{}
	parentSpreadNames := map[uint64]string{}
	type relRow struct {
		UserID       uint64 `gorm:"column:user_id"`
		ParentUserID uint64 `gorm:"column:parent_user_id"`
	}
	var rels []relRow
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_relation").
		Select("user_id, COALESCE(parent_user_id,0) AS parent_user_id").Where("user_id IN ?", userIDs).Find(&rels)
	parentIDs := make([]uint64, 0)
	parentOf := map[uint64]uint64{}
	for _, rel := range rels {
		if rel.ParentUserID > 0 {
			parentOf[rel.UserID] = rel.ParentUserID
			parentIDs = append(parentIDs, rel.ParentUserID)
		}
	}
	grandOf := map[uint64]uint64{}
	if len(parentIDs) > 0 {
		var parentRels []relRow
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_distribution_relation").
			Select("user_id, COALESCE(parent_user_id,0) AS parent_user_id").Where("user_id IN ?", parentIDs).Find(&parentRels)
		for _, rel := range parentRels {
			if rel.ParentUserID > 0 {
				grandOf[rel.UserID] = rel.ParentUserID
				parentIDs = append(parentIDs, rel.ParentUserID)
			}
		}
		var nickRows []userRow
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user").
			Select("id,nickname").Where("id IN ?", parentIDs).Find(&nickRows)
		nickMap := map[uint64]string{}
		for _, n := range nickRows {
			nickMap[n.ID] = n.Nickname
		}
		for uid, pid := range parentOf {
			spreadNames[uid] = nickMap[pid]
			if gid := grandOf[pid]; gid > 0 {
				parentSpreadNames[uid] = nickMap[gid]
			}
		}
	}

	itemsByOrder := map[uint64][]gin.H{}
	firstProduct := map[uint64]gin.H{}
	var items []orderItem
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order_item").Where("order_id IN ?", orderIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	for _, row := range items {
		payLine := row.UnitPrice * float64(row.Quantity)
		item := gin.H{
			"order_product_id": row.ID, "product_id": row.ProductID, "product_info": row.Title,
			"product_image": row.CoverURL, "product_sku": specText(row.Spec),
			"product_price": row.UnitPrice, "cost_price": 0, "pay_price": payLine,
			"product_num": row.Quantity, "total_price": payLine,
		}
		itemsByOrder[row.OrderID] = append(itemsByOrder[row.OrderID], item)
		if _, ok := firstProduct[row.OrderID]; !ok {
			firstProduct[row.OrderID] = item
		}
	}

	out := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		var address recipient
		_ = json.Unmarshal(row.Recipient, &address)
		d := deliveriesByOrder[row.ID]
		payTime := ""
		if row.PaidAt != nil {
			payTime = row.PaidAt.Format("2006-01-02 15:04:05")
		}
		u := users[row.UserID]
		v := views[row.MerchantID]
		statusLabel, statusCode := orderStatusLabel(row)
		paidFlag := paid(row.Status)
		payTypeCode := payType(row.PayChannel)
		payTypeLabel := payTypeText(row.PayChannel, paidFlag)
		pointsYuan := float64(row.PointsAmount) / 100.0
		if row.PointsAmount > 0 && row.PointsAmount < 100 {
			pointsYuan = 0
		}
		storeCategory := "非自营"
		if v.IsTrader == 1 {
			storeCategory = "自营"
		}
		resp := gin.H{
			"order_id": row.ID, "order_sn": row.OrderNo, "group_order_id": row.GroupID,
			"group_order_sn": row.GroupOrderNo, "mer_id": row.MerchantID,
			"mer_name": merchantNames[row.MerchantID], "store_id": row.StoreID,
			"store_name": storeNames[row.StoreID], "store_type_name": emptyDash(v.TypeName),
			"store_category_name": storeCategory, "mer_category_name": emptyDash(v.CateName),
			"uid": row.UserID, "nickname": emptyDash(u.Nickname), "user_phone_mask": maskPhone(u.Mobile),
			"user_deleted": row.UserArchivedAt != nil, "paid": paidFlag, "status": statusCode,
			"status_label": statusLabel, "order_type_label": "普通订单",
			"activity_type": row.ActivityType, "activity_type_label": activityTypeLabel(row.ActivityType),
			"product_type_label": "普通商品",
			"pay_price": row.PayAmount, "total_price": row.TotalAmount,
			"discount_amount": row.DiscountAmount, "platform_coupon": row.DiscountAmount,
			"merchant_coupon": 0, "points_deduction": pointsYuan, "member_discount": 0,
			"freight_price": row.FreightAmount, "total_num": row.TotalQuantity,
			"pay_type": payTypeCode, "pay_type_label": payTypeLabel, "pay_status_label": paidText(paidFlag),
			"pay_time": payTime, "delivery_type": d.DeliveryType,
			"delivery_type_label": deliveryTypeLabel(d.DeliveryType),
			"delivery_name": d.CarrierCode, "delivery_id": d.TrackingNo,
			"real_name": emptyDash(address.Recipient), "user_phone": maskPhone(address.Mobile),
			"user_phone_raw": address.Mobile,
			"user_address": strings.TrimSpace(address.Province + address.City + address.District + " " + address.Detail),
			"user_remark": emptyDash(row.Remark), "merchant_remark": "--",
			"first_brokerage": 0, "second_brokerage": 0, "commission_total": commissions[row.ID],
			"spread_name": emptyDash(spreadNames[row.UserID]),
			"top_spread_name": emptyDash(parentSpreadNames[row.UserID]),
			"create_time": row.CreatedAt.Format("2006-01-02 15:04:05"),
			"product": firstProduct[row.ID],
		}
		if includeItems {
			resp["products"] = itemsByOrder[row.ID]
			if resp["products"] == nil {
				resp["products"] = []gin.H{}
			}
		}
		out = append(out, resp)
	}
	return out, nil
}

func orderStatusLabel(row order) (string, int) {
	if row.HasRefunded == 1 {
		return "已退款", -1
	}
	switch row.Status {
	case "pending_pay":
		return "待付款", -2
	case "paid", "fulfilling", "awaiting_final":
		return "待发货", 0
	case "shipped":
		return "待收货", 1
	case "completed":
		if row.PendingComment == 1 {
			return "待评价", 2
		}
		return "交易完成", 3
	case "cancelled":
		return "已取消", -1
	case "aftersale":
		return "售后中", 4
	case "final_timeout":
		return "已关闭", -1
	default:
		return row.Status, 0
	}
}

func activityTypeLabel(v int) string {
	switch v {
	case 1:
		return "秒杀"
	case 2:
		return "预售"
	case 3:
		return "拼团"
	case 4:
		return "助力"
	default:
		return "普通"
	}
}

func deliveryTypeLabel(v string) string {
	switch v {
	case "express":
		return "快递发货"
	case "pickup":
		return "到店自提"
	case "city":
		return "同城配送"
	case "service":
		return "虚拟发货"
	default:
		return "--"
	}
}

func specText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var asMap map[string]any
	if err := json.Unmarshal(raw, &asMap); err == nil && len(asMap) > 0 {
		parts := make([]string, 0, len(asMap))
		for k, v := range asMap {
			parts = append(parts, k+": "+stringify(v))
		}
		return strings.Join(parts, ", ")
	}
	var asArr []any
	if err := json.Unmarshal(raw, &asArr); err == nil {
		parts := make([]string, 0, len(asArr))
		for _, item := range asArr {
			parts = append(parts, stringify(item))
		}
		return strings.Join(parts, ", ")
	}
	s := strings.Trim(string(raw), `"`)
	if s == "null" || s == "{}" || s == "[]" {
		return ""
	}
	return s
}

func stringify(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case json.Number:
		return t.String()
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

func maskPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return "--"
	}
	runes := []rune(phone)
	if len(runes) < 7 {
		return phone
	}
	return string(runes[:3]) + "****" + string(runes[len(runes)-4:])
}

func emptyDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "--"
	}
	return v
}

func strVal(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
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

func paid(status string) int {
	if status == "pending_pay" {
		return 0
	}
	return 1
}

func paidText(paidFlag int) string {
	if paidFlag == 1 {
		return "已支付"
	}
	return "未支付"
}

func payType(channel string) int {
	switch channel {
	case "wechat":
		return 1
	case "alipay":
		return 2
	case "balance":
		return 0
	case "mock":
		return 7
	default:
		if channel == "" {
			return -1
		}
		return 7
	}
}

func payTypeText(channel string, paidFlag int) string {
	if paidFlag == 0 || channel == "" {
		return "未支付"
	}
	switch channel {
	case "wechat":
		return "微信"
	case "alipay":
		return "支付宝"
	case "balance":
		return "余额"
	case "mock":
		return "模拟支付"
	default:
		return channel
	}
}

func fail(c *gin.Context, message string) {
	response.Fail(c, http.StatusInternalServerError, message)
}
