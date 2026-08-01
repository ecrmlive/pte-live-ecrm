package aftersalepersist

import "github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/aftersale"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ aftersale.Store = (*StoreAdapter)(nil)
