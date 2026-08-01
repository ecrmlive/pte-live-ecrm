// Package merchantapplication projects C-end application events into the
// platform-owned review queue. The only cross-service exchange is NATS.
package merchantapplication

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const Subject = "qixi.business.merchant-application.v1"
const Submitted = "business.merchant_application.submitted"

type envelope struct {
	EventID   uint64          `json:"event_id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}
type payload struct {
	ID              uint64 `json:"id"`
	ApplicantUserID uint64 `json:"applicant_user_id"`
	MerchantName    string `json:"merchant_name"`
	ContactName     string `json:"contact_name"`
	ContactMobile   string `json:"contact_mobile"`
	CategoryName    string `json:"category_name"`
	MerchantType    string `json:"merchant_type"`
	LicenseKey      string `json:"license_key"`
}

func Start(ctx context.Context, adminDB *gorm.DB, natsURL string) (*nats.Conn, error) {
	if natsURL == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_platform_merchant_application_projection"))
	if err != nil {
		return nil, err
	}
	if _, err = nc.Subscribe(Subject, func(msg *nats.Msg) {
		if err := Apply(ctx, adminDB, msg.Data); err != nil {
			log.Printf("merchant application projection failed: %v", err)
		}
	}); err != nil {
		nc.Close()
		return nil, err
	}
	return nc, nc.Flush()
}

func Apply(ctx context.Context, db *gorm.DB, raw []byte) error {
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return err
	}
	if env.EventType != Submitted {
		return nil
	}
	var data payload
	if err := json.Unmarshal(env.Payload, &data); err != nil {
		return err
	}
	if data.ID == 0 || data.ApplicantUserID == 0 || data.MerchantName == "" {
		return nil
	}
	return db.WithContext(ctx).Exec(`INSERT INTO qixi_crm_a_merchant_application
    (source_application_id, applicant_user_id, merchant_name, contact_name, contact_mobile, category_name, merchant_type, license_key, status, created_at)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'pending', NOW())
    ON DUPLICATE KEY UPDATE merchant_name=VALUES(merchant_name), contact_name=VALUES(contact_name), contact_mobile=VALUES(contact_mobile), category_name=VALUES(category_name), merchant_type=VALUES(merchant_type), license_key=VALUES(license_key)`,
		data.ID, data.ApplicantUserID, data.MerchantName, data.ContactName, data.ContactMobile, data.CategoryName, data.MerchantType, data.LicenseKey).Error
}
