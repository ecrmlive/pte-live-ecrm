package chat

import (
	"context"
	"testing"

	"gorm.io/gorm"
)

type chatStoreFake struct {
	threads  map[uint]*Thread
	messages []Message
	nextID   uint
	config   *MerchantIMConfig
}

func newChatStoreFake() *chatStoreFake { return &chatStoreFake{threads: map[uint]*Thread{}, nextID: 1} }
func (s *chatStoreFake) FindThreadByMerUID(_ context.Context, merID, uid uint) (*Thread, error) {
	for _, row := range s.threads {
		if row.MerID == merID && row.UID == uid {
			copy := *row
			return &copy, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (s *chatStoreFake) GetThread(_ context.Context, id uint) (*Thread, error) {
	row, ok := s.threads[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *row
	return &copy, nil
}
func (s *chatStoreFake) CreateThread(_ context.Context, row *Thread) error {
	row.ThreadID = s.nextID
	s.nextID++
	copy := *row
	s.threads[row.ThreadID] = &copy
	return nil
}
func (s *chatStoreFake) UpdateThread(_ context.Context, row *Thread) error {
	copy := *row
	s.threads[row.ThreadID] = &copy
	return nil
}
func (s *chatStoreFake) ListThreadsByMer(_ context.Context, merID uint, _, _ int) ([]Thread, int64, error) {
	var rows []Thread
	for _, row := range s.threads {
		if row.MerID == merID {
			rows = append(rows, *row)
		}
	}
	return rows, int64(len(rows)), nil
}
func (s *chatStoreFake) ListThreadsByUID(_ context.Context, uid uint, _, _ int) ([]Thread, int64, error) {
	var rows []Thread
	for _, row := range s.threads {
		if row.UID == uid {
			rows = append(rows, *row)
		}
	}
	return rows, int64(len(rows)), nil
}
func (s *chatStoreFake) CreateMessage(_ context.Context, row *Message) error {
	row.MsgID = uint64(len(s.messages) + 1)
	s.messages = append(s.messages, *row)
	return nil
}
func (s *chatStoreFake) ListMessages(_ context.Context, threadID uint, _, _ int) ([]Message, int64, error) {
	var rows []Message
	for _, row := range s.messages {
		if row.ThreadID == threadID {
			rows = append(rows, row)
		}
	}
	return rows, int64(len(rows)), nil
}
func (s *chatStoreFake) EnsureIdentity(context.Context, string, uint, string, int64) (*Identity, error) {
	return &Identity{}, nil
}
func (s *chatStoreFake) FindCustomerServiceID(context.Context, uint) (uint, error) { return 7001, nil }
func (s *chatStoreFake) LoadUserNickname(context.Context, uint) (string, error) {
	return "七禧体验用户", nil
}
func (s *chatStoreFake) LoadMerName(context.Context, uint) (string, error) {
	return "七禧服饰旗舰店", nil
}
func (s *chatStoreFake) LoadMerchantIMConfig(context.Context, uint) (*MerchantIMConfig, error) {
	return s.config, nil
}

func TestOpenThreadUsesBusinessSessionStateAndReopens(t *testing.T) {
	store := newChatStoreFake()
	svc := NewService(store, IMSettings{})
	opened, err := svc.OpenUserThread(context.Background(), 9101, 1)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if opened.Status != StatusOpen || opened.ServiceID != 7001 || opened.ImConversationID == "" {
		t.Fatalf("unexpected thread: %#v", opened)
	}
	opened.Status = StatusClosed
	if err := store.UpdateThread(context.Background(), opened); err != nil {
		t.Fatal(err)
	}
	reopened, err := svc.OpenUserThread(context.Background(), 9101, 1)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if reopened.Status != StatusOpen || reopened.LastMsg != "会话已重新开启" {
		t.Fatalf("thread not reopened: %#v", reopened)
	}
	if len(store.messages) != 1 || store.messages[0].Content != "您好，客服稍后为您服务" {
		t.Fatalf("system metadata not persisted: %#v", store.messages)
	}
}

func TestOrderCardMetadataTracksUnreadWithoutPersistingText(t *testing.T) {
	store := newChatStoreFake()
	svc := NewService(store, IMSettings{})
	thread, err := svc.OpenUserThread(context.Background(), 9101, 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SendUser(context.Background(), 9101, thread.ThreadID, SendInput{MsgType: MsgText, Content: "你好"}); err != ErrTextViaCS {
		t.Fatalf("text err=%v", err)
	}
	message, err := svc.SendUser(context.Background(), 9101, thread.ThreadID, SendInput{MsgType: MsgOrder, Content: `{"order_id":88001,"title":"七禧演示订单"}`})
	if err != nil {
		t.Fatalf("order card: %v", err)
	}
	if message.MsgType != MsgOrder {
		t.Fatalf("message=%#v", message)
	}
	updated, _ := store.GetThread(context.Background(), thread.ThreadID)
	if updated.ServiceUnread != 1 || updated.LastMsg == "" {
		t.Fatalf("unread not updated: %#v", updated)
	}
}

func TestMerchantIMConfigOverridesOnlyClientRouting(t *testing.T) {
	store := newChatStoreFake()
	store.config = &MerchantIMConfig{SDKAppID: "30101", APIPublicURL: "https://im.example.invalid/api", WSPublicURL: "wss://im.example.invalid/ws"}
	svc := NewService(store, IMSettings{Mode: "remote", APIBase: "http://im", AppID: "30001", Token: "server-only", APIPublicURL: "https://default.invalid", WSPublicURL: "wss://default.invalid"})
	im, err := svc.effectiveIMSettingsForMerchant(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if im.AppID != "30101" || im.APIPublicURL != "https://im.example.invalid/api" || im.WSPublicURL != "wss://im.example.invalid/ws" {
		t.Fatalf("merchant im config was not applied: %#v", im)
	}
	if im.Token != "server-only" {
		t.Fatal("server-only integration token must not be replaced by merchant client routing")
	}
}
