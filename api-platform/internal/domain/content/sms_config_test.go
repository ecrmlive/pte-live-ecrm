package content

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

type smsStore struct{ cache *Cache }

func (s *smsStore) ListNotices(context.Context, bool, int, int) ([]Notice, int64, error) {
	return nil, 0, nil
}
func (s *smsStore) GetNotice(context.Context, uint) (*Notice, error) {
	return nil, gorm.ErrRecordNotFound
}
func (s *smsStore) CreateNotice(context.Context, *Notice) error  { return nil }
func (s *smsStore) UpdateNotice(context.Context, *Notice) error  { return nil }
func (s *smsStore) SoftDeleteNotice(context.Context, uint) error { return nil }
func (s *smsStore) GetCache(_ context.Context, key string) (*Cache, error) {
	if s.cache == nil || s.cache.Key != key {
		return nil, gorm.ErrRecordNotFound
	}
	return s.cache, nil
}
func (s *smsStore) UpsertCache(_ context.Context, row *Cache) error { s.cache = row; return nil }

func TestSMSStubRejectsSecretBearingConfig(t *testing.T) {
	svc := NewService(&smsStore{})
	_, err := svc.SaveSMSConfig(context.Background(), `{"enabled":true,"provider":"stub","api_key":"not-allowed"}`)
	if !errors.Is(err, ErrBadParam) {
		t.Fatalf("secret-bearing config error = %v, want ErrBadParam", err)
	}
}

func TestSMSStubCanonicalizesAndNeverEchoesLegacySecret(t *testing.T) {
	store := &smsStore{}
	svc := NewService(store)
	got, err := svc.SaveSMSConfig(context.Background(), `{"remark":"中文模拟短信开关","provider":"stub","enabled":true,"sign":"七禧商城"}`)
	if err != nil {
		t.Fatalf("save safe stub: %v", err)
	}
	want := `{"enabled":true,"provider":"stub","sign":"七禧商城","remark":"中文模拟短信开关"}`
	if got != want || store.cache == nil || store.cache.Result != want {
		t.Fatalf("canonical config = %q, stored = %#v", got, store.cache)
	}
	store.cache.Result = `{"provider":"real","secret":"must-not-leak"}`
	safe, err := svc.GetSMSConfig(context.Background())
	if err != nil {
		t.Fatalf("get safe stub: %v", err)
	}
	if safe == store.cache.Result || safe != marshalSMSConfig(defaultSMSConfig()) {
		t.Fatalf("legacy secret must not be returned, got %q", safe)
	}
}
