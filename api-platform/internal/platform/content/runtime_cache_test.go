package content

import "testing"

func TestCacheScopePattern(t *testing.T) {
	tests := []struct {
		scope   string
		pattern string
		ok      bool
	}{
		{scope: "all", pattern: "ecrm:platform:*", ok: true},
		{scope: "store", pattern: "ecrm:platform:store:*", ok: true},
		{scope: "config", pattern: "ecrm:platform:config:*", ok: true},
		{scope: "replace_domain", pattern: "ecrm:platform:config:*", ok: true},
		{scope: "flushdb", ok: false},
	}
	for _, tt := range tests {
		pattern, ok := cacheScopePattern(tt.scope)
		if pattern != tt.pattern || ok != tt.ok {
			t.Fatalf("scope %q = (%q, %v), want (%q, %v)", tt.scope, pattern, ok, tt.pattern, tt.ok)
		}
	}
}

func TestNormalizeResourceDomain(t *testing.T) {
	valid := []string{"https://cdn.example.com", "http://localhost:9000/"}
	for _, raw := range valid {
		if _, ok := normalizeResourceDomain(raw); !ok {
			t.Fatalf("expected valid resource domain: %s", raw)
		}
	}
	for _, raw := range []string{"cdn.example.com", "https://cdn.example.com/path", "ftp://cdn.example.com", "https://cdn.example.com?x=1"} {
		if _, ok := normalizeResourceDomain(raw); ok {
			t.Fatalf("expected invalid resource domain: %s", raw)
		}
	}
}
