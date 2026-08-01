package contentpersist

import "github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/content"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ content.Store = (*StoreAdapter)(nil)
