// Package paymentconfig transfers the platform-owned payment configuration to
// the business runtime as an encrypted read model. The business API never
// queries qixi_crm_admin directly and neither admin nor merchant APIs return
// the encrypted payload to a browser.
package paymentconfig

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const runtimeKey = "payment.runtime.v1"

var ErrNotConfigured = errors.New("支付平台配置尚未完成")

type Values map[string]string

type runtimeRow struct {
	ConfigKey   string `gorm:"column:config_key;primaryKey"`
	ConfigValue []byte `gorm:"column:config_value"`
}

func (runtimeRow) TableName() string { return "qixi_crm_b_config" }

type Store struct {
	db   *gorm.DB
	aead cipher.AEAD
}

func NewStore(db *gorm.DB, masterSecret string) (*Store, error) {
	if db == nil || strings.TrimSpace(masterSecret) == "" {
		return nil, ErrNotConfigured
	}
	key := sha256.Sum256([]byte("qixi-mergers/payment-runtime/v1:" + masterSecret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Store{db: db, aead: aead}, nil
}

// Publish only runs in the trusted platform process after a payment setting
// save. It writes an encrypted business read model, never a plaintext secret.
func (s *Store) Publish(ctx context.Context, values map[string]string) error {
	return s.publish(ctx, runtimeKey, values)
}

// PublishStore is a one-way encrypted projection of a store-owned payment
// account. It is intentionally separate from the platform key above.
func (s *Store) PublishStore(ctx context.Context, storeID uint, values map[string]string) error {
	if storeID == 0 {
		return ErrNotConfigured
	}
	return s.publish(ctx, storeKey(storeID), values)
}

func (s *Store) publish(ctx context.Context, key string, values map[string]string) error {
	plain, err := json.Marshal(values)
	if err != nil {
		return err
	}
	ciphertext, err := s.encrypt(plain)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(map[string]string{"ciphertext": ciphertext})
	if err != nil {
		return err
	}
	row := runtimeRow{ConfigKey: key, ConfigValue: payload}
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "config_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"config_value", "updated_at"}),
	}).Create(&row).Error
}

func (s *Store) Load(ctx context.Context) (Values, error) {
	return s.load(ctx, runtimeKey)
}

func (s *Store) LoadStore(ctx context.Context, storeID uint) (Values, error) {
	if storeID == 0 {
		return nil, ErrNotConfigured
	}
	return s.load(ctx, storeKey(storeID))
}

func (s *Store) load(ctx context.Context, key string) (Values, error) {
	var row runtimeRow
	if err := s.db.WithContext(ctx).Where("config_key = ?", key).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotConfigured
		}
		return nil, err
	}
	var payload struct {
		Ciphertext string `json:"ciphertext"`
	}
	if err := json.Unmarshal(row.ConfigValue, &payload); err != nil || payload.Ciphertext == "" {
		return nil, ErrNotConfigured
	}
	plain, err := s.decrypt(payload.Ciphertext)
	if err != nil {
		return nil, ErrNotConfigured
	}
	var values Values
	if err := json.Unmarshal(plain, &values); err != nil {
		return nil, ErrNotConfigured
	}
	return values, nil
}

func storeKey(storeID uint) string {
	return "payment.store." + strconv.FormatUint(uint64(storeID), 10) + ".v1"
}

func (s *Store) encrypt(plain []byte) (string, error) {
	nonce := make([]byte, s.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(append(nonce, s.aead.Seal(nil, nonce, plain, nil)...)), nil
}

// Seal/Open are used by the merchant service to persist its own credential
// bundle in qixi_crm_merchant before publishing the encrypted read model.
func (s *Store) Seal(values map[string]string) (string, error) {
	plain, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return s.encrypt(plain)
}

func (s *Store) Open(ciphertext string) (Values, error) {
	plain, err := s.decrypt(ciphertext)
	if err != nil {
		return nil, ErrNotConfigured
	}
	var values Values
	if err := json.Unmarshal(plain, &values); err != nil {
		return nil, ErrNotConfigured
	}
	return values, nil
}

func (s *Store) decrypt(value string) ([]byte, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil || len(raw) < s.aead.NonceSize() {
		return nil, ErrNotConfigured
	}
	return s.aead.Open(nil, raw[:s.aead.NonceSize()], raw[s.aead.NonceSize():], nil)
}

func ChannelReady(values Values, channel string) bool {
	prefix := strings.ToLower(strings.TrimSpace(channel)) + "_"
	if values[prefix+"enabled"] != "true" {
		return false
	}
	switch strings.TrimSuffix(prefix, "_") {
	case "wechat":
		return values["wechat_app_id"] != "" && values["wechat_mch_id"] != "" && values["wechat_api_v3_key"] != "" && values["wechat_private_key"] != "" && values["wechat_notify_url"] != ""
	case "alipay":
		return values["alipay_app_id"] != "" && values["alipay_private_key"] != "" && values["alipay_public_key"] != "" && values["alipay_notify_url"] != ""
	default:
		return false
	}
}

// StoreChannelReady validates a merchant-owned credential bundle. Its field
// names intentionally have no platform prefix, preventing accidental fallback
// from a merchant payment account to the platform account.
func StoreChannelReady(values Values, channel string) bool {
	if values["enabled"] != "true" {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(channel)) {
	case "wechat":
		return values["app_id"] != "" && values["mch_id"] != "" && values["api_v3_key"] != "" && values["serial_no"] != "" && values["private_key"] != "" && values["notify_url"] != ""
	case "alipay":
		return values["app_id"] != "" && values["private_key"] != "" && values["public_key"] != "" && values["notify_url"] != ""
	default:
		return false
	}
}
