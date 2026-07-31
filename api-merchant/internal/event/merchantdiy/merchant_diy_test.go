package merchantdiy

import (
	"encoding/json"
	"testing"
)

func TestDIYEventEnvelopeContainsPublicPageFieldsOnly(t *testing.T) {
	payload, err := json.Marshal(map[string]any{"page_id": 12, "store_id": 3, "page_type": "home", "name": "演示店铺", "document": map[string]any{"page": map[string]any{}, "items": []any{}}, "status": "published", "is_active": true})
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(Envelope{EventID: 18, EventType: Upserted, Payload: payload})
	if err != nil {
		t.Fatal(err)
	}
	if len(wire) == 0 || containsSensitive(string(wire)) {
		t.Fatalf("unexpected outbox envelope: %s", wire)
	}
}

func containsSensitive(value string) bool {
	for _, word := range []string{"password", "secret", "token", "private_key"} {
		for i := 0; i+len(word) <= len(value); i++ {
			if value[i:i+len(word)] == word {
				return true
			}
		}
	}
	return false
}
