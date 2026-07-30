package diy

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, f ListFilter) ([]Page, int64, error)
	Get(ctx context.Context, id uint) (*Page, error)
	GetActiveHome(ctx context.Context, merID uint) (*Page, error)
	Create(ctx context.Context, p *Page) error
	Update(ctx context.Context, p *Page) error
	ClearActive(ctx context.Context, merID uint, isDiy int8) error
	SoftDelete(ctx context.Context, id uint) error
	ListCategories(ctx context.Context, isMer int8) ([]PageCategory, error)
	GetCategory(ctx context.Context, id uint) (*PageCategory, error)
	CreateCategory(ctx context.Context, category *PageCategory) error
	UpdateCategory(ctx context.Context, category *PageCategory) error
	DeleteCategory(ctx context.Context, id uint) error
	CountCategoryChildren(ctx context.Context, id uint) (int64, error)
	CountLinksByCategory(ctx context.Context, id uint) (int64, error)
	ListLinks(ctx context.Context, f LinkListFilter) ([]PageLink, int64, error)
	GetLink(ctx context.Context, id uint) (*PageLink, error)
	CreateLink(ctx context.Context, link *PageLink) error
	UpdateLink(ctx context.Context, link *PageLink) error
	DeleteLink(ctx context.Context, id uint) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) List(ctx context.Context, f ListFilter) (*PageResult, error) {
	list, total, err := s.store.List(ctx, f)
	if err != nil {
		return nil, err
	}
	for i := range list {
		d := list[i].ParseDoc()
		list[i].Doc = &d
	}
	page, limit := f.Page, f.Limit
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
	d := p.ParseDoc()
	p.Doc = &d
	return p, nil
}

func (s *Service) GetActiveHome(ctx context.Context, merID uint) (*Page, error) {
	p, err := s.store.GetActiveHome(ctx, merID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	d := p.ParseDoc()
	p.Doc = &d
	return p, nil
}

func (s *Service) EditorBootstrap(ctx context.Context, id, merID uint) (*EditorBootstrap, error) {
	defs := loadDefaults()
	out := &EditorBootstrap{
		DefaultData: defs.DefaultItems,
		DefaultPage: defs.DefaultPage,
		JSONData: PageDoc{
			Page:  cloneMap(defs.DefaultPage),
			Items: []map[string]any{},
		},
		Opts: map[string]any{},
	}
	if id == 0 {
		return out, nil
	}
	p, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if p.MerID != merID {
		return nil, ErrNotFound
	}
	if p.Doc != nil {
		out.JSONData = *p.Doc
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, merID uint, in SaveInput) (*Page, error) {
	name := strings.TrimSpace(in.Name)
	doc, err := resolveDoc(in)
	if err != nil {
		return nil, err
	}
	if name == "" {
		if n, _ := doc.Page["params"].(map[string]any); n != nil {
			if ns, ok := n["name"].(string); ok {
				name = strings.TrimSpace(ns)
			}
		}
	}
	if name == "" {
		return nil, ErrBadParam
	}
	raw, err := marshalDoc(doc)
	if err != nil {
		return nil, err
	}
	p := &Page{
		Version:      "2.0",
		Name:         name,
		Title:        strings.TrimSpace(in.Title),
		CoverImage:   strings.TrimSpace(in.CoverImage),
		TemplateName: strings.TrimSpace(in.TemplateName),
		Type:         0,
		IsShow:       1,
		IsDiy:        1,
		MerID:        merID,
		Value:        raw,
		Status:       0,
		ColorPicker:  strings.TrimSpace(in.ColorPicker),
		BgPic:        strings.TrimSpace(in.BgPic),
	}
	if p.Title == "" {
		p.Title = name
	}
	if p.TemplateName == "" {
		p.TemplateName = "home"
	}
	applyChrome(p, in)
	// 微页面是独立可访问的已发布内容，不参与首页的单选启用状态。
	if p.IsDiy == 0 {
		p.Status = 1
		p.IsShow = 1
	} else if in.Status != nil && *in.Status == 1 {
		if err := s.store.ClearActive(ctx, merID, p.IsDiy); err != nil {
			return nil, err
		}
		p.Status = 1
		p.IsDefault = 1
		p.IsShow = 1
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
	doc, err := resolveDoc(in)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		name = p.Name
	}
	raw, err := marshalDoc(doc)
	if err != nil {
		return nil, err
	}
	p.Name = name
	if t := strings.TrimSpace(in.Title); t != "" {
		p.Title = t
	}
	if tn := strings.TrimSpace(in.TemplateName); tn != "" {
		p.TemplateName = tn
	}
	if c := strings.TrimSpace(in.CoverImage); c != "" {
		p.CoverImage = c
	}
	p.Value = raw
	p.Version = "2.0"
	applyChrome(p, in)
	if p.IsDiy == 0 {
		p.Status = 1
		p.IsShow = 1
	} else if in.Status != nil && *in.Status == 1 {
		if err := s.store.ClearActive(ctx, merID, p.IsDiy); err != nil {
			return nil, err
		}
		p.Status = 1
		p.IsDefault = 1
		p.IsShow = 1
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
	if p.IsDiy == 0 {
		return nil, ErrBadParam
	}
	if err := s.store.ClearActive(ctx, merID, p.IsDiy); err != nil {
		return nil, err
	}
	p.Status = 1
	p.IsDefault = 1
	p.IsShow = 1
	if err := s.store.Update(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s *Service) Copy(ctx context.Context, id, merID uint) (*Page, error) {
	src, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if src.MerID != merID {
		return nil, ErrNotFound
	}
	return s.cloneToMer(ctx, src, merID, src.Name+" 副本")
}

func (s *Service) ApplyDefault(ctx context.Context, templateID, merID uint) (*Page, error) {
	src, err := s.Get(ctx, templateID)
	if err != nil {
		return nil, err
	}
	if src.MerID != 0 {
		return nil, ErrBadParam
	}
	return s.cloneToMer(ctx, src, merID, src.Name)
}

func (s *Service) cloneToMer(ctx context.Context, src *Page, merID uint, name string) (*Page, error) {
	doc := src.ParseDoc()
	raw, err := marshalDoc(doc)
	if err != nil {
		return nil, err
	}
	p := &Page{
		Version:      "2.0",
		Name:         name,
		Title:        src.Title,
		CoverImage:   src.CoverImage,
		TemplateName: src.TemplateName,
		Type:         src.Type,
		IsShow:       0,
		IsDiy:        src.IsDiy,
		IsBgColor:    src.IsBgColor,
		IsBgPic:      src.IsBgPic,
		ColorPicker:  src.ColorPicker,
		BgPic:        src.BgPic,
		BgTabVal:     src.BgTabVal,
		MerID:        merID,
		Value:        raw,
		Status:       0,
		IsDefault:    0,
	}
	if src.IsDiy == 0 {
		p.Status = 1
		p.IsShow = 1
	}
	if err := s.store.Create(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, p.ID)
}

func (s *Service) Recovery(ctx context.Context, id, merID uint) (*Page, error) {
	p, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if p.MerID != merID {
		return nil, ErrNotFound
	}
	defs := loadDefaults()
	doc := PageDoc{Page: cloneMap(defs.DefaultPage), Items: []map[string]any{}}
	raw, err := marshalDoc(doc)
	if err != nil {
		return nil, err
	}
	p.Value = raw
	p.Version = "2.0"
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
	if p.IsDiy == 1 && p.Status == 1 {
		return ErrBadParam
	}
	return s.store.SoftDelete(ctx, id)
}

func (s *Service) ListCategories(ctx context.Context, isMer int8) ([]PageCategory, error) {
	rows, err := s.store.ListCategories(ctx, isMer)
	if err != nil {
		return nil, err
	}
	return buildCategoryTree(rows), nil
}

func (s *Service) CreateCategory(ctx context.Context, isMer int8, in CategoryInput) (*PageCategory, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || isMer < 0 || isMer > 1 {
		return nil, ErrBadParam
	}
	level, err := s.categoryLevel(ctx, in.PID, isMer)
	if err != nil {
		return nil, err
	}
	status := int8(1)
	if in.Status != nil {
		status = *in.Status
	}
	category := &PageCategory{
		PID: in.PID, Name: name, Type: "link", Sort: in.Sort,
		Status: status, Level: level, IsMer: isMer,
	}
	if err := s.store.CreateCategory(ctx, category); err != nil {
		return nil, err
	}
	return s.store.GetCategory(ctx, category.ID)
}

func (s *Service) UpdateCategory(ctx context.Context, id uint, isMer int8, in CategoryInput) (*PageCategory, error) {
	category, err := s.categoryForScope(ctx, id, isMer)
	if err != nil {
		return nil, err
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		category.Name = name
	}
	if in.PID != category.PID {
		if in.PID == id {
			return nil, ErrBadParam
		}
		children, err := s.store.CountCategoryChildren(ctx, id)
		if err != nil {
			return nil, err
		}
		// 移动非叶子节点需要同时重算整棵子树层级；当前接口明确拒绝，避免三级限制被绕过或形成环。
		if children > 0 {
			return nil, ErrBadParam
		}
		level, err := s.categoryLevel(ctx, in.PID, isMer)
		if err != nil {
			return nil, err
		}
		category.PID, category.Level = in.PID, level
	}
	category.Sort = in.Sort
	if in.Status != nil {
		category.Status = *in.Status
	}
	if err := s.store.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}
	return s.store.GetCategory(ctx, id)
}

func (s *Service) DeleteCategory(ctx context.Context, id uint, isMer int8) error {
	if _, err := s.categoryForScope(ctx, id, isMer); err != nil {
		return err
	}
	children, err := s.store.CountCategoryChildren(ctx, id)
	if err != nil {
		return err
	}
	links, err := s.store.CountLinksByCategory(ctx, id)
	if err != nil {
		return err
	}
	if children > 0 || links > 0 {
		return ErrBadParam
	}
	return s.store.DeleteCategory(ctx, id)
}

func (s *Service) ListLinks(ctx context.Context, f LinkListFilter) (*LinkResult, error) {
	rows, total, err := s.store.ListLinks(ctx, f)
	if err != nil {
		return nil, err
	}
	page, limit := f.Page, f.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &LinkResult{List: rows, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateLink(ctx context.Context, isMer int8, in LinkInput) (*PageLink, error) {
	name, rawURL := strings.TrimSpace(in.Name), strings.TrimSpace(in.URL)
	if name == "" || !isAllowedLinkURL(rawURL) {
		return nil, ErrBadParam
	}
	if _, err := s.categoryForScope(ctx, in.CateID, isMer); err != nil {
		return nil, err
	}
	status := int8(1)
	if in.Status != nil {
		status = *in.Status
	}
	link := &PageLink{CateID: in.CateID, Type: 1, Name: name, URL: rawURL,
		Param: strings.TrimSpace(in.Param), Example: strings.TrimSpace(in.Example),
		Status: status, Sort: in.Sort, IsMer: isMer}
	if err := s.store.CreateLink(ctx, link); err != nil {
		return nil, err
	}
	return s.store.GetLink(ctx, link.ID)
}

func (s *Service) UpdateLink(ctx context.Context, id uint, isMer int8, in LinkInput) (*PageLink, error) {
	link, err := s.linkForScope(ctx, id, isMer)
	if err != nil {
		return nil, err
	}
	if in.CateID != 0 && in.CateID != link.CateID {
		if _, err := s.categoryForScope(ctx, in.CateID, isMer); err != nil {
			return nil, err
		}
		link.CateID = in.CateID
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		link.Name = name
	}
	if rawURL := strings.TrimSpace(in.URL); rawURL != "" {
		if !isAllowedLinkURL(rawURL) {
			return nil, ErrBadParam
		}
		link.URL = rawURL
	}
	link.Param, link.Example, link.Sort = strings.TrimSpace(in.Param), strings.TrimSpace(in.Example), in.Sort
	if in.Status != nil {
		link.Status = *in.Status
	}
	if err := s.store.UpdateLink(ctx, link); err != nil {
		return nil, err
	}
	return s.store.GetLink(ctx, id)
}

func (s *Service) DeleteLink(ctx context.Context, id uint, isMer int8) error {
	if _, err := s.linkForScope(ctx, id, isMer); err != nil {
		return err
	}
	return s.store.DeleteLink(ctx, id)
}

func (s *Service) categoryForScope(ctx context.Context, id uint, isMer int8) (*PageCategory, error) {
	category, err := s.store.GetCategory(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if category.IsMer != isMer || category.Type != "link" {
		return nil, ErrNotFound
	}
	return category, nil
}

func (s *Service) linkForScope(ctx context.Context, id uint, isMer int8) (*PageLink, error) {
	link, err := s.store.GetLink(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if link.IsMer != isMer {
		return nil, ErrNotFound
	}
	return link, nil
}

func (s *Service) categoryLevel(ctx context.Context, parentID uint, isMer int8) (int8, error) {
	if parentID == 0 {
		return 1, nil
	}
	parent, err := s.categoryForScope(ctx, parentID, isMer)
	if err != nil {
		return 0, err
	}
	if parent.Level >= 3 {
		return 0, ErrBadParam
	}
	return parent.Level + 1, nil
}

func isAllowedLinkURL(raw string) bool {
	if strings.HasPrefix(raw, "/pages/") {
		return true
	}
	parsed, err := url.ParseRequestURI(raw)
	return err == nil && parsed.Host != "" && (parsed.Scheme == "https" || parsed.Scheme == "http")
}

func buildCategoryTree(rows []PageCategory) []PageCategory {
	byParent := make(map[uint][]PageCategory, len(rows))
	for _, row := range rows {
		row.Children = []PageCategory{}
		byParent[row.PID] = append(byParent[row.PID], row)
	}
	var collect func(uint) []PageCategory
	collect = func(parentID uint) []PageCategory {
		children := byParent[parentID]
		for i := range children {
			children[i].Children = collect(children[i].ID)
		}
		return children
	}
	return collect(0)
}

func (s *Service) ListDefaults(ctx context.Context, page, limit int) (*PageResult, error) {
	one := int8(1)
	return s.List(ctx, ListFilter{MerID: 0, IsDiy: &one, Page: page, Limit: limit, Status: nil})
}

func applyChrome(p *Page, in SaveInput) {
	if in.IsDiy != nil {
		p.IsDiy = *in.IsDiy
	}
	if in.Type != nil {
		p.Type = *in.Type
	}
	if in.IsShow != nil {
		p.IsShow = *in.IsShow
	}
	if in.IsBgColor != nil {
		p.IsBgColor = *in.IsBgColor
	}
	if in.IsBgPic != nil {
		p.IsBgPic = *in.IsBgPic
	}
	if in.BgTabVal != nil {
		p.BgTabVal = *in.BgTabVal
	}
	if in.ColorPicker != "" {
		p.ColorPicker = strings.TrimSpace(in.ColorPicker)
	}
	if in.BgPic != "" {
		p.BgPic = strings.TrimSpace(in.BgPic)
	}
}

func resolveDoc(in SaveInput) (PageDoc, error) {
	if in.Doc != nil {
		doc := *in.Doc
		if doc.Page == nil {
			doc.Page = map[string]any{}
		}
		if doc.Items == nil {
			doc.Items = []map[string]any{}
		}
		return doc, nil
	}
	if len(in.Value) > 0 {
		var doc PageDoc
		if err := json.Unmarshal(in.Value, &doc); err == nil && doc.Items != nil {
			if doc.Page == nil {
				doc.Page = map[string]any{}
			}
			return doc, nil
		}
		var legacy LegacyValue
		if err := json.Unmarshal(in.Value, &legacy); err == nil {
			return legacyToDoc(legacy, in.Name, in.Title), nil
		}
		return PageDoc{}, ErrBadParam
	}
	return emptyDoc(), nil
}

func marshalDoc(doc PageDoc) (string, error) {
	if doc.Page == nil {
		doc.Page = map[string]any{}
	}
	if doc.Items == nil {
		doc.Items = []map[string]any{}
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func cloneMap(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return map[string]any{}
	}
	out := map[string]any{}
	_ = json.Unmarshal(b, &out)
	return out
}
