package promotionpersist

import "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/promotion"

type StoreAdapter struct {
	*Repo
}

func NewStoreAdapter(repo *Repo) *StoreAdapter {
	return &StoreAdapter{Repo: repo}
}

var _ promotion.Store = (*StoreAdapter)(nil)
