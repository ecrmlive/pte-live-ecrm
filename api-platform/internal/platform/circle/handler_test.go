package circle

import "testing"

func TestValidAgentPasswordReset(t *testing.T) {
	valid := &resetAgentPasswordInput{
		Password:       "LocalDemoPassword12",
		Reason:         "虚构中文验收：代理本人已完成身份核验",
		IdempotencyKey: "agent-password-demo-001",
	}
	if !validAgentPasswordReset(valid) {
		t.Fatal("expected valid password reset command")
	}
	for _, input := range []*resetAgentPasswordInput{
		{Password: "short", Reason: valid.Reason, IdempotencyKey: valid.IdempotencyKey},
		{Password: valid.Password, Reason: "x", IdempotencyKey: valid.IdempotencyKey},
		{Password: valid.Password, Reason: valid.Reason, IdempotencyKey: "short"},
		nil,
	} {
		if validAgentPasswordReset(input) {
			t.Fatalf("invalid password reset command accepted: %#v", input)
		}
	}
}
