package promotionpersist

import "github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/promotion"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ promotion.Store = (*StoreAdapter)(nil)
