package combination

import "context"

// TradeBridge 适配 trade.CombinationHook。
type TradeBridge struct{ svc *Service }

func NewTradeBridge(svc *Service) *TradeBridge { return &TradeBridge{svc: svc} }

func (b *TradeBridge) Quote(ctx context.Context, productGroupID uint) (price float64, productID, merID uint, storeName, image, merName string, err error) {
	g, err := b.svc.Quote(ctx, productGroupID)
	if err != nil {
		return 0, 0, 0, "", "", "", err
	}
	return g.Price, g.ProductID, g.MerID, g.StoreName, g.Image, g.MerName, nil
}

func (b *TradeBridge) ProductCost(ctx context.Context, productID uint) (float64, error) {
	return b.svc.ProductCost(ctx, productID)
}

func (b *TradeBridge) BeginJoin(ctx context.Context, uid, productGroupID, joinBuyingID uint, nickname string) (uint, bool, error) {
	return b.svc.BeginJoin(ctx, uid, productGroupID, joinBuyingID, nickname)
}

func (b *TradeBridge) AttachMember(ctx context.Context, buyingID, productGroupID, uid, orderID uint, isLeader bool, nickname string) error {
	return b.svc.AttachMember(ctx, buyingID, productGroupID, uid, orderID, isLeader, nickname)
}

func (b *TradeBridge) OnOrderPaid(ctx context.Context, orderID uint) (bool, []uint, error) {
	return b.svc.OnOrderPaid(ctx, orderID)
}

func (b *TradeBridge) CancelUnpaid(ctx context.Context, orderID uint) error {
	return b.svc.CancelUnpaid(ctx, orderID)
}
