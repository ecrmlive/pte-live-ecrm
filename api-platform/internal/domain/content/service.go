package content

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Store interface {
	ListNotices(ctx context.Context, onlyShow bool, page, limit int) ([]Notice, int64, error)
	GetNotice(ctx context.Context, id uint) (*Notice, error)
	CreateNotice(ctx context.Context, n *Notice) error
	UpdateNotice(ctx context.Context, n *Notice) error
	SoftDeleteNotice(ctx context.Context, id uint) error
	GetCache(ctx context.Context, key string) (*Cache, error)
	UpsertCache(ctx context.Context, row *Cache) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListApp(ctx context.Context, page, limit int) (*PageResult[Notice], error) {
	list, total, err := s.store.ListNotices(ctx, true, page, limit)
	if err != nil {
		return nil, err
	}
	page, limit = normalize(page, limit)
	return &PageResult[Notice]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListAdmin(ctx context.Context, page, limit int) (*PageResult[Notice], error) {
	list, total, err := s.store.ListNotices(ctx, false, page, limit)
	if err != nil {
		return nil, err
	}
	page, limit = normalize(page, limit)
	return &PageResult[Notice]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*Notice, error) {
	n, err := s.store.GetNotice(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return n, nil
}

func (s *Service) Create(ctx context.Context, in NoticeInput) (*Notice, error) {
	title := strings.TrimSpace(in.Title)
	content := strings.TrimSpace(in.Content)
	if title == "" || content == "" {
		return nil, ErrBadParam
	}
	n := &Notice{Title: title, Content: content, IsShow: 1, Sort: in.Sort}
	if in.IsShow != nil {
		n.IsShow = *in.IsShow
	}
	if err := s.store.CreateNotice(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) Update(ctx context.Context, id uint, in NoticeInput) (*Notice, error) {
	n, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(in.Title)
	content := strings.TrimSpace(in.Content)
	if title == "" || content == "" {
		return nil, ErrBadParam
	}
	n.Title = title
	n.Content = content
	n.Sort = in.Sort
	if in.IsShow != nil {
		n.IsShow = *in.IsShow
	}
	if err := s.store.UpdateNotice(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.store.SoftDeleteNotice(ctx, id)
}

// AgreeCatalog 平台可维护的协议键（对齐 CRMEB CacheRepository 常用项；key≤32）。
func AgreeCatalog() []AgreeMeta {
	return []AgreeMeta{
		{Key: "sys_user_agree", Label: "用户协议"},
		{Key: "sys_userr_privacy", Label: "隐私政策"},
		{Key: "sys_svip", Label: "付费会员协议"},
		{Key: "sys_product_presell_agree", Label: "预售协议"},
		{Key: "business_entry_agree", Label: "商户入驻协议"},
		{Key: "promoter_explain", Label: "分销说明"},
		{Key: "sys_coupon_agree", Label: "优惠券使用说明"},
		{Key: "sys_extension_agree", Label: "佣金说明"},
		{Key: "sys_brokerage", Label: "分销等级规则"},
		{Key: "sys_about_us", Label: "关于我们"},
		{Key: "sys_refund_agree", Label: "退款协议"},
		{Key: "sys_cancel_agree", Label: "取消订单说明"},
		{Key: "sys_recharge_agree", Label: "充值协议"},
		{Key: "sys_integral_agree", Label: "积分规则"},
		{Key: "mer_settle_agree", Label: "商户结算说明"},
		{Key: "sys_lottery_agree", Label: "抽奖活动说明"},
		{Key: "sys_deposit_agree", Label: "保证金说明"},
		// 店铺设置 / 说明提示（对齐 CRMEB sys_merchant_type / sys_merchant_category）
		{Key: "sys_merchant_type", Label: "店铺类型说明"},
		{Key: "sys_merchant_category", Label: "店铺分类说明"},
	}
}

func (s *Service) ListAgreements(ctx context.Context) ([]AgreeView, error) {
	out := make([]AgreeView, 0, len(AgreeCatalog()))
	for _, m := range AgreeCatalog() {
		row, err := s.store.GetCache(ctx, m.Key)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				out = append(out, AgreeView{Key: m.Key, Label: m.Label, Content: ""})
				continue
			}
			return nil, err
		}
		out = append(out, AgreeView{Key: m.Key, Label: m.Label, Content: row.Result})
	}
	return out, nil
}

func (s *Service) GetAgreement(ctx context.Context, key string) (*AgreeView, error) {
	meta, ok := agreeMeta(key)
	if !ok {
		return nil, ErrAgreeNotFound
	}
	row, err := s.store.GetCache(ctx, key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &AgreeView{Key: meta.Key, Label: meta.Label, Content: ""}, nil
		}
		return nil, err
	}
	return &AgreeView{Key: meta.Key, Label: meta.Label, Content: row.Result}, nil
}

func (s *Service) SaveAgreement(ctx context.Context, key string, in AgreeSaveInput) (*AgreeView, error) {
	meta, ok := agreeMeta(key)
	if !ok {
		return nil, ErrAgreeNotFound
	}
	content := strings.TrimSpace(in.Content)
	if content == "" {
		return nil, ErrBadParam
	}
	if err := s.store.UpsertCache(ctx, &Cache{Key: key, ExpireTime: 0, Result: content}); err != nil {
		return nil, err
	}
	return &AgreeView{Key: meta.Key, Label: meta.Label, Content: content}, nil
}

func agreeMeta(key string) (AgreeMeta, bool) {
	key = strings.TrimSpace(key)
	for _, m := range AgreeCatalog() {
		if m.Key == key {
			return m, true
		}
	}
	return AgreeMeta{}, false
}

const smsConfigKey = "sms_config"

// smsStubConfig intentionally has no provider credentials. The current product
// baseline only exposes a non-delivery stub; real SMS channel secrets belong in
// an external secret manager, never in the unified-admin cache or response.
type smsStubConfig struct {
	Enabled  bool   `json:"enabled"`
	Provider string `json:"provider"`
	Sign     string `json:"sign"`
	Remark   string `json:"remark"`
}

func defaultSMSConfig() smsStubConfig {
	return smsStubConfig{Enabled: false, Provider: "stub", Sign: "七禧商城", Remark: "未配置真实短信通道"}
}

func marshalSMSConfig(config smsStubConfig) string {
	data, _ := json.Marshal(config)
	return string(data)
}

func parseSMSConfig(raw string) (smsStubConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return smsStubConfig{}, ErrBadParam
	}
	for key := range fields {
		switch key {
		case "enabled", "provider", "sign", "remark":
		default:
			return smsStubConfig{}, ErrBadParam
		}
	}
	var config smsStubConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil || strings.TrimSpace(config.Provider) != "stub" || len([]rune(config.Sign)) > 64 || len([]rune(config.Remark)) > 500 {
		return smsStubConfig{}, ErrBadParam
	}
	config.Provider = "stub"
	config.Sign = strings.TrimSpace(config.Sign)
	config.Remark = strings.TrimSpace(config.Remark)
	if config.Sign == "" {
		config.Sign = defaultSMSConfig().Sign
	}
	if config.Remark == "" {
		config.Remark = defaultSMSConfig().Remark
	}
	return config, nil
}

// GetSMSConfig returns only a validated non-secret stub configuration. Legacy
// malformed or secret-bearing cache values are never echoed back to browsers.
func (s *Service) GetSMSConfig(ctx context.Context) (string, error) {
	row, err := s.store.GetCache(ctx, smsConfigKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return marshalSMSConfig(defaultSMSConfig()), nil
		}
		return "", err
	}
	config, err := parseSMSConfig(row.Result)
	if err != nil {
		return marshalSMSConfig(defaultSMSConfig()), nil
	}
	return marshalSMSConfig(config), nil
}

func (s *Service) SaveSMSConfig(ctx context.Context, raw string) (string, error) {
	config, err := parseSMSConfig(raw)
	if err != nil {
		return "", err
	}
	canonical := marshalSMSConfig(config)
	if err := s.store.UpsertCache(ctx, &Cache{Key: smsConfigKey, ExpireTime: 0, Result: canonical}); err != nil {
		return "", err
	}
	return canonical, nil
}

const shopConfigKey = "mall_shop_config"
const payConfigKey = "mall_pay_config"
const wechatAppConfigKey = "wechat_app_config"

type shopConfig struct {
	SiteName               string `json:"site_name"`
	SiteURL                string `json:"site_url"`
	OrderAutoCancelMinutes int    `json:"order_auto_cancel_minutes"`
	OrderAutoReceiveDays   int    `json:"order_auto_receive_days"`
	Enabled                bool   `json:"enabled"`
	Remark                 string `json:"remark"`
}

type payConfig struct {
	WechatEnabled  bool   `json:"wechat_enabled"`
	AlipayEnabled  bool   `json:"alipay_enabled"`
	BalanceEnabled bool   `json:"balance_enabled"`
	Remark         string `json:"remark"`
}

type wechatAppConfig struct {
	AppName string `json:"app_name"`
	Enabled bool   `json:"enabled"`
	Remark  string `json:"remark"`
}

func defaultShopConfig() shopConfig {
	return shopConfig{
		SiteName:               "七禧商城",
		SiteURL:                "",
		OrderAutoCancelMinutes: 30,
		OrderAutoReceiveDays:   7,
		Enabled:                true,
		Remark:                 "",
	}
}

func defaultPayConfig() payConfig {
	return payConfig{
		WechatEnabled:  false,
		AlipayEnabled:  false,
		BalanceEnabled: true,
		Remark:         "不含支付密钥；真实凭据请通过云服务配置或密钥管理维护",
	}
}

func defaultWechatAppConfig() wechatAppConfig {
	return wechatAppConfig{
		AppName: "七禧商城公众号",
		Enabled: false,
		Remark:  "不含 AppSecret、Token 或 EncodingAESKey；真实凭据请通过云服务配置维护",
	}
}

func marshalShopConfig(config shopConfig) string {
	data, _ := json.Marshal(config)
	return string(data)
}

func marshalPayConfig(config payConfig) string {
	data, _ := json.Marshal(config)
	return string(data)
}

func marshalWechatAppConfig(config wechatAppConfig) string {
	data, _ := json.Marshal(config)
	return string(data)
}

func parseShopConfig(raw string) (shopConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return shopConfig{}, ErrBadParam
	}
	for key := range fields {
		switch key {
		case "site_name", "site_url", "order_auto_cancel_minutes", "order_auto_receive_days", "enabled", "remark":
		default:
			return shopConfig{}, ErrBadParam
		}
	}
	var config shopConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return shopConfig{}, ErrBadParam
	}
	config.SiteName = strings.TrimSpace(config.SiteName)
	config.SiteURL = strings.TrimSpace(config.SiteURL)
	config.Remark = strings.TrimSpace(config.Remark)
	if config.SiteName == "" {
		config.SiteName = defaultShopConfig().SiteName
	}
	if config.OrderAutoCancelMinutes < 0 || config.OrderAutoReceiveDays < 0 ||
		len([]rune(config.SiteName)) > 128 || len([]rune(config.SiteURL)) > 256 || len([]rune(config.Remark)) > 500 {
		return shopConfig{}, ErrBadParam
	}
	return config, nil
}

func isSensitiveConfigKey(key string) bool {
	lower := strings.ToLower(strings.TrimSpace(key))
	for _, part := range []string{"secret", "key", "cert", "password", "token"} {
		if strings.Contains(lower, part) {
			return true
		}
	}
	return false
}

func parsePayConfig(raw string) (payConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return payConfig{}, ErrBadParam
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return payConfig{}, ErrBadParam
		}
		switch key {
		case "wechat_enabled", "alipay_enabled", "balance_enabled", "remark":
		default:
			return payConfig{}, ErrBadParam
		}
	}
	var config payConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil || len([]rune(strings.TrimSpace(config.Remark))) > 500 {
		return payConfig{}, ErrBadParam
	}
	config.Remark = strings.TrimSpace(config.Remark)
	if config.Remark == "" {
		config.Remark = defaultPayConfig().Remark
	}
	return config, nil
}

func parseWechatAppConfig(raw string) (wechatAppConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return wechatAppConfig{}, ErrBadParam
	}
	for key := range fields {
		if isSensitiveConfigKey(key) {
			return wechatAppConfig{}, ErrBadParam
		}
		switch key {
		case "app_name", "enabled", "remark":
		default:
			return wechatAppConfig{}, ErrBadParam
		}
	}
	var config wechatAppConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil || len([]rune(strings.TrimSpace(config.AppName))) > 64 || len([]rune(strings.TrimSpace(config.Remark))) > 500 {
		return wechatAppConfig{}, ErrBadParam
	}
	config.AppName = strings.TrimSpace(config.AppName)
	config.Remark = strings.TrimSpace(config.Remark)
	if config.AppName == "" {
		config.AppName = defaultWechatAppConfig().AppName
	}
	if config.Remark == "" {
		config.Remark = defaultWechatAppConfig().Remark
	}
	return config, nil
}

func (s *Service) GetShopConfig(ctx context.Context) (string, error) {
	row, err := s.store.GetCache(ctx, shopConfigKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return marshalShopConfig(defaultShopConfig()), nil
		}
		return "", err
	}
	config, err := parseShopConfig(row.Result)
	if err != nil {
		return marshalShopConfig(defaultShopConfig()), nil
	}
	return marshalShopConfig(config), nil
}

func (s *Service) SaveShopConfig(ctx context.Context, raw string) (string, error) {
	config, err := parseShopConfig(raw)
	if err != nil {
		return "", err
	}
	canonical := marshalShopConfig(config)
	if err := s.store.UpsertCache(ctx, &Cache{Key: shopConfigKey, ExpireTime: 0, Result: canonical}); err != nil {
		return "", err
	}
	return canonical, nil
}

func (s *Service) GetPayConfig(ctx context.Context) (string, error) {
	row, err := s.store.GetCache(ctx, payConfigKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return marshalPayConfig(defaultPayConfig()), nil
		}
		return "", err
	}
	config, err := parsePayConfig(row.Result)
	if err != nil {
		return marshalPayConfig(defaultPayConfig()), nil
	}
	return marshalPayConfig(config), nil
}

func (s *Service) SavePayConfig(ctx context.Context, raw string) (string, error) {
	config, err := parsePayConfig(raw)
	if err != nil {
		return "", err
	}
	canonical := marshalPayConfig(config)
	if err := s.store.UpsertCache(ctx, &Cache{Key: payConfigKey, ExpireTime: 0, Result: canonical}); err != nil {
		return "", err
	}
	return canonical, nil
}

func (s *Service) GetWechatAppConfig(ctx context.Context) (string, error) {
	row, err := s.store.GetCache(ctx, wechatAppConfigKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return marshalWechatAppConfig(defaultWechatAppConfig()), nil
		}
		return "", err
	}
	config, err := parseWechatAppConfig(row.Result)
	if err != nil {
		return marshalWechatAppConfig(defaultWechatAppConfig()), nil
	}
	return marshalWechatAppConfig(config), nil
}

func (s *Service) SaveWechatAppConfig(ctx context.Context, raw string) (string, error) {
	config, err := parseWechatAppConfig(raw)
	if err != nil {
		return "", err
	}
	canonical := marshalWechatAppConfig(config)
	if err := s.store.UpsertCache(ctx, &Cache{Key: wechatAppConfigKey, ExpireTime: 0, Result: canonical}); err != nil {
		return "", err
	}
	return canonical, nil
}

// marginConfigKey 对齐 CRMEB systemConfig：margin_remind_switch / margin_remind_day。
const marginConfigKey = "margin_remind_config"

type marginConfig struct {
	MarginRemindSwitch bool `json:"margin_remind_switch"`
	MarginRemindDay    int  `json:"margin_remind_day"`
}

func defaultMarginConfig() marginConfig {
	return marginConfig{
		MarginRemindSwitch: false,
		MarginRemindDay:    30,
	}
}

func marshalMarginConfig(config marginConfig) string {
	data, _ := json.Marshal(config)
	return string(data)
}

func parseMarginConfig(raw string) (marginConfig, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &fields); err != nil {
		return marginConfig{}, ErrBadParam
	}
	for key := range fields {
		switch key {
		case "margin_remind_switch", "margin_remind_day":
		default:
			return marginConfig{}, ErrBadParam
		}
	}
	var config marginConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return marginConfig{}, ErrBadParam
	}
	if config.MarginRemindDay < 0 || config.MarginRemindDay > 3650 {
		return marginConfig{}, ErrBadParam
	}
	return config, nil
}

func (s *Service) GetMarginConfig(ctx context.Context) (string, error) {
	row, err := s.store.GetCache(ctx, marginConfigKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return marshalMarginConfig(defaultMarginConfig()), nil
		}
		return "", err
	}
	config, err := parseMarginConfig(row.Result)
	if err != nil {
		return marshalMarginConfig(defaultMarginConfig()), nil
	}
	return marshalMarginConfig(config), nil
}

func (s *Service) SaveMarginConfig(ctx context.Context, raw string) (string, error) {
	config, err := parseMarginConfig(raw)
	if err != nil {
		return "", err
	}
	canonical := marshalMarginConfig(config)
	if err := s.store.UpsertCache(ctx, &Cache{Key: marginConfigKey, ExpireTime: 0, Result: canonical}); err != nil {
		return "", err
	}
	return canonical, nil
}

const (
	PriceDescriptionCacheKey = "product_price_desc"
	ActivityLabelCacheKey    = "product_activity_label"
)

func (s *Service) GetCacheList(ctx context.Context, key string) ([]CacheListItem, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrBadParam
	}
	row, err := s.store.GetCache(ctx, key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []CacheListItem{}, nil
		}
		return nil, err
	}
	if strings.TrimSpace(row.Result) == "" {
		return []CacheListItem{}, nil
	}
	var list []CacheListItem
	if err := json.Unmarshal([]byte(row.Result), &list); err != nil {
		return nil, ErrBadParam
	}
	return list, nil
}

func (s *Service) SaveCacheList(ctx context.Context, key string, items []CacheListItem) ([]CacheListItem, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrBadParam
	}
	clean := make([]CacheListItem, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			return nil, ErrBadParam
		}
		id := strings.TrimSpace(item.ID)
		if id == "" {
			id = name
		}
		if _, ok := seen[id]; ok {
			return nil, ErrBadParam
		}
		seen[id] = struct{}{}
		clean = append(clean, CacheListItem{
			ID: id, Name: name, Enabled: item.Enabled,
			Remark: strings.TrimSpace(item.Remark),
		})
	}
	raw, err := json.Marshal(clean)
	if err != nil {
		return nil, err
	}
	if err := s.store.UpsertCache(ctx, &Cache{Key: key, ExpireTime: 0, Result: string(raw)}); err != nil {
		return nil, err
	}
	return clean, nil
}

func normalize(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}
