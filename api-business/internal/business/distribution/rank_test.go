package distribution

import "testing"

func TestNormalizeRankMetric(t *testing.T) {
	cases := []struct {
		input string
		want  string
		ok    bool
	}{
		{"", rankMetricCommission, true},
		{"commission", rankMetricCommission, true},
		{"promoters", rankMetricPromoters, true},
		{"mobile", "", false},
	}
	for _, item := range cases {
		got, ok := normalizeRankMetric(item.input)
		if got != item.want || ok != item.ok {
			t.Fatalf("normalizeRankMetric(%q)=(%q,%t), want (%q,%t)", item.input, got, ok, item.want, item.ok)
		}
	}
}

func TestMaskNickname(t *testing.T) {
	cases := map[string]string{
		"张三":    "张**",
		"  李  ": "李**",
		"":      "推广用户",
	}
	for input, want := range cases {
		if got := maskNickname(input); got != want {
			t.Fatalf("maskNickname(%q)=%q, want %q", input, got, want)
		}
	}
}
