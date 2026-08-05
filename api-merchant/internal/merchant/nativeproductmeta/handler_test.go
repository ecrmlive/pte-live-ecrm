package nativeproductmeta

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestLabelRouteContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(nil).Register(r)
	got := map[string]bool{}
	for _, route := range r.Routes() {
		got[route.Method+" "+route.Path] = true
	}
	for _, want := range []string{"GET /product/labels", "POST /product/labels", "PUT /product/labels/:id", "DELETE /product/labels/:id"} {
		if !got[want] {
			t.Errorf("missing %s", want)
		}
	}
}
func TestChineseLabelFixture(t *testing.T) {
	if !valid(tagInput{Name: "七禧甄选", Info: "中文模拟商品标签", Sort: 10}) {
		t.Fatal("中文商品标签应通过校验")
	}
}
