package order

import (
	"errors"
	"testing"
)

func TestBuildCheckoutUsesCentsForOnePaymentSubject(t *testing.T) {
	checkout, err := BuildCheckout([]CartLine{
		{CartID: 1, StoreID: 2, MerchantID: 1, UnitCents: 1999, Quantity: 2, Stock: 2, SaleStatus: 1},
		{CartID: 2, StoreID: 2, MerchantID: 1, UnitCents: 100, Quantity: 3, Stock: 9, SaleStatus: 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(checkout.Stores) != 1 || checkout.TotalCents != 4298 || checkout.TotalQty != 5 {
		t.Fatalf("unexpected checkout: %#v", checkout)
	}
}

func TestBuildCheckoutRejectsMixedPaymentSubject(t *testing.T) {
	_, err := BuildCheckout([]CartLine{
		{CartID: 1, StoreID: 0, MerchantID: 0, UnitCents: 100, Quantity: 1, Stock: 1, SaleStatus: 1},
		{CartID: 2, StoreID: 2, MerchantID: 1, UnitCents: 100, Quantity: 1, Stock: 1, SaleStatus: 1},
	})
	if !errors.Is(err, ErrMixedPaySubject) {
		t.Fatalf("err = %v", err)
	}
}

func TestBuildCheckoutRejectsUnavailableAndMixedActivity(t *testing.T) {
	_, err := BuildCheckout([]CartLine{{StoreID: 1, UnitCents: 1, Quantity: 2, Stock: 1, SaleStatus: 1}})
	if !errors.Is(err, ErrUnavailableCart) {
		t.Fatalf("err = %v", err)
	}
	_, err = BuildCheckout([]CartLine{{StoreID: 1, UnitCents: 1, Quantity: 1, Stock: 1, SaleStatus: 1, ProductType: 0}, {StoreID: 1, UnitCents: 1, Quantity: 1, Stock: 1, SaleStatus: 1, ProductType: 1}})
	if !errors.Is(err, ErrMixedActivity) {
		t.Fatalf("err = %v", err)
	}
}
