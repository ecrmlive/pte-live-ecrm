package authjwt

import (
	"strconv"
	"testing"
	"time"
)

func TestIssueCUserCarriesScopeChannelAndSubject(t *testing.T) {
	mgr := NewManager("test-secret", time.Hour, 2*time.Hour)
	pair, err := mgr.IssueCUser(42, "demo-account", "harmony")
	if err != nil {
		t.Fatalf("IssueCUser() error = %v", err)
	}
	claims, err := mgr.ParseExpect(pair.AccessToken, PortalApp, TokenAccess)
	if err != nil {
		t.Fatalf("ParseExpect() error = %v", err)
	}
	if claims.Scope != ScopeCUser || claims.PrincipalType != PrincipalCUser || claims.PrincipalID != 42 || claims.Roles[0] != "customer" || claims.ClientPlatform != "harmony" || claims.Channel != "harmony" || claims.Subject != strconv.Itoa(42) || claims.UID != 42 || claims.MerchantAppID != "" || claims.IMSDKAppID != "" || claims.AuthContext != ContextPlatform || claims.IdentityVersion != 1 || claims.ID == "" || claims.SessionID == "" || len(claims.Audience) != 1 || claims.Audience[0] != audienceForPortal(PortalApp) || claims.NotBefore == nil {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}

func TestParseExpectRejectsWrongAudience(t *testing.T) {
	mgr := NewManager("test-secret", time.Hour, 2*time.Hour)
	pair, err := mgr.IssueCUser(42, "demo-account", "pc")
	if err != nil {
		t.Fatalf("IssueCUser() error = %v", err)
	}
	if _, err := mgr.ParseExpect(pair.AccessToken, PortalPlatform, TokenAccess); err != ErrWrongPortal {
		t.Fatalf("ParseExpect() error = %v, want portal mismatch", err)
	}
}

func TestIssueAdminConsoleCarriesRolesAndDataScopeVersion(t *testing.T) {
	mgr := NewManager("test-secret", time.Hour, 2*time.Hour)
	pair, err := mgr.IssueAdminConsole(7, "operator", []string{"operations", "region"}, 3)
	if err != nil {
		t.Fatalf("IssueAdminConsole() error = %v", err)
	}
	claims, err := mgr.ParseExpect(pair.AccessToken, PortalPlatform, TokenAccess)
	if err != nil {
		t.Fatalf("ParseExpect() error = %v", err)
	}
	if claims.Scope != ScopeAdminConsole || claims.DataScopeVersion != 3 || len(claims.Roles) != 2 || claims.Roles[0] != "operations" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}

func TestIssueStoreConsoleCarriesMerchantStoreAndRole(t *testing.T) {
	mgr := NewManager("test-secret", time.Minute, time.Hour)
	pair, err := mgr.IssueStoreConsole(11, 22, 33, "qixi.store.demo.1", "1400000001", "shop-owner", "owner")
	if err != nil {
		t.Fatalf("IssueStoreConsole() error = %v", err)
	}
	claims, err := mgr.ParseExpect(pair.AccessToken, PortalMerchant, TokenAccess)
	if err != nil {
		t.Fatalf("ParseExpect() error = %v", err)
	}
	if claims.Scope != ScopeStoreConsole || claims.PrincipalType != PrincipalStoreAccount || claims.PrincipalID != 11 || claims.ClientPlatform != "merchant_web" || claims.AdminID != 11 || claims.MerID != 22 || claims.StoreID != 33 || claims.MerchantAppID != "qixi.store.demo.1" || claims.StoreAppID != claims.MerchantAppID || claims.IMSDKAppID != "1400000001" || claims.ID == "" || claims.SessionID == "" || len(claims.Roles) != 1 || claims.Roles[0] != "owner" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}

func TestIssueCUserStoreContextCarriesResolvedMerchantAndIMIdentity(t *testing.T) {
	mgr := NewManager("test-secret", time.Minute, time.Hour)
	pair, err := mgr.IssueCUserStoreContext(11, "consumer", "mini_program", 7, 22, 33, "qixi.store.demo.1", "1400000001")
	if err != nil {
		t.Fatalf("IssueCUserStoreContext() error = %v", err)
	}
	claims, err := mgr.ParseExpect(pair.AccessToken, PortalApp, TokenAccess)
	if err != nil {
		t.Fatalf("ParseExpect() error = %v", err)
	}
	if claims.AuthContext != ContextStore || claims.PrincipalType != PrincipalCUser || claims.UID != 11 || claims.MerID != 22 || claims.StoreID != 33 || claims.MerchantAppID != "qixi.store.demo.1" || claims.IMSDKAppID != "1400000001" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}
