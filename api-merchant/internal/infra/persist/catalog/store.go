package catalogpersist

import (
	"context"

	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/catalog"
)

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

func (s *StoreAdapter) ListProducts(ctx context.Context, status *int8, keyword string, merID *uint, page, limit int) ([]catalog.Product, int64, error) {
	return s.Repo.ListProducts(ctx, ListProductFilter{
		Status:  status,
		Keyword: keyword,
		MerID:   merID,
		Page:    page,
		Limit:   limit,
	})
}

var _ catalog.Store = (*StoreAdapter)(nil)
