package tradepersist

import "github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/trade"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ trade.Store = (*StoreAdapter)(nil)
