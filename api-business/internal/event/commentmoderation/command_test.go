package commentmoderation

import (
	"errors"
	"testing"
)

func TestCommentModerationStatesAreExplicit(t *testing.T) {
	if next, ok := nextStatus("pending", "publish"); !ok || next != "published" {
		t.Fatal("pending comment must publish")
	}
	if next, ok := nextStatus("published", "hide"); !ok || next != "hidden" {
		t.Fatal("published comment must hide")
	}
	if _, ok := nextStatus("published", "publish"); ok {
		t.Fatal("published comment must not be re-published without a new state")
	}
	if _, ok := nextStatus("hidden", "hide"); ok {
		t.Fatal("hidden comment must not be re-hidden")
	}
	if !valid(command{CommentID: 8801, Action: "hide", OperatorID: 3, IdempotencyKey: "comment-hide-8801"}) {
		t.Fatal("valid comment command must pass")
	}
	if !valid(command{Action: "create_virtual", ProductID: 1001, Score: 5, Content: "虚构中文虚拟评论", VirtualAuthorName: "演示用户小满", Sort: 80, OperatorID: 3, IdempotencyKey: "comment-create-virtual-1001"}) {
		t.Fatal("valid virtual comment command must pass")
	}
	if !valid(command{CommentID: 8803, Action: "update_virtual", Score: 5, Content: "保留原图片的中文编辑", VirtualAuthorName: "演示用户小满", Sort: 81, MediaSet: false, OperatorID: 3, IdempotencyKey: "comment-update-without-media-8803"}) {
		t.Fatal("edit without explicit image replacement must pass")
	}
	if valid(command{CommentID: 8803, Action: "update_virtual", Score: 6, Content: "无效评分", VirtualAuthorName: "演示用户小满", OperatorID: 3, IdempotencyKey: "comment-update-virtual-8803"}) {
		t.Fatal("invalid virtual comment score must fail")
	}
	if !isDuplicate(errors.New("Error 1062: Duplicate entry")) {
		t.Fatal("duplicate-key error must be recognized for idempotent replay")
	}
	if !isUnknownColumn(errors.New("Error 1054: Unknown column 'virtual_author_avatar' in 'field list'")) {
		t.Fatal("missing virtual_author_avatar column must be recognized")
	}
	if isUnknownColumn(errors.New("Error 1062: Duplicate entry")) {
		t.Fatal("duplicate-key must not be treated as schema mismatch")
	}
}
