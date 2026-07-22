package openapi

import (
	"testing"
	"time"
)

func TestSignAndVerify(t *testing.T) {
	ak, sk, unique := "demo_mer1_ak", "demo_mer1_sk", "u-test-1"
	exp := time.Now().Unix()
	sig, err := SignPolicy(ak, sk, unique, exp)
	if err != nil {
		t.Fatal(err)
	}
	if len(sig) != 64 {
		t.Fatalf("sig len=%d want 64", len(sig))
	}
	if err := VerifySignature(ak, sk, unique, sig, exp, time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := VerifySignature(ak, sk, unique, "deadbeef", exp, time.Now()); err == nil {
		t.Fatal("want unauthorized")
	}
	old := exp - 400
	sigOld, _ := SignPolicy(ak, sk, unique, old)
	if err := VerifySignature(ak, sk, unique, sigOld, old, time.Now()); err == nil {
		t.Fatal("want expired")
	}
}
