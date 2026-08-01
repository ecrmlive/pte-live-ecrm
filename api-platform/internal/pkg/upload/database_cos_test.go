package upload

import (
	"context"
	"errors"
	"strings"
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

func TestDatabaseCOSRejectsDisabledConfiguration(t *testing.T) {
	store := DatabaseCOS{Resolver: resolverStub{"cos": {"enabled": "false"}}}
	_, _, err := store.Save("platform", nil)
	if err == nil || !strings.Contains(err.Error(), "未启用") {
		t.Fatalf("err=%v", err)
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
