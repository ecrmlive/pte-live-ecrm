package community

import (
	"context"
	"testing"
)

type listFilterStore struct {
	Store
	filter ListFilter
}

func (s *listFilterStore) ListPosts(_ context.Context, filter ListFilter) ([]Post, int64, error) {
	s.filter = filter
	return []Post{}, 0, nil
}

func TestListAppCombinesTopicAndCategoryFilters(t *testing.T) {
	store := &listFilterStore{}
	result, err := NewService(store).ListApp(context.Background(), 8, 3, 2, 15)
	if err != nil {
		t.Fatalf("ListApp error = %v", err)
	}
	if result.Page != 2 || result.Limit != 15 {
		t.Fatalf("page result = %d,%d, want 2,15", result.Page, result.Limit)
	}
	if store.filter.TopicID != 8 || store.filter.CategoryID != 3 || !store.filter.OnlyPublic {
		t.Fatalf("filter = %#v, want public topic=8 category=3", store.filter)
	}
}
