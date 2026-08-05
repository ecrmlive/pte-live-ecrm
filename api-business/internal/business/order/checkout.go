// Package order contains the business-database order state machine.
package order

import (
	"errors"
	"sort"
)

var (
	ErrEmptyCart       = errors.New("请选择商品")
	ErrUnavailableCart = errors.New("存在失效或库存不足的商品")
	ErrMixedActivity   = errors.New("活动商品不能与其他类型商品混合购买")
	ErrMixedPaySubject = errors.New("平台商品与不同店铺商品不能合并支付，请分别结算")
)

// CartLine is the immutable consumption-view snapshot used to prepare an order.
// Amounts are integer cents: no float is allowed in order calculation.
type CartLine struct {
	CartID                  uint64
	ProductID               uint64
	MerchantSKUID           uint64
	SKUKey                  string
	SpecSnapshot            string
	MerchantID              uint64
	MerchantName            string
	StoreID                 uint64
	StoreName               string
	IntegralEnabled         bool
	IntegralPointsPerYuan   int64
	IntegralMaxDeductionBps int64
	Title                   string
	CoverURL                string
	ListCents               int64
	UnitCents               int64
	SVIPApplied             bool
	Quantity                int
	Stock                   int
	SaleStatus              int8
	ProductType             int8 // 0=normal; non-zero is an activity order line.
}

type StoreCheckout struct {
	MerchantID     uint64
	MerchantName   string
	StoreID        uint64
	StoreName      string
	IntegralPolicy IntegralPolicy
	Lines          []CartLine
	TotalCents     int64
	TotalQty       int
}

type Checkout struct {
	Stores     []StoreCheckout
	TotalCents int64
	TotalQty   int
}

// BuildCheckout validates a cart snapshot then deterministically groups it by store.
func BuildCheckout(lines []CartLine) (Checkout, error) {
	if len(lines) == 0 {
		return Checkout{}, ErrEmptyCart
	}
	activity := int8(-1)
	var storeID uint64
	storeSet := false
	byStore := make(map[uint64]*StoreCheckout, len(lines))
	for _, line := range lines {
		if line.Quantity < 1 || line.UnitCents < 0 || line.Stock < line.Quantity || line.SaleStatus != 1 {
			return Checkout{}, ErrUnavailableCart
		}
		if activity == -1 {
			activity = line.ProductType
		} else if activity != line.ProductType && (activity != 0 || line.ProductType != 0) {
			return Checkout{}, ErrMixedActivity
		}
		if !storeSet {
			storeID = line.StoreID
			storeSet = true
		} else if storeID != line.StoreID {
			return Checkout{}, ErrMixedPaySubject
		}
		store := byStore[line.StoreID]
		if store == nil {
			store = &StoreCheckout{MerchantID: line.MerchantID, MerchantName: line.MerchantName, StoreID: line.StoreID, StoreName: line.StoreName, IntegralPolicy: IntegralPolicy{Enabled: line.IntegralEnabled, PointsPerYuan: line.IntegralPointsPerYuan, MaxDeductionBps: line.IntegralMaxDeductionBps}, Lines: []CartLine{}}
			byStore[line.StoreID] = store
		}
		store.Lines = append(store.Lines, line)
		store.TotalCents += line.UnitCents * int64(line.Quantity)
		store.TotalQty += line.Quantity
	}
	out := Checkout{Stores: make([]StoreCheckout, 0, len(byStore))}
	for _, store := range byStore {
		out.Stores = append(out.Stores, *store)
		out.TotalCents += store.TotalCents
		out.TotalQty += store.TotalQty
	}
	sort.Slice(out.Stores, func(i, j int) bool { return out.Stores[i].StoreID < out.Stores[j].StoreID })
	return out, nil
}
