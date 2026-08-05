package order

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	merchantstock "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantstock"
	"gorm.io/gorm"
)

var (
	ErrAddressOwnership = errors.New("收货地址不存在或无权访问")
	ErrIdempotencyKey   = errors.New("缺少幂等键")
	ErrOrderRemark      = errors.New("订单备注不能超过200个字符")
)

type CreateInput struct {
	CartIDs        []uint64
	AddressID      uint64
	CouponUserIDs  []uint64
	UseIntegral    bool
	IdempotencyKey string
	Remark         string
	DeliveryType   string
}
type CreatedOrder struct {
	GroupOrderID  uint64
	GroupOrderNo  string
	PayCents      int64
	TotalQuantity int
}

func validIdempotencyKey(value string) bool { return len(value) >= 12 && len(value) <= 128 }

func normalizeOrderRemark(value string) (string, error) {
	value = strings.TrimSpace(value)
	if utf8.RuneCountInString(value) > 200 {
		return "", ErrOrderRemark
	}
	return value, nil
}

func Create(ctx context.Context, db *gorm.DB, userID uint64, input CreateInput) (CreatedOrder, error) {
	input.IdempotencyKey = strings.TrimSpace(input.IdempotencyKey)
	remark, err := normalizeOrderRemark(input.Remark)
	if err != nil {
		return CreatedOrder{}, err
	}
	input.Remark = remark
	deliveryType, err := normalizeDeliveryType(input.DeliveryType)
	if err != nil {
		return CreatedOrder{}, err
	}
	input.DeliveryType = deliveryType
	if !validIdempotencyKey(input.IdempotencyKey) {
		return CreatedOrder{}, ErrIdempotencyKey
	}
	if input.AddressID == 0 {
		return CreatedOrder{}, ErrAddressOwnership
	}
	var existing groupRow
	err = db.WithContext(ctx).Where("user_id = ? AND idempotency_key = ?", userID, input.IdempotencyKey).First(&existing).Error
	if err == nil {
		return CreatedOrder{GroupOrderID: existing.ID, GroupOrderNo: existing.OrderNo, PayCents: int64(existing.PayAmount * 100), TotalQuantity: existing.TotalQuantity}, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return CreatedOrder{}, err
	}
	checkout, err := LoadCheckout(ctx, db, userID, input.CartIDs)
	if err != nil {
		return CreatedOrder{}, err
	}
	address, err := loadAddress(ctx, db, userID, input.AddressID)
	if err != nil {
		return CreatedOrder{}, err
	}
	addressJSON, err := json.Marshal(address)
	if err != nil {
		return CreatedOrder{}, err
	}
	created := CreatedOrder{}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		pricing, err := ResolveCoupons(ctx, tx, userID, checkout, input.CouponUserIDs, true)
		if err != nil {
			return err
		}
		// 显式写入状态
		quote, err := ResolveIntegral(ctx, tx, userID, checkout, pricing.DiscountCents, input.UseIntegral, true)
		if err != nil {
			return err
		}
		// 显式写入状态而不依赖数据库默认值。GORM 会把 Go 的零值字符串
		// 一并写入，MySQL 严格模式下 enum('pending', ...) 不能接受空字符串。
		group := groupRow{OrderNo: orderNo("G"), UserID: userID, TotalAmount: money(checkout.TotalCents), DiscountAmount: money(pricing.DiscountCents + quote.DiscountCents), PayAmount: money(checkout.TotalCents - pricing.DiscountCents - quote.DiscountCents), PointsAmount: quote.PointsUsed, PayStatus: "pending", TotalQuantity: checkout.TotalQty, RecipientSnapshot: string(addressJSON), IdempotencyKey: input.IdempotencyKey, Remark: strings.TrimSpace(input.Remark)}
		if err := tx.Create(&group).Error; err != nil {
			return err
		}
		if err := deductIntegral(tx, userID, group, quote); err != nil {
			return err
		}
		if err := lockCouponsToOrder(tx, userID, group.ID, pricing); err != nil {
			return err
		}
		reservationExpiry := time.Now().UTC().Add(30 * time.Minute)
		for _, store := range checkout.Stores {
			order := orderRow{GroupOrderID: group.ID, OrderNo: orderNo("S"), MerchantID: store.MerchantID, MerchantNameSnapshot: store.MerchantName, StoreID: store.StoreID, StoreNameSnapshot: store.StoreName, UserID: userID, TotalAmount: money(store.TotalCents), DiscountAmount: money(pricing.DiscountCents + quote.DiscountCents), PayAmount: money(store.TotalCents - pricing.DiscountCents - quote.DiscountCents), PointsAmount: quote.PointsUsed, TotalQuantity: store.TotalQty, RecipientSnapshot: string(addressJSON), Remark: group.Remark, Status: "pending_pay"}
			if err := tx.Create(&order).Error; err != nil {
				return err
			}
			for _, line := range store.Lines {
				if line.MerchantSKUID == 0 {
					return ErrUnavailableCart
				}
				item := orderItemRow{OrderID: order.ID, ProductID: line.ProductID, MerchantSKUID: line.MerchantSKUID, SKUKey: line.SKUKey, TitleSnapshot: line.Title, CoverURLSnapshot: line.CoverURL, SpecSnapshot: normalizeSpecSnapshot(line.SpecSnapshot), UnitPrice: money(line.UnitCents), Quantity: line.Quantity}
				if err := tx.Create(&item).Error; err != nil {
					return err
				}
				if err := merchantstock.EnqueueReserve(tx, order.ID, store.StoreID, line.MerchantSKUID, line.Quantity, reservationExpiry); err != nil {
					return err
				}
			}
			if input.DeliveryType == "pickup" || input.DeliveryType == "service" || input.DeliveryType == "city" {
				status := "awaiting"
				if input.DeliveryType == "city" {
					status = "pending_dispatch"
				}
				if err := tx.Exec(
					`INSERT INTO qixi_crm_b_order_delivery (order_id,delivery_type,status) VALUES (?,?,?)`,
					order.ID, input.DeliveryType, status,
				).Error; err != nil {
					return err
				}
			}
		}
		created = CreatedOrder{GroupOrderID: group.ID, GroupOrderNo: group.OrderNo, PayCents: checkout.TotalCents - pricing.DiscountCents - quote.DiscountCents, TotalQuantity: checkout.TotalQty}
		return nil
	})
	return created, err
}

type addressSnapshot struct {
	Recipient string `json:"recipient"`
	Mobile    string `json:"mobile"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Detail    string `json:"detail"`
	PostCode  int    `json:"post_code"`
}

func (addressSnapshot) TableName() string { return "qixi_crm_b_user_address" }
func loadAddress(ctx context.Context, db *gorm.DB, uid, id uint64) (addressSnapshot, error) {
	var address addressSnapshot
	err := db.WithContext(ctx).Where("id = ? AND user_id = ?", id, uid).First(&address).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return address, ErrAddressOwnership
	}
	return address, err
}

type groupRow struct {
	ID                uint64     `gorm:"column:id"`
	OrderNo           string     `gorm:"column:order_no"`
	UserID            uint64     `gorm:"column:user_id"`
	TotalAmount       float64    `gorm:"column:total_amount"`
	DiscountAmount    float64    `gorm:"column:discount_amount"`
	PayAmount         float64    `gorm:"column:pay_amount"`
	PayStatus         string     `gorm:"column:pay_status"`
	ActivityType      int        `gorm:"column:activity_type"`
	PointsAmount      int64      `gorm:"column:points_amount"`
	TotalQuantity     int        `gorm:"column:total_quantity"`
	RecipientSnapshot string     `gorm:"column:recipient_snapshot"`
	IdempotencyKey    string     `gorm:"column:idempotency_key"`
	Remark            string     `gorm:"column:remark"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
	ArchivedAt        *time.Time `gorm:"column:user_archived_at"`
}

func (groupRow) TableName() string { return "qixi_crm_b_group_order" }

type orderRow struct {
	ID                   uint64  `gorm:"column:id"`
	GroupOrderID         uint64  `gorm:"column:group_order_id"`
	OrderNo              string  `gorm:"column:order_no"`
	MerchantID           uint64  `gorm:"column:merchant_id"`
	MerchantNameSnapshot string  `gorm:"column:merchant_name_snapshot"`
	StoreID              uint64  `gorm:"column:store_id"`
	StoreNameSnapshot    string  `gorm:"column:store_name_snapshot"`
	UserID               uint64  `gorm:"column:user_id"`
	TotalAmount          float64 `gorm:"column:total_amount"`
	DiscountAmount       float64 `gorm:"column:discount_amount"`
	PayAmount            float64 `gorm:"column:pay_amount"`
	TotalQuantity        int     `gorm:"column:total_quantity"`
	RecipientSnapshot    string  `gorm:"column:recipient_snapshot"`
	Remark               string  `gorm:"column:remark"`
	Status               string  `gorm:"column:status"`
	ActivityType         int     `gorm:"column:activity_type"`
	PointsAmount         int64   `gorm:"column:points_amount"`
}

func (orderRow) TableName() string { return "qixi_crm_b_order" }

type orderItemRow struct {
	ID               uint64  `gorm:"column:id"`
	OrderID          uint64  `gorm:"column:order_id"`
	ProductID        uint64  `gorm:"column:product_id"`
	MerchantSKUID    uint64  `gorm:"column:merchant_sku_id"`
	SKUKey           string  `gorm:"column:sku_key"`
	TitleSnapshot    string  `gorm:"column:title_snapshot"`
	CoverURLSnapshot string  `gorm:"column:cover_url_snapshot"`
	SpecSnapshot     string  `gorm:"column:spec_snapshot"`
	UnitPrice        float64 `gorm:"column:unit_price"`
	Quantity         int     `gorm:"column:quantity"`
	Commented        bool    `gorm:"-"`
}

func (orderItemRow) TableName() string { return "qixi_crm_b_order_item" }
func money(cents int64) float64        { return float64(cents) / 100 }
func orderNo(prefix string) string {
	var nonce [6]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return fmt.Sprintf("%s%s", prefix, time.Now().UTC().Format("20060102150405.000000000"))
	}
	return fmt.Sprintf("%s%s%x", prefix, time.Now().UTC().Format("20060102150405"), nonce)
}
