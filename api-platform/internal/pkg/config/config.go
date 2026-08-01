package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Databases DatabasesConfig `yaml:"databases"`
	// MySQL 是旧单库字段，仅供兼容读取；本服务必须使用 databases.admin。
	MySQL   MySQLConfig   `yaml:"mysql"`
	Redis   RedisConfig   `yaml:"redis"`
	NATS    NATSConfig    `yaml:"nats"`
	Etcd    EtcdConfig    `yaml:"etcd"`
	JWT     JWTConfig     `yaml:"jwt"`
	Job     JobConfig     `yaml:"job"`
	Upload  UploadConfig  `yaml:"upload"`
	Payment PaymentConfig `yaml:"payment"`
	IM      IMConfig      `yaml:"im"`
}

// DatabaseScope 显式约束目标 API 的数据库所有权，禁止以一个 DSN 跨三个 CRM 库直接查询。
type DatabaseScope string

const (
	DatabaseAdmin    DatabaseScope = "admin"
	DatabaseBusiness DatabaseScope = "business"
	DatabaseMerchant DatabaseScope = "merchant"
)

type DatabasesConfig struct {
	Admin    MySQLConfig `yaml:"admin"`
	Business MySQLConfig `yaml:"business"`
	Merchant MySQLConfig `yaml:"merchant"`
}

// IMConfig 客服 IM 桥。api_base 仅服务端 S2S 使用；api_public_url、ws_public_url 返回给 H5/小程序 SDK。
type IMConfig struct {
	Mode             string `yaml:"mode"` // local | remote
	APIBase          string `yaml:"api_base"`
	APIPublicURL     string `yaml:"api_public_url"`
	WSPublicURL      string `yaml:"ws_public_url"`
	AppID            string `yaml:"app_id"` // 业务 app_id，默认 30001
	IntegrationToken string `yaml:"integration_token"`
}

// PaymentConfig 第三方支付（沙箱验签闭环；真实 SDK 后续替换回调验签）。
type PaymentConfig struct {
	Sandbox      bool   `yaml:"sandbox"`
	NotifySecret string `yaml:"notify_secret"`
	Wechat       bool   `yaml:"wechat"`
	Alipay       bool   `yaml:"alipay"`
}

type UploadConfig struct {
	Dir        string `yaml:"dir"`
	PublicBase string `yaml:"public_base"`
}

type JobConfig struct {
	TickSeconds      int `yaml:"tick_seconds"`
	UnpaidTTLMinutes int `yaml:"unpaid_ttl_minutes"`
	UnpaidBatch      int `yaml:"unpaid_batch"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"`
}

type MySQLConfig struct {
	DSN string `yaml:"dsn"`
}

// DSNFor 返回所属数据库的唯一连接串。调用方必须使用与自身服务边界一致的 scope。
func (c Config) DSNFor(scope DatabaseScope) (string, error) {
	var dsn string
	switch scope {
	case DatabaseAdmin:
		dsn = c.Databases.Admin.DSN
	case DatabaseBusiness:
		dsn = c.Databases.Business.DSN
	case DatabaseMerchant:
		dsn = c.Databases.Merchant.DSN
		// 旧 local 配置仅声明 admin/business 两个 DSN。为保证旧配置升级后
		// 仍能启动，开发态可由 business 库名推导 merchant 库名；正式环境
		// 仍应显式填写 databases.merchant.dsn，避免隐式跨实例推导。
		if dsn == "" && c.Server.Mode == "debug" {
			dsn = strings.Replace(c.Databases.Business.DSN, "qixi_crm_business", "qixi_crm_merchant", 1)
		}
	default:
		return "", fmt.Errorf("unsupported database scope %q", scope)
	}
	if dsn == "" {
		return "", fmt.Errorf("databases.%s.dsn must be configured", scope)
	}
	return dsn, nil
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type NATSConfig struct {
	URL string `yaml:"url"`
}

type EtcdConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

type JWTConfig struct {
	Secret          string `yaml:"secret"`
	AccessTTLHours  int    `yaml:"access_ttl_hours"`
	RefreshTTLHours int    `yaml:"refresh_ttl_hours"`
}

func (j JWTConfig) AccessTTL() time.Duration {
	h := j.AccessTTLHours
	if h <= 0 {
		h = 12
	}
	return time.Duration(h) * time.Hour
}

func (j JWTConfig) RefreshTTL() time.Duration {
	h := j.RefreshTTLHours
	if h <= 0 {
		h = 168
	}
	return time.Duration(h) * time.Hour
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":8080"
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "debug"
	}
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("jwt.secret must be configured in %s", path)
	}
	if cfg.JWT.AccessTTLHours <= 0 {
		cfg.JWT.AccessTTLHours = 12
	}
	if cfg.JWT.RefreshTTLHours <= 0 {
		cfg.JWT.RefreshTTLHours = 168
	}
	if cfg.Job.TickSeconds <= 0 {
		cfg.Job.TickSeconds = 30
	}
	if cfg.Job.UnpaidTTLMinutes <= 0 {
		cfg.Job.UnpaidTTLMinutes = 30
	}
	if cfg.Job.UnpaidBatch <= 0 {
		cfg.Job.UnpaidBatch = 50
	}
	if cfg.Upload.Dir == "" {
		cfg.Upload.Dir = "data/uploads"
	}
	if cfg.Upload.PublicBase == "" {
		cfg.Upload.PublicBase = "/uploads"
	}
	// 未显式配置时：沙箱开启 wechat/alipay，便于本地演示
	if !cfg.Payment.Wechat && !cfg.Payment.Alipay && cfg.Payment.Sandbox {
		cfg.Payment.Wechat = true
		cfg.Payment.Alipay = true
	}
	if cfg.IM.Mode == "local" {
		return nil, fmt.Errorf("im.mode=local is not supported; configure pte-live-im endpoint in %s", path)
	}
	return &cfg, nil
}
