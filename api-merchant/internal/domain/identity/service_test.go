package identity

import (
	"context"
	"errors"
	"testing"
)

func TestBuildMenuTree(t *testing.T) {
	rows := []SystemMenu{
		{MenuID: 1, PID: 0, Path: "/a", MenuName: "A", Sort: 10},
		{MenuID: 2, PID: 1, Path: "/a/b", MenuName: "B", Sort: 5},
		{MenuID: 3, PID: 0, Path: "/c", MenuName: "C", Sort: 20},
	}
	tree := buildMenuTree(rows)
	if len(tree) != 2 {
		t.Fatalf("roots=%d want 2", len(tree))
	}
	if tree[0].MenuID != 3 {
		t.Fatalf("first root=%d want 3 (higher sort)", tree[0].MenuID)
	}
	if len(tree[1].Children) != 1 || tree[1].Children[0].MenuID != 2 {
		t.Fatalf("child under A missing")
	}
}

type stubStore struct {
	Store
	users map[string]*User
	seq   uint
}

func (s *stubStore) FindUserByAccount(_ context.Context, account string) (*User, error) {
	u, ok := s.users[account]
	if !ok {
		return nil, errors.New("record not found")
	}
	return u, nil
}

func (s *stubStore) FindUserByID(_ context.Context, id uint) (*User, error) {
	for _, u := range s.users {
		if u.UID == id {
			return u, nil
		}
	}
	return nil, errors.New("record not found")
}

func (s *stubStore) CreateUser(_ context.Context, user *User) error {
	s.seq++
	user.UID = s.seq
	cp := *user
	s.users[user.Account] = &cp
	return nil
}

func (s *stubStore) TouchUserLogin(_ context.Context, _ uint, _ string) error { return nil }

func (s *stubStore) ListUsers(_ context.Context, page, limit int) ([]User, int64, error) {
	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, *u)
	}
	return out, int64(len(out)), nil
}

func (s *stubStore) UpdateUserSvip(_ context.Context, u *User) error {
	for k, cur := range s.users {
		if cur.UID == u.UID {
			cur.IsSvip = u.IsSvip
			cur.SvipEndtime = u.SvipEndtime
			s.users[k] = cur
			return nil
		}
	}
	return errors.New("record not found")
}

func TestRegisterAppWeakPassword(t *testing.T) {
	svc := NewService(&stubStore{users: map[string]*User{}})
	_, err := svc.RegisterApp(context.Background(), "a", "123", "", "127.0.0.1")
	if !errors.Is(err, ErrWeakPassword) {
		t.Fatalf("want ErrWeakPassword got %v", err)
	}
}
