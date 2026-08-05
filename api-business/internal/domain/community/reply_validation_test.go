package community

import (
	"context"
	"strings"
	"testing"
)

func TestCreateReplyRejectsBlankAndOverlongContent(t *testing.T) {
	svc := NewService(nil)
	for _, content := range []string{"   ", strings.Repeat("中", 1001)} {
		if _, err := svc.CreateReply(context.Background(), 1, 1, CreateReplyInput{Content: content}); err != ErrBadParam {
			t.Fatalf("content length %d error = %v, want %v", len([]rune(content)), err, ErrBadParam)
		}
	}
}

func TestCreatePostRejectsOverlongUnicodeTitleOrContent(t *testing.T) {
	svc := NewService(nil)
	cases := []CreatePostInput{
		{Title: strings.Repeat("中", 101), Content: "中文正文"},
		{Title: "中文标题", Content: strings.Repeat("中", 5001)},
	}
	for _, input := range cases {
		if _, err := svc.CreatePost(context.Background(), 1, input); err != ErrBadParam {
			t.Fatalf("post title/content length must be rejected, error = %v", err)
		}
	}
	if !validPostText(strings.Repeat("中", 100), strings.Repeat("中", 5000)) {
		t.Fatal("Unicode boundaries must be accepted")
	}
}
