package operationlog

import "testing"

func TestDateRange(t *testing.T) {
	start, end, ok := dateRange("2026-08-01", "2026-08-01")
	if !ok || start.IsZero() || end.Sub(start).Hours() != 24 {
		t.Fatalf("range start=%v end=%v ok=%v", start, end, ok)
	}
	for _, values := range [][2]string{{"bad", ""}, {"2026-08-03", "2026-08-01"}} {
		if _, _, ok := dateRange(values[0], values[1]); ok {
			t.Fatalf("invalid range accepted: %v", values)
		}
	}
}
