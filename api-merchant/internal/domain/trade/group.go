package trade

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ActivityTypeCombination 拼团（product_type=4）
const ActivityTypeCombination int8 = 4

// OrderStatusGrouping 拼团中（成团后转为待发货 0）
const OrderStatusGrouping int8 = 9

// CombinationHook 拼团活动钩子。
type CombinationHook interface {
	Quote(ctx context.Context, productGroupID uint) (price float64, productID, merID uint, storeName, image, merName string, err error)
	ProductCost(ctx context.Context, productID uint) (float64, error)
	BeginJoin(ctx context.Context, uid, productGroupID, joinBuyingID uint, nickname string) (buyingID uint, isLeader bool, err error)
	AttachMember(ctx context.Context, buyingID, productGroupID, uid, orderID uint, isLeader bool, nickname string) error
	OnOrderPaid(ctx context.Context, orderID uint) (success bool, orderIDs []uint, err error)
	CancelUnpaid(ctx context.Context, orderID uint) error
}

func (s *Service) SetCombination(h CombinationHook) { s.combo = h }

type GroupInput struct {
	ProductGroupID uint `json:"product_group_id"`
	GroupBuyingID  uint `json:"group_buying_id"` // 0=开团
	CartNum        uint `json:"cart_num"`
	AddressID      uint `json:"address_id"`
}

type GroupCheckResult struct {
	ProductGroupID uint    `json:"product_group_id"`
	GroupBuyingID  uint    `json:"group_buying_id"`
	ProductID      uint    `json:"product_id"`
	StoreName      string  `json:"store_name"`
	Image          string  `json:"image"`
	MerID          uint    `json:"mer_id"`
	MerName        string  `json:"mer_name"`
	CartNum        uint    `json:"cart_num"`
	Price          float64 `json:"price"`
	PayPrice       float64 `json:"pay_price"`
	ActivityType   int8    `json:"activity_type"`
}

func (s *Service) GroupCheck(ctx context.Context, uid uint, in GroupInput) (*GroupCheckResult, error) {
	_ = uid
	if s.combo == nil || in.ProductGroupID == 0 {
		return nil, ErrBadParam
	}
	num := in.CartNum
	if num == 0 {
		num = 1
	}
	price, productID, merID, storeName, image, merName, err := s.combo.Quote(ctx, in.ProductGroupID)
	if err != nil {
		return nil, err
	}
	pay := round2(price * float64(num))
	return &GroupCheckResult{
		ProductGroupID: in.ProductGroupID, GroupBuyingID: in.GroupBuyingID,
		ProductID: productID, StoreName: storeName, Image: image,
		MerID: merID, MerName: merName, CartNum: num,
		Price: price, PayPrice: pay, ActivityType: ActivityTypeCombination,
	}, nil
}

func (s *Service) GroupCreate(ctx context.Context, uid uint, in GroupInput) (*GroupOrder, error) {
	if s.combo == nil {
		return nil, ErrBadParam
	}
	if in.AddressID == 0 {
		return nil, ErrAddressRequired
	}
	check, err := s.GroupCheck(ctx, uid, in)
	if err != nil {
		return nil, err
	}
	addr, err := s.cartSvc.GetAddress(ctx, uid, in.AddressID)
	if err != nil {
		return nil, err
	}
	cost, err := s.combo.ProductCost(ctx, check.ProductID)
	if err != nil {
		return nil, err
	}
	nickname := strings.TrimSpace(addr.RealName)
	buyingID, isLeader, err := s.combo.BeginJoin(ctx, uid, in.ProductGroupID, in.GroupBuyingID, nickname)
	if err != nil {
		return nil, err
	}
	fullAddr := addr.FullAddress()
	num := check.CartNum

	var created *GroupOrder
	err = s.store.WithTx(func(tx Store) error {
		g := &GroupOrder{
			GroupOrderSN: genSN("CG"),
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
			ActivityType: ActivityTypeCombination,
		}
		if err := tx.CreateGroupOrder(ctx, g); err != nil {
			return err
		}
		o := &StoreOrder{
			GroupOrderID: g.GroupOrderID,
			OrderSN:      genSN("CO"),
			UID:          uid,
			RealName:     addr.RealName,
			UserPhone:    addr.Phone,
			UserAddress:  fullAddr,
			TotalNum:     int(num),
			TotalPrice:   check.PayPrice,
			PayPrice:     check.PayPrice,
			Paid:         0,
			Status:       OrderStatusGrouping,
			MerID:        check.MerID,
			Cost:         cost * float64(num),
			ActivityType: ActivityTypeCombination,
			Mark:         fmt.Sprintf("拼团buying=%d", buyingID),
		}
		if err := tx.CreateStoreOrder(ctx, o); err != nil {
			return err
		}
		info, _ := json.Marshal(map[string]any{
			"store_name": check.StoreName, "image": check.Image,
			"price": check.Price, "product_group_id": check.ProductGroupID,
			"group_buying_id": buyingID,
		})
		p := &OrderProduct{
			OrderID: o.OrderID, UID: uid, ProductID: check.ProductID,
			ProductNum: int(num), ProductType: ActivityTypeCombination,
			ActivityID: check.ProductGroupID, Cost: cost,
			ProductPrice: check.Price, TotalPrice: check.PayPrice,
			ProductInfo: string(info),
		}
		if err := tx.CreateOrderProduct(ctx, p); err != nil {
			return err
		}
		if err := s.combo.AttachMember(ctx, buyingID, check.ProductGroupID, uid, o.OrderID, isLeader, nickname); err != nil {
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
