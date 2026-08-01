package merchantpersist

import (
	"context"

	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/merchant"
)

// StoreAdapter 适配 domain.Store
type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

func (s *StoreAdapter) ListMerchants(ctx context.Context, keyword string, status *int8, regionIDs []uint, page, limit int) ([]merchant.Merchant, int64, error) {
	return s.Repo.ListMerchants(ctx, ListMerchantsFilter{Keyword: keyword, Status: status, RegionIDs: regionIDs, Page: page, Limit: limit})
}

func (s *StoreAdapter) ListIntentions(ctx context.Context, keyword string, status *int8, page, limit int) ([]merchant.Intention, int64, error) {
	return s.Repo.ListIntentions(ctx, ListIntentionFilter{Keyword: keyword, Status: status, Page: page, Limit: limit})
}

func (s *StoreAdapter) WithTx(fn func(tx merchant.Store) error) error {
	return s.Repo.WithTx(func(tx *Repo) error {
		return fn(&StoreAdapter{Repo: tx})
	})
}

var _ merchant.Store = (*StoreAdapter)(nil)
