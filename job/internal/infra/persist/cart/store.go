package cartpersist

import "github.com/crmlive/qixi-live-ecrm/job/internal/domain/cart"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ cart.Store = (*StoreAdapter)(nil)
