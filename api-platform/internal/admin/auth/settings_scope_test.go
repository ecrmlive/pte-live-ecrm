package auth

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestChineseCustomerServiceAdminFixture(t *testing.T) {
	// 中文模拟数据必须可完整进入后台账号配置，不使用真实姓名、手机号或凭据。
	raw := []byte(`{"account":"cs_demo_zh","real_name":"客服张敏","roles":"customer_service","service_store_ids":"1001,1002"}`)
	var req adminSaveRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		t.Fatal(err)
	}
	if got := normalizedRoleCodes(req); !reflect.DeepEqual(got, []string{"customer_service"}) {
		t.Fatalf("roles = %v", got)
	}
	if got, err := parseIDs(req.ServiceStoreIDs, "客服授权店铺 ID"); err != nil || !reflect.DeepEqual(got, []uint64{1001, 1002}) {
		t.Fatalf("store ids = %v, err = %v", got, err)
	}
}

func TestChineseMerchantAdminFixture(t *testing.T) {
	raw := []byte(`{"account":"merchant_demo_zh","real_name":"商户运营李航","roles":"merchant","merchant_ids":"2001,2002"}`)
	var req adminSaveRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		t.Fatal(err)
	}
	if got := normalizedRoleCodes(req); !reflect.DeepEqual(got, []string{"merchant"}) {
		t.Fatalf("roles = %v", got)
	}
	if got, err := parseIDs(req.MerchantIDs, "授权商户 ID"); err != nil || !reflect.DeepEqual(got, []uint64{2001, 2002}) {
		t.Fatalf("merchant ids = %v, err = %v", got, err)
	}
}

func TestParseServiceQueueScope(t *testing.T) {
	raw, err := json.Marshal(serviceQueueScope{StoreIDs: []uint64{1001, 1002}})
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseScopeIDs(raw, true)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []uint64{1001, 1002}) {
		t.Fatalf("store ids = %v", got)
	}
}

func TestParseLegacyServiceQueueScope(t *testing.T) {
	got, err := parseScopeIDs(json.RawMessage(`[1001,1002]`), true)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []uint64{1001, 1002}) {
		t.Fatalf("legacy store ids = %v", got)
	}
}

func TestParseIDsRejectsInvalidCustomerServiceStore(t *testing.T) {
	if _, err := parseIDs("1001,中文店铺", "客服授权店铺 ID"); err == nil {
		t.Fatal("invalid store id must be rejected")
	}
}

func TestAllowsAdminDeletionProtectsCurrentAndLastPlatform(t *testing.T) {
	if err := allowsAdminDeletion(9001, 9001, []string{"customer_service"}, 1); err == nil {
		t.Fatal("current logged-in administrator must not be deleted")
	}
	if err := allowsAdminDeletion(9001, 9002, []string{"platform"}, 0); err == nil {
		t.Fatal("last active platform administrator must not be deleted")
	}
	if err := allowsAdminDeletion(9001, 9302, []string{"customer_service"}, 1); err != nil {
		t.Fatalf("a historical virtual customer-service account should be logically deletable: %v", err)
	}
}
