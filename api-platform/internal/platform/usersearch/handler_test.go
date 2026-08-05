package usersearch

import "testing"

func TestSearchRecordFilters(t *testing.T) {
	if !validDateRange("2026-08-01", "2026-08-04") || validDateRange("2026-08-04", "2026-08-01") || validDateRange("2026/08/01", "") {
		t.Fatal("date range validation mismatch")
	}
	if !validSource("pc") || !validSource("h5") || validSource("unknown") {
		t.Fatal("source validation mismatch")
	}
}

func TestCSVCellPreventsFormulaInjection(t *testing.T) {
	if got := csvCell("=恶意公式"); got != "'=恶意公式" {
		t.Fatalf("got %q", got)
	}
	if got := csvCell("夏季凉鞋"); got != "夏季凉鞋" {
		t.Fatalf("got %q", got)
	}
}

func TestExportFingerprintAndDateEnd(t *testing.T) {
	in := query{UserID: 9101, Keyword: "夏季亚麻衬衫", Source: "pc", StartDate: "2026-08-03", EndDate: "2026-08-04"}
	if got := fingerprint(in); len(got) != 64 || got != fingerprint(in) {
		t.Fatalf("fingerprint=%q", got)
	}
	if got := nextDate("2026-08-04"); got != "2026-08-05 00:00:00" {
		t.Fatalf("next date=%q", got)
	}
}
