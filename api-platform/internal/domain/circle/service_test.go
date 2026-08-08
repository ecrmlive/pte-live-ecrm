package circle

import (
	"context"
	"testing"
	"time"
)

type serviceStoreStub struct {
	agent         *Agent
	updatedAgent  *Agent
	revokedKey    string
	revokedReason string
}

func (s *serviceStoreStub) ListCircles(context.Context, CircleListFilter, int, int) ([]Circle, int64, error) {
	return nil, 0, nil
}
func (s *serviceStoreStub) GetCircle(context.Context, uint) (*Circle, error) { return nil, nil }
func (s *serviceStoreStub) CreateCircle(context.Context, *Circle) error      { return nil }
func (s *serviceStoreStub) UpdateCircle(context.Context, *Circle) error      { return nil }
func (s *serviceStoreStub) DeleteCircle(context.Context, uint) error         { return nil }
func (s *serviceStoreStub) CountCircleChildren(context.Context, uint) (int64, error) {
	return 0, nil
}
func (s *serviceStoreStub) ListAgents(context.Context, string, *int8, *int8, int, int) ([]Agent, int64, error) {
	if s.agent == nil {
		return nil, 0, nil
	}
	return []Agent{*s.agent}, 1, nil
}
func (s *serviceStoreStub) GetAgent(_ context.Context, id uint) (*Agent, error) {
	if s.agent == nil || s.agent.CircleAgentID != id {
		return nil, ErrNotFound
	}
	copy := *s.agent
	return &copy, nil
}
func (s *serviceStoreStub) CreateAgent(context.Context, *Agent) error { return nil }
func (s *serviceStoreStub) UpdateAgent(_ context.Context, row *Agent) error {
	copy := *row
	s.updatedAgent = &copy
	s.agent = &copy
	return nil
}
func (s *serviceStoreStub) AuditAgent(context.Context, uint, int8, string, uint, time.Time) error {
	return nil
}
func (s *serviceStoreStub) RevokeAgent(_ context.Context, id uint, reason, key string, _ uint, _ time.Time) (bool, error) {
	if s.agent == nil || s.agent.CircleAgentID != id {
		return false, ErrNotFound
	}
	if s.revokedKey == key {
		if s.revokedReason != reason {
			return false, ErrCommandConflict
		}
		return true, nil
	}
	s.revokedKey, s.revokedReason = key, reason
	s.agent.Status = AgentRevoked
	return false, nil
}

func TestListAgentsDoesNotReturnPaymentCredentials(t *testing.T) {
	store := &serviceStoreStub{agent: &Agent{
		CircleAgentID:  1,
		PaymentAccount: "local-demo-account",
		PaymentBank:    "local-demo-bank",
		PaymentQRImg:   "/demo/payment-qr.png",
	}}
	res, err := NewService(store).ListAgents(context.Background(), "", nil, nil, 1, 20)
	if err != nil {
		t.Fatalf("list agents: %v", err)
	}
	if len(res.List) != 1 || !res.List[0].PaymentConfigured {
		t.Fatalf("payment configuration should be indicated without exposing values: %+v", res.List)
	}
	if res.List[0].PaymentAccount != "" || res.List[0].PaymentBank != "" || res.List[0].PaymentQRImg != "" {
		t.Fatal("list agents must not return payment account, bank, or QR data")
	}
}

func TestUpdateAgentKeepsPaymentCredentialsWhenInputIsBlank(t *testing.T) {
	store := &serviceStoreStub{agent: &Agent{
		CircleAgentID:  1,
		Name:           "虚构区域代理",
		Phone:          "13900000010",
		Status:         AgentPending,
		PaymentAccount: "local-demo-account",
		PaymentBank:    "local-demo-bank",
		PaymentQRImg:   "/demo/payment-qr.png",
	}}
	_, err := NewService(store).UpdateAgent(context.Background(), 1, AgentInput{
		Name: "虚构区域代理（更新）", Phone: "13900000010", PaymentMethod: 1,
	})
	if err != nil {
		t.Fatalf("update agent: %v", err)
	}
	if store.updatedAgent == nil {
		t.Fatal("expected persistent update")
	}
	if store.updatedAgent.PaymentAccount != "local-demo-account" || store.updatedAgent.PaymentBank != "local-demo-bank" || store.updatedAgent.PaymentQRImg != "/demo/payment-qr.png" {
		t.Fatalf("blank update must retain write-only payment data: %+v", store.updatedAgent)
	}
}

func TestRevokeAgentRequiresReasonAndReplaysSameCommand(t *testing.T) {
	store := &serviceStoreStub{agent: &Agent{CircleAgentID: 11, Status: AgentApproved}}
	svc := NewService(store)
	if _, err := svc.RevokeAgent(context.Background(), 11, "x", "revoke-001", 9); err != ErrBadParam {
		t.Fatalf("short revoke reason error=%v", err)
	}
	replayed, err := svc.RevokeAgent(context.Background(), 11, "虚构中文撤销原因", "revoke-001", 9)
	if err != nil || replayed || store.agent.Status != AgentRevoked {
		t.Fatalf("first revoke replayed=%v status=%d err=%v", replayed, store.agent.Status, err)
	}
	replayed, err = svc.RevokeAgent(context.Background(), 11, "虚构中文撤销原因", "revoke-001", 9)
	if err != nil || !replayed {
		t.Fatalf("revoke retry replayed=%v err=%v", replayed, err)
	}
}
