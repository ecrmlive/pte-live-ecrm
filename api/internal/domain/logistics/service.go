package logistics

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListExpress(ctx context.Context, page, limit int, showOnly bool) ([]Express, int64, error)
	GetExpress(ctx context.Context, id uint) (*Express, error)
	CreateExpress(ctx context.Context, row *Express) error
	UpdateExpress(ctx context.Context, row *Express) error
	SoftDeleteExpress(ctx context.Context, id uint) error

	ListCity(ctx context.Context, parentID *uint) ([]City, error)

	ListTemplate(ctx context.Context, merID uint, page, limit int) ([]ShippingTemplate, int64, error)
	GetTemplate(ctx context.Context, merID, id uint) (*ShippingTemplate, error)
	ListRegions(ctx context.Context, templateID uint) ([]Region, error)
	CreateTemplate(ctx context.Context, row *ShippingTemplate, regions []Region) error
	UpdateTemplate(ctx context.Context, row *ShippingTemplate, regions []Region) error
	SoftDeleteTemplate(ctx context.Context, merID, id uint) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListExpress(ctx context.Context, page, limit int, showOnly bool) (*PageResult[Express], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListExpress(ctx, page, limit, showOnly)
	if err != nil {
		return nil, err
	}
	return &PageResult[Express]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetExpressName(ctx context.Context, id uint) (string, error) {
	row, err := s.store.GetExpress(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}
	return row.Name, nil
}

func (s *Service) CreateExpress(ctx context.Context, in ExpressInput) (*Express, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrBadParam
	}
	row := &Express{
		Name: name, Code: strings.TrimSpace(in.Code), Sort: in.Sort,
		IsShow: 1, CreateTime: time.Now(),
	}
	if in.IsShow != nil {
		row.IsShow = *in.IsShow
	}
	if err := s.store.CreateExpress(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateExpress(ctx context.Context, id uint, in ExpressInput) (*Express, error) {
	row, err := s.store.GetExpress(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		row.Name = name
	}
	if code := strings.TrimSpace(in.Code); code != "" {
		row.Code = code
	}
	row.Sort = in.Sort
	if in.IsShow != nil {
		row.IsShow = *in.IsShow
	}
	if err := s.store.UpdateExpress(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) DeleteExpress(ctx context.Context, id uint) error {
	if _, err := s.store.GetExpress(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteExpress(ctx, id)
}

func (s *Service) ListCity(ctx context.Context, parentID *uint) ([]City, error) {
	return s.store.ListCity(ctx, parentID)
}

func (s *Service) ListTemplate(ctx context.Context, merID uint, page, limit int) (*PageResult[ShippingTemplate], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListTemplate(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[ShippingTemplate]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetTemplate(ctx context.Context, merID, id uint) (*ShippingTemplate, error) {
	if merID == 0 || id == 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.GetTemplate(ctx, merID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	regions, err := s.store.ListRegions(ctx, id)
	if err != nil {
		return nil, err
	}
	row.Regions = regions
	return row, nil
}

func (s *Service) CreateTemplate(ctx context.Context, merID uint, in TemplateInput) (*ShippingTemplate, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrBadParam
	}
	typ := in.Type
	if typ == 0 {
		typ = 1
	}
	row := &ShippingTemplate{
		MerID: merID, Name: name, Type: typ, Appoint: in.Appoint,
		Sort: in.Sort, CreateTime: time.Now(),
	}
	regions := toRegions(0, in.Regions)
	if err := s.store.CreateTemplate(ctx, row, regions); err != nil {
		return nil, err
	}
	row.Regions = regions
	return row, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, merID, id uint, in TemplateInput) (*ShippingTemplate, error) {
	row, err := s.store.GetTemplate(ctx, merID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		row.Name = name
	}
	if in.Type > 0 {
		row.Type = in.Type
	}
	row.Appoint = in.Appoint
	row.Sort = in.Sort
	regions := toRegions(id, in.Regions)
	if err := s.store.UpdateTemplate(ctx, row, regions); err != nil {
		return nil, err
	}
	row.Regions = regions
	return row, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, merID, id uint) error {
	if _, err := s.store.GetTemplate(ctx, merID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteTemplate(ctx, merID, id)
}

func toRegions(templateID uint, in []RegionInput) []Region {
	if len(in) == 0 {
		return []Region{{
			TemplateID: templateID, First: 1, FirstPrice: 0, Continue: 1, ContinuePrice: 0,
		}}
	}
	out := make([]Region, 0, len(in))
	for _, r := range in {
		first := r.First
		if first <= 0 {
			first = 1
		}
		cont := r.Continue
		if cont <= 0 {
			cont = 1
		}
		out = append(out, Region{
			TemplateID: templateID, CityIDs: strings.TrimSpace(r.CityIDs),
			First: first, FirstPrice: r.FirstPrice, Continue: cont, ContinuePrice: r.ContinuePrice,
		})
	}
	return out
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
