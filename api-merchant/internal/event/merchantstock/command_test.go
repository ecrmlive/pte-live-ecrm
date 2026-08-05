package merchantstock

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestStockCommandValidationAndFailClosedDatabase(t *testing.T) {
	reserve := command{Action: "reserve", OrderID: 7001, StoreID: 1, MerchantSKUID: 61001, Quantity: 2, ExpiresAt: time.Now().UTC().Add(30 * time.Minute), IdempotencyKey: "reserve-order-7001-sku-61001"}
	if !validCommand(reserve) {
		t.Fatal("valid reserve command rejected")
	}
	reserve.Action, reserve.ExpiresAt = "confirm", time.Time{}
	if !validCommand(reserve) {
		t.Fatal("valid confirm command rejected")
	}
	reserve.Action = "restock"
	if !validCommand(reserve) {
		t.Fatal("valid refund restock command rejected")
	}
	reserve.IdempotencyKey = "短"
	if validCommand(reserve) {
		t.Fatal("weak idempotency key accepted")
	}
	wire, err := json.Marshal(command{Action: "reserve", OrderID: 7001, StoreID: 1, MerchantSKUID: 61001, Quantity: 1, ExpiresAt: time.Now().UTC().Add(30 * time.Minute), IdempotencyKey: "reserve-order-7001-sku-61001"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := ApplyCommand(context.Background(), nil, wire)
	if err == nil || result.Code != "failed" {
		t.Fatalf("unavailable merchant database must fail closed, got %#v, %v", result, err)
	}
}
