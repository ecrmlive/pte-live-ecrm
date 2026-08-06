package content

import (
	"context"
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
		{Key: "sys_about_us", Label: "关于我们"},
		{Key: "sys_refund_agree", Label: "退款协议"},
		{Key: "sys_cancel_agree", Label: "取消订单说明"},
		{Key: "sys_recharge_agree", Label: "充值协议"},
		{Key: "sys_integral_agree", Label: "积分规则"},
		{Key: "mer_settle_agree", Label: "商户结算说明"},
		{Key: "sys_lottery_agree", Label: "抽奖活动说明"},
		{Key: "sys_deposit_agree", Label: "保证金说明"},
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

// GetSMSConfig 短信通道配置 stub（存 qixi_m_admin_cache，未发真实短信）。
func (s *Service) GetSMSConfig(ctx context.Context) (string, error) {
	row, err := s.store.GetCache(ctx, smsConfigKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return `{"enabled":false,"provider":"stub","sign":"七禧商城","remark":"未配置"}`, nil
		}
		return "", err
	}
	return row.Result, nil
}

func (s *Service) SaveSMSConfig(ctx context.Context, raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ErrBadParam
	}
	if err := s.store.UpsertCache(ctx, &Cache{Key: smsConfigKey, ExpireTime: 0, Result: raw}); err != nil {
		return "", err
	}
	return raw, nil
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
