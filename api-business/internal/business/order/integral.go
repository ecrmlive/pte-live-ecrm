package order

// IntegralPolicy is the store-published rule used by a normal-order checkout.
// PointsPerYuan means how many points are worth one yuan.  MaxDeductionBps is
// expressed in basis points so pricing never relies on floating point values.
type IntegralPolicy struct {
	Enabled         bool
	PointsPerYuan   int64
	MaxDeductionBps int64
}

// IntegralQuote is the server-authoritative result after coupons.  PointsUsed
// is persisted with the order and must be debited atomically when it is made.
type IntegralQuote struct {
	DiscountCents int64
	PointsUsed    int64
}

// QuoteIntegral applies the CRMEB-style normal-order deduction after coupons.
// A disabled or malformed policy is deliberately a no-op: clients cannot turn
// points into a discount merely by sending a request flag.
func QuoteIntegral(policy IntegralPolicy, availablePoints, payableCents int64, requested bool) IntegralQuote {
	if !requested || !policy.Enabled || policy.PointsPerYuan <= 0 || policy.MaxDeductionBps <= 0 || payableCents <= 0 || availablePoints <= 0 {
		return IntegralQuote{}
	}
	if policy.MaxDeductionBps > 10000 {
		policy.MaxDeductionBps = 10000
	}
	maxByPolicy := payableCents * policy.MaxDeductionBps / 10000
	maxByPoints := availablePoints * 100 / policy.PointsPerYuan
	discount := maxByPolicy
	if maxByPoints < discount {
		discount = maxByPoints
	}
	if discount <= 0 {
		return IntegralQuote{}
	}
	// Round points up so the quoted cash discount can never exceed the point
	// value consumed. maxByPoints above guarantees this stays within balance.
	used := (discount*policy.PointsPerYuan + 99) / 100
	if used > availablePoints {
		return IntegralQuote{}
	}
	return IntegralQuote{DiscountCents: discount, PointsUsed: used}
}
