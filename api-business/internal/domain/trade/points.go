package trade

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// PointsInput 积分商城核对/下单（独立于 v2，禁止混用券）
type PointsInput struct {
	ProductID         uint   `json:"product_id"`
	ProductAttrUnique string `json:"product_attr_unique"`
	CartNum           uint   `json:"cart_num"`
	AddressID         uint   `json:"address_id"`
}

type PointsCheckResult struct {
	ProductID    uint    `json:"product_id"`
	StoreName    string  `json:"store_name"`
	Image        string  `json:"image"`
	CartNum      uint    `json:"cart_num"`
	Integral     int     `json:"integral"` // 应付积分 = product.integral * cart_num
	UserIntegral int     `json:"user_integral"`
	PayPrice     float64 `json:"pay_price"` // 积分单现金为 0
	MerID        uint    `json:"mer_id"`
	MerName      string  `json:"mer_name"`
	ActivityType int8    `json:"activity_type"`
}

// V3Check 积分商城核对：入口 /order/v3/check（≠ /v3/order/*）
func (s *Service) V3Check(ctx context.Context, uid uint, in PointsInput) (*PointsCheckResult, error) {
	return s.buildPointsCheck(ctx, uid, in)
}

// V3Create 积分商城下单：入口 /order/v3/create（创建时不扣积分）
func (s *Service) V3Create(ctx context.Context, uid uint, in PointsInput) (*GroupOrder, error) {
	if in.AddressID == 0 {
		return nil, ErrAddressRequired
	}
	check, err := s.buildPointsCheck(ctx, uid, in)
	if err != nil {
		return nil, err
	}
	if check.Integral > check.UserIntegral {
		return nil, ErrIntegralNotEnough
	}
	addr, err := s.cartSvc.GetAddress(ctx, uid, in.AddressID)
	if err != nil {
		return nil, err
	}
	pv, err := s.store.LoadPointsProduct(ctx, in.ProductID, strings.TrimSpace(in.ProductAttrUnique))
	if err != nil {
		return nil, err
	}
	fullAddr := addr.FullAddress()
	num := in.CartNum
	if num == 0 {
		num = 1
	}

	var created *GroupOrder
	err = s.store.WithTx(func(tx Store) error {
		g := &GroupOrder{
			GroupOrderSN:  genSN("PG"),
			UID:           uid,
			TotalPostage:  0,
			TotalPrice:    0,
			TotalNum:      int(num),
			Integral:      check.Integral,
			IntegralPrice: 0,
			GiveIntegral:  0,
			RealName:      addr.RealName,
			UserPhone:     addr.Phone,
			UserAddress:   fullAddr,
			PayPrice:      0,
			PayPostage:    0,
			Cost:          pv.Cost * float64(num),
			Paid:          0,
			ActivityType:  ActivityTypePoints,
		}
		if err := tx.CreateGroupOrder(ctx, g); err != nil {
			return err
		}
		o := &StoreOrder{
			GroupOrderID: g.GroupOrderID,
			OrderSN:      genSN("PO"),
			UID:          uid,
			RealName:     addr.RealName,
			UserPhone:    addr.Phone,
			UserAddress:  fullAddr,
			CartID:       "",
			TotalNum:     int(num),
			TotalPrice:   0,
			PayPrice:     0,
			Paid:         0,
			Status:       OrderStatusAwaitShip,
			MerID:        pv.MerID,
			Cost:         pv.Cost * float64(num),
			Integral:     check.Integral,
			ActivityType: ActivityTypePoints,
			VerifyCode:   ensureVerifyCode(""),
		}
		if err := tx.CreateStoreOrder(ctx, o); err != nil {
			return err
		}
		info, _ := json.Marshal(map[string]interface{}{
			"product_id":   pv.ProductID,
			"store_name":   pv.StoreName,
			"image":        pv.Image,
			"price":        pv.Price,
			"integral":     pv.Integral,
			"unique":       pv.Unique,
			"product_type": pv.ProductType,
		})
		op := &OrderProduct{
			OrderID:      o.OrderID,
			UID:          uid,
			CartID:       0,
			ProductID:    pv.ProductID,
			ProductSKU:   trimSKU(pv.Unique),
			ProductNum:   int(num),
			ProductType:  int8(pv.ProductType),
			Cost:         pv.Cost,
			ProductPrice: pv.Price,
			TotalPrice:   0,
			ProductInfo:  string(info),
		}
		if err := tx.CreateOrderProduct(ctx, op); err != nil {
			return err
		}
		o.MerName = pv.MerName
		o.Products = []OrderProduct{*op}
		g.Orders = []StoreOrder{*o}
		created = g
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

// PointsPay 积分单支付：入口 /order/points/pay/:id（扣积分 + MarkPaid，幂等）
func (s *Service) PointsPay(ctx context.Context, uid, groupOrderID uint) (*GroupOrder, error) {
	g, err := s.store.GetGroupOrder(ctx, groupOrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if g.UID != uid || g.IsDel == 1 {
		return nil, ErrForbidden
	}
	if g.ActivityType != ActivityTypePoints {
		return nil, ErrNotPointsProduct
	}
	if g.Paid == 1 {
		return s.attachGroup(ctx, g)
	}

	children, err := s.store.ListStoreOrdersByGroup(ctx, groupOrderID)
	if err != nil {
		return nil, err
	}
	orderIDs := make([]uint, 0, len(children))
	for _, c := range children {
		orderIDs = append(orderIDs, c.OrderID)
	}
	products, err := s.store.ListOrderProductsByOrders(ctx, orderIDs)
	if err != nil {
		return nil, err
	}

	err = s.store.WithTx(func(tx Store) error {
		cur, err := tx.GetGroupOrder(ctx, groupOrderID)
		if err != nil {
			return err
		}
		if cur.Paid == 1 {
			return nil
		}
		if cur.Integral <= 0 {
			return ErrBadParam
		}
		link := fmt.Sprintf("points:%d", groupOrderID)
		dup, err := tx.HasBill(ctx, uid, "integral", "deduction", link)
		if err != nil {
			return err
		}
		if !dup {
			if err := tx.DeductUserIntegral(ctx, uid, cur.Integral); err != nil {
				return err
			}
			bal, err := tx.GetUserIntegral(ctx, uid)
			if err != nil {
				return err
			}
			if err := tx.CreateBill(ctx, &UserBill{
				UID:      uid,
				LinkID:   link,
				PM:       BillPMOut,
				Title:    "积分兑换",
				Category: "integral",
				Type:     "deduction",
				Number:   float64(cur.Integral),
				Balance:  float64(bal),
				Mark:     fmt.Sprintf("积分商城订单%d", groupOrderID),
				Status:   1,
			}); err != nil {
				return err
			}
		}
		now := time.Now()
		if err := tx.MarkGroupPaid(ctx, groupOrderID, PayTypeIntegral, now); err != nil {
			return err
		}
		if err := tx.MarkChildrenPaid(ctx, groupOrderID, PayTypeIntegral, now); err != nil {
			return err
		}
		for _, p := range products {
			n := uint(p.ProductNum)
			if n == 0 {
				continue
			}
			if err := tx.DeductProductStock(ctx, p.ProductID, n); err != nil {
				return err
			}
			if err := tx.DeductSKUStock(ctx, p.ProductID, skuUnique(p), n); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	g, err = s.store.GetGroupOrder(ctx, groupOrderID)
	if err != nil {
		return nil, err
	}
	return s.attachGroup(ctx, g)
}

func (s *Service) buildPointsCheck(ctx context.Context, uid uint, in PointsInput) (*PointsCheckResult, error) {
	if in.ProductID == 0 {
		return nil, ErrBadParam
	}
	num := in.CartNum
	if num == 0 {
		num = 1
	}
	pv, err := s.store.LoadPointsProduct(ctx, in.ProductID, strings.TrimSpace(in.ProductAttrUnique))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotPointsProduct
		}
		return nil, err
	}
	if pv.Stock < num {
		return nil, ErrStockNotEnough
	}
	if pv.Integral <= 0 {
		return nil, ErrBadParam
	}
	bal, err := s.store.GetUserIntegral(ctx, uid)
	if err != nil {
		return nil, err
	}
	need := pv.Integral * int(num)
	return &PointsCheckResult{
		ProductID:    pv.ProductID,
		StoreName:    pv.StoreName,
		Image:        pv.Image,
		CartNum:      num,
		Integral:     need,
		UserIntegral: bal,
		PayPrice:     0,
		MerID:        pv.MerID,
		MerName:      pv.MerName,
		ActivityType: ActivityTypePoints,
	}, nil
}
