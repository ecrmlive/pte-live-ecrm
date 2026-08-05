package order

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/cloudconfig"
	merchantstock "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantstock"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/paymentconfig"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/wechatpayv3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	db              *gorm.DB
	configs         *paymentconfig.Store
	platformConfigs *cloudconfig.Service
	wechatClient    *wechatpayv3.Client
	allowMock       bool
}

func NewHandler(db *gorm.DB, configs *paymentconfig.Store, allowMock bool, platformConfigs ...*cloudconfig.Service) *Handler {
	var platformConfig *cloudconfig.Service
	if len(platformConfigs) > 0 {
		platformConfig = platformConfigs[0]
	}
	return &Handler{db: db, configs: configs, platformConfigs: platformConfig, wechatClient: &wechatpayv3.Client{}, allowMock: allowMock}
}
func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/v2/order/check", h.Check)
	r.POST("/v2/order/create", h.Create)
	r.POST("/order/pay/:id", h.Pay)
	r.GET("/order/pay/:id/channels", h.PaymentChannels)
	r.GET("/orders", h.List)
	r.GET("/orders/:id", h.GetGroup)
	r.POST("/orders/:id/cancel", h.Cancel)
	r.DELETE("/orders/:id", h.Archive)
	r.GET("/order/:id", h.GetOrder)
	r.GET("/order/:id/delivery", h.Delivery)
	r.POST("/order/:id/confirm-receipt", h.ConfirmReceipt)
}

type checkRequest struct {
	CartIDs       []uint64 `json:"cart_ids"`
	CouponUserIDs []uint64 `json:"coupon_user_ids"`
	UseIntegral   bool     `json:"use_integral"`
}
type createRequest struct {
	CartIDs        []uint64 `json:"cart_ids"`
	AddressID      uint64   `json:"address_id"`
	Mark           string   `json:"mark"`
	CouponUserIDs  []uint64 `json:"coupon_user_ids"`
	UseIntegral    bool     `json:"use_integral"`
	IdempotencyKey string   `json:"idempotency_key"`
}

func (h *Handler) Check(c *gin.Context) {
	var req checkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		bad(c, "参数错误")
		return
	}
	checkout, err := LoadCheckout(c.Request.Context(), h.db, uint64(middleware.UID(c)), req.CartIDs)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	pricing, err := ResolveCoupons(c.Request.Context(), h.db, uint64(middleware.UID(c)), checkout, req.CouponUserIDs, false)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	quote, err := ResolveIntegral(c.Request.Context(), h.db, uint64(middleware.UID(c)), checkout, pricing.DiscountCents, req.UseIntegral, false)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, checkResponse(checkout, pricing, quote))
}
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		bad(c, "参数错误")
		return
	}
	req.IdempotencyKey = strings.TrimSpace(req.IdempotencyKey)
	if !validIdempotencyKey(req.IdempotencyKey) {
		bad(c, ErrIdempotencyKey.Error())
		return
	}
	uid := uint64(middleware.UID(c))
	created, err := Create(c.Request.Context(), h.db, uid, CreateInput{CartIDs: req.CartIDs, AddressID: req.AddressID, CouponUserIDs: req.CouponUserIDs, UseIntegral: req.UseIntegral, IdempotencyKey: req.IdempotencyKey, Remark: req.Mark})
	if err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, createdResponse(created, false))
}
func (h *Handler) Pay(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		bad(c, "订单 ID 错误")
		return
	}
	var req struct {
		PayType string `json:"pay_type"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.PayType == "" {
		bad(c, "必须选择支付方式")
		return
	}
	if req.PayType == "wechat" {
		userID := uint64(middleware.UID(c))
		claims := middleware.ClaimsFrom(c)
		var intent wechatPaymentIntent
		if claims != nil && claims.ClientPlatform == "mini_program" {
			intent, err = h.createWechatJSAPIPay(c.Request.Context(), userID, id)
		} else if claims != nil && claims.ClientPlatform == "h5" {
			intent, err = h.createWechatH5Pay(c.Request.Context(), userID, id, c.ClientIP())
		} else {
			intent, err = h.createWechatNativePay(c.Request.Context(), userID, id)
		}
		if err != nil {
			writeOrderError(c, err)
			return
		}
		response.OK(c, intent)
		return
	}
	if req.PayType == "balance" {
		created, err := PayBalance(c.Request.Context(), h.db, uint64(middleware.UID(c)), id)
		if err != nil {
			writeOrderError(c, err)
			return
		}
		response.OK(c, createdResponse(created, true))
		return
	}
	if req.PayType != "mock" {
		writeOrderError(c, ErrPayChannel)
		return
	}
	if !h.allowMock {
		writeOrderError(c, ErrPayChannel)
		return
	}
	created, err := PayMock(c.Request.Context(), h.db, uint64(middleware.UID(c)), id)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, createdResponse(created, true))
}
func (h *Handler) PaymentChannels(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		bad(c, "订单 ID 错误")
		return
	}
	if h.configs == nil {
		response.OK(c, gin.H{"list": []PaymentChannelView{{Channel: "balance", Enabled: true}, {Channel: "wechat"}, {Channel: "alipay"}}})
		return
	}
	list, err := availablePaymentChannels(c.Request.Context(), h.db, h.configs, h.platformConfigs, uint64(middleware.UID(c)), id)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, gin.H{"list": list})
}
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	payStatus, err := NormalizeGroupPayStatus(c.Query("pay_status"))
	if err != nil {
		bad(c, "订单状态筛选参数错误")
		return
	}
	fulfillmentStatus, err := NormalizeGroupFulfillmentStatus(c.Query("fulfillment_status"))
	if err != nil {
		bad(c, "订单履约状态筛选参数错误")
		return
	}
	if payStatus != "" && fulfillmentStatus != "" {
		bad(c, "付款状态与履约状态不能同时筛选")
		return
	}
	keyword, err := NormalizeOrderListKeyword(c.Query("keyword"))
	if err != nil {
		bad(c, "订单搜索关键词不合法")
		return
	}
	rows, total, err := ListGroups(c.Request.Context(), h.db, uint64(middleware.UID(c)), payStatus, fulfillmentStatus, keyword, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "订单查询失败")
		return
	}
	summaries, err := loadGroupListSummaries(c.Request.Context(), h.db, uint64(middleware.UID(c)), rows)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "订单状态汇总失败")
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		out := groupResponse(row)
		summary := summaries[row.ID]
		out["fulfillment_status"] = summary.FulfillmentStatus
		out["has_uncommented_item"] = summary.HasUncommentedItem
		list = append(list, out)
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}
func (h *Handler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		bad(c, "订单 ID 错误")
		return
	}
	detail, err := GetGroup(c.Request.Context(), h.db, uint64(middleware.UID(c)), id)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	deliveries, err := deliverySummaries(c.Request.Context(), h.db, detail.Orders)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "配送信息查询失败")
		return
	}
	reservations, err := reservationSummaries(c.Request.Context(), h.db, detail.Orders)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "预约信息查询失败")
		return
	}
	orders := make([]gin.H, 0, len(detail.Orders))
	for _, order := range detail.Orders {
		items := make([]gin.H, 0, len(detail.Items[order.ID]))
		for _, item := range detail.Items[order.ID] {
			items = append(items, orderItemResponse(item))
		}
		out := orderResponse(order, items)
		if delivery, ok := deliveries[order.ID]; ok {
			for key, value := range delivery {
				out[key] = value
			}
		}
		if reservation, ok := reservations[order.ID]; ok {
			for key, value := range reservation {
				out[key] = value
			}
		}
		orders = append(orders, out)
	}
	out := groupResponse(detail.Group)
	out["orders"] = orders
	out["can_archive"] = CanArchiveGroup(detail.Group, detail.Orders)
	response.OK(c, out)
}
func (h *Handler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		bad(c, "订单 ID 错误")
		return
	}
	if err := CancelPending(c.Request.Context(), h.db, uint64(middleware.UID(c)), id); err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

type deliveryRow struct {
	OrderID      uint64     `gorm:"column:order_id"`
	DeliveryType string     `gorm:"column:delivery_type"`
	CarrierCode  string     `gorm:"column:carrier_code"`
	TrackingNo   string     `gorm:"column:tracking_no"`
	Status       string     `gorm:"column:status"`
	DeliveredAt  *time.Time `gorm:"column:delivered_at"`
}

func deliverySummary(row deliveryRow) gin.H {
	name := row.CarrierCode
	if name == "" {
		name = "快递配送"
	}
	return gin.H{"delivery_name": name, "delivery_id": row.TrackingNo, "delivery_type": row.DeliveryType, "delivery_status": row.Status, "delivered_at": row.DeliveredAt}
}

func (h *Handler) ConfirmReceipt(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		bad(c, "订单 ID 错误")
		return
	}
	uid := uint64(middleware.UID(c))
	err = h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var order orderRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ?", id, uid).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOrderOwnership
			}
			return err
		}
		if order.Status == "completed" {
			return nil
		}
		if order.Status != "shipped" {
			return ErrOrderNotReceivable
		}
		result := tx.Model(&orderRow{}).Where("id = ? AND user_id = ? AND status = ?", id, uid, "shipped").Update("status", "completed")
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrOrderNotReceivable
		}
		return enqueueSettlementAccrual(tx, order)
	})
	if err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, gin.H{"order_id": id, "status": "completed", "received": true})
}

func (h *Handler) Delivery(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		bad(c, "订单 ID 错误")
		return
	}
	if _, _, err = GetOrder(c.Request.Context(), h.db, uint64(middleware.UID(c)), id); err != nil {
		writeOrderError(c, err)
		return
	}
	var row deliveryRow
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_order_delivery").Where("order_id = ?", id).Order("id DESC").First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "订单暂无配送信息")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "配送信息查询失败")
		return
	}
	response.OK(c, deliverySummary(row))
}

func (h *Handler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		bad(c, "订单 ID 错误")
		return
	}
	order, items, err := GetOrder(c.Request.Context(), h.db, uint64(middleware.UID(c)), id)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	out := make([]gin.H, 0, len(items))
	for _, item := range items {
		out = append(out, orderItemResponse(item))
	}
	deliveries, derr := deliverySummaries(c.Request.Context(), h.db, []orderRow{order})
	if derr != nil {
		response.Fail(c, http.StatusInternalServerError, "配送信息查询失败")
		return
	}
	reservations, rerr := reservationSummaries(c.Request.Context(), h.db, []orderRow{order})
	if rerr != nil {
		response.Fail(c, http.StatusInternalServerError, "预约信息查询失败")
		return
	}
	responseBody := orderResponse(order, out)
	if delivery, ok := deliveries[order.ID]; ok {
		for key, value := range delivery {
			responseBody[key] = value
		}
	}
	if reservation, ok := reservations[order.ID]; ok {
		for key, value := range reservation {
			responseBody[key] = value
		}
	}
	response.OK(c, responseBody)
}

func orderItemResponse(item orderItemRow) gin.H {
	return gin.H{"order_item_id": item.ID, "product_id": item.ProductID, "product_attr_unique": item.SKUKey, "store_name": item.TitleSnapshot, "image": item.CoverURLSnapshot, "spec_text": specText(item.SpecSnapshot), "product_num": item.Quantity, "product_price": item.UnitPrice, "commented": item.Commented}
}

func checkResponse(checkout Checkout, pricing CouponPricing, quotes ...IntegralQuote) gin.H {
	quote := IntegralQuote{}
	if len(quotes) > 0 {
		quote = quotes[0]
	}
	merchants := make([]gin.H, 0, len(checkout.Stores))
	for _, store := range checkout.Stores {
		items := make([]gin.H, 0, len(store.Lines))
		for _, line := range store.Lines {
			items = append(items, gin.H{"cart_id": line.CartID, "product_id": line.ProductID, "product_attr_unique": line.SKUKey, "store_name": line.Title, "image": line.CoverURL, "spec_text": specText(line.SpecSnapshot), "price": money(line.UnitCents), "original_price": money(line.ListCents), "svip_applied": line.SVIPApplied, "cart_num": line.Quantity, "subtotal": money(line.UnitCents * int64(line.Quantity))})
		}
		merchants = append(merchants, gin.H{"mer_id": store.MerchantID, "mer_name": store.MerchantName, "total_price": money(store.TotalCents), "pay_price": money(store.TotalCents - pricing.DiscountCents - quote.DiscountCents), "coupon_price": money(pricing.DiscountCents), "total_num": store.TotalQty, "items": items})
	}
	return gin.H{"merchants": merchants, "total_price": money(checkout.TotalCents), "pay_price": money(checkout.TotalCents - pricing.DiscountCents - quote.DiscountCents), "total_num": checkout.TotalQty, "total_postage": 0, "coupon_price": money(pricing.DiscountCents), "integral_price": money(quote.DiscountCents), "integral": quote.PointsUsed, "integral_enabled": len(checkout.Stores) == 1 && checkout.Stores[0].IntegralPolicy.Enabled}
}
func createdResponse(created CreatedOrder, paid bool) gin.H {
	return gin.H{"group_order_id": created.GroupOrderID, "group_order_sn": created.GroupOrderNo, "pay_price": money(created.PayCents), "total_num": created.TotalQuantity, "pay_status": map[bool]string{false: "pending", true: "paid"}[paid], "paid": map[bool]int{false: 0, true: 1}[paid], "pay_type": 7}
}
func groupResponse(group groupRow) gin.H {
	return gin.H{"group_order_id": group.ID, "group_order_sn": group.OrderNo, "pay_price": group.PayAmount, "total_num": group.TotalQuantity, "pay_status": group.PayStatus, "paid": map[bool]int{false: 0, true: 1}[group.PayStatus == "paid"], "activity_type": group.ActivityType, "points_amount": group.PointsAmount, "create_time": group.CreatedAt.Format("2006-01-02 15:04:05")}
}
func orderResponse(order orderRow, items []gin.H) gin.H {
	return gin.H{"order_id": order.ID, "order_sn": order.OrderNo, "mer_id": order.MerchantID, "mer_name": order.MerchantNameSnapshot, "pay_price": order.PayAmount, "total_num": order.TotalQuantity, "status": order.Status, "products": items, "activity_type": order.ActivityType, "points_amount": order.PointsAmount}
}
func bad(c *gin.Context, message string) { response.Fail(c, http.StatusBadRequest, message) }
func writeOrderError(c *gin.Context, err error) {
	switch err {
	case ErrEmptyCart, ErrUnavailableCart, ErrMixedActivity, ErrMixedPaySubject, ErrAddressOwnership, ErrIdempotencyKey, ErrCartOwnership, ErrPayChannel, ErrStoreChannelDisabled, ErrCouponOwnership, ErrCouponConflict, ErrCouponMinNotMet, ErrOrderRemark, ErrOrderNotCancellable, ErrOrderNotPayable, ErrInsufficientBalance, ErrOrderNotReceivable, ErrOrderNotArchivable, ErrH5PaymentConfig, ErrPaymentProcessing, merchantstock.ErrReservationsPending, merchantstock.ErrReservationsFailed:
		bad(c, err.Error())
	case ErrOrderOwnership:
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "订单服务异常")
	}
}

func deliverySummaries(ctx context.Context, db *gorm.DB, orders []orderRow) (map[uint64]gin.H, error) {
	ids := make([]uint64, 0, len(orders))
	for _, order := range orders {
		ids = append(ids, order.ID)
	}
	out := make(map[uint64]gin.H)
	if len(ids) == 0 {
		return out, nil
	}
	var rows []deliveryRow
	if err := db.WithContext(ctx).Table("qixi_crm_b_order_delivery").Where("order_id IN ?", ids).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		if _, exists := out[row.OrderID]; !exists {
			out[row.OrderID] = deliverySummary(row)
		}
	}
	return out, nil
}

type reservationDetailRow struct {
	OrderID         uint64    `gorm:"column:order_id"`
	ReservationDate time.Time `gorm:"column:reservation_date"`
	StartTime       string    `gorm:"column:start_time"`
	EndTime         string    `gorm:"column:end_time"`
	VerifyCode      string    `gorm:"column:verify_code"`
}

func reservationSummaries(ctx context.Context, db *gorm.DB, orders []orderRow) (map[uint64]gin.H, error) {
	ids := make([]uint64, 0, len(orders))
	for _, order := range orders {
		if order.ActivityType == reservationOrderActivityType {
			ids = append(ids, order.ID)
		}
	}
	out := make(map[uint64]gin.H)
	if len(ids) == 0 {
		return out, nil
	}
	var rows []reservationDetailRow
	if err := db.WithContext(ctx).Table("qixi_crm_b_reservation_booking AS b").
		Select("b.order_id,b.booking_date AS reservation_date,b.verify_code,s.start_time,s.end_time").
		Joins("JOIN qixi_crm_b_reservation_slot AS s ON s.attr_reservation_id = b.slot_id").
		Where("b.order_id IN ? AND b.status = 1", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.OrderID] = gin.H{"reservation_date": row.ReservationDate.Format("2006-01-02"), "reservation_time_part": row.StartTime + "-" + row.EndTime, "verify_code": row.VerifyCode}
	}
	return out, nil
}
