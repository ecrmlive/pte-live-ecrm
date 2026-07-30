package merchantim

import (
	"encoding/json"
	"testing"
)

func TestEnvelopeCarriesNoSecret(t *testing.T) {
	payload, err := json.Marshal(Payload{MerchantID: 9, SDKAppID: "140000009", Status: "enabled", IsActive: true, APIPublicURL: "https://im.example.test", WSPublicURL: "wss://im.example.test/ws", PTEProfileID: "merchant-9"})
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(Envelope{EventID: 1, EventType: Activated, Payload: payload})
	if err != nil {
		t.Fatal(err)
	}
	if string(wire) == "" || contains(string(wire), "token") || contains(string(wire), "user_sig") || contains(string(wire), "secret") {
		t.Fatalf("unsafe event: %s", wire)
	}
}

func contains(value, expected string) bool {
	for i := 0; i+len(expected) <= len(value); i++ {
		if value[i:i+len(expected)] == expected {
			return true
		}
	}
	return false
}
