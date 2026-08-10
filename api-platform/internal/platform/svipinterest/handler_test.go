package svipinterest

import "testing"

func TestInterestValidation(t *testing.T) {
	ok := &input{
		Name:        "会员专属价",
		DisplayName: "会员尊享专属价",
		Description: "会员尊享专属价",
		IconURL:     "https://picsum.photos/seed/qixi-svip-price-off/80/80",
		OnIconURL:   "https://picsum.photos/seed/qixi-svip-price-on/80/80",
		Status:      1,
		Sort:        10,
	}
	if !valid(ok, false) {
		t.Fatal("valid interest rejected")
	}
	if valid(&input{
		Name: "错误图标", DisplayName: "展示", Description: "简介",
		IconURL: "http://unsafe.example/icon.png", OnIconURL: "https://example.com/a.png",
		Status: 1,
	}, false) {
		t.Fatal("unsafe icon accepted")
	}
	if valid(&input{
		Name: "更新", DisplayName: "展示", Description: "简介",
		IconURL: "/demo/a.png", OnIconURL: "/demo/b.png", Status: 1,
	}, true) {
		t.Fatal("zero version update accepted")
	}
	if valid(&input{
		Name: "缺展示名", DisplayName: "", Description: "简介",
		IconURL: "/demo/a.png", OnIconURL: "/demo/b.png", Status: 1,
	}, false) {
		t.Fatal("empty display name accepted")
	}
}
