package nativeproductmeta

import "testing"

func TestChineseCategoryFixture(t *testing.T) {
	if !validCategory(categoryInput{Name: "七禧夏装", Sort: 20}) {
		t.Fatal("中文商品分类应通过校验")
	}
}
