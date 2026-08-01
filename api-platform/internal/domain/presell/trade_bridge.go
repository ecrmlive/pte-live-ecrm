package presell

import (
	"context"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/trade"
)

// TradeBridge 适配 trade.PresellHook。
type TradeBridge struct{ svc *Service }

func NewTradeBridge(svc *Service) *TradeBridge { return &TradeBridge{svc: svc} }

func (b *TradeBridge) Quote(ctx context.Context, productPresellID uint) (*trade.PresellQuote, error) {
	p, err := b.svc.Quote(ctx, productPresellID)
	if err != nil {
		return nil, err
	}
	fs, fe := parseFinalWindow(p.FinalStartTime, p.FinalEndTime)
	return &trade.PresellQuote{
		ProductPresellID: p.ProductPresellID,
		Price:            p.Price,
		DownPrice:        p.DownPrice,
		FinalPrice:       p.FinalPrice,
		ProductID:        p.ProductID,
		MerID:            p.MerID,
		StoreName:        p.StoreName,
		Image:            p.Image,
		MerName:          p.MerName,
		Stock:            p.Stock,
		PresellType:      p.PresellType,
		FinalStart:       fs,
		FinalEnd:         fe,
	}, nil
}

func (b *TradeBridge) ProductCost(ctx context.Context, productID uint) (float64, error) {
	return b.svc.ProductCost(ctx, productID)
}

func (b *TradeBridge) ReserveStock(ctx context.Context, productPresellID uint, num int) error {
	return b.svc.ReserveStock(ctx, productPresellID, num)
}

func (b *TradeBridge) RestoreStock(ctx context.Context, productPresellID uint, num int) error {
	return b.svc.RestoreStock(ctx, productPresellID, num)
}

func (b *TradeBridge) OnOrderPaid(ctx context.Context, productPresellID uint, num int) error {
	return b.svc.OnOrderPaid(ctx, productPresellID, num)
}

var _ trade.PresellHook = (*TradeBridge)(nil)
