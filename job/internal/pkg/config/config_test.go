package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadJWTSecretFromAppYAML(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "app.yaml")
	if err := os.WriteFile(configPath, []byte("jwt:\n  secret: shared-jwt-for-test-only\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.JWT.Secret != "shared-jwt-for-test-only" {
		t.Fatalf("JWT secret = %q, want app YAML value", cfg.JWT.Secret)
	}
}

func TestLoadAllowsNoJWTSecret(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "app.yaml")
	if err := os.WriteFile(configPath, []byte("server:\n  addr: ':8080'\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.JWT.Secret != "" {
		t.Fatalf("JWT secret = %q, want empty for job", cfg.JWT.Secret)
	}
}

func TestLoadRejectsLocalIMMode(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "app.yaml")
	data := []byte("jwt:\n  secret: test-only-secret\nim:\n  mode: local\n")
	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := Load(configPath); err == nil {
		t.Fatal("Load() error = nil, want local IM mode error")
	}
}

func TestDSNForRequiresMatchingScope(t *testing.T) {
	cfg := Config{Databases: DatabasesConfig{
		Business: MySQLConfig{DSN: "business-dsn"},
	}}
	got, err := cfg.DSNFor(DatabaseBusiness)
	if err != nil || got != "business-dsn" {
		t.Fatalf("DSNFor(business) = %q, %v", got, err)
	}
	if _, err := cfg.DSNFor(DatabaseAdmin); err == nil {
		t.Fatal("DSNFor(admin) error = nil, want missing admin DSN error")
	}
}
