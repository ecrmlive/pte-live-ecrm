package cs

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Store interface {
	ListReplies(ctx context.Context, merID uint, onlyOn bool, page, limit int) ([]Reply, int64, error)
	GetReply(ctx context.Context, merID, id uint) (*Reply, error)
	CreateReply(ctx context.Context, row *Reply) error
	UpdateReply(ctx context.Context, row *Reply) error
	DeleteReply(ctx context.Context, merID, id uint) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListAdmin(ctx context.Context, merID uint, page, limit int) (*PageResult[Reply], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	list, total, err := s.store.ListReplies(ctx, merID, false, page, limit)
	if err != nil {
		return nil, err
	}
	page, limit = normalize(page, limit)
	return &PageResult[Reply]{List: list, Total: total, Page: page, Limit: limit}, nil
}

// ListEnabled 客服工作台：仅启用项。
func (s *Service) ListEnabled(ctx context.Context, merID uint) ([]Reply, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	list, _, err := s.store.ListReplies(ctx, merID, true, 1, 100)
	return list, err
}

func (s *Service) Create(ctx context.Context, merID uint, in ReplyInput) (*Reply, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	keyword := strings.TrimSpace(in.Keyword)
	content := strings.TrimSpace(in.Content)
	if keyword == "" || content == "" {
		return nil, ErrBadParam
	}
	typ := in.Type
	if typ == 0 {
		typ = 1
	}
	if typ != 1 && typ != 2 {
		return nil, ErrBadParam
	}
	row := &Reply{MerID: merID, Type: typ, Keyword: keyword, Content: content, Status: 1}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.CreateReply(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) Update(ctx context.Context, merID, id uint, in ReplyInput) (*Reply, error) {
	row, err := s.getOwned(ctx, merID, id)
	if err != nil {
		return nil, err
	}
	keyword := strings.TrimSpace(in.Keyword)
	content := strings.TrimSpace(in.Content)
	if keyword == "" || content == "" {
		return nil, ErrBadParam
	}
	typ := in.Type
	if typ == 0 {
		typ = row.Type
	}
	if typ != 1 && typ != 2 {
		return nil, ErrBadParam
	}
	row.Type = typ
	row.Keyword = keyword
	row.Content = content
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.UpdateReply(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) Delete(ctx context.Context, merID, id uint) error {
	if _, err := s.getOwned(ctx, merID, id); err != nil {
		return err
	}
	return s.store.DeleteReply(ctx, merID, id)
}

func (s *Service) getOwned(ctx context.Context, merID, id uint) (*Reply, error) {
	if merID == 0 || id == 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.GetReply(ctx, merID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if row.MerID != merID {
		return nil, ErrForbidden
	}
	return row, nil
}

func normalize(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}
