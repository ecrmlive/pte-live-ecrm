package cartpersist

import "github.com/crmlive/pte-live-ecrm/api-business/internal/domain/cart"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ cart.Store = (*StoreAdapter)(nil)
