package paynotify

import "testing"

func TestMakeAndVerifyToken(t *testing.T) {
	secret := "test-secret"
	tok := MakeToken(secret, "wechat", "GO123", 9, 42, 12.5)
	if tok == "" {
		t.Fatal("empty token")
	}
	if !VerifyToken(secret, "wechat", "GO123", 9, 42, 12.5, tok) {
		t.Fatal("verify failed")
	}
	if VerifyToken(secret, "alipay", "GO123", 9, 42, 12.5, tok) {
		t.Fatal("channel mismatch should fail")
	}
	if VerifyToken("other", "wechat", "GO123", 9, 42, 12.5, tok) {
		t.Fatal("secret mismatch should fail")
	}
	if VerifyToken(secret, "wechat", "GO123", 9, 42, 12.51, tok) {
		t.Fatal("amount mismatch should fail")
	}
}
