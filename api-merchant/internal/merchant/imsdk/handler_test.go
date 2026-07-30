package imsdk

import (
	"testing"
)

func TestIMEventPayloadContainsNoCredential(t *testing.T) {
	row := binding{MerchantID: 12, SDKAppID: "14001234", Name: "客服生产", Status: statusEnabled, IsActive: true, APIPublicURL: "https://im.example.test", WSPublicURL: "wss://im.example.test/ws", PTEProfileID: "merchant-12-production"}
	payload, err := imEventPayload(row)
	if err != nil {
		t.Fatal(err)
	}
	text := string(payload)
	if !contains(text, "sdk_app_id") || contains(text, "token") || contains(text, "user_sig") || contains(text, "secret") {
		t.Fatalf("unsafe IM outbox payload: %s", text)
	}
}

func contains(value, expected string) bool {
	for index := 0; index+len(expected) <= len(value); index++ {
		if value[index:index+len(expected)] == expected {
			return true
		}
	}
	return false
}
