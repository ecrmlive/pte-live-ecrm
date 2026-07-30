package assist

import "context"

// TradeBridge 适配 trade.AssistHook。
type TradeBridge struct{ svc *Service }

func NewTradeBridge(svc *Service) *TradeBridge { return &TradeBridge{svc: svc} }

func (b *TradeBridge) Quote(ctx context.Context, setID, uid uint) (price float64, productID, merID, assistID uint, storeName, image, merName string, stock int, err error) {
	return b.svc.QuoteForOrder(ctx, setID, uid)
}

func (b *TradeBridge) ProductCost(ctx context.Context, productID uint) (float64, error) {
	return b.svc.ProductCost(ctx, productID)
}

func (b *TradeBridge) ReserveStock(ctx context.Context, assistID uint, num int) error {
	return b.svc.ReserveStock(ctx, assistID, num)
}

func (b *TradeBridge) RestoreStock(ctx context.Context, assistID uint, num int) error {
	return b.svc.RestoreStock(ctx, assistID, num)
}

func (b *TradeBridge) MarkSetPaid(ctx context.Context, setID uint) error {
	return b.svc.MarkSetPaid(ctx, setID)
}
