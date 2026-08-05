package merchantonboarding

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestApplyIntegrationConcurrentChineseOnboardingIsIdempotent(t *testing.T) {
	dsn := os.Getenv("ECRM_ONBOARDING_MERCHANT_TEST_DSN")
	if dsn == "" {
		t.Skip("set ECRM_ONBOARDING_MERCHANT_TEST_DSN to run merchant onboarding integration test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated merchant database: %v", err)
	}

	const applicationID uint = 987680099
	const account = "验收店主987680099"
	const appID = "qixi.store.onboard.987680099"
	ctx := context.Background()
	t.Cleanup(func() {
		_ = db.WithContext(ctx).Table("qixi_crm_m_account").Where("username = ?", account).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_m_store").Where("app_id = ?", appID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_m_merchant").Where("name = ?", "七禧入驻并发验收茶铺").Delete(nil).Error
	})

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("local-onboarding-acceptance-only"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate local test password hash: %v", err)
	}
	wire, err := json.Marshal(request{
		ApplicationID: applicationID,
		RegionID:      987680099,
		MerchantName:  "七禧入驻并发验收茶铺",
		ContactName:   "王小明",
		ContactMobile: "13800000099",
		Account:       account,
		PasswordHash:  string(passwordHash),
	})
	if err != nil {
		t.Fatalf("marshal onboarding command: %v", err)
	}

	const calls = 8
	results := make(chan result, calls)
	errs := make(chan error, calls)
	var wg sync.WaitGroup
	for range calls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, applyErr := Apply(ctx, db, wire)
			if applyErr != nil {
				errs <- applyErr
				return
			}
			results <- out
		}()
	}
	wg.Wait()
	close(results)
	close(errs)
	for applyErr := range errs {
		t.Fatalf("concurrent onboarding command should converge, got %v", applyErr)
	}

	var first result
	for out := range results {
		if out.MerchantID == 0 || out.StoreID == 0 || out.Account != account {
			t.Fatalf("unexpected onboarding response: %+v", out)
		}
		if first.MerchantID == 0 {
			first = out
			continue
		}
		if out.MerchantID != first.MerchantID || out.StoreID != first.StoreID {
			t.Fatalf("duplicate deliveries returned different provisioning results: first=%+v current=%+v", first, out)
		}
	}

	var stores, owners int64
	if err := db.WithContext(ctx).Table("qixi_crm_m_store").Where("app_id = ?", appID).Count(&stores).Error; err != nil {
		t.Fatalf("count provisioned stores: %v", err)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_m_account").Where("username = ? AND role_code = 'owner'", account).Count(&owners).Error; err != nil {
		t.Fatalf("count provisioned owner accounts: %v", err)
	}
	if stores != 1 || owners != 1 {
		t.Fatalf("expected exactly one Chinese merchant/store/owner projection, stores=%d owners=%d", stores, owners)
	}
}

func TestStartIntegrationNATSRequestReplyChineseProvisioning(t *testing.T) {
	dsn := os.Getenv("ECRM_ONBOARDING_MERCHANT_TEST_DSN")
	natsURL := os.Getenv("ECRM_ONBOARDING_NATS_URL")
	if dsn == "" || natsURL == "" {
		t.Skip("set ECRM_ONBOARDING_MERCHANT_TEST_DSN and ECRM_ONBOARDING_NATS_URL to run onboarding request/reply integration test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated merchant database: %v", err)
	}
	const applicationID uint = 987680098
	const account = "验收NATS店主987680098"
	const appID = "qixi.store.onboard.987680098"
	ctx := context.Background()
	t.Cleanup(func() {
		_ = db.WithContext(ctx).Table("qixi_crm_m_account").Where("username = ?", account).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_m_store").Where("app_id = ?", appID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_m_merchant").Where("name = ?", "七禧NATS入驻验收花店").Delete(nil).Error
	})

	subscriber, err := Start(ctx, db, natsURL)
	if err != nil {
		t.Fatalf("start merchant onboarding NATS subscriber: %v", err)
	}
	if subscriber == nil {
		t.Fatal("expected onboarding NATS subscriber")
	}
	t.Cleanup(subscriber.Close)
	publisher, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_onboarding_acceptance"))
	if err != nil {
		t.Fatalf("connect local NATS publisher: %v", err)
	}
	t.Cleanup(publisher.Close)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("local-nats-onboarding-acceptance-only"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate local test password hash: %v", err)
	}
	wire, err := json.Marshal(request{
		ApplicationID: applicationID,
		RegionID:      987680098,
		MerchantName:  "七禧NATS入驻验收花店",
		ContactName:   "李华",
		ContactMobile: "13800000098",
		Account:       account,
		PasswordHash:  string(passwordHash),
	})
	if err != nil {
		t.Fatalf("marshal onboarding command: %v", err)
	}
	request := func() result {
		msg, requestErr := publisher.Request(Subject, wire, 2*time.Second)
		if requestErr != nil {
			t.Fatalf("NATS request/reply onboarding command: %v", requestErr)
		}
		var out result
		if err := json.Unmarshal(msg.Data, &out); err != nil {
			t.Fatalf("decode onboarding response: %v", err)
		}
		if out.Error != "" || out.MerchantID == 0 || out.StoreID == 0 || out.Account != account {
			t.Fatalf("invalid NATS onboarding response: %+v", out)
		}
		return out
	}
	first, repeated := request(), request()
	if repeated.MerchantID != first.MerchantID || repeated.StoreID != first.StoreID {
		t.Fatalf("NATS redelivery returned different provisioning result: first=%+v repeated=%+v", first, repeated)
	}
}
