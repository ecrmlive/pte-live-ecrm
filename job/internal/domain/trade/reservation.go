package trade

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ActivityTypeReservation 预约服务订单标记（订单 activity_type；商品 type=4）
const ActivityTypeReservation int8 = 30

// ReservationHook 预约校验与计数。
type ReservationHook interface {
	ValidateBook(ctx context.Context, productID, slotID uint, date string) (
		price float64, merID uint, storeName, image, merName, timePart string, cost float64, err error,
	)
	AfterBooked(ctx context.Context, slotID uint) error
}

func (s *Service) SetReservation(h ReservationHook) { s.reserve = h }

type ReservationInput struct {
	ProductID uint   `json:"product_id"`
	SlotID    uint   `json:"slot_id"`
	Date      string `json:"date"`
	AddressID uint   `json:"address_id"`
	Mark      string `json:"mark"`
}

type ReservationCheckResult struct {
	ProductID  uint    `json:"product_id"`
	SlotID     uint    `json:"slot_id"`
	Date       string  `json:"date"`
	TimePart   string  `json:"time_part"`
	StoreName  string  `json:"store_name"`
	Image      string  `json:"image"`
	MerID      uint    `json:"mer_id"`
	MerName    string  `json:"mer_name"`
	PayPrice   float64 `json:"pay_price"`
	VerifyHint string  `json:"verify_hint"`
}

func (s *Service) ReservationCheck(ctx context.Context, uid uint, in ReservationInput) (*ReservationCheckResult, error) {
	_ = uid
	if s.reserve == nil || in.ProductID == 0 || in.SlotID == 0 {
		return nil, ErrBadParam
	}
	price, merID, storeName, image, merName, timePart, _, err := s.reserve.ValidateBook(ctx, in.ProductID, in.SlotID, in.Date)
	if err != nil {
		return nil, err
	}
	return &ReservationCheckResult{
		ProductID: in.ProductID, SlotID: in.SlotID, Date: in.Date, TimePart: timePart,
		StoreName: storeName, Image: image, MerID: merID, MerName: merName,
		PayPrice: price, VerifyHint: "支付后凭核销码到店核销",
	}, nil
}

func (s *Service) ReservationCreate(ctx context.Context, uid uint, in ReservationInput) (*GroupOrder, error) {
	if s.reserve == nil {
		return nil, ErrBadParam
	}
	if in.AddressID == 0 {
		return nil, ErrAddressRequired
	}
	check, err := s.ReservationCheck(ctx, uid, in)
	if err != nil {
		return nil, err
	}
	price, merID, storeName, image, _, timePart, cost, err := s.reserve.ValidateBook(ctx, in.ProductID, in.SlotID, in.Date)
	if err != nil {
		return nil, err
	}
	_ = merID
	addr, err := s.cartSvc.GetAddress(ctx, uid, in.AddressID)
	if err != nil {
		return nil, err
	}
	fullAddr := addr.FullAddress()
	mark := strings.TrimSpace(in.Mark)
	if mark == "" {
		mark = fmt.Sprintf("预约 %s %s", in.Date, timePart)
	}

	var created *GroupOrder
	err = s.store.WithTx(func(tx Store) error {
		g := &GroupOrder{
			GroupOrderSN: genSN("RG"),
			UID:          uid,
			TotalPrice:   price,
			TotalNum:     1,
			RealName:     addr.RealName,
			UserPhone:    addr.Phone,
			UserAddress:  fullAddr,
			PayPrice:     price,
			Cost:         cost,
			Paid:         0,
			ActivityType: ActivityTypeReservation,
		}
		if err := tx.CreateGroupOrder(ctx, g); err != nil {
			return err
		}
		o := &StoreOrder{
			GroupOrderID:        g.GroupOrderID,
			OrderSN:             genSN("RO"),
			UID:                 uid,
			RealName:            addr.RealName,
			UserPhone:           addr.Phone,
			UserAddress:         fullAddr,
			TotalNum:            1,
			TotalPrice:          price,
			PayPrice:            price,
			Paid:                0,
			Status:              OrderStatusAwaitShip,
			MerID:               check.MerID,
			Cost:                cost,
			ActivityType:        ActivityTypeReservation,
			VerifyCode:          genVerifyCode(),
			ReservationDate:     in.Date,
			ReservationID:       in.SlotID,
			ReservationTimePart: timePart,
			Mark:                mark,
		}
		if err := tx.CreateStoreOrder(ctx, o); err != nil {
			return err
		}
		info, _ := json.Marshal(map[string]any{
			"store_name": storeName, "image": image, "price": price,
			"reservation_date": in.Date, "time_part": timePart, "slot_id": in.SlotID,
		})
		p := &OrderProduct{
			OrderID: o.OrderID, UID: uid, ProductID: in.ProductID,
			ProductNum: 1, ProductType: ActivityTypeReservation,
			ActivityID: in.SlotID, Cost: cost,
			ProductPrice: price, TotalPrice: price, ProductInfo: string(info),
		}
		if err := tx.CreateOrderProduct(ctx, p); err != nil {
			return err
		}
		if err := s.reserve.AfterBooked(ctx, in.SlotID); err != nil {
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

