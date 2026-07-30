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
