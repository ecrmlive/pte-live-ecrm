package auth

import (
	"strings"
	"testing"
)

func TestNormalizeNicknameKeepsChineseText(t *testing.T) {
	value, err := normalizeNickname("  陈小满  ")
	if err != nil || value != "陈小满" {
		t.Fatalf("normalizeNickname() = %q, %v", value, err)
	}
	for _, raw := range []string{"", "   ", strings.Repeat("测", 65)} {
		if _, err := normalizeNickname(raw); err == nil {
			t.Fatalf("%q should be rejected", raw)
		}
	}
}
