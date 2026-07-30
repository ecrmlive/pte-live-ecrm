package trade

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ActivityTypePresell 预售（product_type=2 / activity_type=2）
const ActivityTypePresell int8 = 2

// OrderStatusAwaitFinal 定金已付待付尾款
const OrderStatusAwaitFinal int8 = 10

// OrderStatusFinalTimeout 尾款超时未付
const OrderStatusFinalTimeout int8 = 11

// PresellQuote 预售报价（全款/定金）。
type PresellQuote struct {
	ProductPresellID uint
	Price            float64
	DownPrice        float64
	FinalPrice       float64
	ProductID        uint
	MerID            uint
	StoreName        string
	Image            string
	MerName          string
	Stock            int
	PresellType      int
	FinalStart       time.Time
	FinalEnd         time.Time
}

// PresellHook 预售钩子。
type PresellHook interface {
	Quote(ctx context.Context, productPresellID uint) (*PresellQuote, error)
	ProductCost(ctx context.Context, productID uint) (float64, error)
	ReserveStock(ctx context.Context, productPresellID uint, num int) error
	RestoreStock(ctx context.Context, productPresellID uint, num int) error
	OnOrderPaid(ctx context.Context, productPresellID uint, num int) error
}

func (s *Service) SetPresell(h PresellHook) { s.presell = h }

// PresellOrder 定金预售尾款单（qixi_m_app_presell_order）。
type PresellOrder struct {
	PresellOrderID   uint       `gorm:"column:presell_order_id;primaryKey" json:"presell_order_id"`
	PresellOrderSN   string     `gorm:"column:presell_order_sn" json:"presell_order_sn"`
	UID              uint       `gorm:"column:uid" json:"uid"`
	MerID            uint       `gorm:"column:mer_id" json:"mer_id"`
	OrderID          uint       `gorm:"column:order_id" json:"order_id"`
	ProductPresellID uint       `gorm:"column:product_presell_id" json:"product_presell_id"`
	FinalStartTime   time.Time  `gorm:"column:final_start_time" json:"final_start_time"`
	FinalEndTime     time.Time  `gorm:"column:final_end_time" json:"final_end_time"`
	Paid             int8       `gorm:"column:paid" json:"paid"`
	Status           int8       `gorm:"column:status" json:"status"`
	PayType          int8       `gorm:"column:pay_type" json:"pay_type"`
	PayPrice         float64    `gorm:"column:pay_price" json:"pay_price"`
	PayTime          *time.Time `gorm:"column:pay_time" json:"pay_time,omitempty"`
	CreateTime       time.Time  `gorm:"column:create_time" json:"create_time"`

	StoreName string `gorm:"-" json:"store_name,omitempty"`
	OrderSN   string `gorm:"-" json:"order_sn,omitempty"`
}

func (PresellOrder) TableName() string { return "qixi_m_app_presell_order" }

type PresellInput struct {
	ProductPresellID uint `json:"product_presell_id"`
	CartNum          uint `json:"cart_num"`
	AddressID        uint `json:"address_id"`
}

type PresellCheckResult struct {
	ProductPresellID uint      `json:"product_presell_id"`
	ProductID        uint      `json:"product_id"`
	StoreName        string    `json:"store_name"`
	Image            string    `json:"image"`
	MerID            uint      `json:"mer_id"`
	MerName          string    `json:"mer_name"`
	CartNum          uint      `json:"cart_num"`
	Price            float64   `json:"price"`
	DownPrice        float64   `json:"down_price"`
	FinalPrice       float64   `json:"final_price"`
	PayPrice         float64   `json:"pay_price"`
	TotalPrice       float64   `json:"total_price"`
	Stock            int       `json:"stock"`
	ActivityType     int8      `json:"activity_type"`
	PresellType      int       `json:"presell_type"`
	FinalStartTime   time.Time `json:"final_start_time"`
	FinalEndTime     time.Time `json:"final_end_time"`
}

type PresellFinalPage struct {
	List  []PresellOrder `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

func (s *Service) PresellCheck(ctx context.Context, uid uint, in PresellInput) (*PresellCheckResult, error) {
	_ = uid
	if s.presell == nil || in.ProductPresellID == 0 {
		return nil, ErrBadParam
	}
	num := in.CartNum
	if num == 0 {
		num = 1
	}
	q, err := s.presell.Quote(ctx, in.ProductPresellID)
	if err != nil {
		return nil, err
	}
	if q.Stock < int(num) {
		return nil, ErrStockNotEnough
	}
	full := round2(q.Price * float64(num))
	pay := full
	if q.PresellType == 2 {
		pay = round2(q.DownPrice * float64(num))
	}
	return &PresellCheckResult{
		ProductPresellID: in.ProductPresellID,
		ProductID:        q.ProductID, StoreName: q.StoreName, Image: q.Image,
		MerID: q.MerID, MerName: q.MerName, CartNum: num,
		Price: q.Price, DownPrice: q.DownPrice, FinalPrice: q.FinalPrice,
		PayPrice: pay, TotalPrice: full, Stock: q.Stock,
		ActivityType: ActivityTypePresell, PresellType: q.PresellType,
		FinalStartTime: q.FinalStart, FinalEndTime: q.FinalEnd,
	}, nil
}

func (s *Service) PresellCreate(ctx context.Context, uid uint, in PresellInput) (*GroupOrder, error) {
	if s.presell == nil {
		return nil, ErrBadParam
	}
	if in.AddressID == 0 {
		return nil, ErrAddressRequired
	}
	check, err := s.PresellCheck(ctx, uid, in)
	if err != nil {
		return nil, err
	}
	addr, err := s.cartSvc.GetAddress(ctx, uid, in.AddressID)
	if err != nil {
		return nil, err
	}
	cost, err := s.presell.ProductCost(ctx, check.ProductID)
	if err != nil {
		return nil, err
	}
	q, err := s.presell.Quote(ctx, in.ProductPresellID)
	if err != nil {
		return nil, err
	}
	fullAddr := addr.FullAddress()
	num := check.CartNum

	var created *GroupOrder
	err = s.store.WithTx(func(tx Store) error {
		if err := s.presell.ReserveStock(ctx, in.ProductPresellID, int(num)); err != nil {
			return err
		}
		g := &GroupOrder{
			GroupOrderSN: genSN("PS"),
			UID:          uid,
			TotalPostage: 0,
			TotalPrice:   check.TotalPrice,
			TotalNum:     int(num),
			RealName:     addr.RealName,
			UserPhone:    addr.Phone,
			UserAddress:  fullAddr,
			PayPrice:     check.PayPrice,
			Cost:         cost * float64(num),
			Paid:         0,
			ActivityType: ActivityTypePresell,
		}
		if err := tx.CreateGroupOrder(ctx, g); err != nil {
			return err
		}
		mark := fmt.Sprintf("全款预售=%d", in.ProductPresellID)
		if check.PresellType == 2 {
			mark = fmt.Sprintf("定金预售=%d", in.ProductPresellID)
		}
		o := &StoreOrder{
			GroupOrderID: g.GroupOrderID,
			OrderSN:      genSN("PO"),
			UID:          uid,
			RealName:     addr.RealName,
			UserPhone:    addr.Phone,
			UserAddress:  fullAddr,
			TotalNum:     int(num),
			TotalPrice:   check.TotalPrice,
			PayPrice:     check.PayPrice,
			Paid:         0,
			Status:       OrderStatusAwaitShip,
			MerID:        check.MerID,
			Cost:         cost * float64(num),
			ActivityType: ActivityTypePresell,
			Mark:         mark,
		}
		if err := tx.CreateStoreOrder(ctx, o); err != nil {
			return err
		}
		info, _ := json.Marshal(map[string]any{
			"store_name": check.StoreName, "image": check.Image,
			"price": check.Price, "product_presell_id": in.ProductPresellID,
			"presell_type": check.PresellType,
			"down_price":   q.DownPrice, "final_price": q.FinalPrice,
		})
		p := &OrderProduct{
			OrderID: o.OrderID, UID: uid, ProductID: check.ProductID,
			ProductNum: int(num), ProductType: ActivityTypePresell,
			ActivityID: in.ProductPresellID, Cost: cost,
			ProductPrice: check.Price, TotalPrice: check.TotalPrice,
			ProductInfo: string(info),
		}
		if err := tx.CreateOrderProduct(ctx, p); err != nil {
			return err
		}
		if check.PresellType == 2 {
			unitFinal := q.FinalPrice
			if unitFinal <= 0 {
				unitFinal = round2(q.Price - q.DownPrice)
			}
			po := &PresellOrder{
				PresellOrderSN:   genSN("PF"),
				UID:              uid,
				MerID:            check.MerID,
				OrderID:          o.OrderID,
				ProductPresellID: in.ProductPresellID,
				FinalStartTime:   q.FinalStart,
				FinalEndTime:     q.FinalEnd,
				Paid:             0,
				Status:           1,
				PayPrice:         round2(unitFinal * float64(num)),
				CreateTime:       time.Now(),
			}
			if err := tx.CreatePresellOrder(ctx, po); err != nil {
				return err
			}
		}
		created = g
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.attachGroup(ctx, created)
}

// GetPresellFinal 尾款单详情（过期则置 11）。
func (s *Service) GetPresellFinal(ctx context.Context, uid, id uint) (*PresellOrder, error) {
	if uid == 0 || id == 0 {
		return nil, ErrBadParam
	}
	po, err := s.store.GetPresellOrder(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPresellFinalNotFound
		}
		return nil, err
	}
	if po.UID != uid {
		return nil, ErrForbidden
	}
	if po.Paid == 0 && po.Status == 1 && time.Now().After(po.FinalEndTime) {
		_ = s.expirePresellFinal(ctx, po)
		return nil, ErrPresellFinalTimeout
	}
	if o, err := s.store.GetStoreOrder(ctx, po.OrderID); err == nil && o != nil {
		po.OrderSN = o.OrderSN
	}
	if products, err := s.store.ListOrderProductsByOrder(ctx, po.OrderID); err == nil && len(products) > 0 {
		var meta map[string]any
		if json.Unmarshal([]byte(products[0].ProductInfo), &meta) == nil {
			if name, ok := meta["store_name"].(string); ok {
				po.StoreName = name
			}
		}
	}
	return po, nil
}

// ListPresellFinals 待付尾款列表；过期自动置 status=11 并还活动库存。
func (s *Service) ListPresellFinals(ctx context.Context, uid uint, page, limit int) (*PresellFinalPage, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	list, total, err := s.store.ListPresellOrdersByUID(ctx, uid, true, page, limit)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	out := make([]PresellOrder, 0, len(list))
	for i := range list {
		po := list[i]
		if po.Paid == 0 && po.Status == 1 && now.After(po.FinalEndTime) {
			_ = s.expirePresellFinal(ctx, &po)
			continue
		}
		if o, err := s.store.GetStoreOrder(ctx, po.OrderID); err == nil && o != nil {
			po.OrderSN = o.OrderSN
		}
		if products, err := s.store.ListOrderProductsByOrder(ctx, po.OrderID); err == nil && len(products) > 0 {
			var meta map[string]any
			if json.Unmarshal([]byte(products[0].ProductInfo), &meta) == nil {
				if name, ok := meta["store_name"].(string); ok {
					po.StoreName = name
				}
			}
		}
		out = append(out, po)
	}
	return &PresellFinalPage{List: out, Total: total, Page: page, Limit: limit}, nil
}

// PresellFinalPay 支付尾款；过期则 status→11、作废尾款单、还活动库存。
func (s *Service) PresellFinalPay(ctx context.Context, uid, presellOrderID uint, payTypeStr string) (*PresellOrder, error) {
	if uid == 0 || presellOrderID == 0 {
		return nil, ErrBadParam
	}
	payType, err := parsePayType(payTypeStr)
	if err != nil {
		return nil, err
	}
	if payType == PayTypeIntegral {
		return nil, ErrInvalidPayType
	}
	// 预售尾款当前只支持余额扣款；mock 仅允许在显式 sandbox 中做本地闭环。
	// 微信/支付宝不能在没有官方预下单与回调验签的情况下直接改为已支付。
	if payType == PayTypeWechat || payType == PayTypeAlipay || (payType == PayTypeMock && !s.payment.Sandbox) {
		return nil, ErrPaymentConfig
	}
	po, err := s.store.GetPresellOrder(ctx, presellOrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPresellFinalNotFound
		}
		return nil, err
	}
	if po.UID != uid {
		return nil, ErrForbidden
	}
	if po.Paid == 1 {
		return nil, ErrPresellFinalPaid
	}
	if po.Status != 1 {
		return nil, ErrPresellFinalInvalid
	}
	now := time.Now()
	if now.After(po.FinalEndTime) {
		_ = s.expirePresellFinal(ctx, po)
		return nil, ErrPresellFinalTimeout
	}
	if now.Before(po.FinalStartTime) {
		return nil, ErrPresellFinalNotOpen
	}
	o, err := s.store.GetStoreOrder(ctx, po.OrderID)
	if err != nil {
		return nil, err
	}
	if o.UID != uid || o.Paid != 1 || o.Status != OrderStatusAwaitFinal {
		return nil, ErrBadStatus
	}

	err = s.store.WithTx(func(tx Store) error {
		cur, err := tx.GetPresellOrder(ctx, presellOrderID)
		if err != nil {
			return err
		}
		if cur.Paid == 1 {
			return nil
		}
		if cur.Status != 1 {
			return ErrPresellFinalInvalid
		}
		if time.Now().After(cur.FinalEndTime) {
			return ErrPresellFinalTimeout
		}
		if payType == PayTypeBalance {
			bal, err := tx.GetUserBalance(ctx, uid)
			if err != nil {
				return err
			}
			if bal < cur.PayPrice {
				return ErrBalanceNotEnough
			}
			if err := tx.DeductUserBalance(ctx, uid, cur.PayPrice); err != nil {
				return err
			}
		}
		at := time.Now()
		if err := tx.MarkPresellOrderPaid(ctx, cur.PresellOrderID, payType, at); err != nil {
			return err
		}
		return tx.UpdateOrderStatus(ctx, cur.OrderID, OrderStatusAwaitShip)
	})
	if err != nil {
		if errors.Is(err, ErrPresellFinalTimeout) {
			_ = s.expirePresellFinal(ctx, po)
		}
		return nil, err
	}
	return s.store.GetPresellOrder(ctx, presellOrderID)
}

func (s *Service) expirePresellFinal(ctx context.Context, po *PresellOrder) error {
	if po == nil || po.Paid == 1 {
		return nil
	}
	return s.store.WithTx(func(tx Store) error {
		cur, err := tx.GetPresellOrder(ctx, po.PresellOrderID)
		if err != nil {
			return err
		}
		if cur.Paid == 1 || cur.Status != 1 {
			return nil
		}
		if err := tx.InvalidatePresellOrder(ctx, cur.PresellOrderID); err != nil {
			return err
		}
		if err := tx.UpdateOrderStatus(ctx, cur.OrderID, OrderStatusFinalTimeout); err != nil {
			return err
		}
		if s.presell != nil {
			o, err := tx.GetStoreOrder(ctx, cur.OrderID)
			if err != nil {
				return err
			}
			n := o.TotalNum
			if n <= 0 {
				n = 1
			}
			if err := s.presell.RestoreStock(ctx, cur.ProductPresellID, n); err != nil {
				return err
			}
		}
		return nil
	})
}
