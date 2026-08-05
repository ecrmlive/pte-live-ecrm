package community

import (
	"context"
	"errors"
	"testing"
)

// userDeleteStore embeds Store so this test only implements methods used by the use case.
// The nil embedded store is safe because no unrelated method is called.
type userDeleteStore struct {
	Store
	post    *Post
	deleted uint
}

func (s *userDeleteStore) GetPost(_ context.Context, id uint) (*Post, error) {
	if s.post == nil || s.post.CommunityID != id {
		return nil, ErrNotFound
	}
	copy := *s.post
	return &copy, nil
}

func (s *userDeleteStore) LoadUserNickname(_ context.Context, _ uint) (string, error) {
	return "演示用户", nil
}

func (s *userDeleteStore) LoadTopicName(_ context.Context, _ uint) (string, error) {
	return "", nil
}

func (s *userDeleteStore) LoadCateName(_ context.Context, _ uint) (string, error) {
	return "", nil
}

func (s *userDeleteStore) SoftDeletePost(_ context.Context, id uint) error {
	s.deleted = id
	return nil
}

func TestDeleteUserPostOnlyAuthor(t *testing.T) {
	cases := []struct {
		name    string
		uid     uint
		want    error
		deleted uint
	}{
		{name: "owner", uid: 11, want: nil, deleted: 7},
		{name: "other user", uid: 12, want: ErrNotFound, deleted: 0},
		{name: "zero user", uid: 0, want: ErrBadParam, deleted: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &userDeleteStore{post: &Post{CommunityID: 7, UID: 11}}
			err := NewService(store).DeleteUserPost(context.Background(), tc.uid, 7)
			if !errors.Is(err, tc.want) {
				t.Fatalf("error = %v, want %v", err, tc.want)
			}
			if store.deleted != tc.deleted {
				t.Fatalf("deleted = %d, want %d", store.deleted, tc.deleted)
			}
		})
	}
}
