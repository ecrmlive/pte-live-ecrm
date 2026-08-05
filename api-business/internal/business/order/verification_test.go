package order

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestNormalizeDeliveryType(t *testing.T) {
	cases := []struct {
		in   string
		want string
		err  bool
	}{
		{"", "express", false},
		{"EXPRESS", "express", false},
		{"pickup", "pickup", false},
		{"1", "pickup", false},
		{"service", "service", false},
		{"city", "city", false},
		{"drone", "", true},
	}
	for _, tc := range cases {
		got, err := normalizeDeliveryType(tc.in)
		if tc.err {
			if !errors.Is(err, ErrDeliveryType) {
				t.Fatalf("normalizeDeliveryType(%q) err=%v, want ErrDeliveryType", tc.in, err)
			}
			continue
		}
		if err != nil || got != tc.want {
			t.Fatalf("normalizeDeliveryType(%q)=(%q,%v), want %q", tc.in, got, err, tc.want)
		}
	}
}

func TestHashVerifyCodeStable(t *testing.T) {
	a := hashVerifyCode(" 1234567890 ")
	b := hashVerifyCode("1234567890")
	if a == "" || a != b || len(a) != 64 {
		t.Fatalf("hash unstable: %q vs %q", a, b)
	}
}

func TestGenerateVerifyCodeFormat(t *testing.T) {
	code := generateVerifyCode()
	if len(code) != 10 {
		t.Fatalf("code length=%d want 10 (%q)", len(code), code)
	}
	for _, ch := range code {
		if ch < '0' || ch > '9' {
			t.Fatalf("non-digit in code %q", code)
		}
	}
}

func TestIssueVerificationsRejectsEmptyInput(t *testing.T) {
	if err := issueVerificationsForPaidGroup(nil, 0); !errors.Is(err, gorm.ErrInvalidData) {
		t.Fatalf("err=%v, want gorm.ErrInvalidData", err)
	}
}
