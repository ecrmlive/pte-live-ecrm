package auth

import "testing"

func TestCancellationBlockersRequireConfirmation(t *testing.T) {
	if (CancellationBlockers{}).RequiresConfirmation() {
		t.Fatal("empty blockers should not require confirmation")
	}
	blockers := CancellationBlockers{Balance: 12.5, Points: 20, Commission: 3.2, ActiveOrderCount: 1, OpenRefundCount: 1}
	if !blockers.RequiresConfirmation() {
		t.Fatal("business blockers must require confirmation")
	}
	if got := len(blockers.Messages()); got != 5 {
		t.Fatalf("messages=%d", got)
	}
}

func TestCancellationTokenHashDoesNotExposeToken(t *testing.T) {
	token := "演示确认令牌"
	hash := accountCancellationTokenHash(token)
	if hash == token || len(hash) != 64 {
		t.Fatalf("unexpected hash %q", hash)
	}
	if accountCancellationTokenHash(" "+token+" ") != hash {
		t.Fatal("token hash should normalize surrounding whitespace")
	}
}
