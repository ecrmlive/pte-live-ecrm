package merchant

import (
	"context"
	"errors"
	"math"
	"strings"

	"gorm.io/gorm"
)

type Store interface {
	ListMerchants(ctx context.Context, keyword string, status *int8, scope MerchantScope, page, limit int) ([]Merchant, int64, error)
	GetMerchant(ctx context.Context, id uint) (*Merchant, error)
	UpdateMerchantStatus(ctx context.Context, id uint, status, merState int8) error
	UpdateSvipCouponMerge(ctx context.Context, merID uint, merge int8) error
	UpdateShopProfile(ctx context.Context, merID uint, merName, realName, merPhone, merAddress, merInfo string) error
	UpsertMerchantView(ctx context.Context, m *Merchant) error

	ListIntentions(ctx context.Context, keyword string, status *int8, regionIDs []uint, page, limit int) ([]Intention, int64, error)
	GetIntention(ctx context.Context, id uint, regionIDs []uint) (*Intention, error)
	SaveIntention(ctx context.Context, row *Intention) error
	AssignIntentionRegion(ctx context.Context, id, regionID uint) (bool, error)

	ListCategories(ctx context.Context) ([]Category, error)
	CreateCategory(ctx context.Context, c *Category) error
	UpdateCategory(ctx context.Context, c *Category) error
	DeleteCategory(ctx context.Context, id uint) error

	WithTx(fn func(tx Store) error) error
}

type Service struct {
	store Store
}

// MerchantScope is resolved from qixi_crm_a_data_scope. A nil pair denotes
// full-platform supervision; non-nil empty slices intentionally deny access.
type MerchantScope struct {
	MerchantIDs []uint
	RegionIDs   []uint
}

func NewService(store Store) *Service { return &Service{store: store} }

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

func (s *Service) ListMerchants(ctx context.Context, keyword string, status *int8, scope MerchantScope, page, limit int) (*PageResult[Merchant], error) {
	list, total, err := s.store.ListMerchants(ctx, keyword, status, scope, page, limit)
	if err != nil {
		return nil, err
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &PageResult[Merchant]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) RequireMerchantScope(ctx context.Context, id uint, scope MerchantScope) error {
	if scope.MerchantIDs == nil && scope.RegionIDs == nil {
		return nil
	}
	if len(scope.MerchantIDs) == 0 && len(scope.RegionIDs) == 0 {
		return ErrNotFound
	}
	m, err := s.GetMerchant(ctx, id)
	if err != nil {
		return err
	}
	for _, merchantID := range scope.MerchantIDs {
		if m.MerID == merchantID {
			return nil
		}
	}
	for _, regionID := range scope.RegionIDs {
		if m.RegionID == regionID {
			return nil
		}
	}
	return ErrNotFound
}

func (s *Service) GetMerchant(ctx context.Context, id uint) (*Merchant, error) {
	m, err := s.store.GetMerchant(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return m, nil
}

func (s *Service) SetMerchantEnabled(ctx context.Context, id uint, enabled bool) error {
	if _, err := s.GetMerchant(ctx, id); err != nil {
		return err
	}
	status, state := int8(0), int8(0)
	if enabled {
		status, state = 1, 1
	}
	return s.store.UpdateMerchantStatus(ctx, id, status, state)
}

type ShopProfileInput struct {
	MerName    string `json:"mer_name"`
	RealName   string `json:"real_name"`
	MerPhone   string `json:"mer_phone"`
	MerAddress string `json:"mer_address"`
	MerInfo    string `json:"mer_info"`
}

func (s *Service) UpdateShopProfile(ctx context.Context, merID uint, in ShopProfileInput) (*Merchant, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	m, err := s.GetMerchant(ctx, merID)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(in.MerName)
	if name == "" {
		name = m.MerName
	}
	realName := strings.TrimSpace(in.RealName)
	if realName == "" {
		realName = m.RealName
	}
	phone := strings.TrimSpace(in.MerPhone)
	if phone == "" {
		phone = m.MerPhone
	}
	addr := strings.TrimSpace(in.MerAddress)
	if addr == "" {
		addr = m.MerAddress
	}
	info := strings.TrimSpace(in.MerInfo)
	if info == "" {
		info = m.MerInfo
	}
	if err := s.store.UpdateShopProfile(ctx, merID, name, realName, phone, addr, info); err != nil {
		return nil, err
	}
	return s.GetMerchant(ctx, merID)
}

func (s *Service) GetSvipConfig(ctx context.Context, merID uint) (*SvipConfig, error) {
	m, err := s.GetMerchant(ctx, merID)
	if err != nil {
		return nil, err
	}
	return &SvipConfig{MerID: m.MerID, SvipCouponMerge: m.SvipCouponMerge}, nil
}

func (s *Service) UpdateSvipConfig(ctx context.Context, merID uint, merge int8) (*SvipConfig, error) {
	if merge != 0 && merge != 1 {
		return nil, ErrBadParam
	}
	if _, err := s.GetMerchant(ctx, merID); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSvipCouponMerge(ctx, merID, merge); err != nil {
		return nil, err
	}
	return s.GetSvipConfig(ctx, merID)
}

func (s *Service) ListIntentions(ctx context.Context, keyword string, status *int8, regionIDs []uint, page, limit int) (*PageResult[Intention], error) {
	list, total, err := s.store.ListIntentions(ctx, keyword, status, regionIDs, page, limit)
	if err != nil {
		return nil, err
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &PageResult[Intention]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetIntention(ctx context.Context, id uint, regionIDs []uint) (*Intention, error) {
	row, err := s.store.GetIntention(ctx, id, regionIDs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

type AuditIntentionInput struct {
	Status   int8   `json:"status"`
	FailMsg  string `json:"fail_msg"`
	Mark     string `json:"mark"`
	Account  string `json:"account"`
	Password string `json:"password"`
	RegionID uint   `json:"region_id"`
}

type AuditIntentionResult struct {
	Intention *Intention `json:"intention"`
	MerID     uint       `json:"mer_id,omitempty"`
	Account   string     `json:"account,omitempty"`
}

func (s *Service) AssignIntentionRegion(ctx context.Context, id, regionID uint) (*Intention, error) {
	if id == 0 || regionID == 0 {
		return nil, ErrBadParam
	}
	row, err := s.GetIntention(ctx, id, nil)
	if err != nil {
		return nil, err
	}
	if row.Status != IntentionPending {
		return nil, ErrAlreadyAudited
	}
	updated, err := s.store.AssignIntentionRegion(ctx, id, regionID)
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, ErrAlreadyAudited
	}
	row.CircleID = regionID
	return row, nil
}

// FinalizeIntentionApproval records a store that has already been created by
// api-merchant's idempotent onboarding command. Keeping the command outside
// this transaction preserves the three-database ownership boundary.
func (s *Service) FinalizeIntentionApproval(ctx context.Context, id uint, in AuditIntentionInput, merchantID uint, regionIDs []uint) (*AuditIntentionResult, error) {
	if merchantID == 0 || in.RegionID == 0 || strings.TrimSpace(in.Account) == "" {
		return nil, ErrBadParam
	}
	out := AuditIntentionResult{MerID: merchantID, Account: strings.TrimSpace(in.Account)}
	err := s.store.WithTx(func(tx Store) error {
		row, err := tx.GetIntention(ctx, id, regionIDs)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if row.Status != IntentionPending {
			return ErrAlreadyAudited
		}
		m := &Merchant{MerID: merchantID, CategoryID: row.MerchantCategoryID, RegionID: in.RegionID, MerName: row.MerName, RealName: row.Name, MerPhone: row.Phone, Mark: "入驻审核通过", Status: 1, MerState: 1, IsAudit: 1}
		if err := tx.UpsertMerchantView(ctx, m); err != nil {
			return err
		}
		row.Status, row.MerID, row.CircleID, row.Mark, row.FailMsg = IntentionApproved, merchantID, in.RegionID, strings.TrimSpace(in.Mark), ""
		if err := tx.SaveIntention(ctx, row); err != nil {
			return err
		}
		out.Intention = row
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *Service) AuditIntention(ctx context.Context, id uint, in AuditIntentionInput, regionIDs []uint) (*AuditIntentionResult, error) {
	if in.Status != IntentionApproved && in.Status != IntentionRejected {
		return nil, ErrBadStatus
	}
	if in.Status == IntentionRejected && strings.TrimSpace(in.FailMsg) == "" {
		return nil, ErrRejectNeedMsg
	}

	var out AuditIntentionResult
	err := s.store.WithTx(func(tx Store) error {
		row, err := tx.GetIntention(ctx, id, regionIDs)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if row.Status != IntentionPending {
			return ErrAlreadyAudited
		}
		row.Mark = strings.TrimSpace(in.Mark)
		if in.Status == IntentionRejected {
			row.Status = IntentionRejected
			row.FailMsg = strings.TrimSpace(in.FailMsg)
			if err := tx.SaveIntention(ctx, row); err != nil {
				return err
			}
			out.Intention = row
			return nil
		}

		return ErrProvisionRequired
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *Service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.store.ListCategories(ctx)
}

func (s *Service) CreateCategory(ctx context.Context, name string, rate float64) (*Category, error) {
	name = strings.TrimSpace(name)
	if name == "" || !validCommissionRate(rate) {
		return nil, ErrBadParam
	}
	c := &Category{CategoryName: name, CommissionRate: rate}
	if err := s.store.CreateCategory(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id uint, name string, rate float64) error {
	name = strings.TrimSpace(name)
	if id == 0 || name == "" || !validCommissionRate(rate) {
		return ErrBadParam
	}
	return s.store.UpdateCategory(ctx, &Category{
		MerchantCategoryID: id,
		CategoryName:       name,
		CommissionRate:     rate,
	})
}

func (s *Service) DeleteCategory(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrBadParam
	}
	return s.store.DeleteCategory(ctx, id)
}

func validCommissionRate(rate float64) bool {
	if math.IsNaN(rate) || math.IsInf(rate, 0) || rate < 0 || rate > 100 {
		return false
	}
	return math.Abs(rate-math.Round(rate*100)/100) < 0.000001
}
