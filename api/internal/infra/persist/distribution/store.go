package distributionpersist

import "github.com/qixi-live/qixi-live-mergers/api/internal/domain/distribution"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ distribution.Store = (*StoreAdapter)(nil)
