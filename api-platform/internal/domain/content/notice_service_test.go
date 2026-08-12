package content

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

type noticeStore struct {
	created  *Notice
	scopeIDs []uint
}

func (s *noticeStore) ListNotices(context.Context, bool, NoticeListFilter) ([]Notice, int64, error) {
	return nil, 0, nil
}
func (s *noticeStore) GetNotice(_ context.Context, id uint) (*Notice, error) {
	if s.created == nil || s.created.NoticeID != id {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *s.created
	return &copy, nil
}
func (s *noticeStore) CreateNotice(_ context.Context, n *Notice, ids []uint) error {
	n.NoticeID = 101
	s.created = n
	s.scopeIDs = append([]uint(nil), ids...)
	return nil
}
func (s *noticeStore) UpdateNotice(_ context.Context, n *Notice, ids []uint) error {
	s.created = n
	s.scopeIDs = append([]uint(nil), ids...)
	return nil
}
func (s *noticeStore) ListNoticeScopes(_ context.Context, ids []uint) ([]NoticeScope, error) {
	if len(ids) == 0 || s.created == nil || s.created.ScopeType == NoticeScopeAll {
		return nil, nil
	}
	items := make([]NoticeScope, 0, len(s.scopeIDs))
	for _, id := range s.scopeIDs {
		items = append(items, NoticeScope{NoticeID: s.created.NoticeID, ScopeID: id, ScopeType: s.created.ScopeType, Name: "关联对象"})
	}
	return items, nil
}
func (s *noticeStore) UpdateNoticeStatus(context.Context, uint, int8) error { return nil }
func (s *noticeStore) SoftDeleteNotice(context.Context, uint) error         { return nil }
func (s *noticeStore) GetCache(context.Context, string) (*Cache, error) {
	return nil, gorm.ErrRecordNotFound
}
func (s *noticeStore) UpsertCache(context.Context, *Cache) error { return nil }

func TestNoticeCreateUsesAssociationScope(t *testing.T) {
	store := &noticeStore{}
	svc := NewService(store)
	created, err := svc.Create(context.Background(), NoticeInput{
		Title:     "店铺服务通知",
		Content:   "<p>本公告仅发送给已关联店铺。</p>",
		ScopeType: NoticeScopeStoreName,
		ScopeIDs:  []uint{12, 12, 18},
	})
	if err != nil {
		t.Fatalf("create notice: %v", err)
	}
	if created.ScopeType != NoticeScopeStoreName || len(store.scopeIDs) != 2 || store.scopeIDs[0] != 12 || store.scopeIDs[1] != 18 {
		t.Fatalf("notice scope = %#v, saved IDs = %#v", created, store.scopeIDs)
	}
}

func TestNoticeCreateRejectsRawOrMissingAssociation(t *testing.T) {
	svc := NewService(&noticeStore{})
	cases := []NoticeInput{
		{Title: "", Content: "<p>内容</p>", ScopeType: NoticeScopeAll},
		{Title: "超过二十个字符的公告标题超过二十个字符的额外内容", Content: "<p>内容</p>", ScopeType: NoticeScopeAll},
		{Title: "指定店铺", Content: "<p>内容</p>", ScopeType: NoticeScopeStoreName},
		{Title: "错误范围", Content: "<p>内容</p>", ScopeType: "store_id", ScopeIDs: []uint{1}},
	}
	for _, input := range cases {
		if _, err := svc.Create(context.Background(), input); !errors.Is(err, ErrBadParam) {
			t.Fatalf("input %#v error = %v, want ErrBadParam", input, err)
		}
	}
}
