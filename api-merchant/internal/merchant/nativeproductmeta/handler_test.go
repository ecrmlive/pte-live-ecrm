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

func TestParameterTemplateRouteContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(nil).Register(r)
	got := map[string]bool{}
	for _, route := range r.Routes() {
		got[route.Method+" "+route.Path] = true
	}
	for _, want := range []string{
		"GET /product/parameter-templates",
		"GET /product/parameter-templates/:id",
		"POST /product/parameter-templates",
		"PUT /product/parameter-templates/:id",
		"DELETE /product/parameter-templates/:id",
	} {
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

func TestNormalizeParameterTemplateChinese(t *testing.T) {
	req := int8(0)
	name, cateID, isRequired, raw, sort, ok := normalizeParameterTemplate(parameterTemplateInput{
		TemplateName: "测试",
		CateID:       1,
		IsRequired:   &req,
		Sort:         150,
		Params: []parameterItem{
			{Name: "材质", Values: []string{"棉", "涤纶"}, Required: 1, Sort: 10},
		},
	})
	if !ok || name != "测试" || cateID != 1 || isRequired != 0 || sort != 150 || len(raw) == 0 {
		t.Fatalf("normalize failed: ok=%v name=%q cate=%d req=%d sort=%d raw=%s", ok, name, cateID, isRequired, sort, string(raw))
	}
}
