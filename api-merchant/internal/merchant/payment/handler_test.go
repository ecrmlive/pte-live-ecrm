package payment

import "testing"

func TestNormalizeChannel(t *testing.T) {
	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{input: "wechat", want: channelWechat, ok: true},
		{input: " ALIPAY ", want: channelAlipay, ok: true},
		{input: "mock", ok: false},
	}
	for _, tt := range tests {
		got, ok := normalizeChannel(tt.input)
		if got != tt.want || ok != tt.ok {
			t.Fatalf("normalizeChannel(%q) = (%q, %v), want (%q, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}
