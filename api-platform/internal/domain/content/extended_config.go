package content

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/gorm"
)

const (
	StorageConfigKey          = "storage_config"
	UserSetupConfigKey        = "user_setup_config"
	TransferSettingsConfigKey = "transfer_settings_config"
	RoutineAppConfigKey       = "routine_app_config"
	WechatReplyConfigKey      = "wechat_reply_config"
	WechatMenusConfigKey      = "wechat_menus_config"
	WechatTemplateConfigKey   = "wechat_template_config"
	WechatNewsConfigKey       = "wechat_news_config"
)

var appStubConfigKeys = map[string]struct{}{
	RoutineAppConfigKey:     {},
	WechatReplyConfigKey:    {},
	WechatMenusConfigKey:    {},
	WechatTemplateConfigKey: {},
	WechatNewsConfigKey:     {},
}

type appStubConfig struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Remark  string `json:"remark"`
}

type storageConfig struct {
	Provider   string `json:"provider"`
	Region     string `json:"region"`
	BucketName string `json:"bucket_name"`
	Enabled    bool   `json:"enabled"`
	Remark     string `json:"remark"`
}

type userSetupConfig struct {
	RegisterEnabled bool   `json:"register_enabled"`
	MobileRequired  bool   `json:"mobile_required"`
	InviteRequired  bool   `json:"invite_required"`
	Remark          string `json:"remark"`
}

type transferSettingsConfig struct {
	Enabled   bool    `json:"enabled"`
	MinAmount float64 `json:"min_amount"`
	Remark    string  `json:"remark"`
}

func defaultAppStubConfig(name, remark string) appStubConfig {
	return appStubConfig{Name: name, Enabled: false, Remark: remark}
}

func defaultStorageConfig() storageConfig {
	return storageConfig{
		Provider:   "cos",
		Region:     "",
		BucketName: "",
		Enabled:    false,
		Remark:     "不含 SecretId/SecretKey；真实对象存储凭据请通过云服务配置维护",
	}
}

func defaultUserSetupConfig() userSetupConfig {
	return userSetupConfig{
		RegisterEnabled: true,
		MobileRequired:  true,
		InviteRequired:  false,
		Remark:          "仅保存注册开关与校验规则；不含短信或第三方登录密钥",
	}
}

func defaultTransferSettingsConfig() transferSettingsConfig {
	return transferSettingsConfig{
		Enabled:   false,
		MinAmount: 1,
		Remark:    "仅保存转账监管开关与最低金额；真实打款凭据不在后台保存",
	}
}

func appStubDefault(key string) appStubConfig {
	switch key {
	case RoutineAppConfigKey:
		return defaultAppStubConfig("小程序", "不含 AppId/AppSecret；真实凭据请通过云服务配置维护")
	case WechatReplyConfigKey:
		return defaultAppStubConfig("关键词回复", "不含 Token/EncodingAESKey；真实凭据请通过云服务配置维护")
	case WechatMenusConfigKey:
		return defaultAppStubConfig("自定义菜单", "仅保存菜单配置开关；菜单 JSON 不含密钥")
	case WechatTemplateConfigKey:
		return defaultAppStubConfig("模板消息", "不含模板 ID 密钥；真实凭据请通过云服务配置维护")
	case WechatNewsConfigKey:
		return defaultAppStubConfig("图文消息", "不含素材密钥；真实凭据请通过云服务配置维护")
	default:
		return defaultAppStubConfig("应用配置", "")
	}
}

func parseAppStubConfig(raw string, fallback appStubConfig) (appStubConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return appStubConfig{}, ErrBadParam
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return appStubConfig{}, ErrBadParam
		}
		switch key {
		case "name", "enabled", "remark":
		default:
			return appStubConfig{}, ErrBadParam
		}
	}
	var config appStubConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil ||
		len([]rune(strings.TrimSpace(config.Name))) > 64 ||
		len([]rune(strings.TrimSpace(config.Remark))) > 500 {
		return appStubConfig{}, ErrBadParam
	}
	config.Name = strings.TrimSpace(config.Name)
	config.Remark = strings.TrimSpace(config.Remark)
	if config.Name == "" {
		config.Name = fallback.Name
	}
	if config.Remark == "" {
		config.Remark = fallback.Remark
	}
	return config, nil
}

func (s *Service) GetAppStubConfig(ctx context.Context, key string) (string, error) {
	if _, ok := appStubConfigKeys[key]; !ok {
		return "", ErrBadParam
	}
	fallback := appStubDefault(key)
	row, err := s.store.GetCache(ctx, key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			data, _ := json.Marshal(fallback)
			return string(data), nil
		}
		return "", err
	}
	config, err := parseAppStubConfig(row.Result, fallback)
	if err != nil {
		data, _ := json.Marshal(fallback)
		return string(data), nil
	}
	data, _ := json.Marshal(config)
	return string(data), nil
}

func (s *Service) SaveAppStubConfig(ctx context.Context, key, raw string) (string, error) {
	if _, ok := appStubConfigKeys[key]; !ok {
		return "", ErrBadParam
	}
	config, err := parseAppStubConfig(raw, appStubDefault(key))
	if err != nil {
		return "", err
	}
	data, _ := json.Marshal(config)
	if err := s.store.UpsertCache(ctx, &Cache{Key: key, ExpireTime: 0, Result: string(data)}); err != nil {
		return "", err
	}
	return string(data), nil
}

func parseStorageConfig(raw string) (storageConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return storageConfig{}, ErrBadParam
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return storageConfig{}, ErrBadParam
		}
		switch key {
		case "provider", "region", "bucket_name", "enabled", "remark":
		default:
			return storageConfig{}, ErrBadParam
		}
	}
	var config storageConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return storageConfig{}, ErrBadParam
	}
	config.Provider = strings.TrimSpace(config.Provider)
	config.Region = strings.TrimSpace(config.Region)
	config.BucketName = strings.TrimSpace(config.BucketName)
	config.Remark = strings.TrimSpace(config.Remark)
	if config.Provider == "" {
		config.Provider = defaultStorageConfig().Provider
	}
	if config.Remark == "" {
		config.Remark = defaultStorageConfig().Remark
	}
	if len([]rune(config.Provider)) > 32 || len([]rune(config.Region)) > 64 ||
		len([]rune(config.BucketName)) > 128 || len([]rune(config.Remark)) > 500 {
		return storageConfig{}, ErrBadParam
	}
	return config, nil
}

func (s *Service) GetStorageConfig(ctx context.Context) (string, error) {
	return getJSONConfig(s, ctx, StorageConfigKey, defaultStorageConfig(), parseStorageConfig)
}

func (s *Service) SaveStorageConfig(ctx context.Context, raw string) (string, error) {
	return saveJSONConfig(s, ctx, StorageConfigKey, raw, parseStorageConfig)
}

func parseUserSetupConfig(raw string) (userSetupConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return userSetupConfig{}, ErrBadParam
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return userSetupConfig{}, ErrBadParam
		}
		switch key {
		case "register_enabled", "mobile_required", "invite_required", "remark":
		default:
			return userSetupConfig{}, ErrBadParam
		}
	}
	var config userSetupConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil ||
		len([]rune(strings.TrimSpace(config.Remark))) > 500 {
		return userSetupConfig{}, ErrBadParam
	}
	config.Remark = strings.TrimSpace(config.Remark)
	if config.Remark == "" {
		config.Remark = defaultUserSetupConfig().Remark
	}
	return config, nil
}

func (s *Service) GetUserSetupConfig(ctx context.Context) (string, error) {
	return getJSONConfig(s, ctx, UserSetupConfigKey, defaultUserSetupConfig(), parseUserSetupConfig)
}

func (s *Service) SaveUserSetupConfig(ctx context.Context, raw string) (string, error) {
	return saveJSONConfig(s, ctx, UserSetupConfigKey, raw, parseUserSetupConfig)
}

func parseTransferSettingsConfig(raw string) (transferSettingsConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return transferSettingsConfig{}, ErrBadParam
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return transferSettingsConfig{}, ErrBadParam
		}
		switch key {
		case "enabled", "min_amount", "remark":
		default:
			return transferSettingsConfig{}, ErrBadParam
		}
	}
	var config transferSettingsConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil ||
		len([]rune(strings.TrimSpace(config.Remark))) > 500 || config.MinAmount < 0 {
		return transferSettingsConfig{}, ErrBadParam
	}
	config.Remark = strings.TrimSpace(config.Remark)
	if config.Remark == "" {
		config.Remark = defaultTransferSettingsConfig().Remark
	}
	return config, nil
}

func (s *Service) GetTransferSettingsConfig(ctx context.Context) (string, error) {
	return getJSONConfig(s, ctx, TransferSettingsConfigKey, defaultTransferSettingsConfig(), parseTransferSettingsConfig)
}

func (s *Service) SaveTransferSettingsConfig(ctx context.Context, raw string) (string, error) {
	return saveJSONConfig(s, ctx, TransferSettingsConfigKey, raw, parseTransferSettingsConfig)
}

func (s *Service) ClearMaintainCache(ctx context.Context) error {
	_ = ctx
	return nil
}

type jsonConfigParser[T any] func(string) (T, error)

func getJSONConfig[T any](s *Service, ctx context.Context, key string, fallback T, parse jsonConfigParser[T]) (string, error) {
	row, err := s.store.GetCache(ctx, key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			data, _ := json.Marshal(fallback)
			return string(data), nil
		}
		return "", err
	}
	config, err := parse(row.Result)
	if err != nil {
		data, _ := json.Marshal(fallback)
		return string(data), nil
	}
	data, _ := json.Marshal(config)
	return string(data), nil
}

func saveJSONConfig[T any](s *Service, ctx context.Context, key, raw string, parse jsonConfigParser[T]) (string, error) {
	config, err := parse(raw)
	if err != nil {
		return "", err
	}
	data, _ := json.Marshal(config)
	if err := s.store.UpsertCache(ctx, &Cache{Key: key, ExpireTime: 0, Result: string(data)}); err != nil {
		return "", err
	}
	return string(data), nil
}
