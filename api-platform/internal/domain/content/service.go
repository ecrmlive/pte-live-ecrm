package content

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
)

type Store interface {
	ListNotices(ctx context.Context, onlyShow bool, filter NoticeListFilter) ([]Notice, int64, error)
	GetNotice(ctx context.Context, id uint) (*Notice, error)
	CreateNotice(ctx context.Context, n *Notice, scopeIDs []uint) error
	UpdateNotice(ctx context.Context, n *Notice, scopeIDs []uint) error
	ListNoticeScopes(ctx context.Context, noticeIDs []uint) ([]NoticeScope, error)
	UpdateNoticeStatus(ctx context.Context, id uint, isShow int8) error
	SoftDeleteNotice(ctx context.Context, id uint) error
	GetCache(ctx context.Context, key string) (*Cache, error)
	UpsertCache(ctx context.Context, row *Cache) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListApp(ctx context.Context, page, limit int) (*PageResult[Notice], error) {
	return s.list(ctx, true, NoticeListFilter{Page: page, Limit: limit})
}

func (s *Service) ListAdmin(ctx context.Context, filter NoticeListFilter) (*PageResult[Notice], error) {
	return s.list(ctx, false, filter)
}

func (s *Service) list(ctx context.Context, onlyShow bool, filter NoticeListFilter) (*PageResult[Notice], error) {
	filter.Page, filter.Limit = normalize(filter.Page, filter.Limit)
	list, total, err := s.store.ListNotices(ctx, onlyShow, filter)
	if err != nil {
		return nil, err
	}
	if err := s.hydrateScopes(ctx, list); err != nil {
		return nil, err
	}
	return &PageResult[Notice]{List: list, Total: total, Page: filter.Page, Limit: filter.Limit}, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*Notice, error) {
	n, err := s.store.GetNotice(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if err := s.hydrateNotice(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) Create(ctx context.Context, in NoticeInput) (*Notice, error) {
	title, content, scopeType, scopeIDs, err := normalizeNoticeInput(in)
	if err != nil {
		return nil, ErrBadParam
	}
	n := &Notice{Title: title, Content: content, IsShow: 1, ScopeType: scopeType}
	if in.IsShow != nil {
		n.IsShow = *in.IsShow
	}
	if err := s.store.CreateNotice(ctx, n, scopeIDs); err != nil {
		return nil, err
	}
	n.ScopeIDs = scopeIDs
	if err := s.hydrateNotice(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) Update(ctx context.Context, id uint, in NoticeInput) (*Notice, error) {
	n, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	title, content, scopeType, scopeIDs, validateErr := normalizeNoticeInput(in)
	if validateErr != nil {
		return nil, ErrBadParam
	}
	n.Title = title
	n.Content = content
	n.ScopeType = scopeType
	if in.IsShow != nil {
		n.IsShow = *in.IsShow
	}
	if err := s.store.UpdateNotice(ctx, n, scopeIDs); err != nil {
		return nil, err
	}
	n.ScopeIDs = scopeIDs
	if err := s.hydrateNotice(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id uint, in NoticeStatusInput) error {
	if in.IsShow != 0 && in.IsShow != 1 {
		return ErrBadParam
	}
	if _, err := s.store.GetNotice(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.UpdateNoticeStatus(ctx, id, in.IsShow)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.store.SoftDeleteNotice(ctx, id)
}

func normalizeNoticeInput(in NoticeInput) (string, string, NoticeScopeType, []uint, error) {
	title := strings.TrimSpace(in.Title)
	content := strings.TrimSpace(in.Content)
	if title == "" || utf8.RuneCountInString(title) > 20 || content == "" {
		return "", "", "", nil, ErrBadParam
	}
	scopeType := in.ScopeType
	if scopeType == "" {
		scopeType = NoticeScopeAll
	}
	if scopeType != NoticeScopeAll && scopeType != NoticeScopeStoreName && scopeType != NoticeScopeStoreType && scopeType != NoticeScopeStoreCategory {
		return "", "", "", nil, ErrBadParam
	}
	if scopeType == NoticeScopeAll {
		return title, content, scopeType, nil, nil
	}
	seen := make(map[uint]struct{}, len(in.ScopeIDs))
	ids := make([]uint, 0, len(in.ScopeIDs))
	for _, id := range in.ScopeIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 || len(ids) > 200 {
		return "", "", "", nil, ErrBadParam
	}
	return title, content, scopeType, ids, nil
}

func (s *Service) hydrateScopes(ctx context.Context, notices []Notice) error {
	ids := make([]uint, 0, len(notices))
	for _, item := range notices {
		if item.NoticeID > 0 && item.ScopeType != NoticeScopeAll {
			ids = append(ids, item.NoticeID)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	scopes, err := s.store.ListNoticeScopes(ctx, ids)
	if err != nil {
		return err
	}
	byNotice := make(map[uint][]NoticeScope)
	for _, scope := range scopes {
		byNotice[scope.NoticeID] = append(byNotice[scope.NoticeID], scope)
	}
	for i := range notices {
		notices[i].ScopeItems = byNotice[notices[i].NoticeID]
		notices[i].ScopeIDs = make([]uint, 0, len(notices[i].ScopeItems))
		for _, scope := range notices[i].ScopeItems {
			notices[i].ScopeIDs = append(notices[i].ScopeIDs, scope.ScopeID)
		}
	}
	return nil
}

func (s *Service) hydrateNotice(ctx context.Context, notice *Notice) error {
	list := []Notice{*notice}
	if err := s.hydrateScopes(ctx, list); err != nil {
		return err
	}
	*notice = list[0]
	return nil
}

// AgreeCatalog 是所有已被平台业务使用的协议键。保留既有键，避免其它协议维护页失效。
func AgreeCatalog() []AgreeMeta {
	return []AgreeMeta{
		{Key: "sys_user_agree", Label: "用户协议"},
		{Key: "sys_userr_privacy", Label: "隐私政策"},
		{Key: "the_cancellation_prompt", Label: "注销提示"},
		{Key: "platform_rule", Label: "平台规则"},
		{Key: "sys_intention_agree", Label: "店铺入驻申请协议"},
		{Key: "circle_entry_agree", Label: "代理入驻申请协议"},
		{Key: "the_cancellation_msg", Label: "注销声明"},
		{Key: "sys_certificate", Label: "资质证照"},
		{Key: "sys_svip", Label: "付费会员协议"},
		{Key: "sys_receipt_agree", Label: "发票说明"},
		{Key: "sys_product_presell_agree", Label: "预售协议"},
		{Key: "business_entry_agree", Label: "商户入驻申请协议"},
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

// AgreementSettingsCatalog 是“协议设置”页面固定展示的产品协议。顺序与平台设置菜单一致。
func AgreementSettingsCatalog() []AgreeMeta {
	return []AgreeMeta{
		{Key: "sys_user_agree", Label: "用户协议"},
		{Key: "sys_userr_privacy", Label: "隐私政策"},
		{Key: "the_cancellation_prompt", Label: "注销提示"},
		{Key: "platform_rule", Label: "平台规则"},
		{Key: "sys_intention_agree", Label: "店铺入驻申请协议"},
		{Key: "circle_entry_agree", Label: "代理入驻申请协议"},
		{Key: "business_entry_agree", Label: "商户入驻申请协议"},
		{Key: "the_cancellation_msg", Label: "注销声明"},
		{Key: "sys_about_us", Label: "关于我们"},
		{Key: "sys_certificate", Label: "资质证照"},
	}
}

func (s *Service) ListAgreements(ctx context.Context) ([]AgreeView, error) {
	catalog := AgreementSettingsCatalog()
	out := make([]AgreeView, 0, len(catalog))
	for _, m := range catalog {
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

const shopConfigKey = "mall_shop_config"
const payConfigKey = "mall_pay_config"
const wechatAppConfigKey = "wechat_app_config"

type shopConfig struct {
	// 保留历史字段，供店铺入驻邀请链接读取；商城设置页面不展示它们。
	SiteName string `json:"site_name,omitempty"`
	SiteURL  string `json:"site_url,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
	Remark   string `json:"remark,omitempty"`

	AutoParseClipboard        bool `json:"auto_parse_clipboard"`
	ArrivalNoticeEnabled      bool `json:"arrival_notice_enabled"`
	ProductCommentEnabled     bool `json:"product_comment_enabled"`
	AutoPositiveReviewEnabled bool `json:"auto_positive_review_enabled"`
	DefaultCopyTimes          int  `json:"default_copy_times"`

	OrderAutoCancelMinutes int      `json:"order_auto_cancel_minutes"`
	OrderAutoReceiveDays   int      `json:"order_auto_receive_days"`
	AfterSaleDays          int      `json:"after_sale_days"`
	MerchantRefundAutoDays int      `json:"merchant_refund_auto_days"`
	RefundReasons          []string `json:"refund_reasons"`
	PlatformRightsEnabled  bool     `json:"platform_rights_enabled"`
	PlatformRightsDays     int      `json:"platform_rights_days"`
	MergePaymentEnabled    bool     `json:"merge_payment_enabled"`

	MerchantApplyEnabled          bool   `json:"merchant_apply_enabled"`
	MerchantQualificationRequired bool   `json:"merchant_qualification_required"`
	MerchantMarginBadgeEnabled    bool   `json:"merchant_margin_badge_enabled"`
	MerchantMarginBadgeImage      string `json:"merchant_margin_badge_image"`
	MerchantCategoryLimit         int    `json:"merchant_category_limit"`

	MallShowStores               bool   `json:"mall_show_stores"`
	MallRecommendEnabled         bool   `json:"mall_recommend_enabled"`
	MallRecommendDistanceEnabled bool   `json:"mall_recommend_distance_enabled"`
	MallRecommendSort            string `json:"mall_recommend_sort"`
	LiveStreamAutoApprove        bool   `json:"live_stream_auto_approve"`
	LiveProductAutoApprove       bool   `json:"live_product_auto_approve"`
	HotRankingEnabled            bool   `json:"hot_ranking_enabled"`
	HotRankingCategoryLevel      int    `json:"hot_ranking_category_level"`
	HotRankingRefreshHours       int    `json:"hot_ranking_refresh_hours"`
	MallSearchMode               string `json:"mall_search_mode"`

	ProductRankingPeriod string `json:"product_ranking_period"`
	ProductRankingMetric string `json:"product_ranking_metric"`
	ShopRankingPeriod    string `json:"shop_ranking_period"`
	ShopRankingMetric    string `json:"shop_ranking_metric"`
	DashboardDisplayName string `json:"dashboard_display_name"`
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
		SiteName:                      "七禧商城",
		SiteURL:                       "",
		Enabled:                       true,
		Remark:                        "",
		AutoParseClipboard:            true,
		ArrivalNoticeEnabled:          true,
		ProductCommentEnabled:         true,
		AutoPositiveReviewEnabled:     true,
		DefaultCopyTimes:              8,
		OrderAutoCancelMinutes:        15,
		OrderAutoReceiveDays:          7,
		AfterSaleDays:                 1,
		MerchantRefundAutoDays:        1,
		RefundReasons:                 []string{"商品质量问题", "不想要了", "未收到货"},
		PlatformRightsEnabled:         true,
		PlatformRightsDays:            1,
		MergePaymentEnabled:           true,
		MerchantApplyEnabled:          true,
		MerchantQualificationRequired: true,
		MerchantMarginBadgeEnabled:    false,
		MerchantMarginBadgeImage:      "",
		MerchantCategoryLimit:         5,
		MallShowStores:                true,
		MallRecommendEnabled:          true,
		MallRecommendDistanceEnabled:  true,
		MallRecommendSort:             "star",
		LiveStreamAutoApprove:         false,
		LiveProductAutoApprove:        false,
		HotRankingEnabled:             true,
		HotRankingCategoryLevel:       2,
		HotRankingRefreshHours:        24,
		MallSearchMode:                "fuzzy",
		ProductRankingPeriod:          "month",
		ProductRankingMetric:          "sales_amount",
		ShopRankingPeriod:             "month",
		ShopRankingMetric:             "product_count",
		DashboardDisplayName:          "数据大屏",
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
		if isSensitiveConfigKey(key) {
			return shopConfig{}, ErrBadParam
		}
		switch key {
		case "site_name", "site_url", "enabled", "remark",
			"auto_parse_clipboard", "arrival_notice_enabled", "product_comment_enabled", "auto_positive_review_enabled", "default_copy_times",
			"order_auto_cancel_minutes", "order_auto_receive_days", "after_sale_days", "merchant_refund_auto_days", "refund_reasons", "platform_rights_enabled", "platform_rights_days", "merge_payment_enabled",
			"merchant_apply_enabled", "merchant_qualification_required", "merchant_margin_badge_enabled", "merchant_margin_badge_image", "merchant_category_limit",
			"mall_show_stores", "mall_recommend_enabled", "mall_recommend_distance_enabled", "mall_recommend_sort", "live_stream_auto_approve", "live_product_auto_approve", "hot_ranking_enabled", "hot_ranking_category_level", "hot_ranking_refresh_hours", "mall_search_mode",
			"product_ranking_period", "product_ranking_metric", "shop_ranking_period", "shop_ranking_metric", "dashboard_display_name":
		default:
			return shopConfig{}, ErrBadParam
		}
	}
	config := defaultShopConfig()
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return shopConfig{}, ErrBadParam
	}
	config.SiteName = strings.TrimSpace(config.SiteName)
	config.SiteURL = strings.TrimSpace(config.SiteURL)
	config.Remark = strings.TrimSpace(config.Remark)
	if config.SiteName == "" {
		config.SiteName = defaultShopConfig().SiteName
	}
	config.MerchantMarginBadgeImage = strings.TrimSpace(config.MerchantMarginBadgeImage)
	config.DashboardDisplayName = strings.TrimSpace(config.DashboardDisplayName)
	if config.DashboardDisplayName == "" {
		config.DashboardDisplayName = defaultShopConfig().DashboardDisplayName
	}
	if config.OrderAutoCancelMinutes < 0 || config.OrderAutoCancelMinutes > 10080 ||
		config.OrderAutoReceiveDays < 0 || config.OrderAutoReceiveDays > 365 ||
		config.AfterSaleDays < 0 || config.AfterSaleDays > 365 ||
		config.MerchantRefundAutoDays < 0 || config.MerchantRefundAutoDays > 365 ||
		config.PlatformRightsDays < 0 || config.PlatformRightsDays > 365 ||
		config.DefaultCopyTimes < 0 || config.DefaultCopyTimes > 1000000 ||
		config.MerchantCategoryLimit < 0 || config.MerchantCategoryLimit > 10000 ||
		config.HotRankingRefreshHours < 1 || config.HotRankingRefreshHours > 720 ||
		config.HotRankingCategoryLevel < 1 || config.HotRankingCategoryLevel > 3 ||
		len([]rune(config.SiteName)) > 128 || len([]rune(config.SiteURL)) > 256 || len([]rune(config.Remark)) > 500 ||
		len([]rune(config.MerchantMarginBadgeImage)) > 512 || len([]rune(config.DashboardDisplayName)) > 64 ||
		!oneOf(config.MallRecommendSort, "default", "star", "created_at") ||
		!oneOf(config.MallSearchMode, "fuzzy", "split") ||
		!oneOf(config.ProductRankingPeriod, "today", "week", "month") ||
		!oneOf(config.ProductRankingMetric, "sales_quantity", "sales_amount") ||
		!oneOf(config.ShopRankingPeriod, "today", "week", "month") ||
		!oneOf(config.ShopRankingMetric, "sales_amount", "product_count") {
		return shopConfig{}, ErrBadParam
	}
	config.RefundReasons = normalizeRefundReasons(config.RefundReasons)
	if len(config.RefundReasons) == 0 || len(config.RefundReasons) > 50 {
		return shopConfig{}, ErrBadParam
	}
	return config, nil
}

func oneOf(value string, options ...string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}

func normalizeRefundReasons(reasons []string) []string {
	result := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		reason = strings.TrimSpace(reason)
		if reason == "" || len([]rune(reason)) > 100 {
			continue
		}
		result = append(result, reason)
	}
	return result
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
