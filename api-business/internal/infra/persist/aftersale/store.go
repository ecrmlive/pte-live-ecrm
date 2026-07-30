package aftersalepersist

import "github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/aftersale"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ aftersale.Store = (*StoreAdapter)(nil)
