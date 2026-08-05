package merchant_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"

	domainmerchant "github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/merchant"
	merchantpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/merchant"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestFinalizeIntentionApprovalIntegrationConcurrentProjection(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_MERCHANT_ONBOARDING_ADMIN_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_MERCHANT_ONBOARDING_ADMIN_TEST_DSN to run platform onboarding projection integration test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated admin database: %v", err)
	}
	const merchantID uint = 987680099
	const regionID uint = 987680099
	const sourceApplicationID uint = 987680099
	const account = "中文入驻店主987680099"
	const merchantName = "七禧平台入驻投影验收茶铺"
	cleanup := func() {
		_ = db.Table("qixi_crm_a_merchant_application").Where("source_application_id = ?", sourceApplicationID).Delete(nil).Error
		_ = db.Table("qixi_crm_a_merchant_view").Where("merchant_id = ?", merchantID).Delete(nil).Error
	}
	cleanup()

	created := db.Table("qixi_crm_a_merchant_application").Create(map[string]any{
		"merchant_name":         merchantName,
		"contact_name":          "王小明",
		"contact_mobile":        "13800000099",
		"region_id":             regionID,
		"source_application_id": sourceApplicationID,
		"category_name":         "中文入驻验收分类",
		"merchant_type":         "企业店",
		"status":                "pending",
	})
	if created.Error != nil {
		t.Fatalf("create pending Chinese application: %v", created.Error)
	}
	var applicationRow struct{ ID uint }
	if err := db.Table("qixi_crm_a_merchant_application").Select("id").Where("source_application_id = ?", sourceApplicationID).Take(&applicationRow).Error; err != nil {
		t.Fatalf("load pending Chinese application id: %v", err)
	}
	applicationID := applicationRow.ID
	t.Cleanup(cleanup)

	service := domainmerchant.NewService(merchantpersist.NewStoreAdapter(merchantpersist.NewRepo(db)))
	in := domainmerchant.AuditIntentionInput{RegionID: regionID, Account: account, Mark: "中文并发入驻审核通过"}
	const calls = 6
	var wg sync.WaitGroup
	results := make(chan error, calls)
	for range calls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, finalizeErr := service.FinalizeIntentionApproval(context.Background(), applicationID, in, merchantID, nil)
			results <- finalizeErr
		}()
	}
	wg.Wait()
	close(results)

	var success, alreadyAudited int
	for finalizeErr := range results {
		switch {
		case finalizeErr == nil:
			success++
		case errors.Is(finalizeErr, domainmerchant.ErrAlreadyAudited):
			alreadyAudited++
		default:
			t.Fatalf("unexpected concurrent finalize error: %v", finalizeErr)
		}
	}
	if success != 1 || alreadyAudited != calls-1 {
		t.Fatalf("expected one finalization and %d idempotency rejections, success=%d already_audited=%d", calls-1, success, alreadyAudited)
	}

	var application struct {
		Status   string
		RegionID uint
	}
	if err := db.Table("qixi_crm_a_merchant_application").Select("status, region_id").Where("id = ?", applicationID).Take(&application).Error; err != nil {
		t.Fatalf("load finalized application: %v", err)
	}
	if application.Status != "approved" || application.RegionID != regionID {
		t.Fatalf("unexpected application projection: %+v", application)
	}
	var view struct {
		MerchantID    uint
		MerchantName  string
		ContactName   string
		ContactMobile string
		RegionID      uint
		Status        int8
	}
	if err := db.Table("qixi_crm_a_merchant_view").Select("merchant_id, merchant_name, contact_name, contact_mobile, region_id, status").Where("merchant_id = ?", merchantID).Take(&view).Error; err != nil {
		t.Fatalf("load merchant supervision projection: %v", err)
	}
	if view.MerchantID != merchantID || view.MerchantName != merchantName || view.ContactName != "王小明" || view.ContactMobile != "13800000099" || view.RegionID != regionID || view.Status != 1 {
		t.Fatalf("unexpected Chinese merchant supervision projection: %+v", view)
	}
}
