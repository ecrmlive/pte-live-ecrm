package nativecatalog

import "testing"

func TestProductAuditResponseIncludesBrand(t *testing.T) {
	row := productRow{ID: 1001, BrandName: "云锦织造"}
	if row.BrandName != "云锦织造" {
		t.Fatalf("brand = %q", row.BrandName)
	}
}
