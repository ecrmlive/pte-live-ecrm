package merchantapplication

import "testing"

func TestApplyReviewRejectsMalformedAndNonTerminalMessagesBeforeDatabaseAccess(t *testing.T) {
	if err := ApplyReview(nil, nil, []byte("not-json")); err == nil {
		t.Fatal("malformed event must be rejected")
	}
	cases := []string{
		`{"event_id":1,"event_type":"ignored","payload":{}}`,
		`{"event_id":1,"event_type":"platform.merchant_application.reviewed","payload":{"source_application_id":0,"status":"approved"}}`,
		`{"event_id":1,"event_type":"platform.merchant_application.reviewed","payload":{"source_application_id":11,"status":"pending"}}`,
	}
	for _, raw := range cases {
		if err := ApplyReview(nil, nil, []byte(raw)); err != nil {
			t.Fatalf("event should be ignored before database access: %s, err=%v", raw, err)
		}
	}
}
