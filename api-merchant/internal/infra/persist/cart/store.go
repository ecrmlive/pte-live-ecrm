package cartpersist

import "github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/cart"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ cart.Store = (*StoreAdapter)(nil)
