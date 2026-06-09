package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ============================================================================
// Configuration
// ============================================================================

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type EpayConfig struct {
	PartnerID string `yaml:"partner_id"` // 商户ID（对应 fast-api 的 EpayId）
	Key       string `yaml:"key"`        // 签名密钥（对应 fast-api 的 EpayKey）
}

type AlipayConfig struct {
	AppID           string `yaml:"app_id"`            // 支付宝应用 APPID
	PrivateKey      string `yaml:"private_key"`        // 商户私钥（PKCS1 格式）
	AliPayPublicKey string `yaml:"alipay_public_key"`  // 支付宝公钥
	IsProduction    bool   `yaml:"is_production"`      // true=正式环境, false=沙箱
	NotifyURL       string `yaml:"notify_url"`         // 支付宝异步通知地址（必须公网可访问）
	ReturnURL       string `yaml:"return_url"`         // 支付宝同步跳转地址
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Epay    EpayConfig    `yaml:"epay"`
	Alipay  AlipayConfig  `yaml:"alipay"`
	Wechat  WechatConfig  `yaml:"wechat"`
}

// WechatConfig holds WeChat Pay APIv3 credentials
// is_configured is false by default — only set to true in main when all fields are non-empty
type WechatConfig struct {
	MchID      string `yaml:"mch_id"`      // 商户号
	APIv3Key   string `yaml:"api_v3_key"`  // APIv3 密钥（32位）
	SerialNo   string `yaml:"serial_no"`   // 商户证书序列号
	PrivateKey string `yaml:"private_key"` // 商户私钥（PKCS8 格式）
	NotifyURL  string `yaml:"notify_url"`  // 微信支付异步通知地址
	Configured bool   `yaml:"-"`           // 运行时标记，不持久化到 yaml
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = "config.yaml"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file %s: %w", path, err)
	}

	// Defaults
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}

	return &cfg, nil
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
