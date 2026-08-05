package productmeta

import "testing"

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
	if validGuarantee(guaranteeInput{Name: "正品保障", Status: -1}) {
		t.Fatal("invalid guarantee status must fail")
	}
	valid := parameterInput{Name: "演示容量", Values: []string{" 小杯 ", "中杯", "大杯"}, Status: 1}
	if !validParameter(&valid) {
		t.Fatal("unique Chinese parameter values must pass")
	}
	if got := normalizeParameterValues(valid.Values); got[0] != "小杯" {
		t.Fatalf("values must be normalized: %#v", got)
	}
	if validParameter(&parameterInput{Name: "演示", Values: []string{"小杯", "小杯"}, Status: 1}) {
		t.Fatal("duplicate parameter values must fail")
	}
}
