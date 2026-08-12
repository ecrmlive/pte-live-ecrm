package nativeconfigitem

import "testing"

func TestNormalizeGroupInput(t *testing.T) {
	name, key, description, fields, sort, err := normalizeGroupInput(dataGroupInput{
		Name:        "首页轮播图",
		GroupKey:    "home_banner",
		Description: "首页展示图片",
		Fields: []map[string]any{{
			"field": "image",
			"name":  "图片",
			"type":  "image",
		}},
		Sort: groupDataIntPtr(10),
	}, 0)
	if err != nil {
		t.Fatalf("normalize group input: %v", err)
	}
	if name != "首页轮播图" || key != "home_banner" || description != "首页展示图片" || sort != 10 || fields == "" {
		t.Fatalf("unexpected normalized group: %q %q %q %q %d", name, key, description, fields, sort)
	}

	_, _, _, _, _, err = normalizeGroupInput(dataGroupInput{Name: "测试", GroupKey: "中文key"}, 0)
	if err == nil {
		t.Fatal("expected invalid group key error")
	}
}

func TestNormalizeGroupItemInput(t *testing.T) {
	data, sort, status, err := normalizeGroupItemInput(dataGroupItemInput{
		Data:   map[string]any{"title": "热销"},
		Sort:   groupDataIntPtr(2),
		Status: groupDataIntPtr(1),
	}, 0, 0)
	if err != nil {
		t.Fatalf("normalize group item input: %v", err)
	}
	if data != `{"title":"热销"}` || sort != 2 || status != 1 {
		t.Fatalf("unexpected normalized item: %q %d %d", data, sort, status)
	}
}

func groupDataIntPtr(value int) *int { return &value }
