package feedback

import "testing"

func TestFeedbackResponseKeepsChineseText(t *testing.T) {
	got := (row{ID: 9, Type: "功能建议", Content: "希望商品搜索支持更多筛选", Status: "pending"}).view()
	if got["feedback_id"] != uint64(9) || got["type"] != "功能建议" || got["content"] != "希望商品搜索支持更多筛选" || got["status"] != "pending" {
		t.Fatalf("unexpected feedback response: %#v", got)
	}
}

func TestPositiveID(t *testing.T) {
	if _, ok := positiveID("0"); ok {
		t.Fatal("zero must be rejected")
	}
	if id, ok := positiveID("8"); !ok || id != 8 {
		t.Fatalf("id=%d ok=%v", id, ok)
	}
}
