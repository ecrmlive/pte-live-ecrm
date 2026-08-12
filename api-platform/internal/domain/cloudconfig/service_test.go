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

func TestPublicBootstrapConfigCanBeReadButCannotSeedSecret(t *testing.T) {
	store := newMemoryStore()
	store.rows["cos:enabled"] = Config{GroupKey: "cos", ConfigKey: "enabled", Ciphertext: "true", KeyVersion: keyVersionBootstrapPublic}
	store.rows["cos:bucket"] = Config{GroupKey: "cos", ConfigKey: "bucket", Ciphertext: "demo-bucket", KeyVersion: keyVersionBootstrapPublic}
	store.rows["tencent_account:secret_id"] = Config{GroupKey: "tencent_account", ConfigKey: "secret_id", Ciphertext: "forbidden", KeyVersion: keyVersionBootstrapPublic}
	svc, err := NewService(store, "test-master-key")
	if err != nil {
		t.Fatal(err)
	}
	values, err := svc.Values(context.Background(), "cos")
	if err != nil || values["bucket"] != "demo-bucket" || values["enabled"] != "true" {
		t.Fatalf("values=%v err=%v", values, err)
	}
	if _, err := svc.Values(context.Background(), "tencent_account"); err == nil {
		t.Fatal("expected bootstrap secret rejection")
	}
}

func TestIgnoredKeySQLBootstrapCanSeedSecret(t *testing.T) {
	store := newMemoryStore()
	store.rows["tencent_account:secret_id"] = Config{GroupKey: "tencent_account", ConfigKey: "secret_id", Ciphertext: "local-secret", KeyVersion: keyVersionBootstrapLocal}
	store.rows["tencent_account:region"] = Config{GroupKey: "tencent_account", ConfigKey: "region", Ciphertext: "ap-shanghai", KeyVersion: keyVersionBootstrapLocal}
	svc, err := NewService(store, "test-master-key")
	if err != nil {
		t.Fatal(err)
	}
	values, err := svc.Values(context.Background(), "tencent_account")
	if err != nil || values["secret_id"] != "local-secret" || values["region"] != "ap-shanghai" {
		t.Fatalf("values=%v err=%v", values, err)
	}
}

func TestTencentSMSConfigEncryptsAppKeyAndKeepsPublicValues(t *testing.T) {
	store := newMemoryStore()
	svc, err := NewService(store, "test-master-key")
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.Save(context.Background(), "tencent_sms", SaveInput{Values: map[string]string{
		"enabled":      "true",
		"sdk_app_id":   "1401165606",
		"app_key":      "local-app-key",
		"sign_id":      "711884",
		"sign_content": "杭州乐成体育",
		"template_id":  "2701987",
	}}, 9)
	if err != nil {
		t.Fatal(err)
	}
	if store.rows["tencent_sms:app_key"].Ciphertext == "local-app-key" {
		t.Fatal("App Key must be encrypted")
	}
	view, err := svc.Get(context.Background(), "tencent_sms")
	if err != nil || view.Values["app_key"] != SecretMasked || view.Values["sign_content"] != "杭州乐成体育" {
		t.Fatalf("view=%#v err=%v", view, err)
	}
}

func TestTencentSMSRequiresAppKeyOnFirstSaveButKeepsConfiguredSecret(t *testing.T) {
	store := newMemoryStore()
	svc, err := NewService(store, "test-master-key")
	if err != nil {
		t.Fatal(err)
	}
	values := map[string]string{
		"enabled":      "true",
		"sdk_app_id":   "1401165606",
		"sign_id":      "711884",
		"sign_content": "杭州乐成体育",
		"template_id":  "2701987",
	}
	if _, err = svc.Save(context.Background(), "tencent_sms", SaveInput{Values: values}, 9); err != ErrBadValue {
		t.Fatalf("first save without app key err=%v, want %v", err, ErrBadValue)
	}
	values["app_key"] = "local-app-key"
	if _, err = svc.Save(context.Background(), "tencent_sms", SaveInput{Values: values}, 9); err != nil {
		t.Fatal(err)
	}
	values["app_key"] = SecretMasked
	if _, err = svc.Save(context.Background(), "tencent_sms", SaveInput{Values: values}, 9); err != nil {
		t.Fatalf("masked configured app key should be preserved: %v", err)
	}
}
