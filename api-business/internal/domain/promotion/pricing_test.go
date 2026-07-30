package promotion

import "testing"

func TestCalcQuote_StoreThenPlatform(t *testing.T) {
	mers := []MerTotal{
		{MerID: 1, TotalPrice: 40},
		{MerID: 2, TotalPrice: 20},
	}
	store := map[uint]CouponUser{
		1: {CouponUserID: 11, MerID: 1, CouponPrice: 3, UseMinPrice: 20},
	}
	platform := &CouponUser{CouponUserID: 99, CouponPrice: 5, UseMinPrice: 30}

	q := CalcQuote(mers, store, platform)
	if q.TotalPrice != 60 {
		t.Fatalf("total=%v", q.TotalPrice)
	}
	if q.MerStoreDiscount[1] != 3 {
		t.Fatalf("store disc=%v", q.MerStoreDiscount[1])
	}
	// after store: 37+20=57 >=30 → platform 5
	if q.PlatformDiscount != 5 {
		t.Fatalf("platform=%v", q.PlatformDiscount)
	}
	if q.CouponPrice != 8 {
		t.Fatalf("coupon=%v", q.CouponPrice)
	}
	if q.PayPrice != 52 {
		t.Fatalf("pay=%v", q.PayPrice)
	}
	shareSum := q.MerPlatformShare[1] + q.MerPlatformShare[2]
	if round2(shareSum) != 5 {
		t.Fatalf("share sum=%v", shareSum)
	}
}

func TestCalcQuote_MinNotMet(t *testing.T) {
	mers := []MerTotal{{MerID: 1, TotalPrice: 10}}
	platform := &CouponUser{CouponUserID: 1, CouponPrice: 5, UseMinPrice: 30}
	q := CalcQuote(mers, nil, platform)
	if q.PlatformDiscount != 0 || q.PayPrice != 10 {
		t.Fatalf("unexpected %#v", q)
	}
}
