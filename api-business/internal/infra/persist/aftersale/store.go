package aftersalepersist

import "github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/aftersale"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ aftersale.Store = (*StoreAdapter)(nil)
