package profitsharing

import (
	"strings"
	"testing"
)

func TestReviewNoteKeepsChineseAuditWithinColumnBoundary(t *testing.T) {
	if !validReviewNote("虚构中文审核说明：资料齐全，可进入下一步。") {
		t.Fatal("valid Chinese review note must pass")
	}
	if !validReviewNote(strings.Repeat("审", 500)) {
		t.Fatal("500-rune review note must fit the schema column")
	}
	if validReviewNote(strings.Repeat("审", 501)) {
		t.Fatal("oversized review note must be rejected before persistence")
	}
}
