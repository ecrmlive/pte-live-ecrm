package merchantpersist

import (
	"context"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/merchant"
)

// StoreAdapter 适配 domain.Store
type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

func (s *StoreAdapter) ListMerchants(ctx context.Context, filter merchant.ListFilter, scope merchant.MerchantScope) ([]merchant.Merchant, int64, error) {
	return s.Repo.ListMerchants(ctx, ListMerchantsFilter{
		Keyword:     filter.Keyword,
		Status:      filter.Status,
		CategoryID:  filter.CategoryID,
		TypeID:      filter.TypeID,
		RegionID:    filter.RegionID,
		IsBest:      filter.IsBest,
		OfflinePay:  filter.OfflinePay,
		DateFrom:    filter.DateFrom,
		DateTo:      filter.DateTo,
		MerchantIDs: scope.MerchantIDs,
		RegionIDs:   scope.RegionIDs,
		Page:        filter.Page,
		Limit:       filter.Limit,
	})
}

func (s *StoreAdapter) ListIntentions(ctx context.Context, keyword string, status *int8, regionIDs []uint, page, limit int) ([]merchant.Intention, int64, error) {
	return s.Repo.ListIntentions(ctx, ListIntentionFilter{Keyword: keyword, Status: status, RegionIDs: regionIDs, Page: page, Limit: limit})
}

func (s *StoreAdapter) GetIntention(ctx context.Context, id uint, regionIDs []uint) (*merchant.Intention, error) {
	return s.Repo.GetIntention(ctx, id, regionIDs)
}

func (s *StoreAdapter) AssignIntentionRegion(ctx context.Context, id, regionID uint) (bool, error) {
	return s.Repo.AssignIntentionRegion(ctx, id, regionID)
}

func (s *StoreAdapter) UpsertMerchantView(ctx context.Context, row *merchant.Merchant) error {
	return s.Repo.UpsertMerchantView(ctx, row)
}

func (s *StoreAdapter) UpdateMerchant(ctx context.Context, m *merchant.Merchant) error {
	return s.Repo.UpdateMerchant(ctx, m)
}

func (s *StoreAdapter) CreateMerchant(ctx context.Context, m *merchant.Merchant) error {
	return s.Repo.CreateMerchant(ctx, m)
}

func (s *StoreAdapter) WithTx(fn func(tx merchant.Store) error) error {
	return s.Repo.WithTx(func(tx *Repo) error {
		return fn(&StoreAdapter{Repo: tx})
	})
}

var _ merchant.Store = (*StoreAdapter)(nil)
