package trade

import "math"

// IntegralRate 积分兑换：多少积分 = 1 元；最多抵扣应付的该比例。
const (
	IntegralPerYuan   = 100
	IntegralMaxRatio  = 0.2
	GiveIntegralRatio = 1 // 实付 1 元赠 1 积分（向下取整）
)

type MerAmount struct {
	MerID      uint
	TotalPrice float64
	Coupon     float64
	PayPrice   float64
}

type PricingInput struct {
	MerAmounts       []MerAmount
	ShopCouponMerID  uint
	ShopCouponAmount float64
	PlatformCoupon   float64
	UseIntegral      int
	UserIntegral     int
}

type PricingResult struct {
	Merchants      []MerAmount
	TotalPrice     float64
	ShopCoupon     float64
	PlatformCoupon float64
	CouponPrice    float64
	Integral       int
	IntegralPrice  float64
	PayPrice       float64
	GiveIntegral   int
}

// ApplyPricing §8 顺序：商品价 → 店铺券 → 平台券 → 积分抵扣（无运费）。
func ApplyPricing(in PricingInput) PricingResult {
	out := PricingResult{Merchants: make([]MerAmount, len(in.MerAmounts))}
	copy(out.Merchants, in.MerAmounts)
	for i := range out.Merchants {
		out.TotalPrice += out.Merchants[i].TotalPrice
		out.Merchants[i].PayPrice = out.Merchants[i].TotalPrice
	}
	// 店铺券
	if in.ShopCouponAmount > 0 && in.ShopCouponMerID > 0 {
		for i := range out.Merchants {
			if out.Merchants[i].MerID != in.ShopCouponMerID {
				continue
			}
			d := in.ShopCouponAmount
			if d > out.Merchants[i].PayPrice {
				d = out.Merchants[i].PayPrice
			}
			out.Merchants[i].Coupon = d
			out.Merchants[i].PayPrice -= d
			out.ShopCoupon = d
			break
		}
	}
	// 平台券按各商户应付比例分摊
	remain := 0.0
	for i := range out.Merchants {
		remain += out.Merchants[i].PayPrice
	}
	plat := in.PlatformCoupon
	if plat > remain {
		plat = remain
	}
	if plat > 0 && remain > 0 {
		allocated := 0.0
		for i := range out.Merchants {
			if i == len(out.Merchants)-1 {
				d := plat - allocated
				if d > out.Merchants[i].PayPrice {
					d = out.Merchants[i].PayPrice
				}
				out.Merchants[i].Coupon += d
				out.Merchants[i].PayPrice -= d
				break
			}
			d := round2(plat * (out.Merchants[i].PayPrice / remain))
			if d > out.Merchants[i].PayPrice {
				d = out.Merchants[i].PayPrice
			}
			out.Merchants[i].Coupon += d
			out.Merchants[i].PayPrice -= d
			allocated += d
		}
		out.PlatformCoupon = plat
	}
	out.CouponPrice = out.ShopCoupon + out.PlatformCoupon
	payBeforeIntegral := 0.0
	for i := range out.Merchants {
		payBeforeIntegral += out.Merchants[i].PayPrice
	}
	// 积分
	maxYuan := round2(payBeforeIntegral * IntegralMaxRatio)
	wantYuan := float64(in.UseIntegral) / float64(IntegralPerYuan)
	if wantYuan > maxYuan {
		wantYuan = maxYuan
	}
	maxByBal := float64(in.UserIntegral) / float64(IntegralPerYuan)
	if wantYuan > maxByBal {
		wantYuan = maxByBal
	}
	wantYuan = round2(wantYuan)
	if wantYuan > payBeforeIntegral {
		wantYuan = payBeforeIntegral
	}
	out.IntegralPrice = wantYuan
	out.Integral = int(math.Ceil(wantYuan * float64(IntegralPerYuan)))
	if out.Integral > in.UserIntegral {
		out.Integral = in.UserIntegral
	}
	// 分摊积分抵扣
	if out.IntegralPrice > 0 && payBeforeIntegral > 0 {
		left := out.IntegralPrice
		for i := range out.Merchants {
			if i == len(out.Merchants)-1 {
				d := left
				if d > out.Merchants[i].PayPrice {
					d = out.Merchants[i].PayPrice
				}
				out.Merchants[i].PayPrice = round2(out.Merchants[i].PayPrice - d)
				break
			}
			d := round2(out.IntegralPrice * (out.Merchants[i].PayPrice / payBeforeIntegral))
			if d > out.Merchants[i].PayPrice {
				d = out.Merchants[i].PayPrice
			}
			out.Merchants[i].PayPrice = round2(out.Merchants[i].PayPrice - d)
			left = round2(left - d)
		}
	}
	for i := range out.Merchants {
		out.PayPrice += out.Merchants[i].PayPrice
	}
	out.PayPrice = round2(out.PayPrice)
	out.GiveIntegral = int(out.PayPrice * GiveIntegralRatio)
	return out
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
