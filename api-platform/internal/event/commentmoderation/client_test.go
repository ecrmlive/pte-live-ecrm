package commentmoderation

import (
	"context"
	"testing"
)

func TestCommandValidationFailsClosed(t *testing.T) {
	valid := Command{CommentID: 8801, Action: "publish", OperatorID: 9, IdempotencyKey: "comment-publish-8801", Note: "虚构中文审核通过"}
	if !Valid(valid) {
		t.Fatal("valid moderation command must pass")
	}
	if !Valid(Command{Action: "create_virtual", ProductID: 1001, Score: 5, Content: "虚构中文虚拟评论", VirtualAuthorName: "演示用户小满", Sort: 80, OperatorID: 9, IdempotencyKey: "comment-create-virtual-1001"}) {
		t.Fatal("valid virtual-comment creation command must pass")
	}
	if !Valid(Command{Action: "create_virtual", ProductID: 1001, Score: 5, Content: "虚构中文虚拟评论", VirtualAuthorName: "演示用户小满", Sort: 80, Media: []string{"/demo/comment-1.png"}, OperatorID: 9, IdempotencyKey: "comment-create-virtual-media-1001"}) {
		t.Fatal("registered image url command must pass")
	}
	if !Valid(Command{CommentID: 8803, Action: "sort_virtual", Sort: 90, OperatorID: 9, IdempotencyKey: "comment-sort-virtual-8803"}) {
		t.Fatal("valid virtual-comment sort command must pass")
	}
	for _, value := range []Command{{CommentID: 0, Action: "publish", OperatorID: 9, IdempotencyKey: "comment-publish-8801"}, {CommentID: 1, Action: "delete", OperatorID: 9, IdempotencyKey: "comment-delete-1"}, {Action: "create_virtual", ProductID: 1001, Score: 0, Content: "评论", VirtualAuthorName: "小满", OperatorID: 9, IdempotencyKey: "comment-create-virtual-1001"}, {CommentID: 1, Action: "hide", OperatorID: 9, IdempotencyKey: "短"}} {
		if Valid(value) {
			t.Fatalf("invalid command must fail: %#v", value)
		}
	}
	if Valid(Command{Action: "create_virtual", ProductID: 1001, Score: 5, Content: "评论", VirtualAuthorName: "小满", Sort: 1, Media: make([]string, 10), OperatorID: 9, IdempotencyKey: "comment-too-many-media"}) {
		t.Fatal("more than nine images must fail")
	}
}

func TestCommandRejectsUnavailableClient(t *testing.T) {
	if _, err := (*Client)(nil).Dispatch(context.Background(), Command{CommentID: 8801, Action: "hide", OperatorID: 3, IdempotencyKey: "comment-hide-8801"}); err == nil {
		t.Fatal("unavailable command client must fail fast")
	}
}
