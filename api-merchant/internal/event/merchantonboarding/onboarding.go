// Package merchantonboarding owns the only write path for platform-approved
// merchant onboarding. It is idempotent by the deterministic store app_id.
package merchantonboarding

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const Subject = "qixi.platform.merchant-onboarding.v1"

type request struct {
	ApplicationID uint   `json:"application_id"`
	RegionID      uint   `json:"region_id"`
	MerchantName  string `json:"merchant_name"`
	ContactName   string `json:"contact_name"`
	ContactMobile string `json:"contact_mobile"`
	Account       string `json:"account"`
	PasswordHash  string `json:"password_hash"`
}

type result struct {
	MerchantID uint   `json:"merchant_id"`
	StoreID    uint   `json:"store_id"`
	Account    string `json:"account"`
	Error      string `json:"error,omitempty"`
}

type storeRow struct {
	ID         uint `gorm:"column:id"`
	MerchantID uint `gorm:"column:merchant_id"`
}

type merchantRow struct {
	ID       uint   `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	RegionID uint   `gorm:"column:region_id"`
	Status   int    `gorm:"column:status"`
}

func (merchantRow) TableName() string { return "qixi_crm_m_merchant" }

type provisionStoreRow struct {
	ID         uint   `gorm:"column:id"`
	MerchantID uint   `gorm:"column:merchant_id"`
	AppID      string `gorm:"column:app_id"`
	Name       string `gorm:"column:name"`
	Status     int    `gorm:"column:status"`
}

func (provisionStoreRow) TableName() string { return "qixi_crm_m_store" }

func Start(ctx context.Context, merchantDB *gorm.DB, natsURL string) (*nats.Conn, error) {
	if strings.TrimSpace(natsURL) == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_merchant_onboarding"))
	if err != nil {
		return nil, err
	}
	_, err = nc.QueueSubscribe(Subject, "pte_live_ecrm_merchant_onboarding", func(msg *nats.Msg) {
		out, applyErr := Apply(ctx, merchantDB, msg.Data)
		if applyErr != nil {
			out.Error = "开通店铺失败"
		}
		wire, _ := json.Marshal(out)
		_ = msg.Respond(wire)
	})
	if err != nil {
		nc.Close()
		return nil, err
	}
	return nc, nc.Flush()
}

func Apply(ctx context.Context, db *gorm.DB, raw []byte) (result, error) {
	var in request
	if err := json.Unmarshal(raw, &in); err != nil {
		return result{}, err
	}
	in.MerchantName = strings.TrimSpace(in.MerchantName)
	in.ContactName = strings.TrimSpace(in.ContactName)
	in.ContactMobile = strings.TrimSpace(in.ContactMobile)
	in.Account = strings.TrimSpace(in.Account)
	if in.ApplicationID == 0 || in.RegionID == 0 || in.MerchantName == "" || in.ContactName == "" || in.ContactMobile == "" || in.Account == "" || !strings.HasPrefix(in.PasswordHash, "$2") {
		return result{}, errors.New("onboarding command is incomplete")
	}
	appID := fmt.Sprintf("qixi.store.onboard.%d", in.ApplicationID)
	out := result{Account: in.Account}
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing storeRow
		err := tx.Table("qixi_crm_m_store").Where("app_id = ?", appID).Take(&existing).Error
		if err == nil {
			out.MerchantID, out.StoreID = existing.MerchantID, existing.ID
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		merchant := merchantRow{Name: in.MerchantName, RegionID: in.RegionID, Status: 1}
		if err := tx.Create(&merchant).Error; err != nil {
			return err
		}
		if merchant.ID == 0 {
			return errors.New("merchant id was not returned")
		}
		store := provisionStoreRow{MerchantID: merchant.ID, AppID: appID, Name: in.MerchantName, Status: 1}
		if err := tx.Create(&store).Error; err != nil {
			return err
		}
		if store.ID == 0 {
			return errors.New("store id was not returned")
		}
		if err := tx.Table("qixi_crm_m_account").Create(map[string]any{"store_id": store.ID, "username": in.Account, "password_hash": in.PasswordHash, "role_code": "owner", "display_name": in.ContactName, "phone": in.ContactMobile, "status": 1}).Error; err != nil {
			return err
		}
		out.MerchantID, out.StoreID = merchant.ID, store.ID
		return nil
	})
	// Two deliveries of the same command may both miss the initial lookup. The
	// unique app_id makes only one transaction create the store; the loser must
	// return that completed provisioning result instead of surfacing a retryable
	// duplicate-key error to the platform reviewer.
	if isDuplicateKey(err) {
		var existing storeRow
		lookupErr := db.WithContext(ctx).Table("qixi_crm_m_store").Where("app_id = ?", appID).Take(&existing).Error
		if lookupErr == nil {
			out.MerchantID, out.StoreID = existing.MerchantID, existing.ID
			return out, nil
		}
	}
	return out, err
}

func isDuplicateKey(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
