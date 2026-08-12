package databasebackup

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestValidTableName(t *testing.T) {
	tests := []struct {
		scope databaseScope
		name  string
		want  bool
	}{
		{scopeAdmin, "qixi_crm_a_admin_user", true},
		{scopeBusiness, "qixi_crm_b_order", true},
		{scopeAdmin, "qixi_crm_b_order", false},
		{scopeAdmin, "qixi_crm_a_bad-name", false},
		{scopeBusiness, "information_schema.tables", false},
	}
	for _, test := range tests {
		if got := validTableName(test.scope, test.name); got != test.want {
			t.Fatalf("validTableName(%q, %q) = %v, want %v", test.scope, test.name, got, test.want)
		}
	}
}

func TestBackupPathIsContained(t *testing.T) {
	base := t.TempDir()
	if !isSafeBackupPath(base, filepath.Join(base, "ecrm_admin.sql")) {
		t.Fatal("expected backup path inside directory to be accepted")
	}
	if isSafeBackupPath(base, filepath.Join(base, "..", "outside.sql")) {
		t.Fatal("expected path outside directory to be rejected")
	}
	if isSafeBackupPath(base, filepath.Join(base, "backup.txt")) {
		t.Fatal("expected non-sql backup path to be rejected")
	}
}

func TestBackupFileName(t *testing.T) {
	name := backupFileName(scopeAdmin, []string{"qixi_crm_a_admin_user", "qixi_crm_a_menu"})
	if !strings.HasPrefix(name, "ecrm_admin_") || !strings.HasSuffix(name, "_02-tables.sql") {
		t.Fatalf("unexpected backup name %q", name)
	}
}
