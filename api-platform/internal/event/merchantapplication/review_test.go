package merchantapplication

import "testing"

func TestEnqueueReviewIgnoresNonTerminalOrUnlinkedDecision(t *testing.T) {
	cases := []ReviewPayload{
		{SourceApplicationID: 0, Status: "approved"},
		{SourceApplicationID: 18, Status: "pending"},
		{SourceApplicationID: 18, Status: ""},
	}
	for _, payload := range cases {
		if err := EnqueueReview(nil, payload); err != nil {
			t.Fatalf("payload should be ignored without database access: %#v, err=%v", payload, err)
		}
	}
}
