package order

import (
	"encoding/json"
	"sort"
	"strings"
)

// normalizeSpecSnapshot keeps the exact JSON snapshot from the business SKU
// projection. A malformed legacy value is never allowed to break order writes.
func normalizeSpecSnapshot(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || !json.Valid([]byte(value)) {
		return "{}"
	}
	return value
}

// specText is a stable, human-readable representation used by cart checkout
// and order responses. Sorting prevents MySQL JSON key order from changing the
// customer-visible snapshot text.
func specText(snapshot string) string {
	specs := map[string]string{}
	if json.Unmarshal([]byte(normalizeSpecSnapshot(snapshot)), &specs) != nil {
		return "默认规格"
	}
	keys := make([]string, 0, len(specs))
	for key := range specs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		if value := strings.TrimSpace(specs[key]); value != "" {
			parts = append(parts, key+"："+value)
		}
	}
	if len(parts) == 0 {
		return "默认规格"
	}
	return strings.Join(parts, "；")
}
