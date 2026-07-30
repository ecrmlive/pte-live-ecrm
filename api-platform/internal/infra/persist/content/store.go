package contentpersist

import "github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/content"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ content.Store = (*StoreAdapter)(nil)
