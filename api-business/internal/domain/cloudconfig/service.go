package cloudconfig

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

var ErrBadGroup = errors.New("不支持的配置分组")
var ErrBadField = errors.New("不支持的配置字段")
var ErrBadValue = errors.New("配置字段值不合法")

const (
	keyVersionEncrypted       = "v1"
	keyVersionBootstrapPublic = "bootstrap-public-v1"
	keyVersionBootstrapLocal  = "bootstrap-local-v1"
)

type Store interface {
	ListByGroup(ctx context.Context, group string) ([]Config, error)
	Upsert(ctx context.Context, row *Config) error
}

type Service struct {
	store Store
	aead  cipher.AEAD
}

func NewService(store Store, masterSecret string) (*Service, error) {
	key := sha256.Sum256([]byte("qixi-cloud-config/v1:" + masterSecret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Service{store: store, aead: aead}, nil
}

// Catalog 是平台可管理的服务端配置白名单，避免任意键注入。
func Catalog() []GroupMeta {
	return []GroupMeta{
		{Key: "payment", Label: "微信支付 / 支付宝", Fields: []FieldMeta{
			{Key: "wechat_enabled", Label: "启用微信支付"}, {Key: "alipay_enabled", Label: "启用支付宝支付"},
			{Key: "wechat_app_id", Label: "微信 AppID"}, {Key: "wechat_mch_id", Label: "微信商户号"},
			{Key: "wechat_api_v3_key", Label: "微信 APIv3 密钥", Secret: true}, {Key: "wechat_serial_no", Label: "微信商户证书序列号"},
			{Key: "wechat_private_key", Label: "微信商户私钥", Secret: true}, {Key: "wechat_merchant_cert", Label: "微信商户证书", Secret: true},
			{Key: "wechat_public_key_id", Label: "微信支付公钥 ID"}, {Key: "wechat_public_key", Label: "微信支付公钥", Secret: true},
			{Key: "wechat_platform_cert_serial", Label: "微信平台证书序列号"}, {Key: "wechat_platform_cert", Label: "微信平台证书", Secret: true},
			{Key: "wechat_notify_url", Label: "微信支付回调地址"},
			{Key: "alipay_app_id", Label: "支付宝 AppID"}, {Key: "alipay_private_key", Label: "支付宝应用私钥", Secret: true},
			{Key: "alipay_public_key", Label: "支付宝公钥"}, {Key: "alipay_seller_id", Label: "支付宝卖家 ID"},
			{Key: "alipay_notify_url", Label: "支付宝回调地址"}, {Key: "alipay_gateway", Label: "支付宝网关地址"}, {Key: "alipay_sign_type", Label: "支付宝签名算法", Hint: "默认 RSA2"},
		}},
		{Key: "wechat_mini_program", Label: "微信小程序", Fields: []FieldMeta{
			{Key: "enabled", Label: "启用微信小程序"}, {Key: "app_id", Label: "小程序 AppID", Required: true},
			{Key: "app_secret", Label: "小程序 AppSecret", Secret: true, Required: true},
		}},
		{Key: "tencent_account", Label: "腾讯云账号", Fields: []FieldMeta{
			{Key: "secret_id", Label: "SecretId", Secret: true, Required: true}, {Key: "secret_key", Label: "SecretKey", Secret: true, Required: true},
			{Key: "region", Label: "默认地域", Hint: "例如 ap-guangzhou"},
		}},
		{Key: "cos", Label: "腾讯云 COS", Fields: []FieldMeta{
			{Key: "enabled", Label: "启用 COS"}, {Key: "bucket", Label: "存储桶"}, {Key: "region", Label: "地域"},
			{Key: "base_url", Label: "CDN / 自定义域名"}, {Key: "key_prefix", Label: "对象键前缀"},
		}},
		{Key: "lvb", Label: "腾讯云 LVB 直播", Fields: []FieldMeta{
			{Key: "push_domain", Label: "推流域名"}, {Key: "play_domain", Label: "拉流域名"}, {Key: "push_key", Label: "推流鉴权 Key", Secret: true},
			{Key: "pull_auth_key", Label: "拉流鉴权 Key", Secret: true}, {Key: "app_name", Label: "应用名称", Hint: "默认 live"},
			{Key: "push_url_template", Label: "推流 URL 模板"}, {Key: "play_url_template", Label: "拉流 URL 模板"},
		}},
		{Key: "vod", Label: "腾讯云 VOD", Fields: []FieldMeta{
			{Key: "app_id", Label: "VOD AppID"}, {Key: "sub_app_id", Label: "子应用 ID"}, {Key: "procedure", Label: "任务流名称"},
			{Key: "callback_url", Label: "事件通知地址"}, {Key: "signature", Label: "上传签名 / 密钥", Secret: true},
		}},
		{Key: "live_license", Label: "腾讯直播 / 播放器 License", Fields: []FieldMeta{
			{Key: "push_license_name", Label: "推流 License 名称"}, {Key: "push_license_package", Label: "推流包名"},
			{Key: "push_license_url", Label: "推流 License URL"}, {Key: "push_license_key", Label: "推流 License Key", Secret: true},
			{Key: "player_license_name", Label: "播放器 License 名称"}, {Key: "player_license_package", Label: "播放器包名"},
			{Key: "player_license_domain", Label: "播放器 License 域名"}, {Key: "player_license_url", Label: "播放器 License URL"},
			{Key: "player_license_key", Label: "播放器 License Key", Secret: true},
		}},
		{Key: "lvb_callback", Label: "腾讯 LVB 推拉流回调", Fields: []FieldMeta{
			{Key: "template_name", Label: "模板名称"}, {Key: "template_id", Label: "模板 ID"}, {Key: "callback_secret", Label: "回调密钥", Secret: true},
			{Key: "callback_types", Label: "回调类型", Hint: "publish,disconnect,exception"}, {Key: "publish_callback_url", Label: "推流回调地址"},
			{Key: "disconnect_callback_url", Label: "断流回调地址"}, {Key: "exception_callback_url", Label: "推流异常回调地址"},
			{Key: "bound_domain", Label: "模板绑定域名"},
		}},
		{Key: "im", Label: "即时通讯服务", Fields: []FieldMeta{
			{Key: "mode", Label: "接入模式", Hint: "仅 remote"}, {Key: "api_base", Label: "IM 服务端 API 地址", Hint: "仅容器/服务端 S2S 使用"},
			{Key: "api_public_url", Label: "IM 客户端 API 公网地址", Hint: "H5/小程序可访问，禁止填写 Docker 服务名"},
			{Key: "ws_public_url", Label: "WebSocket 公网地址"}, {Key: "app_id", Label: "IM AppID"},
			{Key: "integration_token", Label: "IM 服务令牌", Secret: true},
		}},
	}
}

// Values 仅供服务端运行时适配器读取完整配置，绝不能直接作为 HTTP 响应返回。
func (s *Service) Values(ctx context.Context, group string) (map[string]string, error) {
	meta, ok := groupMeta(group)
	if !ok {
		return nil, ErrBadGroup
	}
	rows, err := s.store.ListByGroup(ctx, meta.Key)
	if err != nil {
		return nil, err
	}
	byKey := make(map[string]Config, len(rows))
	for _, row := range rows {
		byKey[row.ConfigKey] = row
	}
	values := make(map[string]string, len(meta.Fields))
	for _, field := range meta.Fields {
		if row, exists := byKey[field.Key]; exists {
			plain, err := s.decode(row, field)
			if err != nil {
				return nil, err
			}
			values[field.Key] = plain
		}
	}
	return values, nil
}

func (s *Service) List(ctx context.Context) ([]GroupView, error) {
	metas := Catalog()
	out := make([]GroupView, 0, len(metas))
	for _, meta := range metas {
		view, err := s.Get(ctx, meta.Key)
		if err != nil {
			return nil, err
		}
		out = append(out, *view)
	}
	return out, nil
}

func (s *Service) Get(ctx context.Context, group string) (*GroupView, error) {
	meta, ok := groupMeta(group)
	if !ok {
		return nil, ErrBadGroup
	}
	rows, err := s.store.ListByGroup(ctx, meta.Key)
	if err != nil {
		return nil, err
	}
	byKey := make(map[string]Config, len(rows))
	for _, row := range rows {
		byKey[row.ConfigKey] = row
	}
	values := make(map[string]string, len(meta.Fields))
	configured := false
	var latest *time.Time
	for _, field := range meta.Fields {
		row, exists := byKey[field.Key]
		if !exists {
			values[field.Key] = ""
			continue
		}
		if latest == nil || row.UpdateTime.After(*latest) {
			at := row.UpdateTime
			latest = &at
		}
		if field.Secret {
			values[field.Key] = SecretMasked
		} else {
			plain, err := s.decode(row, field)
			if err != nil {
				return nil, err
			}
			values[field.Key] = plain
		}
	}
	for _, field := range meta.Fields {
		if field.Required && values[field.Key] == "" {
			return &GroupView{GroupKey: meta.Key, Label: meta.Label, Fields: meta.Fields, Values: values, UpdatedAt: latest, Configured: false}, nil
		}
	}
	for _, field := range meta.Fields {
		if values[field.Key] != "" {
			configured = true
			break
		}
	}
	return &GroupView{GroupKey: meta.Key, Label: meta.Label, Fields: meta.Fields, Values: values, UpdatedAt: latest, Configured: configured}, nil
}

func (s *Service) Save(ctx context.Context, group string, in SaveInput, adminID uint) (*GroupView, error) {
	meta, ok := groupMeta(group)
	if !ok {
		return nil, ErrBadGroup
	}
	if len(in.Values) == 0 {
		return nil, ErrBadValue
	}
	allowed := make(map[string]FieldMeta, len(meta.Fields))
	for _, field := range meta.Fields {
		allowed[field.Key] = field
	}
	keys := make([]string, 0, len(in.Values))
	for key := range in.Values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		field, ok := allowed[key]
		if !ok {
			return nil, ErrBadField
		}
		value := strings.TrimSpace(in.Values[key])
		if field.Secret && (value == "" || value == SecretMasked) {
			continue
		}
		if len(value) > 16*1024 {
			return nil, ErrBadValue
		}
		ciphertext, err := s.encrypt(value)
		if err != nil {
			return nil, err
		}
		row := &Config{GroupKey: meta.Key, ConfigKey: key, Ciphertext: ciphertext, KeyVersion: keyVersionEncrypted, IsSecret: field.Secret, UpdatedBy: adminID}
		if err := s.store.Upsert(ctx, row); err != nil {
			return nil, err
		}
	}
	return s.Get(ctx, meta.Key)
}

func groupMeta(key string) (GroupMeta, bool) {
	for _, meta := range Catalog() {
		if meta.Key == strings.TrimSpace(key) {
			return meta, true
		}
	}
	return GroupMeta{}, false
}
func (s *Service) encrypt(plain string) (string, error) {
	nonce := make([]byte, s.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(append(nonce, s.aead.Seal(nil, nonce, []byte(plain), nil)...)), nil
}

func (s *Service) decode(row Config, field FieldMeta) (string, error) {
	if row.KeyVersion == keyVersionBootstrapLocal {
		return row.Ciphertext, nil
	}
	if row.KeyVersion == keyVersionBootstrapPublic {
		if field.Secret {
			return "", errors.New("密钥不得使用可提交的 key.sql 明文初始化")
		}
		return row.Ciphertext, nil
	}
	return s.decrypt(row.Ciphertext)
}
func (s *Service) decrypt(raw string) (string, error) {
	data, err := base64.RawStdEncoding.DecodeString(raw)
	if err != nil {
		return "", fmt.Errorf("decode config ciphertext: %w", err)
	}
	if len(data) < s.aead.NonceSize() {
		return "", errors.New("invalid config ciphertext")
	}
	plain, err := s.aead.Open(nil, data[:s.aead.NonceSize()], data[s.aead.NonceSize():], nil)
	if err != nil {
		return "", fmt.Errorf("decrypt config: %w", err)
	}
	return string(plain), nil
}
