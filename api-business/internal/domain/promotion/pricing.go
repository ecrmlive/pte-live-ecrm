package promotion

import "math"

// CalcQuote 店铺券 → 平台券（FUNCTIONAL-TRUTH §8.2 顺序子集，不含 SVIP/积分）。
func CalcQuote(mers []MerTotal, storeByMer map[uint]CouponUser, platform *CouponUser) QuoteResult {
	out := QuoteResult{
		MerStoreDiscount: map[uint]float64{},
		MerCouponUserID:  map[uint]uint{},
		MerPlatformShare: map[uint]float64{},
	}
	var total float64
	for _, m := range mers {
		total += m.TotalPrice
		out.MerStoreDiscount[m.MerID] = 0
		out.MerPlatformShare[m.MerID] = 0
	}
	out.TotalPrice = round2(total)

	remainByMer := map[uint]float64{}
	for _, m := range mers {
		remain := m.TotalPrice
		if cu, ok := storeByMer[m.MerID]; ok {
			if m.TotalPrice+1e-9 >= float64(cu.UseMinPrice) {
				d := math.Min(cu.CouponPrice, m.TotalPrice)
				d = round2(d)
				out.MerStoreDiscount[m.MerID] = d
				out.MerCouponUserID[m.MerID] = cu.CouponUserID
				remain -= d
			}
		}
		if remain < 0 {
			remain = 0
		}
		remainByMer[m.MerID] = remain
	}

	var afterStore float64
	for _, m := range mers {
		afterStore += remainByMer[m.MerID]
	}
	afterStore = round2(afterStore)

	if platform != nil && afterStore+1e-9 >= float64(platform.UseMinPrice) {
		pd := math.Min(platform.CouponPrice, afterStore)
		pd = round2(pd)
		out.PlatformCouponUserID = platform.CouponUserID
		out.PlatformDiscount = pd
		// 按剩余金额比例分摊到子单
		allocated := 0.0
		for i, m := range mers {
			base := remainByMer[m.MerID]
			var share float64
			if i == len(mers)-1 {
				share = round2(pd - allocated)
			} else if afterStore > 0 {
				share = round2(pd * base / afterStore)
				allocated += share
			}
			out.MerPlatformShare[m.MerID] = share
		}
	}

	out.CouponPrice = round2(sumMap(out.MerStoreDiscount) + out.PlatformDiscount)
	pay := out.TotalPrice - out.CouponPrice
	if pay < 0 {
		pay = 0
	}
	out.PayPrice = round2(pay)
	return out
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func sumMap(m map[uint]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}
