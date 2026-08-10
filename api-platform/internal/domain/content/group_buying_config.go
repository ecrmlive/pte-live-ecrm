package content

import (
	"context"
	"encoding/json"
	"strings"
)

// GroupBuyingConfigKey 对齐 CRMEB ConfigOthers::getGroupBuying / setGroupBuying。
// 存储：qixi_crm_a_setting_cache.key = group_buying_config
const GroupBuyingConfigKey = "group_buying_config"

type groupBuyingConfig struct {
	// FictiStatus 虚拟成团启用：0 关闭 / 1 启用（对齐 CRMEB ficti_status）
	FictiStatus int `json:"ficti_status"`
	// GroupBuyingRate 真实成团最小比例 0～100（对齐 CRMEB group_buying_rate）
	GroupBuyingRate int `json:"group_buying_rate"`
}

func defaultGroupBuyingConfig() groupBuyingConfig {
	return groupBuyingConfig{
		FictiStatus:     1,
		GroupBuyingRate: 30,
	}
}

func parseGroupBuyingConfig(raw string) (groupBuyingConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return groupBuyingConfig{}, ErrBadParam
	}
	allowed := map[string]struct{}{
		"ficti_status":      {},
		"group_buying_rate": {},
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return groupBuyingConfig{}, ErrBadParam
		}
		if _, ok := allowed[key]; !ok {
			return groupBuyingConfig{}, ErrBadParam
		}
	}
	var config groupBuyingConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return groupBuyingConfig{}, ErrBadParam
	}
	if err := validateGroupBuyingConfig(&config); err != nil {
		return groupBuyingConfig{}, err
	}
	return config, nil
}

func validateGroupBuyingConfig(config *groupBuyingConfig) error {
	if config.FictiStatus != 0 && config.FictiStatus != 1 {
		return ErrBadParam
	}
	if config.GroupBuyingRate < 0 || config.GroupBuyingRate > 100 {
		return ErrBadParam
	}
	return nil
}

func (s *Service) GetGroupBuyingConfig(ctx context.Context) (string, error) {
	return getJSONConfig(s, ctx, GroupBuyingConfigKey, defaultGroupBuyingConfig(), parseGroupBuyingConfig)
}

func (s *Service) SaveGroupBuyingConfig(ctx context.Context, raw string) (string, error) {
	return saveJSONConfig(s, ctx, GroupBuyingConfigKey, raw, parseGroupBuyingConfig)
}
