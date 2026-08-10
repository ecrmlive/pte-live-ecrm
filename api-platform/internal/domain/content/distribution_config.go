package content

import (
	"context"
	"encoding/json"
	"math"
	"strings"
)

// DistributionConfigKey 对齐 CRMEB ConfigOthers::update（systemForm/Basics/distribution_tabs）。
const DistributionConfigKey = "distribution_config"

type distributionConfig struct {
	ExtensionStatus    bool     `json:"extension_status"`
	ExtensionSelf      bool     `json:"extension_self"`
	ExtensionLimit     bool     `json:"extension_limit"`
	ExtensionLimitDay  int      `json:"extension_limit_day"`
	PromoterType       int      `json:"promoter_type"`
	PromoterLowMoney   float64  `json:"promoter_low_money"`
	ExtensionPop       int      `json:"extension_pop"`
	ExtensionOneRate   float64  `json:"extension_one_rate"`
	ExtensionTwoRate   float64  `json:"extension_two_rate"`
	UserExtractMin     float64  `json:"user_extract_min"`
	LockBrokerageTimer int      `json:"lock_brokerage_timer"`
	SysExtensionType   int      `json:"sys_extension_type"`
	WithdrawType       []string `json:"withdraw_type"`
	ExtractSwitch      int      `json:"extract_switch"`
	TransferSceneID    int      `json:"transfer_scene_id"`
	MaxBagNumber      int      `json:"max_bag_number"`
}

func defaultDistributionConfig() distributionConfig {
	return distributionConfig{
		ExtensionStatus:    false,
		ExtensionSelf:      false,
		ExtensionLimit:     false,
		ExtensionLimitDay:  15,
		PromoterType:       0,
		PromoterLowMoney:   0,
		ExtensionPop:       0,
		ExtensionOneRate:   0,
		ExtensionTwoRate:   0,
		UserExtractMin:     1,
		LockBrokerageTimer: 0,
		SysExtensionType:   0,
		WithdrawType:       []string{"1"},
		ExtractSwitch:      1,
		TransferSceneID:    0,
		MaxBagNumber:      10,
	}
}

func parseDistributionConfig(raw string) (distributionConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return distributionConfig{}, ErrBadParam
	}
	allowed := map[string]struct{}{
		"extension_status":     {},
		"extension_self":       {},
		"extension_limit":      {},
		"extension_limit_day":  {},
		"promoter_type":        {},
		"promoter_low_money":   {},
		"extension_pop":        {},
		"extension_one_rate":   {},
		"extension_two_rate":   {},
		"user_extract_min":     {},
		"lock_brokerage_timer": {},
		"sys_extension_type":   {},
		"withdraw_type":        {},
		"extract_switch":       {},
		"transfer_scene_id":    {},
		"max_bag_number":       {},
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return distributionConfig{}, ErrBadParam
		}
		if _, ok := allowed[key]; !ok {
			return distributionConfig{}, ErrBadParam
		}
	}
	var config distributionConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return distributionConfig{}, ErrBadParam
	}
	if err := validateDistributionConfig(&config); err != nil {
		return distributionConfig{}, err
	}
	return config, nil
}

func validateDistributionConfig(config *distributionConfig) error {
	if config.ExtensionLimitDay <= 0 || config.ExtensionLimitDay > 3650 {
		return ErrBadParam
	}
	if config.PromoterType < 0 || config.PromoterType > 3 {
		return ErrBadParam
	}
	if config.PromoterType == 3 && config.PromoterLowMoney <= 0 {
		return ErrBadParam
	}
	if config.PromoterLowMoney < 0 || config.PromoterLowMoney > 1_000_000 {
		return ErrBadParam
	}
	if config.ExtensionPop < 0 || config.ExtensionPop > 3 {
		return ErrBadParam
	}
	if config.ExtensionOneRate < 0 || config.ExtensionTwoRate < 0 {
		return ErrBadParam
	}
	if config.ExtensionOneRate > 1 || config.ExtensionTwoRate > 1 {
		return ErrBadParam
	}
	// 对齐 CRMEB：一级比例不能小于二级；之和不能超过 1。
	if round4(config.ExtensionOneRate) < round4(config.ExtensionTwoRate) {
		return ErrBadParam
	}
	if round4(config.ExtensionOneRate+config.ExtensionTwoRate) > 1 {
		return ErrBadParam
	}
	if config.UserExtractMin < 0 || config.UserExtractMin > 1_000_000 {
		return ErrBadParam
	}
	if config.LockBrokerageTimer < 0 || config.LockBrokerageTimer > 3650 {
		return ErrBadParam
	}
	if config.SysExtensionType < 0 || config.SysExtensionType > 2 {
		return ErrBadParam
	}
	if config.ExtractSwitch < 1 || config.ExtractSwitch > 2 {
		return ErrBadParam
	}
	if config.TransferSceneID < 0 || config.TransferSceneID > 999999 {
		return ErrBadParam
	}
	if config.MaxBagNumber < 0 || config.MaxBagNumber > 9999 {
		return ErrBadParam
	}
	config.ExtensionOneRate = round4(config.ExtensionOneRate)
	config.ExtensionTwoRate = round4(config.ExtensionTwoRate)
	config.UserExtractMin = round2(config.UserExtractMin)
	config.PromoterLowMoney = round2(config.PromoterLowMoney)
	config.WithdrawType = normalizeWithdrawType(config.WithdrawType)
	return nil
}

func normalizeWithdrawType(in []string) []string {
	allowed := map[string]struct{}{"0": {}, "1": {}, "2": {}, "4": {}}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, item := range in {
		item = strings.TrimSpace(item)
		if _, ok := allowed[item]; !ok {
			continue
		}
		if _, dup := seen[item]; dup {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	if len(out) == 0 {
		return []string{"1"}
	}
	return out
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func (s *Service) GetDistributionConfig(ctx context.Context) (string, error) {
	return getJSONConfig(s, ctx, DistributionConfigKey, defaultDistributionConfig(), parseDistributionConfig)
}

func (s *Service) SaveDistributionConfig(ctx context.Context, raw string) (string, error) {
	return saveJSONConfig(s, ctx, DistributionConfigKey, raw, parseDistributionConfig)
}
