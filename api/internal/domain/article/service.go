package article

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListCategory(ctx context.Context) ([]Category, error)
	GetCategory(ctx context.Context, id uint) (*Category, error)
	CreateCategory(ctx context.Context, row *Category) error
	UpdateCategory(ctx context.Context, row *Category) error
	SoftDeleteCategory(ctx context.Context, id uint) error

	ListArticle(ctx context.Context, page, limit int, cid uint, publishedOnly bool) ([]Article, int64, error)
	GetArticle(ctx context.Context, id uint) (*Article, error)
	CreateArticle(ctx context.Context, row *Article) error
	UpdateArticle(ctx context.Context, row *Article) error
	SoftDeleteArticle(ctx context.Context, id uint) error
	IncrVisit(ctx context.Context, id uint) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.store.ListCategory(ctx)
}

func (s *Service) CreateCategory(ctx context.Context, in CategoryInput) (*Category, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, ErrBadParam
	}
	row := &Category{Title: title, Status: 1, Sort: in.Sort}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.CreateCategory(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id uint, in CategoryInput) (*Category, error) {
	row, err := s.store.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if title := strings.TrimSpace(in.Title); title != "" {
		row.Title = title
	}
	row.Sort = in.Sort
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.UpdateCategory(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id uint) error {
	if _, err := s.store.GetCategory(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteCategory(ctx, id)
}

func (s *Service) ListAdmin(ctx context.Context, page, limit int, cid uint) (*PageResult[Article], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListArticle(ctx, page, limit, cid, false)
	if err != nil {
		return nil, err
	}
	return &PageResult[Article]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListApp(ctx context.Context, page, limit int, cid uint) (*PageResult[Article], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListArticle(ctx, page, limit, cid, true)
	if err != nil {
		return nil, err
	}
	return &PageResult[Article]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetAdmin(ctx context.Context, id uint) (*Article, error) {
	row, err := s.store.GetArticle(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (s *Service) GetApp(ctx context.Context, id uint) (*Article, error) {
	row, err := s.store.GetArticle(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if row.Status != 1 {
		return nil, ErrNotFound
	}
	_ = s.store.IncrVisit(ctx, id)
	row.Visit++
	return row, nil
}

func (s *Service) Create(ctx context.Context, in ArticleInput) (*Article, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" || strings.TrimSpace(in.Content) == "" {
		return nil, ErrBadParam
	}
	row := &Article{
		CID: in.CID, Title: title, Author: strings.TrimSpace(in.Author),
		Image: strings.TrimSpace(in.Image), Synopsis: strings.TrimSpace(in.Synopsis),
		Content: in.Content, Sort: in.Sort, Status: 1, CreateTime: time.Now(),
	}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.CreateArticle(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) Update(ctx context.Context, id uint, in ArticleInput) (*Article, error) {
	row, err := s.store.GetArticle(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if title := strings.TrimSpace(in.Title); title != "" {
		row.Title = title
	}
	if content := strings.TrimSpace(in.Content); content != "" {
		row.Content = in.Content
	}
	row.CID = in.CID
	row.Author = strings.TrimSpace(in.Author)
	row.Image = strings.TrimSpace(in.Image)
	row.Synopsis = strings.TrimSpace(in.Synopsis)
	row.Sort = in.Sort
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.UpdateArticle(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if _, err := s.store.GetArticle(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteArticle(ctx, id)
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
