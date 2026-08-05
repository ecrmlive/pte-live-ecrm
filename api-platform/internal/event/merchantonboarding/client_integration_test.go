package merchantonboarding

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestProvisionIntegrationNATSRequestReplyContract(t *testing.T) {
	natsURL := os.Getenv("ECRM_ONBOARDING_NATS_URL")
	if natsURL == "" {
		t.Skip("set ECRM_ONBOARDING_NATS_URL to run onboarding client request/reply integration test")
	}
	responder, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_onboarding_client_acceptance"))
	if err != nil {
		t.Fatalf("connect local NATS responder: %v", err)
	}
	t.Cleanup(responder.Close)
	received := make(chan Request, 1)
	if _, err := responder.Subscribe(Subject, func(msg *nats.Msg) {
		var input Request
		if err := json.Unmarshal(msg.Data, &input); err != nil {
			t.Errorf("decode platform onboarding request: %v", err)
			return
		}
		received <- input
		body, _ := json.Marshal(Result{MerchantID: 987680097, StoreID: 987680097, Account: input.Account})
		_ = msg.Respond(body)
	}); err != nil {
		t.Fatalf("subscribe onboarding responder: %v", err)
	}
	if err := responder.Flush(); err != nil {
		t.Fatalf("flush onboarding responder: %v", err)
	}

	client, err := New(natsURL)
	if err != nil {
		t.Fatalf("create platform onboarding client: %v", err)
	}
	t.Cleanup(client.Close)
	out, err := client.Provision(context.Background(), Request{
		ApplicationID: 987680097,
		RegionID:      987680097,
		MerchantName:  "七禧平台请求验收书店",
		ContactName:   "赵敏",
		ContactMobile: "13800000097",
		Account:       "平台验收店主987680097",
		PasswordHash:  "$2b$12$local-onboarding-contract-hash-only",
	})
	if err != nil {
		t.Fatalf("provision through NATS: %v", err)
	}
	if out.MerchantID != 987680097 || out.StoreID != 987680097 || out.Account != "平台验收店主987680097" {
		t.Fatalf("unexpected onboarding result: %+v", out)
	}
	select {
	case in := <-received:
		if in.MerchantName != "七禧平台请求验收书店" || in.ContactName != "赵敏" || in.ContactMobile != "13800000097" || in.RegionID != 987680097 {
			t.Fatalf("platform onboarding request lost Chinese fields: %+v", in)
		}
	case <-time.After(time.Second):
		t.Fatal("NATS responder did not receive platform onboarding request")
	}
}
