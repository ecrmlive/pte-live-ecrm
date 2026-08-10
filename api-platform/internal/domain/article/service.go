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
	UpdateCategoryStatus(ctx context.Context, id uint, status int8) error
	SoftDeleteCategory(ctx context.Context, id uint) error

	ListArticle(ctx context.Context, page, limit int, cid uint, title string, publishedOnly bool) ([]Article, int64, error)
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
	row := &Category{
		Title:      title,
		Info:       strings.TrimSpace(in.Info),
		Image:      strings.TrimSpace(in.Image),
		Status:     1,
		Sort:       in.Sort,
		CreateTime: time.Now(),
	}
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
	row.Info = strings.TrimSpace(in.Info)
	row.Image = strings.TrimSpace(in.Image)
	row.Sort = in.Sort
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.UpdateCategory(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) SetCategoryStatus(ctx context.Context, id uint, status int8) error {
	if status != 0 && status != 1 {
		return ErrBadParam
	}
	if _, err := s.store.GetCategory(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.UpdateCategoryStatus(ctx, id, status)
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

func (s *Service) ListAdmin(ctx context.Context, page, limit int, cid uint, title string) (*PageResult[Article], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListArticle(ctx, page, limit, cid, strings.TrimSpace(title), false)
	if err != nil {
		return nil, err
	}
	return &PageResult[Article]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListApp(ctx context.Context, page, limit int, cid uint) (*PageResult[Article], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListArticle(ctx, page, limit, cid, "", true)
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
	author := strings.TrimSpace(in.Author)
	image := strings.TrimSpace(in.Image)
	content := strings.TrimSpace(in.Content)
	if title == "" || author == "" || image == "" || content == "" || in.CID == 0 {
		return nil, ErrBadParam
	}
	if _, err := s.store.GetCategory(ctx, in.CID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBadParam
		}
		return nil, err
	}
	row := &Article{
		CID: in.CID, Title: title, Author: author,
		Image: image, Synopsis: strings.TrimSpace(in.Synopsis),
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
	// 列表「是否显示」开关：仅传 status 时走轻量更新。
	if strings.TrimSpace(in.Title) == "" && strings.TrimSpace(in.Content) == "" &&
		in.CID == 0 && strings.TrimSpace(in.Author) == "" && strings.TrimSpace(in.Image) == "" &&
		in.Status != nil {
		row.Status = *in.Status
		if err := s.store.UpdateArticle(ctx, row); err != nil {
			return nil, err
		}
		return row, nil
	}
	title := strings.TrimSpace(in.Title)
	author := strings.TrimSpace(in.Author)
	image := strings.TrimSpace(in.Image)
	content := strings.TrimSpace(in.Content)
	if title == "" || author == "" || image == "" || content == "" || in.CID == 0 {
		return nil, ErrBadParam
	}
	if _, err := s.store.GetCategory(ctx, in.CID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBadParam
		}
		return nil, err
	}
	row.Title = title
	row.Content = in.Content
	row.CID = in.CID
	row.Author = author
	row.Image = image
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
