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

func (s *StoreAdapter) ListMerchants(ctx context.Context, keyword string, status *int8, scope merchant.MerchantScope, page, limit int) ([]merchant.Merchant, int64, error) {
	return s.Repo.ListMerchants(ctx, ListMerchantsFilter{Keyword: keyword, Status: status, MerchantIDs: scope.MerchantIDs, RegionIDs: scope.RegionIDs, Page: page, Limit: limit})
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

func (s *StoreAdapter) WithTx(fn func(tx merchant.Store) error) error {
	return s.Repo.WithTx(func(tx *Repo) error {
		return fn(&StoreAdapter{Repo: tx})
	})
}

var _ merchant.Store = (*StoreAdapter)(nil)
