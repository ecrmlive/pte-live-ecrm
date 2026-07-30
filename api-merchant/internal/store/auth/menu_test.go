package auth

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMerchantMenuJSONUsesFrontendFieldNames(t *testing.T) {
	value, err := json.Marshal(menu{
		ID:       21,
		ParentID: 20,
		Name:     "订单管理",
		Path:     "/order/list",
	})
	if err != nil {
		t.Fatalf("marshal menu: %v", err)
	}
	body := string(value)
	for _, key := range []string{`"id":21`, `"parent_id":20`, `"name":"订单管理"`, `"path":"/order/list"`} {
		if !strings.Contains(body, key) {
			t.Fatalf("menu json missing %s: %s", key, body)
		}
	}
	if strings.Contains(body, `"ID"`) || strings.Contains(body, `"ParentID"`) {
		t.Fatalf("menu json must not expose Go field names: %s", body)
	}
}
