package catalog

import "testing"

func TestParseHomeDecorationReadsPublishedDocumentShape(t *testing.T) {
	document := `{
		"page":{"params":{"title":"七禧商城"}},
		"items":[
			{"type":"banner","data":[
				{"imgName":"夏日精选","imgUrl":"/demo/hero.png","linkUrl":"/pages/goods/list"},
				{"imgName":"未配置图片","imgUrl":"","linkUrl":"/pages/goods/list"}
			]},
			{"type":"navBar","data":[{"text":"全部商品","linkUrl":"/pages/goods/list"}]}
		]
	}`

	got := parseHomeDecoration(4001, "默认标题", document)
	if got.PageID != 4001 || got.Title != "七禧商城" {
		t.Fatalf("page metadata = %#v", got)
	}
	if len(got.Banners) != 1 || got.Banners[0]["image"] != "/demo/hero.png" {
		t.Fatalf("banners = %#v", got.Banners)
	}
	if len(got.Menus) != 1 || got.Menus[0]["title"] != "全部商品" {
		t.Fatalf("menus = %#v", got.Menus)
	}
}

func TestParseHomeDecorationKeepsEmptyShapesOnMalformedDocument(t *testing.T) {
	got := parseHomeDecoration(4001, "七禧商城", "not-json")
	if got.Title != "七禧商城" || len(got.Banners) != 0 || len(got.Menus) != 0 {
		t.Fatalf("unexpected malformed result: %#v", got)
	}
}
