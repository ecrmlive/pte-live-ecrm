package catalog

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Store interface {
	ListPlatformCategories(ctx context.Context) ([]Category, error)
	GetCategory(ctx context.Context, id uint) (*Category, error)
	CreateCategory(ctx context.Context, c *Category) error
	UpdateCategory(ctx context.Context, c *Category) error
	DeleteCategory(ctx context.Context, id uint) error

	ListBrands(ctx context.Context) ([]Brand, error)
	CreateBrand(ctx context.Context, b *Brand) error
	UpdateBrand(ctx context.Context, b *Brand) error
	DeleteBrand(ctx context.Context, id uint) error

	ListProducts(ctx context.Context, status *int8, keyword string, merID *uint, page, limit int) ([]Product, int64, error)
	GetProduct(ctx context.Context, id uint) (*Product, error)
	UpdateProductAudit(ctx context.Context, id uint, status int8, refusal string) error
	MerchantName(ctx context.Context, merID uint) (string, error)
	CategoryName(ctx context.Context, cateID int) (string, error)
	MerchantNeedAudit(ctx context.Context, merID uint) (bool, error)

	CreateProduct(ctx context.Context, p *Product) error
	UpdateProduct(ctx context.Context, p *Product) error
	SoftDeleteProduct(ctx context.Context, id, merID uint) error
	UpsertDefaultSKU(ctx context.Context, v *AttrValue) error
	ListDefaultSKUs(ctx context.Context, productIDs []uint) (map[uint]AttrValue, error)

	ListOnSaleProducts(ctx context.Context, cateID int, keyword string, merID *uint, page, limit int) ([]Product, int64, error)
	ListPointsProducts(ctx context.Context, page, limit int) ([]Product, int64, error)
	ListCategoryIDsUnder(ctx context.Context, cateID int) ([]int, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

func (s *Service) CategoryTree(ctx context.Context) ([]CategoryNode, error) {
	rows, err := s.store.ListPlatformCategories(ctx)
	if err != nil {
		return nil, err
	}
	return buildCategoryTree(rows), nil
}

func (s *Service) CreateCategory(ctx context.Context, pid uint, name string, sort int, isShow int8) (*Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrNameRequired
	}
	level := uint(0)
	path := "/"
	if pid > 0 {
		parent, err := s.store.GetCategory(ctx, pid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		level = parent.Level + 1
		path = parent.Path
	}
	c := &Category{
		PID:      pid,
		CateName: name,
		Sort:     sort,
		IsShow:   isShow,
		Level:    level,
		MerID:    0,
		Type:     0,
		Path:     path,
	}
	if err := s.store.CreateCategory(ctx, c); err != nil {
		return nil, err
	}
	c.Path = fmt.Sprintf("%s%d/", path, c.StoreCategoryID)
	_ = s.store.UpdateCategory(ctx, c)
	return c, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id uint, name string, sort int, isShow int8) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrNameRequired
	}
	c, err := s.store.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	c.CateName = name
	c.Sort = sort
	c.IsShow = isShow
	return s.store.UpdateCategory(ctx, c)
}

func (s *Service) DeleteCategory(ctx context.Context, id uint) error {
	return s.store.DeleteCategory(ctx, id)
}

func (s *Service) ListBrands(ctx context.Context) ([]Brand, error) {
	return s.store.ListBrands(ctx)
}

func (s *Service) CreateBrand(ctx context.Context, name string, sort int, isShow int8) (*Brand, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrNameRequired
	}
	b := &Brand{BrandName: name, Sort: sort, IsShow: isShow}
	if err := s.store.CreateBrand(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) UpdateBrand(ctx context.Context, id uint, name string, sort int, isShow int8) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrNameRequired
	}
	return s.store.UpdateBrand(ctx, &Brand{
		BrandID:   id,
		BrandName: name,
		Sort:      sort,
		IsShow:    isShow,
	})
}

func (s *Service) DeleteBrand(ctx context.Context, id uint) error {
	return s.store.DeleteBrand(ctx, id)
}

func (s *Service) ListProducts(ctx context.Context, status *int8, keyword string, merID *uint, page, limit int) (*PageResult[Product], error) {
	list, total, err := s.store.ListProducts(ctx, status, keyword, merID, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if name, err := s.store.MerchantName(ctx, list[i].MerID); err == nil {
			list[i].MerName = name
		}
		if name, err := s.store.CategoryName(ctx, list[i].CateID); err == nil {
			list[i].CateName = name
		}
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &PageResult[Product]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetProduct(ctx context.Context, id uint) (*Product, error) {
	p, err := s.store.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name, err := s.store.MerchantName(ctx, p.MerID); err == nil {
		p.MerName = name
	}
	if name, err := s.store.CategoryName(ctx, p.CateID); err == nil {
		p.CateName = name
	}
	return p, nil
}

func (s *Service) AuditProduct(ctx context.Context, id uint, status int8, refusal string) error {
	if status != ProductStatusApproved && status != ProductStatusRejected && status != ProductStatusOff {
		return ErrBadStatus
	}
	if status == ProductStatusRejected && strings.TrimSpace(refusal) == "" {
		return ErrRejectNeedMsg
	}
	if _, err := s.GetProduct(ctx, id); err != nil {
		return err
	}
	if status != ProductStatusRejected {
		refusal = ""
	}
	return s.store.UpdateProductAudit(ctx, id, status, strings.TrimSpace(refusal))
}

func (s *Service) ListMerchantProducts(ctx context.Context, merID uint, status *int8, keyword string, page, limit int) (*PageResult[Product], error) {
	return s.ListProducts(ctx, status, keyword, &merID, page, limit)
}

func (s *Service) GetMerchantProduct(ctx context.Context, merID, id uint) (*Product, error) {
	p, err := s.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	if p.MerID != merID {
		return nil, ErrForbidden
	}
	return p, nil
}

func (s *Service) CreateMerchantProduct(ctx context.Context, merID uint, in ProductSaveInput) (*Product, error) {
	if merID == 0 {
		return nil, ErrForbidden
	}
	if err := validateSave(in); err != nil {
		return nil, err
	}
	needAudit, err := s.store.MerchantNeedAudit(ctx, merID)
	if err != nil {
		return nil, err
	}
	status := ProductStatusApproved
	if needAudit {
		status = ProductStatusPending
	}
	isShow := uint8(0)
	if in.IsShow != nil {
		isShow = *in.IsShow
	} else if status == ProductStatusApproved {
		isShow = 1
	}
	unit := strings.TrimSpace(in.UnitName)
	if unit == "" {
		unit = "件"
	}
	delivery := strings.TrimSpace(in.DeliveryWay)
	if delivery == "" {
		delivery = "2"
	}
	p := &Product{
		MerID:         merID,
		StoreName:     strings.TrimSpace(in.StoreName),
		StoreInfo:     strings.TrimSpace(in.StoreInfo),
		Keyword:       strings.TrimSpace(in.Keyword),
		BrandID:       in.BrandID,
		IsShow:        isShow,
		Status:        status,
		MerStatus:     1,
		CateID:        in.CateID,
		UnitName:      unit,
		Price:         in.Price,
		OtPrice:       in.OtPrice,
		Stock:         in.Stock,
		ProductType:   0,
		SpecType:      in.SpecType,
		Image:         in.Image,
		SliderImage:   in.SliderImage,
		DeliveryWay:   delivery,
		Type:          in.Type,
		MerSvipStatus: 1,
	}
	if in.SvipPriceType != nil {
		p.SvipPriceType = *in.SvipPriceType
	}
	if in.SvipPrice != nil {
		p.SvipPrice = *in.SvipPrice
	}
	if in.MerSvipStatus != nil {
		p.MerSvipStatus = *in.MerSvipStatus
	}
	if err := s.store.CreateProduct(ctx, p); err != nil {
		return nil, err
	}
	if err := s.upsertDefaultSKU(ctx, p, in.Cost); err != nil {
		return nil, err
	}
	return s.GetMerchantProduct(ctx, merID, p.ProductID)
}

func (s *Service) UpdateMerchantProduct(ctx context.Context, merID, id uint, in ProductSaveInput) (*Product, error) {
	if err := validateSave(in); err != nil {
		return nil, err
	}
	p, err := s.GetMerchantProduct(ctx, merID, id)
	if err != nil {
		return nil, err
	}
	needAudit, err := s.store.MerchantNeedAudit(ctx, merID)
	if err != nil {
		return nil, err
	}
	p.StoreName = strings.TrimSpace(in.StoreName)
	p.StoreInfo = strings.TrimSpace(in.StoreInfo)
	p.Keyword = strings.TrimSpace(in.Keyword)
	p.BrandID = in.BrandID
	p.CateID = in.CateID
	p.UnitName = strings.TrimSpace(in.UnitName)
	if p.UnitName == "" {
		p.UnitName = "件"
	}
	p.Price = in.Price
	p.OtPrice = in.OtPrice
	p.Stock = in.Stock
	p.Image = in.Image
	p.SliderImage = in.SliderImage
	p.Type = in.Type
	p.SpecType = in.SpecType
	if in.DeliveryWay != "" {
		p.DeliveryWay = in.DeliveryWay
	}
	if in.IsShow != nil {
		p.IsShow = *in.IsShow
	}
	if in.SvipPriceType != nil {
		p.SvipPriceType = *in.SvipPriceType
	}
	if in.SvipPrice != nil {
		p.SvipPrice = *in.SvipPrice
	}
	if in.MerSvipStatus != nil {
		p.MerSvipStatus = *in.MerSvipStatus
	}
	if needAudit {
		p.Status = ProductStatusPending
		p.Refusal = ""
	}
	if err := s.store.UpdateProduct(ctx, p); err != nil {
		return nil, err
	}
	if err := s.upsertDefaultSKU(ctx, p, in.Cost); err != nil {
		return nil, err
	}
	return s.GetMerchantProduct(ctx, merID, id)
}

func (s *Service) DeleteMerchantProduct(ctx context.Context, merID, id uint) error {
	if _, err := s.GetMerchantProduct(ctx, merID, id); err != nil {
		return err
	}
	return s.store.SoftDeleteProduct(ctx, id, merID)
}

func (s *Service) SetMerchantProductShow(ctx context.Context, merID, id uint, show bool) error {
	p, err := s.GetMerchantProduct(ctx, merID, id)
	if err != nil {
		return err
	}
	if show {
		if p.Status != ProductStatusApproved {
			return ErrNotOnSale
		}
		p.IsShow = 1
	} else {
		p.IsShow = 0
	}
	return s.store.UpdateProduct(ctx, p)
}

func (s *Service) SetMerchantProductStock(ctx context.Context, merID, id uint, stock uint) error {
	p, err := s.GetMerchantProduct(ctx, merID, id)
	if err != nil {
		return err
	}
	p.Stock = stock
	if err := s.store.UpdateProduct(ctx, p); err != nil {
		return err
	}
	// -1 表示保留既有默认 SKU 成本价，库存调整不能把成本价重置为 0。
	return s.upsertDefaultSKU(ctx, p, -1)
}

func (s *Service) ListAppCategories(ctx context.Context) ([]Category, error) {
	rows, err := s.store.ListPlatformCategories(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Category, 0, len(rows))
	for _, r := range rows {
		if r.IsShow == 1 {
			out = append(out, r)
		}
	}
	return out, nil
}

func (s *Service) ListAppProducts(ctx context.Context, cateID int, keyword string, merID *uint, page, limit int) (*PageResult[Product], error) {
	list, total, err := s.store.ListOnSaleProducts(ctx, cateID, keyword, merID, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if name, err := s.store.MerchantName(ctx, list[i].MerID); err == nil {
			list[i].MerName = name
		}
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &PageResult[Product]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListPointsProducts(ctx context.Context, page, limit int) (*PageResult[Product], error) {
	list, total, err := s.store.ListPointsProducts(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if name, err := s.store.MerchantName(ctx, list[i].MerID); err == nil {
			list[i].MerName = name
		}
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &PageResult[Product]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetAppProduct(ctx context.Context, id uint) (*Product, error) {
	p, err := s.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	if p.Status != ProductStatusApproved || p.IsShow != 1 || p.MerStatus != 1 {
		return nil, ErrNotOnSale
	}
	return p, nil
}

func (s *Service) AppHome(ctx context.Context) (banners []map[string]interface{}, hot []Product, err error) {
	page, err := s.ListAppProducts(ctx, 0, "", nil, 1, 8)
	if err != nil {
		return nil, nil, err
	}
	banners = []map[string]interface{}{
		{"id": 1, "title": "多商户入驻 · 精选好物", "image": ""},
		{"id": 2, "title": "同城配送 · 品质生活", "image": ""},
	}
	return banners, page.List, nil
}

func (s *Service) upsertDefaultSKU(ctx context.Context, p *Product, cost float64) error {
	if cost < 0 {
		existing, err := s.store.ListDefaultSKUs(ctx, []uint{p.ProductID})
		if err != nil {
			return err
		}
		if sku, ok := existing[p.ProductID]; ok {
			cost = sku.Cost
		} else {
			cost = 0
		}
	}
	sku := &AttrValue{
		ProductID: p.ProductID,
		Detail:    "{}",
		SKU:       "默认",
		Stock:     p.Stock,
		Cost:      cost,
		OtPrice:   p.OtPrice,
		Price:     p.Price,
		SvipPrice: p.SvipPrice,
		Unique:    fmt.Sprintf("sku%010d", p.ProductID),
		IsShow:    1,
	}
	return s.store.UpsertDefaultSKU(ctx, sku)
}

func validateSave(in ProductSaveInput) error {
	if strings.TrimSpace(in.StoreName) == "" {
		return ErrNameRequired
	}
	if in.CateID <= 0 {
		return ErrCateRequired
	}
	if in.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}

func buildCategoryTree(rows []Category) []CategoryNode {
	nodes := make(map[uint]*CategoryNode, len(rows))
	order := make([]uint, 0, len(rows))
	for _, row := range rows {
		nodes[row.StoreCategoryID] = &CategoryNode{
			StoreCategoryID: row.StoreCategoryID,
			PID:             row.PID,
			CateName:        row.CateName,
			Sort:            row.Sort,
			IsShow:          row.IsShow,
			Level:           row.Level,
			Children:        []CategoryNode{},
		}
		order = append(order, row.StoreCategoryID)
	}
	rootIDs := make([]uint, 0)
	for _, id := range order {
		n := nodes[id]
		if n.PID == 0 || nodes[n.PID] == nil {
			rootIDs = append(rootIDs, id)
			continue
		}
		parent := nodes[n.PID]
		parent.Children = append(parent.Children, *n)
	}
	out := make([]CategoryNode, 0, len(rootIDs))
	for _, id := range rootIDs {
		out = append(out, *nodes[id])
	}
	for i := range out {
		children := make([]CategoryNode, 0, len(out[i].Children))
		for _, ch := range out[i].Children {
			if full, ok := nodes[ch.StoreCategoryID]; ok {
				children = append(children, *full)
			} else {
				children = append(children, ch)
			}
		}
		out[i].Children = children
	}
	return out
}
