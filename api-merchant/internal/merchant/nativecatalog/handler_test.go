package nativecatalog

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProductConsoleRouteContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(nil, nil).Register(r)

	actual := map[string]bool{}
	for _, route := range r.Routes() {
		actual[route.Method+" "+route.Path] = true
	}
	for _, want := range []string{
		"GET /product-categories",
		"GET /products",
		"GET /products/recycle-bin",
		"POST /products",
		"PUT /products/:id",
		"DELETE /products/:id",
		"POST /products/:id/restore",
		"PUT /products/:id/show",
		"PUT /products/:id/stock",
	} {
		if !actual[want] {
			t.Errorf("missing route %s", want)
		}
	}
}

func TestChineseProductFixtureValidation(t *testing.T) {
	fixture := saveRequest{
		CateID:    10101,
		StoreName: "七禧夏日真丝衬衫",
		StoreInfo: "中文模拟数据：用于店铺商品、库存和回收站验收。",
		UnitName:  "件",
		Price:     299,
		Stock:     20,
	}
	if !valid(fixture) {
		t.Fatal("中文模拟商品应通过参数校验")
	}
	if !valid(saveRequest{CateID: 10101, StoreName: "七禧会员专享包", Price: 469, Stock: 20, SVIPPriceType: 2, SVIPPrice: 429}) {
		t.Fatal("固定会员价应通过参数校验")
	}
	if valid(saveRequest{CateID: 10101, StoreName: "七禧会员专享包", Price: 469, Stock: 20, SVIPPriceType: 2, SVIPPrice: 499}) {
		t.Fatal("高于售价的固定会员价必须被拒绝")
	}
	if got := statusName("-2"); got != "off_sale" {
		t.Fatalf("statusName(-2) = %q, want off_sale", got)
	}
}
