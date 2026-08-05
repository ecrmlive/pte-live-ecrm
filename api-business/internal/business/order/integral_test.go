package order

import "testing"

func TestQuoteIntegralUsesIntegerPolicyAndBalanceCaps(t *testing.T) {
	policy := IntegralPolicy{Enabled: true, PointsPerYuan: 100, MaxDeductionBps: 2000}
	quote := QuoteIntegral(policy, 5000, 9200, true)
	if quote.DiscountCents != 1840 || quote.PointsUsed != 1840 {
		t.Fatalf("quote=%+v, want 1840 cents / 1840 points", quote)
	}

	quote = QuoteIntegral(policy, 137, 9200, true)
	if quote.DiscountCents != 137 || quote.PointsUsed != 137 {
		t.Fatalf("balance cap quote=%+v, want 137 cents / 137 points", quote)
	}
}

func TestQuoteIntegralRejectsClientOnlyOrMalformedPolicy(t *testing.T) {
	for _, policy := range []IntegralPolicy{
		{},
		{Enabled: true, PointsPerYuan: 0, MaxDeductionBps: 2000},
		{Enabled: true, PointsPerYuan: 100, MaxDeductionBps: 0},
	} {
		if quote := QuoteIntegral(policy, 1000, 1000, true); quote != (IntegralQuote{}) {
			t.Fatalf("unexpected quote for policy %+v: %+v", policy, quote)
		}
	}
	if quote := QuoteIntegral(IntegralPolicy{Enabled: true, PointsPerYuan: 100, MaxDeductionBps: 2000}, 1000, 1000, false); quote != (IntegralQuote{}) {
		t.Fatalf("client request flag must be required: %+v", quote)
	}
}
