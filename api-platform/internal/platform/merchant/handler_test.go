package merchant

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/merchant"
)

func TestUintScopeDropsZeroAndKeepsChineseFixtureIDs(t *testing.T) {
	got := uints([]uint64{0, 10, 20})
	if !reflect.DeepEqual(got, []uint{10, 20}) {
		t.Fatalf("区域主管王小明的范围 = %#v", got)
	}
}

func TestMerchantCategoryCommissionRateFailsClosed(t *testing.T) {
	svc := merchant.NewService(nil)
	if _, err := svc.CreateCategory(t.Context(), "中文演示分类", 100.01); !errors.Is(err, merchant.ErrBadParam) {
		t.Fatalf("out-of-range commission error = %v, want ErrBadParam", err)
	}
	if _, err := svc.CreateCategory(t.Context(), "中文演示分类", 8.123); !errors.Is(err, merchant.ErrBadParam) {
		t.Fatalf("fractional-cent commission error = %v, want ErrBadParam", err)
	}
	if err := svc.UpdateCategory(t.Context(), 0, "中文演示分类", 8.5); !errors.Is(err, merchant.ErrBadParam) {
		t.Fatalf("zero category id error = %v, want ErrBadParam", err)
	}
}

func TestMerchantOperateActionLabel(t *testing.T) {
	cases := map[string]string{
		"PUT /api/platform/v1/merchants/:id/status":    "开启/关闭",
		"PUT /api/platform/v1/merchants/:id/recommend": "推荐变更",
		"POST /api/platform/v1/merchants":              "新增",
		"PUT /api/platform/v1/merchants/:id":           "编辑",
		"DELETE /api/platform/v1/merchants/:id":        "删除",
		"": "操作",
	}
	for in, want := range cases {
		if got := merchantOperateActionLabel(in); got != want {
			t.Fatalf("action %q label = %q, want %q", in, got, want)
		}
	}
}

func TestMerchantOperateTerminalAndRoleLabel(t *testing.T) {
	if got := merchantOperateTerminal("platform"); got != "平台操作" {
		t.Fatalf("platform terminal = %q", got)
	}
	if got := merchantOperateTerminal("merchant"); got != "商户操作" {
		t.Fatalf("merchant terminal = %q", got)
	}
	if got := merchantOperateRoleLabel("platform,region"); got != "平台管理员/区域管理员" {
		t.Fatalf("role label = %q", got)
	}
}

func TestParseOperateDateRange(t *testing.T) {
	start, end, ok := parseOperateDateRange("2026-01-05", "2026-01-05")
	if !ok || start.IsZero() || end.IsZero() {
		t.Fatalf("same-day range should be valid")
	}
	if !end.Equal(start.AddDate(0, 0, 1)) {
		t.Fatalf("end should be exclusive next day, got %v", end)
	}
	if _, _, ok := parseOperateDateRange("2026-01-06", "2026-01-05"); ok {
		t.Fatalf("inverted range must fail")
	}
	if _, _, ok := parseOperateDateRange("bad", ""); ok {
		t.Fatalf("bad start must fail")
	}
	emptyStart, emptyEnd, ok := parseOperateDateRange("", "")
	if !ok || !emptyStart.IsZero() || !emptyEnd.IsZero() {
		t.Fatalf("empty range should be zero times")
	}
	_ = time.Time{}
}
