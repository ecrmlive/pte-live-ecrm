// Package integralorder exposes platform supervision for points-mall orders
// (CRMEB admin.points.Order / activity_type=20).
package integralorder

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	// RequireAdminMenu 仅认 kind=button；page 码 marketing.integral.orders 只用于导航。
	menuOrdersRead   = "marketing.integral.orders.read"
	menuOrdersManage = "marketing.integral.orders.manage"
	activityPoints   = 20
	exportLimit      = 5000
)

type Handler struct {
	businessDB *gorm.DB
	adminDB    *gorm.DB
}

func NewHandler(businessDB, adminDB *gorm.DB) *Handler {
	return &Handler{businessDB: businessDB, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	read := middleware.RequireAdminMenu(h.adminDB, menuOrdersRead)
	manage := middleware.RequireAdminMenu(h.adminDB, menuOrdersManage)
	r.GET("/integral/orders", access, read, h.List)
	r.GET("/integral/orders/:id", access, read, h.Detail)
	r.POST("/integral/orders/export", access, read, h.Export)
	r.POST("/integral/orders/:id/delivery", access, manage, h.Delivery)
	r.DELETE("/integral/orders/:id", access, manage, h.Delete)
}

type listFilter struct {
	Status     string
	DateFrom   string
	DateTo     string
	SearchType string
	Keyword    string
}

type orderRow struct {
	ID              uint64          `gorm:"column:id"`
	GroupOrderID    uint64          `gorm:"column:group_order_id"`
	OrderNo         string          `gorm:"column:order_no"`
	MerchantID      uint64          `gorm:"column:merchant_id"`
	MerchantName    string          `gorm:"column:merchant_name_snapshot"`
	StoreID         uint64          `gorm:"column:store_id"`
	StoreName       string          `gorm:"column:store_name_snapshot"`
	UserID          uint64          `gorm:"column:user_id"`
	TotalAmount     float64         `gorm:"column:total_amount"`
	PayAmount       float64         `gorm:"column:pay_amount"`
	FreightAmount   float64         `gorm:"column:freight_amount"`
	TotalQuantity   int             `gorm:"column:total_quantity"`
	PointsAmount    int64           `gorm:"column:points_amount"`
	Recipient       json.RawMessage `gorm:"column:recipient_snapshot"`
	Remark          string          `gorm:"column:remark"`
	MerchantRemark  string          `gorm:"column:merchant_remark"`
	IsSystemDel     int8            `gorm:"column:is_system_del"`
	Status          string          `gorm:"column:status"`
	PaidAt          *time.Time      `gorm:"column:paid_at"`
	CreatedAt       time.Time       `gorm:"column:created_at"`
	PayChannel      string          `gorm:"column:pay_channel"`
	UserArchivedAt  *time.Time      `gorm:"column:user_archived_at"`
	GroupRemark     string          `gorm:"column:group_remark"`
}

type itemRow struct {
	ID        uint64          `gorm:"column:id"`
	OrderID   uint64          `gorm:"column:order_id"`
	ProductID uint64          `gorm:"column:product_id"`
	Title     string          `gorm:"column:title_snapshot"`
	CoverURL  string          `gorm:"column:cover_url_snapshot"`
	Spec      json.RawMessage `gorm:"column:spec_snapshot"`
	UnitPrice float64         `gorm:"column:unit_price"`
	Quantity  int             `gorm:"column:quantity"`
}

type deliveryRow struct {
	OrderID      uint64     `gorm:"column:order_id"`
	DeliveryType string     `gorm:"column:delivery_type"`
	CarrierCode  string     `gorm:"column:carrier_code"`
	TrackingNo   string     `gorm:"column:tracking_no"`
	Status       string     `gorm:"column:status"`
	DeliveredAt  *time.Time `gorm:"column:delivered_at"`
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

type deliveryInput struct {
	DeliveryType string `json:"delivery_type"`
	DeliveryName string `json:"delivery_name"`
	DeliveryID   string `json:"delivery_id"`
	Remark       string `json:"remark"`
}

type exportInput struct {
	Status     string `json:"status"`
	DateFrom   string `json:"date_from"`
	DateTo     string `json:"date_to"`
	SearchType string `json:"search_type"`
	Keyword    string `json:"keyword"`
}

func (h *Handler) List(c *gin.Context) {
	page, limit := queryfilter.Page(c)
	f := parseListFilter(c)
	q := h.baseQuery(c, f)
	var total int64
	if err := q.Distinct("o.id").Count(&total).Error; err != nil {
		fail(c, "积分订单查询失败")
		return
	}
	rows := make([]orderRow, 0)
	if err := q.Select(orderSelect).
		Order("o.id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&rows).Error; err != nil {
		fail(c, "积分订单查询失败")
		return
	}
	list, err := h.hydrate(c, rows, false)
	if err != nil {
		fail(c, "积分订单查询失败")
		return
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "订单参数错误")
		return
	}
	q := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order AS o").
		Joins("INNER JOIN qixi_crm_b_group_order AS g ON g.id = o.group_order_id").
		Where("o.id = ? AND o.activity_type = ?", id, activityPoints)
	var row orderRow
	if err := q.Select(orderSelect).Scan(&row).Error; err != nil || row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "积分订单不存在")
		return
	}
	list, err := h.hydrate(c, []orderRow{row}, true)
	if err != nil || len(list) == 0 {
		fail(c, "积分订单详情加载失败")
		return
	}
	response.OK(c, list[0])
}

func (h *Handler) Export(c *gin.Context) {
	var in exportInput
	_ = c.ShouldBindJSON(&in)
	f := listFilter{
		Status:     strings.TrimSpace(in.Status),
		DateFrom:   strings.TrimSpace(in.DateFrom),
		DateTo:     strings.TrimSpace(in.DateTo),
		SearchType: strings.TrimSpace(in.SearchType),
		Keyword:    strings.TrimSpace(in.Keyword),
	}
	if f.Status == "" {
		f.Status = strings.TrimSpace(c.Query("status"))
	}
	if f.DateFrom == "" {
		f.DateFrom = strings.TrimSpace(c.Query("date_from"))
	}
	if f.DateTo == "" {
		f.DateTo = strings.TrimSpace(c.Query("date_to"))
	}
	if f.SearchType == "" {
		f.SearchType = strings.TrimSpace(c.Query("search_type"))
	}
	if f.Keyword == "" {
		f.Keyword = strings.TrimSpace(c.Query("keyword"))
	}
	q := h.baseQuery(c, f)
	rows := make([]orderRow, 0)
	if err := q.Select(orderSelect).Order("o.id DESC").Limit(exportLimit).Scan(&rows).Error; err != nil {
		fail(c, "积分订单导出失败")
		return
	}
	list, err := h.hydrate(c, rows, true)
	if err != nil {
		fail(c, "积分订单导出失败")
		return
	}
	content, err := ordersCSV(list)
	if err != nil {
		fail(c, "积分订单导出失败")
		return
	}
	response.OK(c, gin.H{
		"file_name":  "积分订单列表_" + time.Now().Format("20060102150405") + ".csv",
		"content":    content,
		"row_count":  len(list),
		"truncated": len(list) == exportLimit,
	})
}

func (h *Handler) Delivery(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "订单参数错误")
		return
	}
	var in deliveryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "发货参数错误")
		return
	}
	deliveryType := strings.TrimSpace(in.DeliveryType)
	if deliveryType == "" {
		deliveryType = "express"
	}
	if deliveryType != "express" && deliveryType != "city" && deliveryType != "pickup" && deliveryType != "service" {
		response.Fail(c, http.StatusBadRequest, "发货类型错误")
		return
	}
	name := strings.TrimSpace(in.DeliveryName)
	track := strings.TrimSpace(in.DeliveryID)
	if deliveryType != "pickup" && deliveryType != "service" && (name == "" || track == "") {
		response.Fail(c, http.StatusBadRequest, "请填写物流公司和运单号")
		return
	}

	ctx := c.Request.Context()
	var row orderRow
	err := h.businessDB.WithContext(ctx).
		Table("qixi_crm_b_order AS o").
		Joins("INNER JOIN qixi_crm_b_group_order AS g ON g.id = o.group_order_id").
		Where("o.id = ? AND o.activity_type = ? AND o.is_system_del = 0", id, activityPoints).
		Select(orderSelect).Scan(&row).Error
	if err != nil || row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "积分订单不存在")
		return
	}
	if row.Status != "paid" && row.Status != "fulfilling" && row.Status != "awaiting_final" {
		response.Fail(c, http.StatusBadRequest, "当前订单状态不可发货")
		return
	}

	err = h.businessDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exist int64
		if err := tx.Table("qixi_crm_b_order_delivery").Where("order_id = ?", id).Count(&exist).Error; err != nil {
			return err
		}
		payload := map[string]any{
			"order_id":      id,
			"delivery_type": deliveryType,
			"carrier_code":  name,
			"tracking_no":   track,
			"status":        "shipped",
		}
		if exist > 0 {
			if err := tx.Table("qixi_crm_b_order_delivery").Where("order_id = ?", id).Updates(payload).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Table("qixi_crm_b_order_delivery").Create(payload).Error; err != nil {
				return err
			}
		}
		updates := map[string]any{"status": "shipped"}
		if remark := strings.TrimSpace(in.Remark); remark != "" {
			updates["merchant_remark"] = remark
		}
		return tx.Table("qixi_crm_b_order").Where("id = ?", id).Updates(updates).Error
	})
	if err != nil {
		fail(c, "发货失败")
		return
	}
	response.OK(c, gin.H{"order_id": id, "status": "shipped"})
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "订单参数错误")
		return
	}
	result := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order").
		Where("id = ? AND activity_type = ? AND is_system_del = 0", id, activityPoints).
		Update("is_system_del", 1)
	if result.Error != nil {
		fail(c, "删除失败")
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "订单不存在或已删除")
		return
	}
	response.OK(c, gin.H{"order_id": id})
}

const orderSelect = `o.id,o.group_order_id,o.order_no,o.merchant_id,o.merchant_name_snapshot,o.store_id,o.store_name_snapshot,
o.user_id,o.total_amount,o.pay_amount,o.freight_amount,o.total_quantity,o.points_amount,o.recipient_snapshot,
o.remark,COALESCE(o.merchant_remark,'') AS merchant_remark,COALESCE(o.is_system_del,0) AS is_system_del,
o.status,o.paid_at,o.created_at,COALESCE(g.pay_channel,'') AS pay_channel,g.user_archived_at,COALESCE(g.remark,'') AS group_remark`

func (h *Handler) baseQuery(c *gin.Context, f listFilter) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order AS o").
		Joins("INNER JOIN qixi_crm_b_group_order AS g ON g.id = o.group_order_id").
		Where("o.activity_type = ?", activityPoints)

	status := strings.TrimSpace(f.Status)
	switch status {
	case "", "all":
		q = q.Where("o.is_system_del = 0")
	case "-10", "deleted":
		// CRMEB：status=-10 → 用户已删除
		q = q.Where("o.is_system_del = 0 AND g.user_archived_at IS NOT NULL")
	case "0", "unshipped":
		q = q.Where("o.is_system_del = 0 AND o.status IN ?", []string{"paid", "fulfilling", "awaiting_final"})
	case "1", "unreceived":
		q = q.Where("o.is_system_del = 0 AND o.status = ?", "shipped")
	case "2", "unevaluated":
		// 本库无独立评价态，先按已完成展示；后续可接评价表
		q = q.Where("o.is_system_del = 0 AND o.status = ?", "completed")
	case "3", "completed":
		q = q.Where("o.is_system_del = 0 AND o.status = ?", "completed")
	case "-1", "refunded":
		q = q.Where("o.is_system_del = 0 AND o.status IN ?", []string{"cancelled", "aftersale"})
	default:
		q = q.Where("o.is_system_del = 0")
	}

	q = applyCreateTimeRange(q, f.DateFrom, f.DateTo, "o.created_at")
	return applySearch(q, f.SearchType, f.Keyword)
}

func applySearch(q *gorm.DB, searchType, keyword string) *gorm.DB {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return q
	}
	like := "%" + keyword + "%"
	switch strings.TrimSpace(searchType) {
	case "order_sn":
		return q.Where("o.order_no LIKE ?", like)
	case "real_name":
		return q.Where("JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.recipient')) LIKE ?", like)
	case "phone":
		return q.Where("JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.mobile')) LIKE ?", like)
	case "uid":
		return q.Where("CAST(o.user_id AS CHAR) = ?", keyword)
	default: // all
		return q.Where(`o.order_no LIKE ?
			OR JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.recipient')) LIKE ?
			OR JSON_UNQUOTE(JSON_EXTRACT(o.recipient_snapshot,'$.mobile')) LIKE ?
			OR CAST(o.user_id AS CHAR) = ?
			OR EXISTS (
				SELECT 1 FROM qixi_crm_b_order_item oi
				WHERE oi.order_id = o.id AND oi.title_snapshot LIKE ?
			)`, like, like, like, keyword, like)
	}
}

func (h *Handler) hydrate(c *gin.Context, rows []orderRow, includeItems bool) ([]gin.H, error) {
	if len(rows) == 0 {
		return []gin.H{}, nil
	}
	orderIDs := make([]uint64, 0, len(rows))
	userIDs := make([]uint64, 0, len(rows))
	seenUser := map[uint64]struct{}{}
	for _, row := range rows {
		orderIDs = append(orderIDs, row.ID)
		if _, ok := seenUser[row.UserID]; !ok {
			seenUser[row.UserID] = struct{}{}
			userIDs = append(userIDs, row.UserID)
		}
	}

	users := map[uint64]userRow{}
	if len(userIDs) > 0 {
		var list []userRow
		if err := h.businessDB.WithContext(c.Request.Context()).
			Table("qixi_crm_b_user").Select("id,nickname,mobile").
			Where("id IN ?", userIDs).Find(&list).Error; err != nil {
			return nil, err
		}
		for _, u := range list {
			users[u.ID] = u
		}
	}

	var items []itemRow
	if err := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order_item").Where("order_id IN ?", orderIDs).
		Find(&items).Error; err != nil {
		return nil, err
	}
	itemsByOrder := map[uint64][]gin.H{}
	firstProduct := map[uint64]gin.H{}
	for _, it := range items {
		payload := gin.H{
			"order_product_id": it.ID,
			"product_id":       it.ProductID,
			"product_info":     it.Title,
			"product_image":    it.CoverURL,
			"product_sku":      specText(it.Spec),
			"product_price":    it.UnitPrice,
			"product_num":      it.Quantity,
			"total_price":      it.UnitPrice * float64(it.Quantity),
		}
		itemsByOrder[it.OrderID] = append(itemsByOrder[it.OrderID], payload)
		if _, ok := firstProduct[it.OrderID]; !ok {
			firstProduct[it.OrderID] = payload
		}
	}

	var deliveries []deliveryRow
	_ = h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order_delivery").Where("order_id IN ?", orderIDs).
		Find(&deliveries).Error
	deliveryByOrder := map[uint64]deliveryRow{}
	for _, d := range deliveries {
		deliveryByOrder[d.OrderID] = d
	}

	out := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		var address recipient
		_ = json.Unmarshal(row.Recipient, &address)
		u := users[row.UserID]
		d := deliveryByOrder[row.ID]
		statusLabel, statusCode := statusLabelOf(row)
		userRemark := strings.TrimSpace(row.Remark)
		if userRemark == "" {
			userRemark = strings.TrimSpace(row.GroupRemark)
		}
		resp := gin.H{
			"order_id":             row.ID,
			"order_sn":             row.OrderNo,
			"group_order_id":       row.GroupOrderID,
			"mer_id":               row.MerchantID,
			"mer_name":             row.MerchantName,
			"store_id":             row.StoreID,
			"store_name":           row.StoreName,
			"uid":                  row.UserID,
			"nickname":             emptyDash(u.Nickname),
			"user_deleted":         row.UserArchivedAt != nil,
			"status":               statusCode,
			"status_label":         statusLabel,
			"points_amount":        row.PointsAmount,
			"pay_amount":           row.PayAmount,
			"pay_price":            row.PayAmount,
			"total_price":          row.TotalAmount,
			"freight_price":        row.FreightAmount,
			"total_num":            row.TotalQuantity,
			"real_name":            emptyDash(address.Recipient),
			"user_phone":           maskPhone(address.Mobile),
			"user_address":         strings.TrimSpace(address.Province + address.City + address.District + " " + address.Detail),
			"user_remark":          emptyDash(userRemark),
			"merchant_remark":      emptyDash(row.MerchantRemark),
			"delivery_type":        d.DeliveryType,
			"delivery_type_label":  deliveryTypeLabel(d.DeliveryType),
			"delivery_name":        d.CarrierCode,
			"delivery_id":          d.TrackingNo,
			"can_deliver":          canDeliver(row),
			"can_delete":           row.IsSystemDel == 0,
			"create_time":          row.CreatedAt.Format("2006-01-02 15:04:05"),
			"product":              firstProduct[row.ID],
		}
		if includeItems {
			products := itemsByOrder[row.ID]
			if products == nil {
				products = []gin.H{}
			}
			resp["products"] = products
		}
		out = append(out, resp)
	}
	return out, nil
}

func statusLabelOf(row orderRow) (string, int) {
	if row.IsSystemDel == 1 {
		return "已删除", -10
	}
	if row.UserArchivedAt != nil {
		return "已删除", -10
	}
	switch row.Status {
	case "pending_pay":
		return "待付款", -2
	case "paid", "fulfilling", "awaiting_final":
		return "待发货", 0
	case "shipped":
		return "待收货", 1
	case "completed":
		return "已完成", 3
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

func canDeliver(row orderRow) bool {
	if row.IsSystemDel == 1 || row.UserArchivedAt != nil {
		return false
	}
	return row.Status == "paid" || row.Status == "fulfilling" || row.Status == "awaiting_final"
}

func parseListFilter(c *gin.Context) listFilter {
	return listFilter{
		Status:     strings.TrimSpace(c.Query("status")),
		DateFrom:   strings.TrimSpace(c.Query("date_from")),
		DateTo:     strings.TrimSpace(c.Query("date_to")),
		SearchType: strings.TrimSpace(c.Query("search_type")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
	}
}

func applyCreateTimeRange(q *gorm.DB, from, to, column string) *gorm.DB {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	if from != "" {
		if len(from) == 10 {
			from += " 00:00:00"
		}
		q = q.Where(column+" >= ?", from)
	}
	if to != "" {
		if len(to) == 10 {
			to += " 23:59:59"
		}
		q = q.Where(column+" <= ?", to)
	}
	return q
}

func ordersCSV(list []gin.H) (string, error) {
	buf := &bytes.Buffer{}
	// UTF-8 BOM for Excel
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(buf)
	_ = w.Write([]string{"订单编号", "订单状态", "收货人", "商品信息", "兑换积分", "兑换金额", "下单时间", "用户备注", "店铺备注"})
	for _, row := range list {
		productInfo := ""
		if p, ok := row["product"].(gin.H); ok && p != nil {
			productInfo = strings.TrimSpace(stringify(p["product_info"]) + " " + stringify(p["product_sku"]))
		}
		status := stringify(row["status_label"])
		if row["user_deleted"] == true {
			status += "(用户已删除)"
		}
		_ = w.Write([]string{
			stringify(row["order_sn"]),
			status,
			stringify(row["real_name"]),
			productInfo,
			stringify(row["points_amount"]),
			stringify(row["pay_amount"]),
			stringify(row["create_time"]),
			stringify(row["user_remark"]),
			stringify(row["merchant_remark"]),
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
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
		return strings.Join(parts, "，")
	}
	var asStr string
	if err := json.Unmarshal(raw, &asStr); err == nil {
		return asStr
	}
	return string(raw)
}

func stringify(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

func emptyDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "--"
	}
	return v
}

func maskPhone(mobile string) string {
	mobile = strings.TrimSpace(mobile)
	if len(mobile) < 7 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[len(mobile)-4:]
}

func fail(c *gin.Context, message string) {
	response.Fail(c, http.StatusInternalServerError, message)
}
