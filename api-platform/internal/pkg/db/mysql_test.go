package db

import (
	"strings"
	"testing"
	"time"
)

func TestNormalizeShanghaiDSN(t *testing.T) {
	got, err := NormalizeShanghaiDSN("qixi_crm:secret@tcp(pte_live_mysql:3306)/qixi_crm_admin?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Fatalf("NormalizeShanghaiDSN: %v", err)
	}
	if !strings.Contains(got, "loc=Asia%2FShanghai") && !strings.Contains(got, "loc=Asia/Shanghai") {
		t.Fatalf("expected Asia/Shanghai loc, got %q", got)
	}
	if !strings.Contains(got, "parseTime=true") && !strings.Contains(got, "parseTime=True") {
		t.Fatalf("expected parseTime, got %q", got)
	}
	if !strings.Contains(got, "time_zone") {
		t.Fatalf("expected time_zone param, got %q", got)
	}
}

func TestEnsureShanghaiTimezone(t *testing.T) {
	EnsureShanghaiTimezone()
	now := time.Now()
	name, offset := now.Zone()
	if offset != 8*3600 {
		t.Fatalf("expected +08:00 offset, got name=%s offset=%d", name, offset)
	}
}
