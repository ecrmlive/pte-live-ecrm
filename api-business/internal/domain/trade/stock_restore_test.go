package trade

import (
	"context"
	"testing"
)

type stubPresellHook struct {
	restoreID  uint
	restoreNum int
}

func (h *stubPresellHook) Quote(ctx context.Context, productPresellID uint) (*PresellQuote, error) {
	return &PresellQuote{ProductPresellID: productPresellID}, nil
}
func (h *stubPresellHook) ProductCost(ctx context.Context, productID uint) (float64, error) {
	return 0, nil
}
func (h *stubPresellHook) ReserveStock(ctx context.Context, productPresellID uint, num int) error {
	return nil
}
func (h *stubPresellHook) RestoreStock(ctx context.Context, productPresellID uint, num int) error {
	h.restoreID, h.restoreNum = productPresellID, num
	return nil
}
func (h *stubPresellHook) OnOrderPaid(ctx context.Context, productPresellID uint, num int) error {
	return nil
}

type stubAssistHook struct {
	restoreID  uint
	restoreNum int
}

func (h *stubAssistHook) Quote(ctx context.Context, setID, uid uint) (float64, uint, uint, uint, string, string, string, int, error) {
	return 0, 0, 0, 0, "", "", "", 0, nil
}
func (h *stubAssistHook) ProductCost(ctx context.Context, productID uint) (float64, error) {
	return 0, nil
}
func (h *stubAssistHook) ReserveStock(ctx context.Context, assistID uint, num int) error {
	return nil
}
func (h *stubAssistHook) RestoreStock(ctx context.Context, assistID uint, num int) error {
	h.restoreID, h.restoreNum = assistID, num
	return nil
}
func (h *stubAssistHook) MarkSetPaid(ctx context.Context, setID uint) error { return nil }

func TestRestoreActivityStock_PresellAndAssist(t *testing.T) {
	ph := &stubPresellHook{}
	ah := &stubAssistHook{}
	svc := &Service{presell: ph, assist: ah}
	ctx := context.Background()

	if err := svc.restoreActivityStock(ctx, &OrderProduct{
		ProductType: ActivityTypePresell, ActivityID: 7, ProductNum: 2,
	}); err != nil {
		t.Fatal(err)
	}
	if ph.restoreID != 7 || ph.restoreNum != 2 {
		t.Fatalf("presell restore id=%d num=%d", ph.restoreID, ph.restoreNum)
	}

	if err := svc.restoreActivityStock(ctx, &OrderProduct{
		ProductType: ActivityTypeAssist, ProductNum: 3,
		ProductInfo: `{"product_assist_id":9}`,
	}); err != nil {
		t.Fatal(err)
	}
	if ah.restoreID != 9 || ah.restoreNum != 3 {
		t.Fatalf("assist restore id=%d num=%d", ah.restoreID, ah.restoreNum)
	}

	ah.restoreID, ah.restoreNum = 0, 0
	if err := svc.restoreActivityStock(ctx, &OrderProduct{ProductType: 0, ProductNum: 1}); err != nil {
		t.Fatal(err)
	}
	if ah.restoreID != 0 || ph.restoreID != 7 {
		t.Fatalf("normal line should not touch assist; assist=%d", ah.restoreID)
	}
}
