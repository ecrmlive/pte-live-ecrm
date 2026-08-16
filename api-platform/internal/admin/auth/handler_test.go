package auth

import "testing"

func TestProfileDerivesRegionCompatibilityFlag(t *testing.T) {
	h := NewHandler(nil, nil)
	user := &adminUser{ID: 12, Username: "region-admin", DisplayName: "区域管理员", DataScopeVersion: 4}
	got := h.profile(user, []string{"operations", "region"})
	if got.IsAgent != 1 {
		t.Fatalf("IsAgent = %d, want 1", got.IsAgent)
	}
	if got.Account != user.Username || got.AdminID != user.ID || got.DataScopeVersion != 4 {
		t.Fatalf("profile compatibility fields are not populated: %#v", got)
	}
}

func TestResolveMenuScopeKeepsPortalMenusSeparated(t *testing.T) {
	tests := []struct {
		name      string
		roleTypes []string
		requested string
		want      string
		wantErr   error
	}{
		{name: "platform", roleTypes: []string{"platform"}, want: "platform"},
		{name: "merchant", roleTypes: []string{"merchant"}, want: "merchant"},
		{name: "region", roleTypes: []string{"region"}, want: "region"},
		{name: "platform wins by default", roleTypes: []string{"merchant", "platform"}, want: "platform"},
		{name: "explicit merchant", roleTypes: []string{"merchant", "platform"}, requested: "merchant", want: "merchant"},
		{name: "forbid unassigned scope", roleTypes: []string{"platform"}, requested: "region", wantErr: errMenuScopeForbidden},
		{name: "reject invalid scope", roleTypes: []string{"platform"}, requested: "all", wantErr: errInvalidMenuScope},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveMenuScope(tt.roleTypes, tt.requested)
			if err != tt.wantErr {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("scope = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCircleAgentIDNullability(t *testing.T) {
	if got := derefCircleAgentID(nil); got != 0 {
		t.Fatalf("nil circle agent ID = %d, want 0", got)
	}
	id := uint64(23)
	if got := derefCircleAgentID(&id); got != 23 {
		t.Fatalf("circle agent ID = %d, want 23", got)
	}
	if got := nullableCircleAgentID(0); got != nil {
		t.Fatalf("zero circle agent ID = %#v, want nil", got)
	}
	if got := nullableCircleAgentID(id); got != id {
		t.Fatalf("circle agent ID = %#v, want %d", got, id)
	}
}

func TestLoginAuditValuesAreBoundedWithoutCredentials(t *testing.T) {
	if got := nullableAdminID(0); got != nil {
		t.Fatalf("unknown login must not invent admin id: %#v", got)
	}
	if got := nullableAdminID(12); got != uint64(12) {
		t.Fatalf("admin id=%#v", got)
	}
	if got := boundedAuditValue("  平台演示账号  ", 64); got != "平台演示账号" {
		t.Fatalf("audit value=%q", got)
	}
	if got := boundedAuditValue("123456", 4); got != "1234" {
		t.Fatalf("bounded value=%q", got)
	}
}
