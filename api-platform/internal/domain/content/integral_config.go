package content

import (
	"context"
	"encoding/json"
	"math"
	"strings"
)

// IntegralConfigKey 对齐 CRMEB UserIntegral::getConfig / saveConfig。
// 存储：qixi_crm_a_setting_cache.key = integral_config
// 说明正文 rule 对齐 CRMEB CacheRepository::INTEGRAL_RULE（sys_integral_rule），一并写入本 JSON。
const IntegralConfigKey = "integral_config"

type integralConfig struct {
	// IntegralStatus 积分开关：0 关闭 / 1 开启（对齐 CRMEB integral_status）
	IntegralStatus int `json:"integral_status"`
	// IntegralMoney 积分抵用金额（元）：1 积分抵多少金额（对齐 CRMEB integral_money）
	IntegralMoney float64 `json:"integral_money"`
	// IntegralOrderRate 下单赠送积分比例：消费 1 元赠送多少积分（对齐 CRMEB integral_order_rate）
	IntegralOrderRate float64 `json:"integral_order_rate"`
	// IntegralFreeze 下单赠送积分冻结期（天）（对齐 CRMEB integral_freeze）
	IntegralFreeze int `json:"integral_freeze"`
	// IntegralClearTime 积分清除时间设置（月）（对齐 CRMEB integral_clear_time）
	IntegralClearTime int `json:"integral_clear_time"`
	// IntegralUserGive 邀请好友赠送积分（分）（对齐 CRMEB integral_user_give）
	IntegralUserGive int `json:"integral_user_give"`
	// IntegralCommunityGive 发布种草可获得积分（分）（对齐 CRMEB integral_community_give）
	IntegralCommunityGive int `json:"integral_community_give"`
	// IntegralCommunityGiveLimit 发布种草篇数限量（对齐 CRMEB integral_community_give_limit）
	IntegralCommunityGiveLimit int `json:"integral_community_give_limit"`
	// Rule 积分说明（富文本，对齐 CRMEB sys_integral_rule / rule）
	Rule string `json:"rule"`
}

func defaultIntegralConfig() integralConfig {
	return integralConfig{
		IntegralStatus:             1,
		IntegralMoney:              0.1,
		IntegralOrderRate:          1,
		IntegralFreeze:             0,
		IntegralClearTime:          24,
		IntegralUserGive:           50,
		IntegralCommunityGive:      10,
		IntegralCommunityGiveLimit: 10,
		Rule:                       "",
	}
}

func parseIntegralConfig(raw string) (integralConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return integralConfig{}, ErrBadParam
	}
	allowed := map[string]struct{}{
		"integral_status":               {},
		"integral_money":                {},
		"integral_order_rate":           {},
		"integral_freeze":               {},
		"integral_clear_time":           {},
		"integral_user_give":            {},
		"integral_community_give":       {},
		"integral_community_give_limit": {},
		"rule":                          {},
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return integralConfig{}, ErrBadParam
		}
		if _, ok := allowed[key]; !ok {
			return integralConfig{}, ErrBadParam
		}
	}
	var config integralConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return integralConfig{}, ErrBadParam
	}
	if err := validateIntegralConfig(&config); err != nil {
		return integralConfig{}, err
	}
	return config, nil
}

func validateIntegralConfig(config *integralConfig) error {
	if config.IntegralStatus != 0 && config.IntegralStatus != 1 {
		return ErrBadParam
	}
	if config.IntegralMoney < 0 || math.IsNaN(config.IntegralMoney) || math.IsInf(config.IntegralMoney, 0) {
		return ErrBadParam
	}
	if config.IntegralOrderRate < 0 || math.IsNaN(config.IntegralOrderRate) || math.IsInf(config.IntegralOrderRate, 0) {
		return ErrBadParam
	}
	if config.IntegralFreeze < 0 {
		return ErrBadParam
	}
	if config.IntegralClearTime < 0 {
		return ErrBadParam
	}
	if config.IntegralUserGive < 0 {
		return ErrBadParam
	}
	if config.IntegralCommunityGive < 0 || config.IntegralCommunityGive > 9999 {
		return ErrBadParam
	}
	if config.IntegralCommunityGiveLimit < 0 || config.IntegralCommunityGiveLimit > 9999 {
		return ErrBadParam
	}
	return nil
}

func (s *Service) GetIntegralConfig(ctx context.Context) (string, error) {
	return getJSONConfig(s, ctx, IntegralConfigKey, defaultIntegralConfig(), parseIntegralConfig)
}

func (s *Service) SaveIntegralConfig(ctx context.Context, raw string) (string, error) {
	return saveJSONConfig(s, ctx, IntegralConfigKey, raw, parseIntegralConfig)
}
