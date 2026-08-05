package merchant

import (
	"errors"
	"reflect"
	"testing"

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
