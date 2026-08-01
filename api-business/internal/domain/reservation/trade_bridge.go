package reservation

import (
	"context"
	"fmt"
)

// TradeBridge 适配 trade.ReservationHook。
type TradeBridge struct{ svc *Service }

func NewTradeBridge(svc *Service) *TradeBridge { return &TradeBridge{svc: svc} }

func (b *TradeBridge) ValidateBook(ctx context.Context, productID, slotID uint, date string) (
	price float64, merID uint, storeName, image, merName, timePart string, cost float64, err error,
) {
	p, slot, err := b.svc.ValidateBook(ctx, productID, slotID, date)
	if err != nil {
		return 0, 0, "", "", "", "", 0, err
	}
	cost = p.Price // 演示：成本近似用售价；可后续接 SKU cost
	if cost <= 0 {
		cost = 1
	}
	return p.Price, p.MerID, p.StoreName, p.Image, p.MerName,
		fmt.Sprintf("%s-%s", slot.StartTime, slot.EndTime), cost * 0.3, nil
}

func (b *TradeBridge) AfterBooked(ctx context.Context, productID, slotID uint, date string, orderID, uid uint) error {
	return b.svc.AfterBooked(ctx, productID, slotID, date, orderID, uid)
}
