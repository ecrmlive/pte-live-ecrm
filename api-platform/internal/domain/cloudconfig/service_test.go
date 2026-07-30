package cloudconfig

import (
	"context"
	"testing"
)

type memoryStore struct{ rows map[string]Config }

func newMemoryStore() *memoryStore { return &memoryStore{rows: map[string]Config{}} }
func (m *memoryStore) ListByGroup(_ context.Context, group string) ([]Config, error) {
	out := []Config{}
	for _, row := range m.rows {
		if row.GroupKey == group {
			out = append(out, row)
		}
	}
	return out, nil
}
func (m *memoryStore) Upsert(_ context.Context, row *Config) error {
	copy := *row
	m.rows[row.GroupKey+":"+row.ConfigKey] = copy
	return nil
}

func TestSecretIsEncryptedAndMasked(t *testing.T) {
	store := newMemoryStore()
	svc, err := NewService(store, "test-master-key")
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.Save(context.Background(), "tencent_account", SaveInput{Values: map[string]string{"secret_id": "id-secret", "secret_key": "key-secret", "region": "ap-guangzhou"}}, 9)
	if err != nil {
		t.Fatal(err)
	}
	if store.rows["tencent_account:secret_key"].Ciphertext == "key-secret" {
		t.Fatal("secret stored as plaintext")
	}
	view, err := svc.Get(context.Background(), "tencent_account")
	if err != nil {
		t.Fatal(err)
	}
	if got := view.Values["secret_key"]; got != SecretMasked {
		t.Fatalf("secret view=%q", got)
	}
	if got := view.Values["region"]; got != "ap-guangzhou" {
		t.Fatalf("region=%q", got)
	}
	_, err = svc.Save(context.Background(), "tencent_account", SaveInput{Values: map[string]string{"secret_key": SecretMasked}}, 9)
	if err != nil {
		t.Fatal(err)
	}
	plain, err := svc.decrypt(store.rows["tencent_account:secret_key"].Ciphertext)
	if err != nil || plain != "key-secret" {
		t.Fatalf("secret should be preserved")
	}
}
