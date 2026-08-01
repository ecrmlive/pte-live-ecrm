package financepersist

import "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/finance"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ finance.Store = (*StoreAdapter)(nil)
