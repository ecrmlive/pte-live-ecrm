package feedbackmoderation

import "testing"

func TestFeedbackStates(t *testing.T) {
	if n, ok := next("pending", "reply"); !ok || n != "replied" {
		t.Fatal("pending must reply")
	}
	if n, ok := next("replied", "close"); !ok || n != "closed" {
		t.Fatal("replied must close")
	}
	if _, ok := next("closed", "reply"); ok {
		t.Fatal("closed feedback must be terminal")
	}
	if n, ok := next("closed", "delete"); !ok || n != "deleted" {
		t.Fatal("closed feedback must be deletable")
	}
	if !valid(command{FeedbackID: 1, Action: "reply", Reply: "已安排处理", OperatorID: 2, IdempotencyKey: "feedback-reply-1"}) {
		t.Fatal("valid reply rejected")
	}
	if !valid(command{FeedbackID: 1, Action: "delete", OperatorID: 2, IdempotencyKey: "feedback-delete-1"}) {
		t.Fatal("valid delete rejected")
	}
	if !valid(command{Action: "category_create", Name: "功能建议", Sort: 10, Status: 1, OperatorID: 2, IdempotencyKey: "feedback-category-create-1"}) {
		t.Fatal("valid category creation rejected")
	}
	if valid(command{Action: "category_update", CategoryID: 3, Name: "分类", Sort: -1, Status: 1, OperatorID: 2, IdempotencyKey: "feedback-category-update-3"}) {
		t.Fatal("negative category sort accepted")
	}
}
