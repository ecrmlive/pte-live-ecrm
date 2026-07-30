package order

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/paymentconfig"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

type Handler struct {
	db        *gorm.DB
	configs   *paymentconfig.Store
	allowMock bool
}

func NewHandler(db *gorm.DB, configs *paymentconfig.Store, allowMock bool) *Handler {
	return &Handler{db: db, configs: configs, allowMock: allowMock}
}
func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/v2/order/check", h.Check)
	r.POST("/v2/order/create", h.Create)
	r.POST("/order/pay/:id", h.Pay)
	r.GET("/order/pay/:id/channels", h.PaymentChannels)
	r.GET("/orders", h.List)
	r.GET("/orders/:id", h.GetGroup)
	r.GET("/order/:id", h.GetOrder)
}

type checkRequest struct {
	CartIDs []uint64 `json:"cart_ids"`
}
type createRequest struct {
	CartIDs        []uint64 `json:"cart_ids"`
	AddressID      uint64   `json:"address_id"`
	Mark           string   `json:"mark"`
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
	response.OK(c, checkResponse(checkout))
}
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		bad(c, "参数错误")
		return
	}
	uid := uint64(middleware.UID(c))
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		req.IdempotencyKey = derivedKey(uid, req)
	}
	created, err := Create(c.Request.Context(), h.db, uid, CreateInput{CartIDs: req.CartIDs, AddressID: req.AddressID, IdempotencyKey: req.IdempotencyKey, Remark: req.Mark})
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
	if req.PayType != "mock" {
		if h.configs == nil {
			writeOrderError(c, ErrStoreChannelDisabled)
			return
		}
		if err := AssertPaymentChannelAvailable(c.Request.Context(), h.db, h.configs, uint64(middleware.UID(c)), id, req.PayType); err != nil {
			writeOrderError(c, err)
			return
		}
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
		response.OK(c, gin.H{"list": []PaymentChannelView{{Channel: "wechat"}, {Channel: "alipay"}}})
		return
	}
	list, err := AvailablePaymentChannels(c.Request.Context(), h.db, h.configs, uint64(middleware.UID(c)), id)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, gin.H{"list": list})
}
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	rows, total, err := ListGroups(c.Request.Context(), h.db, uint64(middleware.UID(c)), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "订单查询失败")
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, groupResponse(row))
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
	orders := make([]gin.H, 0, len(detail.Orders))
	for _, order := range detail.Orders {
		items := make([]gin.H, 0, len(detail.Items[order.ID]))
		for _, item := range detail.Items[order.ID] {
			items = append(items, gin.H{"product_id": item.ProductID, "product_attr_unique": item.SKUKey, "store_name": item.TitleSnapshot, "image": item.CoverURLSnapshot, "product_num": item.Quantity, "product_price": item.UnitPrice})
		}
		orders = append(orders, orderResponse(order, items))
	}
	out := groupResponse(detail.Group)
	out["orders"] = orders
	response.OK(c, out)
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
		out = append(out, gin.H{"product_id": item.ProductID, "product_attr_unique": item.SKUKey, "store_name": item.TitleSnapshot, "image": item.CoverURLSnapshot, "product_num": item.Quantity, "product_price": item.UnitPrice})
	}
	response.OK(c, orderResponse(order, out))
}

func checkResponse(checkout Checkout) gin.H {
	merchants := make([]gin.H, 0, len(checkout.Stores))
	for _, store := range checkout.Stores {
		items := make([]gin.H, 0, len(store.Lines))
		for _, line := range store.Lines {
			items = append(items, gin.H{"cart_id": line.CartID, "product_id": line.ProductID, "product_attr_unique": line.SKUKey, "store_name": line.Title, "image": line.CoverURL, "price": money(line.UnitCents), "cart_num": line.Quantity, "subtotal": money(line.UnitCents * int64(line.Quantity))})
		}
		merchants = append(merchants, gin.H{"mer_id": store.MerchantID, "mer_name": store.MerchantName, "total_price": money(store.TotalCents), "pay_price": money(store.TotalCents), "total_num": store.TotalQty, "items": items})
	}
	return gin.H{"merchants": merchants, "total_price": money(checkout.TotalCents), "pay_price": money(checkout.TotalCents), "total_num": checkout.TotalQty, "total_postage": 0, "coupon_price": 0}
}
func createdResponse(created CreatedOrder, paid bool) gin.H {
	return gin.H{"group_order_id": created.GroupOrderID, "group_order_sn": created.GroupOrderNo, "pay_price": money(created.PayCents), "total_num": created.TotalQuantity, "paid": map[bool]int{false: 0, true: 1}[paid], "pay_type": 7}
}
func groupResponse(group groupRow) gin.H {
	return gin.H{"group_order_id": group.ID, "group_order_sn": group.OrderNo, "pay_price": group.PayAmount, "total_num": group.TotalQuantity, "paid": map[bool]int{false: 0, true: 1}[group.PayStatus == "paid"]}
}
func orderResponse(order orderRow, items []gin.H) gin.H {
	return gin.H{"order_id": order.ID, "order_sn": order.OrderNo, "mer_id": order.MerchantID, "mer_name": order.MerchantNameSnapshot, "pay_price": order.PayAmount, "total_num": order.TotalQuantity, "status": order.Status, "products": items}
}
func derivedKey(uid uint64, req createRequest) string {
	return fmt.Sprintf("auto:%x", sha256.Sum256([]byte(fmt.Sprintf("%d:%d:%v:%s", uid, req.AddressID, req.CartIDs, req.Mark))))
}
func bad(c *gin.Context, message string) { response.Fail(c, http.StatusBadRequest, message) }
func writeOrderError(c *gin.Context, err error) {
	switch err {
	case ErrEmptyCart, ErrUnavailableCart, ErrMixedActivity, ErrMixedPaySubject, ErrAddressOwnership, ErrIdempotencyKey, ErrCartOwnership, ErrPayChannel, ErrStoreChannelDisabled:
		bad(c, err.Error())
	case ErrOrderOwnership:
		response.Fail(c, http.StatusNotFound, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "订单服务异常")
	}
}
