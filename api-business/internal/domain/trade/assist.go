package trade

import (
	"context"
	"encoding/json"
	"fmt"
)

// ActivityTypeAssist 助力（product_type=3 / activity_type=3）
const ActivityTypeAssist int8 = 3

// AssistHook 助力钩子。
type AssistHook interface {
	Quote(ctx context.Context, setID, uid uint) (price float64, productID, merID, assistID uint, storeName, image, merName string, stock int, err error)
	ProductCost(ctx context.Context, productID uint) (float64, error)
	ReserveStock(ctx context.Context, assistID uint, num int) error
	RestoreStock(ctx context.Context, assistID uint, num int) error
	MarkSetPaid(ctx context.Context, setID uint) error
}

func (s *Service) SetAssist(h AssistHook) { s.assist = h }

type AssistInput struct {
	ProductAssistSetID uint `json:"product_assist_set_id"`
	CartNum            uint `json:"cart_num"`
	AddressID          uint `json:"address_id"`
}

type AssistCheckResult struct {
	ProductAssistSetID uint    `json:"product_assist_set_id"`
	ProductAssistID    uint    `json:"product_assist_id"`
	ProductID          uint    `json:"product_id"`
	StoreName          string  `json:"store_name"`
	Image              string  `json:"image"`
	MerID              uint    `json:"mer_id"`
	MerName            string  `json:"mer_name"`
	CartNum            uint    `json:"cart_num"`
	Price              float64 `json:"price"`
	PayPrice           float64 `json:"pay_price"`
	Stock              int     `json:"stock"`
	ActivityType       int8    `json:"activity_type"`
}

func (s *Service) AssistCheck(ctx context.Context, uid uint, in AssistInput) (*AssistCheckResult, error) {
	if s.assist == nil || in.ProductAssistSetID == 0 {
		return nil, ErrBadParam
	}
	num := in.CartNum
	if num == 0 {
		num = 1
	}
	if num != 1 {
		return nil, ErrBadParam // 助力单次购买 1 件
	}
	price, productID, merID, assistID, storeName, image, merName, stock, err := s.assist.Quote(ctx, in.ProductAssistSetID, uid)
	if err != nil {
		return nil, err
	}
	if stock < int(num) {
		return nil, ErrStockNotEnough
	}
	return &AssistCheckResult{
		ProductAssistSetID: in.ProductAssistSetID, ProductAssistID: assistID,
		ProductID: productID, StoreName: storeName, Image: image,
		MerID: merID, MerName: merName, CartNum: num,
		Price: price, PayPrice: round2(price * float64(num)), Stock: stock,
		ActivityType: ActivityTypeAssist,
	}, nil
}

func (s *Service) AssistCreate(ctx context.Context, uid uint, in AssistInput) (*GroupOrder, error) {
	if s.assist == nil {
		return nil, ErrBadParam
	}
	if in.AddressID == 0 {
		return nil, ErrAddressRequired
	}
	check, err := s.AssistCheck(ctx, uid, in)
	if err != nil {
		return nil, err
	}
	addr, err := s.cartSvc.GetAddress(ctx, uid, in.AddressID)
	if err != nil {
		return nil, err
	}
	cost, err := s.assist.ProductCost(ctx, check.ProductID)
	if err != nil {
		return nil, err
	}
	fullAddr := addr.FullAddress()
	num := check.CartNum

	var created *GroupOrder
	err = s.store.WithTx(func(tx Store) error {
		if err := s.assist.ReserveStock(ctx, check.ProductAssistID, int(num)); err != nil {
			return err
		}
		g := &GroupOrder{
			GroupOrderSN: genSN("AS"),
			UID:          uid,
			TotalPostage: 0,
			TotalPrice:   check.PayPrice,
			TotalNum:     int(num),
			RealName:     addr.RealName,
			UserPhone:    addr.Phone,
			UserAddress:  fullAddr,
			PayPrice:     check.PayPrice,
			Cost:         cost * float64(num),
			Paid:         0,
			ActivityType: ActivityTypeAssist,
		}
		if err := tx.CreateGroupOrder(ctx, g); err != nil {
			return err
		}
		o := &StoreOrder{
			GroupOrderID: g.GroupOrderID,
			OrderSN:      genSN("AO"),
			UID:          uid,
			RealName:     addr.RealName,
			UserPhone:    addr.Phone,
			UserAddress:  fullAddr,
			TotalNum:     int(num),
			TotalPrice:   check.PayPrice,
			PayPrice:     check.PayPrice,
			Paid:         0,
			Status:       OrderStatusAwaitShip,
			MerID:        check.MerID,
			Cost:         cost * float64(num),
			ActivityType: ActivityTypeAssist,
			Mark:         fmt.Sprintf("助力=%d/set=%d", check.ProductAssistID, in.ProductAssistSetID),
		}
		if err := tx.CreateStoreOrder(ctx, o); err != nil {
			return err
		}
		info, _ := json.Marshal(map[string]any{
			"store_name": check.StoreName, "image": check.Image,
			"price": check.Price, "product_assist_id": check.ProductAssistID,
			"product_assist_set_id": in.ProductAssistSetID,
		})
		p := &OrderProduct{
			OrderID: o.OrderID, UID: uid, ProductID: check.ProductID,
			ProductNum: int(num), ProductType: ActivityTypeAssist,
			ActivityID: in.ProductAssistSetID, Cost: cost,
			ProductPrice: check.Price, TotalPrice: check.PayPrice,
			ProductInfo: string(info),
		}
		if err := tx.CreateOrderProduct(ctx, p); err != nil {
			return err
		}
		created = g
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.attachGroup(ctx, created)
}
