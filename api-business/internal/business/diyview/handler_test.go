package diyview

import (
	"encoding/json"
	"testing"
)

func TestPayloadKeepsRenderedDocumentAndLegacyBlocks(t *testing.T) {
	document, err := json.Marshal(map[string]any{
		"page": map[string]any{"params": map[string]any{"name": "测试店铺"}},
		"items": []map[string]any{
			{"type": "banner", "data": []map[string]any{{"imgName": "夏日上新", "imgUrl": "https://cdn.example.test/banner.png", "linkUrl": "/pages/goods/list"}}},
			{"type": "navBar", "data": []map[string]any{{"text": "全部商品", "imgUrl": "", "linkUrl": "/pages/goods/list"}}},
		},
		"_qixi": map[string]any{"title": "测试店铺首页"},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := payload(pageView{PageID: 7, StoreID: 3, PageType: "home", Name: "测试页", Document: document})
	if got["title"] != "测试店铺首页" || got["store_id"] != uint64(3) {
		t.Fatalf("unexpected payload: %#v", got)
	}
	if len(got["banners"].([]map[string]any)) != 1 || len(got["menus"].([]map[string]any)) != 1 {
		t.Fatalf("legacy blocks missing: %#v", got)
	}
}
