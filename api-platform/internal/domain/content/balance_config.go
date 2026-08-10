package content

import (
	"context"
	"encoding/json"
	"math"
	"strings"
)

// BalanceConfigKey 对齐 CRMEB 余额设置（systemForm/Basics/balance）。
// 存储：qixi_crm_a_setting_cache.key = balance_config
// 字段 key 对齐 CRMEB systemConfig：
// balance_func_status / recharge_switch / store_user_min_recharge / recharge_attention
const BalanceConfigKey = "balance_config"

const defaultRechargeAttention = "1、账户充值仅限用于购买商城内商品，不可提现\n2、账户充值成功后，一般1～5分钟到账\n3、如有疑问，请联系客服"

type balanceConfig struct {
	// BalanceFuncStatus 余额功能：0 关闭 / 1 开启（对齐 CRMEB balance_func_status）
	BalanceFuncStatus int `json:"balance_func_status"`
	// RechargeSwitch 余额充值开关：0 关闭 / 1 开启（对齐 CRMEB recharge_switch）
	RechargeSwitch int `json:"recharge_switch"`
	// StoreUserMinRecharge 用户最低充值金额（对齐 CRMEB store_user_min_recharge）
	StoreUserMinRecharge float64 `json:"store_user_min_recharge"`
	// RechargeAttention 充值注意事项，多行文本（对齐 CRMEB recharge_attention）
	RechargeAttention string `json:"recharge_attention"`
}

func defaultBalanceConfig() balanceConfig {
	return balanceConfig{
		BalanceFuncStatus:    1,
		RechargeSwitch:       1,
		StoreUserMinRecharge: 1,
		RechargeAttention:    defaultRechargeAttention,
	}
}

func parseBalanceConfig(raw string) (balanceConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return balanceConfig{}, ErrBadParam
	}
	allowed := map[string]struct{}{
		"balance_func_status":     {},
		"recharge_switch":         {},
		"store_user_min_recharge": {},
		"recharge_attention":      {},
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return balanceConfig{}, ErrBadParam
		}
		if _, ok := allowed[key]; !ok {
			return balanceConfig{}, ErrBadParam
		}
	}
	var config balanceConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return balanceConfig{}, ErrBadParam
	}
	if err := validateBalanceConfig(&config); err != nil {
		return balanceConfig{}, err
	}
	return config, nil
}

func validateBalanceConfig(config *balanceConfig) error {
	if config.BalanceFuncStatus != 0 && config.BalanceFuncStatus != 1 {
		return ErrBadParam
	}
	if config.RechargeSwitch != 0 && config.RechargeSwitch != 1 {
		return ErrBadParam
	}
	if config.StoreUserMinRecharge < 0 ||
		math.IsNaN(config.StoreUserMinRecharge) ||
		math.IsInf(config.StoreUserMinRecharge, 0) {
		return ErrBadParam
	}
	return nil
}

func (s *Service) GetBalanceConfig(ctx context.Context) (string, error) {
	return getJSONConfig(s, ctx, BalanceConfigKey, defaultBalanceConfig(), parseBalanceConfig)
}

func (s *Service) SaveBalanceConfig(ctx context.Context, raw string) (string, error) {
	return saveJSONConfig(s, ctx, BalanceConfigKey, raw, parseBalanceConfig)
}
