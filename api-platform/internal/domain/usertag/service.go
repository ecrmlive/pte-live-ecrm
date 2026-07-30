package usertag

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListLabel(ctx context.Context, page, limit int) ([]Label, int64, error)
	GetLabel(ctx context.Context, id uint) (*Label, error)
	CreateLabel(ctx context.Context, row *Label) error
	UpdateLabel(ctx context.Context, row *Label) error
	SoftDeleteLabel(ctx context.Context, id uint) error

	ListGroup(ctx context.Context, page, limit int) ([]Group, int64, error)
	GetGroup(ctx context.Context, id uint) (*Group, error)
	CreateGroup(ctx context.Context, row *Group) error
	UpdateGroup(ctx context.Context, row *Group) error
	SoftDeleteGroup(ctx context.Context, id uint) error

	ListUserLabels(ctx context.Context, uid uint) ([]Label, error)
	ReplaceUserLabels(ctx context.Context, uid uint, labelIDs []uint) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListLabels(ctx context.Context, page, limit int) (*PageResult[Label], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListLabel(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Label]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateLabel(ctx context.Context, in LabelInput) (*Label, error) {
	name := strings.TrimSpace(in.LabelName)
	if name == "" {
		return nil, ErrBadParam
	}
	row := &Label{LabelName: name, Sort: in.Sort, CreateTime: time.Now()}
	if err := s.store.CreateLabel(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateLabel(ctx context.Context, id uint, in LabelInput) (*Label, error) {
	row, err := s.store.GetLabel(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.LabelName); name != "" {
		row.LabelName = name
	}
	row.Sort = in.Sort
	if err := s.store.UpdateLabel(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) DeleteLabel(ctx context.Context, id uint) error {
	if _, err := s.store.GetLabel(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteLabel(ctx, id)
}

func (s *Service) ListGroups(ctx context.Context, page, limit int) (*PageResult[Group], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListGroup(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Group]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateGroup(ctx context.Context, in GroupInput) (*Group, error) {
	name := strings.TrimSpace(in.GroupName)
	if name == "" {
		return nil, ErrBadParam
	}
	row := &Group{GroupName: name, Sort: in.Sort, CreateTime: time.Now()}
	if err := s.store.CreateGroup(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateGroup(ctx context.Context, id uint, in GroupInput) (*Group, error) {
	row, err := s.store.GetGroup(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.GroupName); name != "" {
		row.GroupName = name
	}
	row.Sort = in.Sort
	if err := s.store.UpdateGroup(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) DeleteGroup(ctx context.Context, id uint) error {
	if _, err := s.store.GetGroup(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteGroup(ctx, id)
}

func (s *Service) ListUserLabels(ctx context.Context, uid uint) ([]Label, error) {
	if uid == 0 {
		return nil, ErrBadParam
	}
	return s.store.ListUserLabels(ctx, uid)
}

func (s *Service) MarkUser(ctx context.Context, in MarkInput) error {
	if in.UID == 0 {
		return ErrBadParam
	}
	return s.store.ReplaceUserLabels(ctx, in.UID, in.LabelIDs)
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
