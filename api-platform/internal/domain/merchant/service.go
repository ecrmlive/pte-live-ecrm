package merchant

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Store interface {
	ListMerchants(ctx context.Context, filter ListFilter, scope MerchantScope) ([]Merchant, int64, error)
	GetMerchant(ctx context.Context, id uint) (*Merchant, error)
	UpdateMerchantStatus(ctx context.Context, id uint, status, merState int8) error
	UpdateSvipCouponMerge(ctx context.Context, merID uint, merge int8) error
	UpdateMerchant(ctx context.Context, m *Merchant) error
	CreateMerchant(ctx context.Context, m *Merchant) error
	UpsertMerchantView(ctx context.Context, m *Merchant) error

	ListIntentions(ctx context.Context, keyword string, status *int8, regionIDs []uint, page, limit int, dateFrom, dateTo string, categoryID, typeID *uint) ([]Intention, int64, error)
	GetIntention(ctx context.Context, id uint, regionIDs []uint) (*Intention, error)
	SaveIntention(ctx context.Context, row *Intention) error
	DeleteIntention(ctx context.Context, id uint, regionIDs []uint) (bool, error)
	AssignIntentionRegion(ctx context.Context, id, regionID uint) (bool, error)

	ListCategories(ctx context.Context) ([]Category, error)
	CreateCategory(ctx context.Context, c *Category) error
	UpdateCategory(ctx context.Context, c *Category) error
	DeleteCategory(ctx context.Context, id uint) error

	WithTx(fn func(tx Store) error) error
}

// ListFilter 对齐 CRMEB 店铺列表筛选项。
type ListFilter struct {
	Keyword    string
	Status     *int8
	CategoryID *uint
	TypeID     *uint
	RegionID   *uint
	IsBest     *int8
	OfflinePay *int8
	DateFrom   string
	DateTo     string
	Page       int
	Limit      int
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

func (s *Service) ListMerchants(ctx context.Context, filter ListFilter, scope MerchantScope) (*PageResult[Merchant], error) {
	list, total, err := s.store.ListMerchants(ctx, filter, scope)
	if err != nil {
		return nil, err
	}
	page, limit := filter.Page, filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
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
	MerName          string  `json:"mer_name"`
	OwnerName        string  `json:"owner_name"`
	RealName         string  `json:"real_name"`
	MerPhone         string  `json:"mer_phone"`
	MerAddress       string  `json:"mer_address"`
	MerInfo          string  `json:"mer_info"`
	MerKeyword       string  `json:"mer_keyword"`
	Mark             string  `json:"mark"`
	CategoryID       uint    `json:"category_id"`
	TypeID           uint    `json:"type_id"`
	BusinessID       uint    `json:"business_id"`
	RegionID         uint    `json:"region_id"`
	IsBest           *bool   `json:"is_best"`
	OfflinePay       *bool   `json:"offline_pay"`
	IsTrader         *bool   `json:"is_trader"`
	IsAudit          *bool   `json:"is_audit"`
	IsBroRoom        *bool   `json:"is_bro_room"`
	IsBroGoods       *bool   `json:"is_bro_goods"`
	CommissionSwitch *bool   `json:"commission_switch"`
	CommissionRate   *float64 `json:"commission_rate"`
	MerAccount       string  `json:"mer_account"`
	MerPassword      string  `json:"mer_password"`
	SubMchid         string  `json:"sub_mchid"`
	ApplymentID      string  `json:"applyment_id"`
	CareCount        *int    `json:"care_count"`
	CareFicti        *int    `json:"care_ficti"`
	GoodsTypes       []int   `json:"goods_types"`
	PlatformCategoryIDs []uint `json:"platform_category_ids"`
	MerStar          *int8   `json:"mer_star"`
	MerAvatar        string  `json:"mer_avatar"`
	Sort             *int    `json:"sort"`
	Status           *bool   `json:"status"`
	StoreGroupIDs    []uint  `json:"store_group_ids"`
}

func (s *Service) UpdateShopProfile(ctx context.Context, merID uint, in ShopProfileInput) (*Merchant, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	if strings.TrimSpace(in.MerAvatar) == "" {
		return nil, ErrBadParam
	}
	m, err := s.GetMerchant(ctx, merID)
	if err != nil {
		return nil, err
	}
	applyProfile(m, in)
	if err := s.store.UpdateMerchant(ctx, m); err != nil {
		return nil, err
	}
	return s.GetMerchant(ctx, merID)
}

func (s *Service) CreateMerchant(ctx context.Context, in ShopProfileInput) (*Merchant, error) {
	name := strings.TrimSpace(in.MerName)
	if name == "" {
		return nil, ErrBadParam
	}
	if strings.TrimSpace(in.MerAvatar) == "" {
		return nil, ErrBadParam
	}
	m := &Merchant{
		MerName: name,
		Status:  1,
		MerState: 1,
		IsAudit: 1,
		MerStar: 5,
	}
	applyProfile(m, in)
	m.MerName = name
	if err := s.store.CreateMerchant(ctx, m); err != nil {
		return nil, err
	}
	return s.GetMerchant(ctx, m.MerID)
}

func (s *Service) SetMerchantRecommend(ctx context.Context, id uint, enabled bool) error {
	m, err := s.GetMerchant(ctx, id)
	if err != nil {
		return err
	}
	if enabled {
		m.IsBest = 1
	} else {
		m.IsBest = 0
	}
	return s.store.UpdateMerchant(ctx, m)
}

func applyBool8(dst *int8, src *bool) {
	if src == nil {
		return
	}
	if *src {
		*dst = 1
	} else {
		*dst = 0
	}
}

func applyProfile(m *Merchant, in ShopProfileInput) {
	if v := strings.TrimSpace(in.MerName); v != "" {
		m.MerName = v
	}
	m.RealName = strings.TrimSpace(in.RealName)
	m.MerPhone = strings.TrimSpace(in.MerPhone)
	m.MerAddress = strings.TrimSpace(in.MerAddress)
	m.MerInfo = strings.TrimSpace(in.MerInfo)
	m.MerKeyword = strings.TrimSpace(in.MerKeyword)
	m.Mark = strings.TrimSpace(in.Mark)
	m.MerAccount = strings.TrimSpace(in.MerAccount)
	m.SubMchid = strings.TrimSpace(in.SubMchid)
	m.ApplymentID = strings.TrimSpace(in.ApplymentID)
	m.CategoryID = in.CategoryID
	m.TypeID = in.TypeID
	m.BusinessID = in.BusinessID
	m.RegionID = in.RegionID
	m.StoreGroupIDs = in.StoreGroupIDs
	applyBool8(&m.IsBest, in.IsBest)
	applyBool8(&m.OfflinePay, in.OfflinePay)
	applyBool8(&m.IsTrader, in.IsTrader)
	applyBool8(&m.IsAudit, in.IsAudit)
	applyBool8(&m.IsBroRoom, in.IsBroRoom)
	applyBool8(&m.IsBroGoods, in.IsBroGoods)
	applyBool8(&m.CommissionSwitch, in.CommissionSwitch)
	if in.CommissionRate != nil {
		m.CommissionRate = *in.CommissionRate
	}
	if in.CareCount != nil {
		m.CareCount = *in.CareCount
	}
	if in.CareFicti != nil {
		m.CareFicti = *in.CareFicti
	}
	if in.MerStar != nil {
		m.MerStar = *in.MerStar
	}
	m.MerAvatar = strings.TrimSpace(in.MerAvatar)
	if in.GoodsTypes != nil {
		m.GoodsTypes = append([]int(nil), in.GoodsTypes...)
		m.GoodsType = joinInts(in.GoodsTypes)
	}
	if in.PlatformCategoryIDs != nil {
		m.PlatformCategoryIDList = append([]uint(nil), in.PlatformCategoryIDs...)
		m.PlatformCategoryIDs = joinUints(in.PlatformCategoryIDs)
	}
	if in.Sort != nil {
		m.Sort = *in.Sort
	}
	if in.Status != nil {
		if *in.Status {
			m.Status, m.MerState = 1, 1
		} else {
			m.Status, m.MerState = 0, 0
		}
	}
}

func joinInts(vals []int) string {
	if len(vals) == 0 {
		return ""
	}
	parts := make([]string, 0, len(vals))
	for _, v := range vals {
		parts = append(parts, strconv.Itoa(v))
	}
	return strings.Join(parts, ",")
}

func joinUints(vals []uint) string {
	if len(vals) == 0 {
		return ""
	}
	parts := make([]string, 0, len(vals))
	for _, v := range vals {
		parts = append(parts, strconv.FormatUint(uint64(v), 10))
	}
	return strings.Join(parts, ",")
}

func parseIntsCSV(raw string) []int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []int{}
	}
	parts := strings.Split(raw, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err == nil {
			out = append(out, n)
		}
	}
	return out
}

func parseUintsCSV(raw string) []uint {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []uint{}
	}
	parts := strings.Split(raw, ",")
	out := make([]uint, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.ParseUint(strings.TrimSpace(p), 10, 64)
		if err == nil {
			out = append(out, uint(n))
		}
	}
	return out
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

func (s *Service) ListIntentions(ctx context.Context, keyword string, status *int8, regionIDs []uint, page, limit int, dateFrom, dateTo string, categoryID, typeID *uint) (*PageResult[Intention], error) {
	list, total, err := s.store.ListIntentions(ctx, keyword, status, regionIDs, page, limit, dateFrom, dateTo, categoryID, typeID)
	if err != nil {
		return nil, err
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
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

func (s *Service) DeleteIntention(ctx context.Context, id uint, regionIDs []uint) error {
	if id == 0 {
		return ErrBadParam
	}
	if _, err := s.GetIntention(ctx, id, regionIDs); err != nil {
		return err
	}
	deleted, err := s.store.DeleteIntention(ctx, id, regionIDs)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
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
