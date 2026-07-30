package upload

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"
)

type resolverStub map[string]map[string]string

func (s resolverStub) Values(_ context.Context, group string) (map[string]string, error) {
	values, ok := s[group]
	if !ok {
		return nil, errors.New("group not found")
	}
	return values, nil
}

type storeStub struct{ called bool }

func (s *storeStub) Save(_ string, _ *multipart.FileHeader) (string, string, error) {
	s.called = true
	return "/uploads/demo.png", "demo.png", nil
}

func TestDatabaseCOSFallsBackWhenDisabled(t *testing.T) {
	fallback := &storeStub{}
	store := DatabaseCOS{Resolver: resolverStub{
		"cos": {"enabled": "false"},
	}, Fallback: fallback}
	url, name, err := store.Save("platform", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !fallback.called || url != "/uploads/demo.png" || name != "demo.png" {
		t.Fatalf("fallback=%t url=%q name=%q", fallback.called, url, name)
	}
}

func TestDatabaseCOSEnabledRequiresCompleteCredentials(t *testing.T) {
	store := DatabaseCOS{Resolver: resolverStub{
		"cos":             {"enabled": "true", "bucket": "bucket", "region": "ap-guangzhou"},
		"tencent_account": {"secret_id": "id"},
	}}
	_, _, err := store.Save("platform", nil)
	if err == nil {
		t.Fatal("expected incomplete configuration error")
	}
}
