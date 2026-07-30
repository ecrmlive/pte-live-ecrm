// Package nativeorder implements the platform order supervision read model.
//
// It deliberately reads the new business and merchant databases instead of the
// retired qixi_m_* projection.  Fulfilment, payment and after-sale mutations
// remain owned by their respective services; the platform order page is a
// supervision page and is read-only by design.
package nativeorder

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/response"
	"gorm.io/gorm"
)

type Handler struct {
	businessDB *gorm.DB
	merchantDB *gorm.DB
	identity   *identity.Service
}

func NewHandler(businessDB, merchantDB *gorm.DB, identity *identity.Service) *Handler {
	return &Handler{businessDB: businessDB, merchantDB: merchantDB, identity: identity}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/orders", h.list)
	r.GET("/orders/:id", h.get)
}

type order struct {
	ID            uint64          `gorm:"column:id"`
	GroupID       uint64          `gorm:"column:group_order_id"`
	OrderNo       string          `gorm:"column:order_no"`
	MerchantID    uint64          `gorm:"column:merchant_id"`
	StoreID       uint64          `gorm:"column:store_id"`
	PayAmount     float64         `gorm:"column:pay_amount"`
	TotalAmount   float64         `gorm:"column:total_amount"`
	TotalQuantity int             `gorm:"column:total_quantity"`
	Recipient     json.RawMessage `gorm:"column:recipient_snapshot"`
	Status        string          `gorm:"column:status"`
	CreatedAt     time.Time       `gorm:"column:created_at"`
	PaidAt        *time.Time      `gorm:"column:paid_at"`
	PayChannel    string          `gorm:"column:pay_channel"`
}

type orderItem struct {
	ID        uint64          `gorm:"column:id"`
	OrderID   uint64          `gorm:"column:order_id"`
	ProductID uint64          `gorm:"column:product_id"`
	Title     string          `gorm:"column:title_snapshot"`
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
		response.Fail(c, http.StatusUnauthorized, "登录已失效")
		return
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.OK(c, gin.H{"list": []gin.H{}, "total": 0, "page": page, "limit": limit})
		return
	}

	q := h.base(c, merchantIDs)
	if paid := c.Query("paid"); paid == "0" {
		q = q.Where("o.status = 'pending_pay'")
	} else if paid == "1" {
		q = q.Where("o.status <> 'pending_pay'")
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
	items, err := h.responses(c, rows, false)
	if err != nil {
		fail(c, "加载订单失败")
		return
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func (h *Handler) get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "订单 ID 错误")
		return
	}
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "登录已失效")
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
	if err != nil {
		fail(c, "加载订单失败")
		return
	}
	response.OK(c, items[0])
}

func (h *Handler) base(c *gin.Context, merchantIDs []uint64) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_order AS o").
		Select("o.id,o.group_order_id,o.order_no,o.merchant_id,o.store_id,o.pay_amount,o.total_amount,o.total_quantity,o.recipient_snapshot,o.status,o.created_at,o.paid_at,g.pay_channel").
		Joins("JOIN qixi_crm_b_group_order AS g ON g.id = o.group_order_id")
	if merchantIDs != nil {
		q = q.Where("o.merchant_id IN ?", merchantIDs)
	}
	return q
}

// merchantScope maps the unified-admin region assignment to current merchant
// records. nil means the platform account has full access; an empty slice means
// a region account is intentionally denied until it receives an assignment.
func (h *Handler) merchantScope(c *gin.Context) ([]uint64, error) {
	regionIDs, err := h.identity.PlatformRegionScope(c.Request.Context(), middleware.AdminID(c))
	if err != nil || regionIDs == nil {
		return nil, err
	}
	if len(regionIDs) == 0 {
		return []uint64{}, nil
	}
	var rows []merchant
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_merchant").
		Select("id,name,region_id").Where("region_id IN ?", regionIDs).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

func (h *Handler) responses(c *gin.Context, rows []order, includeItems bool) ([]gin.H, error) {
	if len(rows) == 0 {
		return []gin.H{}, nil
	}
	orderIDs := make([]uint64, 0, len(rows))
	merchantIDs := make([]uint64, 0, len(rows))
	storeIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		orderIDs = append(orderIDs, row.ID)
		merchantIDs = append(merchantIDs, row.MerchantID)
		storeIDs = append(storeIDs, row.StoreID)
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
	var deliveries []delivery
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order_delivery").Where("order_id IN ?", orderIDs).Find(&deliveries).Error; err != nil {
		return nil, err
	}
	deliveriesByOrder := map[uint64]delivery{}
	for _, row := range deliveries {
		deliveriesByOrder[row.OrderID] = row
	}
	itemsByOrder := map[uint64][]gin.H{}
	if includeItems {
		var items []orderItem
		if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order_item").Where("order_id IN ?", orderIDs).Find(&items).Error; err != nil {
			return nil, err
		}
		for _, row := range items {
			itemsByOrder[row.OrderID] = append(itemsByOrder[row.OrderID], gin.H{
				"order_product_id": row.ID, "product_id": row.ProductID, "product_info": row.Title,
				"product_sku": string(row.Spec), "product_price": row.UnitPrice, "product_num": row.Quantity,
				"total_price": row.UnitPrice * float64(row.Quantity),
			})
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
		out = append(out, gin.H{
			"order_id": row.ID, "order_sn": row.OrderNo, "mer_id": row.MerchantID, "mer_name": merchantNames[row.MerchantID],
			"store_id": row.StoreID, "store_name": storeNames[row.StoreID], "paid": paid(row.Status), "status": orderStatus(row.Status),
			"pay_price": row.PayAmount, "total_price": row.TotalAmount, "total_num": row.TotalQuantity, "pay_type": payType(row.PayChannel),
			"pay_time": payTime, "delivery_type": d.DeliveryType, "delivery_name": d.CarrierCode, "delivery_id": d.TrackingNo,
			"user_phone": address.Mobile, "user_address": strings.TrimSpace(address.Province + address.City + address.District + " " + address.Detail),
			"create_time": row.CreatedAt.Format("2006-01-02 15:04:05"), "products": itemsByOrder[row.ID],
		})
	}
	return out, nil
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
func orderStatus(status string) int {
	switch status {
	case "shipped":
		return 1
	case "completed":
		return 3
	case "cancelled":
		return -1
	default:
		return 0
	}
}
func payType(channel string) int {
	switch channel {
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
func fail(c *gin.Context, message string) { response.Fail(c, http.StatusInternalServerError, message) }
