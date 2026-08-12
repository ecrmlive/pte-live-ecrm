package operationlog

import "testing"

func TestSplitAction(t *testing.T) {
	method, path := splitAction("POST /api/platform/v1/setting/admins/:id")
	if method != "POST" || path != "/api/platform/v1/setting/admins/:id" {
		t.Fatalf("splitAction() = %q, %q", method, path)
	}
	method, path = splitAction("更新配置")
	if method != "" || path != "" {
		t.Fatalf("non route action must stay empty, got %q, %q", method, path)
	}
}

func TestOperationPermissionName(t *testing.T) {
	if got := operationPermissionName("customer-service"); got != "客服管理" {
		t.Fatalf("permission name=%q", got)
	}
	if got := operationPermissionName("unknown"); got != "平台管理" {
		t.Fatalf("fallback permission name=%q", got)
	}
}
