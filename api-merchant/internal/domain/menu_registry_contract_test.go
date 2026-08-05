package domain

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// TestMerchantRoutableMenuHasConcreteVbenComponent ensures every seeded
// store console page (is_route=1 with an ecrm view path) is registered in
// admin-merchant registry.ts. Unregistered leaves must stay hidden, never
// fall back to placeholder.
func TestMerchantRoutableMenuHasConcreteVbenComponent(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	seed, err := os.ReadFile(filepath.Join(root, "sql/merchant/init_data.sql"))
	if err != nil {
		t.Fatalf("read merchant menu seed: %v", err)
	}
	registry, err := os.ReadFile(filepath.Join(root, "admin-merchant/src/views/ecrm/registry.ts"))
	if err != nil {
		t.Fatalf("read merchant Vben registry: %v", err)
	}
	// Match: (id,parent,'code','name','/path','views/ecrm/...',...,is_menu,is_route,...)
	// is_route=1 pages only.
	pattern := regexp.MustCompile(
		`\(\d+,\d+,'[^']+','[^']+','(/[^']+)','views/ecrm/[^']*',[^,]*,1,1,\d+,1\)`,
	)
	paths := pattern.FindAllStringSubmatch(string(seed), -1)
	if len(paths) == 0 {
		t.Fatal("merchant menu seed contains no routable page routes")
	}
	regText := string(registry)
	for _, match := range paths {
		path := match[1]
		if !strings.Contains(regText, "'"+path+"': 'ecrm/") {
			t.Fatalf("merchant menu page %q has no concrete Vben component in registry.ts", path)
		}
	}
}
