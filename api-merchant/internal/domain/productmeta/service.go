package productmeta

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListLabel(ctx context.Context, merID uint, page, limit int) ([]Label, int64, error)
	GetLabel(ctx context.Context, merID, id uint) (*Label, error)
	CreateLabel(ctx context.Context, row *Label) error
	UpdateLabel(ctx context.Context, row *Label) error
	SoftDeleteLabel(ctx context.Context, merID, id uint) error

	ListGuarantee(ctx context.Context, merID uint, page, limit int, keyword string, status *int8) ([]Guarantee, int64, error)
	GetGuarantee(ctx context.Context, merID, id uint) (*Guarantee, error)
	CreateGuarantee(ctx context.Context, row *Guarantee) error
	UpdateGuarantee(ctx context.Context, row *Guarantee) error
	SoftDeleteGuarantee(ctx context.Context, merID, id uint) error

	ListAttrTemplate(ctx context.Context, merID uint, page, limit int, keyword string) ([]AttrTemplate, int64, error)
	GetAttrTemplate(ctx context.Context, merID, id uint) (*AttrTemplate, error)
	CreateAttrTemplate(ctx context.Context, row *AttrTemplate) error
	UpdateAttrTemplate(ctx context.Context, row *AttrTemplate) error
	SoftDeleteAttrTemplate(ctx context.Context, merID, id uint) error

	ListLabelPlatform(ctx context.Context, page, limit int) ([]Label, int64, error)
	ListGuaranteePlatform(ctx context.Context, page, limit int) ([]Guarantee, int64, error)
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListLabels(ctx context.Context, merID uint, page, limit int) (*PageResult[Label], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListLabel(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Label]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListLabelsPlatform(ctx context.Context, page, limit int) (*PageResult[Label], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListLabelPlatform(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Label]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateLabel(ctx context.Context, merID uint, in LabelInput) (*Label, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrBadParam
	}
	row := &Label{MerID: merID, Name: name, Info: strings.TrimSpace(in.Info), Sort: in.Sort, Status: 1, CreateTime: time.Now()}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.CreateLabel(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateLabel(ctx context.Context, merID, id uint, in LabelInput) (*Label, error) {
	row, err := s.store.GetLabel(ctx, merID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		row.Name = name
	}
	row.Info = strings.TrimSpace(in.Info)
	row.Sort = in.Sort
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.UpdateLabel(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) DeleteLabel(ctx context.Context, merID, id uint) error {
	if _, err := s.store.GetLabel(ctx, merID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteLabel(ctx, merID, id)
}

func (s *Service) ListGuarantees(ctx context.Context, merID uint, page, limit int, keyword string, status *int8) (*PageResult[Guarantee], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListGuarantee(ctx, merID, page, limit, keyword, status)
	if err != nil {
		return nil, err
	}
	return &PageResult[Guarantee]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListGuaranteesPlatform(ctx context.Context, page, limit int) (*PageResult[Guarantee], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListGuaranteePlatform(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Guarantee]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateGuarantee(ctx context.Context, merID uint, in GuaranteeInput) (*Guarantee, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrBadParam
	}
	row := &Guarantee{MerID: merID, Name: name, Content: strings.TrimSpace(in.Content), Sort: in.Sort, Status: 1, CreateTime: time.Now()}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.CreateGuarantee(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateGuarantee(ctx context.Context, merID, id uint, in GuaranteeInput) (*Guarantee, error) {
	row, err := s.store.GetGuarantee(ctx, merID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		row.Name = name
	}
	row.Content = strings.TrimSpace(in.Content)
	row.Sort = in.Sort
	if in.Status != nil {
		row.Status = *in.Status
	}
	if err := s.store.UpdateGuarantee(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) DeleteGuarantee(ctx context.Context, merID, id uint) error {
	if _, err := s.store.GetGuarantee(ctx, merID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteGuarantee(ctx, merID, id)
}

func (s *Service) ListAttrTemplates(ctx context.Context, merID uint, page, limit int, keyword string) (*PageResult[AttrTemplate], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListAttrTemplate(ctx, merID, page, limit, keyword)
	if err != nil {
		return nil, err
	}
	return &PageResult[AttrTemplate]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateAttrTemplate(ctx context.Context, merID uint, in AttrTemplateInput) (*AttrTemplate, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	name := strings.TrimSpace(in.TemplateName)
	if name == "" {
		return nil, ErrBadParam
	}
	val := strings.TrimSpace(in.TemplateValue)
	if val == "" {
		val = "[]"
	}
	row := &AttrTemplate{MerID: merID, TemplateName: name, TemplateValue: val, Sort: in.Sort, CreateTime: time.Now()}
	if err := s.store.CreateAttrTemplate(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateAttrTemplate(ctx context.Context, merID, id uint, in AttrTemplateInput) (*AttrTemplate, error) {
	row, err := s.store.GetAttrTemplate(ctx, merID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.TemplateName); name != "" {
		row.TemplateName = name
	}
	if val := strings.TrimSpace(in.TemplateValue); val != "" {
		row.TemplateValue = val
	}
	row.Sort = in.Sort
	if err := s.store.UpdateAttrTemplate(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) DeleteAttrTemplate(ctx context.Context, merID, id uint) error {
	if _, err := s.store.GetAttrTemplate(ctx, merID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteAttrTemplate(ctx, merID, id)
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
