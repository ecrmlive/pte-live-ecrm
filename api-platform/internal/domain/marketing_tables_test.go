package domain_test

import (
	"testing"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/combination"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/presell"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/seckill"
)

func TestMarketingModelsUseBusinessSchema(t *testing.T) {
	cases := map[string]string{
		"seckill activity":     (seckill.Active{}).TableName(),
		"seckill time":         (seckill.TimeSlot{}).TableName(),
		"combination activity": (combination.ProductGroup{}).TableName(),
		"combination buying":   (combination.Buying{}).TableName(),
		"combination member":   (combination.Member{}).TableName(),
		"presell activity":     (presell.ProductPresell{}).TableName(),
		"presell order":        (presell.PresellOrder{}).TableName(),
	}
	for name, table := range cases {
		if len(table) < len("qixi_crm_b_") || table[:len("qixi_crm_b_")] != "qixi_crm_b_" {
			t.Fatalf("%s table = %q, want qixi_crm_b_ prefix", name, table)
		}
	}
}
