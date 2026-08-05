package nativecatalog

import (
	"strings"
	"testing"
)

func TestProductBrandKeepsChineseFixtureAndBoundsLength(t *testing.T) {
	if !valid(saveRequest{CateID: 101, StoreName: "中文模拟针织衫", BrandName: "云锦织造", Price: 299, Stock: 10}) {
		t.Fatal("Chinese brand fixture should be accepted")
	}
	tooLong := strings.Repeat("品", 65)
	if valid(saveRequest{CateID: 101, StoreName: "中文模拟针织衫", BrandName: tooLong, Price: 299, Stock: 10}) {
		t.Fatal("brand name over 64 runes must be rejected")
	}
	if valid(saveRequest{CateID: 101, StoreName: "中文模拟针织衫", BrandName: "云锦\n织造", Price: 299, Stock: 10}) {
		t.Fatal("brand name containing a line break must be rejected")
	}
}
