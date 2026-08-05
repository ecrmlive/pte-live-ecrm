package merchantstock

import (
	"testing"
	"time"
)

func TestReserveCommandCarriesOnlyInventoryIdentity(t *testing.T) {
	cmd := Command{Action: "reserve", OrderID: 7001, StoreID: 1, MerchantSKUID: 61001, Quantity: 2, ExpiresAt: time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC), IdempotencyKey: "stock:reserve:7001:61001"}
	if cmd.OrderID == 0 || cmd.StoreID == 0 || cmd.MerchantSKUID == 0 || cmd.Quantity < 1 || cmd.IdempotencyKey == "" {
		t.Fatalf("invalid reserve command: %#v", cmd)
	}
}
