package content

import (
	"context"
	"encoding/json"
	"math"
	"regexp"
	"strings"
)

// UserSetupConfigKey 对齐 CRMEB 用户设置 /user/setup_user
// （UserInfo::saveAll + User::saveRegisterConfig）。
// 存储：qixi_crm_a_setting_cache.key = user_setup_config
const UserSetupConfigKey = "user_setup_config"

var userSetupFieldNameRe = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)

var userSetupFieldTypes = map[string]string{
	"input":   "文本",
	"int":     "数字",
	"phone":   "手机号",
	"date":    "时间",
	"radio":   "单选",
	"address": "地址",
	"id_card": "身份证",
	"email":   "邮箱",
}

type userSetupField struct {
	ID        int64    `json:"id"`
	Field     string   `json:"field"`
	Title     string   `json:"title"`
	IsUsed    int      `json:"is_used"`
	IsRequire int      `json:"is_require"`
	IsShow    int      `json:"is_show"`
	Type      string   `json:"type"`
	Msg       string   `json:"msg"`
	Content   []string `json:"content,omitempty"`
	IsDefault int      `json:"is_default"`
	Sort      int      `json:"sort"`
}

type userSetupCoupon struct {
	CouponID     int64   `json:"coupon_id"`
	Title        string  `json:"title"`
	CouponType   int     `json:"coupon_type"`
	CouponPrice  float64 `json:"coupon_price"`
	UseMinPrice  float64 `json:"use_min_price"`
	CouponTime   int     `json:"coupon_time"`
	IsTimeout    int     `json:"is_timeout"`
	UseStartTime string  `json:"use_start_time,omitempty"`
	UseEndTime   string  `json:"use_end_time,omitempty"`
}

// 对齐 CRMEB systemConfig：user_default_avatar / is_phone_login /
// first_avatar_switch / open_update_info / wechat_phone_switch /
// newcomer_status / register_popup_pic / register_*_status / register_give_*
type userSetupConfig struct {
	UserDefaultAvatar      string            `json:"user_default_avatar"`
	Fields                 []userSetupField  `json:"fields"`
	IsPhoneLogin           int               `json:"is_phone_login"`
	FirstAvatarSwitch      int               `json:"first_avatar_switch"`
	OpenUpdateInfo         int               `json:"open_update_info"`
	WechatPhoneSwitch      int               `json:"wechat_phone_switch"`
	NewcomerStatus         int               `json:"newcomer_status"`
	RegisterPopupPic       string            `json:"register_popup_pic"`
	RegisterMoneyStatus    int               `json:"register_money_status"`
	RegisterGiveMoney      float64           `json:"register_give_money"`
	RegisterIntegralStatus int               `json:"register_integral_status"`
	RegisterGiveIntegral   int               `json:"register_give_integral"`
	RegisterCouponStatus   int               `json:"register_coupon_status"`
	RegisterGiveCoupon     []userSetupCoupon `json:"register_give_coupon"`
}

func defaultUserSetupFields() []userSetupField {
	// 对齐 CRMEB 用户设置截图常见系统字段：姓名/性别/生日/地址/备注/身份证（实名认证）
	return []userSetupField{
		{ID: 1, Field: "real_name", Title: "姓名", IsUsed: 1, IsRequire: 0, IsShow: 1, Type: "input", Msg: "请填写真实姓名", IsDefault: 1, Sort: 0},
		{ID: 2, Field: "sex", Title: "性别", IsUsed: 1, IsRequire: 0, IsShow: 1, Type: "radio", Msg: "请选择性别", Content: []string{"男", "女", "保密"}, IsDefault: 1, Sort: 1},
		{ID: 3, Field: "birthday", Title: "生日", IsUsed: 1, IsRequire: 0, IsShow: 1, Type: "date", Msg: "请选择生日", IsDefault: 1, Sort: 2},
		{ID: 4, Field: "address", Title: "地址", IsUsed: 1, IsRequire: 0, IsShow: 1, Type: "address", Msg: "请选择地址", IsDefault: 1, Sort: 3},
		{ID: 5, Field: "mark", Title: "备注", IsUsed: 1, IsRequire: 0, IsShow: 0, Type: "input", Msg: "请填写备注", IsDefault: 1, Sort: 4},
		{ID: 6, Field: "id_card", Title: "身份证（实名认证）", IsUsed: 0, IsRequire: 0, IsShow: 0, Type: "id_card", Msg: "请填写身份证号", IsDefault: 1, Sort: 5},
	}
}

func defaultUserSetupConfig() userSetupConfig {
	return userSetupConfig{
		UserDefaultAvatar:      "",
		Fields:                 defaultUserSetupFields(),
		IsPhoneLogin:           0,
		FirstAvatarSwitch:      1,
		OpenUpdateInfo:         1,
		WechatPhoneSwitch:      0,
		NewcomerStatus:         0,
		RegisterPopupPic:       "",
		RegisterMoneyStatus:    0,
		RegisterGiveMoney:      0,
		RegisterIntegralStatus: 0,
		RegisterGiveIntegral:   0,
		RegisterCouponStatus:   0,
		RegisterGiveCoupon:     []userSetupCoupon{},
	}
}

func parseUserSetupConfig(raw string) (userSetupConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return userSetupConfig{}, ErrBadParam
	}
	allowed := map[string]struct{}{
		"user_default_avatar":      {},
		"fields":                   {},
		"is_phone_login":           {},
		"first_avatar_switch":      {},
		"open_update_info":         {},
		"wechat_phone_switch":      {},
		"newcomer_status":          {},
		"register_popup_pic":       {},
		"register_money_status":    {},
		"register_give_money":      {},
		"register_integral_status": {},
		"register_give_integral":   {},
		"register_coupon_status":   {},
		"register_give_coupon":     {},
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return userSetupConfig{}, ErrBadParam
		}
		if _, ok := allowed[key]; !ok {
			return userSetupConfig{}, ErrBadParam
		}
	}
	var config userSetupConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return userSetupConfig{}, ErrBadParam
	}
	if err := validateUserSetupConfig(&config); err != nil {
		return userSetupConfig{}, err
	}
	return config, nil
}

func validateUserSetupConfig(config *userSetupConfig) error {
	if config == nil {
		return ErrBadParam
	}
	config.UserDefaultAvatar = strings.TrimSpace(config.UserDefaultAvatar)
	config.RegisterPopupPic = strings.TrimSpace(config.RegisterPopupPic)
	if len([]rune(config.UserDefaultAvatar)) > 512 || len([]rune(config.RegisterPopupPic)) > 512 {
		return ErrBadParam
	}
	for _, v := range []int{
		config.IsPhoneLogin,
		config.FirstAvatarSwitch,
		config.OpenUpdateInfo,
		config.WechatPhoneSwitch,
		config.NewcomerStatus,
		config.RegisterMoneyStatus,
		config.RegisterIntegralStatus,
		config.RegisterCouponStatus,
	} {
		if v != 0 && v != 1 {
			return ErrBadParam
		}
	}
	if config.RegisterGiveMoney < 0 ||
		math.IsNaN(config.RegisterGiveMoney) ||
		math.IsInf(config.RegisterGiveMoney, 0) ||
		config.RegisterGiveMoney > 99999 {
		return ErrBadParam
	}
	if config.RegisterGiveIntegral < 0 || config.RegisterGiveIntegral > 999999 {
		return ErrBadParam
	}
	// 注册有礼关闭时不强制校验赠送礼品明细（与 CRMEB 隐藏子项一致）
	if config.NewcomerStatus == 1 {
		if config.RegisterMoneyStatus == 1 && config.RegisterGiveMoney <= 0 {
			return ErrBadParam
		}
		if config.RegisterIntegralStatus == 1 && config.RegisterGiveIntegral <= 0 {
			return ErrBadParam
		}
	}
	if config.Fields == nil {
		config.Fields = []userSetupField{}
	}
	if len(config.Fields) == 0 {
		return ErrBadParam
	}
	seenField := map[string]struct{}{}
	seenID := map[int64]struct{}{}
	maxID := int64(0)
	for i := range config.Fields {
		f := &config.Fields[i]
		f.Field = strings.ToLower(strings.TrimSpace(f.Field))
		f.Title = strings.TrimSpace(f.Title)
		f.Msg = strings.TrimSpace(f.Msg)
		f.Type = strings.TrimSpace(f.Type)
		if f.ID <= 0 {
			return ErrBadParam
		}
		if _, ok := seenID[f.ID]; ok {
			return ErrBadParam
		}
		seenID[f.ID] = struct{}{}
		if f.ID > maxID {
			maxID = f.ID
		}
		if !userSetupFieldNameRe.MatchString(f.Field) {
			return ErrBadParam
		}
		if _, ok := seenField[f.Field]; ok {
			return ErrBadParam
		}
		seenField[f.Field] = struct{}{}
		if f.Title == "" || len([]rune(f.Title)) > 64 {
			return ErrBadParam
		}
		if f.Msg == "" || len([]rune(f.Msg)) > 255 {
			return ErrBadParam
		}
		if _, ok := userSetupFieldTypes[f.Type]; !ok {
			return ErrBadParam
		}
		for _, flag := range []int{f.IsUsed, f.IsRequire, f.IsShow, f.IsDefault} {
			if flag != 0 && flag != 1 {
				return ErrBadParam
			}
		}
		if f.IsUsed == 0 {
			f.IsRequire = 0
			f.IsShow = 0
		}
		if f.Type == "radio" {
			opts := make([]string, 0, len(f.Content))
			for _, opt := range f.Content {
				opt = strings.TrimSpace(opt)
				if opt == "" {
					continue
				}
				if len([]rune(opt)) > 32 {
					return ErrBadParam
				}
				opts = append(opts, opt)
			}
			if len(opts) < 2 {
				return ErrBadParam
			}
			f.Content = opts
		} else {
			f.Content = nil
		}
		f.Sort = i
	}
	// 系统默认字段缺失时补齐（兼容旧配置仅含 email/phone、无 mark 的情况）
	defaults := defaultUserSetupFields()
	for _, d := range defaults {
		if _, ok := seenField[d.Field]; ok {
			continue
		}
		maxID++
		added := d
		added.ID = maxID
		added.Sort = len(config.Fields)
		config.Fields = append(config.Fields, added)
		seenField[d.Field] = struct{}{}
		seenID[added.ID] = struct{}{}
	}
	for i := range config.Fields {
		f := &config.Fields[i]
		for _, d := range defaults {
			if f.Field == d.Field {
				f.IsDefault = 1
				f.Type = d.Type
				break
			}
		}
	}

	if config.RegisterGiveCoupon == nil {
		config.RegisterGiveCoupon = []userSetupCoupon{}
	}
	if config.NewcomerStatus == 1 &&
		config.RegisterCouponStatus == 1 &&
		len(config.RegisterGiveCoupon) == 0 {
		return ErrBadParam
	}
	seenCoupon := map[int64]struct{}{}
	for i := range config.RegisterGiveCoupon {
		c := &config.RegisterGiveCoupon[i]
		if c.CouponID <= 0 {
			return ErrBadParam
		}
		if _, ok := seenCoupon[c.CouponID]; ok {
			return ErrBadParam
		}
		seenCoupon[c.CouponID] = struct{}{}
		c.Title = strings.TrimSpace(c.Title)
		c.UseStartTime = strings.TrimSpace(c.UseStartTime)
		c.UseEndTime = strings.TrimSpace(c.UseEndTime)
		if c.Title == "" || len([]rune(c.Title)) > 128 {
			return ErrBadParam
		}
		if c.CouponPrice < 0 || c.UseMinPrice < 0 || c.CouponTime < 0 {
			return ErrBadParam
		}
		if c.IsTimeout != 0 && c.IsTimeout != 1 {
			return ErrBadParam
		}
	}
	return nil
}

func (s *Service) GetUserSetupConfig(ctx context.Context) (string, error) {
	return getJSONConfig(s, ctx, UserSetupConfigKey, defaultUserSetupConfig(), parseUserSetupConfig)
}

func (s *Service) SaveUserSetupConfig(ctx context.Context, raw string) (string, error) {
	return saveJSONConfig(s, ctx, UserSetupConfigKey, raw, parseUserSetupConfig)
}
