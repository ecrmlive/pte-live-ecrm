package notification

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

type memoryStore struct{ row Config }

func (s *memoryStore) List(_ context.Context, _ Audience, _, _ int) ([]Config, int64, error) {
	return []Config{s.row}, 1, nil
}
func (s *memoryStore) Get(_ context.Context, id uint) (*Config, error) {
	if id != s.row.NotificationID {
		return nil, gorm.ErrRecordNotFound
	}
	copy := s.row
	return &copy, nil
}
func (s *memoryStore) Save(_ context.Context, config *Config) error { s.row = *config; return nil }

func TestSavePersistsAllChannelDefaults(t *testing.T) {
	store := &memoryStore{row: Config{NotificationID: 1, Audience: AudienceMember}}
	got, err := NewService(store).Save(context.Background(), 1, SaveInput{
		WechatEnabled: 1, MiniEnabled: 1, SMSEnabled: 0,
		WechatText: "订单已发货", MiniText: "订单已发货", SMSText: "",
	})
	if err != nil {
		t.Fatalf("save notification: %v", err)
	}
	if got.WechatEnabled != 1 || got.MiniEnabled != 1 || got.SMSEnabled != 0 || store.row.MiniText != "订单已发货" {
		t.Fatalf("unexpected saved configuration: %#v", store.row)
	}
}

func TestSaveRejectsEnabledChannelWithoutText(t *testing.T) {
	store := &memoryStore{row: Config{NotificationID: 1, Audience: AudienceMember}}
	_, err := NewService(store).Save(context.Background(), 1, SaveInput{WechatEnabled: 1})
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("error = %v, want ErrBadParam", err)
	}
}

func TestSyncNeverPretendsSuccessWithoutCredentials(t *testing.T) {
	err := NewService(&memoryStore{}).Sync(context.Background(), AudienceStore, ChannelWechat)
	if !errors.Is(err, ErrSyncUnavailable) {
		t.Fatalf("error = %v, want ErrSyncUnavailable", err)
	}
}
