package nativerefund

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type orderMeta struct {
	ID                   uint64 `gorm:"column:id"`
	OrderNo              string `gorm:"column:order_no"`
	MerchantID           uint64 `gorm:"column:merchant_id"`
	MerchantNameSnapshot string `gorm:"column:merchant_name_snapshot"`
	StoreID              uint64 `gorm:"column:store_id"`
	StoreNameSnapshot    string `gorm:"column:store_name_snapshot"`
	UserID               uint64 `gorm:"column:user_id"`
	Remark               string `gorm:"column:remark"`
}

type aftersaleItem struct {
	ID          uint64  `gorm:"column:id"`
	RefundID    uint64  `gorm:"column:refund_id"`
	OrderItemID uint64  `gorm:"column:order_item_id"`
	Quantity    int     `gorm:"column:quantity"`
	Amount      float64 `gorm:"column:amount"`
}

type orderItemRow struct {
	ID        uint64          `gorm:"column:id"`
	OrderID   uint64          `gorm:"column:order_id"`
	ProductID uint64          `gorm:"column:product_id"`
	Title     string          `gorm:"column:title_snapshot"`
	CoverURL  string          `gorm:"column:cover_url_snapshot"`
	Spec      json.RawMessage `gorm:"column:spec_snapshot"`
	UnitPrice float64         `gorm:"column:unit_price"`
	Quantity  int             `gorm:"column:quantity"`
}

type merchantViewRow struct {
	MerchantID uint64 `gorm:"column:merchant_id"`
	IsTrader   int    `gorm:"column:is_trader"`
}

type userRow struct {
	ID       uint64 `gorm:"column:id"`
	Nickname string `gorm:"column:nickname"`
	Mobile   string `gorm:"column:mobile"`
}

func tabStatuses(tab string) []string {
	switch strings.TrimSpace(tab) {
	case "", "all":
		return nil
	case "0", "applied", "audit":
		return []string{"applied", "merchant_handling"}
	case "-1", "rejected", "refuse":
		return []string{"rejected"}
	case "1", "approved", "agree":
		// CRMEB tab「审核通过」= legacy status 1 (待退货); include refunding for仅退款 in-flight.
		return []string{"awaiting_return", "refunding"}
	case "2", "awaiting_receipt", "backgood":
		return []string{"awaiting_receipt"}
	case "4", "dispute", "platform", "platform_intervene":
		return []string{"platform_intervene"}
	case "3", "completed", "end", "refunded":
		return []string{"refunded"}
	default:
		if validRefundStatus(tab) {
			return []string{tab}
		}
		return nil
	}
}

func statusLabel(code string) string {
	switch code {
	case "applied", "merchant_handling":
		return "待审核"
	case "awaiting_return":
		return "待退货"
	case "awaiting_receipt":
		return "待收货"
	case "refunding":
		return "退款中"
	case "refunded":
		return "已退款"
	case "platform_intervene":
		return "纠纷中"
	case "rejected":
		return "审核未通过"
	case "cancelled":
		return "已取消"
	default:
		return refundStatusName(code)
	}
}

func refundTypeLabel(value string) string {
	if value == "return_and_refund" {
		return "退货退款"
	}
	return "仅退款"
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
	v = strings.TrimSpace(v)
	if v == "" {
		return "--"
	}
	return v
}

func specText(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return "默认"
	}
	var asMap map[string]any
	if err := json.Unmarshal(raw, &asMap); err == nil && len(asMap) > 0 {
		parts := make([]string, 0, len(asMap))
		for k, v := range asMap {
			parts = append(parts, strings.TrimSpace(k)+":"+strings.TrimSpace(toString(v)))
		}
		if len(parts) > 0 {
			return strings.Join(parts, " ")
		}
	}
	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil && strings.TrimSpace(asString) != "" {
		return asString
	}
	return "默认"
}

func toString(v any) string {
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

func parseAttachments(raw json.RawMessage) []string {
	if len(raw) == 0 || string(raw) == "null" {
		return []string{}
	}
	var urls []string
	if err := json.Unmarshal(raw, &urls); err == nil {
		out := make([]string, 0, len(urls))
		for _, u := range urls {
			if s := strings.TrimSpace(u); s != "" {
				out = append(out, s)
			}
		}
		return out
	}
	var objects []map[string]any
	if err := json.Unmarshal(raw, &objects); err == nil {
		out := make([]string, 0, len(objects))
		for _, obj := range objects {
			for _, key := range []string{"url", "src", "path"} {
				if s := strings.TrimSpace(toString(obj[key])); s != "" && s != "null" {
					out = append(out, s)
					break
				}
			}
		}
		return out
	}
	return []string{}
}

func actorRole(actorType string) string {
	switch actorType {
	case "user":
		return "用户"
	case "merchant":
		return "店铺"
	case "platform":
		return "平台"
	case "system":
		return "系统"
	default:
		return actorType
	}
}

func eventContent(fromStatus, toStatus, reason string) string {
	reason = strings.TrimSpace(reason)
	if reason != "" {
		return reason
	}
	key := fromStatus + "->" + toStatus
	switch key {
	case "->applied":
		return "提交退款申请"
	case "applied->awaiting_return", "merchant_handling->awaiting_return", "platform_intervene->awaiting_return":
		return "同意退货退款"
	case "applied->refunding", "merchant_handling->refunding", "platform_intervene->refunding", "awaiting_receipt->refunding":
		return "同意仅退款"
	case "awaiting_return->awaiting_receipt":
		return "用户寄回商品"
	case "refunding->refunded", "->refunded":
		return "退款成功"
	case "applied->rejected", "merchant_handling->rejected", "platform_intervene->rejected":
		return "审核未通过"
	case "applied->cancelled", "awaiting_return->cancelled":
		return "用户取消退款"
	case "rejected->platform_intervene":
		return "申请平台介入"
	default:
		if fromStatus == "" {
			if toStatus == "applied" {
				return "提交退款申请"
			}
			if toStatus == "refunded" {
				return "退款成功"
			}
			return statusLabel(toStatus)
		}
		return statusLabel(fromStatus) + " → " + statusLabel(toStatus)
	}
}

func (h *Handler) applyListFilters(c *gin.Context, q *gorm.DB, ids []uint64) *gorm.DB {
	tab := strings.TrimSpace(c.Query("tab_status"))
	if tab == "" {
		tab = strings.TrimSpace(c.Query("status"))
	}
	if statuses := tabStatuses(tab); len(statuses) == 1 {
		q = q.Where("r.status = ?", statuses[0])
	} else if len(statuses) > 1 {
		q = q.Where("r.status IN ?", statuses)
	}

	if sn := strings.TrimSpace(c.Query("refund_order_sn")); sn != "" {
		q = q.Where("r.refund_no LIKE ?", "%"+sn+"%")
	}
	if orderSn := strings.TrimSpace(c.Query("order_sn")); orderSn != "" {
		q = q.Where("o.order_no LIKE ?", "%"+orderSn+"%")
	}
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		q = q.Where("r.created_at >= ?", from+" 00:00:00")
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		q = q.Where("r.created_at <= ?", to+" 23:59:59")
	}

	userType := strings.TrimSpace(c.Query("user_search_type"))
	userKeyword := strings.TrimSpace(c.Query("user_search_keyword"))
	if userKeyword == "" {
		// legacy single-field filters
		if nick := strings.TrimSpace(c.Query("nickname")); nick != "" {
			userType, userKeyword = "nickname", nick
		} else if phone := strings.TrimSpace(c.Query("phone")); phone != "" {
			userType, userKeyword = "phone", phone
		} else if uid := strings.TrimSpace(c.Query("uid")); uid != "" {
			userType, userKeyword = "uid", uid
		} else if realName := strings.TrimSpace(c.Query("real_name")); realName != "" {
			userType, userKeyword = "nickname", realName
		}
	}
	if userKeyword != "" {
		switch userType {
		case "uid", "user_id":
			if id, err := strconv.ParseUint(userKeyword, 10, 64); err == nil && id > 0 {
				q = q.Where("o.user_id = ?", id)
			} else {
				q = q.Where("1 = 0")
			}
		case "phone", "mobile":
			q = q.Where("EXISTS (SELECT 1 FROM qixi_crm_b_user u WHERE u.id = o.user_id AND u.mobile LIKE ?)", "%"+userKeyword+"%")
		default: // nickname
			q = q.Where("EXISTS (SELECT 1 FROM qixi_crm_b_user u WHERE u.id = o.user_id AND u.nickname LIKE ?)", "%"+userKeyword+"%")
		}
	}

	isTrader := strings.TrimSpace(c.Query("is_trader"))
	if isTrader == "0" || isTrader == "1" {
		traderIDs := h.traderMerchantIDs(c, isTrader == "1")
		if len(traderIDs) == 0 {
			q = q.Where("1 = 0")
		} else if ids != nil {
			filtered := intersectIDs(ids, traderIDs)
			if len(filtered) == 0 {
				q = q.Where("1 = 0")
			} else {
				q = q.Where("o.merchant_id IN ?", filtered)
			}
		} else {
			q = q.Where("o.merchant_id IN ?", traderIDs)
		}
	}
	return q
}

func (h *Handler) traderMerchantIDs(c *gin.Context, wantTrader bool) []uint64 {
	flag := 0
	if wantTrader {
		flag = 1
	}
	var rows []merchantViewRow
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
		Select("merchant_id,is_trader").Where("is_trader = ?", flag).Scan(&rows)
	out := make([]uint64, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.MerchantID)
	}
	return out
}

func intersectIDs(a, b []uint64) []uint64 {
	set := make(map[uint64]struct{}, len(b))
	for _, id := range b {
		set[id] = struct{}{}
	}
	out := make([]uint64, 0)
	for _, id := range a {
		if _, ok := set[id]; ok {
			out = append(out, id)
		}
	}
	return out
}

func (h *Handler) loadMerchantTraderMap(c *gin.Context, merchantIDs []uint64) map[uint64]int {
	out := map[uint64]int{}
	if len(merchantIDs) == 0 {
		return out
	}
	var rows []merchantViewRow
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
		Select("merchant_id,is_trader").Where("merchant_id IN ?", merchantIDs).Scan(&rows)
	for _, row := range rows {
		out[row.MerchantID] = row.IsTrader
	}
	return out
}

func (h *Handler) enrichRows(c *gin.Context, rows []refund, includeDetail bool) []gin.H {
	if len(rows) == 0 {
		return []gin.H{}
	}
	orderIDs := make([]uint64, 0, len(rows))
	refundIDs := make([]uint64, 0, len(rows))
	userIDs := make([]uint64, 0, len(rows))
	merchantIDs := make([]uint64, 0, len(rows))
	seenOrder := map[uint64]struct{}{}
	seenUser := map[uint64]struct{}{}
	seenMer := map[uint64]struct{}{}
	for _, row := range rows {
		refundIDs = append(refundIDs, row.ID)
		if _, ok := seenOrder[row.OrderID]; !ok {
			seenOrder[row.OrderID] = struct{}{}
			orderIDs = append(orderIDs, row.OrderID)
		}
		if row.UserID > 0 {
			if _, ok := seenUser[row.UserID]; !ok {
				seenUser[row.UserID] = struct{}{}
				userIDs = append(userIDs, row.UserID)
			}
		}
		if row.MerchantID > 0 {
			if _, ok := seenMer[row.MerchantID]; !ok {
				seenMer[row.MerchantID] = struct{}{}
				merchantIDs = append(merchantIDs, row.MerchantID)
			}
		}
	}

	orders := map[uint64]orderMeta{}
	var orderRows []orderMeta
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select("id,order_no,merchant_id,merchant_name_snapshot,store_id,store_name_snapshot,user_id,remark").
		Where("id IN ?", orderIDs).Scan(&orderRows)
	for _, row := range orderRows {
		orders[row.ID] = row
		if row.UserID > 0 {
			if _, ok := seenUser[row.UserID]; !ok {
				seenUser[row.UserID] = struct{}{}
				userIDs = append(userIDs, row.UserID)
			}
		}
		if row.MerchantID > 0 {
			if _, ok := seenMer[row.MerchantID]; !ok {
				seenMer[row.MerchantID] = struct{}{}
				merchantIDs = append(merchantIDs, row.MerchantID)
			}
		}
	}

	users := map[uint64]userRow{}
	if len(userIDs) > 0 {
		var userRows []userRow
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user").
			Select("id,nickname,mobile").Where("id IN ?", userIDs).Find(&userRows)
		for _, row := range userRows {
			users[row.ID] = row
		}
	}

	traderMap := h.loadMerchantTraderMap(c, merchantIDs)

	itemsByRefund := map[uint64][]gin.H{}
	refundNumByID := map[uint64]int{}
	var aftersaleItems []aftersaleItem
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_aftersale_item").
		Where("refund_id IN ?", refundIDs).Find(&aftersaleItems)
	orderItemIDs := make([]uint64, 0, len(aftersaleItems))
	for _, item := range aftersaleItems {
		orderItemIDs = append(orderItemIDs, item.OrderItemID)
		refundNumByID[item.RefundID] += item.Quantity
	}
	orderItems := map[uint64]orderItemRow{}
	if len(orderItemIDs) > 0 {
		var oiRows []orderItemRow
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order_item").
			Where("id IN ?", orderItemIDs).Find(&oiRows)
		for _, row := range oiRows {
			orderItems[row.ID] = row
		}
	}
	for _, item := range aftersaleItems {
		oi := orderItems[item.OrderItemID]
		payLine := oi.UnitPrice * float64(oi.Quantity)
		product := gin.H{
			"refund_product_id": item.ID,
			"order_product_id":  item.OrderItemID,
			"product_id":        oi.ProductID,
			"product_info":      emptyDash(oi.Title),
			"product_image":     oi.CoverURL,
			"product_sku":       specText(oi.Spec),
			"product_price":     oi.UnitPrice,
			"product_num":       oi.Quantity,
			"refund_num":        item.Quantity,
			"pay_price":         payLine,
			"refund_price":      item.Amount,
		}
		itemsByRefund[item.RefundID] = append(itemsByRefund[item.RefundID], product)
	}
	// Fallback: if no aftersale_item rows, surface order items for display (read-only).
	missingOrders := make([]uint64, 0)
	seenMissing := map[uint64]struct{}{}
	for _, row := range rows {
		if len(itemsByRefund[row.ID]) == 0 {
			if _, ok := seenMissing[row.OrderID]; !ok {
				seenMissing[row.OrderID] = struct{}{}
				missingOrders = append(missingOrders, row.OrderID)
			}
		}
	}
	fallbackByOrder := map[uint64][]gin.H{}
	if len(missingOrders) > 0 {
		var oiRows []orderItemRow
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order_item").
			Where("order_id IN ?", missingOrders).Find(&oiRows)
		for _, oi := range oiRows {
			payLine := oi.UnitPrice * float64(oi.Quantity)
			fallbackByOrder[oi.OrderID] = append(fallbackByOrder[oi.OrderID], gin.H{
				"refund_product_id": 0,
				"order_product_id":  oi.ID,
				"product_id":        oi.ProductID,
				"product_info":      emptyDash(oi.Title),
				"product_image":     oi.CoverURL,
				"product_sku":       specText(oi.Spec),
				"product_price":     oi.UnitPrice,
				"product_num":       oi.Quantity,
				"refund_num":        oi.Quantity,
				"pay_price":         payLine,
				"refund_price":      0,
			})
		}
	}

	evidences := map[uint64][]string{}
	userNotes := map[uint64]string{}
	type evidenceScan struct {
		RefundID    uint64          `gorm:"column:refund_id"`
		Content     string          `gorm:"column:content"`
		Attachments json.RawMessage `gorm:"column:attachments"`
	}
	// User remarks are needed on list expand; attachments only for detail.
	if len(refundIDs) > 0 {
		var evRows []evidenceScan
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_aftersale_evidence").
			Select("refund_id,content,attachments").
			Where("refund_id IN ? AND actor_type = ?", refundIDs, "user").
			Order("id ASC").Scan(&evRows)
		for _, row := range evRows {
			if strings.TrimSpace(row.Content) != "" && userNotes[row.RefundID] == "" {
				userNotes[row.RefundID] = row.Content
			}
			if includeDetail {
				evidences[row.RefundID] = append(evidences[row.RefundID], parseAttachments(row.Attachments)...)
			}
		}
	}

	initiators := map[uint64]string{}
	if includeDetail && len(refundIDs) > 0 {
		type firstEvent struct {
			RefundID  uint64 `gorm:"column:refund_id"`
			ActorType string `gorm:"column:actor_type"`
			ActorID   uint64 `gorm:"column:actor_id"`
		}
		var events []firstEvent
		_ = h.businessDB.WithContext(c.Request.Context()).Raw(`
			SELECT e.refund_id, e.actor_type, e.actor_id
			FROM qixi_crm_b_refund_event e
			INNER JOIN (
				SELECT refund_id, MIN(id) AS min_id
				FROM qixi_crm_b_refund_event
				WHERE refund_id IN ?
				GROUP BY refund_id
			) x ON x.min_id = e.id
		`, refundIDs).Scan(&events)
		actorUserIDs := make([]uint64, 0)
		for _, e := range events {
			if e.ActorType == "user" && e.ActorID > 0 {
				actorUserIDs = append(actorUserIDs, e.ActorID)
			}
		}
		actorUsers := map[uint64]string{}
		if len(actorUserIDs) > 0 {
			var rows []userRow
			_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user").
				Select("id,nickname").Where("id IN ?", actorUserIDs).Find(&rows)
			for _, row := range rows {
				actorUsers[row.ID] = row.Nickname
			}
		}
		for _, e := range events {
			switch e.ActorType {
			case "user":
				if name := actorUsers[e.ActorID]; name != "" {
					initiators[e.RefundID] = name
				} else {
					initiators[e.RefundID] = "用户"
				}
			case "merchant":
				initiators[e.RefundID] = "店铺"
			case "platform":
				initiators[e.RefundID] = "平台"
			default:
				initiators[e.RefundID] = actorRole(e.ActorType)
			}
		}
	}

	out := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		base := view(row)
		ord := orders[row.OrderID]
		uid := row.UserID
		if uid == 0 {
			uid = ord.UserID
		}
		u := users[uid]
		merName := ord.MerchantNameSnapshot
		if merName == "" {
			merName = "店铺 #" + strconv.FormatUint(row.MerchantID, 10)
		}
		storeName := ord.StoreNameSnapshot
		if storeName == "" {
			storeName = merName
		}
		storeCategory := "非自营"
		if traderMap[row.MerchantID] == 1 {
			storeCategory = "自营"
		}
		products := itemsByRefund[row.ID]
		if len(products) == 0 {
			products = fallbackByOrder[row.OrderID]
		}
		if products == nil {
			products = []gin.H{}
		}
		refundNum := refundNumByID[row.ID]
		if refundNum == 0 {
			for _, p := range products {
				refundNum += intVal(p["refund_num"])
			}
		}
		// For list money display on fallback products, attribute whole refund amount to first line.
		if len(itemsByRefund[row.ID]) == 0 && len(products) > 0 && row.Amount > 0 {
			products[0]["refund_price"] = row.Amount
		}

		productTotal := 0.0
		for _, p := range products {
			qty := float64(intVal(p["refund_num"]))
			if qty <= 0 {
				qty = float64(intVal(p["product_num"]))
			}
			if qty <= 0 {
				qty = 1
			}
			productTotal += floatVal(p["product_price"]) * qty
		}

		userDeleted := uid > 0
		if uid > 0 {
			if _, ok := users[uid]; ok {
				userDeleted = false
			}
		}
		nick := emptyDash(u.Nickname)
		if userDeleted {
			nick = "用户已被删除"
		}

		base["order_sn"] = emptyDash(ord.OrderNo)
		base["mer_name"] = merName
		base["store_name"] = storeName
		base["store_category_name"] = storeCategory
		base["is_trader"] = traderMap[row.MerchantID]
		base["uid"] = uid
		if userDeleted && uid == 0 {
			base["uid"] = -1
		}
		base["nickname"] = nick
		base["user_deleted"] = userDeleted
		base["user_phone_mask"] = maskPhone(u.Mobile)
		base["status_label"] = statusLabel(row.Status)
		base["refund_type_label"] = refundTypeLabel(row.RefundType)
		base["refund_num"] = refundNum
		base["product_total_price"] = productTotal
		base["refund_method"] = "原路返回"
		base["products"] = products
		base["user_remark"] = emptyDash(userNotes[row.ID])
		base["merchant_remark"] = "--"
		base["order_user_remark"] = emptyDash(ord.Remark)
		if includeDetail {
			base["refund_evidence"] = evidences[row.ID]
			if evidences[row.ID] == nil {
				base["refund_evidence"] = []string{}
			}
			initiator := initiators[row.ID]
			if initiator == "" {
				if userDeleted {
					initiator = "用户已被删除"
				} else {
					initiator = emptyDash(u.Nickname)
				}
			}
			base["refund_initiator"] = initiator
			var shipment returnShipment
			if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund_return_shipment").
				Where("refund_id = ?", row.ID).Scan(&shipment).Error; err == nil && shipment.TrackingNo != "" {
				base["return_shipment"] = gin.H{
					"carrier_name": shipment.CarrierName, "tracking_no": shipment.TrackingNo,
					"remark": shipment.Remark, "submitted_at": shipment.SubmittedAt.Format("2006-01-02 15:04:05"),
				}
			}
		}
		out = append(out, base)
	}
	return out
}

func intVal(v any) int {
	switch t := v.(type) {
	case int:
		return t
	case int64:
		return int(t)
	case float64:
		return int(t)
	case json.Number:
		n, _ := t.Int64()
		return int(n)
	case string:
		n, _ := strconv.Atoi(t)
		return n
	default:
		return 0
	}
}

func floatVal(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case json.Number:
		n, _ := t.Float64()
		return n
	case string:
		n, _ := strconv.ParseFloat(t, 64)
		return n
	default:
		return 0
	}
}

func operatorLabel(actorType string, actorID uint64, nicknames map[uint64]string) string {
	switch actorType {
	case "system":
		return "系统"
	case "user":
		if name := nicknames[actorID]; name != "" {
			return "ID:" + strconv.FormatUint(actorID, 10) + "，昵称：" + name
		}
		return "用户#" + strconv.FormatUint(actorID, 10)
	case "merchant":
		if name := nicknames[actorID]; name != "" {
			return "ID:" + strconv.FormatUint(actorID, 10) + "，昵称：" + name
		}
		return "店铺账号#" + strconv.FormatUint(actorID, 10)
	case "platform":
		return "平台管理员#" + strconv.FormatUint(actorID, 10)
	default:
		return strconv.FormatUint(actorID, 10)
	}
}
