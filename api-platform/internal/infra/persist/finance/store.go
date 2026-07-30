package financepersist

import "github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/finance"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ finance.Store = (*StoreAdapter)(nil)
