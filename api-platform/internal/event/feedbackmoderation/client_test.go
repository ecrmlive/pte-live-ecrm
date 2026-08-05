package feedbackmoderation

import "testing"

func TestFeedbackCommandValidation(t *testing.T) {
	if !Valid(Command{FeedbackID: 1, Action: "close", OperatorID: 2, IdempotencyKey: "feedback-close-1"}) {
		t.Fatal("valid close rejected")
	}
	if !Valid(Command{FeedbackID: 1, Action: "delete", OperatorID: 2, IdempotencyKey: "feedback-delete-1"}) {
		t.Fatal("valid delete rejected")
	}
	if !Valid(Command{Action: "category_create", Name: "订单问题", Sort: 20, Status: 1, OperatorID: 2, IdempotencyKey: "feedback-category-create-2"}) {
		t.Fatal("valid category creation rejected")
	}
	if Valid(Command{Action: "category_status", CategoryID: 2, Status: 2, OperatorID: 2, IdempotencyKey: "feedback-category-status-2"}) {
		t.Fatal("invalid category status accepted")
	}
	if Valid(Command{FeedbackID: 1, Action: "reply", OperatorID: 2, IdempotencyKey: "feedback-reply-1"}) {
		t.Fatal("empty reply accepted")
	}
}
