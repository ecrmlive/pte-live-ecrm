package contentpersist

import "github.com/crmlive/pte-live-ecrm/api-business/internal/domain/content"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ content.Store = (*StoreAdapter)(nil)
