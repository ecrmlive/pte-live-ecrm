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
	GetCategoryByEnname(ctx context.Context, merID uint, enname string) (*Category, error)

	List(ctx context.Context, userType int, cateID uint, systemOnly bool, attachmentType *int8, page, limit int) ([]Attachment, int64, error)
	Get(ctx context.Context, id uint) (*Attachment, error)
	Create(ctx context.Context, a *Attachment) error
	Delete(ctx context.Context, id uint, userType int) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListCategories(ctx context.Context, merID uint, mediaType *int8) ([]Category, error) {
	if merID == 0 {
		if err := s.EnsureSystemCategories(ctx); err != nil {
			return nil, err
		}
	}
	list, err := s.store.ListCategories(ctx, merID)
	if err != nil || mediaType == nil || merID != 0 {
		return list, err
	}
	// 平台侧按图片/视频系统 enname 过滤系统分类；自定义分类始终保留。
	allowed := make(map[string]struct{}, len(SystemCategories))
	for _, en := range SystemCategoryEnnames(mediaType) {
		allowed[en] = struct{}{}
	}
	out := make([]Category, 0, len(list))
	for _, row := range list {
		if row.IsSystem == 1 || IsSystemCategoryEnname(row.AttachmentCategoryEnname) {
			if _, ok := allowed[row.AttachmentCategoryEnname]; !ok {
				continue
			}
		}
		out = append(out, row)
	}
	return out, nil
}

// EnsureSystemCategories 幂等写入平台固定系统分类。
func (s *Service) EnsureSystemCategories(ctx context.Context) error {
	for _, spec := range SystemCategories {
		row, err := s.store.GetCategoryByEnname(ctx, 0, spec.EnName)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			c := &Category{
				AttachmentCategoryName:   spec.Name,
				AttachmentCategoryEnname: spec.EnName,
				Sort:                     spec.Sort,
				MerID:                    0,
				IsSystem:                 1,
				CreateTime:               time.Now(),
			}
			if err := s.store.CreateCategory(ctx, c); err != nil {
				return err
			}
			continue
		}
		needUpdate := row.AttachmentCategoryName != spec.Name ||
			row.Sort != spec.Sort ||
			row.IsSystem != 1
		if !needUpdate {
			continue
		}
		row.AttachmentCategoryName = spec.Name
		row.Sort = spec.Sort
		row.IsSystem = 1
		if err := s.store.UpdateCategory(ctx, row); err != nil {
			return err
		}
	}
	return nil
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
	if merID == 0 && IsSystemCategoryEnname(en) {
		return nil, ErrSystemCategory
	}
	for _, spec := range SystemCategories {
		if merID == 0 && name == spec.Name {
			return nil, ErrSystemCategory
		}
	}
	c := &Category{
		PID: in.PID, AttachmentCategoryName: name, AttachmentCategoryEnname: en,
		Sort: in.Sort, MerID: merID, IsSystem: 0, CreateTime: time.Now(),
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
	if c.IsSystem == 1 || IsSystemCategoryEnname(c.AttachmentCategoryEnname) {
		return nil, ErrSystemCategory
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		c.AttachmentCategoryName = name
	}
	if en := strings.TrimSpace(in.EnName); en != "" {
		if IsSystemCategoryEnname(en) {
			return nil, ErrSystemCategory
		}
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
	if c.IsSystem == 1 || IsSystemCategoryEnname(c.AttachmentCategoryEnname) {
		return ErrSystemCategory
	}
	return s.store.DeleteCategory(ctx, id, merID)
}

func (s *Service) List(ctx context.Context, userType int, cateID uint, systemOnly bool, attachmentType *int8, page, limit int) (*PageResult[Attachment], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	list, total, err := s.store.List(ctx, userType, cateID, systemOnly, attachmentType, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Attachment]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateFile(ctx context.Context, userType int, userID, cateID uint, name, src string, attachmentType int8, isSystem bool) (*Attachment, error) {
	name = strings.TrimSpace(name)
	src = strings.TrimSpace(src)
	if name == "" || src == "" || (attachmentType != 0 && attachmentType != 1) {
		return nil, ErrBadParam
	}
	var cat *Category
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
		cat = c
	}
	sysFlag := int8(0)
	if isSystem {
		// 仅允许把系统预置素材写入固定系统分类
		if cat == nil || cat.IsSystem != 1 || !IsSystemCategoryEnname(cat.AttachmentCategoryEnname) {
			return nil, ErrBadParam
		}
		sysFlag = 1
	}
	a := &Attachment{
		AttachmentCategoryID: cateID, AttachmentName: name, AttachmentSrc: src,
		UploadType: 1, UserType: userType, UserID: userID,
		CreateTime: time.Now(), AttachmentType: attachmentType, IsSystem: sysFlag,
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
