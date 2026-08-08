package productmeta

import "testing"

func TestNormalizeStoreParams(t *testing.T) {
	ok, valid := normalizeStoreParams([]storeParameterItem{
		{Name: "材质", Values: []string{" 棉 ", "涤纶"}, Required: 1, Sort: 10},
	})
	if !valid || len(ok) != 1 || ok[0].Values[0] != "棉" {
		t.Fatalf("normalize store params failed: %#v valid=%v", ok, valid)
	}
	if _, valid := normalizeStoreParams([]storeParameterItem{
		{Name: "材质", Values: []string{"棉", "棉"}},
	}); valid {
		t.Fatal("duplicate values must fail")
	}
	if _, valid := normalizeStoreParams(nil); valid {
		t.Fatal("empty params must fail")
	}
}

func TestProductMetadataValidationRejectsUnsafeOrAmbiguousInputs(t *testing.T) {
	if !validLabel(labelInput{Name: "虚构中文标签", Status: 1}) {
		t.Fatal("valid label must pass")
	}
	if validLabel(labelInput{Name: "", Status: 1}) || validLabel(labelInput{Name: "演示", Status: 2}) {
		t.Fatal("invalid labels must fail")
	}
	if !validGuarantee(guaranteeInput{Name: "正品保障", Content: "虚构中文保障说明", Status: 1}) {
		t.Fatal("valid guarantee must pass")
	}
	if validGuarantee(guaranteeInput{Name: "正品保障", Content: "", Status: 1}) {
		t.Fatal("empty guarantee content must fail")
	}
	if validGuarantee(guaranteeInput{Name: "正品保障", Content: "说明", Status: -1}) {
		t.Fatal("invalid guarantee status must fail")
	}
	valid := parameterInput{
		Name:    "演示容量",
		CateIDs: []uint64{7605},
		Params: []parameterItem{
			{Name: "容量", Values: []string{" 小杯 ", "中杯", "大杯"}, Required: 0, Sort: 10},
		},
		Status: 1,
	}
	if !validParameter(&valid) {
		t.Fatal("unique Chinese parameter values must pass")
	}
	if got := normalizeParameterItems(valid.Params); got[0].Values[0] != "小杯" {
		t.Fatalf("values must be normalized: %#v", got)
	}
	if validParameter(&parameterInput{
		Name:    "演示",
		CateIDs: []uint64{1},
		Params:  []parameterItem{{Name: "容量", Values: []string{"小杯", "小杯"}}},
		Status:  1,
	}) {
		t.Fatal("duplicate parameter values must fail")
	}
	if validParameter(&parameterInput{
		Name:   "无分类",
		Params: []parameterItem{{Name: "容量", Values: []string{"小杯"}}},
		Status: 1,
	}) {
		t.Fatal("platform template must require categories")
	}
	if !validPriceRule(&priceRuleInput{
		Name:    "集成灶",
		CateIDs: []uint64{1, 2},
		Content: "<p><strong>演示价格说明</strong></p>",
		Status:  1,
	}) {
		t.Fatal("valid price rule must pass")
	}
	if !validPriceRule(&priceRuleInput{
		Name:    "全部商品",
		Content: "<p>默认全部</p>",
		Status:  1,
	}) {
		t.Fatal("empty categories (default all) must pass")
	}
	if validPriceRule(&priceRuleInput{Name: "空内容", Content: "<p></p>", Status: 1}) {
		t.Fatal("empty rich-text content must fail")
	}
	if validPriceRule(&priceRuleInput{Name: "", Content: "<p>x</p>", Status: 1}) {
		t.Fatal("empty name must fail")
	}
}
