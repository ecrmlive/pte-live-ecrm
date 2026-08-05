package order

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSpecTextKeepsChineseSKUValues(t *testing.T) {
	if got := specText(`{"颜色":"晨雾灰","尺码":"40"}`); got != "尺码：40；颜色：晨雾灰" {
		t.Fatalf("spec text = %q", got)
	}
	if got := normalizeSpecSnapshot("not-json"); got != "{}" {
		t.Fatalf("invalid snapshot = %q", got)
	}
}

func TestCheckResponseKeepsSelectedSKUText(t *testing.T) {
	response := checkResponse(Checkout{Stores: []StoreCheckout{{MerchantID: 1, Lines: []CartLine{{CartID: 7, ProductID: 1003, SKUKey: "61003", Title: "轻量缓震跑步鞋", SpecSnapshot: `{"颜色":"晨雾灰","尺码":"40"}`, UnitCents: 36900, Quantity: 1, Stock: 24, SaleStatus: 1}}}}}, CouponPricing{})
	merchants := response["merchants"].([]gin.H)
	items := merchants[0]["items"].([]gin.H)
	if got := items[0]["spec_text"]; got != "尺码：40；颜色：晨雾灰" {
		t.Fatalf("checkout SKU text = %q", got)
	}
	if got := response["integral_price"]; got != float64(0) {
		t.Fatalf("integral_price=%v, want 0", got)
	}
	quoted := checkResponse(Checkout{}, CouponPricing{}, IntegralQuote{DiscountCents: 1840, PointsUsed: 1840})
	if got := quoted["integral_price"]; got != 18.4 {
		t.Fatalf("quoted integral_price=%v, want 18.4", got)
	}
	if got := quoted["integral"]; got != int64(1840) {
		t.Fatalf("quoted integral=%v, want 1840", got)
	}
}
