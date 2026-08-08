package content

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// agentZoneConfigKey 平台「区域代理 → 代理设置」：默认三级提成 + 代理申请表单字段。
// 对齐 CRMEB system config 组 circle_config（one/two/three_agent_commission + agent_application_form）。
const agentZoneConfigKey = "agent_zone_setting"

type agentZoneConfig struct {
	OneAgentCommission   float64              `json:"one_agent_commission"`
	TwoAgentCommission   float64              `json:"two_agent_commission"`
	ThreeAgentCommission float64              `json:"three_agent_commission"`
	FormFields           []merchantApplyField `json:"form_fields"`
}

func defaultAgentZoneConfig() agentZoneConfig {
	return agentZoneConfig{
		OneAgentCommission:   0,
		TwoAgentCommission:   0,
		ThreeAgentCommission: 0,
		FormFields:           []merchantApplyField{},
	}
}

func marshalAgentZoneConfig(config agentZoneConfig) string {
	if config.FormFields == nil {
		config.FormFields = []merchantApplyField{}
	}
	data, _ := json.Marshal(config)
	return string(data)
}

func parseAgentZoneConfig(raw string) (agentZoneConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return agentZoneConfig{}, ErrBadParam
	}
	for key := range fields {
		switch key {
		case "one_agent_commission", "two_agent_commission", "three_agent_commission", "form_fields":
		default:
			return agentZoneConfig{}, ErrBadParam
		}
	}
	var config agentZoneConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return agentZoneConfig{}, ErrBadParam
	}
	if !isValidCommissionRate(config.OneAgentCommission) ||
		!isValidCommissionRate(config.TwoAgentCommission) ||
		!isValidCommissionRate(config.ThreeAgentCommission) {
		return agentZoneConfig{}, ErrBadParam
	}
	// 对齐 CRMEB：一级 ≥ 二级 ≥ 三级
	if config.OneAgentCommission < config.TwoAgentCommission ||
		config.TwoAgentCommission < config.ThreeAgentCommission {
		return agentZoneConfig{}, ErrBadParam
	}
	normalized, err := normalizeApplyFormFields(config.FormFields)
	if err != nil {
		return agentZoneConfig{}, err
	}
	config.FormFields = normalized
	return config, nil
}

func isValidCommissionRate(v float64) bool {
	return v >= 0 && v <= 100
}

func (s *Service) GetAgentZoneConfig(ctx context.Context) (string, error) {
	row, err := s.store.GetCache(ctx, agentZoneConfigKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return marshalAgentZoneConfig(defaultAgentZoneConfig()), nil
		}
		return "", err
	}
	config, err := parseAgentZoneConfig(row.Result)
	if err != nil {
		return marshalAgentZoneConfig(defaultAgentZoneConfig()), nil
	}
	return marshalAgentZoneConfig(config), nil
}

func (s *Service) SaveAgentZoneConfig(ctx context.Context, raw string) (string, error) {
	config, err := parseAgentZoneConfig(raw)
	if err != nil {
		return "", err
	}
	canonical := marshalAgentZoneConfig(config)
	if err := s.store.UpsertCache(ctx, &Cache{Key: agentZoneConfigKey, ExpireTime: 0, Result: canonical}); err != nil {
		return "", err
	}
	return canonical, nil
}
