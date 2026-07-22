package diy

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, merID uint, page, limit int) ([]Page, int64, error)
	Get(ctx context.Context, id uint) (*Page, error)
	GetActiveHome(ctx context.Context, merID uint) (*Page, error)
	Create(ctx context.Context, p *Page) error
	Update(ctx context.Context, p *Page) error
	ClearActive(ctx context.Context, merID uint) error
	SoftDelete(ctx context.Context, id uint) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) List(ctx context.Context, merID uint, page, limit int) (*PageResult, error) {
	list, total, err := s.store.List(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		v := list[i].ParseValue()
		list[i].Parsed = &v
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &PageResult{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*Page, error) {
	p, err := s.store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	v := p.ParseValue()
	p.Parsed = &v
	return p, nil
}

// GetActiveHome 当前启用的平台/商户首页；无则返回 nil。
func (s *Service) GetActiveHome(ctx context.Context, merID uint) (*Page, error) {
	p, err := s.store.GetActiveHome(ctx, merID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	v := p.ParseValue()
	p.Parsed = &v
	return p, nil
}

func (s *Service) Create(ctx context.Context, merID uint, in SaveInput) (*Page, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrBadParam
	}
	raw, err := marshalValue(in.Value)
	if err != nil {
		return nil, err
	}
	p := &Page{
		Version:      "1.0",
		Name:         name,
		Title:        strings.TrimSpace(in.Title),
		TemplateName: strings.TrimSpace(in.TemplateName),
		Type:         0,
		IsShow:       1,
		IsDiy:        1,
		MerID:        merID,
		Value:        raw,
		Status:       0,
	}
	if p.TemplateName == "" {
		p.TemplateName = "home"
	}
	if in.Status != nil && *in.Status == 1 {
		if err := s.store.ClearActive(ctx, merID); err != nil {
			return nil, err
		}
		p.Status = 1
		p.IsDefault = 1
	}
	if err := s.store.Create(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, p.ID)
}

func (s *Service) Update(ctx context.Context, id, merID uint, in SaveInput) (*Page, error) {
	p, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if p.MerID != merID {
		return nil, ErrNotFound
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrBadParam
	}
	raw, err := marshalValue(in.Value)
	if err != nil {
		return nil, err
	}
	p.Name = name
	p.Title = strings.TrimSpace(in.Title)
	if tn := strings.TrimSpace(in.TemplateName); tn != "" {
		p.TemplateName = tn
	}
	p.Value = raw
	if in.Status != nil && *in.Status == 1 {
		if err := s.store.ClearActive(ctx, merID); err != nil {
			return nil, err
		}
		p.Status = 1
		p.IsDefault = 1
	} else if in.Status != nil {
		p.Status = *in.Status
	}
	if err := s.store.Update(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s *Service) SetActive(ctx context.Context, id, merID uint) (*Page, error) {
	p, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if p.MerID != merID {
		return nil, ErrNotFound
	}
	if err := s.store.ClearActive(ctx, merID); err != nil {
		return nil, err
	}
	p.Status = 1
	p.IsDefault = 1
	if err := s.store.Update(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id, merID uint) error {
	p, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if p.MerID != merID {
		return ErrNotFound
	}
	return s.store.SoftDelete(ctx, id)
}

func marshalValue(v PageValue) (string, error) {
	if v.Banners == nil {
		v.Banners = []BannerItem{}
	}
	if v.Menus == nil {
		v.Menus = []MenuItem{}
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
