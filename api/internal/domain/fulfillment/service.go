package fulfillment

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListStaff(ctx context.Context, merID uint, page, limit int) ([]Staff, int64, error)
	GetStaff(ctx context.Context, merID, id uint) (*Staff, error)
	CreateStaff(ctx context.Context, row *Staff) error
	UpdateStaff(ctx context.Context, row *Staff) error
	SoftDeleteStaff(ctx context.Context, merID, id uint) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) List(ctx context.Context, merID uint, page, limit int) (*PageResult[Staff], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListStaff(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Staff]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Create(ctx context.Context, merID uint, in StaffInput) (*Staff, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	name := strings.TrimSpace(in.Name)
	phone := strings.TrimSpace(in.Phone)
	if name == "" || phone == "" {
		return nil, ErrBadParam
	}
	row := &Staff{MerID: merID, Name: name, Phone: phone, Status: 1, CreateTime: time.Now()}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.CreateStaff(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) Update(ctx context.Context, merID, id uint, in StaffInput) (*Staff, error) {
	row, err := s.store.GetStaff(ctx, merID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		row.Name = name
	}
	if phone := strings.TrimSpace(in.Phone); phone != "" {
		row.Phone = phone
	}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.UpdateStaff(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) Delete(ctx context.Context, merID, id uint) error {
	if _, err := s.store.GetStaff(ctx, merID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteStaff(ctx, merID, id)
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
