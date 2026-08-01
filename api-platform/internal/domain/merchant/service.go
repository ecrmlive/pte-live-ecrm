package merchant

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Store interface {
	ListMerchants(ctx context.Context, keyword string, status *int8, regionIDs []uint, page, limit int) ([]Merchant, int64, error)
	GetMerchant(ctx context.Context, id uint) (*Merchant, error)
	UpdateMerchantStatus(ctx context.Context, id uint, status, merState int8) error
	UpdateSvipCouponMerge(ctx context.Context, merID uint, merge int8) error
	UpdateShopProfile(ctx context.Context, merID uint, merName, realName, merPhone, merAddress, merInfo string) error
	CreateMerchant(ctx context.Context, m *Merchant) error
	CreateMerchantAdmin(ctx context.Context, a *identity.MerchantAdmin) error

	ListIntentions(ctx context.Context, keyword string, status *int8, page, limit int) ([]Intention, int64, error)
	GetIntention(ctx context.Context, id uint) (*Intention, error)
	SaveIntention(ctx context.Context, row *Intention) error

	ListCategories(ctx context.Context) ([]Category, error)
	CreateCategory(ctx context.Context, c *Category) error
	UpdateCategory(ctx context.Context, c *Category) error
	DeleteCategory(ctx context.Context, id uint) error

	WithTx(fn func(tx Store) error) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

func (s *Service) ListMerchants(ctx context.Context, keyword string, status *int8, regionIDs []uint, page, limit int) (*PageResult[Merchant], error) {
	list, total, err := s.store.ListMerchants(ctx, keyword, status, regionIDs, page, limit)
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

func (s *Service) RequireMerchantScope(ctx context.Context, id uint, regionIDs []uint) error {
	if regionIDs == nil {
		return nil
	}
	if len(regionIDs) == 0 {
		return ErrNotFound
	}
	m, err := s.GetMerchant(ctx, id)
	if err != nil {
		return err
	}
	for _, regionID := range regionIDs {
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

func (s *Service) ListIntentions(ctx context.Context, keyword string, status *int8, page, limit int) (*PageResult[Intention], error) {
	list, total, err := s.store.ListIntentions(ctx, keyword, status, page, limit)
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

func (s *Service) GetIntention(ctx context.Context, id uint) (*Intention, error) {
	row, err := s.store.GetIntention(ctx, id)
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

func (s *Service) AuditIntention(ctx context.Context, id uint, in AuditIntentionInput) (*AuditIntentionResult, error) {
	if in.Status != IntentionApproved && in.Status != IntentionRejected {
		return nil, ErrBadStatus
	}
	if in.Status == IntentionRejected && strings.TrimSpace(in.FailMsg) == "" {
		return nil, ErrRejectNeedMsg
	}

	var out AuditIntentionResult
	err := s.store.WithTx(func(tx Store) error {
		row, err := tx.GetIntention(ctx, id)
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

		m := &Merchant{
			CategoryID: row.MerchantCategoryID,
			RegionID:   in.RegionID,
			MerName:    row.MerName,
			RealName:   row.Name,
			MerPhone:   row.Phone,
			Mark:       "入驻审核通过",
			Status:     1,
			MerState:   1,
			IsAudit:    1,
		}
		if err := tx.CreateMerchant(ctx, m); err != nil {
			return err
		}
		account := strings.TrimSpace(in.Account)
		if account == "" {
			account = fmt.Sprintf("mer%d", m.MerID)
		}
		pwd := strings.TrimSpace(in.Password)
		if pwd == "" {
			pwd = "admin123"
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		admin := &identity.MerchantAdmin{
			MerID:    m.MerID,
			Account:  account,
			Pwd:      string(hash),
			RealName: row.Name,
			Phone:    row.Phone,
			Roles:    "2",
			Level:    0,
			Status:   1,
		}
		if err := tx.CreateMerchantAdmin(ctx, admin); err != nil {
			return err
		}
		row.Status = IntentionApproved
		row.MerID = m.MerID
		row.FailMsg = ""
		if err := tx.SaveIntention(ctx, row); err != nil {
			return err
		}
		out.Intention = row
		out.MerID = m.MerID
		out.Account = account
		return nil
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
	if name == "" {
		return nil, errors.New("分类名称不能为空")
	}
	c := &Category{CategoryName: name, CommissionRate: rate}
	if err := s.store.CreateCategory(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id uint, name string, rate float64) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("分类名称不能为空")
	}
	return s.store.UpdateCategory(ctx, &Category{
		MerchantCategoryID: id,
		CategoryName:       name,
		CommissionRate:     rate,
	})
}

func (s *Service) DeleteCategory(ctx context.Context, id uint) error {
	return s.store.DeleteCategory(ctx, id)
}
