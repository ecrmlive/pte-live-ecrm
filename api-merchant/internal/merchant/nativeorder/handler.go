// Package nativeorder serves store fulfillment from qixi_crm_b_* only.
package nativeorder

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	db       *gorm.DB
	identity *identity.Service
}

func NewHandler(db *gorm.DB, id *identity.Service) *Handler { return &Handler{db: db, identity: id} }
func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/orders", h.list)
	r.GET("/orders/:id", h.get)
	r.POST("/orders/:id/delivery", middleware.RequireMerchantMenu(h.identity, identity.MerPermOrderDeliver), h.deliver)
}

type order struct {
	ID            uint64          `gorm:"column:id"`
	GroupID       uint64          `gorm:"column:group_order_id"`
	OrderNo       string          `gorm:"column:order_no"`
	StoreID       uint64          `gorm:"column:store_id"`
	PayAmount     float64         `gorm:"column:pay_amount"`
	TotalQuantity int             `gorm:"column:total_quantity"`
	Recipient     json.RawMessage `gorm:"column:recipient_snapshot"`
	Remark        string          `gorm:"column:remark"`
	Status        string          `gorm:"column:status"`
	CreatedAt     time.Time       `gorm:"column:created_at"`
	PayChannel    string          `gorm:"column:pay_channel"`
}
type recipient struct {
	Recipient string `json:"recipient"`
	Mobile    string `json:"mobile"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Detail    string `json:"detail"`
}
type delivery struct {
	OrderID      uint64 `gorm:"column:order_id"`
	DeliveryType string `gorm:"column:delivery_type"`
	CarrierCode  string `gorm:"column:carrier_code"`
	TrackingNo   string `gorm:"column:tracking_no"`
}
type item struct {
	ID        uint64          `gorm:"column:id"`
	OrderID   uint64          `gorm:"column:order_id"`
	ProductID uint64          `gorm:"column:product_id"`
	Title     string          `gorm:"column:title_snapshot"`
	Spec      json.RawMessage `gorm:"column:spec_snapshot"`
	UnitPrice float64         `gorm:"column:unit_price"`
	Quantity  int             `gorm:"column:quantity"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := page(c)
	q := h.base(c)
	if paid := c.Query("paid"); paid == "0" {
		q = q.Where("o.status = 'pending_pay'")
	} else if paid == "1" {
		q = q.Where("o.status <> 'pending_pay'")
	}
	if status := c.Query("status"); status != "" {
		q = q.Where("o.status = ?", fromStatus(status))
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "查询订单失败")
		return
	}
	var rows []order
	if err := q.Order("o.created_at DESC,o.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "查询订单失败")
		return
	}
	result, err := h.responses(c, rows, false)
	if err != nil {
		fail(c, "加载订单失败")
		return
	}
	response.OK(c, gin.H{"list": result, "total": total, "page": page, "limit": limit})
}
func (h *Handler) get(c *gin.Context) {
	id := parseID(c)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "订单 ID 错误")
		return
	}
	var row order
	if err := h.base(c).Where("o.id = ?", id).Scan(&row).Error; err != nil {
		fail(c, "查询订单失败")
		return
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "订单不存在")
		return
	}
	result, err := h.responses(c, []order{row}, true)
	if err != nil {
		fail(c, "加载订单失败")
		return
	}
	response.OK(c, result[0])
}
func (h *Handler) deliver(c *gin.Context) {
	id := parseID(c)
	var req struct {
		DeliveryID   string `json:"delivery_id"`
		DeliveryName string `json:"delivery_name"`
		DeliveryType string `json:"delivery_type"`
	}
	if id == 0 || c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.DeliveryID) == "" || strings.TrimSpace(req.DeliveryName) == "" {
		response.Fail(c, http.StatusBadRequest, "发货参数不合法")
		return
	}
	if req.DeliveryType != "express" && req.DeliveryType != "local" {
		response.Fail(c, http.StatusBadRequest, "配送类型不合法")
		return
	}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var row order
		if err := tx.Table("qixi_crm_b_order AS o").Select("o.id,o.status").Where("o.id = ? AND o.store_id = ?", id, middleware.StoreID(c)).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&row).Error; err != nil {
			return err
		}
		if row.ID == 0 {
			return gorm.ErrRecordNotFound
		}
		if row.Status != "paid" && row.Status != "fulfilling" {
			return errStatus
		}
		if err := tx.Table("qixi_crm_b_order").Where("id = ?", id).Update("status", "shipped").Error; err != nil {
			return err
		}
		return tx.Exec(`INSERT INTO qixi_crm_b_order_delivery (order_id,delivery_type,carrier_code,tracking_no,status,delivered_at) VALUES (?, ?, ?, ?, 'shipped', NOW()) ON DUPLICATE KEY UPDATE delivery_type=VALUES(delivery_type),carrier_code=VALUES(carrier_code),tracking_no=VALUES(tracking_no),status='shipped',delivered_at=NOW()`, id, req.DeliveryType, strings.TrimSpace(req.DeliveryName), strings.TrimSpace(req.DeliveryID)).Error
	})
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, http.StatusNotFound, "订单不存在")
		return
	}
	if err == errStatus {
		response.Fail(c, http.StatusConflict, "当前订单不可发货")
		return
	}
	if err != nil {
		fail(c, "订单发货失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}
func (h *Handler) base(c *gin.Context) *gorm.DB {
	return h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_order AS o").Select("o.id,o.group_order_id,o.order_no,o.store_id,o.pay_amount,o.total_quantity,o.recipient_snapshot,o.remark,o.status,o.created_at,g.pay_channel").Joins("JOIN qixi_crm_b_group_order AS g ON g.id = o.group_order_id").Where("o.store_id = ?", middleware.StoreID(c))
}
func (h *Handler) responses(c *gin.Context, rows []order, details bool) ([]gin.H, error) {
	if len(rows) == 0 {
		return []gin.H{}, nil
	}
	ids := make([]uint64, 0, len(rows))
	for _, x := range rows {
		ids = append(ids, x.ID)
	}
	var deliveries []delivery
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_order_delivery").Where("order_id IN ?", ids).Find(&deliveries).Error; err != nil {
		return nil, err
	}
	deliveryByID := map[uint64]delivery{}
	for _, x := range deliveries {
		deliveryByID[x.OrderID] = x
	}
	itemsByID := map[uint64][]gin.H{}
	if details {
		var items []item
		if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_order_item").Where("order_id IN ?", ids).Find(&items).Error; err != nil {
			return nil, err
		}
		for _, x := range items {
			itemsByID[x.OrderID] = append(itemsByID[x.OrderID], gin.H{"order_product_id": x.ID, "product_id": x.ProductID, "product_info": x.Title, "product_sku": string(x.Spec), "product_price": x.UnitPrice, "product_num": x.Quantity, "total_price": x.UnitPrice * float64(x.Quantity)})
		}
	}
	out := make([]gin.H, 0, len(rows))
	for _, x := range rows {
		var r recipient
		_ = json.Unmarshal(x.Recipient, &r)
		d := deliveryByID[x.ID]
		out = append(out, gin.H{"order_id": x.ID, "order_sn": x.OrderNo, "paid": paid(x.Status), "status": status(x.Status), "pay_price": x.PayAmount, "pay_type": payType(x.PayChannel), "total_num": x.TotalQuantity, "real_name": r.Recipient, "user_phone": r.Mobile, "user_address": strings.TrimSpace(r.Province + r.City + r.District + " " + r.Detail), "mark": x.Remark, "delivery_type": d.DeliveryType, "delivery_name": d.CarrierCode, "delivery_id": d.TrackingNo, "create_time": x.CreatedAt.Format("2006-01-02 15:04:05"), "products": itemsByID[x.ID]})
	}
	return out, nil
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
func paid(v string) int {
	if v == "pending_pay" {
		return 0
	}
	return 1
}
func status(v string) int {
	switch v {
	case "shipped":
		return 1
	case "completed":
		return 3
	default:
		return 0
	}
}
func fromStatus(v string) string {
	switch v {
	case "1":
		return "shipped"
	case "3":
		return "completed"
	default:
		return "paid"
	}
}
func payType(v string) int {
	switch v {
	case "wechat":
		return 1
	case "alipay":
		return 2
	case "balance":
		return 0
	default:
		return 7
	}
}
func fail(c *gin.Context, msg string) { response.Fail(c, http.StatusInternalServerError, msg) }

var errStatus = fmtError("invalid order status")

type fmtError string

func (e fmtError) Error() string { return string(e) }
