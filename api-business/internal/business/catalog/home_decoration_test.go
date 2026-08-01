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
			{"type":"navBar","data":[{"text":"全部商品","linkUrl":"/pages/goods/list"}]},
			{"type":"product","params":{"source":"auto","auto":{"category":101,"showNum":8,"productSort":"sales"}}},
			{"type":"product","params":{"source":"auto","auto":{"category":102,"showNum":99,"productSort":"price"}}},
			{"type":"product","params":{"source":"choice","auto":{"category":103,"showNum":8,"productSort":"sales"}}}
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
	if len(got.DisplayTypes) != 2 {
		t.Fatalf("display types = %#v", got.DisplayTypes)
	}
	if got.DisplayTypes[0] != (homeDisplayType{CategoryID: 101, InitialLimit: 8, Sort: "sales"}) {
		t.Fatalf("first display type = %#v", got.DisplayTypes[0])
	}
	if got.DisplayTypes[1] != (homeDisplayType{CategoryID: 102, InitialLimit: 24, Sort: "price"}) {
		t.Fatalf("second display type = %#v", got.DisplayTypes[1])
	}
}

func TestParseHomeDecorationKeepsEmptyShapesOnMalformedDocument(t *testing.T) {
	got := parseHomeDecoration(4001, "七禧商城", "not-json")
	if got.Title != "七禧商城" || len(got.Banners) != 0 || len(got.Menus) != 0 || len(got.DisplayTypes) != 0 {
		t.Fatalf("unexpected malformed result: %#v", got)
	}
}
