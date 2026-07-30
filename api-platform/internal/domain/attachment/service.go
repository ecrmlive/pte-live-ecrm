package attachment

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListCategories(ctx context.Context, merID uint) ([]Category, error)
	CreateCategory(ctx context.Context, c *Category) error
	UpdateCategory(ctx context.Context, c *Category) error
	DeleteCategory(ctx context.Context, id, merID uint) error
	GetCategory(ctx context.Context, id uint) (*Category, error)

	List(ctx context.Context, userType int, cateID uint, attachmentType *int8, page, limit int) ([]Attachment, int64, error)
	Get(ctx context.Context, id uint) (*Attachment, error)
	Create(ctx context.Context, a *Attachment) error
	Delete(ctx context.Context, id uint, userType int) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListCategories(ctx context.Context, merID uint) ([]Category, error) {
	return s.store.ListCategories(ctx, merID)
}

func (s *Service) CreateCategory(ctx context.Context, merID uint, in CategoryInput) (*Category, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrBadParam
	}
	en := strings.TrimSpace(in.EnName)
	if en == "" {
		en = "cate"
	}
	c := &Category{
		PID: in.PID, AttachmentCategoryName: name, AttachmentCategoryEnname: en,
		Sort: in.Sort, MerID: merID, CreateTime: time.Now(),
	}
	if err := s.store.CreateCategory(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) UpdateCategory(ctx context.Context, merID, id uint, in CategoryInput) (*Category, error) {
	c, err := s.store.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if c.MerID != merID {
		return nil, ErrForbidden
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		c.AttachmentCategoryName = name
	}
	if en := strings.TrimSpace(in.EnName); en != "" {
		c.AttachmentCategoryEnname = en
	}
	c.PID = in.PID
	c.Sort = in.Sort
	if err := s.store.UpdateCategory(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCategory(ctx context.Context, merID, id uint) error {
	c, err := s.store.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if c.MerID != merID {
		return ErrForbidden
	}
	return s.store.DeleteCategory(ctx, id, merID)
}

func (s *Service) List(ctx context.Context, userType int, cateID uint, attachmentType *int8, page, limit int) (*PageResult[Attachment], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	list, total, err := s.store.List(ctx, userType, cateID, attachmentType, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Attachment]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateFile(ctx context.Context, userType int, userID, cateID uint, name, src string, attachmentType int8) (*Attachment, error) {
	name = strings.TrimSpace(name)
	src = strings.TrimSpace(src)
	if name == "" || src == "" || (attachmentType != 0 && attachmentType != 1) {
		return nil, ErrBadParam
	}
	if cateID > 0 {
		c, err := s.store.GetCategory(ctx, cateID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrBadParam
			}
			return nil, err
		}
		// userType=0 平台；>0 为 mer_id
		wantMer := uint(0)
		if userType > 0 {
			wantMer = uint(userType)
		}
		if c.MerID != wantMer {
			return nil, ErrForbidden
		}
	}
	a := &Attachment{
		AttachmentCategoryID: cateID, AttachmentName: name, AttachmentSrc: src,
		UploadType: 1, UserType: userType, UserID: userID,
		CreateTime: time.Now(), AttachmentType: attachmentType,
	}
	if err := s.store.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) Delete(ctx context.Context, userType int, id uint) error {
	a, err := s.store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if a.UserType != userType {
		return ErrForbidden
	}
	return s.store.Delete(ctx, id, userType)
}
