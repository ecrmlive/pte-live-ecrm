package trade

import "testing"

func TestApplyPricing_ShopThenPlatformThenIntegral(t *testing.T) {
	in := PricingInput{
		MerAmounts: []MerAmount{
			{MerID: 1, TotalPrice: 60},
			{MerID: 2, TotalPrice: 40},
		},
		ShopCouponMerID:  1,
		ShopCouponAmount: 5,
		PlatformCoupon:   3,
		UseIntegral:      1000, // 最多 10 元，但上限 20%
		UserIntegral:     5000,
	}
	// 商品 100 → 店券5 → 95 → 平台3 → 92 → 积分最多 18.4，请求10 → 10
	out := ApplyPricing(in)
	if out.TotalPrice != 100 {
		t.Fatalf("total=%v", out.TotalPrice)
	}
	if out.ShopCoupon != 5 || out.PlatformCoupon != 3 {
		t.Fatalf("coupons shop=%v plat=%v", out.ShopCoupon, out.PlatformCoupon)
	}
	if out.CouponPrice != 8 {
		t.Fatalf("coupon_price=%v", out.CouponPrice)
	}
	if out.IntegralPrice != 10 {
		t.Fatalf("integral_price=%v want 10", out.IntegralPrice)
	}
	if out.PayPrice != 82 {
		t.Fatalf("pay=%v want 82", out.PayPrice)
	}
}

func TestActivityTypeConstants_CRMEB(t *testing.T) {
	// CRMEB：0普通 1秒杀 2预售 3助力 4拼团；积分商城 activity_type=20
	if ActivityTypeSeckill != 1 {
		t.Fatalf("seckill=%d", ActivityTypeSeckill)
	}
	if ActivityTypePresell != 2 {
		t.Fatalf("presell=%d", ActivityTypePresell)
	}
	if ActivityTypeAssist != 3 {
		t.Fatalf("assist=%d", ActivityTypeAssist)
	}
	if ActivityTypeCombination != 4 {
		t.Fatalf("combination=%d want 4 (not 2=presell)", ActivityTypeCombination)
	}
	if ActivityTypePoints != 20 {
		t.Fatalf("points=%d", ActivityTypePoints)
	}
	if OrderStatusGrouping != 9 {
		t.Fatalf("grouping status=%d", OrderStatusGrouping)
	}
	if OrderStatusAwaitFinal != 10 {
		t.Fatalf("await final=%d", OrderStatusAwaitFinal)
	}
	if OrderStatusFinalTimeout != 11 {
		t.Fatalf("final timeout=%d", OrderStatusFinalTimeout)
	}
}

func TestApplyPricing_NoCoupon(t *testing.T) {
	out := ApplyPricing(PricingInput{
		MerAmounts: []MerAmount{{MerID: 1, TotalPrice: 29.9}},
	})
	if out.PayPrice != 29.9 || out.CouponPrice != 0 {
		t.Fatalf("got pay=%v coupon=%v", out.PayPrice, out.CouponPrice)
	}
}

func TestAssistIDFromProductInfo(t *testing.T) {
	if got := assistIDFromProductInfo(`{"product_assist_id":12,"product_assist_set_id":3}`); got != 12 {
		t.Fatalf("got=%d", got)
	}
	if got := assistIDFromProductInfo(`{}`); got != 0 {
		t.Fatalf("empty got=%d", got)
	}
	if got := assistIDFromProductInfo(""); got != 0 {
		t.Fatalf("blank got=%d", got)
	}
}
