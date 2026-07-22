package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MySQL   MySQLConfig   `yaml:"mysql"`
	Redis   RedisConfig   `yaml:"redis"`
	NATS    NATSConfig    `yaml:"nats"`
	Etcd    EtcdConfig    `yaml:"etcd"`
	Minio   MinioConfig   `yaml:"minio"`
	JWT     JWTConfig     `yaml:"jwt"`
	Job     JobConfig     `yaml:"job"`
	Upload  UploadConfig  `yaml:"upload"`
	Payment PaymentConfig `yaml:"payment"`
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
	TickSeconds       int `yaml:"tick_seconds"`
	UnpaidTTLMinutes  int `yaml:"unpaid_ttl_minutes"`
	UnpaidBatch       int `yaml:"unpaid_batch"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"`
}

type MySQLConfig struct {
	DSN string `yaml:"dsn"`
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

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	UseSSL    bool   `yaml:"use_ssl"`
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
		cfg.JWT.Secret = "qixi-mergers-dev-jwt-secret-change-me"
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
	if cfg.Payment.NotifySecret == "" {
		cfg.Payment.NotifySecret = "qixi-pay-notify-dev-secret-change-me"
	}
	// 未显式配置时：沙箱开启 wechat/alipay，便于本地演示
	if !cfg.Payment.Wechat && !cfg.Payment.Alipay && cfg.Payment.Sandbox {
		cfg.Payment.Wechat = true
		cfg.Payment.Alipay = true
	}
	return &cfg, nil
}
