package content

import (
	"context"
	"testing"
)

func TestNormalizeWechatMenus(t *testing.T) {
	ok, err := normalizeWechatMenus([]WechatMenuButton{
		{Name: "商城", SubButton: []WechatMenuButton{
			{Name: "首页", Type: "view", URL: "https://example.invalid/"},
		}},
		{Name: "客服", Type: "click", Key: "CS"},
	})
	if err != nil || len(ok) != 2 {
		t.Fatalf("normalize ok failed: %#v %v", ok, err)
	}
	if ok[0].Type != "" || len(ok[0].SubButton) != 1 {
		t.Fatalf("parent with children should clear leaf fields: %#v", ok[0])
	}

	if _, err := normalizeWechatMenus(nil); err == nil {
		t.Fatal("empty menus must fail")
	}
	if _, err := normalizeWechatMenus([]WechatMenuButton{
		{Name: "坏链", Type: "view", URL: "javascript:alert(1)"},
	}); err == nil {
		t.Fatal("unsafe url must fail")
	}
}

func TestGetSaveWechatMenusRoundTrip(t *testing.T) {
	store := &mallSettingStore{}
	svc := NewService(store)
	saved, err := svc.SaveWechatMenus(context.Background(), []WechatMenuButton{
		{Name: "商城", Type: "view", URL: "https://example.invalid/mall"},
	})
	if err != nil || len(saved) != 1 {
		t.Fatalf("save failed: %#v %v", saved, err)
	}
	got, err := svc.GetWechatMenus(context.Background())
	if err != nil || len(got) != 1 || got[0].Name != "商城" {
		t.Fatalf("get failed: %#v %v", got, err)
	}
}
