package auth

import (
	"strings"
	"testing"
)

func TestValidNewPassword(t *testing.T) {
	for _, value := range []string{"安全密码1234", "Abcd1234"} {
		if !validNewPassword(value) {
			t.Fatalf("expected valid password %q", value)
		}
	}
	for _, value := range []string{"短密码1", "", strings.Repeat("a", 129)} {
		if validNewPassword(value) {
			t.Fatalf("expected invalid password %q", value)
		}
	}
}
