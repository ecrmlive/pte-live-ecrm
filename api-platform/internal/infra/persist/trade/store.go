package tradepersist

import "github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/trade"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ trade.Store = (*StoreAdapter)(nil)
